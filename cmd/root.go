package cmd

import (
	"fmt"
	"os"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/spf13/cobra"
)

var (
	serverFlag  string
	outputFlag  string
	debugFlag   bool
	maxFlag     int
	dryRunFlag  bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cupi-cli",
	Short: "Cisco Unity Connection CLI management tool",
	Long: `cupi-cli is a command-line tool for querying and managing Cisco Unity Connection (CUC) servers.
It provides access to CUPI REST, PAWS, AST, and DIME APIs for voicemail user, distribution list, and system management.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debugFlag {
			os.Setenv("CUPI_DEBUG", "1")
		}
		if dryRunFlag {
			os.Setenv("CUPI_DRY_RUN", "1")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&serverFlag, "server", "", "Server to use (overrides default)")
	rootCmd.PersistentFlags().StringVar(&outputFlag, "output", "table", "Output format (json, table, csv, raw)")
	rootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, "Print API requests/responses")
	rootCmd.PersistentFlags().IntVar(&maxFlag, "max", 0, "Maximum number of results to return (0 = no limit)")
	rootCmd.PersistentFlags().BoolVar(&dryRunFlag, "dry-run", false, "Print what would be sent without making any changes")

	// Add subcommands
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(usersCmd)
	rootCmd.AddCommand(distlistsCmd)
	rootCmd.AddCommand(handlersCmd)
	rootCmd.AddCommand(cosCmd)
	rootCmd.AddCommand(templatesCmd)
	rootCmd.AddCommand(schedulesCmd)
	rootCmd.AddCommand(systemCmd)
	rootCmd.AddCommand(pawsCmd)
	rootCmd.AddCommand(astCmd)
	rootCmd.AddCommand(dimeCmd)
}

// resolveServer returns the server name to use.
// Priority: --server flag > default server in config > error
func resolveServer(cmd *cobra.Command) (string, error) {
	if serverFlag != "" {
		return serverFlag, nil
	}

	cfg, err := loadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.DefaultServer == "" {
		return "", fmt.Errorf("no default server configured; use --server flag or run 'cupi-cli auth login'")
	}

	return cfg.DefaultServer, nil
}

func loadConfig() (*appconfig.Config, error) {
	return appconfig.LoadConfig()
}
