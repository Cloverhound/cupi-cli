package cmd

import (
	"github.com/spf13/cobra"
)

var handlersCmd = &cobra.Command{
	Use:   "handlers",
	Short: "Manage CUC call handlers",
	Long:  `Commands for listing, viewing, creating, updating, and removing CUC call handlers.`,
}

func init() {
	handlersCmd.AddCommand(handlersListCmd)
	handlersCmd.AddCommand(handlersGetCmd)
	handlersCmd.AddCommand(handlersAddCmd)
	handlersCmd.AddCommand(handlersUpdateCmd)
	handlersCmd.AddCommand(handlersRemoveCmd)
}
