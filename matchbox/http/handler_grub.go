package http

import (
	"bytes"
	"net/http"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"

	"go.uber.org/zap"
)

// grubHandler returns a handler which renders a GRUB2 config for the
// requester.
func (s *Server) grubHandler() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		core := s.core

		ctx, err := s.unwrapContext(req.Context())
		if err != nil {
			s.logger.Info("context not valid", zap.Error(err))
			http.NotFound(w, req)
			return
		}

		templateID, present := ctx.Profile.Template["grub"]
		if !present {
			templateID = "default-grub"
		}
		tmpl, err := core.TemplateGet(ctx, &pb.TemplateGetRequest{Id: templateID})
		if err != nil {
			s.logger.Info("template not found",
				zap.String("template", templateID),
				zap.String("group", ctx.Group.Id),
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
