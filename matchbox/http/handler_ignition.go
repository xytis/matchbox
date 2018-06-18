package http

import (
	"bytes"
	"fmt"
	"net/http"

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
		core := s.core

		ctx, err := s.unwrapContext(req.Context())
		logger := s.logger.With(zap.String("request-id", ctx.RequestID))
		if err != nil {
			logger.Debug("context not valid", zap.Error(err))
			http.Error(w, fmt.Sprintf(`404 context build error: %v`, err), http.StatusNotFound)
			return
		}

		templateID, present := ctx.Profile.Template["ignition"]
		if !present {
			logger.Debug("template binding for ignition is not set",
				zap.String("group", ctx.Group.Id),
				zap.String("profile", ctx.Profile.Id),
			)
			http.Error(w, fmt.Sprintf(`404 template binding for "ignition" is not set in profile "%s"`, ctx.Profile.Id), http.StatusNotFound)
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
		_, report, err := ignition.Parse(buf.Bytes())
		if err != nil {
			logger.Debug("ignition parsing failed", zap.Error(err), zap.String("report", report.String()))
			http.Error(w, fmt.Sprintf("404 ignition parsing failed: %v\nreport: %s", err, report.String()), http.StatusNotFound)
		}
		w.Header().Set(contentType, jsonContentType)
		if _, err := buf.WriteTo(w); err != nil {
			logger.Error("error writing to response", zap.Error(err))
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
