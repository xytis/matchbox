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
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a machine profile",
		Long:  `Create a machine profile`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			client := mustClientFromCmd(cmd)
			flags := cmd.Flags()
			if filename, _ := flags.GetString("filename"); filename != "" {
				profile, err := loadProfile(filename)
				if err != nil {
					exitWithError(ExitError, err)
				}
				req := &pb.ProfilePutRequest{Profile: profile}
				_, err = client.Profiles.ProfilePut(context.TODO(), req)
				if err != nil {
					exitWithError(ExitError, err)
				}
			} else if id, _ := flags.GetString("id"); id != "" {
				name, _ := flags.GetString("name")
				templates, _ := flags.GetStringSlice("templates")
				meta, _ := flags.GetString("metadata")
				profile := &storagepb.Profile{
					Id:       id,
					Name:     name,
					Template: stringSliceToMap(templates),
					Metadata: []byte(meta),
				}
				req := &pb.ProfilePutRequest{Profile: profile}
				_, err := client.Profiles.ProfilePut(context.TODO(), req)
				if err != nil {
					exitWithError(ExitError, err)
				}
			}
		},
	}

	cmd.Flags().StringP("filename", "f", "", "filename to use to create a Profile")
	cmd.MarkFlagFilename("filename", "json")

	cmd.Flags().String("id", "", "Profile ID")
	cmd.Flags().String("name", "", "Profile name")
	cmd.Flags().StringSlice("templates", []string{}, "Template bindings (type=template_id)")
	cmd.Flags().String("metadata", "", "Profile metadata")

	return cmd
}

func loadProfile(filename string) (*storagepb.Profile, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return storagepb.ParseProfile(data)
}
