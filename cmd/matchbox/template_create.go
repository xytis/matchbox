package main

import (
	"io/ioutil"

	"context"
	"github.com/spf13/cobra"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

// NewTemplateCreateCommand creates groups
func NewTemplateCreateCommand() *cobra.Command {
	var filename string
	cmd := &cobra.Command{
		Use:   "create --file FILENAME",
		Short: "Create a template",
		Long:  `Create a template`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if len(filename) == 0 {
				cmd.Help()
				return
			}
			client := mustClientFromCmd(cmd)
			template, err := loadTemplate(filename)
			if err != nil {
				exitWithError(ExitError, err)
			}
			req := &pb.TemplatePutRequest{Template: template}
			_, err = client.Templates.TemplatePut(context.TODO(), req)
			if err != nil {
				exitWithError(ExitError, err)
			}
		},
	}

	cmd.Flags().StringVarP(&filename, "filename", "f", "", "filename to use to create a Group")
	cmd.MarkFlagRequired("filename")
	cmd.MarkFlagFilename("filename", "json")

	return cmd
}

func loadTemplate(filename string) (*storagepb.Template, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return storagepb.ParseTemplate(data)
}
