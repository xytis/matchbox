package main

import (
	"github.com/coreos/matchbox/matchbox/server/config"
	// "github.com/docker/go-connections/tlsconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type daemonOptions struct {
	version bool

	logger *logrus.Logger

	serverConfig *config.Config
	daemonConfig *DaemonConfig
}

// newDaemonOptions returns a new daemonOptions
func newDaemonOptions(logger *logrus.Logger) *daemonOptions {
	return &daemonOptions{
		logger:       logger,
		daemonConfig: NewDaemonConfig(),
		serverConfig: config.NewConfig(logger),
	}
}

// InstallFlags adds flags for the common options on the FlagSet
func (o *daemonOptions) InstallFlags(flags *pflag.FlagSet) {
	// Daemon flags
	flags.StringP("log-level", "l", "info", `Set the logging level ("debug"|"info"|"warn"|"error"|"fatal")`)
	flags.String("http-address", "127.0.0.1:8080", "HTTP listen address")
	flags.String("rpc-address", "", "RPC listen address")

	flags.Bool("tls", true, "Enable TLS encryption on RPC")
	flags.String("tls-key", "", "Path to TLS key file")
	flags.String("tls-cert", "", "Path to TLS cert file")

	flags.Bool("tls-verify", true, "Enable TLS verification for RPC clients")
	flags.String("tls-ca", "", "Path to TLS CA file to verify clients")

	//Storage flags
	flags.String("store-backend", "filesystem", `Select storage backend to use ("filesystem|etcd")`)

	flags.String("store-filesystem-root", "/var/lib/matchbox", `Root catalog for filesystem storage`)

	flags.StringSlice("store-etcd-endpoints", []string{"127.0.0.1:6237"}, `Etcd endpoints for connecting`)
	flags.String("store-etcd-password", "", `Etcd password`)
	flags.String("store-etcd-username", "", `Etcd username`)
	flags.String("store-etcd-prefix", "", `Etcd prefix to use`)
	flags.Bool("store-etcd-tls", false, `Etcd use TLS`)
	flags.String("store-etcd-tls-ca", "", `Etcd server CA to trust`)
	flags.String("store-etcd-tls-cert", "", `Etcd cert file`)
	flags.Bool("store-etcd-tls-insecure-skip-verify", false, `Etcd skip verification`)
	flags.String("store-etcd-tls-key", "", `Etcd key file`)
	flags.String("store-etcd-tls-server-name", "", `Etcd CN to verify server against`)

	//Signature flags
	flags.String("signature-keyring", "", `Path to a private keyring file (use ENV or config file to provide passphrase to it)`)
}

// ExtractConfig aquires required values from configuration
func (o *daemonOptions) ExtractConfig(cfg *viper.Viper) {
	// logging setup
	lvl, err := logrus.ParseLevel(cfg.GetString("log-level"))
	if err != nil {
		o.logger.Fatalf("invalid log-level: %v", err)
	}
	o.logger.Level = lvl

	o.daemonConfig.HTTPAddress = cfg.GetString("http-address")
	o.daemonConfig.RPCAddress = cfg.GetString("rpc-address")
	o.daemonConfig.TLS = cfg.GetBool("tls")
	o.daemonConfig.TLSKeyFile = cfg.GetString("tls-key")
	o.daemonConfig.TLSCertFile = cfg.GetString("tls-cert")
	o.daemonConfig.TLSClientVerify = cfg.GetBool("tls-verify")
	o.daemonConfig.TLSClientCAFile = cfg.GetString("tls-ca")

	o.daemonConfig.SignatureKeyring = cfg.GetString("signature-keyring")
	o.daemonConfig.SignaturePassphase = cfg.GetString("signature-passphrase")

	if err := o.daemonConfig.Validate(); err != nil {
		o.logger.Fatalf("invalid configuration: %v", err)
	}

	o.serverConfig.StoreBackend = cfg.GetString("store-backend")

	o.serverConfig.FileStoreConfig.Root = cfg.GetString("store-filesystem-root")

	o.serverConfig.EtcdStoreConfig.Endpoints = cfg.GetStringSlice("store-etcd-endpoints")
	o.serverConfig.EtcdStoreConfig.Username = cfg.GetString("store-etcd-username")
	o.serverConfig.EtcdStoreConfig.Password = cfg.GetString("store-etcd-password")
	o.serverConfig.EtcdStoreConfig.Prefix = cfg.GetString("store-etcd-prefix")
	o.serverConfig.EtcdStoreConfig.TLS = cfg.GetBool("store-etcd-tls")
	o.serverConfig.EtcdStoreConfig.TLSCAFile = cfg.GetString("store-etcd-tls-ca")
	o.serverConfig.EtcdStoreConfig.TLSCertFile = cfg.GetString("store-etcd-tls-cert")
	o.serverConfig.EtcdStoreConfig.TLSInsecureSkipVerify = cfg.GetBool("store-etcd-tls-insecure-skip-verify")
	o.serverConfig.EtcdStoreConfig.TLSKeyFile = cfg.GetString("store-etcd-tls-key")
	o.serverConfig.EtcdStoreConfig.TLSServerName = cfg.GetString("store-etcd-tls-server-name")

	if err := o.serverConfig.Validate(); err != nil {
		o.logger.Fatalf("invalid configuration: %v", err)
	}
}
