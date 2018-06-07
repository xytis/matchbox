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
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a machine group",
		Long:  `Create a machine group`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			client := mustClientFromCmd(cmd)
			flags := cmd.Flags()
			if filename, _ := flags.GetString("filename"); filename != "" {
				group, err := loadGroup(filename)
				if err != nil {
					exitWithError(ExitError, err)
				}
				req := &pb.GroupPutRequest{Group: group}
				_, err = client.Groups.GroupPut(context.TODO(), req)
				if err != nil {
					exitWithError(ExitError, err)
				}
			} else if id, _ := flags.GetString("id"); id != "" {
				name, _ := flags.GetString("name")
				profile, _ := flags.GetString("profile")
				selectors, _ := flags.GetStringSlice("selectors")
				meta, _ := flags.GetString("metadata")
				group := &storagepb.Group{
					Id:       id,
					Name:     name,
					Profile:  profile,
					Selector: stringSliceToMap(selectors),
					Metadata: []byte(meta),
				}
				req := &pb.GroupPutRequest{Group: group}
				_, err := client.Groups.GroupPut(context.TODO(), req)
				if err != nil {
					exitWithError(ExitError, err)
				}
			} else {
				cmd.Help()
			}
		},
	}

	cmd.Flags().StringP("filename", "f", "", "filename to use to create a Group")
	cmd.MarkFlagFilename("filename", "json")

	cmd.Flags().String("id", "", "Group ID")
	cmd.Flags().String("name", "", "Group name")
	cmd.Flags().String("profile", "", "Profile asigned to group")
	cmd.Flags().StringSlice("selectors", []string{}, "Selectors for group (key=value)")
	cmd.Flags().String("metadata", "", "Group metadata")

	return cmd
}

func loadGroup(filename string) (*storagepb.Group, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return storagepb.ParseGroup(data)
}
