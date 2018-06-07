package main

import (
	"github.com/spf13/cobra"
)

// NewGroupCommand creates group command tree
func NewGroupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
		Short: "Manage machine groups",
		Long:  `List and describe machine groups`,
	}
	cmd.AddCommand(NewGroupCreateCommand())
	cmd.AddCommand(NewGroupDescribeCommand())
	cmd.AddCommand(NewGroupListCommand())

	return cmd
}
