package config

import (
	"crypto/tls"

	"github.com/coreos/matchbox/matchbox/tlsutil"

	"github.com/pkg/errors"
)

// EtcdStoreConfig initializes a etcdStore.
type EtcdStoreConfig struct {
	Endpoints             []string
	Password              string
	Username              string
	Prefix                string
	TLS                   bool
	TLSCAFile             string
	TLSCertFile           string
	TLSInsecureSkipVerify bool
	TLSKeyFile            string
	TLSServerName         string
}

// NewEtcdStoreConfig creates an empty config
func NewEtcdStoreConfig() *EtcdStoreConfig {
	return &EtcdStoreConfig{}
}

// Validate checks if given config is viable
func (c *EtcdStoreConfig) Validate() error {
	if len(c.Endpoints) == 0 {
		return errors.New("missing etcd endpoints slice")
	}

	if c.TLS {
		if _, err := tlsutil.NewCert(c.TLSCertFile, c.TLSKeyFile, nil); err != nil {
			return err
		}
		if c.TLSCAFile != "" {
			if _, err := tlsutil.NewCertPool([]string{c.TLSCAFile}); err != nil {
				return err
			}
		}
	}

	return nil
}

// BuildTLSConfig creates tls.Config for Etcd Client
func (c *EtcdStoreConfig) BuildTLSConfig() (*tls.Config, error) {
	pool, err := tlsutil.NewCertPool([]string{c.TLSCAFile})
	if err != nil {
		return nil, err
	}

	cert, err := tlsutil.NewCert(c.TLSCertFile, c.TLSKeyFile, nil)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		GetClientCertificate: tlsutil.StaticClientCertificate(cert),
		InsecureSkipVerify:   c.TLSInsecureSkipVerify,
		MinVersion:           tls.VersionTLS12,
		RootCAs:              pool,
		ServerName:           c.TLSServerName,
	}, nil
}
