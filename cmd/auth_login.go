package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/discovery"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	loginHost     string
	loginUsername string
	loginPassword string
	loginServer   string
	loginPort     int
	loginDefault  bool
)

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate to a CUC server",
	Long: `Login to a Cisco Unity Connection server and save configuration.

This command will:
1. Test connectivity to the CUC host using CUPI REST API
2. Detect the CUC version (if available)
3. Validate CUPI credentials
4. Save server configuration
5. Store password in system keyring

Example:
  cupi auth login --host cuc.example.com --username admin --server prod --default`,
	RunE: runAuthLogin,
}

func runAuthLogin(cmd *cobra.Command, args []string) error {
	if loginHost == "" {
		return fmt.Errorf("--host is required")
	}
	if loginUsername == "" {
		return fmt.Errorf("--username is required")
	}

	// Prompt for password if not provided
	var password string
	if loginPassword == "" {
		fmt.Fprint(os.Stderr, "Password: ")
		pwBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		password = string(pwBytes)
		fmt.Fprintln(os.Stderr)
	} else {
		password = loginPassword
	}

	port := loginPort
	if port == 0 {
		port = 443
	}

	fmt.Fprintf(os.Stderr, "Testing connectivity to %s:%d...\n", loginHost, port)

	version, err := discovery.TestCUPIAuth(loginHost, port, loginUsername, password)
	if err != nil {
		return fmt.Errorf("connectivity/authentication check failed: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Authentication successful\n")
	if version != "" {
		fmt.Fprintf(os.Stderr, "CUC version: %s\n", version)
	}

	// Determine server key for config
	serverKey := loginServer
	if serverKey == "" {
		serverKey = loginHost
	}

	// Load existing config
	cfg, err := appconfig.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create/update server config
	serverCfg := appconfig.ServerConfig{
		Host:    loginHost,
		Port:    port,
		Version: version,
		Credentials: map[string]appconfig.CredentialConfig{
			auth.CredTypeCUPI: {Username: loginUsername},
		},
	}

	cfg.Servers[serverKey] = serverCfg

	// Set as default if requested or if it's the first server
	if loginDefault || len(cfg.Servers) == 1 {
		cfg.DefaultServer = serverKey
	}

	if dryRunFlag {
		fmt.Printf("[DRY RUN] Would save server '%s' to ~/.cupi-cli/config.json\n", serverKey)
		fmt.Printf("[DRY RUN]   host=%s  port=%d  version=%s\n", loginHost, port, version)
		fmt.Printf("[DRY RUN] Would store CUPI password in keyring for '%s'\n", loginHost)
		if loginDefault || len(cfg.Servers) == 0 {
			fmt.Printf("[DRY RUN] Would set '%s' as default server\n", serverKey)
		}
		return nil
	}

	// Save config
	if err := appconfig.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Server configuration saved\n")

	// Store password in keyring
	if err := auth.StorePassword(loginHost, auth.CredTypeCUPI, password); err != nil {
		return fmt.Errorf("failed to store password: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Password stored in keyring\n")

	fmt.Fprintf(os.Stderr, "Successfully authenticated to %s\n", serverKey)
	if loginDefault || len(cfg.Servers) == 1 {
		fmt.Fprintf(os.Stderr, "Set as default server\n")
	}

	return nil
}

func init() {
	authLoginCmd.Flags().StringVar(&loginHost, "host", "", "CUC host IP or hostname (required)")
	authLoginCmd.Flags().IntVar(&loginPort, "port", 443, "CUC HTTPS port")
	authLoginCmd.Flags().StringVar(&loginUsername, "username", "", "CUPI username (required)")
	authLoginCmd.Flags().StringVar(&loginPassword, "password", "", "CUPI password (will prompt if not provided)")
	authLoginCmd.Flags().StringVar(&loginServer, "server", "", "Server name for config (defaults to hostname)")
	authLoginCmd.Flags().BoolVar(&loginDefault, "default", false, "Set as default server")
}
