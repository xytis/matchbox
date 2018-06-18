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
		core := s.core
		ctx, err := s.unwrapContext(req.Context())
		logger := s.logger.With(zap.String("request-id", ctx.RequestID))
		if err != nil {
			logger.Debug("context not valid", zap.Error(err))
			http.Error(w, fmt.Sprintf(`404 context build error: %v`, err), http.StatusNotFound)
			return
		}

		templateID, present := ctx.Profile.Template["grub"]
		if !present {
			logger.Debug("template binding for grub is not set",
				zap.String("group", ctx.Group.Id),
				zap.String("profile", ctx.Profile.Id),
			)
			http.Error(w, fmt.Sprintf(`404 template binding for "grub" is not set in profile "%s"`, ctx.Profile.Id), http.StatusNotFound)
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
