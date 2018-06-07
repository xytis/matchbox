package main

import (
	"fmt"
	"os"

	"context"
	"github.com/spf13/cobra"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
)

// NewProfileListCommand creates groups
func NewProfileListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List machine profiles",
		Long:  `List machine profiles`,
		Run: func(cmd *cobra.Command, args []string) {

			tw := newTabWriter(os.Stdout)
			defer tw.Flush()
			// legend
			fmt.Fprintf(tw, "ID\tPROFILE NAME\t")

			client := mustClientFromCmd(cmd)
			resp, err := client.Profiles.ProfileList(context.TODO(), &pb.ProfileListRequest{})
			if err != nil {
				return
			}
			for _, profile := range resp.Profiles {
				fmt.Fprintf(tw, "%s\t%s\n", profile.Id, profile.Name)
			}
		},
	}

	return cmd
}
