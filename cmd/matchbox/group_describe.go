package main

import (
	"encoding/json"
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
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tw := newTabWriter(os.Stdout)
			defer tw.Flush()

			client := mustClientFromCmd(cmd)
			flags := cmd.Flags()
			request := &pb.GroupGetRequest{
				Id: args[0],
			}
			resp, err := client.Groups.GroupGet(context.TODO(), request)
			if err != nil {
				return
			}
			g := resp.Group
			if ok, _ := flags.GetBool("json"); !ok {
				fmt.Fprintf(tw, "ID:\t%s\n", g.Id)
				fmt.Fprintf(tw, "Name:\t%s\n", g.Name)
				fmt.Fprintf(tw, "Selectors:\t%s\n", g.SelectorString())
				fmt.Fprintf(tw, "Profile:\t%s\n", g.Profile)
				fmt.Fprintf(tw, "Metadata:\n%s\n", g.MetadataPrettyString())
			} else {
				j, err := json.Marshal(g)
				if err != nil {
					fmt.Fprintf(os.Stdout, `{"success": "false", "error": "%v"}`, err)
				}
				fmt.Fprint(os.Stdout, string(j))
			}
		},
	}

	cmd.Flags().Bool("json", false, "Output JSON")

	return cmd
}
