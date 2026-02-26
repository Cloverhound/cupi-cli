package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var usersRolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "Manage user roles",
}

var usersRolesListCmd = &cobra.Command{
	Use:   "list <alias-or-objectId>",
	Short: "List roles assigned to a user",
	Args:  cobra.ExactArgs(1),
	RunE:  runUsersRolesList,
}

var usersRolesAddCmd = &cobra.Command{
	Use:   "add <alias-or-objectId> <roleObjectId>",
	Short: "Add a role to a user",
	Args:  cobra.ExactArgs(2),
	RunE:  runUsersRolesAdd,
}

var usersRolesRemoveCmd = &cobra.Command{
	Use:   "remove <alias-or-objectId> <roleObjectId>",
	Short: "Remove a role from a user",
	Args:  cobra.ExactArgs(2),
	RunE:  runUsersRolesRemove,
}

func runUsersRolesList(cmd *cobra.Command, args []string) error {
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

	roles, err := client.ListUserRoles(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId)
	if err != nil {
		return err
	}

	var rows []map[string]string
	for _, r := range roles {
		rows = append(rows, map[string]string{
			"objectId":     r.ObjectId,
			"roleObjectId": r.RoleObjectId,
			"userObjectId": r.UserObjectId,
		})
	}

	return output.Print(rows, outputFlag)
}

func runUsersRolesAdd(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]
	roleObjectId := args[1]

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

	if err := client.AddUserRole(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, roleObjectId); err != nil {
		return err
	}

	fmt.Printf("Added role %s to user %s\n", roleObjectId, u.Alias)
	return nil
}

func runUsersRolesRemove(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]
	roleObjectId := args[1]

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

	if err := client.RemoveUserRole(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, roleObjectId); err != nil {
		return err
	}

	fmt.Printf("Removed role %s from user %s\n", roleObjectId, u.Alias)
	return nil
}

func init() {
	usersRolesCmd.AddCommand(usersRolesListCmd)
	usersRolesCmd.AddCommand(usersRolesAddCmd)
	usersRolesCmd.AddCommand(usersRolesRemoveCmd)
}
