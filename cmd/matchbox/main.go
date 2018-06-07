package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"strings"

	"github.com/coreos/matchbox/matchbox/client"
	"github.com/coreos/matchbox/matchbox/tlsutil"

	"github.com/spf13/cobra"
)

func newCLICommand() *cobra.Command {
	opts := newCLIOptions()

	cmd := &cobra.Command{
		Use:   "matchbox",
		Short: "A command line client for the matchbox service.",
		Long: `A CLI for the matchbox Service

    To get help about a resource or command, run "matchbox help resource"`,
	}

	opts.InstallGlobalFlags(cmd.PersistentFlags())

	return cmd
}

func main() {
	cmd := newCLICommand()
	cmd.SetOutput(os.Stdout)

	cmd.AddCommand(NewVersionCommand())
	cmd.AddCommand(NewGroupCommand())
	cmd.AddCommand(NewProfileCommand())
	cmd.AddCommand(NewTemplateCommand())

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func mustClientFromCmd(cmd *cobra.Command) *client.Client {
	opts := newCLIOptions()
	if err := opts.Parse(cmd.Flags()); err != nil {
		exitWithError(ExitBadArgs, err)
	}

	var tlscfg *tls.Config
	if opts.TLS {
		tlscfg = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: false,
		}
		if opts.TLSServerVerify {
			pool, err := tlsutil.NewCertPool([]string{opts.TLSServerCAFile})
			if err != nil {
				exitWithError(ExitError, err)
			}
			cert, err := tlsutil.NewCert(opts.TLSCertFile, opts.TLSKeyFile)
			if err != nil {
				exitWithError(ExitError, err)
			}
			tlscfg.RootCAs = pool
			tlscfg.GetClientCertificate = tlsutil.StaticClientCertificate(cert)
		}
	}

	cfg := &client.Config{
		Endpoint:    opts.RPCAddress,
		DialTimeout: opts.DialTimeout,
		TLS:         tlscfg,
	}

	// gRPC client
	client, err := client.New(cfg)
	if err != nil {
		exitWithError(ExitBadConnection, err)
	}
	return client
}

func stringSliceToMap(slice []string) map[string]string {
	result := make(map[string]string)
	for _, s := range slice {
		kv := strings.SplitN(s, "=", 2)
		if len(kv) < 2 {
			kv = append(kv, "")
		}
		result[kv[0]] = kv[1]
	}
	return result
}
