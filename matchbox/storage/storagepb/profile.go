package storagepb

import (
	"encoding/json"
)

// ParseProfile parses bytes into a Profile.
func ParseProfile(data []byte) (*Profile, error) {
	profile := new(Profile)
	err := json.Unmarshal(data, profile)
	return profile, err
}

// AssertValid validates a Profile. Returns nil if there are no validation
// errors.
func (p *Profile) AssertValid() error {
	// Id is required
	if p.Id == "" {
		return ErrIdRequired
	}
	return nil
}

func (p *Profile) Copy() *Profile {
	templates := make(map[string]string)
	for k, v := range p.Template {
		templates[k] = v
	}
	return &Profile{
		Id:       p.Id,
		Name:     p.Name,
		Template: templates,
		Metadata: p.Metadata,
	}
}

// ToRichProfile converts a Profile into a RichProfile suitable for writing and
// user manipulation.
func (g *Profile) ToRichProfile() (*RichProfile, error) {
	metadata := make(map[string]interface{})
	if g.Metadata != nil {
		err := json.Unmarshal(g.Metadata, &metadata)
		if err != nil {
			return nil, err
		}
	}
	//TODO: Not a full copy?!
	return &RichProfile{
		Id:       g.Id,
		Name:     g.Name,
		Template: g.Template,
		Metadata: metadata,
	}, nil
}

// RichProfile is parsed representation of stored Profile
type RichProfile struct {
	// machine readable Id
	Id string `json:"id,omitempty"`
	// Human readable name
	Name string `json:"name,omitempty"`
	// Template bindings
	Template map[string]string `json:"selector,omitempty"`
	// Metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ToProfile converts a user provided RichProfile into a Profile which can be
// serialized as a protocol buffer.
func (rg *RichProfile) ToProfile() (*Profile, error) {
	var metadata []byte
	if rg.Metadata != nil {
		var err error
		metadata, err = json.Marshal(rg.Metadata)
		if err != nil {
			return nil, err
		}
	}
	//TODO: Not a full copy?!
	return &Profile{
		Id:       rg.Id,
		Name:     rg.Name,
		Template: rg.Template,
		Metadata: metadata,
	}, nil
}
