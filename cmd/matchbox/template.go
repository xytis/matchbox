package main

import (
	"github.com/spf13/cobra"
)

// NewTemplateCommand creates template command tree
func NewTemplateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Manage Templates",
		Long:  `Manage Templates`,
	}
	cmd.AddCommand(NewTemplateImportCommand())
	cmd.AddCommand(NewTemplateCreateCommand())
	//cmd.AddCommand(NewTemplateDescribeCommand())
	//cmd.AddCommand(NewTemplateListCommand())

	return cmd
}
