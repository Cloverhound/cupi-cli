package cmd

import (
	"fmt"
	"strconv"

	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var astTriggeredFlag bool

var astAlertsCmd = &cobra.Command{
	Use:   "alerts",
	Short: "Show system alerts",
	Long: `Display system alerts from the AST API.

Use --triggered to show only currently triggered alerts.
Use the 'get' sub-command to retrieve detailed info for a specific alert.`,
	RunE: runASTAlerts,
}

func runASTAlerts(cmd *cobra.Command, args []string) error {
	srv, user, pass, err := resolveCredentials(cmd, "cupi")
	if err != nil {
		return err
	}

	alerts, err := client.GetASTAlerts(srv.Host, user, pass, astTriggeredFlag)
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

// ast alerts get <alertID> — retrieve detailed info for a specific alert
var astAlertGetCmd = &cobra.Command{
	Use:   "get <alertID>",
	Short: "Get detailed information for a specific alert",
	Long: `Retrieve detailed information for a specific alert by its AlertID.

The AlertID can be found in the output of 'cupi ast alerts'.

Example:
  cupi ast alerts get "HighMboxUsage"
  cupi ast alerts get "LowDiskSpace"`,
	Args: cobra.ExactArgs(1),
	RunE: runASTAlertGet,
}

func runASTAlertGet(cmd *cobra.Command, args []string) error {
	srv, user, pass, err := resolveCredentials(cmd, "cupi")
	if err != nil {
		return err
	}

	detail, err := client.GetASTAlertDetail(srv.Host, user, pass, args[0])
	if err != nil {
		return fmt.Errorf("failed to get alert detail: %w", err)
	}

	rows := []map[string]string{
		{
			"alertID":           detail.AlertID,
			"displayName":       detail.DisplayName,
			"group":             detail.Group,
			"description":       detail.Description,
			"enabled":           strconv.FormatBool(detail.IsEnabled),
			"triggered":         strconv.FormatBool(detail.IsTriggered),
			"withinSafeRange":   strconv.FormatBool(detail.IsWithinSafeRange),
			"thresholdType":     detail.ThresholdType,
			"thresholdValue":    detail.ThresholdValue,
			"severity":          detail.Severity,
		},
	}

	return output.Print(rows, outputFlag)
}

func init() {
	astAlertsCmd.Flags().BoolVar(&astTriggeredFlag, "triggered", false, "Show only triggered alerts")
	astAlertsCmd.AddCommand(astAlertGetCmd)
}
