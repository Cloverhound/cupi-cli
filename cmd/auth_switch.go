package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/spf13/cobra"
)

var authSwitchCmd = &cobra.Command{
	Use:   "switch <server>",
	Short: "Switch default server",
	Long:  `Set the default server for subsequent commands.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runAuthSwitch,
}

func runAuthSwitch(cmd *cobra.Command, args []string) error {
	serverName := args[0]

	cfg, err := appconfig.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if _, ok := cfg.Servers[serverName]; !ok {
		return fmt.Errorf("server '%s' not found", serverName)
	}

	if dryRunFlag {
		fmt.Printf("[DRY RUN] Would set default server to '%s' in ~/.cupi-cli/config.json\n", serverName)
		return nil
	}

	cfg.DefaultServer = serverName

	if err := appconfig.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Default server set to '%s'\n", serverName)
	return nil
}

func init() {
}
