package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/spf13/cobra"
)

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	Long:  `Display configured servers and their credential status.`,
	RunE:  runAuthStatus,
}

func runAuthStatus(cmd *cobra.Command, args []string) error {
	cfg, err := appconfig.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if len(cfg.Servers) == 0 {
		fmt.Println("No servers configured. Run 'cupi auth login' to get started.")
		return nil
	}

	for serverKey, server := range cfg.Servers {
		marker := ""
		if serverKey == cfg.DefaultServer {
			marker = " (default)"
		}

		fmt.Printf("server=%s%s  host=%s  port=%d  version=%s\n",
			serverKey, marker, server.Host, server.Port, server.Version)

		credTypes := []string{auth.CredTypeCUPI, auth.CredTypeApplication, auth.CredTypePlatform}
		for _, credType := range credTypes {
			cred, ok := server.Credentials[credType]
			username := ""
			if ok {
				username = cred.Username
			}

			status := "[not set]"
			_, err := auth.GetPassword(server.Host, credType)
			if err == nil {
				status = "[set]"
			}

			fmt.Printf("  %-11s: %-20s %s\n", credType, username, status)
		}
		fmt.Println()
	}

	return nil
}

func init() {
}
