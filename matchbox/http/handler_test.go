package http

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/coreos/matchbox/matchbox/server"
	"github.com/coreos/matchbox/matchbox/storage"
	fake "github.com/coreos/matchbox/matchbox/storage/testfakes"
)

func prepareEmptyServer() *Server {
	store := fake.NewFixedStore()

	core := server.NewServer(store)
	srv := NewServer(&Config{Logger: zap.NewNop(), Core: core})
	return srv
}

func prepareServerWithStore(store storage.Store) *Server {
	core := server.NewServer(store)
	srv := NewServer(&Config{Logger: zap.NewNop(), Core: core})
	return srv
}

func prepareServer() *Server {
	store := fake.NewFixedStore()

	store.GroupPut(fake.Group())

	store.ProfilePut(fake.Profile())

	store.TemplatePut(fake.GrubTemplate())
	store.TemplatePut(fake.IPXETemplate())
	store.TemplatePut(fake.IgnitionTemplate())
	store.TemplatePut(fake.CustomTemplate())

	core := server.NewServer(store)
	srv := NewServer(&Config{Logger: zap.NewExample(), Core: core})
	return srv
}

func TestHandler_MissingCtxValues(t *testing.T) {
	srv := prepareEmptyServer()

	// All handlers should act the same on context issues
	handlers := []struct {
		name    string
		handler http.Handler
		vars    map[string]string
	}{
		{"grub", srv.grubHandler(), nil},
		{"ipxe", srv.ipxeHandler(), nil},
		{"ignition", srv.ignitionHandler(), nil},
		{"template", srv.templateHandler(), map[string]string{"selector": "custom"}},
	}

	cases := []struct {
		ctx context.Context
		err string
	}{
		{createFakeContext(context.Background(), nil, fake.Group(), fake.Profile()), "Context missing parsed Labels"},
		{createFakeContext(context.Background(), fake.Labels(), nil, fake.Profile()), "Context missing a Group"},
		{createFakeContext(context.Background(), fake.Labels(), fake.Group(), nil), "Context missing a Profile"},
	}

	for _, h := range handlers {
		for _, c := range cases {
			handler := wrapFakeContext(c.ctx, h.vars, h.handler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/", nil)
			handler.ServeHTTP(w, req)
			assert.Equalf(t, http.StatusNotFound, w.Code, `%s handler`, h.name)

			res := w.Result()
			body, _ := ioutil.ReadAll(res.Body)
			assert.Containsf(t, string(body), c.err, `%s handler`, h.name)
		}
	}
}

func TestHandler_MissingTemplateBindings(t *testing.T) {
	srv := prepareEmptyServer()

	// All handlers should act the same on profile lacking a binding
	handlers := []struct {
		name    string
		handler http.Handler
		vars    map[string]string
	}{
		{"grub", srv.grubHandler(), nil},
		{"ipxe", srv.ipxeHandler(), nil},
		{"ignition", srv.ignitionHandler(), nil},
		{"template", srv.templateHandler(), map[string]string{"selector": "custom"}},
	}

	profile := fake.Profile()
	// Clear any bindings
	profile.Template = map[string]string{}

	ctx := createFakeContext(context.Background(), fake.Labels(), fake.Group(), profile)

	for _, h := range handlers {
		handler := wrapFakeContext(ctx, h.vars, h.handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		handler.ServeHTTP(w, req)
		assert.Equalf(t, http.StatusNotFound, w.Code, `%s handler`, h.name)

		res := w.Result()
		body, _ := ioutil.ReadAll(res.Body)
		assert.Containsf(t, string(body), "template binding for", `%s handler`, h.name)
	}
}

func TestHandler_MissingTemplateValues(t *testing.T) {
	srv := prepareServer()

	// All handlers should act the same on profile lacking a binding
	handlers := []struct {
		name    string
		handler http.Handler
		vars    map[string]string
	}{
		{"grub", srv.grubHandler(), nil},
		{"ipxe", srv.ipxeHandler(), nil},
		{"ignition", srv.ignitionHandler(), nil},
		{"template", srv.templateHandler(), map[string]string{"selector": "custom"}},
	}

	//Clear metadata
	labels := fake.LabelsEmpty()
	group := fake.GroupNoMetadata()
	profile := fake.ProfileNoMetadata()

	ctx := createFakeContext(context.Background(), labels, group, profile)

	for _, h := range handlers {
		handler := wrapFakeContext(ctx, h.vars, h.handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		handler.ServeHTTP(w, req)
		assert.Equalf(t, http.StatusNotFound, w.Code, `%s handler`, h.name)

		res := w.Result()
		body, _ := ioutil.ReadAll(res.Body)
		assert.Containsf(t, string(body), "template binding for", `%s handler`, h.name)
	}
}

func TestHandler_WriteError(t *testing.T) {
	srv := prepareServer()

	// All handlers should act the same on write errors
	handlers := []struct {
		name    string
		handler http.Handler
		vars    map[string]string
	}{
		{"grub", srv.grubHandler(), nil},
		{"ipxe", srv.ipxeHandler(), nil},
		{"ignition", srv.ignitionHandler(), nil},
		{"template", srv.templateHandler(), map[string]string{"selector": "custom"}},
	}

	ctx := createFakeContext(context.Background(), fake.Labels(), fake.Group(), fake.Profile())

	for _, h := range handlers {
		handler := wrapFakeContext(ctx, h.vars, h.handler)

		w := NewUnwriteableResponseWriter()
		req, _ := http.NewRequest("GET", "/", nil)
		handler.ServeHTTP(w, req)
		if !assert.Equalf(t, http.StatusInternalServerError, w.Code, `%s handler`, h.name) {
			assert.Emptyf(t, w.Body.String(), `%s handler`, h.name)
		}
	}
}
