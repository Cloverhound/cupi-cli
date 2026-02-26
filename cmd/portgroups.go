package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	pgDisplayName string
	pgMWIOnCode   string
	pgMWIOffCode  string
	pgQuery       string
)

var portGroupsCmd = &cobra.Command{
	Use:   "portgroups",
	Short: "Manage port groups",
}

var portGroupsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List port groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		items, err := client.ListPortGroups(srv.Host, srv.Port, cupiUser, cupiPass, pgQuery, maxFlag)
		if err != nil {
			return err
		}
		var rows []map[string]string
		for _, v := range items {
			rows = append(rows, map[string]string{
				"objectId":    v.ObjectId,
				"displayName": v.DisplayName,
				"mwiOnCode":   v.MWIOnCode,
				"mwiOffCode":  v.MWIOffCode,
			})
		}
		return output.Print(rows, outputFlag)
	},
}

var portGroupsGetCmd = &cobra.Command{
	Use:   "get <name-or-id>",
	Short: "Get a port group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		item, err := client.GetPortGroup(srv.Host, srv.Port, cupiUser, cupiPass, args[0])
		if err != nil {
			return err
		}
		return output.Print(item, outputFlag)
	},
}

var portGroupsUpdateCmd = &cobra.Command{
	Use:   "update <name-or-id>",
	Short: "Update a port group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{}
		if pgDisplayName != "" {
			fields["DisplayName"] = pgDisplayName
		}
		if pgMWIOnCode != "" {
			fields["MWIOnCode"] = pgMWIOnCode
		}
		if pgMWIOffCode != "" {
			fields["MWIOffCode"] = pgMWIOffCode
		}
		if len(fields) == 0 {
			return fmt.Errorf("no fields to update; use --help to see available flags")
		}
		if err := client.UpdatePortGroup(srv.Host, srv.Port, cupiUser, cupiPass, args[0], fields); err != nil {
			return err
		}
		fmt.Printf("Updated port group %s\n", args[0])
		return nil
	},
}

func init() {
	portGroupsListCmd.Flags().StringVar(&pgQuery, "query", "", "Filter query")
	portGroupsUpdateCmd.Flags().StringVar(&pgDisplayName, "display-name", "", "Display name")
	portGroupsUpdateCmd.Flags().StringVar(&pgMWIOnCode, "mwi-on-code", "", "MWI on code")
	portGroupsUpdateCmd.Flags().StringVar(&pgMWIOffCode, "mwi-off-code", "", "MWI off code")

	portGroupsCmd.AddCommand(portGroupsListCmd)
	portGroupsCmd.AddCommand(portGroupsGetCmd)
	portGroupsCmd.AddCommand(portGroupsUpdateCmd)
}
