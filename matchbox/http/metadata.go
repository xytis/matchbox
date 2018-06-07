package http

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

const plainContentType = "plain/text"

// metadataHandler returns a handler that responds with the metadata env file
// matching the request.
func (s *Server) metadataHandler() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		var err error

		_, err = groupFromContext(ctx)
		if err != nil {
			http.NotFound(w, req)
			return
		}

		_, err = profileFromContext(ctx)
		if err != nil {
			http.NotFound(w, req)
			return
		}

		// collect data for rendering
		metadata, err := mergeMetadata(ctx)
		if err != nil {
			s.logger.Warn("metadata not merged", zap.Error(err))
			http.NotFound(w, req)
			return
		}

		w.Header().Set(contentType, plainContentType)
		renderAsEnvFile(w, "", metadata)
	}
	return http.HandlerFunc(fn)
}

// renderAsEnvFile writes map data into a KEY=value\n "env file" format,
// descending recursively into nested maps and prepending parent keys.
//
// For example, {"outer":{"inner":"val"}} -> OUTER_INNER=val). Note that
// structure is lost in this transformation, the inverse transfom has two
// possible outputs.
func renderAsEnvFile(w io.Writer, prefix string, root map[string]interface{}) {
	for key, value := range root {
		name := prefix + key
		switch val := value.(type) {
		case string, bool, float64:
			// simple JSON unmarshal types
			fmt.Fprintf(w, "%s=%v\n", strings.ToUpper(name), val)
		case map[string]string:
			m := map[string]interface{}{}
			for k, v := range val {
				m[k] = v
			}
			renderAsEnvFile(w, name+"_", m)
		case map[string]interface{}:
			renderAsEnvFile(w, name+"_", val)
		}
	}
}
