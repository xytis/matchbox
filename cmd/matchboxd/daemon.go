package main

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	web "github.com/coreos/matchbox/matchbox/http"
	"github.com/coreos/matchbox/matchbox/rpc"
	"github.com/coreos/matchbox/matchbox/server"
	"github.com/coreos/matchbox/matchbox/tlsutil"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// DaemonConfig holds information about daemon settings
type DaemonConfig struct {
	HTTPAddress string
	RPCAddress  string

	//RPC config
	TLS             bool
	TLSKeyFile      string
	TLSCertFile     string
	TLSClientVerify bool
	TLSClientCAFile string

	//HTTP Config
	AssetsDir string

	//Signature Config
	SignatureKeyring   string
	SignaturePassphase string
}

// NewDaemonConfig returns empty config
func NewDaemonConfig() *DaemonConfig {
	return &DaemonConfig{}
}

// Validate checks if given config is valid
func (c *DaemonConfig) Validate() error {
	if c.TLS {
		if c.TLSKeyFile == "" {
			return errors.New("tls key file not provided")
		}
		if c.TLSCertFile == "" {
			return errors.New("tls cert file not provided")
		}
		if _, err := tlsutil.NewCert(c.TLSCertFile, c.TLSKeyFile, nil); err != nil {
			return err
		}
		if c.TLSClientVerify {
			if c.TLSClientCAFile == "" {
				return errors.New("client trust CA file not provided")
			}
			if _, err := tlsutil.NewCertPool([]string{c.TLSClientCAFile}); err != nil {
				return err
			}
		}
	}
	if c.AssetsDir != "" {
		if finfo, err := os.Stat(c.AssetsDir); err != nil || !finfo.IsDir() {
			return errors.Errorf("assets path %s is invalid", c.AssetsDir)
		}
	}
	return nil
}

func runDaemon(opts *daemonOptions) error {
	daemon := NewDaemon()
	return daemon.start(opts)
}

// Daemon groups all underlying components and holds their state
// TODO: it should react to reload and other system events
type Daemon struct {
	core server.Server
	web  *web.Server
	http *http.Server
	rpc  *grpc.Server
}

// NewDaemon returns a daemon
func NewDaemon() *Daemon {
	return &Daemon{}
}

type keepalive struct {
	*net.TCPListener
}

func (ln keepalive) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func (d *Daemon) start(opts *daemonOptions) error {
	d.core = server.NewServer(opts.serverConfig)

	cfg := opts.daemonConfig

	var tc *tls.Config

	if cfg.TLS {
		cert, err := tlsutil.NewCert(cfg.TLSCertFile, cfg.TLSKeyFile, nil)
		if err != nil {
			return errors.Wrap(err, "certificate")
		}

		pool, err := tlsutil.NewCertPool([]string{cfg.TLSClientCAFile})
		if err != nil {
			return err
		}
		clientAuthType := tls.RequireAndVerifyClientCert
		if !cfg.TLSClientVerify {
			clientAuthType = tls.RequestClientCert
		}
		tc = &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
			GetCertificate:           tlsutil.StaticServerCertificate(cert),
			ClientAuth:               clientAuthType,
			ClientCAs:                pool,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			},
		}
	}

	d.rpc = rpc.NewServer(&rpc.Config{
		Core: d.core,
		TLS:  tc,
	})

	d.web = web.NewServer(&web.Config{
		Core:   d.core,
		Logger: opts.logger,
	})

	if cfg.RPCAddress != "" {
		listen, err := net.Listen("tcp", cfg.RPCAddress)
		if err != nil {
			return err
		}
		opts.logger.Infof("rpc: listening on %s", cfg.RPCAddress)
		go d.rpc.Serve(listen)
		defer d.rpc.Stop()
	}

	if cfg.HTTPAddress != "" {
		d.http = &http.Server{Addr: cfg.HTTPAddress, Handler: d.web.HTTPHandler()}
		go func() {
			opts.logger.Infof("http: listening on %s", cfg.HTTPAddress)
			if err := d.http.ListenAndServe(); err != nil {
				opts.logger.Fatal(err)
			}
		}()
	}

	done := make(chan bool, 1)
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		// Should handle stuff here
		opts.logger.Infoln("shutting down")
		d.http.Shutdown(context.Background())
		done <- true
	}()

	<-done
	return nil
}
