package main

import (
	"fmt"
	"os"

	"context"
	"github.com/spf13/cobra"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
)

// NewTemplateListCommand creates groups
func NewTemplateListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List templates",
		Long:  `List templates`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {

			tw := newTabWriter(os.Stdout)
			defer tw.Flush()
			// legend
			fmt.Fprintf(tw, "ID\tTEMPLATE NAME\t\n")

			client := mustClientFromCmd(cmd)
			resp, err := client.Templates.TemplateList(context.TODO(), &pb.TemplateListRequest{})
			if err != nil {
				return
			}
			for _, template := range resp.Templates {
				fmt.Fprintf(tw, "%s\t%s\n", template.Id, template.Name)
			}
		},
	}

	return cmd
}
