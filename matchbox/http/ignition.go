package http

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
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

		templateID, present := profile.Template["ignition"]
		if !present {
			templateID = "default-ignition"
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
		_, report, err := ignition.Parse(buf.Bytes())
		if err != nil {
			s.logger.Warn("ignition parsing failed", zap.Error(err), zap.String("report", report.String()))
		}
		w.Header().Set(contentType, jsonContentType)
		if _, err := buf.WriteTo(w); err != nil {
			s.logger.Error("error writing to response", zap.Error(err))
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
