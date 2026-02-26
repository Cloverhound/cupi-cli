package cmd

import (
	"fmt"
	"strconv"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var astTriggeredFlag bool

var astAlertsCmd = &cobra.Command{
	Use:   "alerts",
	Short: "Show system alerts",
	Long: `Display system alerts from the AST API.

Use --triggered to show only currently triggered alerts.`,
	RunE: runASTAlerts,
}

func runASTAlerts(cmd *cobra.Command, args []string) error {
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

	alerts, err := client.GetASTAlerts(serverCfg.Host, user, pass, astTriggeredFlag)
	if err != nil {
		return fmt.Errorf("failed to get alerts: %w", err)
	}

	var rows []map[string]string
	for _, a := range alerts {
		rows = append(rows, map[string]string{
			"alertID":         a.AlertID,
			"displayName":     a.DisplayName,
			"group":           a.Group,
			"triggered":       strconv.FormatBool(a.IsTriggered),
			"enabled":         strconv.FormatBool(a.IsEnabled),
			"withinSafeRange": strconv.FormatBool(a.IsWithinSafeRange),
			"timestamp":       a.Timestamp,
		})
	}

	return output.Print(rows, outputFlag)
}

func init() {
	astAlertsCmd.Flags().BoolVar(&astTriggeredFlag, "triggered", false, "Show only triggered alerts")
}
