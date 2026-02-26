package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/spf13/cobra"
)

var handlersRemoveCmd = &cobra.Command{
	Use:   "remove <name-or-objectId>",
	Short: "Remove a CUC call handler",
	Args:  cobra.ExactArgs(1),
	RunE:  runHandlersRemove,
}

func runHandlersRemove(cmd *cobra.Command, args []string) error {
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

	if err := client.DeleteCallHandler(serverCfg.Host, serverCfg.Port, user, pass, nameOrID); err != nil {
		return fmt.Errorf("failed to remove call handler: %w", err)
	}

	fmt.Printf("Removed call handler: %s\n", nameOrID)
	return nil
}

func init() {
}
