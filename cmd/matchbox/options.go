package main

import (
	"time"

	"github.com/coreos/matchbox/matchbox/tlsutil"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

type cliOptions struct {
	RPCAddress  string
	DialTimeout time.Duration

	TLS             bool
	TLSKeyFile      string
	TLSCertFile     string
	TLSServerVerify bool
	TLSServerCAFile string
}

func newCLIOptions() *cliOptions {
	return &cliOptions{}
}

// Validate checks if given config is valid
func (o *cliOptions) Validate() error {
	if o.TLS {
		if o.TLSKeyFile == "" {
			return errors.New("tls key file not provided")
		}
		if o.TLSCertFile == "" {
			return errors.New("tls cert file not provided")
		}
		if _, err := tlsutil.NewCert(o.TLSCertFile, o.TLSKeyFile); err != nil {
			return err
		}
		if o.TLSServerVerify {
			if o.TLSServerCAFile == "" {
				return errors.New("server trust CA file not provided")
			}
			if _, err := tlsutil.NewCertPool([]string{o.TLSServerCAFile}); err != nil {
				return err
			}
		}
	}
	return nil
}

// Parse gathers values from flags
func (o *cliOptions) Parse(flags *pflag.FlagSet) (err error) {
	if o.RPCAddress, err = flags.GetString("address"); err != nil {
		return err
	}
	if o.DialTimeout, err = flags.GetDuration("dial-timeout"); err != nil {
		return err
	}

	if o.TLS, err = flags.GetBool("tls"); err != nil {
		return err
	}
	if o.TLSKeyFile, err = flags.GetString("tls-key"); err != nil {
		return err
	}
	if o.TLSCertFile, err = flags.GetString("tls-cert"); err != nil {
		return err
	}

	if o.TLSServerVerify, err = flags.GetBool("tls-verify"); err != nil {
		return err
	}
	if o.TLSServerCAFile, err = flags.GetString("tls-ca"); err != nil {
		return err
	}

	return o.Validate()
}

// InstallGlobalFlags adds flags for the common options on the FlagSet
func (o *cliOptions) InstallGlobalFlags(flags *pflag.FlagSet) {
	flags.String("address", "127.0.0.1:8081", "Server RPC endpoint")
	flags.Duration("dial-timeout", 5*time.Second, "Connection timeout duration")

	flags.Bool("tls", true, "Enable TLS encryption on RPC")
	flags.String("tls-key", "/etc/matchbox/client.key", "Path to TLS key file")
	flags.String("tls-cert", "/etc/matchbox/client.crt", "Path to TLS cert file")

	flags.Bool("tls-verify", true, "Enable server TLS verification")
	flags.String("tls-ca", "/etc/matchbox/ca.pem", "Path to TLS CA file to verify server")
}
