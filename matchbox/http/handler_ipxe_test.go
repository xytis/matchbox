package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"context"
	"github.com/stretchr/testify/assert"

	"github.com/coreos/matchbox/matchbox/storage/storagepb"
	fake "github.com/coreos/matchbox/matchbox/storage/testfakes"
)

func prepareProfileIPXE() *storagepb.Profile {
	profile := fake.Profile()
	profile.Template["ipxe"] = fake.IPXETemplate().Id
	return profile
}

func TestIPXEInspect(t *testing.T) {
	h := ipxeInspect()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, ipxeBootstrap, w.Body.String())
}

func TestIPXEHandler(t *testing.T) {
	srv := prepareServer()
	h := srv.ipxeHandler()

	ctx := createFakeContext(context.Background(), fake.Labels(), fake.Group(), prepareProfileIPXE())

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
