package cmd

import (
	"fmt"
	"strconv"

	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var astPerfmonCmd = &cobra.Command{
	Use:   "perfmon",
	Short: "Show perfmon objects and counters",
	Long: `Display available perfmon (performance monitor) objects and their associated counters.

Sub-commands allow drilling into a specific object to list its counters or collect live values.`,
	RunE: runASTPerfmon,
}

func runASTPerfmon(cmd *cobra.Command, args []string) error {
	srv, user, pass, err := resolveCredentials(cmd, "cupi")
	if err != nil {
		return err
	}

	objects, err := client.GetASTPerfmonObjects(srv.Host, user, pass)
	if err != nil {
		return fmt.Errorf("failed to get perfmon objects: %w", err)
	}

	var rows []map[string]string
	for _, obj := range objects {
		rows = append(rows, map[string]string{
			"host":         obj.Host,
			"object":       obj.ObjectName,
			"hasInstances": strconv.FormatBool(obj.HasInstances),
			"counterCount": strconv.Itoa(obj.CounterCount),
		})
	}

	return output.Print(rows, outputFlag)
}

// ast perfmon counters <object> — list counters for a specific perfmon object
var astPerfmonCountersCmd = &cobra.Command{
	Use:   "counters <object>",
	Short: "List counters for a specific perfmon object",
	Long: `List all counters available for a specific perfmon object, including descriptions and units.

Example:
  cupi ast perfmon counters "Cisco Unity Connection Voicemail"
  cupi ast perfmon counters "Memory"`,
	Args: cobra.ExactArgs(1),
	RunE: runASTPerfmonCounters,
}

func runASTPerfmonCounters(cmd *cobra.Command, args []string) error {
	srv, user, pass, err := resolveCredentials(cmd, "cupi")
	if err != nil {
		return err
	}

	counters, err := client.GetASTPerfmonCounters(srv.Host, user, pass, args[0])
	if err != nil {
		return fmt.Errorf("failed to get perfmon counters: %w", err)
	}

	var rows []map[string]string
	for _, c := range counters {
		rows = append(rows, map[string]string{
			"host":        c.Host,
			"object":      c.ObjectName,
			"counter":     c.CounterName,
			"unit":        c.Unit,
			"description": c.Description,
			"isExcluded":  strconv.FormatBool(c.IsExcluded),
		})
	}

	return output.Print(rows, outputFlag)
}

// ast perfmon collect <object> — collect real-time counter values
var astPerfmonCollectCmd = &cobra.Command{
	Use:   "collect <object>",
	Short: "Collect real-time counter values for a perfmon object",
	Long: `Collect and display current real-time values for all counters in a specific perfmon object.

Example:
  cupi ast perfmon collect "Cisco Unity Connection Voicemail"
  cupi ast perfmon collect "Memory"
  cupi ast perfmon collect "Processor"`,
	Args: cobra.ExactArgs(1),
	RunE: runASTPerfmonCollect,
}

func runASTPerfmonCollect(cmd *cobra.Command, args []string) error {
	srv, user, pass, err := resolveCredentials(cmd, "cupi")
	if err != nil {
		return err
	}

	data, err := client.GetASTPerfmonData(srv.Host, user, pass, args[0])
	if err != nil {
		return fmt.Errorf("failed to collect perfmon data: %w", err)
	}

	var rows []map[string]string
	for _, d := range data {
		row := map[string]string{
			"host":    d.Host,
			"object":  d.ObjectName,
			"counter": d.CounterName,
			"value":   d.Value,
			"cstatus": d.CStatus,
		}
		if d.Instance != "" {
			row["instance"] = d.Instance
		}
		rows = append(rows, row)
	}

	return output.Print(rows, outputFlag)
}

func init() {
	astPerfmonCmd.AddCommand(astPerfmonCountersCmd)
	astPerfmonCmd.AddCommand(astPerfmonCollectCmd)
}
