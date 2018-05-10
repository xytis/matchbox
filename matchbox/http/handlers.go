package http

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/coreos/matchbox/matchbox/server"
	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
)

// homeHandler shows the server name for rooted requests. Otherwise, a 404 is
// returned.
func homeHandler() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "matchbox\n")
	}
	return http.HandlerFunc(fn)
}

func stripSuffix(suffix string, h http.Handler) http.Handler {
	if suffix == "" {
		return h
	}
	fn := func(w http.ResponseWriter, req *http.Request) {
		p := strings.TrimSuffix(req.URL.Path, suffix)
		r2 := new(http.Request)
		*r2 = *req
		r2.URL = new(url.URL)
		*r2.URL = *req.URL
		r2.URL.Path = p
		h.ServeHTTP(w, r2)
	}
	return http.HandlerFunc(fn)
}

// logRequest logs HTTP requests.
func (s *Server) logRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		s.logger.Infof("HTTP %s %v", req.Method, req.URL)
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

func (s *Server) selectContext(core server.Server, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		labels := labelsFromRequest(s.logger, req)
		ctx = withLabels(ctx, labels)
		if group, err := core.SelectGroup(ctx, &pb.SelectGroupRequest{Labels: labels}); err == nil {
			ctx = withGroup(ctx, group)
			if profile, err := core.ProfileGet(ctx, &pb.ProfileGetRequest{Id: group.Profile}); err == nil {
				ctx = withProfile(ctx, profile)
			}
		}
		next.ServeHTTP(w, req.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// selectGroup selects the Group whose selectors match the query parameters,
// adds the Group to the ctx, and calls the next handler. The next handler
// should handle a missing Group.
func (s *Server) selectGroup(core server.Server, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		attrs := labelsFromRequest(s.logger, req)
		// match machine request
		group, err := core.SelectGroup(ctx, &pb.SelectGroupRequest{Labels: attrs})
		if err == nil {
			// add the Group to the ctx for next handler
			ctx = withGroup(ctx, group)
		}
		next.ServeHTTP(w, req.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// selectProfile selects the Profile for the given query parameters, adds the
// Profile to the ctx, and calls the next handler. The next handler should
// handle a missing profile.
func (s *Server) selectProfile(core server.Server, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		attrs := labelsFromRequest(s.logger, req)
		// match machine request
		profile, err := core.SelectProfile(ctx, &pb.SelectProfileRequest{Labels: attrs})
		if err == nil {
			// add the Profile to the ctx for the next handler
			ctx = withProfile(ctx, profile)
		}
		next.ServeHTTP(w, req.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
