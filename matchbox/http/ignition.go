package http

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	//ct "github.com/coreos/container-linux-config-transpiler/config"
	ignition "github.com/coreos/ignition/config"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
)

// ignitionHandler returns a handler that responds with the Ignition config
// matching the request. The Ignition file referenced in the Profile is parsed
// as raw Ignition (for .ign/.ignition) or rendered from a Container Linux
// Config (YAML) and converted to Ignition. Ignition configs are served as HTTP
// JSON responses.
func (s *Server) ignitionHandler() http.Handler {
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

		templateID, present := profile.Template["ignition"]
		if !present {
			templateID = "default-ignition"
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
		_, report, err := ignition.Parse(buf.Bytes())
		if err != nil {
			s.logger.Warningf("warning parsing Ignition JSON: %s", report.String())
		}
		w.Header().Set(contentType, jsonContentType)
		if _, err := buf.WriteTo(w); err != nil {
			s.logger.Errorf("error writing to response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
		/*
			// Container Linux Config template

			// collect data for rendering
			data, err := collectVariables(req, group)
			if err != nil {
				s.logger.Errorf("error collecting variables: %v", err)
				http.NotFound(w, req)
				return
			}

			// render the template for an Ignition config with data
			var buf bytes.Buffer
			err = s.renderTemplate(&buf, data, contents)
			if err != nil {
				http.NotFound(w, req)
				return
			}

			// Parse bytes into a Container Linux Config
			config, ast, report := ct.Parse(buf.Bytes())
			if report.IsFatal() {
				s.logger.Errorf("error parsing Container Linux config: %s", report.String())
				http.NotFound(w, req)
				return
			}

			// Convert Container Linux Config into an Ignition Config
			ign, report := ct.Convert(config, "", ast)
			if report.IsFatal() {
				s.logger.Errorf("error converting Container Linux config: %s", report.String())
				http.NotFound(w, req)
				return
			}

			s.renderJSON(w, ign)
			return
		*/
	}
	return http.HandlerFunc(fn)
}

// isIgnition returns true if the file should be treated as plain Ignition.
func isIgnition(filename string) bool {
	return strings.HasSuffix(filename, ".ign") || strings.HasSuffix(filename, ".ignition")
}
