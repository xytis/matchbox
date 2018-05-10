package storagepb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testTemplate = &Template{
		Id:       "id",
		Name:     "name",
		Contents: []byte("{{.Value}}"),
	}
)

func TestTemplateParse(t *testing.T) {
	cases := []struct {
		json     string
		template *Template
	}{
		{`{"id": "id", "name": "name", "contents": "e3suVmFsdWV9fQ=="}`, testTemplate},
	}
	for _, c := range cases {
		template, err := ParseTemplate([]byte(c.json))
		assert.Nil(t, err)
		assert.Equal(t, c.template, template)
	}
}

func TestTemplateValidate(t *testing.T) {
	cases := []struct {
		template *Template
		valid    bool
	}{
		{testTemplate, true},
		{&Template{Id: "a1b2c3d4"}, true},
		{&Template{}, false},
	}
	for _, c := range cases {
		valid := c.template.AssertValid() == nil
		assert.Equal(t, c.valid, valid)
	}
}

func TestTemplateCopy(t *testing.T) {
	template := &Template{
		Id:       "id",
		Name:     "Long and descriptive name",
		Contents: []byte("contents"),
	}
	clone := template.Copy()
	// assert that:
	// - Profile fields are copied to the clone
	// - Mutation of the clone does not affect the original
	assert.Equal(t, template.Id, clone.Id)
	assert.Equal(t, template.Name, clone.Name)
	assert.Equal(t, template.Contents, clone.Contents)
}
