package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var dhDisplayName string

var dirhandlersCmd = &cobra.Command{
	Use:   "dirhandlers",
	Short: "Manage directory handlers",
}

var dirhandlersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List directory handlers",
	RunE:  runDirhandlersList,
}

var dirhandlersGetCmd = &cobra.Command{
	Use:   "get <name-or-id>",
	Short: "Get a directory handler",
	Args:  cobra.ExactArgs(1),
	RunE:  runDirhandlersGet,
}

var dirhandlersAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a directory handler",
	RunE:  runDirhandlersAdd,
}

var dirhandlersUpdateCmd = &cobra.Command{
	Use:   "update <name-or-id>",
	Short: "Update a directory handler",
	Args:  cobra.ExactArgs(1),
	RunE:  runDirhandlersUpdate,
}

var dirhandlersRemoveCmd = &cobra.Command{
	Use:   "remove <name-or-id>",
	Short: "Remove a directory handler",
	Args:  cobra.ExactArgs(1),
	RunE:  runDirhandlersRemove,
}

func runDirhandlersList(cmd *cobra.Command, args []string) error {
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

	items, err := client.ListDirectoryHandlers(serverCfg.Host, serverCfg.Port, user, pass, "", maxFlag)
	if err != nil {
		return err
	}

	var rows []map[string]string
	for _, item := range items {
		rows = append(rows, map[string]string{
			"objectId":     item.ObjectId,
			"displayName":  item.DisplayName,
			"dtmfAccessId": item.DtmfAccessId,
		})
	}

	return output.Print(rows, outputFlag)
}

func runDirhandlersGet(cmd *cobra.Command, args []string) error {
	nameOrID := args[0]

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

	item, err := client.GetDirectoryHandler(serverCfg.Host, serverCfg.Port, user, pass, nameOrID)
	if err != nil {
		return err
	}

	return output.Print(item, outputFlag)
}

func runDirhandlersAdd(cmd *cobra.Command, args []string) error {
	if dhDisplayName == "" {
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

	fields := map[string]interface{}{
		"DisplayName": dhDisplayName,
	}

	if err := client.CreateDirectoryHandler(serverCfg.Host, serverCfg.Port, user, pass, fields); err != nil {
		return err
	}

	fmt.Printf("Added directory handler %s\n", dhDisplayName)
	return nil
}

func runDirhandlersUpdate(cmd *cobra.Command, args []string) error {
	nameOrID := args[0]

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

	fields := map[string]interface{}{}
	if dhDisplayName != "" {
		fields["DisplayName"] = dhDisplayName
	}

	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	if err := client.UpdateDirectoryHandler(serverCfg.Host, serverCfg.Port, user, pass, nameOrID, fields); err != nil {
		return err
	}

	fmt.Printf("Updated directory handler %s\n", nameOrID)
	return nil
}

func runDirhandlersRemove(cmd *cobra.Command, args []string) error {
	nameOrID := args[0]

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

	if err := client.DeleteDirectoryHandler(serverCfg.Host, serverCfg.Port, user, pass, nameOrID); err != nil {
		return err
	}

	fmt.Printf("Removed directory handler %s\n", nameOrID)
	return nil
}

func init() {
	dirhandlersAddCmd.Flags().StringVar(&dhDisplayName, "display-name", "", "Display name (required)")
	dirhandlersUpdateCmd.Flags().StringVar(&dhDisplayName, "display-name", "", "Display name")

	dirhandlersCmd.AddCommand(dirhandlersListCmd)
	dirhandlersCmd.AddCommand(dirhandlersGetCmd)
	dirhandlersCmd.AddCommand(dirhandlersAddCmd)
	dirhandlersCmd.AddCommand(dirhandlersUpdateCmd)
	dirhandlersCmd.AddCommand(dirhandlersRemoveCmd)
}
