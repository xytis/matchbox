package http

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coreos/matchbox/matchbox/storage/storagepb"
	fake "github.com/coreos/matchbox/matchbox/storage/testfakes"
)

func TestTemplateHandler(t *testing.T) {
	content := `#foo-bar-baz template
UUID={{.label.uuid}}
SERVICE={{.service_name}}
FOO={{.label.foo}}
`
	expected := `#foo-bar-baz template
UUID=a1s2d3
SERVICE=etcd2
FOO=some-param
`

	template := &storagepb.Template{
		Id:       "fake-template",
		Contents: []byte(content),
	}
	store := fake.NewFixedStore()
	store.TemplatePut(template)

	profile := fake.Profile()
	profile.Template = map[string]string{"template": template.Id}

	labels := map[string]string{}
	labels["uuid"] = "a1s2d3"
	labels["foo"] = "some-param"

	vars := map[string]string{"selector": "template"}

	srv := prepareServerWithStore(store)

	ctx := createFakeContext(context.Background(), labels, fake.Group(), profile)

	h := wrapFakeContext(ctx, vars, srv.templateHandler())
	// assert that:
	// - Generic config is rendered with Group selectors, metadata, and query variables
	assert.HTTPSuccess(t, h.ServeHTTP, "GET", "/", nil, nil)
	assert.Equal(t, expected, assert.HTTPBody(h.ServeHTTP, "GET", "/", nil))
}
