package main

import (
	"fmt"
	"os"

	"context"
	"github.com/spf13/cobra"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
)

// NewGroupDescribeCommand describes groups
func NewGroupDescribeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe GROUP_ID",
		Short: "Describe a machine group",
		Long:  `Describe a machine group`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.Help()
				return
			}

			tw := newTabWriter(os.Stdout)
			defer tw.Flush()

			// legend
			fmt.Fprintf(tw, "ID\tNAME\tSELECTORS\tPROFILE\tMETADATA\n")

			client := mustClientFromCmd(cmd)
			request := &pb.GroupGetRequest{
				Id: args[0],
			}
			resp, err := client.Groups.GroupGet(context.TODO(), request)
			if err != nil {
				return
			}
			g := resp.Group
			fmt.Fprintf(tw, "%s\t%s\t%s\t%#v\t%s\n", g.Id, g.Name, g.Selector, g.Profile, g.Metadata)
		},
	}

	return cmd
}
