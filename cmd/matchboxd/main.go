package main

import (
	"fmt"
	"os"

	"github.com/coreos/matchbox/matchbox/cli"
	"github.com/coreos/matchbox/matchbox/version"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	// Defaults to info logging
	log = logrus.New()
)

func newDaemonCommand() *cobra.Command {
	var flags *pflag.FlagSet
	var cfg *viper.Viper
	opts := newDaemonOptions(log)

	cmd := &cobra.Command{
		Use:           "matchboxd [OPTIONS]",
		Short:         "Provides fire to your boots",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cli.NoArgsRequired,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ExtractConfig(cfg)
			return runDaemon(opts)
		},
		DisableFlagsInUseLine: true,
		Version:               fmt.Sprintf("%s, build %s", version.Version, version.GitCommit),
	}

	flags = cmd.Flags()
	cfg = viper.New()

	cmd.PersistentFlags().BoolP("help", "h", false, "Print usage")

	flags.BoolP("version", "v", false, "Print version information and quit")
	flags.String("config", "/etc/matchbox/matchboxd.yaml", "Matchboxd configuration file")

	opts.InstallFlags(flags)

	cfg.SetConfigName("matchboxd")
	cfg.SetEnvPrefix("matchboxd")
	cfg.AutomaticEnv()
	cfg.BindPFlags(flags)
	cfg.SetConfigFile(cfg.GetString("config"))

	return cmd
}

func main() {
	cmd := newDaemonCommand()
	cmd.SetOutput(os.Stdout)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
