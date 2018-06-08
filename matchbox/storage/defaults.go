package storage

import (
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

const (
	ignition = `{
  "ignition": { "version": "2.2.0" },
  "passwd": {
    "users": [
      {
        "name": "core",
        "sshAuthorizedKeys": [
        {{- if .SSHAuthorizedKeys }}
          {{- range $i, $key := .SSHAuthorizedKeys }}
          "ssh-rsa {{$key}}"{{- if ne ($i+1 len(.SSHAuthorizedKeys)) }},{{ end }}
          {{- end }}
        {{- end }}
        ]
      }
    ]
  }
}`

	ipxe = `#!ipxe
kernel {{.Kernel}}{{range $arg := .Args}} {{$arg}}{{end}}
{{- range $element := .Initrd }}
initrd {{$element}}
{{- end}}
boot
`

	grub = `default=0
fallback=1
timeout=1
menuentry "CoreOS (EFI)" {
  echo "Loading kernel"
  linuxefi "{{.Kernel}}"{{range $arg := .Args}} {{$arg}}{{end}}
  echo "Loading initrd"
  initrdefi {{ range $element := .Initrd }} "{{$element}}"{{end}}
}
menuentry "CoreOS (BIOS)" {
  echo "Loading kernel"
  linux "{{.Kernel}}"{{range $arg := .Args}} {{$arg}}{{end}}
  echo "Loading initrd"
  initrd {{ range $element := .Initrd }} "{{$element}}"{{end}}
}
`
)

// AssertDefaultTemplates ensures that data store contains default templates
func AssertDefaultTemplates(s Store) {
	if _, err := s.TemplateGet("default-ipxe"); err != nil {
		s.TemplatePut(&storagepb.Template{
			Id:       "default-ipxe",
			Name:     "Default IPXE boot configuration",
			Contents: []byte(ipxe),
		})
	}
	if _, err := s.TemplateGet("default-ignition"); err != nil {
		s.TemplatePut(&storagepb.Template{
			Id:       "default-ignition",
			Name:     "Default Ignition configuration",
			Contents: []byte(ignition),
		})
	}
	if _, err := s.TemplateGet("default-grub"); err != nil {
		s.TemplatePut(&storagepb.Template{
			Id:       "default-grub",
			Name:     "Default GRUB template",
			Contents: []byte(grub),
		})
	}
}
