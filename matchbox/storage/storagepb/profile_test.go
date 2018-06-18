package storagepb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testProfile = &Profile{
		Id:   "id",
		Name: "name",
		Template: map[string]string{
			"key": "value",
		},
		Metadata: []byte(`{"key":"value", "map":{"key":"value"}}`),
	}
)

func TestProfileParse(t *testing.T) {
	cases := []struct {
		json    string
		profile *Profile
	}{
		{`{"id":"id", "name":"name", "template":{"key":"value"}, "metadata":"eyJrZXkiOiJ2YWx1ZSIsICJtYXAiOnsia2V5IjoidmFsdWUifX0="}`, testProfile},
	}
	for _, c := range cases {
		profile, err := ParseProfile([]byte(c.json))
		assert.Nil(t, err)
		assert.Equal(t, c.profile, profile)
	}
}

func TestProfileValidate(t *testing.T) {
	cases := []struct {
		profile *Profile
		valid   bool
	}{
		{testProfile, true},
		{&Profile{Id: "a1b2c3d4"}, true},
		{&Profile{}, false},
	}
	for _, c := range cases {
		valid := c.profile.AssertValid() == nil
		assert.Equal(t, c.valid, valid)
	}
}

func TestProfileTemplateString(t *testing.T) {
	profile := Profile{
		Template: map[string]string{
			"a": "b",
			"c": "d",
		},
	}
	expected := "a=b,c=d"
	assert.Equal(t, expected, profile.TemplateString())
}

func TestProfileCopy(t *testing.T) {
	profile := testProfile
	clone := profile.Copy()
	// assert that:
	// - Profile fields are copied to the clone
	// - Mutation of the clone does not affect the original
	assert.Equal(t, profile.Id, clone.Id)
	assert.Equal(t, profile.Name, clone.Name)
	assert.Equal(t, profile.Template, clone.Template)
	assert.Equal(t, profile.Metadata, clone.Metadata)
}
