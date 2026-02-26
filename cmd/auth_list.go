package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var authListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured servers",
	Long:  `Display all configured CUC servers with their details.`,
	RunE:  runAuthList,
}

func runAuthList(cmd *cobra.Command, args []string) error {
	cfg, err := appconfig.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if len(cfg.Servers) == 0 {
		fmt.Println("No servers configured. Run 'cupi auth login' to get started.")
		return nil
	}

	var data []map[string]string
	for serverKey, server := range cfg.Servers {
		defaultMarker := ""
		if serverKey == cfg.DefaultServer {
			defaultMarker = "*"
		}

		credTypes := []string{}
		for credType := range server.Credentials {
			credTypes = append(credTypes, credType)
		}

		credStr := ""
		for _, ct := range credTypes {
			if credStr != "" {
				credStr += ","
			}
			credStr += ct
		}

		data = append(data, map[string]string{
			"server":  serverKey + defaultMarker,
			"host":    server.Host,
			"port":    fmt.Sprintf("%d", server.Port),
			"version": server.Version,
			"creds":   credStr,
		})
	}

	return output.Print(data, outputFlag)
}

func init() {
}
