package http

import (
	"bufio"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

func TestMetadataHandler(t *testing.T) {
	group := &storagepb.Group{
		Id:       "test-group",
		Selector: map[string]string{"mac": "52:54:00:a1:9c:ae"},
		Metadata: []byte(`{"meta":"group-data", "some":{"nested":{"data":"some-value"}}}`),
	}

	profile := &storagepb.Profile{
		Id:       "test-profile",
		Metadata: []byte(`{"meta":"profile-data", "some": {"nested": {"override": "value"}}}`),
	}

	labels := map[string]string{
		"custom": "value",
		"some":   "not-override",
	}

	srv := NewServer(&Config{Logger: zap.NewNop()})
	h := srv.metadataHandler()
	ctx := context.Background()
	ctx = withGroup(ctx, group)
	ctx = withProfile(ctx, profile)
	ctx = withLabels(ctx, labels)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req.WithContext(ctx))
	// assert that:
	// - Group selectors, metadata, and query variables are formatted
	// - nested metadata are namespaced
	// - key names are upper case
	// - key/value pairs are newline separated
	expectedLines := map[string]string{
		// merged metadata
		"META":                 "profile-data",
		"SOME_NESTED_DATA":     "some-value",
		"SOME_NESTED_OVERRIDE": "value",
		// group selector
		//"MAC": "52:54:00:a1:9c:ae",
		// labels
		"LABEL_CUSTOM": "value",
		"LABEL_SOME":   "not-override",
	}
	assert.Equal(t, http.StatusOK, w.Code)
	// convert response (random order) to map (tests compare in order)
	assert.Equal(t, expectedLines, metadataToMap(w.Body.String()))
	assert.Equal(t, plainContentType, w.HeaderMap.Get(contentType))
}

func TestMetadataHandler_MetadataEdgeCases(t *testing.T) {
	srv := NewServer(&Config{Logger: zap.NewNop()})
	h := srv.metadataHandler()
	// groups with different metadata
	cases := []struct {
		group    *storagepb.Group
		expected string
	}{
		{&storagepb.Group{Metadata: []byte(`{"num":3}`)}, "NUM=3\n"},
		{&storagepb.Group{Metadata: []byte(`{"yes":true}`)}, "YES=true\n"},
		{&storagepb.Group{Metadata: []byte(`{"no":false}`)}, "NO=false\n"},
	}
	for _, c := range cases {
		ctx := context.Background()
		ctx = withGroup(ctx, c.group)
		ctx = withProfile(ctx, &storagepb.Profile{})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		h.ServeHTTP(w, req.WithContext(ctx))
		// assert that:
		// - Group metadata key names are upper case
		// - key/value pairs are newline separated
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), c.expected)
		assert.Equal(t, plainContentType, w.HeaderMap.Get(contentType))
	}
}

func TestMetadataHandler_MissingCtxGroup(t *testing.T) {
	srv := NewServer(&Config{Logger: zap.NewNop()})
	h := srv.metadataHandler()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// metadataToMap converts a KEY=val\nKEY=val ResponseWriter body to a map for
// testing purposes.
func metadataToMap(metadata string) map[string]string {
	scanner := bufio.NewScanner(strings.NewReader(metadata))
	data := make(map[string]string)
	for scanner.Scan() {
		token := scanner.Text()
		pair := strings.SplitN(token, "=", 2)
		if len(pair) != 2 {
			continue
		}
		data[pair[0]] = pair[1]
	}
	return data
}
