package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/spf13/cobra"
)

var (
	addDLAlias       string
	addDLDisplayName string
	addDLDtmf        string
)

var distlistsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a CUC distribution list",
	RunE:  runDistlistsAdd,
}

func runDistlistsAdd(cmd *cobra.Command, args []string) error {
	if addDLAlias == "" {
		return fmt.Errorf("--alias is required")
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
		"Alias": addDLAlias,
	}
	if addDLDisplayName != "" {
		fields["DisplayName"] = addDLDisplayName
	}
	if addDLDtmf != "" {
		fields["DtmfAccessId"] = addDLDtmf
	}

	dl, err := client.CreateDistList(serverCfg.Host, serverCfg.Port, user, pass, fields)
	if err != nil {
		return fmt.Errorf("failed to create distribution list: %w", err)
	}

	fmt.Printf("Created distribution list: alias=%s objectId=%s\n", dl.Alias, dl.ObjectId)
	return nil
}

func init() {
	distlistsAddCmd.Flags().StringVar(&addDLAlias, "alias", "", "List alias (required)")
	distlistsAddCmd.Flags().StringVar(&addDLDisplayName, "display-name", "", "Display name")
	distlistsAddCmd.Flags().StringVar(&addDLDtmf, "dtmf", "", "DTMF access ID")
}
