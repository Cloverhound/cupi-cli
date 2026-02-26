package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/spf13/cobra"
)

var (
	updateDLDisplayName string
	updateDLDtmf        string
)

var distlistsUpdateCmd = &cobra.Command{
	Use:   "update <alias-or-objectId>",
	Short: "Update a CUC distribution list",
	Args:  cobra.ExactArgs(1),
	RunE:  runDistlistsUpdate,
}

func runDistlistsUpdate(cmd *cobra.Command, args []string) error {
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

	fields := map[string]interface{}{}
	if updateDLDisplayName != "" {
		fields["DisplayName"] = updateDLDisplayName
	}
	if updateDLDtmf != "" {
		fields["DtmfAccessId"] = updateDLDtmf
	}

	if len(fields) == 0 {
		return fmt.Errorf("at least one field must be specified to update")
	}

	if err := client.UpdateDistList(serverCfg.Host, serverCfg.Port, user, pass, aliasOrID, fields); err != nil {
		return fmt.Errorf("failed to update distribution list: %w", err)
	}

	fmt.Printf("Updated distribution list: %s\n", aliasOrID)
	return nil
}

func init() {
	distlistsUpdateCmd.Flags().StringVar(&updateDLDisplayName, "display-name", "", "Display name")
	distlistsUpdateCmd.Flags().StringVar(&updateDLDtmf, "dtmf", "", "DTMF access ID")
}
