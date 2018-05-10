package cli

import (
	"io/ioutil"
	"path/filepath"

	"context"
	"github.com/spf13/cobra"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

// templateImportCmd creates template from a raw file
var (
	templateImportCmd = &cobra.Command{
		Use:   "import --id ID --filename FILENAME",
		Short: "Import a template",
		Long:  `Import a file as a template`,
		Run:   runTemplateImportCmd,
	}
)

func init() {
	templateCmd.AddCommand(templateImportCmd)
	templateImportCmd.Flags().StringVarP(&flagFilename, "filename", "f", "", "filename to use to create a template")
	templateImportCmd.Flags().StringVarP(&flagID, "id", "i", "", "id to use for template")
	templateImportCmd.MarkFlagRequired("id")
	templateImportCmd.MarkFlagRequired("filename")
}

func runTemplateImportCmd(cmd *cobra.Command, args []string) {
	if len(flagID) == 0 {
		cmd.Help()
		return
	}
	if len(flagFilename) == 0 {
		cmd.Help()
		return
	}
	if err := validateArgs(cmd, args); err != nil {
		return
	}

	client := mustClientFromCmd(cmd)
	data, err := ioutil.ReadFile(flagFilename)
	if err != nil {
		exitWithError(ExitError, err)
	}
	template := &storagepb.Template{}
	template.Id = flagID
	template.Name = filepath.Base(flagFilename)
	template.Contents = data
	if err = template.AssertValid(); err != nil {
		exitWithError(ExitError, err)
	}
	req := &pb.TemplatePutRequest{Template: template}
	_, err = client.Templates.TemplatePut(context.TODO(), req)
	if err != nil {
		exitWithError(ExitError, err)
	}
}
