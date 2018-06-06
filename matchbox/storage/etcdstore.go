package storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/coreos/matchbox/matchbox/storage/config"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"

	etcd "github.com/coreos/etcd/clientv3"
	namespace "github.com/coreos/etcd/clientv3/namespace"
	"go.uber.org/zap"
)

const (
	etcdStoreNamespace = "coreos.matchbox.v1/"
	groupScope         = "/groups/"
	profileScope       = "/profiles/"
	templateScope      = "/templates/"
)

// etcdStore implements ths Store interface.
type etcdStore struct {
	client *etcd.Client
	logger *zap.SugaredLogger
}

// NewEtcdStore returns a new etcd-backed Store.
func NewEtcdStore(config *config.EtcdStoreConfig, logger *zap.Logger) (Store, error) {
	cfg := etcd.Config{
		Endpoints:   config.Endpoints,
		Password:    config.Password,
		Username:    config.Username,
		DialTimeout: 5 * time.Second,
	}
	client, err := etcd.New(cfg)
	if err != nil {
		return nil, err
	}
	client.KV = namespace.NewKV(client.KV, etcdStoreNamespace+config.Prefix)
	client.Watcher = namespace.NewWatcher(client.Watcher, etcdStoreNamespace+config.Prefix)
	client.Lease = namespace.NewLease(client.Lease, etcdStoreNamespace+config.Prefix)
	return &etcdStore{
		client: client,
		logger: logger.Sugar(),
	}, nil
}

// GroupPut writes the given Group.
func (s *etcdStore) GroupPut(group *storagepb.Group) error {
	richGroup, err := group.ToRichGroup()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(richGroup, "", "\t")
	if err != nil {
		return err
	}

	_, err = s.client.Put(context.Background(), groupScope+group.Id, string(data))
	return err
}

// GroupGet returns a machine Group by id.
func (s *etcdStore) GroupGet(id string) (*storagepb.Group, error) {
	resp, err := s.client.Get(context.Background(), groupScope+id)
	if err != nil {
		return nil, err
	}
	if resp.Count == 0 {
		return nil, ErrGroupNotFound
	}
	kv := resp.Kvs[0]
	group, err := storagepb.ParseGroup([]byte(kv.Value))
	if err != nil {
		return nil, err
	}
	return group, err
}

// GroupDelete deletes a machine Group by id.
func (s *etcdStore) GroupDelete(id string) error {
	_, err := s.client.Delete(context.Background(), groupScope+id)
	return err
}

// GroupList lists all machine Groups.
func (s *etcdStore) GroupList() ([]*storagepb.Group, error) {
	resp, err := s.client.Get(context.Background(), groupScope, etcd.WithPrefix())
	if err != nil {
		return nil, err
	}

	groups := make([]*storagepb.Group, 0, resp.Count)
	for _, kv := range resp.Kvs {
		group, err := storagepb.ParseGroup(kv.Value)
		if err == nil {
			groups = append(groups, group)
		} else if s.logger != nil {
			s.logger.Infof("Group %q: %v", kv.Key, err)
		}
	}
	return groups, nil
}

// ProfilePut writes the given Profile.
func (s *etcdStore) ProfilePut(profile *storagepb.Profile) error {
	data, err := json.MarshalIndent(profile, "", "\t")
	if err != nil {
		return err
	}
	_, err = s.client.Put(context.Background(), profileScope+profile.Id, string(data))
	return err
}

// ProfileGet gets a profile by id.
func (s *etcdStore) ProfileGet(id string) (*storagepb.Profile, error) {
	resp, err := s.client.Get(context.Background(), profileScope+id)
	if err != nil {
		return nil, err
	}
	if resp.Count == 0 {
		return nil, ErrProfileNotFound
	}
	kv := resp.Kvs[0]
	profile := new(storagepb.Profile)
	err = json.Unmarshal([]byte(kv.Value), profile)
	if err != nil {
		return nil, err
	}
	if err := profile.AssertValid(); err != nil {
		return nil, err
	}
	return profile, err
}

// ProfileDelete deletes a profile by id.
func (s *etcdStore) ProfileDelete(id string) error {
	_, err := s.client.Delete(context.Background(), profileScope+id)
	return err
}

// ProfileList lists all profiles.
func (s *etcdStore) ProfileList() ([]*storagepb.Profile, error) {
	resp, err := s.client.Get(context.Background(), profileScope, etcd.WithPrefix())
	if err != nil {
		return nil, err
	}
	profiles := make([]*storagepb.Profile, 0, resp.Count)
	for _, kv := range resp.Kvs {
		profile := new(storagepb.Profile)
		err = json.Unmarshal([]byte(kv.Value), profile)
		if err == nil {
			profiles = append(profiles, profile)
		} else if s.logger != nil {
			s.logger.Infof("Profile %q: %v", kv.Key, err)
		}
	}
	return profiles, nil
}

// TemplatePut creates or updates a template.
func (s *etcdStore) TemplatePut(template *storagepb.Template) error {
	data, err := json.MarshalIndent(template, "", "\t")
	if err != nil {
		return err
	}
	_, err = s.client.Put(context.Background(), templateScope+template.Id, string(data))
	return err
}

// TemplateGet gets a template by name.
func (s *etcdStore) TemplateGet(id string) (*storagepb.Template, error) {
	resp, err := s.client.Get(context.Background(), templateScope+id)
	if err != nil {
		return nil, err
	}
	if resp.Count == 0 {
		return nil, ErrTemplateNotFound
	}
	kv := resp.Kvs[0]
	template := new(storagepb.Template)
	err = json.Unmarshal([]byte(kv.Value), template)
	if err != nil {
		return nil, err
	}
	if err := template.AssertValid(); err != nil {
		return nil, err
	}
	return template, err
}

// TemplateDelete deletes a template by name.
func (s *etcdStore) TemplateDelete(id string) error {
	_, err := s.client.Delete(context.Background(), templateScope+id)
	return err
}

// TemplateList lists all profiles.
func (s *etcdStore) TemplateList() ([]*storagepb.Template, error) {
	resp, err := s.client.Get(context.Background(), templateScope, etcd.WithPrefix())
	if err != nil {
		return nil, err
	}
	templates := make([]*storagepb.Template, 0, resp.Count)
	for _, kv := range resp.Kvs {
		template := new(storagepb.Template)
		err = json.Unmarshal([]byte(kv.Value), template)
		if err == nil {
			templates = append(templates, template)
		} else if s.logger != nil {
			s.logger.Infof("Template %q: %v", kv.Key, err)
		}
	}
	return templates, nil
}
