package http

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestRenderJSON(t *testing.T) {
	srv := NewServer(&Config{Logger: zap.NewNop()})
	w := httptest.NewRecorder()
	data := map[string][]string{
		"a": []string{"b", "c"},
	}
	srv.renderJSON(w, data)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonContentType, w.HeaderMap.Get(contentType))
	assert.Equal(t, `{"a":["b","c"]}`, w.Body.String())
}

func TestRenderJSON_EncodingError(t *testing.T) {
	srv := NewServer(&Config{Logger: zap.NewNop()})
	w := httptest.NewRecorder()
	// channels cannot be JSON encoded
	srv.renderJSON(w, make(chan struct{}))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestRenderJSON_EncodeError(t *testing.T) {
	srv := NewServer(&Config{Logger: zap.NewNop()})
	w := httptest.NewRecorder()
	// channels cannot be JSON encoded
	srv.renderJSON(w, make(chan struct{}))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestRenderJSON_WriteError(t *testing.T) {
	srv := NewServer(&Config{Logger: zap.NewNop()})
	w := NewUnwriteableResponseWriter()
	srv.renderJSON(w, map[string]string{"a": "b"})
	if !assert.Equal(t, http.StatusInternalServerError, w.Code) {
		assert.Empty(t, w.Body.String())
	}
}

// UnwritableResponseWriter is a http.ResponseWriter for testing Write
// failures.
type UnwriteableResponseWriter struct {
	*httptest.ResponseRecorder
	Body *bytes.Buffer
}

func NewUnwriteableResponseWriter() *UnwriteableResponseWriter {
	return &UnwriteableResponseWriter{httptest.NewRecorder(), new(bytes.Buffer)}
}

func (w *UnwriteableResponseWriter) Write(buf []byte) (int, error) {
	// Preserve what was written for quicker fixing
	// Note that actual Header is not modified, contrary to ResponseRecorder
	w.Body.Write(buf)
	return 0, fmt.Errorf("Unwriteable ResponseWriter")
}
