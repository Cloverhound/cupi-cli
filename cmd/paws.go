package cmd

import (
	"github.com/spf13/cobra"
)

// pawsCmd is the parent for all PAWS (Platform Administrative Web Service) subcommands.
var pawsCmd = &cobra.Command{
	Use:   "paws",
	Short: "Platform Administrative Web Service (OS-level operations)",
	Long: `Commands for CUC OS-level administration via the Platform Administrative Web Service (PAWS).

PAWS uses platform (OS admin) credentials, which can be set with:
  cupi auth set-credentials --type platform --username <osadmin> --server <name>

Available subcommand groups:
  cupi paws cluster   Cluster node status and replication health
  cupi paws drs       Disaster Recovery System (DRS) backup operations`,
}

func init() {
	pawsCmd.AddCommand(pawsClusterCmd)
	pawsCmd.AddCommand(pawsDRSCmd)
}
