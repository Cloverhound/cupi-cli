package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/spf13/cobra"
)

var (
	logoutType string
)

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove server credentials",
	Long: `Remove credentials for a server from keyring and config.

If --type is not specified, removes the entire server and all credentials.
If --type is specified, removes only that credential type.

Examples:
  cupi auth logout --server prod              # Remove all prod credentials
  cupi auth logout --server prod --type cupi  # Remove only CUPI credentials`,
	RunE: runAuthLogout,
}

func runAuthLogout(cmd *cobra.Command, args []string) error {
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

	if logoutType == "" || logoutType == "all" {
		if dryRunFlag {
			fmt.Printf("[DRY RUN] Would remove server '%s' and all credentials from keyring and config\n", server)
			return nil
		}

		if err := auth.DeleteAllPasswords(serverCfg.Host); err != nil {
			return fmt.Errorf("failed to delete passwords: %w", err)
		}

		delete(cfg.Servers, server)

		if cfg.DefaultServer == server {
			cfg.DefaultServer = ""
			if len(cfg.Servers) > 0 {
				for key := range cfg.Servers {
					cfg.DefaultServer = key
					break
				}
			}
		}

		fmt.Printf("Removed server '%s' and all credentials\n", server)
	} else {
		credTypes := []string{auth.CredTypeCUPI, auth.CredTypeApplication, auth.CredTypePlatform}
		found := false
		for _, ct := range credTypes {
			if ct == logoutType {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("invalid credential type: %s (must be cupi, application, or platform)", logoutType)
		}

		if dryRunFlag {
			fmt.Printf("[DRY RUN] Would remove %s credentials for server '%s' from keyring and config\n", logoutType, server)
			return nil
		}

		if err := auth.DeletePassword(serverCfg.Host, logoutType); err != nil {
			return fmt.Errorf("failed to delete password: %w", err)
		}

		if serverCfg.Credentials != nil {
			delete(serverCfg.Credentials, logoutType)
		}
		cfg.Servers[server] = *serverCfg

		fmt.Printf("Removed %s credentials for server '%s'\n", logoutType, server)
	}

	if err := appconfig.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

func init() {
	authLogoutCmd.Flags().StringVar(&logoutType, "type", "", "Credential type to remove (cupi|application|platform|all, default: all)")
}
