package testfakes

import (
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

// Labels which match fake groups.
func Labels() map[string]string {
	return map[string]string{"uuid": "a1b2c3d4"}
}

// LabelsEmpty contain no labels
func LabelsEmpty() map[string]string {
	return make(map[string]string)
}

// Group is a machine group for testing.
func Group() *storagepb.Group {
	return &storagepb.Group{
		Id:       "test-group",
		Name:     "test group",
		Profile:  "test-profile",
		Selector: map[string]string{"uuid": "a1b2c3d4"},
		Metadata: []byte(`{"greeting":"Hello there","pod_network":"10.2.0.0/16","service_name":"etcd2"}`),
	}
}

// GroupNoMetadata is a Group without any metadata.
func GroupNoMetadata() *storagepb.Group {
	return &storagepb.Group{
		Id:       "test-group-no-metadata",
		Profile:  "test-profile-no-metadata",
		Selector: map[string]string{"uuid": "a1b2c3d4"},
	}
}

// Profile is a machine profile for testing.
func Profile() *storagepb.Profile {
	return &storagepb.Profile{
		Id:   "test-profile",
		Name: "Test Profile",
		Template: map[string]string{
			"grub":     "fake-grub",
			"ipxe":     "fake-ipxe",
			"ignition": "fake-ignition",
			"custom":   "fake-custom",
		},
		Metadata: []byte(`{"args":["a=b","c"],"initrd":["/image/initrd_a", "/image/initrd_b"],"kernel":"/image/kernel"}`),
	}
}

// ProfileNoMetadata is a Profile without any metadata.
func ProfileNoMetadata() *storagepb.Profile {
	return &storagepb.Profile{
		Id: "test-profile-no-metadata",
	}
}

// GrubTemplate is bare grub template with values from profile and group
func GrubTemplate() *storagepb.Template {
	return &storagepb.Template{
		Id:   "fake-grub",
		Name: "Fake Grub Template",
		Contents: []byte(`default=0
fallback=1
timeout=1
menuentry "CoreOS (EFI)" {
  echo "Loading kernel"
  linuxefi "{{.kernel}}"{{range $arg := .args}} {{$arg}}{{end}}
  echo "Loading initrd"
  initrdefi {{range $element := .initrd}} "{{$element}}"{{end}}
}`),
	}
}

// IPXETemplate is bare grub template with values from profile and group
func IPXETemplate() *storagepb.Template {
	return &storagepb.Template{
		Id:   "fake-ipxe",
		Name: "Fake IPXE Template",
		Contents: []byte(`#!ipxe
kernel {{.kernel}}{{range $arg := .args}} {{$arg}}{{end}}{{range $element := .initrd}}
initrd {{$element}}{{end}}
boot
`),
	}
}

// IgnitionTemplate is example ignition template which should pass validation
func IgnitionTemplate() *storagepb.Template {
	return &storagepb.Template{
		Id:   "fake-ignition",
		Name: "Fake Ignition Template",
		Contents: []byte(`{
  "ignition": { "version": "2.2.0" },
  "systemd": {
    "units": [{
      "name": "etcd2.service",
      "enabled": true
    },
    {
      "name": "example.service",
      "enabled": true,
      "contents": "[Service]\nType=oneshot\nExecStart=/usr/bin/echo {{.greeting}}\n\n[Install]\nWantedBy=multi-user.target"
    }]
  }
}`),
	}
}

// CustomTemplate is an example custom template which has no preferred format
func CustomTemplate() *storagepb.Template {
	return &storagepb.Template{
		Id:       "fake-custom",
		Name:     "Fake custom Template",
		Contents: []byte(`UUID: {{.label.uuid}}`),
	}
}
