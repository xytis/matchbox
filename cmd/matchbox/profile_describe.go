package main

import (
	"fmt"
	"os"

	"context"
	"github.com/spf13/cobra"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
)

// NewProfileDescribeCommand creates groups
func NewProfileDescribeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe PROFILE_ID",
		Short: "Describe a machine profile",
		Long:  `Describe a machine profile`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tw := newTabWriter(os.Stdout)
			defer tw.Flush()
			// legend
			fmt.Fprintf(tw, "ID\tNAME\tBINDINGS\tMETADATA\n")

			client := mustClientFromCmd(cmd)
			request := &pb.ProfileGetRequest{
				Id: args[0],
			}
			resp, err := client.Profiles.ProfileGet(context.TODO(), request)
			if err != nil {
				return
			}
			p := resp.Profile
			fmt.Fprintf(tw, "%s\t%s\t%v\t%v\n", p.Id, p.Name, p.Template, p.Metadata)
		},
	}

	return cmd
}
