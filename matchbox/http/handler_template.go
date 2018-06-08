package http

import (
	"bytes"
	"net/http"

	"go.uber.org/zap"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
)

// templateHandler returns a handler that responds with the generic config
// matching the request.
func (s *Server) templateHandler() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		core := s.core

		ctx, err := s.unwrapContext(req.Context())
		if err != nil {
			s.logger.Info("context not valid", zap.Error(err))
			http.NotFound(w, req)
			return
		}

		selector, present := ctx.Labels["template_id"]
		if !present {
			http.Error(w, "Malformed URL, template_id query param must be set", 400)
			return
		}

		templateID, present := ctx.Profile.Template[selector]
		if !present {
			s.logger.Info("profile does not contain requested template",
				zap.String("profile", ctx.Profile.Id),
				zap.String("template", selector),
			)
			http.NotFound(w, req)
			return
		}

		tmpl, err := core.TemplateGet(ctx, &pb.TemplateGetRequest{Id: templateID})
		if err != nil {
			s.logger.Info("template not found",
				zap.String("template", templateID),
				zap.String("profile", ctx.Profile.Id),
			)
			http.NotFound(w, req)
			return
		}

		var buf bytes.Buffer
		if err = Render(&buf, tmpl.Id, string(tmpl.Contents), ctx.Metadata); err != nil {
			s.logger.Error("error rendering template", zap.Error(err))
			http.NotFound(w, req)
			return
		}
		if _, err := buf.WriteTo(w); err != nil {
			s.logger.Error("error writing to response", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	return http.HandlerFunc(fn)
}
