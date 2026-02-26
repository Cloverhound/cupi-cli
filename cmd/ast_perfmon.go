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

var astPerfmonCmd = &cobra.Command{
	Use:   "perfmon",
	Short: "Show perfmon objects and counters",
	Long:  "Display available perfmon (performance monitor) objects and their associated counters.",
	RunE:  runASTPerfmon,
}

func runASTPerfmon(cmd *cobra.Command, args []string) error {
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

	objects, err := client.GetASTPerfmonObjects(serverCfg.Host, user, pass)
	if err != nil {
		return fmt.Errorf("failed to get perfmon objects: %w", err)
	}

	var rows []map[string]string
	for _, obj := range objects {
		rows = append(rows, map[string]string{
			"host":         obj.Host,
			"object":       obj.ObjectName,
			"hasInstances": strconv.FormatBool(obj.HasInstances),
			"counterCount": strconv.Itoa(obj.CounterCount),
		})
	}

	return output.Print(rows, outputFlag)
}
