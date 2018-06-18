package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coreos/matchbox/matchbox/storage/storagepb"
	fake "github.com/coreos/matchbox/matchbox/storage/testfakes"
)

func prepareProfileGrub() *storagepb.Profile {
	profile := fake.Profile()
	profile.Template["grub"] = fake.GrubTemplate().Id
	return profile
}

func TestGrubHandler(t *testing.T) {
	srv := prepareServer()
	h := srv.grubHandler()

	ctx := createFakeContext(context.Background(), fake.Labels(), fake.Group(), prepareProfileGrub())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req.WithContext(ctx))
	// assert that:
	// - the Profile's NetBoot config is rendered as a GRUB2 config
	expectedScript := `default=0
fallback=1
timeout=1
menuentry "CoreOS (EFI)" {
  echo "Loading kernel"
  linuxefi "/image/kernel" a=b c
  echo "Loading initrd"
  initrdefi  "/image/initrd_a" "/image/initrd_b"
}`
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedScript, w.Body.String())
}
