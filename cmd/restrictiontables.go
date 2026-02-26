package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	rtDisplayName  string
	rtQuery        string
	rpNumberPattern string
	rpBlocked      string
)

var restrictionTablesCmd = &cobra.Command{
	Use:   "restrictiontables",
	Short: "Manage restriction tables",
}

var restrictionTablesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List restriction tables",
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		items, err := client.ListRestrictionTables(srv.Host, srv.Port, cupiUser, cupiPass, rtQuery, maxFlag)
		if err != nil {
			return err
		}
		var rows []map[string]string
		for _, v := range items {
			rows = append(rows, map[string]string{
				"objectId":    v.ObjectId,
				"displayName": v.DisplayName,
			})
		}
		return output.Print(rows, outputFlag)
	},
}

var restrictionTablesGetCmd = &cobra.Command{
	Use:   "get <name-or-id>",
	Short: "Get a restriction table",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		item, err := client.GetRestrictionTable(srv.Host, srv.Port, cupiUser, cupiPass, args[0])
		if err != nil {
			return err
		}
		return output.Print(item, outputFlag)
	},
}

var restrictionTablesAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a restriction table",
	RunE: func(cmd *cobra.Command, args []string) error {
		if rtDisplayName == "" {
			return fmt.Errorf("--display-name is required")
		}
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{"DisplayName": rtDisplayName}
		if err := client.CreateRestrictionTable(srv.Host, srv.Port, cupiUser, cupiPass, fields); err != nil {
			return err
		}
		fmt.Printf("Added restriction table %s\n", rtDisplayName)
		return nil
	},
}

var restrictionTablesUpdateCmd = &cobra.Command{
	Use:   "update <name-or-id>",
	Short: "Update a restriction table",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{}
		if rtDisplayName != "" {
			fields["DisplayName"] = rtDisplayName
		}
		if len(fields) == 0 {
			return fmt.Errorf("no fields to update; use --help to see available flags")
		}
		if err := client.UpdateRestrictionTable(srv.Host, srv.Port, cupiUser, cupiPass, args[0], fields); err != nil {
			return err
		}
		fmt.Printf("Updated restriction table %s\n", args[0])
		return nil
	},
}

var restrictionTablesRemoveCmd = &cobra.Command{
	Use:   "remove <name-or-id>",
	Short: "Remove a restriction table",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		if err := client.DeleteRestrictionTable(srv.Host, srv.Port, cupiUser, cupiPass, args[0]); err != nil {
			return err
		}
		fmt.Printf("Removed restriction table %s\n", args[0])
		return nil
	},
}

var rtPatternsCmd = &cobra.Command{
	Use:   "patterns",
	Short: "Manage restriction patterns",
}

var rtPatternsListCmd = &cobra.Command{
	Use:   "list <table-id>",
	Short: "List patterns in a restriction table",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		items, err := client.ListRestrictionPatterns(srv.Host, srv.Port, cupiUser, cupiPass, args[0])
		if err != nil {
			return err
		}
		var rows []map[string]string
		for _, v := range items {
			rows = append(rows, map[string]string{
				"objectId":      v.ObjectId,
				"numberPattern": v.NumberPattern,
				"blocked":       v.Blocked,
			})
		}
		return output.Print(rows, outputFlag)
	},
}

var rtPatternsAddCmd = &cobra.Command{
	Use:   "add <table-id>",
	Short: "Add a pattern to a restriction table",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if rpNumberPattern == "" {
			return fmt.Errorf("--pattern is required")
		}
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{"NumberPattern": rpNumberPattern}
		if rpBlocked != "" {
			fields["Blocked"] = rpBlocked
		}
		if err := client.CreateRestrictionPattern(srv.Host, srv.Port, cupiUser, cupiPass, args[0], fields); err != nil {
			return err
		}
		fmt.Printf("Added pattern %s\n", rpNumberPattern)
		return nil
	},
}

var rtPatternsRemoveCmd = &cobra.Command{
	Use:   "remove <table-id> <pattern-id>",
	Short: "Remove a pattern from a restriction table",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		if err := client.DeleteRestrictionPattern(srv.Host, srv.Port, cupiUser, cupiPass, args[0], args[1]); err != nil {
			return err
		}
		fmt.Printf("Removed pattern %s\n", args[1])
		return nil
	},
}

func init() {
	restrictionTablesListCmd.Flags().StringVar(&rtQuery, "query", "", "Filter query")
	restrictionTablesAddCmd.Flags().StringVar(&rtDisplayName, "display-name", "", "Display name (required)")
	restrictionTablesUpdateCmd.Flags().StringVar(&rtDisplayName, "display-name", "", "Display name")

	rtPatternsAddCmd.Flags().StringVar(&rpNumberPattern, "pattern", "", "Number pattern (required)")
	rtPatternsAddCmd.Flags().StringVar(&rpBlocked, "blocked", "", "Blocked (true/false)")

	rtPatternsCmd.AddCommand(rtPatternsListCmd)
	rtPatternsCmd.AddCommand(rtPatternsAddCmd)
	rtPatternsCmd.AddCommand(rtPatternsRemoveCmd)

	restrictionTablesCmd.AddCommand(restrictionTablesListCmd)
	restrictionTablesCmd.AddCommand(restrictionTablesGetCmd)
	restrictionTablesCmd.AddCommand(restrictionTablesAddCmd)
	restrictionTablesCmd.AddCommand(restrictionTablesUpdateCmd)
	restrictionTablesCmd.AddCommand(restrictionTablesRemoveCmd)
	restrictionTablesCmd.AddCommand(rtPatternsCmd)
}
