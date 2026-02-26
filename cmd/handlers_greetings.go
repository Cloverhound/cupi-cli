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
	greetingEnabled   string
	greetingPlayWhat  string
)

var handlersGreetingsCmd = &cobra.Command{
	Use:   "greetings",
	Short: "Manage call handler greetings",
}

var handlersGreetingsListCmd = &cobra.Command{
	Use:   "list <handler-name-or-id>",
	Short: "List greetings for a handler",
	Args:  cobra.ExactArgs(1),
	RunE:  runHandlersGreetingsList,
}

var handlersGreetingsGetCmd = &cobra.Command{
	Use:   "get <handler-name-or-id> <greeting-type>",
	Short: "Get a specific greeting",
	Args:  cobra.ExactArgs(2),
	RunE:  runHandlersGreetingsGet,
}

var handlersGreetingsUpdateCmd = &cobra.Command{
	Use:   "update <handler-name-or-id> <greeting-type>",
	Short: "Update a greeting",
	Args:  cobra.ExactArgs(2),
	RunE:  runHandlersGreetingsUpdate,
}

func runHandlersGreetingsList(cmd *cobra.Command, args []string) error {
	handlerNameOrID := args[0]

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

	h, err := client.GetCallHandler(serverCfg.Host, serverCfg.Port, user, pass, handlerNameOrID)
	if err != nil {
		return err
	}

	greetings, err := client.ListGreetings(serverCfg.Host, serverCfg.Port, user, pass, h.ObjectId)
	if err != nil {
		return err
	}

	var rows []map[string]string
	for _, g := range greetings {
		rows = append(rows, map[string]string{
			"greetingType":       g.GreetingType,
			"enabled":            g.Enabled,
			"playWhat":           g.PlayWhat,
			"timeExpiresSetFor":  g.TimeExpiresSetFor,
		})
	}

	return output.Print(rows, outputFlag)
}

func runHandlersGreetingsGet(cmd *cobra.Command, args []string) error {
	handlerNameOrID := args[0]
	greetingType := args[1]

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

	h, err := client.GetCallHandler(serverCfg.Host, serverCfg.Port, user, pass, handlerNameOrID)
	if err != nil {
		return err
	}

	g, err := client.GetGreeting(serverCfg.Host, serverCfg.Port, user, pass, h.ObjectId, greetingType)
	if err != nil {
		return err
	}

	return output.Print(g, outputFlag)
}

func runHandlersGreetingsUpdate(cmd *cobra.Command, args []string) error {
	handlerNameOrID := args[0]
	greetingType := args[1]

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

	h, err := client.GetCallHandler(serverCfg.Host, serverCfg.Port, user, pass, handlerNameOrID)
	if err != nil {
		return err
	}

	fields := map[string]interface{}{}
	if greetingEnabled != "" {
		fields["Enabled"] = greetingEnabled
	}
	if greetingPlayWhat != "" {
		fields["PlayWhat"] = greetingPlayWhat
	}

	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	if err := client.UpdateGreeting(serverCfg.Host, serverCfg.Port, user, pass, h.ObjectId, greetingType, fields); err != nil {
		return err
	}

	fmt.Printf("Updated greeting %s for handler %s\n", greetingType, h.DisplayName)
	return nil
}

func init() {
	handlersGreetingsUpdateCmd.Flags().StringVar(&greetingEnabled, "enabled", "", "Enabled (true|false)")
	handlersGreetingsUpdateCmd.Flags().StringVar(&greetingPlayWhat, "play-what", "", "Play what")

	handlersGreetingsCmd.AddCommand(handlersGreetingsListCmd)
	handlersGreetingsCmd.AddCommand(handlersGreetingsGetCmd)
	handlersGreetingsCmd.AddCommand(handlersGreetingsUpdateCmd)
}
