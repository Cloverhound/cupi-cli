package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	cvQuery string
	cvValue string
)

var configValuesCmd = &cobra.Command{
	Use:   "configvalues",
	Short: "Manage configuration values",
}

var configValuesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configuration values",
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		items, err := client.ListConfigValues(srv.Host, srv.Port, cupiUser, cupiPass, cvQuery, maxFlag)
		if err != nil {
			return err
		}
		var rows []map[string]string
		for _, v := range items {
			rows = append(rows, map[string]string{
				"fullName":    v.FullName,
				"value":       v.Value,
				"type":        v.Type,
				"description": v.Description,
			})
		}
		return output.Print(rows, outputFlag)
	},
}

var configValuesGetCmd = &cobra.Command{
	Use:   "get <full-name>",
	Short: "Get a configuration value by full name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		item, err := client.GetConfigValue(srv.Host, srv.Port, cupiUser, cupiPass, args[0])
		if err != nil {
			return err
		}
		return output.Print(item, outputFlag)
	},
}

var configValuesUpdateCmd = &cobra.Command{
	Use:   "update <full-name>",
	Short: "Update a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if cvValue == "" {
			return fmt.Errorf("--value is required")
		}
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		if err := client.UpdateConfigValue(srv.Host, srv.Port, cupiUser, cupiPass, args[0], cvValue); err != nil {
			return err
		}
		fmt.Printf("Updated configuration value %s\n", args[0])
		return nil
	},
}

func init() {
	configValuesListCmd.Flags().StringVar(&cvQuery, "query", "", "Filter query")
	configValuesUpdateCmd.Flags().StringVar(&cvValue, "value", "", "New value (required)")

	configValuesCmd.AddCommand(configValuesListCmd)
	configValuesCmd.AddCommand(configValuesGetCmd)
	configValuesCmd.AddCommand(configValuesUpdateCmd)
}
