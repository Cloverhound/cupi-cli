package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var handlersGetCmd = &cobra.Command{
	Use:   "get <name-or-objectId>",
	Short: "Get a CUC call handler",
	Args:  cobra.ExactArgs(1),
	RunE:  runHandlersGet,
}

func runHandlersGet(cmd *cobra.Command, args []string) error {
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

	h, err := client.GetCallHandler(serverCfg.Host, serverCfg.Port, user, pass, nameOrID)
	if err != nil {
		return fmt.Errorf("failed to get call handler: %w", err)
	}

	data := map[string]interface{}{
		"objectId":    h.ObjectId,
		"displayName": h.DisplayName,
		"dtmf":        h.DtmfAccessId,
		"isPrimary":   h.IsPrimary,
	}

	return output.Print(data, outputFlag)
}

func init() {
}
