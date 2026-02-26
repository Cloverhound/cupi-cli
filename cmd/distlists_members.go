package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var distlistsMembersCmd = &cobra.Command{
	Use:   "members",
	Short: "Manage distribution list members",
}

var distlistsMembersListCmd = &cobra.Command{
	Use:   "list <list-alias-or-objectId>",
	Short: "List members of a distribution list",
	Args:  cobra.ExactArgs(1),
	RunE:  runDistlistsMembersList,
}

var distlistsMembersAddCmd = &cobra.Command{
	Use:   "add <list-alias-or-objectId> <member-objectId>",
	Short: "Add a member to a distribution list",
	Args:  cobra.ExactArgs(2),
	RunE:  runDistlistsMembersAdd,
}

var distlistsMembersRemoveCmd = &cobra.Command{
	Use:   "remove <list-alias-or-objectId> <member-objectId>",
	Short: "Remove a member from a distribution list",
	Args:  cobra.ExactArgs(2),
	RunE:  runDistlistsMembersRemove,
}

func runDistlistsMembersList(cmd *cobra.Command, args []string) error {
	aliasOrID := args[0]

	serverName, err := resolveServer(cmd)
	if err != nil {
		return err
	}

	cfg, err := appconfig.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	serverCfg, err := appconfig.GetServer(cfg, serverName)
	if err != nil {
		return err
	}

	user, pass, err := auth.ResolveCreds(serverCfg, auth.CredTypeCUPI)
	if err != nil {
		return fmt.Errorf("failed to resolve credentials: %w", err)
	}

	members, err := client.ListDistListMembers(serverCfg.Host, serverCfg.Port, user, pass, aliasOrID)
	if err != nil {
		return fmt.Errorf("failed to list members: %w", err)
	}

	var rows []map[string]string
	for _, m := range members {
		rows = append(rows, map[string]string{
			"objectId":       m.ObjectId,
			"memberObjectId": m.MemberObjectId,
			"memberType":     m.MemberType,
		})
	}

	return output.Print(rows, outputFlag)
}

func runDistlistsMembersAdd(cmd *cobra.Command, args []string) error {
	listAliasOrID := args[0]
	memberObjectID := args[1]

	serverName, err := resolveServer(cmd)
	if err != nil {
		return err
	}

	cfg, err := appconfig.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	serverCfg, err := appconfig.GetServer(cfg, serverName)
	if err != nil {
		return err
	}

	user, pass, err := auth.ResolveCreds(serverCfg, auth.CredTypeCUPI)
	if err != nil {
		return fmt.Errorf("failed to resolve credentials: %w", err)
	}

	if err := client.AddDistListMember(serverCfg.Host, serverCfg.Port, user, pass, listAliasOrID, memberObjectID); err != nil {
		return fmt.Errorf("failed to add member: %w", err)
	}

	fmt.Printf("Added member %s to distribution list %s\n", memberObjectID, listAliasOrID)
	return nil
}

func runDistlistsMembersRemove(cmd *cobra.Command, args []string) error {
	listAliasOrID := args[0]
	memberObjectID := args[1]

	serverName, err := resolveServer(cmd)
	if err != nil {
		return err
	}

	cfg, err := appconfig.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	serverCfg, err := appconfig.GetServer(cfg, serverName)
	if err != nil {
		return err
	}

	user, pass, err := auth.ResolveCreds(serverCfg, auth.CredTypeCUPI)
	if err != nil {
		return fmt.Errorf("failed to resolve credentials: %w", err)
	}

	if err := client.RemoveDistListMember(serverCfg.Host, serverCfg.Port, user, pass, listAliasOrID, memberObjectID); err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}

	fmt.Printf("Removed member %s from distribution list %s\n", memberObjectID, listAliasOrID)
	return nil
}

func init() {
	distlistsMembersCmd.AddCommand(distlistsMembersListCmd)
	distlistsMembersCmd.AddCommand(distlistsMembersAddCmd)
	distlistsMembersCmd.AddCommand(distlistsMembersRemoveCmd)
}
