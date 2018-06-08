package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/coreos/matchbox/matchbox/server"
	fake "github.com/coreos/matchbox/matchbox/storage/testfakes"
)

func TestGrubHandler(t *testing.T) {
	core := server.NewServer(fake.NewFixedStore())
	srv := NewServer(&Config{Logger: zap.NewNop(), Core: core})
	h := srv.grubHandler()

	ctx := createFakeContext(context.Background(), map[string]string{}, fake.Profile, fake.Group)

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
}
menuentry "CoreOS (BIOS)" {
  echo "Loading kernel"
  linux "/image/kernel" a=b c
  echo "Loading initrd"
  initrd  "/image/initrd_a" "/image/initrd_b"
}
`
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedScript, w.Body.String())
}
