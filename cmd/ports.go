package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	portDisplayName string
	portEnabled     string
	portAnswerCalls string
	portQuery       string
)

var portsCmd = &cobra.Command{
	Use:   "ports",
	Short: "Manage system ports",
}

var portsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List ports",
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		items, err := client.ListPorts(srv.Host, srv.Port, cupiUser, cupiPass, portQuery, maxFlag)
		if err != nil {
			return err
		}
		var rows []map[string]string
		for _, v := range items {
			rows = append(rows, map[string]string{
				"objectId":    v.ObjectId,
				"displayName": v.DisplayName,
				"portNumber":  v.PortNumber,
				"enabled":     v.Enabled,
				"answerCalls": v.AnswerCalls,
			})
		}
		return output.Print(rows, outputFlag)
	},
}

var portsGetCmd = &cobra.Command{
	Use:   "get <name-or-id>",
	Short: "Get a port",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		item, err := client.GetPort(srv.Host, srv.Port, cupiUser, cupiPass, args[0])
		if err != nil {
			return err
		}
		return output.Print(item, outputFlag)
	},
}

var portsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a port",
	RunE: func(cmd *cobra.Command, args []string) error {
		if portDisplayName == "" {
			return fmt.Errorf("--display-name is required")
		}
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{"DisplayName": portDisplayName}
		if portEnabled != "" {
			fields["Enabled"] = portEnabled
		}
		if portAnswerCalls != "" {
			fields["AnswerCalls"] = portAnswerCalls
		}
		if err := client.CreatePort(srv.Host, srv.Port, cupiUser, cupiPass, fields); err != nil {
			return err
		}
		fmt.Printf("Added port %s\n", portDisplayName)
		return nil
	},
}

var portsUpdateCmd = &cobra.Command{
	Use:   "update <name-or-id>",
	Short: "Update a port",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{}
		if portDisplayName != "" {
			fields["DisplayName"] = portDisplayName
		}
		if portEnabled != "" {
			fields["Enabled"] = portEnabled
		}
		if portAnswerCalls != "" {
			fields["AnswerCalls"] = portAnswerCalls
		}
		if len(fields) == 0 {
			return fmt.Errorf("no fields to update; use --help to see available flags")
		}
		if err := client.UpdatePort(srv.Host, srv.Port, cupiUser, cupiPass, args[0], fields); err != nil {
			return err
		}
		fmt.Printf("Updated port %s\n", args[0])
		return nil
	},
}

var portsRemoveCmd = &cobra.Command{
	Use:   "remove <name-or-id>",
	Short: "Remove a port",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		if err := client.DeletePort(srv.Host, srv.Port, cupiUser, cupiPass, args[0]); err != nil {
			return err
		}
		fmt.Printf("Removed port %s\n", args[0])
		return nil
	},
}

func init() {
	portsListCmd.Flags().StringVar(&portQuery, "query", "", "Filter query")

	portsAddCmd.Flags().StringVar(&portDisplayName, "display-name", "", "Display name (required)")
	portsAddCmd.Flags().StringVar(&portEnabled, "enabled", "", "Enable port (true/false)")
	portsAddCmd.Flags().StringVar(&portAnswerCalls, "answer-calls", "", "Answer calls (true/false)")

	portsUpdateCmd.Flags().StringVar(&portDisplayName, "display-name", "", "Display name")
	portsUpdateCmd.Flags().StringVar(&portEnabled, "enabled", "", "Enable port (true/false)")
	portsUpdateCmd.Flags().StringVar(&portAnswerCalls, "answer-calls", "", "Answer calls (true/false)")

	portsCmd.AddCommand(portsListCmd)
	portsCmd.AddCommand(portsGetCmd)
	portsCmd.AddCommand(portsAddCmd)
	portsCmd.AddCommand(portsUpdateCmd)
	portsCmd.AddCommand(portsRemoveCmd)
}
