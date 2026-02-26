package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "Show CUC system information",
	Long:  `Display system information for the CUC server via CUPI REST API.`,
	RunE:  runSystem,
}

func runSystem(cmd *cobra.Command, args []string) error {
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

	info, err := client.GetSystemInfo(serverCfg.Host, serverCfg.Port, user, pass)
	if err != nil {
		return fmt.Errorf("failed to get system info: %w", err)
	}

	data := map[string]interface{}{
		"displayName":    info.DisplayName,
		"version":        info.Version,
		"serialNumber":   info.SerialNumber,
		"ipAddress":      info.IpAddress,
		"hostname":       info.Hostname,
		"domainName":     info.DomainName,
		"smtpSmartHost":  info.SmtpSmartHost,
		"maxMailboxSize": info.MaxMailboxSize,
	}

	return output.Print(data, outputFlag)
}

func init() {
}
