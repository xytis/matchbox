package testfakes

import (
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

// FixedStore is used for testing purposes.
type FixedStore struct {
	Groups    map[string]*storagepb.Group
	Profiles  map[string]*storagepb.Profile
	Templates map[string]*storagepb.Template
}

// NewFixedStore returns a new FixedStore.
func NewFixedStore() *FixedStore {
	return &FixedStore{
		Groups:    make(map[string]*storagepb.Group),
		Profiles:  make(map[string]*storagepb.Profile),
		Templates: make(map[string]*storagepb.Template),
	}
}

// GroupPut write the given Group the Groups map.
func (s *FixedStore) GroupPut(group *storagepb.Group) error {
	s.Groups[group.Id] = group
	return nil
}

// GroupGet returns the Group from the Groups map with the given id.
func (s *FixedStore) GroupGet(id string) (*storagepb.Group, error) {
	if group, present := s.Groups[id]; present {
		return group, nil
	}
	return nil, errGroupNotFound
}

// GroupDelete deletes the Group from the Groups map with the given id.
func (s *FixedStore) GroupDelete(id string) error {
	delete(s.Groups, id)
	return nil
}

// GroupList returns the groups in the Groups map.
func (s *FixedStore) GroupList() ([]*storagepb.Group, error) {
	groups := make([]*storagepb.Group, len(s.Groups))
	i := 0
	for _, g := range s.Groups {
		groups[i] = g
		i++
	}
	return groups, nil
}

// ProfilePut writes the given Profile to the Profiles map.
func (s *FixedStore) ProfilePut(profile *storagepb.Profile) error {
	s.Profiles[profile.Id] = profile
	return nil
}

// ProfileGet returns the Profile from the Profiles map with the given id.
func (s *FixedStore) ProfileGet(id string) (*storagepb.Profile, error) {
	if profile, present := s.Profiles[id]; present {
		return profile, nil
	}
	return nil, errProfileNotFound
}

// ProfileDelete deletes the Profile from the Profiles map with the given id.
func (s *FixedStore) ProfileDelete(id string) error {
	delete(s.Profiles, id)
	return nil
}

// ProfileList returns the profiles in the Profiles map.
func (s *FixedStore) ProfileList() ([]*storagepb.Profile, error) {
	profiles := make([]*storagepb.Profile, len(s.Profiles))
	i := 0
	for _, p := range s.Profiles {
		profiles[i] = p
		i++
	}
	return profiles, nil
}

// TemplatePut creates or updates a template.
func (s *FixedStore) TemplatePut(template *storagepb.Template) error {
	s.Templates[template.Id] = template
	return nil
}

// TemplateGet gets a template by name.
func (s *FixedStore) TemplateGet(id string) (*storagepb.Template, error) {
	if template, present := s.Templates[id]; present {
		return template, nil
	}
	return nil, errTemplateNotFound
}

// TemplateDelete deletes a template by name.
func (s *FixedStore) TemplateDelete(id string) error {
	delete(s.Templates, id)
	return nil
}

// TemplateList returns the profiles in the Templates map.
func (s *FixedStore) TemplateList() ([]*storagepb.Template, error) {
	templates := make([]*storagepb.Template, len(s.Templates))
	i := 0
	for _, p := range s.Templates {
		templates[i] = p
		i++
	}
	return templates, nil
}
