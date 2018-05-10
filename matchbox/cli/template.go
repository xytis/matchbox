package cli

import (
	"github.com/spf13/cobra"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage Templates",
	Long:  `Manage Templates`,
}

func init() {
	RootCmd.AddCommand(templateCmd)
}
