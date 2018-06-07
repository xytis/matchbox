package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/coreos/matchbox/matchbox/server"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
	fake "github.com/coreos/matchbox/matchbox/storage/testfakes"
)

func TestIgnitionHandler_V2JSON(t *testing.T) {
	content := `{"ignition":{"version":"2.1.0","config":{}},"storage":{},"systemd":{"units":[{"name":"etcd2.service","enable":true}]},"networkd":{},"passwd":{}}`
	profile := &storagepb.Profile{
		Id:       fake.Group.Profile,
		Template: map[string]string{"ignition": fake.Template.Id},
	}
	store := fake.NewFixedStore()
	store.Templates[fake.Template.Id] = &storagepb.Template{Id: fake.Template.Id, Contents: []byte(content)}

	core := server.NewServer(store)
	srv := NewServer(&Config{Logger: zap.NewNop(), Core: core})
	h := srv.ignitionHandler()
	ctx := createFakeContext(context.Background(), map[string]string{}, profile, fake.Group)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req.WithContext(ctx))
	// assert that:
	// - raw Ignition config served directly
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonContentType, w.HeaderMap.Get(contentType))
	assert.Equal(t, content, w.Body.String())
}

/*
func TestIgnitionHandler_V2YAML(t *testing.T) {
	// exercise templating features, not a realistic Container Linux Config template
	content := `
systemd:
  units:
    - name: {{.service_name}}.service
      enable: true
    - name: {{.uuid}}.service
      enable: true
    - name: {{.request.query.foo}}.service
      enable: true
      contents: {{.request.raw_query}}
`
	expectedIgnitionV2 := `{"ignition":{"config":{},"security":{"tls":{}},"timeouts":{},"version":"2.2.0"},"networkd":{},"passwd":{},"storage":{},"systemd":{"units":[{"enable":true,"name":"etcd2.service"},{"enable":true,"name":"a1b2c3d4.service"},{"contents":"foo=some-param\u0026bar=b","enable":true,"name":"some-param.service"}]}}`
	store := &fake.FixedStore{
		Profiles:  map[string]*storagepb.Profile{fake.Group.Profile: testProfileIgnitionYAML},
		Groups:    map[string]*storagepb.Group{fake.Group.Id: fake.Group},
		Templates: map[string]*storagepb.Template{fake.Template.Id: &storagepb.Template{Id: fake.Template.Id, Contents: []byte(content)}},
	}
	logger, _ := logtest.NewNullLogger()
	core := server.NewServer(&server.Config{Store: store})
	srv := NewServer(&Config{Logger: logger, Core: core})
	h := srv.ignitionHandler()
	ctx := withGroup(context.Background(), fake.Group)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/?foo=some-param&bar=b", nil)
	h.ServeHTTP(w, req.WithContext(ctx))
	// assert that:
	// - Container Linux Config template rendered with Group selectors, metadata, and query variables
	// - Transformed to an Ignition config (JSON)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonContentType, w.HeaderMap.Get(contentType))
	assert.Equal(t, expectedIgnitionV2, w.Body.String())
}
*/

func TestIgnitionHandler_MissingCtxProfile(t *testing.T) {
	core := server.NewServer(fake.NewFixedStore())
	srv := NewServer(&Config{Logger: zap.NewNop(), Core: core})
	h := srv.ignitionHandler()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestIgnitionHandler_MissingIgnitionConfig(t *testing.T) {
	core := server.NewServer(fake.NewFixedStore())
	srv := NewServer(&Config{Logger: zap.NewNop(), Core: core})
	h := srv.ignitionHandler()
	ctx := withProfile(context.Background(), fake.Profile)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req.WithContext(ctx))
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestIgnitionHandler_MissingTemplateMetadata(t *testing.T) {
	content := `
ignition_version: 1
systemd:
  units:
    - name: {{.missing_key}}
      enable: true
`
	store := fake.NewFixedStore()
	store.Templates[fake.Template.Id] = &storagepb.Template{Id: fake.Template.Id, Contents: []byte(content)}

	core := server.NewServer(store)
	srv := NewServer(&Config{Logger: zap.NewNop(), Core: core})
	h := srv.ignitionHandler()
	ctx := withGroup(context.Background(), fake.Group)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req.WithContext(ctx))
	// assert that:
	// - Ignition template rendering errors because "missing_key" is not
	// present in the template variables
	assert.Equal(t, http.StatusNotFound, w.Code)
}
