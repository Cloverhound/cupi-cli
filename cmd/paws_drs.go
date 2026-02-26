package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var pawsDRSCmd = &cobra.Command{
	Use:   "drs",
	Short: "Disaster Recovery System (DRS) backup operations",
	Long:  `Commands for initiating and monitoring DRS backups via PAWS.`,
}

// DRS backup flags
var (
	drsSFTPServer   string
	drsSFTPPort     int
	drsSFTPUser     string
	drsSFTPPassword string
	drsSFTPDir      string
)

var pawsDRSBackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Initiate a DRS backup to an SFTP server",
	Long: `Initiate a Disaster Recovery System (DRS) backup to an SFTP server.

Requires platform credentials (cupi auth set-credentials --type platform ...).

Examples:
  cupi paws drs backup --sftp-server 10.0.0.5 --sftp-user backupuser --sftp-password secret --sftp-dir /backups`,
	RunE: runPAWSDRSBackup,
}

var pawsDRSStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the current DRS backup/restore status",
	Long: `Show the current status of a running or recently completed DRS backup or restore operation.

Requires platform credentials (cupi auth set-credentials --type platform ...).

Examples:
  cupi paws drs status
  cupi paws drs status --output json`,
	RunE: runPAWSDRSStatus,
}

func runPAWSDRSBackup(cmd *cobra.Command, args []string) error {
	if drsSFTPServer == "" {
		return fmt.Errorf("--sftp-server is required")
	}
	if drsSFTPUser == "" {
		return fmt.Errorf("--sftp-user is required")
	}
	if drsSFTPPassword == "" {
		return fmt.Errorf("--sftp-password is required")
	}
	if drsSFTPDir == "" {
		return fmt.Errorf("--sftp-dir is required")
	}

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

	user, pass, err := auth.ResolveCreds(serverCfg, auth.CredTypePlatform)
	if err != nil {
		return fmt.Errorf("failed to resolve platform credentials: %w\n\nSet them with: cupi auth set-credentials --type platform --username <osadmin> --server %s", err, serverName)
	}

	result, err := client.InitiateDRSBackup(serverCfg.Host, user, pass,
		drsSFTPServer, drsSFTPPort, drsSFTPUser, drsSFTPPassword, drsSFTPDir)
	if err != nil {
		return fmt.Errorf("failed to initiate DRS backup: %w", err)
	}

	rows := []map[string]string{
		{
			"result":  result.Result,
			"message": result.Message,
		},
	}

	return output.Print(rows, outputFlag)
}

func runPAWSDRSStatus(cmd *cobra.Command, args []string) error {
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

	user, pass, err := auth.ResolveCreds(serverCfg, auth.CredTypePlatform)
	if err != nil {
		return fmt.Errorf("failed to resolve platform credentials: %w\n\nSet them with: cupi auth set-credentials --type platform --username <osadmin> --server %s", err, serverName)
	}

	status, err := client.GetDRSStatus(serverCfg.Host, user, pass)
	if err != nil {
		return fmt.Errorf("failed to get DRS status: %w", err)
	}

	rows := []map[string]string{
		{
			"status":  status.Status,
			"message": status.Message,
		},
	}

	return output.Print(rows, outputFlag)
}

func init() {
	pawsDRSCmd.AddCommand(pawsDRSBackupCmd)
	pawsDRSCmd.AddCommand(pawsDRSStatusCmd)

	pawsDRSBackupCmd.Flags().StringVar(&drsSFTPServer, "sftp-server", "", "SFTP server hostname or IP (required)")
	pawsDRSBackupCmd.Flags().IntVar(&drsSFTPPort, "sftp-port", 22, "SFTP server port")
	pawsDRSBackupCmd.Flags().StringVar(&drsSFTPUser, "sftp-user", "", "SFTP username (required)")
	pawsDRSBackupCmd.Flags().StringVar(&drsSFTPPassword, "sftp-password", "", "SFTP password (required)")
	pawsDRSBackupCmd.Flags().StringVar(&drsSFTPDir, "sftp-dir", "", "SFTP target directory (required)")
}
