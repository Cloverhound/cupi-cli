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
	mwiDisplayName   string
	mwiPhoneSystemId string
	mwiExtension     string
)

var usersMwiCmd = &cobra.Command{
	Use:   "mwi",
	Short: "Manage user message waiting indicators",
}

var usersMwiListCmd = &cobra.Command{
	Use:   "list <alias-or-objectId>",
	Short: "List MWIs for a user",
	Args:  cobra.ExactArgs(1),
	RunE:  runUsersMwiList,
}

var usersMwiGetCmd = &cobra.Command{
	Use:   "get <alias-or-objectId> <mwiObjectId>",
	Short: "Get a specific MWI",
	Args:  cobra.ExactArgs(2),
	RunE:  runUsersMwiGet,
}

var usersMwiAddCmd = &cobra.Command{
	Use:   "add <alias-or-objectId>",
	Short: "Add an MWI to a user",
	Args:  cobra.ExactArgs(1),
	RunE:  runUsersMwiAdd,
}

var usersMwiUpdateCmd = &cobra.Command{
	Use:   "update <alias-or-objectId> <mwiObjectId>",
	Short: "Update an MWI",
	Args:  cobra.ExactArgs(2),
	RunE:  runUsersMwiUpdate,
}

var usersMwiRemoveCmd = &cobra.Command{
	Use:   "remove <alias-or-objectId> <mwiObjectId>",
	Short: "Remove an MWI from a user",
	Args:  cobra.ExactArgs(2),
	RunE:  runUsersMwiRemove,
}

func runUsersMwiList(cmd *cobra.Command, args []string) error {
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

	mwis, err := client.ListMWIs(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId)
	if err != nil {
		return err
	}

	var rows []map[string]string
	for _, m := range mwis {
		rows = append(rows, map[string]string{
			"objectId":            m.ObjectId,
			"displayName":         m.DisplayName,
			"active":              m.Active,
			"mwiExtension":        m.MWIExtension,
			"mediaSwitchObjectId": m.MediaSwitchObjectId,
		})
	}

	return output.Print(rows, outputFlag)
}

func runUsersMwiGet(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]
	mwiObjectId := args[1]

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

	m, err := client.GetMWI(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, mwiObjectId)
	if err != nil {
		return err
	}

	return output.Print(m, outputFlag)
}

func runUsersMwiAdd(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]

	if mwiDisplayName == "" {
		return fmt.Errorf("--display-name is required")
	}
	if mwiPhoneSystemId == "" {
		return fmt.Errorf("--phone-system-id is required")
	}
	if mwiExtension == "" {
		return fmt.Errorf("--extension is required")
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
		"DisplayName":         mwiDisplayName,
		"MediaSwitchObjectId": mwiPhoneSystemId,
		"MWIExtension":        mwiExtension,
	}

	if err := client.CreateMWI(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, fields); err != nil {
		return err
	}

	fmt.Printf("Added MWI %s to user %s\n", mwiDisplayName, u.Alias)
	return nil
}

func runUsersMwiUpdate(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]
	mwiObjectId := args[1]

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
	if mwiDisplayName != "" {
		fields["DisplayName"] = mwiDisplayName
	}
	if mwiPhoneSystemId != "" {
		fields["MediaSwitchObjectId"] = mwiPhoneSystemId
	}
	if mwiExtension != "" {
		fields["MWIExtension"] = mwiExtension
	}

	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	if err := client.UpdateMWI(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, mwiObjectId, fields); err != nil {
		return err
	}

	fmt.Printf("Updated MWI %s for user %s\n", mwiObjectId, u.Alias)
	return nil
}

func runUsersMwiRemove(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]
	mwiObjectId := args[1]

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

	if err := client.DeleteMWI(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, mwiObjectId); err != nil {
		return err
	}

	fmt.Printf("Removed MWI %s from user %s\n", mwiObjectId, u.Alias)
	return nil
}

func init() {
	usersMwiAddCmd.Flags().StringVar(&mwiDisplayName, "display-name", "", "Display name (required)")
	usersMwiAddCmd.Flags().StringVar(&mwiPhoneSystemId, "phone-system-id", "", "Phone system ObjectId (required)")
	usersMwiAddCmd.Flags().StringVar(&mwiExtension, "extension", "", "MWI extension (required)")

	usersMwiUpdateCmd.Flags().StringVar(&mwiDisplayName, "display-name", "", "Display name")
	usersMwiUpdateCmd.Flags().StringVar(&mwiPhoneSystemId, "phone-system-id", "", "Phone system ObjectId")
	usersMwiUpdateCmd.Flags().StringVar(&mwiExtension, "extension", "", "MWI extension")

	usersMwiCmd.AddCommand(usersMwiListCmd)
	usersMwiCmd.AddCommand(usersMwiGetCmd)
	usersMwiCmd.AddCommand(usersMwiAddCmd)
	usersMwiCmd.AddCommand(usersMwiUpdateCmd)
	usersMwiCmd.AddCommand(usersMwiRemoveCmd)
}
