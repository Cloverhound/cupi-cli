package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/spf13/cobra"
)

var (
	dimeOutputFlag string
	dimeNodeFlag   string
)

// dimeCmd is the parent for DIME log collection subcommands.
var dimeCmd = &cobra.Command{
	Use:   "dime",
	Short: "Download log files via DIME (Direct Internet Message Encapsulation)",
	Long: `Commands for downloading log files from CUC nodes via the DIME log collection service.

DIME is CUC's binary SOAP attachment protocol for transferring log files from /var/log.
It uses CUPI credentials (the default credential type).

Available subcommands:
  cupi dime get-file   Download a log file from a CUC node`,
}

var dimeGetFileCmd = &cobra.Command{
	Use:   "get-file <server-path>",
	Short: "Download a log file from a CUC node",
	Long: `Download a log file from a CUC node using the DIME log collection service.

The server-path must be an absolute path on the CUC node.
If the path does not begin with /var/log, /var/log/active/ is prepended automatically.

Use --node to target a specific node; by default the publisher is used.
Writes file bytes to stdout by default (redirect with >) or to a file with --output.

Examples:
  cupi dime get-file /var/log/active/syslog/CiscoSyslog > syslog.txt
  cupi dime get-file /var/log/active/tomcat/catalina.out --output /tmp/catalina.out
  cupi dime get-file syslog/CiscoSyslog --node 172.20.1.5`,
	Args: cobra.ExactArgs(1),
	RunE: runDimeGetFile,
}

func runDimeGetFile(cmd *cobra.Command, args []string) error {
	serverPath := args[0]

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
		return fmt.Errorf("failed to resolve CUPI credentials: %w", err)
	}

	host := serverCfg.Host
	if dimeNodeFlag != "" {
		host = dimeNodeFlag
	}

	data, err := client.GetFile(host, user, pass, serverPath)
	if err != nil {
		return fmt.Errorf("DIME get-file failed: %w", err)
	}

	if dimeOutputFlag != "" {
		if err := os.WriteFile(dimeOutputFlag, data, 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Wrote %d bytes to %s\n", len(data), dimeOutputFlag)
		return nil
	}

	if _, err := os.Stdout.Write(data); err != nil {
		return fmt.Errorf("failed to write to stdout: %w", err)
	}

	localName := filepath.Base(serverPath)
	fmt.Fprintf(os.Stderr, "# %d bytes written to stdout (tip: redirect with > %s)\n", len(data), localName)

	return nil
}

func init() {
	dimeGetFileCmd.Flags().StringVar(&dimeOutputFlag, "output", "", "Local file path to write the downloaded content (default: stdout)")
	dimeGetFileCmd.Flags().StringVar(&dimeNodeFlag, "node", "", "Specific node IP/hostname to target (default: publisher)")
	dimeCmd.AddCommand(dimeGetFileCmd)
}
