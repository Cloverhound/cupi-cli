package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	plDisplayName string
	plNumericId   string
)

var usersPrivatelistsCmd = &cobra.Command{
	Use:   "privatelists",
	Short: "Manage user private distribution lists",
}

var usersPrivatelistsListCmd = &cobra.Command{
	Use:   "list <alias-or-objectId>",
	Short: "List private lists for a user",
	Args:  cobra.ExactArgs(1),
	RunE:  runUsersPrivatelistsList,
}

var usersPrivatelistsGetCmd = &cobra.Command{
	Use:   "get <alias-or-objectId> <listObjectId>",
	Short: "Get a specific private list",
	Args:  cobra.ExactArgs(2),
	RunE:  runUsersPrivatelistsGet,
}

var usersPrivatelistsAddCmd = &cobra.Command{
	Use:   "add <alias-or-objectId>",
	Short: "Add a private list to a user",
	Args:  cobra.ExactArgs(1),
	RunE:  runUsersPrivatelistsAdd,
}

var usersPrivatelistsUpdateCmd = &cobra.Command{
	Use:   "update <alias-or-objectId> <listObjectId>",
	Short: "Update a private list",
	Args:  cobra.ExactArgs(2),
	RunE:  runUsersPrivatelistsUpdate,
}

var usersPrivatelistsRemoveCmd = &cobra.Command{
	Use:   "remove <alias-or-objectId> <listObjectId>",
	Short: "Remove a private list from a user",
	Args:  cobra.ExactArgs(2),
	RunE:  runUsersPrivatelistsRemove,
}

var usersPrivatelistsMembersCmd = &cobra.Command{
	Use:   "members",
	Short: "Manage private list members",
}

var usersPrivatelistsMembersListCmd = &cobra.Command{
	Use:   "list <alias-or-objectId> <listObjectId>",
	Short: "List members of a private list",
	Args:  cobra.ExactArgs(2),
	RunE:  runUsersPrivatelistsMembersList,
}

var usersPrivatelistsMembersAddCmd = &cobra.Command{
	Use:   "add <alias-or-objectId> <listObjectId> <memberObjectId>",
	Short: "Add a member to a private list",
	Args:  cobra.ExactArgs(3),
	RunE:  runUsersPrivatelistsMembersAdd,
}

var usersPrivatelistsMembersRemoveCmd = &cobra.Command{
	Use:   "remove <alias-or-objectId> <listObjectId> <memberObjectId>",
	Short: "Remove a member from a private list",
	Args:  cobra.ExactArgs(3),
	RunE:  runUsersPrivatelistsMembersRemove,
}

func runUsersPrivatelistsList(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]

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

	u, err := client.GetUser(serverCfg.Host, serverCfg.Port, user, pass, userAliasOrID)
	if err != nil {
		return err
	}

	lists, err := client.ListPrivateLists(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId)
	if err != nil {
		return err
	}

	var rows []map[string]string
	for _, pl := range lists {
		rows = append(rows, map[string]string{
			"objectId":    pl.ObjectId,
			"displayName": pl.DisplayName,
			"numericId":   pl.NumericId,
		})
	}

	return output.Print(rows, outputFlag)
}

func runUsersPrivatelistsGet(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]
	listObjectId := args[1]

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

	u, err := client.GetUser(serverCfg.Host, serverCfg.Port, user, pass, userAliasOrID)
	if err != nil {
		return err
	}

	pl, err := client.GetPrivateList(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, listObjectId)
	if err != nil {
		return err
	}

	return output.Print(pl, outputFlag)
}

func runUsersPrivatelistsAdd(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]

	if plDisplayName == "" {
		return fmt.Errorf("--display-name is required")
	}

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

	u, err := client.GetUser(serverCfg.Host, serverCfg.Port, user, pass, userAliasOrID)
	if err != nil {
		return err
	}

	fields := map[string]interface{}{
		"DisplayName": plDisplayName,
	}
	if plNumericId != "" {
		fields["NumericId"] = plNumericId
	}

	if err := client.CreatePrivateList(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, fields); err != nil {
		return err
	}

	fmt.Printf("Added private list %s to user %s\n", plDisplayName, u.Alias)
	return nil
}

