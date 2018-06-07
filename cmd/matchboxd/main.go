package main

import (
	"fmt"
	"os"

	"github.com/coreos/matchbox/matchbox/cli"
	"github.com/coreos/matchbox/matchbox/version"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type pLevel struct {
	zapcore.Level
}

func (*pLevel) Type() string {
	return "string"
}

func newDaemonCommand() *cobra.Command {
	var flags *pflag.FlagSet
	var cfg *viper.Viper
	level := &pLevel{0}
	opts := newDaemonOptions()

	cmd := &cobra.Command{
		Use:           "matchboxd [OPTIONS]",
		Short:         "Provides fire to your boots",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cli.NoArgsRequired,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := zap.New(zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				zapcore.Lock(os.Stdout),
				zap.NewAtomicLevelAt(level.Level),
			))

			if err := opts.ExtractConfig(cfg); err != nil {
				return err
			}
			d := NewDaemon(logger)
			return d.start(opts)
		},
		DisableFlagsInUseLine: true,
		Version:               fmt.Sprintf("%s, build %s", version.Version, version.GitCommit),
	}

	flags = cmd.Flags()
	cfg = viper.New()

	cmd.PersistentFlags().BoolP("help", "h", false, "Print usage")

	flags.VarP(level, "log-level", "l", `Set the logging level ("debug"|"info"|"warn"|"error"|"fatal")`)

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
