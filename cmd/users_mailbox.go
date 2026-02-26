package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	mailboxQuotaWarning string
	mailboxQuotaSend    string
	mailboxQuotaReceive string
)

var usersMailboxCmd = &cobra.Command{
	Use:   "mailbox",
	Short: "Manage user mailbox settings",
}

var usersMailboxGetCmd = &cobra.Command{
	Use:   "get <alias-or-objectId>",
	Short: "Get mailbox attributes",
	Args:  cobra.ExactArgs(1),
	RunE:  runUsersMailboxGet,
}

var usersMailboxUpdateCmd = &cobra.Command{
	Use:   "update <alias-or-objectId>",
	Short: "Update mailbox attributes",
	Args:  cobra.ExactArgs(1),
	RunE:  runUsersMailboxUpdate,
}

func runUsersMailboxGet(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]

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

	u, err := client.GetUser(serverCfg.Host, serverCfg.Port, user, pass, userAliasOrID)
	if err != nil {
		return err
	}

	ma, err := client.GetMailboxAttributes(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId)
	if err != nil {
		return err
	}

	return output.Print(ma, outputFlag)
}

func runUsersMailboxUpdate(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]

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

	u, err := client.GetUser(serverCfg.Host, serverCfg.Port, user, pass, userAliasOrID)
	if err != nil {
		return err
	}

	fields := map[string]interface{}{}
	if mailboxQuotaWarning != "" {
		fields["QuotaWarning"] = mailboxQuotaWarning
	}
	if mailboxQuotaSend != "" {
		fields["SendQuota"] = mailboxQuotaSend
	}
	if mailboxQuotaReceive != "" {
		fields["ReceiveQuota"] = mailboxQuotaReceive
	}

	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	if err := client.UpdateMailboxAttributes(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, fields); err != nil {
		return err
	}

	fmt.Printf("Updated mailbox attributes for user %s\n", u.Alias)
	return nil
}

func init() {
	usersMailboxUpdateCmd.Flags().StringVar(&mailboxQuotaWarning, "quota-warning", "", "Quota warning level (in bytes)")
	usersMailboxUpdateCmd.Flags().StringVar(&mailboxQuotaSend, "quota-send", "", "Send quota (in bytes)")
	usersMailboxUpdateCmd.Flags().StringVar(&mailboxQuotaReceive, "quota-receive", "", "Receive quota (in bytes)")

	usersMailboxCmd.AddCommand(usersMailboxGetCmd)
	usersMailboxCmd.AddCommand(usersMailboxUpdateCmd)
}
