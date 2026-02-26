package cmd

import (
	"fmt"
	"net/url"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	xferAction    string
	xferExtension string
	xferRingCount string
	xferEnabled   string
)

var handlersTransferoptionsCmd = &cobra.Command{
	Use:   "transferoptions",
	Short: "Manage call handler transfer options",
}

var handlersTransferoptionsListCmd = &cobra.Command{
	Use:   "list <handler-name-or-id>",
	Short: "List transfer options",
	Args:  cobra.ExactArgs(1),
	RunE:  runHandlersTransferoptionsList,
}

var handlersTransferoptionsGetCmd = &cobra.Command{
	Use:   "get <handler-name-or-id> <type>",
	Short: "Get a transfer option",
	Args:  cobra.ExactArgs(2),
	RunE:  runHandlersTransferoptionsGet,
}

var handlersTransferoptionsUpdateCmd = &cobra.Command{
	Use:   "update <handler-name-or-id> <type>",
	Short: "Update a transfer option",
	Args:  cobra.ExactArgs(2),
	RunE:  runHandlersTransferoptionsUpdate,
}

func runHandlersTransferoptionsList(cmd *cobra.Command, args []string) error {
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

	options, err := client.ListTransferOptions(serverCfg.Host, serverCfg.Port, user, pass, h.ObjectId)
	if err != nil {
		return err
	}

	var rows []map[string]string
	for _, t := range options {
		rows = append(rows, map[string]string{
			"transferType": t.TransferType,
			"action":       t.Action,
			"extension":    t.Extension,
			"ringCount":    t.RingCount,
			"enabled":      t.Enabled,
		})
	}

	return output.Print(rows, outputFlag)
}

func runHandlersTransferoptionsGet(cmd *cobra.Command, args []string) error {
	handlerNameOrID := args[0]
	transferType := args[1]

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

	t, err := client.GetTransferOption(serverCfg.Host, serverCfg.Port, user, pass, h.ObjectId, transferType)
	if err != nil {
		return err
	}

	return output.Print(t, outputFlag)
}

func runHandlersTransferoptionsUpdate(cmd *cobra.Command, args []string) error {
	handlerNameOrID := args[0]
	transferType := args[1]

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
	if xferAction != "" {
		fields["Action"] = xferAction
	}
	if xferExtension != "" {
		fields["Extension"] = xferExtension
	}
	if xferRingCount != "" {
		fields["RingCount"] = xferRingCount
	}
	if xferEnabled != "" {
		fields["Enabled"] = xferEnabled
	}

	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	if err := client.UpdateTransferOption(serverCfg.Host, serverCfg.Port, user, pass, h.ObjectId, url.PathEscape(transferType), fields); err != nil {
		return err
	}

	fmt.Printf("Updated transfer option %s for handler %s\n", transferType, h.DisplayName)
	return nil
}

func init() {
	handlersTransferoptionsUpdateCmd.Flags().StringVar(&xferAction, "action", "", "Action")
	handlersTransferoptionsUpdateCmd.Flags().StringVar(&xferExtension, "extension", "", "Extension")
	handlersTransferoptionsUpdateCmd.Flags().StringVar(&xferRingCount, "ring-count", "", "Ring count")
	handlersTransferoptionsUpdateCmd.Flags().StringVar(&xferEnabled, "enabled", "", "Enabled (true|false)")

	handlersTransferoptionsCmd.AddCommand(handlersTransferoptionsListCmd)
	handlersTransferoptionsCmd.AddCommand(handlersTransferoptionsGetCmd)
	handlersTransferoptionsCmd.AddCommand(handlersTransferoptionsUpdateCmd)
}
