package main

import (
	"fmt"
	"os"

	"context"
	"github.com/spf13/cobra"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
)

// NewGroupListCommand describes groups
func NewGroupListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List machine groups",
		Long:  `List machine groups`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			tw := newTabWriter(os.Stdout)
			defer tw.Flush()
			// legend
			fmt.Fprintf(tw, "ID\tGROUP NAME\tSELECTORS\tPROFILE\n")

			client := mustClientFromCmd(cmd)
			resp, err := client.Groups.GroupList(context.TODO(), &pb.GroupListRequest{})
			if err != nil {
				return
			}
			for _, group := range resp.Groups {
				fmt.Fprintf(tw, "%s\t%s\t%#v\t%s\n", group.Id, group.Name, group.Selector, group.Profile)
			}
		},
	}

	return cmd
}
