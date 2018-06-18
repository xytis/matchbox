package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/coreos/matchbox/matchbox/server"
	fake "github.com/coreos/matchbox/matchbox/storage/testfakes"
)

func TestSelectContext(t *testing.T) {
	store := fake.NewFixedStore()
	store.GroupPut(fake.Group())
	store.ProfilePut(fake.Profile())

	c := server.NewServer(store)
	srv := NewServer(&Config{Core: c, Logger: zap.NewNop()})
	next := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		labels, err := labelsFromContext(ctx)
		assert.Nil(t, err)
		assert.Equal(t, "a1b2c3d4", labels["uuid"])

		group, err := groupFromContext(ctx)
		assert.Nil(t, err)
		assert.Equal(t, fake.Group(), group)

		profile, err := profileFromContext(ctx)
		assert.Nil(t, err)
		assert.Equal(t, fake.Profile(), profile)

		fmt.Fprintf(w, "next handler called")
	}
	// assert that:
	// - query params are parsed to labels
	// - labels are used to match uuid=a1b2c3d4 to fake.Group
	// - fake.Group is used to match to fake.Profile
	// - all of the above are added to the context
	// - next handler is called
	h := srv.wrapContext(http.HandlerFunc(next))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "?uuid=a1b2c3d4", nil)
	h.ServeHTTP(w, req)
	assert.Equal(t, "next handler called", w.Body.String())
}
