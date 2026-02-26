package cmd

import (
	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage CUC mailbox users",
	Long:  `Commands for listing, viewing, creating, updating, and removing CUC voicemail users via CUPI REST API.`,
}

func init() {
	usersCmd.AddCommand(usersListCmd)
	usersCmd.AddCommand(usersGetCmd)
	usersCmd.AddCommand(usersAddCmd)
	usersCmd.AddCommand(usersUpdateCmd)
	usersCmd.AddCommand(usersRemoveCmd)
}
