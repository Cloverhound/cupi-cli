package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	arDisplayName      string
	arMaxLogonAttempts string
	arLockoutDuration  string
	arQuery            string
)

var authRulesCmd = &cobra.Command{
	Use:   "authrules",
	Short: "Manage authentication rules",
}

var authRulesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List authentication rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		items, err := client.ListAuthRules(srv.Host, srv.Port, cupiUser, cupiPass, arQuery, maxFlag)
		if err != nil {
			return err
		}
		var rows []map[string]string
		for _, v := range items {
			rows = append(rows, map[string]string{
				"objectId":         v.ObjectId,
				"displayName":      v.DisplayName,
				"maxLogonAttempts": v.MaxLogonAttempts,
				"lockoutDuration":  v.LockoutDuration,
			})
		}
		return output.Print(rows, outputFlag)
	},
}

var authRulesGetCmd = &cobra.Command{
	Use:   "get <name-or-id>",
	Short: "Get an authentication rule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		item, err := client.GetAuthRule(srv.Host, srv.Port, cupiUser, cupiPass, args[0])
		if err != nil {
			return err
		}
		return output.Print(item, outputFlag)
	},
}

var authRulesAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an authentication rule",
	RunE: func(cmd *cobra.Command, args []string) error {
		if arDisplayName == "" {
			return fmt.Errorf("--display-name is required")
		}
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{"DisplayName": arDisplayName}
		if arMaxLogonAttempts != "" {
			fields["MaxLogonAttempts"] = arMaxLogonAttempts
		}
		if arLockoutDuration != "" {
			fields["LockoutDuration"] = arLockoutDuration
		}
		if err := client.CreateAuthRule(srv.Host, srv.Port, cupiUser, cupiPass, fields); err != nil {
			return err
		}
		fmt.Printf("Added authentication rule %s\n", arDisplayName)
		return nil
	},
}

var authRulesUpdateCmd = &cobra.Command{
	Use:   "update <name-or-id>",
	Short: "Update an authentication rule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{}
		if arDisplayName != "" {
			fields["DisplayName"] = arDisplayName
		}
		if arMaxLogonAttempts != "" {
			fields["MaxLogonAttempts"] = arMaxLogonAttempts
		}
		if arLockoutDuration != "" {
			fields["LockoutDuration"] = arLockoutDuration
		}
		if len(fields) == 0 {
			return fmt.Errorf("no fields to update; use --help to see available flags")
		}
		if err := client.UpdateAuthRule(srv.Host, srv.Port, cupiUser, cupiPass, args[0], fields); err != nil {
			return err
		}
		fmt.Printf("Updated authentication rule %s\n", args[0])
		return nil
	},
}

var authRulesRemoveCmd = &cobra.Command{
	Use:   "remove <name-or-id>",
	Short: "Remove an authentication rule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		if err := client.DeleteAuthRule(srv.Host, srv.Port, cupiUser, cupiPass, args[0]); err != nil {
			return err
		}
		fmt.Printf("Removed authentication rule %s\n", args[0])
		return nil
	},
}

func init() {
	authRulesListCmd.Flags().StringVar(&arQuery, "query", "", "Filter query")

	authRulesAddCmd.Flags().StringVar(&arDisplayName, "display-name", "", "Display name (required)")
	authRulesAddCmd.Flags().StringVar(&arMaxLogonAttempts, "max-logon-attempts", "", "Max logon attempts")
	authRulesAddCmd.Flags().StringVar(&arLockoutDuration, "lockout-duration", "", "Lockout duration")

	authRulesUpdateCmd.Flags().StringVar(&arDisplayName, "display-name", "", "Display name")
	authRulesUpdateCmd.Flags().StringVar(&arMaxLogonAttempts, "max-logon-attempts", "", "Max logon attempts")
	authRulesUpdateCmd.Flags().StringVar(&arLockoutDuration, "lockout-duration", "", "Lockout duration")

	authRulesCmd.AddCommand(authRulesListCmd)
	authRulesCmd.AddCommand(authRulesGetCmd)
	authRulesCmd.AddCommand(authRulesAddCmd)
	authRulesCmd.AddCommand(authRulesUpdateCmd)
	authRulesCmd.AddCommand(authRulesRemoveCmd)
}
