package http

import (
	"bytes"
	"net/http"

	"github.com/sirupsen/logrus"

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

		templateID, present := profile.Template[selector]
		if !present {
			s.logger.WithFields(logrus.Fields{
				"labels":           labels,
				"group":            group.Id,
				"profile":          profile.Id,
				"profile_template": selector,
			}).Infof("Profile does not contain requested template binding")
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
