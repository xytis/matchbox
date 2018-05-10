package http

import (
	"encoding/json"
	"io"
	"net/http"
	"text/template"
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
		s.logger.Errorf("error JSON encoding: %v", err)
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
		s.logger.Errorf("error writing to response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Render fills buffer with template render, does not allow missing keys
func Render(w io.Writer, name string, content string, data map[string]interface{}) error {
	tmpl := template.New(name).Option("missingkey=error")
	tmpl, err := tmpl.Parse(content)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

func (s *Server) renderTemplate(w io.Writer, data interface{}, contents ...string) (err error) {
	tmpl := template.New("").Option("missingkey=error")
	for _, content := range contents {
		tmpl, err = tmpl.Parse(content)
		if err != nil {
			s.logger.Errorf("error parsing template: %v", err)
			return err
		}
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		s.logger.Errorf("error rendering template: %v", err)
		return err
	}
	return nil
}
