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
	menuAction           string
	menuTargetConv       string
	menuTargetHandler    string
)

var handlersMenuentriesCmd = &cobra.Command{
	Use:   "menuentries",
	Short: "Manage call handler menu entries",
}

var handlersMenuentriesListCmd = &cobra.Command{
	Use:   "list <handler-name-or-id>",
	Short: "List menu entries",
	Args:  cobra.ExactArgs(1),
	RunE:  runHandlersMenuentriesList,
}

var handlersMenuentriesGetCmd = &cobra.Command{
	Use:   "get <handler-name-or-id> <key>",
	Short: "Get a menu entry",
	Args:  cobra.ExactArgs(2),
	RunE:  runHandlersMenuentriesGet,
}

var handlersMenuentriesUpdateCmd = &cobra.Command{
	Use:   "update <handler-name-or-id> <key>",
	Short: "Update a menu entry",
	Args:  cobra.ExactArgs(2),
	RunE:  runHandlersMenuentriesUpdate,
}

func runHandlersMenuentriesList(cmd *cobra.Command, args []string) error {
	handlerNameOrID := args[0]

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

	h, err := client.GetCallHandler(serverCfg.Host, serverCfg.Port, user, pass, handlerNameOrID)
	if err != nil {
		return err
	}

	entries, err := client.ListMenuEntries(serverCfg.Host, serverCfg.Port, user, pass, h.ObjectId)
	if err != nil {
		return err
	}

	var rows []map[string]string
	for _, m := range entries {
		rows = append(rows, map[string]string{
			"touchtoneKey":            m.TouchtoneKey,
			"action":                  m.Action,
			"targetConversation":      m.TargetConversation,
			"targetHandlerObjectId":   m.TargetHandlerObjectId,
		})
	}

	return output.Print(rows, outputFlag)
}

func runHandlersMenuentriesGet(cmd *cobra.Command, args []string) error {
	handlerNameOrID := args[0]
	key := args[1]

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

	h, err := client.GetCallHandler(serverCfg.Host, serverCfg.Port, user, pass, handlerNameOrID)
	if err != nil {
		return err
	}

	m, err := client.GetMenuEntry(serverCfg.Host, serverCfg.Port, user, pass, h.ObjectId, key)
	if err != nil {
		return err
	}

	return output.Print(m, outputFlag)
}

func runHandlersMenuentriesUpdate(cmd *cobra.Command, args []string) error {
	handlerNameOrID := args[0]
	key := args[1]

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

	h, err := client.GetCallHandler(serverCfg.Host, serverCfg.Port, user, pass, handlerNameOrID)
	if err != nil {
		return err
	}

	fields := map[string]interface{}{}
	if menuAction != "" {
		fields["Action"] = menuAction
	}
	if menuTargetConv != "" {
		fields["TargetConversation"] = menuTargetConv
	}
	if menuTargetHandler != "" {
		fields["TargetHandlerObjectId"] = menuTargetHandler
	}

	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	if err := client.UpdateMenuEntry(serverCfg.Host, serverCfg.Port, user, pass, h.ObjectId, key, fields); err != nil {
		return err
	}

	fmt.Printf("Updated menu entry %s for handler %s\n", key, h.DisplayName)
	return nil
}

func init() {
	handlersMenuentriesUpdateCmd.Flags().StringVar(&menuAction, "action", "", "Action")
	handlersMenuentriesUpdateCmd.Flags().StringVar(&menuTargetConv, "target-conversation", "", "Target conversation")
	handlersMenuentriesUpdateCmd.Flags().StringVar(&menuTargetHandler, "target-handler", "", "Target handler ObjectId")

	handlersMenuentriesCmd.AddCommand(handlersMenuentriesListCmd)
	handlersMenuentriesCmd.AddCommand(handlersMenuentriesGetCmd)
	handlersMenuentriesCmd.AddCommand(handlersMenuentriesUpdateCmd)
}
