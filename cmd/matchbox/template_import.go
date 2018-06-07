package main

import (
	"io/ioutil"
	"path/filepath"

	"context"
	"github.com/spf13/cobra"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

// NewTemplateCreateCommand creates groups
func NewTemplateImportCommand() *cobra.Command {
	var filename string
	var templateID string
	cmd := &cobra.Command{
		Use:   "create --file FILENAME",
		Short: "Create a template",
		Long:  `Create a template`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			client := mustClientFromCmd(cmd)
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				exitWithError(ExitError, err)
			}
			template := &storagepb.Template{}
			template.Id = templateID
			template.Name = filepath.Base(filename)
			template.Contents = data
			if err = template.AssertValid(); err != nil {
				exitWithError(ExitError, err)
			}
			req := &pb.TemplatePutRequest{Template: template}
			_, err = client.Templates.TemplatePut(context.TODO(), req)
			if err != nil {
				exitWithError(ExitError, err)
			}
		},
	}

	cmd.Flags().StringVarP(&templateID, "id", "i", "", "id to use for template")
	cmd.Flags().StringVarP(&filename, "filename", "f", "", "filename to template contents")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("filename")

	return cmd
}
