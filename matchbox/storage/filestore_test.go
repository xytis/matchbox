package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coreos/matchbox/matchbox/storage/storagepb"
	fake "github.com/coreos/matchbox/matchbox/storage/testfakes"
)

func TestFileStoreGroupCRUD(t *testing.T) {
	dir, err := setup(&fake.FixedStore{})
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	store := NewFileStore(&FileStoreConfig{Root: dir})
	// assert that:
	// - Group creation was successful
	// - Group can be retrieved by id
	// - Group can be deleted by id
	err = store.GroupPut(fake.Group)
	assert.Nil(t, err)

	group, err := store.GroupGet(fake.Group.Id)
	assert.Nil(t, err)
	assert.Equal(t, fake.Group, group)

	err = store.GroupDelete(fake.Group.Id)
	assert.Nil(t, err)
	_, err = store.GroupGet(fake.Group.Id)
	if assert.Error(t, err) {
		assert.IsType(t, err, &os.PathError{})
	}
}

func TestFileStoreGroupGet(t *testing.T) {
	dir, err := setup(&fake.FixedStore{
		Groups: map[string]*storagepb.Group{
			fake.Group.Id:           fake.Group,
			fake.GroupNoMetadata.Id: fake.GroupNoMetadata,
		},
	})
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	store := NewFileStore(&FileStoreConfig{Root: dir})
	// assert that:
	// - Groups written to the store can be retrieved
	group, err := store.GroupGet(fake.Group.Id)
	assert.Nil(t, err)
	assert.Equal(t, fake.Group, group)
	group, err = store.GroupGet(fake.GroupNoMetadata.Id)
	assert.Nil(t, err)
	assert.Equal(t, fake.GroupNoMetadata, group)
}

func TestFileStoreGroupGet_NoGroup(t *testing.T) {
	dir, err := setup(&fake.FixedStore{})
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	store := NewFileStore(&FileStoreConfig{Root: dir})
	_, err = store.GroupGet("no-such-group")
	if assert.Error(t, err) {
		assert.IsType(t, &os.PathError{}, err)
	}
}

func TestFileStoreGroupList(t *testing.T) {
	dir, err := setup(&fake.FixedStore{
		Groups: map[string]*storagepb.Group{
			fake.Group.Id:           fake.Group,
			fake.GroupNoMetadata.Id: fake.GroupNoMetadata,
		},
	})
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	store := NewFileStore(&FileStoreConfig{Root: dir})
	groups, err := store.GroupList()
	assert.Nil(t, err)
	if assert.Equal(t, 2, len(groups)) {
		assert.Contains(t, groups, fake.Group)
		assert.Contains(t, groups, fake.GroupNoMetadata)
		assert.NotContains(t, groups, &storagepb.Group{})
	}
}

func TestFileStoreProfileCRUD(t *testing.T) {
	dir, err := setup(&fake.FixedStore{})
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	store := NewFileStore(&FileStoreConfig{Root: dir})
	// assert that:
	// - Profile creation was successful
	// - Profile can be retrieved by id
	// - Profile can be deleted by id
	err = store.ProfilePut(fake.Profile)
	assert.Nil(t, err)

	profile, err := store.ProfileGet(fake.Profile.Id)
	assert.Nil(t, err)
	assert.Equal(t, fake.Profile, profile)

	err = store.ProfileDelete(fake.Profile.Id)
	assert.Nil(t, err)
	_, err = store.ProfileGet(fake.Profile.Id)
	if assert.Error(t, err) {
		assert.IsType(t, err, &os.PathError{})
	}
}

func TestFileStoreProfileGet(t *testing.T) {
	dir, err := setup(&fake.FixedStore{
		Profiles: map[string]*storagepb.Profile{fake.Profile.Id: fake.Profile},
	})
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	store := NewFileStore(&FileStoreConfig{Root: dir})
	profile, err := store.ProfileGet(fake.Profile.Id)
	assert.Equal(t, fake.Profile, profile)
	assert.Nil(t, err)
	_, err = store.ProfileGet("no-such-profile")
	if assert.Error(t, err) {
		assert.IsType(t, &os.PathError{}, err)
	}
}

func TestFileStoreProfileList(t *testing.T) {
	dir, err := setup(&fake.FixedStore{
		Profiles: map[string]*storagepb.Profile{fake.Profile.Id: fake.Profile},
	})
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	store := NewFileStore(&FileStoreConfig{Root: dir})
	profiles, err := store.ProfileList()
	assert.Nil(t, err)
	if assert.Equal(t, 1, len(profiles)) {
		assert.Equal(t, fake.Profile, profiles[0])
	}
}

func TestFileStoreTemplateCRUD(t *testing.T) {
	dir, err := setup(&fake.FixedStore{})
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	store := NewFileStore(&FileStoreConfig{Root: dir})
	// assert that:
	// - Generic template creation was successful
	// - Generic template can be retrieved by name
	// - Generic template can be deleted by name
	err = store.TemplatePut(fake.Template)
	assert.Nil(t, err)

	template, err := store.TemplateGet(fake.Template.Id)
	assert.Nil(t, err)
	assert.Equal(t, fake.Template, template)
	assert.Equal(t, fake.Template.Id, template.Id)
	assert.Equal(t, fake.Template.Name, template.Name)
	assert.Equal(t, fake.Template.Contents, template.Contents)

	err = store.TemplateDelete(fake.Template.Id)
	assert.Nil(t, err)
	_, err = store.TemplateGet(fake.Template.Id)
	if assert.Error(t, err) {
		assert.IsType(t, err, &os.PathError{})
	}
}

// setup creates a temp fileStore directory to mirror a given fixedStore
// for testing. Returns the directory tree root. The caller must remove the
// temp directory when finished.
func setup(fixedStore *fake.FixedStore) (root string, err error) {
	root, err = ioutil.TempDir("", "data")
	if err != nil {
		return "", err
	}
	// directories
	profileDir := filepath.Join(root, "profiles")
	groupDir := filepath.Join(root, "groups")
	templateDir := filepath.Join(root, "templates")
	if err := mkdirs(profileDir, groupDir, templateDir); err != nil {
		return root, err
	}
	// files
	for _, profile := range fixedStore.Profiles {
		profileFile := filepath.Join(profileDir, profile.Id+".json")
		data, err := json.MarshalIndent(profile, "", "\t")
		if err != nil {
			return root, err
		}
		err = ioutil.WriteFile(profileFile, []byte(data), defaultFileMode)
		if err != nil {
			return root, err
		}
	}
	for _, group := range fixedStore.Groups {
		groupFile := filepath.Join(groupDir, group.Id+".json")
		richGroup, err := group.ToRichGroup()
		if err != nil {
			return root, err
		}
		data, err := json.MarshalIndent(richGroup, "", "\t")
		if err != nil {
			return root, err
		}
		err = ioutil.WriteFile(groupFile, []byte(data), defaultFileMode)
		if err != nil {
			return root, err
		}
	}
	for _, template := range fixedStore.Templates {
		templateFile := filepath.Join(templateDir, template.Id+".json")
		data, err := json.MarshalIndent(template, "", "\t")
		if err != nil {
			return root, err
		}
		err = ioutil.WriteFile(templateFile, []byte(data), defaultFileMode)
		if err != nil {
			return root, err
		}
	}
	return root, nil
}

// mkdirs creates new directories with the given names and default permission
// bits.
func mkdirs(names ...string) error {
	for _, dir := range names {
		if err := os.Mkdir(dir, defaultDirectoryMode); err != nil {
			return err
		}
	}
	return nil
}
