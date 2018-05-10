package testfakes

import (
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

// BrokenStore returns errors for testing purposes.
type BrokenStore struct{}

// GroupPut returns an error.
func (s *BrokenStore) GroupPut(group *storagepb.Group) error {
	return errIntentional
}

// GroupGet returns an error.
func (s *BrokenStore) GroupGet(id string) (*storagepb.Group, error) {
	return nil, errIntentional
}

// GroupDelete returns an error.
func (s *BrokenStore) GroupDelete(id string) error {
	return errIntentional
}

// GroupList returns an error.
func (s *BrokenStore) GroupList() (groups []*storagepb.Group, err error) {
	return groups, errIntentional
}

// ProfilePut returns an error.
func (s *BrokenStore) ProfilePut(profile *storagepb.Profile) error {
	return errIntentional
}

// ProfileGet returns an error.
func (s *BrokenStore) ProfileGet(id string) (*storagepb.Profile, error) {
	return nil, errIntentional
}

// ProfileDelete returns an error.
func (s *BrokenStore) ProfileDelete(id string) error {
	return errIntentional
}

// ProfileList returns an error.
func (s *BrokenStore) ProfileList() (profiles []*storagepb.Profile, err error) {
	return profiles, errIntentional
}

// TemplatePut returns an error.
func (s *BrokenStore) TemplatePut(*storagepb.Template) error {
	return errIntentional
}

// TemplateGet returns an error.
func (s *BrokenStore) TemplateGet(string) (*storagepb.Template, error) {
	return nil, errIntentional
}

// TemplateDelete returns an error.
func (s *BrokenStore) TemplateDelete(string) error {
	return errIntentional
}
