package cmd

import (
	"github.com/spf13/cobra"
)

var distlistsCmd = &cobra.Command{
	Use:   "distlists",
	Short: "Manage CUC distribution lists",
	Long:  `Commands for listing, viewing, creating, updating, and removing CUC distribution lists.`,
}

func init() {
	distlistsCmd.AddCommand(distlistsListCmd)
	distlistsCmd.AddCommand(distlistsGetCmd)
	distlistsCmd.AddCommand(distlistsAddCmd)
	distlistsCmd.AddCommand(distlistsUpdateCmd)
	distlistsCmd.AddCommand(distlistsRemoveCmd)
	distlistsCmd.AddCommand(distlistsMembersCmd)
}
