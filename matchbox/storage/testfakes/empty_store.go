package testfakes

import (
	"fmt"

	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

// EmptyStore is used for testing purposes.
type EmptyStore struct{}

// GroupPut returns an error writing any Group.
func (s *EmptyStore) GroupPut(group *storagepb.Group) error {
	return fmt.Errorf("emptyStore does not accept Groups")
}

// GroupGet returns a group not found error.
func (s *EmptyStore) GroupGet(id string) (*storagepb.Group, error) {
	return nil, errGroupNotFound
}

// GroupDelete returns a nil error (successful deletion).
func (s *EmptyStore) GroupDelete(id string) error {
	return nil
}

// GroupList returns an empty list of groups.
func (s *EmptyStore) GroupList() (groups []*storagepb.Group, err error) {
	return groups, nil
}

// ProfilePut returns an error writing any Profile.
func (s *EmptyStore) ProfilePut(profile *storagepb.Profile) error {
	return fmt.Errorf("emptyStore does not accept Profiles")
}

// ProfileGet returns a profile not found error.
func (s *EmptyStore) ProfileGet(id string) (*storagepb.Profile, error) {
	return nil, errProfileNotFound
}

// ProfileDelete returns a nil error (successful deletion).
func (s *EmptyStore) ProfileDelete(id string) error {
	return nil
}

// ProfileList returns an empty list of profiles.
func (s *EmptyStore) ProfileList() (profiles []*storagepb.Profile, err error) {
	return profiles, nil
}

// TemplatePut creates or updates a template.
func (s *EmptyStore) TemplatePut(*storagepb.Template) error {
	return fmt.Errorf("emptyStore does not accept templates")
}

// TemplateGet gets a template by name.
func (s *EmptyStore) TemplateGet(id string) (*storagepb.Template, error) {
	return nil, errTemplateNotFound
}

// TemplateDelete deletes a template by name.
func (s *EmptyStore) TemplateDelete(id string) error {
	return nil
}

// TemplateList returns an empty list of profiles.
func (s *EmptyStore) TemplateList() (templates []*storagepb.Template, err error) {
	return templates, nil
}
