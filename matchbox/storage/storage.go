package storage

import (
	"errors"

	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

// Storage errors
var (
	ErrGroupNotFound    = errors.New("storage: No Group found")
	ErrProfileNotFound  = errors.New("storage: No Profile found")
	ErrTemplateNotFound = errors.New("storage: No Template found")
)

// A Store stores machine Groups, Profiles, and Configs.
type Store interface {
	// GroupPut creates or updates a Group.
	GroupPut(group *storagepb.Group) error
	// GroupGet returns a machine Group by id.
	GroupGet(id string) (*storagepb.Group, error)
	// GroupDelete deletes a machine Group by id.
	GroupDelete(id string) error
	// GroupList lists all machine Groups.
	GroupList() ([]*storagepb.Group, error)

	// ProfilePut creates or updates a Profile.
	ProfilePut(profile *storagepb.Profile) error
	// ProfileGet gets a profile by id.
	ProfileGet(id string) (*storagepb.Profile, error)
	// ProfileDelete deletes a profile by id.
	ProfileDelete(id string) error
	// ProfileList lists all profiles.
	ProfileList() ([]*storagepb.Profile, error)

	// TemplatePut creates or updates a template.
	TemplatePut(template *storagepb.Template) error
	// TemplateGet gets a template by name.
	TemplateGet(id string) (*storagepb.Template, error)
	// TemplateDelete deletes a template by name.
	TemplateDelete(id string) error
}
