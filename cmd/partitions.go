package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var partitionName string
var partitionDescription string

var partitionsCmd = &cobra.Command{
	Use:   "partitions",
	Short: "Manage partitions",
}

var partitionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List partitions",
	RunE:  runPartitionsList,
}

var partitionsGetCmd = &cobra.Command{
	Use:   "get <name-or-id>",
	Short: "Get a partition",
	Args:  cobra.ExactArgs(1),
	RunE:  runPartitionsGet,
}

var partitionsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a partition",
	RunE:  runPartitionsAdd,
}

var partitionsUpdateCmd = &cobra.Command{
	Use:   "update <name-or-id>",
	Short: "Update a partition",
	Args:  cobra.ExactArgs(1),
	RunE:  runPartitionsUpdate,
}

var partitionsRemoveCmd = &cobra.Command{
	Use:   "remove <name-or-id>",
	Short: "Remove a partition",
	Args:  cobra.ExactArgs(1),
	RunE:  runPartitionsRemove,
}

func runPartitionsList(cmd *cobra.Command, args []string) error {
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

	items, err := client.ListPartitions(serverCfg.Host, serverCfg.Port, user, pass, "", maxFlag)
	if err != nil {
		return err
	}

	var rows []map[string]string
	for _, item := range items {
		rows = append(rows, map[string]string{
			"objectId":    item.ObjectId,
			"name":        item.Name,
			"description": item.Description,
		})
	}

	return output.Print(rows, outputFlag)
}

func runPartitionsGet(cmd *cobra.Command, args []string) error {
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

	item, err := client.GetPartition(serverCfg.Host, serverCfg.Port, user, pass, nameOrID)
	if err != nil {
		return err
	}

	return output.Print(item, outputFlag)
}

func runPartitionsAdd(cmd *cobra.Command, args []string) error {
	if partitionName == "" {
		return fmt.Errorf("--name is required")
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
		"Name": partitionName,
	}
	if partitionDescription != "" {
		fields["Description"] = partitionDescription
	}

	if err := client.CreatePartition(serverCfg.Host, serverCfg.Port, user, pass, fields); err != nil {
		return err
	}

	fmt.Printf("Added partition %s\n", partitionName)
	return nil
}

func runPartitionsUpdate(cmd *cobra.Command, args []string) error {
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
	if partitionName != "" {
		fields["Name"] = partitionName
	}
	if partitionDescription != "" {
		fields["Description"] = partitionDescription
	}

	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	if err := client.UpdatePartition(serverCfg.Host, serverCfg.Port, user, pass, nameOrID, fields); err != nil {
		return err
	}

	fmt.Printf("Updated partition %s\n", nameOrID)
	return nil
}

func runPartitionsRemove(cmd *cobra.Command, args []string) error {
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

	if err := client.DeletePartition(serverCfg.Host, serverCfg.Port, user, pass, nameOrID); err != nil {
		return err
	}

	fmt.Printf("Removed partition %s\n", nameOrID)
	return nil
}

func init() {
	partitionsAddCmd.Flags().StringVar(&partitionName, "name", "", "Partition name (required)")
	partitionsAddCmd.Flags().StringVar(&partitionDescription, "description", "", "Description")

	partitionsUpdateCmd.Flags().StringVar(&partitionName, "name", "", "Partition name")
	partitionsUpdateCmd.Flags().StringVar(&partitionDescription, "description", "", "Description")

	partitionsCmd.AddCommand(partitionsListCmd)
	partitionsCmd.AddCommand(partitionsGetCmd)
	partitionsCmd.AddCommand(partitionsAddCmd)
	partitionsCmd.AddCommand(partitionsUpdateCmd)
	partitionsCmd.AddCommand(partitionsRemoveCmd)
}
