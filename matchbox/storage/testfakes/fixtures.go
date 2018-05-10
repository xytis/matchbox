package testfakes

import (
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

var (
	// Group is a machine group for testing.
	Group = &storagepb.Group{
		Id:       "test-group",
		Name:     "test group",
		Profile:  "test-profile",
		Selector: map[string]string{"uuid": "a1b2c3d4"},
		Metadata: []byte(`{"pod_network":"10.2.0.0/16","service_name":"etcd2"}`),
	}

	// GroupNoMetadata is a Group without any metadata.
	GroupNoMetadata = &storagepb.Group{
		Id:       "test-group-no-metadata",
		Selector: map[string]string{"uuid": "a1b2c3d4"},
		Metadata: nil,
	}

	// Profile is a machine profile for testing.
	Profile = &storagepb.Profile{
		Id:       "test-profile",
		Name:     "Test Profile",
		Template: map[string]string{"test-ignition": "fake-template"},
		Metadata: []byte(`{"Args":["a=b","c"],"Initrd":["/image/initrd_a", "/image/initrd_b"],"Kernel":"/image/kernel"}`),
	}

	// ProfileNoConfig is a Profile without extra config
	ProfileNoConfig = &storagepb.Profile{
		Id: "test-profile-no-config",
	}

	// Template is ignition template for testing
	Template = &storagepb.Template{
		Id:   "fake-template",
		Name: "Fake Ignition Template",
		Contents: []byte(`ignition_version: 1
systemd:
  units:
    - name: etcd2.service
      enable: true
`),
	}
)
