package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var ihDisplayName string

var inthandlersCmd = &cobra.Command{
	Use:   "inthandlers",
	Short: "Manage interview handlers",
}

var inthandlersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List interview handlers",
	RunE:  runInthandlersList,
}

var inthandlersGetCmd = &cobra.Command{
	Use:   "get <name-or-id>",
	Short: "Get an interview handler",
	Args:  cobra.ExactArgs(1),
	RunE:  runInthandlersGet,
}

var inthandlersAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an interview handler",
	RunE:  runInthandlersAdd,
}

var inthandlersUpdateCmd = &cobra.Command{
	Use:   "update <name-or-id>",
	Short: "Update an interview handler",
	Args:  cobra.ExactArgs(1),
	RunE:  runInthandlersUpdate,
}

var inthandlersRemoveCmd = &cobra.Command{
	Use:   "remove <name-or-id>",
	Short: "Remove an interview handler",
	Args:  cobra.ExactArgs(1),
	RunE:  runInthandlersRemove,
}

func runInthandlersList(cmd *cobra.Command, args []string) error {
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

	items, err := client.ListInterviewHandlers(serverCfg.Host, serverCfg.Port, user, pass, "", maxFlag)
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

func runInthandlersGet(cmd *cobra.Command, args []string) error {
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

	item, err := client.GetInterviewHandler(serverCfg.Host, serverCfg.Port, user, pass, nameOrID)
	if err != nil {
		return err
	}

	return output.Print(item, outputFlag)
}

func runInthandlersAdd(cmd *cobra.Command, args []string) error {
	if ihDisplayName == "" {
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
		"DisplayName": ihDisplayName,
	}

	if err := client.CreateInterviewHandler(serverCfg.Host, serverCfg.Port, user, pass, fields); err != nil {
		return err
	}

	fmt.Printf("Added interview handler %s\n", ihDisplayName)
	return nil
}

func runInthandlersUpdate(cmd *cobra.Command, args []string) error {
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
	if ihDisplayName != "" {
		fields["DisplayName"] = ihDisplayName
	}

	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	if err := client.UpdateInterviewHandler(serverCfg.Host, serverCfg.Port, user, pass, nameOrID, fields); err != nil {
		return err
	}

	fmt.Printf("Updated interview handler %s\n", nameOrID)
	return nil
}

func runInthandlersRemove(cmd *cobra.Command, args []string) error {
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

	if err := client.DeleteInterviewHandler(serverCfg.Host, serverCfg.Port, user, pass, nameOrID); err != nil {
		return err
	}

	fmt.Printf("Removed interview handler %s\n", nameOrID)
	return nil
}

func init() {
	inthandlersAddCmd.Flags().StringVar(&ihDisplayName, "display-name", "", "Display name (required)")
	inthandlersUpdateCmd.Flags().StringVar(&ihDisplayName, "display-name", "", "Display name")

	inthandlersCmd.AddCommand(inthandlersListCmd)
	inthandlersCmd.AddCommand(inthandlersGetCmd)
	inthandlersCmd.AddCommand(inthandlersAddCmd)
	inthandlersCmd.AddCommand(inthandlersUpdateCmd)
	inthandlersCmd.AddCommand(inthandlersRemoveCmd)
}
