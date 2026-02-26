package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var cosCmd = &cobra.Command{
	Use:   "cos",
	Short: "Manage CUC classes of service",
	Long:  `Commands for listing, viewing, and updating CUC classes of service.`,
}

var cosListCmd = &cobra.Command{
	Use:   "list",
	Short: "List classes of service",
	RunE:  runCOSList,
}

var cosGetCmd = &cobra.Command{
	Use:   "get <name-or-objectId>",
	Short: "Get a class of service",
	Args:  cobra.ExactArgs(1),
	RunE:  runCOSGet,
}

var cosUpdateDisplayName string

var cosUpdateCmd = &cobra.Command{
	Use:   "update <name-or-objectId>",
	Short: "Update a class of service",
	Args:  cobra.ExactArgs(1),
	RunE:  runCOSUpdate,
}

func runCOSList(cmd *cobra.Command, args []string) error {
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

	coses, err := client.ListCOS(serverCfg.Host, serverCfg.Port, user, pass, "", maxFlag)
	if err != nil {
		return fmt.Errorf("failed to list classes of service: %w", err)
	}

	var rows []map[string]string
	for _, c := range coses {
		rows = append(rows, map[string]string{
			"displayName": c.DisplayName,
			"objectId":    c.ObjectId,
		})
	}

	return output.Print(rows, outputFlag)
}

func runCOSGet(cmd *cobra.Command, args []string) error {
	nameOrID := args[0]

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

	c, err := client.GetCOS(serverCfg.Host, serverCfg.Port, user, pass, nameOrID)
	if err != nil {
		return fmt.Errorf("failed to get class of service: %w", err)
	}

	data := map[string]interface{}{
		"objectId":    c.ObjectId,
		"displayName": c.DisplayName,
	}

	return output.Print(data, outputFlag)
}

func runCOSUpdate(cmd *cobra.Command, args []string) error {
	nameOrID := args[0]

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

	fields := map[string]interface{}{}
	if cosUpdateDisplayName != "" {
		fields["DisplayName"] = cosUpdateDisplayName
	}

	if len(fields) == 0 {
		return fmt.Errorf("at least one field must be specified to update")
	}

	if err := client.UpdateCOS(serverCfg.Host, serverCfg.Port, user, pass, nameOrID, fields); err != nil {
		return fmt.Errorf("failed to update class of service: %w", err)
	}

	fmt.Printf("Updated class of service: %s\n", nameOrID)
	return nil
}

func init() {
	cosCmd.AddCommand(cosListCmd)
	cosCmd.AddCommand(cosGetCmd)
	cosCmd.AddCommand(cosUpdateCmd)
	cosUpdateCmd.Flags().StringVar(&cosUpdateDisplayName, "display-name", "", "Display name")
}
