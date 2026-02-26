package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var distlistsListQueryFlag string

var distlistsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List CUC distribution lists",
	RunE:  runDistlistsList,
}

func runDistlistsList(cmd *cobra.Command, args []string) error {
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

	lists, err := client.ListDistLists(serverCfg.Host, serverCfg.Port, user, pass, distlistsListQueryFlag, maxFlag)
	if err != nil {
		return fmt.Errorf("failed to list distribution lists: %w", err)
	}

	var rows []map[string]string
	for _, dl := range lists {
		rows = append(rows, map[string]string{
			"alias":       dl.Alias,
			"displayName": dl.DisplayName,
			"dtmf":        dl.DtmfAccessId,
			"objectId":    dl.ObjectId,
		})
	}

	return output.Print(rows, outputFlag)
}

func init() {
	distlistsListCmd.Flags().StringVar(&distlistsListQueryFlag, "query", "", "CUPI query filter")
}
