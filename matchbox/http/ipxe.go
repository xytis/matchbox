package http

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

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
			s.logger.WithFields(logrus.Fields{
				"labels": labels,
			}).Infof("No matching group")
			http.NotFound(w, req)
			return
		}

		profile, err := profileFromContext(ctx)
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"labels": labels,
			}).Infof("No matching profile")
			http.NotFound(w, req)
			return
		}

		metadata, err := mergeMetadata(ctx)
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"labels":  labels,
				"group":   group.Id,
				"profile": profile.Id,
				"error":   err,
			}).Warnf("Issue with metadata")
		}

		templateID, present := profile.Template["ipxe"]
		if !present {
			templateID = "default-ipxe"
		}
		tmpl, err := core.TemplateGet(ctx, &pb.TemplateGetRequest{Id: templateID})
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"labels":  labels,
				"group":   group.Id,
				"profile": profile.Id,
			}).Infof("No template named: %s", templateID)
			http.NotFound(w, req)
			return
		}

		// match was successful
		s.logger.WithFields(logrus.Fields{
			"labels":  labelsFromRequest(nil, req),
			"group":   group.Id,
			"profile": profile.Id,
		}).Debug("Matched an iPXE config")

		var buf bytes.Buffer
		if err = Render(&buf, tmpl.Id, string(tmpl.Contents), metadata); err != nil {
			s.logger.Errorf("error rendering template: %v", err)
			http.NotFound(w, req)
			return
		}
		if _, err := buf.WriteTo(w); err != nil {
			s.logger.Errorf("error writing to response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	return http.HandlerFunc(fn)
}
