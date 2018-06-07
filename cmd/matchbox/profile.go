package main

import (
	"github.com/spf13/cobra"
)

// NewProfileCommand creates profile command tree
func NewProfileCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage machine profiles",
		Long:  `List and describe machine profiles`,
	}
	cmd.AddCommand(NewProfileCreateCommand())
	cmd.AddCommand(NewProfileDescribeCommand())
	cmd.AddCommand(NewProfileListCommand())

	return cmd
}
