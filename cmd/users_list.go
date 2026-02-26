package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var usersListQueryFlag string

var usersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List CUC users",
	Long: `List CUC voicemail users. Use --query to filter results.

Examples:
  cupi users list
  cupi users list --max 50
  cupi users list --query "(alias startswith j)"
  cupi users list --output json | jq '.[0]'`,
	RunE: runUsersList,
}

func runUsersList(cmd *cobra.Command, args []string) error {
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

	users, err := client.ListUsers(serverCfg.Host, serverCfg.Port, user, pass, usersListQueryFlag, maxFlag)
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}

	var rows []map[string]string
	for _, u := range users {
		rows = append(rows, map[string]string{
			"alias":       u.Alias,
			"displayName": u.DisplayName,
			"dtmf":        u.DtmfAccessId,
			"firstName":   u.FirstName,
			"lastName":    u.LastName,
			"objectId":    u.ObjectId,
		})
	}

	return output.Print(rows, outputFlag)
}

func init() {
	usersListCmd.Flags().StringVar(&usersListQueryFlag, "query", "", "CUPI query filter (e.g. \"(alias startswith j)\")")
}
