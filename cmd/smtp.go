package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	smtpSmartHost  string
	smtpPort       string
	smtpUseSsl     string
	smtpServerName string
	smtpUseAuth    string
)

var smtpCmd = &cobra.Command{
	Use:   "smtp",
	Short: "Manage SMTP configuration",
}

var smtpServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage SMTP server (inbound) configuration",
}

var smtpServerGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get SMTP server configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		item, err := client.GetSMTPServerConfig(srv.Host, srv.Port, cupiUser, cupiPass)
		if err != nil {
			return err
		}
		return output.Print(item, outputFlag)
	},
}

var smtpServerUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update SMTP server configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{}
		if smtpSmartHost != "" {
			fields["SmartHost"] = smtpSmartHost
		}
		if smtpPort != "" {
			fields["Port"] = smtpPort
		}
		if smtpUseSsl != "" {
			fields["UseSsl"] = smtpUseSsl
		}
		if len(fields) == 0 {
			return fmt.Errorf("no fields to update; use --help to see available flags")
		}
		if err := client.UpdateSMTPServerConfig(srv.Host, srv.Port, cupiUser, cupiPass, fields); err != nil {
			return err
		}
		fmt.Println("Updated SMTP server configuration")
		return nil
	},
}

var smtpClientCmd = &cobra.Command{
	Use:   "client",
	Short: "Manage SMTP client (outbound) configuration",
}

var smtpClientGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get SMTP client configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		item, err := client.GetSMTPClientConfig(srv.Host, srv.Port, cupiUser, cupiPass)
		if err != nil {
			return err
		}
		return output.Print(item, outputFlag)
	},
}

var smtpClientUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update SMTP client configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{}
		if smtpServerName != "" {
			fields["ServerName"] = smtpServerName
		}
		if smtpPort != "" {
			fields["Port"] = smtpPort
		}
		if smtpUseAuth != "" {
			fields["UseAuth"] = smtpUseAuth
		}
		if len(fields) == 0 {
			return fmt.Errorf("no fields to update; use --help to see available flags")
		}
		if err := client.UpdateSMTPClientConfig(srv.Host, srv.Port, cupiUser, cupiPass, fields); err != nil {
			return err
		}
		fmt.Println("Updated SMTP client configuration")
		return nil
	},
}

func init() {
	smtpServerUpdateCmd.Flags().StringVar(&smtpSmartHost, "smart-host", "", "SMTP smart host")
	smtpServerUpdateCmd.Flags().StringVar(&smtpPort, "port", "", "SMTP port")
	smtpServerUpdateCmd.Flags().StringVar(&smtpUseSsl, "use-ssl", "", "Use SSL (true/false)")

	smtpClientUpdateCmd.Flags().StringVar(&smtpServerName, "server-name", "", "SMTP server hostname")
	smtpClientUpdateCmd.Flags().StringVar(&smtpPort, "port", "", "SMTP port")
	smtpClientUpdateCmd.Flags().StringVar(&smtpUseAuth, "use-auth", "", "Use authentication (true/false)")

	smtpServerCmd.AddCommand(smtpServerGetCmd)
	smtpServerCmd.AddCommand(smtpServerUpdateCmd)

	smtpClientCmd.AddCommand(smtpClientGetCmd)
	smtpClientCmd.AddCommand(smtpClientUpdateCmd)

	smtpCmd.AddCommand(smtpServerCmd)
	smtpCmd.AddCommand(smtpClientCmd)
}
