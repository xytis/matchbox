package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/Sirupsen/logrus"
	logtest "github.com/Sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	"github.com/coreos/matchbox/matchbox/server"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
	fake "github.com/coreos/matchbox/matchbox/storage/testfakes"
)

func TestGenericHandler(t *testing.T) {
	content := `#foo-bar-baz template
UUID={{.uuid}}
SERVICE={{.service_name}}
FOO={{.labels.foo}}
`
	expected := `#foo-bar-baz template
UUID=a1b2c3d4
SERVICE=etcd2
FOO=some-param
`
	store := fake.NewFixedStore()
	store.Templates[fake.Template.Id] = &storagepb.Template{Id: fake.Template.Id, Contents: []byte(content)}

	//logger, _ := logtest.NewNullLogger()
	logger := log.New()
	core := server.NewServer(&server.Config{Store: store})
	srv := NewServer(&Config{Logger: logger, Core: core})
	h := srv.templateHandler()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := createFakeContext(context.Background(), map[string]string{"template_id": "the-template", "foo": "some-param"}, &storagepb.Profile{Id: fake.Profile.Id, Template: map[string]string{"the-template": fake.Template.Id}}, fake.Group)
	h.ServeHTTP(w, req.WithContext(ctx))
	// assert that:
	// - Generic config is rendered with Group selectors, metadata, and query variables
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expected, w.Body.String())
}

func TestGenericHandler_MissingCtxProfile(t *testing.T) {
	logger, _ := logtest.NewNullLogger()
	core := server.NewServer(&server.Config{Store: fake.NewFixedStore()})
	srv := NewServer(&Config{Logger: logger, Core: core})
	h := srv.templateHandler()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := withLabels(context.Background(), map[string]string{"template_id": "any-template"})
	h.ServeHTTP(w, req.WithContext(ctx))
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGenericHandler_MissingTemplate(t *testing.T) {
	logger, _ := logtest.NewNullLogger()
	core := server.NewServer(&server.Config{Store: &fake.EmptyStore{}})
	srv := NewServer(&Config{Logger: logger, Core: core})
	h := srv.templateHandler()
	ctx := createFakeContext(context.Background(), map[string]string{"template_id": "any-template"}, &storagepb.Profile{Template: map[string]string{"any-template": "non-existing"}}, fake.Group)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req.WithContext(ctx))
	assert.Equal(t, http.StatusNotFound, w.Code)
}

/*
func TestGenericHandler_MissingTemplateMetadata(t *testing.T) {
	content := `#foo-bar-baz template
KEY={{.missing_key}}
`
	store := &fake.FixedStore{
		Profiles:       map[string]*storagepb.Profile{fake.Group.Profile: fake.Profile},
		GenericConfigs: map[string]string{fake.Profile.GenericId: content},
	}
	logger, _ := logtest.NewNullLogger()
	srv := NewServer(&Config{Logger: logger})
	c := server.NewServer(&server.Config{Store: store})
	h := srv.cloudHandler(c)
	ctx := withGroup(context.Background(), fake.Group)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req.WithContext(ctx))
	// assert that:
	// - Generic template rendering errors because "missing_key" is not
	// present in the template variables
	assert.Equal(t, http.StatusNotFound, w.Code)
}
*/
