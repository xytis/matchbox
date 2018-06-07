package main

import (
	"context"
	"fmt"
	"runtime"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
	"github.com/coreos/matchbox/matchbox/version"

	"github.com/spf13/cobra"
)

// NewVersionCommand builds version command for client
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version and exit",
		Long:  `Print the version of the bootcmd client`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("runtime: go: %s os: %s arch: %s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
			fmt.Printf("client: version: %s\n", version.Version)
			client := mustClientFromCmd(cmd)
			req := &pb.VersionReportRequest{}
			if resp, err := client.Version.VersionReport(context.TODO(), req); err != nil {
				fmt.Printf("server could not be reached %v\n", err)
			} else {
				fmt.Printf("server: version: %s\n", resp.Version)
			}
		},
	}
}
