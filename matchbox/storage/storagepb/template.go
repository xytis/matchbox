package storagepb

import (
	"encoding/json"
)

// ParseTemplate parses bytes into a Template.
func ParseTemplate(data []byte) (*Template, error) {
	template := new(Template)
	err := json.Unmarshal(data, template)
	return template, err
}

// AssertValid validates a Template. Returns nil if there are no validation
// errors.
func (t *Template) AssertValid() error {
	// Id is required
	if t.Id == "" {
		return ErrIdRequired
	}
	return nil
}

func (t *Template) Copy() *Template {
	contents := make([]byte, len(t.Contents))
	copy(t.Contents, contents)
	return &Template{
		Id:       t.Id,
		Name:     t.Name,
		Contents: contents,
	}
}
