package http

import (
	"bytes"
	"fmt"
	"net/http"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"

	"go.uber.org/zap"
)

// grubHandler returns a handler which renders a GRUB2 config for the
// requester.
func (s *Server) grubHandler() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		core := s.core
		labels, _ := labelsFromContext(ctx)

		group, err := groupFromContext(ctx)
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

		metadata, err := mergeMetadata(ctx)
		if err != nil {
			s.logger.Info("metadata not merged",
				zap.Error(err),
				zap.String("labels", fmt.Sprintf("%v", labels)),
				zap.String("group", group.Id),
				zap.String("profile", profile.Id),
			)
		}

		templateID, present := profile.Template["grub"]
		if !present {
			templateID = "default-grub"
		}
		tmpl, err := core.TemplateGet(ctx, &pb.TemplateGetRequest{Id: templateID})
		if err != nil {
			s.logger.Info("template not found",
				zap.String("template", templateID),
				zap.String("labels", fmt.Sprintf("%v", labels)),
				zap.String("group", group.Id),
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
