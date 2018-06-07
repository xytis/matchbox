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

func TestSelectGroup(t *testing.T) {
	store := fake.NewFixedStore()
	store.Groups[fake.Group.Id] = fake.Group

	srv := NewServer(&Config{Logger: zap.NewNop()})
	c := server.NewServer(store)
	next := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		group, err := groupFromContext(ctx)
		assert.Nil(t, err)
		assert.Equal(t, fake.Group, group)
		fmt.Fprintf(w, "next handler called")
	}
	// assert that:
	// - query params are used to match uuid=a1b2c3d4 to fake.Group
	// - the fake.Group is added to the context
	// - next handler is called
	h := srv.selectGroup(c, http.HandlerFunc(next))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "?uuid=a1b2c3d4", nil)
	h.ServeHTTP(w, req)
	assert.Equal(t, "next handler called", w.Body.String())
}

func TestSelectProfile(t *testing.T) {
	store := fake.NewFixedStore()
	store.Groups[fake.Group.Id] = fake.Group
	store.Profiles[fake.Profile.Id] = fake.Profile

	srv := NewServer(&Config{Logger: zap.NewNop()})
	c := server.NewServer(store)
	next := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		profile, err := profileFromContext(ctx)
		assert.Nil(t, err)
		assert.Equal(t, fake.Profile, profile)
		fmt.Fprintf(w, "next handler called")
	}
	// assert that:
	// - query params are used to match uuid=a1b2c3d4 to fake.Group's fakeProfile
	// - the fake.Profile is added to the context
	// - next handler is called
	h := srv.selectProfile(c, http.HandlerFunc(next))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "?uuid=a1b2c3d4", nil)
	h.ServeHTTP(w, req)
	assert.Equal(t, "next handler called", w.Body.String())
}
