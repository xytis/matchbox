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
	"github.com/coreos/matchbox/matchbox/sign"
	"github.com/coreos/matchbox/matchbox/storage"
	"github.com/coreos/matchbox/matchbox/tlsutil"

	"github.com/pkg/errors"
	"go.uber.org/zap"
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
	SignatureKeyring    string
	SignaturePassphrase string
}

// NewDaemonConfig returns empty config
func NewDaemonConfig() *DaemonConfig {
	return &DaemonConfig{}
}

// Validate checks if given config is valid
func (c *DaemonConfig) Validate() error {
	if c.HTTPAddress == "" && c.RPCAddress == "" {
		return errors.New("server has no purpose, butter not found")
	}
	if c.TLS {
		if c.TLSKeyFile == "" {
			return errors.New("tls key file not provided")
		}
		if c.TLSCertFile == "" {
			return errors.New("tls cert file not provided")
		}
		if _, err := tlsutil.NewCert(c.TLSCertFile, c.TLSKeyFile); err != nil {
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
	if c.SignatureKeyring != "" {
		if _, err := sign.LoadGPGEntity(c.SignatureKeyring, c.SignaturePassphrase); err != nil {
			return errors.Wrap(err, "failed to create signer")
		}
	}
	if c.AssetsDir != "" {
		if finfo, err := os.Stat(c.AssetsDir); err != nil || !finfo.IsDir() {
			return errors.Errorf("assets path %s is invalid", c.AssetsDir)
		}
	}
	return nil
}

// Daemon groups all underlying components and holds their state
// TODO: it should react to reload and other system events
type Daemon struct {
	core server.Server

	web  *web.Server
	http *http.Server
	rpc  *grpc.Server

	logger *zap.Logger
}

// NewDaemon returns a daemon
func NewDaemon(logger *zap.Logger) *Daemon {
	return &Daemon{logger: logger}
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
	store, err := storage.NewStore(opts.storageConfig, d.logger)
	if err != nil {
		return errors.Wrap(err, "failed to create storage")
	}

	d.core = server.NewServer(store)

	cfg := opts.daemonConfig

	var tc *tls.Config

	if cfg.TLS {
		cert, err := tlsutil.NewCert(cfg.TLSCertFile, cfg.TLSKeyFile)
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

	webcfg := &web.Config{
		Core:   d.core,
		Logger: d.logger,
	}

	if cfg.SignatureKeyring != "" {
		gpg, _ := sign.LoadGPGEntity(cfg.SignatureKeyring, cfg.SignaturePassphrase)
		webcfg.Signer = sign.NewGPGSigner(gpg)
		webcfg.ArmoredSigner = sign.NewArmoredGPGSigner(gpg)
		d.logger.Sugar().Infof("singature: using keyring '%s' (passphrase: %t)", cfg.SignatureKeyring, cfg.SignaturePassphrase != "")
	}

	webcfg.AssetsPath = cfg.AssetsDir

	d.web = web.NewServer(webcfg)

	if cfg.RPCAddress != "" {
		listen, err := net.Listen("tcp", cfg.RPCAddress)
		if err != nil {
			return err
		}
		//TODO: Consider altering listener to include keepalive on connection and other things
		d.logger.Sugar().Infof("rpc: listening on %s", cfg.RPCAddress)
		go d.rpc.Serve(listen)
		defer d.rpc.Stop()
	}

	if cfg.HTTPAddress != "" {
		d.http = &http.Server{Addr: cfg.HTTPAddress, Handler: d.web.HTTPHandler()}
		go func() {
			d.logger.Sugar().Infof("http: listening on %s", cfg.HTTPAddress)
			if err := d.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				d.logger.Error("http: server error", zap.Error(err))
			}
		}()
	}

	done := make(chan bool)
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		logger := d.logger
		// Should handle stuff here
		logger.Info("shutting down")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		d.http.SetKeepAlivesEnabled(false)
		if err := d.http.Shutdown(ctx); err != nil {
			logger.Error("could not gracefully shutdown the server", zap.Error(err))
		}
		close(done)
	}()

	<-done
	return nil
}
