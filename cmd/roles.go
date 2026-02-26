package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	roleDisplayName string
	roleDescription string
	roleQuery       string
)

var rolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "Manage custom roles",
}

var rolesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List custom roles",
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		items, err := client.ListCustomRoles(srv.Host, srv.Port, cupiUser, cupiPass, roleQuery, maxFlag)
		if err != nil {
			return err
		}
		var rows []map[string]string
		for _, v := range items {
			rows = append(rows, map[string]string{
				"objectId":    v.ObjectId,
				"displayName": v.DisplayName,
				"description": v.Description,
			})
		}
		return output.Print(rows, outputFlag)
	},
}

var rolesGetCmd = &cobra.Command{
	Use:   "get <name-or-id>",
	Short: "Get a custom role",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		item, err := client.GetCustomRole(srv.Host, srv.Port, cupiUser, cupiPass, args[0])
		if err != nil {
			return err
		}
		return output.Print(item, outputFlag)
	},
}

var rolesAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a custom role",
	RunE: func(cmd *cobra.Command, args []string) error {
		if roleDisplayName == "" {
			return fmt.Errorf("--display-name is required")
		}
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{"DisplayName": roleDisplayName}
		if roleDescription != "" {
			fields["Description"] = roleDescription
		}
		if err := client.CreateCustomRole(srv.Host, srv.Port, cupiUser, cupiPass, fields); err != nil {
			return err
		}
		fmt.Printf("Added custom role %s\n", roleDisplayName)
		return nil
	},
}

var rolesUpdateCmd = &cobra.Command{
	Use:   "update <name-or-id>",
	Short: "Update a custom role",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{}
		if roleDisplayName != "" {
			fields["DisplayName"] = roleDisplayName
		}
		if roleDescription != "" {
			fields["Description"] = roleDescription
		}
		if len(fields) == 0 {
			return fmt.Errorf("no fields to update; use --help to see available flags")
		}
		if err := client.UpdateCustomRole(srv.Host, srv.Port, cupiUser, cupiPass, args[0], fields); err != nil {
			return err
		}
		fmt.Printf("Updated custom role %s\n", args[0])
		return nil
	},
}

var rolesRemoveCmd = &cobra.Command{
	Use:   "remove <name-or-id>",
	Short: "Remove a custom role",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		if err := client.DeleteCustomRole(srv.Host, srv.Port, cupiUser, cupiPass, args[0]); err != nil {
			return err
		}
		fmt.Printf("Removed custom role %s\n", args[0])
		return nil
	},
}

func init() {
	rolesListCmd.Flags().StringVar(&roleQuery, "query", "", "Filter query")

	rolesAddCmd.Flags().StringVar(&roleDisplayName, "display-name", "", "Display name (required)")
	rolesAddCmd.Flags().StringVar(&roleDescription, "description", "", "Description")

	rolesUpdateCmd.Flags().StringVar(&roleDisplayName, "display-name", "", "Display name")
	rolesUpdateCmd.Flags().StringVar(&roleDescription, "description", "", "Description")

	rolesCmd.AddCommand(rolesListCmd)
	rolesCmd.AddCommand(rolesGetCmd)
	rolesCmd.AddCommand(rolesAddCmd)
	rolesCmd.AddCommand(rolesUpdateCmd)
	rolesCmd.AddCommand(rolesRemoveCmd)
}
