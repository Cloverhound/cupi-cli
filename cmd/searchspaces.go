package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	ssName        string
	ssDescription string
	ssQuery       string
	ssPartitionID string
)

var searchSpacesCmd = &cobra.Command{
	Use:   "searchspaces",
	Short: "Manage search spaces",
}

var searchSpacesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List search spaces",
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		items, err := client.ListSearchSpaces(srv.Host, srv.Port, cupiUser, cupiPass, ssQuery, maxFlag)
		if err != nil {
			return err
		}
		var rows []map[string]string
		for _, v := range items {
			rows = append(rows, map[string]string{
				"objectId":    v.ObjectId,
				"name":        v.Name,
				"description": v.Description,
			})
		}
		return output.Print(rows, outputFlag)
	},
}

var searchSpacesGetCmd = &cobra.Command{
	Use:   "get <name-or-id>",
	Short: "Get a search space",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		item, err := client.GetSearchSpace(srv.Host, srv.Port, cupiUser, cupiPass, args[0])
		if err != nil {
			return err
		}
		return output.Print(item, outputFlag)
	},
}

var searchSpacesAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a search space",
	RunE: func(cmd *cobra.Command, args []string) error {
		if ssName == "" {
			return fmt.Errorf("--name is required")
		}
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{"Name": ssName}
		if ssDescription != "" {
			fields["Description"] = ssDescription
		}
		if err := client.CreateSearchSpace(srv.Host, srv.Port, cupiUser, cupiPass, fields); err != nil {
			return err
		}
		fmt.Printf("Added search space %s\n", ssName)
		return nil
	},
}

var searchSpacesUpdateCmd = &cobra.Command{
	Use:   "update <name-or-id>",
	Short: "Update a search space",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{}
		if ssName != "" {
			fields["Name"] = ssName
		}
		if ssDescription != "" {
			fields["Description"] = ssDescription
		}
		if len(fields) == 0 {
			return fmt.Errorf("no fields to update; use --help to see available flags")
		}
		if err := client.UpdateSearchSpace(srv.Host, srv.Port, cupiUser, cupiPass, args[0], fields); err != nil {
			return err
		}
		fmt.Printf("Updated search space %s\n", args[0])
		return nil
	},
}

var searchSpacesRemoveCmd = &cobra.Command{
	Use:   "remove <name-or-id>",
	Short: "Remove a search space",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		if err := client.DeleteSearchSpace(srv.Host, srv.Port, cupiUser, cupiPass, args[0]); err != nil {
			return err
		}
		fmt.Printf("Removed search space %s\n", args[0])
		return nil
	},
}

var ssMembersCmd = &cobra.Command{
	Use:   "members",
	Short: "Manage search space partition members",
}

var ssMembersListCmd = &cobra.Command{
	Use:   "list <search-space-id>",
	Short: "List partition members of a search space",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		items, err := client.ListSearchSpaceMembers(srv.Host, srv.Port, cupiUser, cupiPass, args[0])
		if err != nil {
			return err
		}
		var rows []map[string]string
		for _, v := range items {
			rows = append(rows, map[string]string{
				"objectId":          v.ObjectId,
				"partitionObjectId": v.PartitionObjectId,
			})
		}
		return output.Print(rows, outputFlag)
	},
}

var ssMembersAddCmd = &cobra.Command{
	Use:   "add <search-space-id>",
	Short: "Add a partition to a search space",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if ssPartitionID == "" {
			return fmt.Errorf("--partition-id is required")
		}
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		if err := client.AddSearchSpaceMember(srv.Host, srv.Port, cupiUser, cupiPass, args[0], ssPartitionID); err != nil {
			return err
		}
		fmt.Printf("Added partition %s to search space %s\n", ssPartitionID, args[0])
		return nil
	},
}

var ssMembersRemoveCmd = &cobra.Command{
	Use:   "remove <search-space-id> <member-id>",
	Short: "Remove a partition member from a search space",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		if err := client.RemoveSearchSpaceMember(srv.Host, srv.Port, cupiUser, cupiPass, args[0], args[1]); err != nil {
			return err
		}
		fmt.Printf("Removed member %s from search space %s\n", args[1], args[0])
		return nil
	},
}

func init() {
	searchSpacesListCmd.Flags().StringVar(&ssQuery, "query", "", "Filter query")

	searchSpacesAddCmd.Flags().StringVar(&ssName, "name", "", "Search space name (required)")
	searchSpacesAddCmd.Flags().StringVar(&ssDescription, "description", "", "Description")

	searchSpacesUpdateCmd.Flags().StringVar(&ssName, "name", "", "Search space name")
	searchSpacesUpdateCmd.Flags().StringVar(&ssDescription, "description", "", "Description")

	ssMembersAddCmd.Flags().StringVar(&ssPartitionID, "partition-id", "", "Partition ObjectId (required)")

	ssMembersCmd.AddCommand(ssMembersListCmd)
	ssMembersCmd.AddCommand(ssMembersAddCmd)
	ssMembersCmd.AddCommand(ssMembersRemoveCmd)

	searchSpacesCmd.AddCommand(searchSpacesListCmd)
	searchSpacesCmd.AddCommand(searchSpacesGetCmd)
	searchSpacesCmd.AddCommand(searchSpacesAddCmd)
	searchSpacesCmd.AddCommand(searchSpacesUpdateCmd)
	searchSpacesCmd.AddCommand(searchSpacesRemoveCmd)
	searchSpacesCmd.AddCommand(ssMembersCmd)
}
