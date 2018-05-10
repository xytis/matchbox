package cli

import (
	"io/ioutil"

	"context"
	"github.com/spf13/cobra"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

// templatePutCmd creates and updates templates.
var (
	templatePutCmd = &cobra.Command{
		Use:   "create --file FILENAME",
		Short: "Create a template",
		Long:  `Create a template`,
		Run:   runTemplatePutCmd,
	}
)

func init() {
	templateCmd.AddCommand(templatePutCmd)
	templatePutCmd.Flags().StringVarP(&flagFilename, "filename", "f", "", "filename to use to create a template")
	templatePutCmd.MarkFlagRequired("filename")
}

func runTemplatePutCmd(cmd *cobra.Command, args []string) {
	if len(flagFilename) == 0 {
		cmd.Help()
		return
	}
	if err := validateArgs(cmd, args); err != nil {
		return
	}

	client := mustClientFromCmd(cmd)
	template, err := loadTemplate(flagFilename)
	if err != nil {
		exitWithError(ExitError, err)
	}
	req := &pb.TemplatePutRequest{Template: template}
	_, err = client.Templates.TemplatePut(context.TODO(), req)
	if err != nil {
		exitWithError(ExitError, err)
	}
}

func loadTemplate(filename string) (*storagepb.Template, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return storagepb.ParseTemplate(data)
}
