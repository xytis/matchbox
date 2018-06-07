package main

import (
	"io/ioutil"

	"context"
	"github.com/spf13/cobra"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

// NewGroupCreateCommand creates groups
func NewGroupCreateCommand() *cobra.Command {
	var filename string
	cmd := &cobra.Command{
		Use:   "create --file FILENAME",
		Short: "Create a machine group",
		Long:  `Create a machine group`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if len(filename) == 0 {
				cmd.Help()
				return
			}
			client := mustClientFromCmd(cmd)
			group, err := loadGroup(filename)
			if err != nil {
				exitWithError(ExitError, err)
			}
			req := &pb.GroupPutRequest{Group: group}
			_, err = client.Groups.GroupPut(context.TODO(), req)
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

func loadGroup(filename string) (*storagepb.Group, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return storagepb.ParseGroup(data)
}
