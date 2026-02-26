package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Manage CUC user templates",
	Long:  `Commands for listing and viewing CUC user templates.`,
}

var templatesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List user templates",
	RunE:  runTemplatesList,
}

var templatesGetCmd = &cobra.Command{
	Use:   "get <alias-or-objectId>",
	Short: "Get a user template",
	Args:  cobra.ExactArgs(1),
	RunE:  runTemplatesGet,
}

func runTemplatesList(cmd *cobra.Command, args []string) error {
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

	templates, err := client.ListUserTemplates(serverCfg.Host, serverCfg.Port, user, pass, "", maxFlag)
	if err != nil {
		return fmt.Errorf("failed to list user templates: %w", err)
	}

	var rows []map[string]string
	for _, t := range templates {
		rows = append(rows, map[string]string{
			"alias":       t.Alias,
			"displayName": t.DisplayName,
			"objectId":    t.ObjectId,
		})
	}

	return output.Print(rows, outputFlag)
}

func runTemplatesGet(cmd *cobra.Command, args []string) error {
	aliasOrID := args[0]

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

	t, err := client.GetUserTemplate(serverCfg.Host, serverCfg.Port, user, pass, aliasOrID)
	if err != nil {
		return fmt.Errorf("failed to get user template: %w", err)
	}

	data := map[string]interface{}{
		"objectId":    t.ObjectId,
		"alias":       t.Alias,
		"displayName": t.DisplayName,
	}

	return output.Print(data, outputFlag)
}

func init() {
	templatesCmd.AddCommand(templatesListCmd)
	templatesCmd.AddCommand(templatesGetCmd)
}
