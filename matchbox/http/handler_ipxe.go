package http

import (
	"bytes"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
)

const ipxeBootstrap = `#!ipxe
chain ipxe?uuid=${uuid}&mac=${mac:hexhyp}&domain=${domain}&hostname=${hostname}&serial=${serial}
`

// ipxeInspect returns a handler that responds with the iPXE script to gather
// client machine data and chainload to the ipxeHandler.
func ipxeInspect() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, ipxeBootstrap)
	}
	return http.HandlerFunc(fn)
}

// ipxeBoot returns a handler which renders the iPXE boot script for the
// requester.
func (s *Server) ipxeHandler() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		core := s.core

		ctx, err := s.unwrapContext(req.Context())
		if err != nil {
			s.logger.Info("context not valid", zap.Error(err))
			http.NotFound(w, req)
			return
		}

		templateID, present := ctx.Profile.Template["ipxe"]
		if !present {
			s.logger.Info("template binding for ipxe is not set",
				zap.String("group", ctx.Group.Id),
				zap.String("profile", ctx.Profile.Id),
			)
			http.NotFound(w, req)
			return
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
