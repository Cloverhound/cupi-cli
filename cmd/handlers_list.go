package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var handlersListQueryFlag string

var handlersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List CUC call handlers",
	RunE:  runHandlersList,
}

func runHandlersList(cmd *cobra.Command, args []string) error {
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

	handlers, err := client.ListCallHandlers(serverCfg.Host, serverCfg.Port, user, pass, handlersListQueryFlag, maxFlag)
	if err != nil {
		return fmt.Errorf("failed to list call handlers: %w", err)
	}

	var rows []map[string]string
	for _, h := range handlers {
		rows = append(rows, map[string]string{
			"displayName": h.DisplayName,
			"dtmf":        h.DtmfAccessId,
			"isPrimary":   h.IsPrimary,
			"objectId":    h.ObjectId,
		})
	}

	return output.Print(rows, outputFlag)
}

func init() {
	handlersListCmd.Flags().StringVar(&handlersListQueryFlag, "query", "", "CUPI query filter")
}
