package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	chtDisplayName string
	chtQuery       string
)

var callHandlerTemplatesCmd = &cobra.Command{
	Use:   "chtemplates",
	Short: "Manage call handler templates",
}

var chtListCmd = &cobra.Command{
	Use:   "list",
	Short: "List call handler templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		items, err := client.ListCallHandlerTemplates(srv.Host, srv.Port, cupiUser, cupiPass, chtQuery, maxFlag)
		if err != nil {
			return err
		}
		var rows []map[string]string
		for _, v := range items {
			rows = append(rows, map[string]string{
				"objectId":    v.ObjectId,
				"displayName": v.DisplayName,
				"isDefault":   v.IsDefault,
			})
		}
		return output.Print(rows, outputFlag)
	},
}

var chtGetCmd = &cobra.Command{
	Use:   "get <name-or-id>",
	Short: "Get a call handler template",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		item, err := client.GetCallHandlerTemplate(srv.Host, srv.Port, cupiUser, cupiPass, args[0])
		if err != nil {
			return err
		}
		return output.Print(item, outputFlag)
	},
}

var chtAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a call handler template",
	RunE: func(cmd *cobra.Command, args []string) error {
		if chtDisplayName == "" {
			return fmt.Errorf("--display-name is required")
		}
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{"DisplayName": chtDisplayName}
		if err := client.CreateCallHandlerTemplate(srv.Host, srv.Port, cupiUser, cupiPass, fields); err != nil {
			return err
		}
		fmt.Printf("Added call handler template %s\n", chtDisplayName)
		return nil
	},
}

var chtUpdateCmd = &cobra.Command{
	Use:   "update <name-or-id>",
	Short: "Update a call handler template",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{}
		if chtDisplayName != "" {
			fields["DisplayName"] = chtDisplayName
		}
		if len(fields) == 0 {
			return fmt.Errorf("no fields to update; use --help to see available flags")
		}
		if err := client.UpdateCallHandlerTemplate(srv.Host, srv.Port, cupiUser, cupiPass, args[0], fields); err != nil {
			return err
		}
		fmt.Printf("Updated call handler template %s\n", args[0])
		return nil
	},
}

var chtRemoveCmd = &cobra.Command{
	Use:   "remove <name-or-id>",
	Short: "Remove a call handler template",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		if err := client.DeleteCallHandlerTemplate(srv.Host, srv.Port, cupiUser, cupiPass, args[0]); err != nil {
			return err
		}
		fmt.Printf("Removed call handler template %s\n", args[0])
		return nil
	},
}

func init() {
	chtListCmd.Flags().StringVar(&chtQuery, "query", "", "Filter query")
	chtAddCmd.Flags().StringVar(&chtDisplayName, "display-name", "", "Display name (required)")
	chtUpdateCmd.Flags().StringVar(&chtDisplayName, "display-name", "", "Display name")

	callHandlerTemplatesCmd.AddCommand(chtListCmd)
	callHandlerTemplatesCmd.AddCommand(chtGetCmd)
	callHandlerTemplatesCmd.AddCommand(chtAddCmd)
	callHandlerTemplatesCmd.AddCommand(chtUpdateCmd)
	callHandlerTemplatesCmd.AddCommand(chtRemoveCmd)
}
