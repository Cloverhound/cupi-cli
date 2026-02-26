package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/spf13/cobra"
)

var usersRemoveCmd = &cobra.Command{
	Use:   "remove <alias-or-objectId>",
	Short: "Remove a CUC user",
	Long: `Delete a CUC voicemail user.

Examples:
  cupi users remove jsmith
  cupi --dry-run users remove jsmith`,
	Args: cobra.ExactArgs(1),
	RunE: runUsersRemove,
}

func runUsersRemove(cmd *cobra.Command, args []string) error {
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

	if err := client.DeleteUser(serverCfg.Host, serverCfg.Port, user, pass, aliasOrID); err != nil {
		return fmt.Errorf("failed to remove user: %w", err)
	}

	fmt.Printf("Removed user: %s\n", aliasOrID)
	return nil
}

func init() {
}
