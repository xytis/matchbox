package http

import (
	"encoding/json"
	"io"
	"net/http"
	"text/template"

	"github.com/masterminds/sprig"
	"go.uber.org/zap"
)

const (
	contentType     = "Content-Type"
	jsonContentType = "application/json"
)

// renderJSON encodes structs to JSON, writes the response to the
// ResponseWriter, and logs encoding errors.
func (s *Server) renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		s.logger.Error("JSON encoding failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, js)
}

// writeJSON writes the given bytes with a JSON Content-Type.
func (s *Server) writeJSON(w http.ResponseWriter, data []byte) {
	w.Header().Set(contentType, jsonContentType)
	_, err := w.Write(data)
	if err != nil {
		s.logger.Error("error writing to response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Render fills buffer with template render, does not allow missing keys
func Render(w io.Writer, name string, content string, data map[string]interface{}) error {
	tmpl := template.New(name).Funcs(sprig.TxtFuncMap()).Option("missingkey=error")
	tmpl, err := tmpl.Parse(content)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}
