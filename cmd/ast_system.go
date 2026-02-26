package cmd

import (
	"fmt"
	"strconv"

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
	srv, user, pass, err := resolveCredentials(cmd, "cupi")
	if err != nil {
		return err
	}

	partitions, err := client.GetASTDiskInfo(srv.Host, user, pass)
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
	srv, user, pass, err := resolveCredentials(cmd, "cupi")
	if err != nil {
		return err
	}

	tftpInfos, err := client.GetASTTftpInfo(srv.Host, user, pass)
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
	srv, user, pass, err := resolveCredentials(cmd, "cupi")
	if err != nil {
		return err
	}

	heartbeats, err := client.GetASTHeartbeat(srv.Host, user, pass)
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
