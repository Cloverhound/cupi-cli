package cmd

import (
	"github.com/spf13/cobra"
)

var astCmd = &cobra.Command{
	Use:   "ast",
	Short: "AST system health and performance monitoring",
	Long: `Access AST (Application Server Task) API for system health and performance monitoring.

Provides disk usage, TFTP stats, heartbeat monitoring, alerts, and perfmon counters.
Uses CUPI credentials (the default credential type).`,
}

func init() {
	astCmd.AddCommand(astDiskCmd)
	astCmd.AddCommand(astTftpCmd)
	astCmd.AddCommand(astHeartbeatCmd)
	astCmd.AddCommand(astAlertsCmd)
	astCmd.AddCommand(astPerfmonCmd)
}
