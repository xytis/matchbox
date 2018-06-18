package http

import (
	"context"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

func TestContextProfile(t *testing.T) {
	expectedProfile := &storagepb.Profile{Id: "g1h2i3j4"}
	ctx := withProfile(context.Background(), expectedProfile)
	profile, err := profileFromContext(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectedProfile, profile)
}

func TestContextProfile_Error(t *testing.T) {
	profile, err := profileFromContext(context.Background())
	assert.Nil(t, profile)
	if assert.NotNil(t, err) {
		assert.Equal(t, errNoProfileFromContext, err)
	}
}

func TestContextGroup(t *testing.T) {
	expectedGroup := &storagepb.Group{Name: "test group"}
	ctx := withGroup(context.Background(), expectedGroup)
	group, err := groupFromContext(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectedGroup, group)
}

func TestContextGroup_Error(t *testing.T) {
	group, err := groupFromContext(context.Background())
	assert.Nil(t, group)
	if assert.NotNil(t, err) {
		assert.Equal(t, errNoGroupFromContext, err)
	}
}

func createFakeContext(
	ctx context.Context,
	labels map[string]string,
	group *storagepb.Group,
	profile *storagepb.Profile,
) context.Context {
	if labels != nil {
		ctx = withLabels(ctx, labels)
	}
	if profile != nil {
		ctx = withProfile(ctx, profile)
	}
	if group != nil {
		ctx = withGroup(ctx, group)
	}
	return ctx
}

func wrapFakeContext(ctx context.Context, vars map[string]string, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		req = req.WithContext(ctx)
		if vars != nil {
			req = mux.SetURLVars(req, vars)
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}
