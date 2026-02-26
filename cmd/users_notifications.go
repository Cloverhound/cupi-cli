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
	notifDeviceType string
	notifDisplayName string
	notifActive      string
)

var usersNotificationsCmd = &cobra.Command{
	Use:   "notifications",
	Short: "Manage user notification devices",
}

var usersNotificationsListCmd = &cobra.Command{
	Use:   "list <alias-or-objectId>",
	Short: "List notification devices for a user",
	Args:  cobra.ExactArgs(1),
	RunE:  runUsersNotificationsList,
}

var usersNotificationsGetCmd = &cobra.Command{
	Use:   "get <alias-or-objectId> <type> <deviceObjectId>",
	Short: "Get a specific notification device",
	Args:  cobra.ExactArgs(3),
	RunE:  runUsersNotificationsGet,
}

var usersNotificationsAddCmd = &cobra.Command{
	Use:   "add <alias-or-objectId> <type>",
	Short: "Add a notification device",
	Args:  cobra.ExactArgs(2),
	RunE:  runUsersNotificationsAdd,
}

var usersNotificationsUpdateCmd = &cobra.Command{
	Use:   "update <alias-or-objectId> <type> <deviceObjectId>",
	Short: "Update a notification device",
	Args:  cobra.ExactArgs(3),
	RunE:  runUsersNotificationsUpdate,
}

var usersNotificationsRemoveCmd = &cobra.Command{
	Use:   "remove <alias-or-objectId> <type> <deviceObjectId>",
	Short: "Remove a notification device",
	Args:  cobra.ExactArgs(3),
	RunE:  runUsersNotificationsRemove,
}

func runUsersNotificationsList(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]

	deviceType := notifDeviceType
	if deviceType == "" {
		deviceType = "phonedevices"
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

	devices, err := client.ListNotificationDevices(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, deviceType)
	if err != nil {
		return err
	}

	var rows []map[string]string
	for _, d := range devices {
		rows = append(rows, map[string]string{
			"objectId":    d.ObjectId,
			"displayName": d.DisplayName,
			"active":      d.Active,
		})
	}

	return output.Print(rows, outputFlag)
}

func runUsersNotificationsGet(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]
	deviceType := args[1]
	deviceObjectId := args[2]

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

	d, err := client.GetNotificationDevice(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, deviceType, deviceObjectId)
	if err != nil {
		return err
	}

	return output.Print(d, outputFlag)
}

func runUsersNotificationsAdd(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]
	deviceType := args[1]

	if notifDisplayName == "" {
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
		"DisplayName": notifDisplayName,
	}
	if notifActive != "" {
		fields["Active"] = notifActive
	}

	if err := client.CreateNotificationDevice(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, deviceType, fields); err != nil {
		return err
	}

	fmt.Printf("Added %s notification device %s to user %s\n", deviceType, notifDisplayName, u.Alias)
	return nil
}

func runUsersNotificationsUpdate(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]
	deviceType := args[1]
	deviceObjectId := args[2]

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
	if notifDisplayName != "" {
		fields["DisplayName"] = notifDisplayName
	}
	if notifActive != "" {
		fields["Active"] = notifActive
	}

	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	if err := client.UpdateNotificationDevice(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, deviceType, deviceObjectId, fields); err != nil {
		return err
	}

	fmt.Printf("Updated notification device %s for user %s\n", deviceObjectId, u.Alias)
	return nil
}

func runUsersNotificationsRemove(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]
	deviceType := args[1]
	deviceObjectId := args[2]

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

	if err := client.DeleteNotificationDevice(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, deviceType, deviceObjectId); err != nil {
		return err
	}

	fmt.Printf("Removed notification device %s from user %s\n", deviceObjectId, u.Alias)
	return nil
}

func init() {
	usersNotificationsListCmd.Flags().StringVar(&notifDeviceType, "type", "phonedevices", "Device type (phonedevices|pagerdevices|smtpdevices|htmldevices)")

	usersNotificationsAddCmd.Flags().StringVar(&notifDisplayName, "display-name", "", "Display name (required)")
	usersNotificationsAddCmd.Flags().StringVar(&notifActive, "active", "", "Active status (true|false)")

	usersNotificationsUpdateCmd.Flags().StringVar(&notifDisplayName, "display-name", "", "Display name")
	usersNotificationsUpdateCmd.Flags().StringVar(&notifActive, "active", "", "Active status (true|false)")

	usersNotificationsCmd.AddCommand(usersNotificationsListCmd)
	usersNotificationsCmd.AddCommand(usersNotificationsGetCmd)
	usersNotificationsCmd.AddCommand(usersNotificationsAddCmd)
	usersNotificationsCmd.AddCommand(usersNotificationsUpdateCmd)
	usersNotificationsCmd.AddCommand(usersNotificationsRemoveCmd)
}