func runUsersPrivatelistsUpdate(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]
	listObjectId := args[1]

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

	u, err := client.GetUser(serverCfg.Host, serverCfg.Port, user, pass, userAliasOrID)
	if err != nil {
		return err
	}

	fields := map[string]interface{}{}
	if plDisplayName != "" {
		fields["DisplayName"] = plDisplayName
	}
	if plNumericId != "" {
		fields["NumericId"] = plNumericId
	}

	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	if err := client.UpdatePrivateList(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, listObjectId, fields); err != nil {
		return err
	}

	fmt.Printf("Updated private list %s for user %s\n", listObjectId, u.Alias)
	return nil
}

func runUsersPrivatelistsRemove(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]
	listObjectId := args[1]

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

	u, err := client.GetUser(serverCfg.Host, serverCfg.Port, user, pass, userAliasOrID)
	if err != nil {
		return err
	}

	if err := client.DeletePrivateList(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, listObjectId); err != nil {
		return err
	}

	fmt.Printf("Removed private list %s from user %s\n", listObjectId, u.Alias)
	return nil
}

func runUsersPrivatelistsMembersList(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]
	listObjectId := args[1]

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

	u, err := client.GetUser(serverCfg.Host, serverCfg.Port, user, pass, userAliasOrID)
	if err != nil {
		return err
	}

	members, err := client.ListPrivateListMembers(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, listObjectId)
	if err != nil {
		return err
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

func runUsersPrivatelistsMembersAdd(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]
	listObjectId := args[1]
	memberObjectId := args[2]

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

	u, err := client.GetUser(serverCfg.Host, serverCfg.Port, user, pass, userAliasOrID)
	if err != nil {
		return err
	}

	if err := client.AddPrivateListMember(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, listObjectId, memberObjectId); err != nil {
		return err
	}

	fmt.Printf("Added member %s to private list %s\n", memberObjectId, listObjectId)
	return nil
}

func runUsersPrivatelistsMembersRemove(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]
	listObjectId := args[1]
	memberObjectId := args[2]

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

	u, err := client.GetUser(serverCfg.Host, serverCfg.Port, user, pass, userAliasOrID)
	if err != nil {
		return err
	}

	if err := client.RemovePrivateListMember(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, listObjectId, memberObjectId); err != nil {
		return err
	}

	fmt.Printf("Removed member %s from private list %s\n", memberObjectId, listObjectId)
	return nil
}

func init() {
	usersPrivatelistsAddCmd.Flags().StringVar(&plDisplayName, "display-name", "", "Display name (required)")
	usersPrivatelistsAddCmd.Flags().StringVar(&plNumericId, "dtmf", "", "Numeric ID / DTMF code")

	usersPrivatelistsUpdateCmd.Flags().StringVar(&plDisplayName, "display-name", "", "Display name")
	usersPrivatelistsUpdateCmd.Flags().StringVar(&plNumericId, "dtmf", "", "Numeric ID / DTMF code")

	usersPrivatelistsMembersCmd.AddCommand(usersPrivatelistsMembersListCmd)
	usersPrivatelistsMembersCmd.AddCommand(usersPrivatelistsMembersAddCmd)
	usersPrivatelistsMembersCmd.AddCommand(usersPrivatelistsMembersRemoveCmd)

	usersPrivatelistsCmd.AddCommand(usersPrivatelistsListCmd)
	usersPrivatelistsCmd.AddCommand(usersPrivatelistsGetCmd)
	usersPrivatelistsCmd.AddCommand(usersPrivatelistsAddCmd)
	usersPrivatelistsCmd.AddCommand(usersPrivatelistsUpdateCmd)
	usersPrivatelistsCmd.AddCommand(usersPrivatelistsRemoveCmd)
	usersPrivatelistsCmd.AddCommand(usersPrivatelistsMembersCmd)
}
