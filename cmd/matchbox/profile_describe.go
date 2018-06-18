package main

import (
	"encoding/json"
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

			client := mustClientFromCmd(cmd)
			flags := cmd.Flags()
			request := &pb.ProfileGetRequest{
				Id: args[0],
			}
			resp, err := client.Profiles.ProfileGet(context.TODO(), request)
			if err != nil {
				return
			}
			p := resp.Profile
			if ok, _ := flags.GetBool("json"); !ok {
				fmt.Fprintf(tw, "ID:\t%s\n", p.Id)
				fmt.Fprintf(tw, "Name:\t%s\n", p.Name)
				fmt.Fprintf(tw, "Templates:\t%s\n", p.TemplateString())
				fmt.Fprintf(tw, "Metadata:\n%s\n", p.MetadataPrettyString())
			} else {
				j, err := json.Marshal(p)
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
