package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var pawsClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Cluster node status and replication health",
	Long:  `Commands for querying CUC cluster node status and replication health via PAWS.`,
}

var pawsClusterStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show PAWS cluster node status",
	Long: `Show OS-level cluster node status for all nodes.

Requires platform credentials (cupi auth set-credentials --type platform ...).

Examples:
  cupi paws cluster status
  cupi paws cluster status --output json`,
	RunE: runPAWSClusterStatus,
}

var pawsClusterReplicationCmd = &cobra.Command{
	Use:   "replication",
	Short: "Check cluster replication health",
	Long: `Check whether cluster database replication is healthy via PAWS.

Requires platform credentials (cupi auth set-credentials --type platform ...).

Examples:
  cupi paws cluster replication`,
	RunE: runPAWSClusterReplication,
}

func runPAWSClusterStatus(cmd *cobra.Command, args []string) error {
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

	nodes, err := client.GetClusterStatus(serverCfg.Host, user, pass)
	if err != nil {
		return fmt.Errorf("failed to get PAWS cluster status: %w", err)
	}

	if len(nodes) == 0 {
		fmt.Println("No cluster nodes returned.")
		return nil
	}

	var rows []map[string]string
	for _, n := range nodes {
		rows = append(rows, map[string]string{
			"hostname": n.Hostname,
			"address":  n.Address,
			"status":   n.Status,
			"dbRole":   n.DBRole,
			"type":     n.Type,
		})
	}

	return output.Print(rows, outputFlag)
}

func runPAWSClusterReplication(cmd *cobra.Command, args []string) error {
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

	status, err := client.GetReplicationStatus(serverCfg.Host, user, pass)
	if err != nil {
		return fmt.Errorf("failed to get PAWS replication status: %w", err)
	}

	rows := []map[string]string{
		{
			"replicationOK": fmt.Sprintf("%v", status.OK),
		},
	}

	return output.Print(rows, outputFlag)
}

func init() {
	pawsClusterCmd.AddCommand(pawsClusterStatusCmd)
	pawsClusterCmd.AddCommand(pawsClusterReplicationCmd)
}
