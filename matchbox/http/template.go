package http

import (
	"bytes"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
)

// templateHandler returns a handler that responds with the generic config
// matching the request.
func (s *Server) templateHandler() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		core := s.core
		labels, _ := labelsFromContext(ctx)

		selector, present := labels["template_id"]
		if !present {
			http.Error(w, "Malformed URL, template_id query param must be set", 400)
			return
		}

		_, err := groupFromContext(ctx)
		if err != nil {
			s.logger.Info("group not matched",
				zap.String("labels", fmt.Sprintf("%v", labels)),
			)
			http.NotFound(w, req)
			return
		}

		profile, err := profileFromContext(ctx)
		if err != nil {
			s.logger.Info("profile not matched",
				zap.String("labels", fmt.Sprintf("%v", labels)),
			)
			http.NotFound(w, req)
			return
		}

		templateID, present := profile.Template[selector]
		if !present {
			s.logger.Info("profile does not contain requested template",
				zap.String("profile", profile.Id),
				zap.String("template", selector),
			)
			http.NotFound(w, req)
			return
		}

		metadata, err := mergeMetadata(ctx)
		if err != nil {
			s.logger.Info("metadata not merged",
				zap.Error(err),
				zap.String("labels", fmt.Sprintf("%v", labels)),
				zap.String("profile", profile.Id),
			)
		}

		tmpl, err := core.TemplateGet(ctx, &pb.TemplateGetRequest{Id: templateID})
		if err != nil {
			s.logger.Info("template not found",
				zap.String("template", templateID),
				zap.String("labels", fmt.Sprintf("%v", labels)),
				zap.String("profile", profile.Id),
			)
			http.NotFound(w, req)
			return
		}

		var buf bytes.Buffer
		if err = Render(&buf, tmpl.Id, string(tmpl.Contents), metadata); err != nil {
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
