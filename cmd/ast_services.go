package cmd

import (
	"fmt"
	"os"

	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var astServicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Manage CUC services",
	Long: `List and manage Cisco Unity Connection VOS services via the AST API.

Supports listing all services and performing Start, Stop, or Restart actions.
Service actions require CUPI admin credentials and will affect system operation.`,
}

// ast services list
var astServicesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all services and their status",
	Long:  "Display all VOS services with their current status and startup type.",
	RunE:  runASTServicesList,
}

func runASTServicesList(cmd *cobra.Command, args []string) error {
	srv, user, pass, err := resolveCredentials(cmd, "cupi")
	if err != nil {
		return err
	}

	services, err := client.GetASTServiceList(srv.Host, user, pass)
	if err != nil {
		return fmt.Errorf("failed to get service list: %w", err)
	}

	var rows []map[string]string
	for _, s := range services {
		rows = append(rows, map[string]string{
			"serviceName":   s.ServiceName,
			"serviceStatus": s.ServiceStatus,
			"startupType":   s.StartupType,
			"reasonCode":    s.ReasonCode,
			"nodeName":      s.NodeName,
		})
	}

	return output.Print(rows, outputFlag)
}

// ast services restart <service>
var astServicesRestartCmd = &cobra.Command{
	Use:   "restart <serviceName>",
	Short: "Restart a service",
	Long: `Restart the named VOS service on the CUC node.

Example:
  cupi ast services restart "Cisco Unity Connection Voicemail"
  cupi ast services restart "Cisco Tomcat"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runASTServiceAction(cmd, args[0], "Restart")
	},
}

// ast services start <service>
var astServicesStartCmd = &cobra.Command{
	Use:   "start <serviceName>",
	Short: "Start a service",
	Long: `Start the named VOS service on the CUC node.

Example:
  cupi ast services start "Cisco Unity Connection Voicemail"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runASTServiceAction(cmd, args[0], "Start")
	},
}

// ast services stop <service>
var astServicesStopCmd = &cobra.Command{
	Use:   "stop <serviceName>",
	Short: "Stop a service",
	Long: `Stop the named VOS service on the CUC node.

Example:
  cupi ast services stop "Cisco Unity Connection Voicemail"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runASTServiceAction(cmd, args[0], "Stop")
	},
}

func runASTServiceAction(cmd *cobra.Command, serviceName, action string) error {
	if os.Getenv("CUPI_DRY_RUN") != "" {
		fmt.Printf("[dry-run] Would %s service: %s\n", action, serviceName)
		return nil
	}

	srv, user, pass, err := resolveCredentials(cmd, "cupi")
	if err != nil {
		return err
	}

	if err := client.DoASTServiceAction(srv.Host, user, pass, serviceName, action); err != nil {
		return fmt.Errorf("failed to %s service %q: %w", action, serviceName, err)
	}

	fmt.Printf("%s action sent for service: %s\n", action, serviceName)
	return nil
}

func init() {
	astServicesCmd.AddCommand(astServicesListCmd)
	astServicesCmd.AddCommand(astServicesRestartCmd)
	astServicesCmd.AddCommand(astServicesStartCmd)
	astServicesCmd.AddCommand(astServicesStopCmd)
}
