package main

import (
	"fmt"
	"os"

	"context"
	"github.com/spf13/cobra"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
)

// NewTemplateDescribeCommand creates groups
func NewTemplateDescribeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe TEMPLATE_ID",
		Short: "Describe a template",
		Long:  `Describe a template`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := mustClientFromCmd(cmd)
			request := &pb.TemplateGetRequest{
				Id: args[0],
			}
			resp, err := client.Templates.TemplateGet(context.TODO(), request)
			if err != nil {
				return
			}
			p := resp.Template
			fmt.Fprintf(os.Stdout, "id: %s\nname: %s\n---\n%v\n", p.Id, p.Name, string(p.Contents))
		},
	}

	return cmd
}
