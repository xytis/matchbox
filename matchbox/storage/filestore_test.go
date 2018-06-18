package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/coreos/matchbox/matchbox/storage/config"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
	fake "github.com/coreos/matchbox/matchbox/storage/testfakes"
)

func TestFileStoreGroupCRUD(t *testing.T) {
	dir, err := setup(&fake.FixedStore{})
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	store, _ := NewFileStore(&config.FileStoreConfig{Root: dir}, zap.NewNop())

	// assert that:
	// - Group creation is successful
	// - Multiple groups can be stored
	// - Group can be retrieved by id
	// - Group list can be retrieved
	// - Group can be deleted by id
	// - Non existing group query returns error
	err = store.GroupPut(fake.Group())
	assert.Nil(t, err)

	err = store.GroupPut(fake.GroupNoMetadata())
	assert.Nil(t, err)

	group, err := store.GroupGet(fake.Group().Id)
	assert.Nil(t, err)
	assert.Equal(t, fake.Group(), group)

	groups, err := store.GroupList()
	assert.Nil(t, err)
	if assert.Equal(t, 2, len(groups)) {
		assert.Contains(t, groups, fake.Group())
		assert.Contains(t, groups, fake.GroupNoMetadata())
		assert.NotContains(t, groups, &storagepb.Group{})
	}

	err = store.GroupDelete(fake.Group().Id)
	assert.Nil(t, err)

	_, err = store.GroupGet(fake.Group().Id)
	if assert.Error(t, err) {
		assert.Equal(t, err, ErrGroupNotFound)
	}
}

func TestFileStoreProfileCRUD(t *testing.T) {
	dir, err := setup(&fake.FixedStore{})
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	store, _ := NewFileStore(&config.FileStoreConfig{Root: dir}, zap.NewNop())

	// assert that:
	// - Group creation is successful
	// - Multiple groups can be stored
	// - Group can be retrieved by id
	// - Group list can be retrieved
	// - Group can be deleted by id
	// - Non existing group query returns error
	err = store.ProfilePut(fake.Profile())
	assert.Nil(t, err)

	err = store.ProfilePut(fake.ProfileNoMetadata())
	assert.Nil(t, err)

	profile, err := store.ProfileGet(fake.Profile().Id)
	assert.Nil(t, err)
	assert.Equal(t, fake.Profile(), profile)

	profiles, err := store.ProfileList()
	assert.Nil(t, err)
	if assert.Equal(t, 2, len(profiles)) {
		assert.Contains(t, profiles, fake.Profile())
		assert.Contains(t, profiles, fake.ProfileNoMetadata())
		assert.NotContains(t, profiles, &storagepb.Group{})
	}

	err = store.ProfileDelete(fake.Profile().Id)
	assert.Nil(t, err)

	_, err = store.ProfileGet(fake.Profile().Id)
	if assert.Error(t, err) {
		assert.Equal(t, err, ErrProfileNotFound)
	}
}

func TestFileStoreTemplateCRUD(t *testing.T) {
	dir, err := setup(&fake.FixedStore{})
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	store, _ := NewFileStore(&config.FileStoreConfig{Root: dir}, zap.NewNop())

	// assert that:
	// - Ignition template creation was successful
	// - Ignition template can be retrieved by name
	// - Ignition template can be deleted by name
	err = store.TemplatePut(fake.CustomTemplate())
	assert.Nil(t, err)

	template, err := store.TemplateGet(fake.CustomTemplate().Id)
	assert.Nil(t, err)
	assert.Equal(t, fake.CustomTemplate().Contents, template.Contents)

	err = store.TemplateDelete(fake.CustomTemplate().Id)
	assert.Nil(t, err)
	_, err = store.TemplateGet(fake.CustomTemplate().Id)
	if assert.Error(t, err) {
		assert.Equal(t, err, ErrTemplateNotFound)
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
