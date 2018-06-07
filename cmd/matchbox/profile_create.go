package main

import (
	"io/ioutil"

	"context"
	"github.com/spf13/cobra"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

// NewProfileCreateCommand creates groups
func NewProfileCreateCommand() *cobra.Command {
	var filename string
	cmd := &cobra.Command{
		Use:   "create --file FILENAME",
		Short: "Create a machine profile",
		Long:  `Create a machine profile`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if len(filename) == 0 {
				cmd.Help()
				return
			}
			client := mustClientFromCmd(cmd)
			profile, err := loadProfile(filename)
			if err != nil {
				exitWithError(ExitError, err)
			}
			req := &pb.ProfilePutRequest{Profile: profile}
			_, err = client.Profiles.ProfilePut(context.TODO(), req)
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

func loadProfile(filename string) (*storagepb.Profile, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return storagepb.ParseProfile(data)
}
