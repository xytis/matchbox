package storage

import (
	"context"
	"encoding/json"

	"github.com/Sirupsen/logrus"
	etcd "github.com/coreos/etcd/clientv3"
	namespace "github.com/coreos/etcd/clientv3/namespace"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

type EtcdStoreConfig struct {
	Config etcd.Config
	Prefix string
	Logger *logrus.Logger
}

type etcdStore struct {
	client etcd.Client
	logger *logrus.Logger
}

func NewEtcdStore(config *EtcdStoreConfig) (Store, error) {
	client, err := etcd.New(config.Config)
	if err != nil {
		return nil, err
	}
	client.KV = namespace.NewKV(client.KV, "coreos.matchbox.v1/"+config.Prefix)
	return &etcdStore{
		client: client,
		logger: config.Logger,
	}, nil
}

func (s *etcdStore) GroupPut(group *storagepb.Group) error {
	kapi := etcd.NewKV(s.client)
	richGroup, err := group.ToRichGroup()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(richGroup, "", "\t")
	if err != nil {
		return err
	}

	_, err = kapi.Set(context.Background(), "/groups/"+group.Id, string(data), nil)
	return err
}

func (s *etcdStore) GroupGet(id string) (*storagepb.Group, error) {
	kapi := etcd.NewKV(s.client)
	resp, err := kapi.Get(context.Background(), "/groups/"+id, nil)
	if err != nil {
		return nil, err
	}
	group, err := storagepb.ParseGroup([]byte(resp.Node.Value))
	if err != nil {
		return nil, err
	}
	return group, err
}

func (s *etcdStore) GroupDelete(id string) error {
	kapi := etcd.NewKV(s.client)
	_, err := kapi.Delete(context.Background(), "/groups/"+id, nil)
	return err
}

func (s *etcdStore) GroupList() ([]*storagepb.Group, error) {
	kapi := etcd.NewKV(s.client)
	resp, err := kapi.Get(context.Background(), "/groups/", nil)
	if err != nil {
		return nil, err
	}

	groups := make([]*storagepb.Group, 0, len(resp.Node.Nodes))
	for _, n := range resp.Node.Nodes {
		group, err := storagepb.ParseGroup([]byte(n.Value))
		if err == nil {
			groups = append(groups, group)
		} else if s.logger != nil {
			s.logger.Infof("Group %q: %v", n.Key, err)
		}
	}
	return groups, nil
}

func (s *etcdStore) ProfilePut(profile *storagepb.Profile) error {
	kapi := etcd.NewKV(s.client)
	data, err := json.MarshalIndent(profile, "", "\t")
	if err != nil {
		return err
	}
	_, err = kapi.Set(context.Background(), "/profiles/"+profile.Id, string(data), nil)
	return err
}

func (s *etcdStore) ProfileGet(id string) (*storagepb.Profile, error) {
	kapi := etcd.NewKV(s.client)
	resp, err := kapi.Get(context.Background(), "/profiles/"+id, nil)
	if err != nil {
		return nil, err
	}
	profile := new(storagepb.Profile)
	err = json.Unmarshal([]byte(resp.Node.Value), profile)
	if err != nil {
		return nil, err
	}
	if err := profile.AssertValid(); err != nil {
		return nil, err
	}
	return profile, err
}

func (s *etcdStore) ProfileDelete(id string) error {
	kapi := etcd.NewKV(s.client)
	_, err := kapi.Delete(context.Background(), "/profiles/"+id, nil)
	return err
}

func (s *etcdStore) ProfileList() ([]*storagepb.Profile, error) {
	kapi := etcd.NewKV(s.client)
	resp, err := kapi.Get(context.Background(), "/profiles/", nil)
	if err != nil {
		return nil, err
	}

	profiles := make([]*storagepb.Profile, 0, len(resp.Node.Nodes))
	for _, n := range resp.Node.Nodes {
		profile := new(storagepb.Profile)
		err = json.Unmarshal([]byte(n.Value), profile)
		if err == nil {
			profiles = append(profiles, profile)
		} else if s.logger != nil {
			s.logger.Infof("Profile %q: %v", n.Key, err)
		}
	}
	return profiles, nil
}

func (s *etcdStore) IgnitionPut(name string, config []byte) error {
	kapi := etcd.NewKV(s.client)
	_, err := kapi.Set(context.Background(), "/ignition/"+name, string(config), nil)
	return err
}

func (s *etcdStore) IgnitionGet(name string) (string, error) {
	kapi := etcd.NewKV(s.client)
	resp, err := kapi.Get(context.Background(), "/ignition/"+name, nil)
	return resp.Node.Value, err
}

func (s *etcdStore) IgnitionDelete(name string) error {
	kapi := etcd.NewKV(s.client)
	_, err := kapi.Delete(context.Background(), "/ignition/"+name, nil)
	return err
}

func (s *etcdStore) GenericPut(name string, config []byte) error {
	kapi := etcd.NewKV(s.client)
	_, err := kapi.Set(context.Background(), "/generic/"+name, string(config), nil)
	return err
}

func (s *etcdStore) GenericGet(name string) (string, error) {
	kapi := etcd.NewKV(s.client)
	resp, err := kapi.Get(context.Background(), "/generic/"+name, nil)
	return resp.Node.Value, err
}

func (s *etcdStore) GenericDelete(name string) error {
	kapi := etcd.NewKV(s.client)
	_, err := kapi.Delete(context.Background(), "/generic/"+name, nil)
	return err
}

func (s *etcdStore) CloudGet(name string) (string, error) {
	kapi := etcd.NewKV(s.client)
	resp, err := kapi.Get(context.Background(), "/generic/"+name, nil)
	return resp.Node.Value, err
}
