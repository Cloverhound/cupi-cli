package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	anFirstName   string
	anLastName    string
	anUserObjID   string
)

var alternateNamesCmd = &cobra.Command{
	Use:   "alternatenames",
	Short: "Manage alternate names for users",
}

var alternateNamesListCmd = &cobra.Command{
	Use:   "list <user-object-id>",
	Short: "List alternate names for a user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		items, err := client.ListAlternateNames(srv.Host, srv.Port, cupiUser, cupiPass, args[0])
		if err != nil {
			return err
		}
		var rows []map[string]string
		for _, v := range items {
			rows = append(rows, map[string]string{
				"objectId":  v.ObjectId,
				"firstName": v.FirstName,
				"lastName":  v.LastName,
			})
		}
		return output.Print(rows, outputFlag)
	},
}

var alternateNamesGetCmd = &cobra.Command{
	Use:   "get <object-id>",
	Short: "Get an alternate name by ObjectId",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		item, err := client.GetAlternateName(srv.Host, srv.Port, cupiUser, cupiPass, args[0])
		if err != nil {
			return err
		}
		return output.Print(item, outputFlag)
	},
}

var alternateNamesAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an alternate name",
	RunE: func(cmd *cobra.Command, args []string) error {
		if anUserObjID == "" || anFirstName == "" || anLastName == "" {
			return fmt.Errorf("--user-id, --first-name, and --last-name are required")
		}
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{
			"GlobalUserObjectId": anUserObjID,
			"FirstName":          anFirstName,
			"LastName":           anLastName,
		}
		if err := client.CreateAlternateName(srv.Host, srv.Port, cupiUser, cupiPass, fields); err != nil {
			return err
		}
		fmt.Printf("Added alternate name %s %s\n", anFirstName, anLastName)
		return nil
	},
}

var alternateNamesUpdateCmd = &cobra.Command{
	Use:   "update <object-id>",
	Short: "Update an alternate name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{}
		if anFirstName != "" {
			fields["FirstName"] = anFirstName
		}
		if anLastName != "" {
			fields["LastName"] = anLastName
		}
		if len(fields) == 0 {
			return fmt.Errorf("no fields to update")
		}
		if err := client.UpdateAlternateName(srv.Host, srv.Port, cupiUser, cupiPass, args[0], fields); err != nil {
			return err
		}
		fmt.Printf("Updated alternate name %s\n", args[0])
		return nil
	},
}

var alternateNamesRemoveCmd = &cobra.Command{
	Use:   "remove <object-id>",
	Short: "Remove an alternate name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		if err := client.DeleteAlternateName(srv.Host, srv.Port, cupiUser, cupiPass, args[0]); err != nil {
			return err
		}
		fmt.Printf("Removed alternate name %s\n", args[0])
		return nil
	},
}

func init() {
	alternateNamesAddCmd.Flags().StringVar(&anUserObjID, "user-id", "", "User ObjectId (required)")
	alternateNamesAddCmd.Flags().StringVar(&anFirstName, "first-name", "", "First name (required)")
	alternateNamesAddCmd.Flags().StringVar(&anLastName, "last-name", "", "Last name (required)")

	alternateNamesUpdateCmd.Flags().StringVar(&anFirstName, "first-name", "", "First name")
	alternateNamesUpdateCmd.Flags().StringVar(&anLastName, "last-name", "", "Last name")

	alternateNamesCmd.AddCommand(alternateNamesListCmd)
	alternateNamesCmd.AddCommand(alternateNamesGetCmd)
	alternateNamesCmd.AddCommand(alternateNamesAddCmd)
	alternateNamesCmd.AddCommand(alternateNamesUpdateCmd)
	alternateNamesCmd.AddCommand(alternateNamesRemoveCmd)
}
