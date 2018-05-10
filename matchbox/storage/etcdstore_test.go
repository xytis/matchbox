package storage

import (
	"testing"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/integration"
	"github.com/stretchr/testify/assert"

	"github.com/coreos/matchbox/matchbox/storage/storagepb"
	fake "github.com/coreos/matchbox/matchbox/storage/testfakes"
)

func TestEtcdClientCreation(t *testing.T) {
	cluster := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	defer cluster.Terminate(t)

	c := cluster.RandClient()

	store, err := NewEtcdStore(&EtcdStoreConfig{
		Config: etcd.Config{
			Endpoints: c.Endpoints(),
		},
		Prefix: "test",
	})
	assert.Nil(t, err)

	err = store.TemplatePut(fake.Template)
	assert.Nil(t, err)

	template, err := store.TemplateGet(fake.Template.Id)
	assert.Equal(t, fake.Template, template)
	assert.Equal(t, fake.Template.Id, template.Id)
	assert.Equal(t, fake.Template.Name, template.Name)
	assert.Equal(t, fake.Template.Contents, template.Contents)
}

func TestEtcdStoreGroupCRUD(t *testing.T) {
	cluster := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	defer cluster.Terminate(t)

	store := &etcdStore{
		client: cluster.RandClient(),
	}

	var err error
	// assert that:
	// - Group creation is successful
	// - Multiple groups can be stored
	// - Group can be retrieved by id
	// - Group list can be retrieved
	// - Group can be deleted by id
	// - Non existing group query returns error
	err = store.GroupPut(fake.Group)
	assert.Nil(t, err)

	err = store.GroupPut(fake.GroupNoMetadata)
	assert.Nil(t, err)

	group, err := store.GroupGet(fake.Group.Id)
	assert.Nil(t, err)
	assert.Equal(t, fake.Group, group)

	groups, err := store.GroupList()
	assert.Nil(t, err)
	if assert.Equal(t, 2, len(groups)) {
		assert.Contains(t, groups, fake.Group)
		assert.Contains(t, groups, fake.GroupNoMetadata)
		assert.NotContains(t, groups, &storagepb.Group{})
	}

	err = store.GroupDelete(fake.Group.Id)
	assert.Nil(t, err)

	_, err = store.GroupGet(fake.Group.Id)
	if assert.Error(t, err) {
		assert.Equal(t, err, ErrGroupNotFound)
	}
}

func TestEtcdStoreProfileCRUD(t *testing.T) {
	cluster := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	defer cluster.Terminate(t)

	store := &etcdStore{
		client: cluster.RandClient(),
	}

	var err error
	// assert that:
	// - Group creation is successful
	// - Multiple groups can be stored
	// - Group can be retrieved by id
	// - Group list can be retrieved
	// - Group can be deleted by id
	// - Non existing group query returns error
	err = store.ProfilePut(fake.Profile)
	assert.Nil(t, err)

	err = store.ProfilePut(fake.ProfileNoConfig)
	assert.Nil(t, err)

	profile, err := store.ProfileGet(fake.Profile.Id)
	assert.Nil(t, err)
	assert.Equal(t, fake.Profile, profile)

	profiles, err := store.ProfileList()
	assert.Nil(t, err)
	if assert.Equal(t, 2, len(profiles)) {
		assert.Contains(t, profiles, fake.Profile)
		assert.Contains(t, profiles, fake.ProfileNoConfig)
		assert.NotContains(t, profiles, &storagepb.Group{})
	}

	err = store.ProfileDelete(fake.Profile.Id)
	assert.Nil(t, err)

	_, err = store.ProfileGet(fake.Profile.Id)
	if assert.Error(t, err) {
		assert.Equal(t, err, ErrProfileNotFound)
	}
}

func TestEtcdStoreTemplateCRUD(t *testing.T) {
	cluster := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	defer cluster.Terminate(t)

	store := &etcdStore{
		client: cluster.RandClient(),
	}

	var err error
	// assert that:
	// - Ignition template creation was successful
	// - Ignition template can be retrieved by name
	// - Ignition template can be deleted by name
	err = store.TemplatePut(fake.Template)
	assert.Nil(t, err)

	template, err := store.TemplateGet(fake.Template.Id)
	assert.Nil(t, err)
	assert.Equal(t, fake.Template.Contents, template.Contents)

	err = store.TemplateDelete(fake.Template.Id)
	assert.Nil(t, err)
	_, err = store.TemplateGet(fake.Template.Id)
	if assert.Error(t, err) {
		assert.Equal(t, err, ErrTemplateNotFound)
	}
}
