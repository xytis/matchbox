package http

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
)

// templateHandler returns a handler that responds with the generic config
// matching the request.
func (s *Server) templateHandler() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		core := s.core
		vars := mux.Vars(req)

		ctx, err := s.unwrapContext(req.Context())
		logger := s.logger.With(zap.String("request-id", ctx.RequestID))
		if err != nil {
			logger.Debug("context not valid", zap.Error(err))
			http.Error(w, fmt.Sprintf(`404 context build error: %v`, err), http.StatusNotFound)
			return
		}

		selector, present := vars["selector"]
		if !present {
			http.Error(w, "Malformed URL, please specify /template/{selector}[.sig|.asc]", 400)
			return
		}

		templateID, present := ctx.Profile.Template[selector]
		if !present {
			logger.Debug("template binding for selector is not set",
				zap.String("profile", ctx.Profile.Id),
				zap.String("selector", selector),
			)
			http.Error(w, fmt.Sprintf(`404 template binding for "%s" is not set in profile "%s"`, selector, ctx.Profile.Id), http.StatusNotFound)
			return
		}

		tmpl, err := core.TemplateGet(ctx, &pb.TemplateGetRequest{Id: templateID})
		if err != nil {
			logger.Debug("template not found",
				zap.String("template", templateID),
				zap.String("group", ctx.Group.Id),
				zap.String("profile", ctx.Profile.Id),
			)
			http.Error(w, fmt.Sprintf(`404 template "%s" not found`, templateID), http.StatusNotFound)
			return
		}

		var buf bytes.Buffer
		if err = Render(&buf, tmpl.Id, string(tmpl.Contents), ctx.Metadata); err != nil {
			logger.Debug("template rendering failure", zap.Error(err))
			http.Error(w, fmt.Sprintf("404 template rendering error: %v", err), http.StatusNotFound)
			return
		}
		if _, err := buf.WriteTo(w); err != nil {
			logger.Error("error writing to response", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	return http.HandlerFunc(fn)
}
