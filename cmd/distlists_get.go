package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var distlistsGetCmd = &cobra.Command{
	Use:   "get <alias-or-objectId>",
	Short: "Get a CUC distribution list",
	Args:  cobra.ExactArgs(1),
	RunE:  runDistlistsGet,
}

func runDistlistsGet(cmd *cobra.Command, args []string) error {
	aliasOrID := args[0]

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

	dl, err := client.GetDistList(serverCfg.Host, serverCfg.Port, user, pass, aliasOrID)
	if err != nil {
		return fmt.Errorf("failed to get distribution list: %w", err)
	}

	data := map[string]interface{}{
		"objectId":    dl.ObjectId,
		"alias":       dl.Alias,
		"displayName": dl.DisplayName,
		"dtmf":        dl.DtmfAccessId,
	}

	return output.Print(data, outputFlag)
}

func init() {
}
