package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"context"
	"github.com/spf13/cobra"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

// NewTemplateImportCommand creates groups
func NewTemplateImportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import --id ID --file FILENAME",
		Short: "Import a template",
		Long:  `Import a template from raw file`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			client := mustClientFromCmd(cmd)
			flags := cmd.Flags()

			template := &storagepb.Template{}
			if id, err := flags.GetString("id"); err == nil {
				template.Id = id
			} else {
				exitWithError(ExitError, err)
			}

			if filename, err := flags.GetString("filename"); err == nil && filename != "-" && filename != "" {
				if name, _ := flags.GetString("name"); name == "" {
					template.Name = filepath.Base(filename)
				} else {
					template.Name = name
				}
				data, err := ioutil.ReadFile(filename)
				if err != nil {
					exitWithError(ExitError, err)
				}
				template.Contents = data
			} else if err == nil {
				data, err := ioutil.ReadAll(os.Stdin)
				if err != nil {
					exitWithError(ExitError, err)
				}
				template.Contents = data
			} else {
				exitWithError(ExitError, err)
			}

			if err := template.AssertValid(); err != nil {
				exitWithError(ExitError, err)
			}
			req := &pb.TemplatePutRequest{Template: template}
			_, err := client.Templates.TemplatePut(context.TODO(), req)
			if err != nil {
				exitWithError(ExitError, err)
			}
		},
	}

	cmd.Flags().StringP("id", "i", "", "id to use for template")
	cmd.Flags().String("name", "", "name to use for template")
	cmd.Flags().StringP("filename", "f", "", "filename to template contents")
	cmd.MarkFlagRequired("id")

	return cmd
}
