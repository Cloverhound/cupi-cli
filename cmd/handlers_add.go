package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/spf13/cobra"
)

var (
	addHandlerDisplayName  string
	addHandlerDtmf         string
	addHandlerTemplateObjID string
)

var handlersAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a CUC call handler",
	RunE:  runHandlersAdd,
}

func runHandlersAdd(cmd *cobra.Command, args []string) error {
	if addHandlerDisplayName == "" {
		return fmt.Errorf("--display-name is required")
	}
	if addHandlerTemplateObjID == "" {
		return fmt.Errorf("--template-id is required (use handlers list to find a template ObjectId)")
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

	fields := map[string]interface{}{
		"DisplayName": addHandlerDisplayName,
	}
	if addHandlerDtmf != "" {
		fields["DtmfAccessId"] = addHandlerDtmf
	}

	h, err := client.CreateCallHandler(serverCfg.Host, serverCfg.Port, user, pass, addHandlerTemplateObjID, fields)
	if err != nil {
		return fmt.Errorf("failed to create call handler: %w", err)
	}

	fmt.Printf("Created call handler: displayName=%s objectId=%s\n", h.DisplayName, h.ObjectId)
	return nil
}

func init() {
	handlersAddCmd.Flags().StringVar(&addHandlerDisplayName, "display-name", "", "Display name (required)")
	handlersAddCmd.Flags().StringVar(&addHandlerDtmf, "dtmf", "", "DTMF access ID")
	handlersAddCmd.Flags().StringVar(&addHandlerTemplateObjID, "template-id", "", "Template ObjectId (required)")
}
