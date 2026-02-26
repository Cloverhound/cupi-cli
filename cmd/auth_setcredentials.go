package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	setCredType     string
	setCredUsername string
	setCredPassword string
)

var authSetCredentialsCmd = &cobra.Command{
	Use:   "set-credentials",
	Short: "Set credentials for a credential type",
	Long: `Add or update credentials (cupi, application, or platform) for a server.

Examples:
  cupi auth set-credentials --type application --username app-user
  cupi auth set-credentials --server lab --type platform --username os-admin --password secret`,
	RunE: runAuthSetCredentials,
}

func runAuthSetCredentials(cmd *cobra.Command, args []string) error {
	if setCredType == "" {
		return fmt.Errorf("--type is required (cupi, application, or platform)")
	}
	if setCredUsername == "" {
		return fmt.Errorf("--username is required")
	}

	server, err := resolveServer(cmd)
	if err != nil {
		return err
	}

	cfg, err := appconfig.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	serverCfg, err := appconfig.GetServer(cfg, server)
	if err != nil {
		return err
	}

	// Validate credential type
	if setCredType != auth.CredTypeCUPI && setCredType != auth.CredTypeApplication && setCredType != auth.CredTypePlatform {
		return fmt.Errorf("invalid credential type: %s (must be cupi, application, or platform)", setCredType)
	}

	// Prompt for password if not provided
	var password string
	if setCredPassword == "" {
		fmt.Fprint(os.Stderr, "Password: ")
		pwBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		password = string(pwBytes)
		fmt.Fprintln(os.Stderr)
	} else {
		password = setCredPassword
	}

	if dryRunFlag {
		fmt.Printf("[DRY RUN] Would save %s username '%s' to config for server '%s'\n", setCredType, setCredUsername, server)
		fmt.Printf("[DRY RUN] Would store %s password in keyring for '%s'\n", setCredType, serverCfg.Host)
		return nil
	}

	// Update config
	if serverCfg.Credentials == nil {
		serverCfg.Credentials = make(map[string]appconfig.CredentialConfig)
	}
	serverCfg.Credentials[setCredType] = appconfig.CredentialConfig{
		Username: setCredUsername,
	}
	cfg.Servers[server] = *serverCfg

	if err := appconfig.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Store password in keyring
	if err := auth.StorePassword(serverCfg.Host, setCredType, password); err != nil {
		return fmt.Errorf("failed to store password: %w", err)
	}

	fmt.Printf("Set %s credentials for server '%s'\n", setCredType, server)
	return nil
}

func init() {
	authSetCredentialsCmd.Flags().StringVar(&setCredType, "type", "", "Credential type (cupi, application, or platform, required)")
	authSetCredentialsCmd.Flags().StringVar(&setCredUsername, "username", "", "Username (required)")
	authSetCredentialsCmd.Flags().StringVar(&setCredPassword, "password", "", "Password (will prompt if not provided)")
}
