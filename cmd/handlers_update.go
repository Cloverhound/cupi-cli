package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/spf13/cobra"
)

var (
	updateHandlerDisplayName string
	updateHandlerDtmf        string
)

var handlersUpdateCmd = &cobra.Command{
	Use:   "update <name-or-objectId>",
	Short: "Update a CUC call handler",
	Args:  cobra.ExactArgs(1),
	RunE:  runHandlersUpdate,
}

func runHandlersUpdate(cmd *cobra.Command, args []string) error {
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
	if updateHandlerDisplayName != "" {
		fields["DisplayName"] = updateHandlerDisplayName
	}
	if updateHandlerDtmf != "" {
		fields["DtmfAccessId"] = updateHandlerDtmf
	}

	if len(fields) == 0 {
		return fmt.Errorf("at least one field must be specified to update")
	}

	if err := client.UpdateCallHandler(serverCfg.Host, serverCfg.Port, user, pass, nameOrID, fields); err != nil {
		return fmt.Errorf("failed to update call handler: %w", err)
	}

	fmt.Printf("Updated call handler: %s\n", nameOrID)
	return nil
}

func init() {
	handlersUpdateCmd.Flags().StringVar(&updateHandlerDisplayName, "display-name", "", "Display name")
	handlersUpdateCmd.Flags().StringVar(&updateHandlerDtmf, "dtmf", "", "DTMF access ID")
}
