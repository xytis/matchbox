package http

import (
	"net/http"
	"strings"

	"github.com/coreos/matchbox/matchbox/server"
	"github.com/coreos/matchbox/matchbox/sign"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Config configures a Server.
type Config struct {
	Core   server.Server
	Logger *zap.Logger
	// Path to static assets
	AssetsPath string
	// config signers (.sig and .asc)
	Signer        sign.Signer
	ArmoredSigner sign.Signer
}

// Server serves boot and provisioning configs to machines via HTTP.
type Server struct {
	core          server.Server
	logger        *zap.Logger
	assetsPath    string
	signer        sign.Signer
	armoredSigner sign.Signer
}

// NewServer returns a new Server.
func NewServer(config *Config) *Server {
	return &Server{
		core:          config.Core,
		logger:        config.Logger,
		assetsPath:    config.AssetsPath,
		signer:        config.Signer,
		armoredSigner: config.ArmoredSigner,
	}
}

// HTTPHandler returns a HTTP handler for the server.
func (s *Server) HTTPHandler() http.Handler {
	r := mux.NewRouter()

	// Logging
	r.Use(func(next http.Handler) http.Handler {
		return s.logRequest(next)
	})
	// Context parser
	r.Use(func(next http.Handler) http.Handler {
		return s.selectContext(s.core, next)
	})
	// Signature Handlers
	r.Use(func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, req *http.Request) {
			if s.signer != nil && strings.HasSuffix(req.URL.Path, ".sig") {
				h := stripSuffix(".sig", sign.SignatureHandler(s.signer, next))
				h.ServeHTTP(w, req)
			} else if s.armoredSigner != nil && strings.HasSuffix(req.URL.Path, ".asc") {
				h := stripSuffix(".asc", sign.SignatureHandler(s.armoredSigner, next))
				h.ServeHTTP(w, req)
			} else {
				next.ServeHTTP(w, req)
			}
		}
		return http.HandlerFunc(fn)
	})
	// matchbox version
	r.Handle("/", homeHandler())
	// Boot via GRUB
	r.Handle("/grub", s.grubHandler())
	// Boot via iPXE
	r.Handle("/boot.ipxe", ipxeInspect())
	r.Handle("/boot.ipxe.0", ipxeInspect())
	r.Handle("/ipxe", s.ipxeHandler())
	// Ignition Config
	r.Handle("/ignition", s.ignitionHandler())
	// Template
	r.Handle("/template", s.templateHandler())
	// Metadata
	r.Handle("/metadata", s.metadataHandler())

	// kernel, initrd, and TLS assets
	if s.assetsPath != "" {
		r.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(s.assetsPath))))
	}

	return r
}
