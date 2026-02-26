package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	psDisplayName string
	psQuery       string
	axlServerName string
	axlPort       string
)

var phoneSystemsCmd = &cobra.Command{
	Use:   "phonesystems",
	Short: "Manage phone systems",
}

var phoneSystemsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List phone systems",
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		items, err := client.ListPhoneSystems(srv.Host, srv.Port, cupiUser, cupiPass, psQuery, maxFlag)
		if err != nil {
			return err
		}
		var rows []map[string]string
		for _, v := range items {
			rows = append(rows, map[string]string{
				"objectId":        v.ObjectId,
				"displayName":     v.DisplayName,
				"callManagerType": v.CallManagerType,
			})
		}
		return output.Print(rows, outputFlag)
	},
}

var phoneSystemsGetCmd = &cobra.Command{
	Use:   "get <name-or-id>",
	Short: "Get a phone system",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		item, err := client.GetPhoneSystem(srv.Host, srv.Port, cupiUser, cupiPass, args[0])
		if err != nil {
			return err
		}
		return output.Print(item, outputFlag)
	},
}

var phoneSystemsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a phone system",
	RunE: func(cmd *cobra.Command, args []string) error {
		if psDisplayName == "" {
			return fmt.Errorf("--display-name is required")
		}
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{"DisplayName": psDisplayName}
		if err := client.CreatePhoneSystem(srv.Host, srv.Port, cupiUser, cupiPass, fields); err != nil {
			return err
		}
		fmt.Printf("Added phone system %s\n", psDisplayName)
		return nil
	},
}

var phoneSystemsUpdateCmd = &cobra.Command{
	Use:   "update <name-or-id>",
	Short: "Update a phone system",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{}
		if psDisplayName != "" {
			fields["DisplayName"] = psDisplayName
		}
		if len(fields) == 0 {
			return fmt.Errorf("no fields to update; use --help to see available flags")
		}
		if err := client.UpdatePhoneSystem(srv.Host, srv.Port, cupiUser, cupiPass, args[0], fields); err != nil {
			return err
		}
		fmt.Printf("Updated phone system %s\n", args[0])
		return nil
	},
}

var axlServersCmd = &cobra.Command{
	Use:   "axlservers",
	Short: "Manage AXL servers for a phone system",
}

var axlServersListCmd = &cobra.Command{
	Use:   "list <phone-system-id>",
	Short: "List AXL servers for a phone system",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		items, err := client.ListAXLServers(srv.Host, srv.Port, cupiUser, cupiPass, args[0])
		if err != nil {
			return err
		}
		var rows []map[string]string
		for _, v := range items {
			rows = append(rows, map[string]string{
				"objectId":   v.ObjectId,
				"serverName": v.ServerName,
				"port":       v.Port,
			})
		}
		return output.Print(rows, outputFlag)
	},
}

var axlServersAddCmd = &cobra.Command{
	Use:   "add <phone-system-id>",
	Short: "Add an AXL server to a phone system",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if axlServerName == "" {
			return fmt.Errorf("--server-name is required")
		}
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{"ServerName": axlServerName}
		if axlPort != "" {
			fields["Port"] = axlPort
		}
		if err := client.CreateAXLServer(srv.Host, srv.Port, cupiUser, cupiPass, args[0], fields); err != nil {
			return err
		}
		fmt.Printf("Added AXL server %s\n", axlServerName)
		return nil
	},
}

var axlServersRemoveCmd = &cobra.Command{
	Use:   "remove <phone-system-id> <axl-server-id>",
	Short: "Remove an AXL server from a phone system",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		if err := client.DeleteAXLServer(srv.Host, srv.Port, cupiUser, cupiPass, args[0], args[1]); err != nil {
			return err
		}
		fmt.Printf("Removed AXL server %s\n", args[1])
		return nil
	},
}

func init() {
	phoneSystemsListCmd.Flags().StringVar(&psQuery, "query", "", "Filter query")
	phoneSystemsAddCmd.Flags().StringVar(&psDisplayName, "display-name", "", "Display name (required)")
	phoneSystemsUpdateCmd.Flags().StringVar(&psDisplayName, "display-name", "", "Display name")

	axlServersAddCmd.Flags().StringVar(&axlServerName, "server-name", "", "AXL server hostname or IP (required)")
	axlServersAddCmd.Flags().StringVar(&axlPort, "port", "", "AXL port (default 8443)")

	axlServersCmd.AddCommand(axlServersListCmd)
	axlServersCmd.AddCommand(axlServersAddCmd)
	axlServersCmd.AddCommand(axlServersRemoveCmd)

	phoneSystemsCmd.AddCommand(phoneSystemsListCmd)
	phoneSystemsCmd.AddCommand(phoneSystemsGetCmd)
	phoneSystemsCmd.AddCommand(phoneSystemsAddCmd)
	phoneSystemsCmd.AddCommand(phoneSystemsUpdateCmd)
	phoneSystemsCmd.AddCommand(axlServersCmd)
}
