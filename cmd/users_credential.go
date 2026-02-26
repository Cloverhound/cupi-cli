package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var credentialValue string

var usersCredentialCmd = &cobra.Command{
	Use:   "credential",
	Short: "Manage user credentials",
}

var usersCredentialGetCmd = &cobra.Command{
	Use:   "get <alias-or-objectId> <pin|password>",
	Short: "Get credential info",
	Args:  cobra.ExactArgs(2),
	RunE:  runUsersCredentialGet,
}

var usersCredentialUnlockCmd = &cobra.Command{
	Use:   "unlock <alias-or-objectId> <pin|password>",
	Short: "Unlock a credential",
	Args:  cobra.ExactArgs(2),
	RunE:  runUsersCredentialUnlock,
}

var usersCredentialSetCmd = &cobra.Command{
	Use:   "set <alias-or-objectId> <pin|password>",
	Short: "Set a new credential value",
	Args:  cobra.ExactArgs(2),
	RunE:  runUsersCredentialSet,
}

func runUsersCredentialGet(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]
	credType := args[1]

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

	cred, err := client.GetCredential(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, credType)
	if err != nil {
		return err
	}

	return output.Print(cred, outputFlag)
}

func runUsersCredentialUnlock(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]
	credType := args[1]

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

	if err := client.UnlockCredential(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, credType); err != nil {
		return err
	}

	fmt.Printf("Unlocked %s credential for user %s\n", credType, u.Alias)
	return nil
}

func runUsersCredentialSet(cmd *cobra.Command, args []string) error {
	userAliasOrID := args[0]
	credType := args[1]

	if credentialValue == "" {
		return fmt.Errorf("--value is required")
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

	user, pass, err := auth.ResolveCreds(serverCfg, auth.CredTypeCUPI)
	if err != nil {
		return fmt.Errorf("failed to resolve credentials: %w", err)
	}

	u, err := client.GetUser(serverCfg.Host, serverCfg.Port, user, pass, userAliasOrID)
	if err != nil {
		return err
	}

	if err := client.SetCredential(serverCfg.Host, serverCfg.Port, user, pass, u.ObjectId, credType, credentialValue); err != nil {
		return err
	}

	fmt.Printf("Set %s credential for user %s\n", credType, u.Alias)
	return nil
}

func init() {
	usersCredentialSetCmd.Flags().StringVar(&credentialValue, "value", "", "New credential value (required)")

	usersCredentialCmd.AddCommand(usersCredentialGetCmd)
	usersCredentialCmd.AddCommand(usersCredentialUnlockCmd)
	usersCredentialCmd.AddCommand(usersCredentialSetCmd)
}
