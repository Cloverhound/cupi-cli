package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var usersGetCmd = &cobra.Command{
	Use:   "get <alias-or-objectId>",
	Short: "Get a CUC user",
	Long: `Get details for a CUC voicemail user by alias or ObjectId.

Examples:
  cupi users get jsmith
  cupi users get 12345678-1234-1234-1234-123456789abc
  cupi users get jsmith --output json`,
	Args: cobra.ExactArgs(1),
	RunE: runUsersGet,
}

func runUsersGet(cmd *cobra.Command, args []string) error {
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

	u, err := client.GetUser(serverCfg.Host, serverCfg.Port, user, pass, aliasOrID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	data := map[string]interface{}{
		"objectId":    u.ObjectId,
		"alias":       u.Alias,
		"displayName": u.DisplayName,
		"dtmf":        u.DtmfAccessId,
		"firstName":   u.FirstName,
		"lastName":    u.LastName,
		"department":  u.Department,
		"title":       u.Title,
	}

	return output.Print(data, outputFlag)
}

func init() {
}
