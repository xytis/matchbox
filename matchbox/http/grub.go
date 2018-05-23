package http

import (
	"bytes"
	"net/http"

	"github.com/sirupsen/logrus"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
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

		templateID, present := profile.Template["grub"]
		if !present {
			templateID = "default-grub"
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
