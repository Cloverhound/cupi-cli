package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var schedulesCmd = &cobra.Command{
	Use:   "schedules",
	Short: "Manage CUC schedules",
	Long:  `Commands for listing and viewing CUC schedules.`,
}

var schedulesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List schedules",
	RunE:  runSchedulesList,
}

var schedulesGetCmd = &cobra.Command{
	Use:   "get <name-or-objectId>",
	Short: "Get a schedule",
	Args:  cobra.ExactArgs(1),
	RunE:  runSchedulesGet,
}

func runSchedulesList(cmd *cobra.Command, args []string) error {
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

	schedules, err := client.ListSchedules(serverCfg.Host, serverCfg.Port, user, pass, "", maxFlag)
	if err != nil {
		return fmt.Errorf("failed to list schedules: %w", err)
	}

	var rows []map[string]string
	for _, s := range schedules {
		rows = append(rows, map[string]string{
			"displayName": s.DisplayName,
			"isHoliday":   s.IsHoliday,
			"objectId":    s.ObjectId,
		})
	}

	return output.Print(rows, outputFlag)
}

func runSchedulesGet(cmd *cobra.Command, args []string) error {
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

	s, err := client.GetSchedule(serverCfg.Host, serverCfg.Port, user, pass, nameOrID)
	if err != nil {
		return fmt.Errorf("failed to get schedule: %w", err)
	}

	data := map[string]interface{}{
		"objectId":    s.ObjectId,
		"displayName": s.DisplayName,
		"isHoliday":   s.IsHoliday,
	}

	return output.Print(data, outputFlag)
}

func init() {
	schedulesCmd.AddCommand(schedulesListCmd)
	schedulesCmd.AddCommand(schedulesGetCmd)
}
