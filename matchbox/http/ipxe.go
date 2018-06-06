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

		templateID, present := profile.Template["ipxe"]
		if !present {
			templateID = "default-ipxe"
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
