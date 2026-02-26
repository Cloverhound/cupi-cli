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

// ast disk command
var astDiskCmd = &cobra.Command{
	Use:   "disk",
	Short: "Show disk partition usage",
	Long:  "Display disk partition information including percentage used, total size, and used space for all partitions.",
	RunE:  runASTDisk,
}

func runASTDisk(cmd *cobra.Command, args []string) error {
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

	partitions, err := client.GetASTDiskInfo(serverCfg.Host, user, pass)
	if err != nil {
		return fmt.Errorf("failed to get disk info: %w", err)
	}

	var rows []map[string]string
	for _, p := range partitions {
		rows = append(rows, map[string]string{
			"node":           p.Node,
			"partition":      p.Name,
			"percentageUsed": strconv.Itoa(p.PercentageUsed),
			"totalMbytes":    strconv.Itoa(p.TotalMbytes),
			"usedMbytes":     strconv.Itoa(p.UsedMbytes),
		})
	}

	return output.Print(rows, outputFlag)
}

// ast tftp command
var astTftpCmd = &cobra.Command{
	Use:   "tftp",
	Short: "Show TFTP server information",
	Long:  "Display TFTP server statistics including total requests and aborted requests.",
	RunE:  runASTTftp,
}

func runASTTftp(cmd *cobra.Command, args []string) error {
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

	tftpInfos, err := client.GetASTTftpInfo(serverCfg.Host, user, pass)
	if err != nil {
		return fmt.Errorf("failed to get TFTP info: %w", err)
	}

	var rows []map[string]string
	for _, t := range tftpInfos {
		rows = append(rows, map[string]string{
			"node":          t.Node,
			"totalRequests": strconv.Itoa(t.TotalRequests),
			"aborted":       strconv.Itoa(t.Aborted),
		})
	}

	return output.Print(rows, outputFlag)
}

// ast heartbeat command
var astHeartbeatCmd = &cobra.Command{
	Use:   "heartbeat",
	Short: "Show heartbeat rates",
	Long:  "Display heartbeat rates for CM nodes and TFTP servers.",
	RunE:  runASTHeartbeat,
}

func runASTHeartbeat(cmd *cobra.Command, args []string) error {
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

	heartbeats, err := client.GetASTHeartbeat(serverCfg.Host, user, pass)
	if err != nil {
		return fmt.Errorf("failed to get heartbeat info: %w", err)
	}

	var rows []map[string]string
	for _, h := range heartbeats {
		rows = append(rows, map[string]string{
			"type": h.Type,
			"node": h.Node,
			"rate": strconv.Itoa(h.Rate),
		})
	}

	return output.Print(rows, outputFlag)
}
