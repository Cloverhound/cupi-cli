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
	altextsIndex     int
	altextsPartition string
)

var usersAltextsCmd = &cobra.Command{
	Use:   "altexts",
	Short: "Manage user alternate extensions",
}

var usersAltextsListCmd = &cobra.Command{
	Use:   "list <alias-or-objectId>",
	Short: "List alternate extensions for a user",
	Args:  cobra.ExactArgs(1),
	RunE:  runUsersAltextsList,
}

var usersAltextsGetCmd = &cobra.Command{
	Use:   "get <alias-or-objectId> <objectId>",
	Short: "Get a specific alternate extension",
	Args:  cobra.ExactArgs(2),
	RunE:  runUsersAltextsGet,
}

var usersAltextsAddCmd = &cobra.Command{
	Use:   "add <alias-or-objectId> --dtmf <ext>",
	Short: "Add an alternate extension to a user",
	Args:  cobra.ExactArgs(1),
	RunE:  runUsersAltextsAdd,
}

var usersAltextsUpdateCmd = &cobra.Command{
	Use:   "update <alias-or-objectId> <objectId>",
	Short: "Update an alternate extension",
	Args:  cobra.ExactArgs(2),
	RunE:  runUsersAltextsUpdate,
}

var usersAltextsRemoveCmd = &cobra.Command{
	Use:   "remove <alias-or-objectId> <objectId>",
	Short: "Remove an alternate extension from a user",
	Args:  cobra.ExactArgs(2),
	RunE:  runUsersAltextsRemove,
}

func runUsersAltextsList(cmd *cobra.Command, args []string) error {
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

	altexts, err := client.ListAlternateExtensions(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId)
	if err != nil {
		return err
	}

	var rows []map[string]string
	for _, ae := range altexts {
		rows = append(rows, map[string]string{
			"objectId":          ae.ObjectId,
			"dtmfAccessId":      ae.DtmfAccessId,
			"idIndex":           fmt.Sprintf("%d", ae.IdIndex),
			"partitionObjectId": ae.PartitionObjectId,
		})
	}

	return output.Print(rows, outputFlag)
}

func runUsersAltextsGet(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]
	objectId := args[1]

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

	ae, err := client.GetAlternateExtension(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, objectId)
	if err != nil {
		return err
	}

	return output.Print(ae, outputFlag)
}

func runUsersAltextsAdd(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]

	if addUserDtmf == "" {
		return fmt.Errorf("--dtmf is required")
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
		"DtmfAccessId": addUserDtmf,
	}
	if altextsIndex > 0 {
		fields["IdIndex"] = altextsIndex
	}
	if altextsPartition != "" {
		fields["PartitionObjectId"] = altextsPartition
	}

	if err := client.CreateAlternateExtension(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, fields); err != nil {
		return err
	}

	fmt.Printf("Added alternate extension %s to user %s\n", addUserDtmf, u.Alias)
	return nil
}

func runUsersAltextsUpdate(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]
	objectId := args[1]

	if addUserDtmf == "" {
		return fmt.Errorf("--dtmf is required")
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
		"DtmfAccessId": addUserDtmf,
	}

	if err := client.UpdateAlternateExtension(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, objectId, fields); err != nil {
		return err
	}

	fmt.Printf("Updated alternate extension %s for user %s\n", objectId, u.Alias)
	return nil
}

func runUsersAltextsRemove(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]
	objectId := args[1]

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

	if err := client.DeleteAlternateExtension(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, objectId); err != nil {
		return err
	}

	fmt.Printf("Removed alternate extension %s from user %s\n", objectId, u.Alias)
	return nil
}

func init() {
	usersAltextsAddCmd.Flags().StringVar(&addUserDtmf, "dtmf", "", "DTMF access ID (required)")
	usersAltextsAddCmd.Flags().IntVar(&altextsIndex, "index", 0, "Index for the extension")
	usersAltextsAddCmd.Flags().StringVar(&altextsPartition, "partition-id", "", "Partition ObjectId")

	usersAltextsUpdateCmd.Flags().StringVar(&addUserDtmf, "dtmf", "", "DTMF access ID (required)")

	usersAltextsCmd.AddCommand(usersAltextsListCmd)
	usersAltextsCmd.AddCommand(usersAltextsGetCmd)
	usersAltextsCmd.AddCommand(usersAltextsAddCmd)
	usersAltextsCmd.AddCommand(usersAltextsUpdateCmd)
	usersAltextsCmd.AddCommand(usersAltextsRemoveCmd)
}
