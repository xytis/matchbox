package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"context"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	"github.com/coreos/matchbox/matchbox/server"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
	fake "github.com/coreos/matchbox/matchbox/storage/testfakes"
)

func TestIPXEInspect(t *testing.T) {
	h := ipxeInspect()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, ipxeBootstrap, w.Body.String())
}

func TestIPXEHandler(t *testing.T) {
	logger, _ := logtest.NewNullLogger()
	core := server.NewServer(&server.Config{Store: fake.NewFixedStore()})
	srv := NewServer(&Config{Logger: logger, Core: core})
	h := srv.ipxeHandler()
	ctx := createFakeContext(context.Background(), map[string]string{}, fake.Profile, fake.Group)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req.WithContext(ctx))
	// assert that:
	// - the Profile's NetBoot config is rendered as an iPXE script
	expectedScript := `#!ipxe
kernel /image/kernel a=b c
initrd /image/initrd_a
initrd /image/initrd_b
boot
`
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedScript, w.Body.String())
}

func TestIPXEHandler_MissingCtxProfile(t *testing.T) {
	logger, _ := logtest.NewNullLogger()
	core := server.NewServer(&server.Config{Store: fake.NewFixedStore()})
	srv := NewServer(&Config{Logger: logger, Core: core})
	h := srv.ipxeHandler()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestIPXEHandler_RenderTemplateError(t *testing.T) {
	logger, _ := logtest.NewNullLogger()
	core := server.NewServer(&server.Config{Store: fake.NewFixedStore()})
	srv := NewServer(&Config{Logger: logger, Core: core})
	h := srv.ipxeHandler()
	//Profile with missing metadata produces template render error
	ctx := createFakeContext(context.Background(), map[string]string{}, &storagepb.Profile{Id: fake.Profile.Id}, fake.Group)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req.WithContext(ctx))
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestIPXEHandler_WriteError(t *testing.T) {
	logger, _ := logtest.NewNullLogger()
	core := server.NewServer(&server.Config{Store: fake.NewFixedStore()})
	srv := NewServer(&Config{Logger: logger, Core: core})
	h := srv.ipxeHandler()
	ctx := createFakeContext(context.Background(), map[string]string{}, fake.Profile, fake.Group)
	w := NewUnwriteableResponseWriter()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req.WithContext(ctx))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Empty(t, w.Body.String())
}
