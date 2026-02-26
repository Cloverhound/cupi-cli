package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/Cloverhound/cupi-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	rrDisplayName   string
	rrRouteAction   string
	rrEnabled       string
	rrQuery         string
	rrCondOperator  string
	rrCondParameter string
	rrCondOperand   string
)

var routingRulesCmd = &cobra.Command{
	Use:   "routingrules",
	Short: "Manage routing rules",
}

var routingRulesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List routing rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		items, err := client.ListRoutingRules(srv.Host, srv.Port, cupiUser, cupiPass, rrQuery, maxFlag)
		if err != nil {
			return err
		}
		var rows []map[string]string
		for _, v := range items {
			rows = append(rows, map[string]string{
				"objectId":    v.ObjectId,
				"displayName": v.DisplayName,
				"type":        v.Type,
				"routeAction": v.RouteAction,
				"enabled":     v.Enabled,
			})
		}
		return output.Print(rows, outputFlag)
	},
}

var routingRulesGetCmd = &cobra.Command{
	Use:   "get <name-or-id>",
	Short: "Get a routing rule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		item, err := client.GetRoutingRule(srv.Host, srv.Port, cupiUser, cupiPass, args[0])
		if err != nil {
			return err
		}
		return output.Print(item, outputFlag)
	},
}

var routingRulesAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a routing rule",
	RunE: func(cmd *cobra.Command, args []string) error {
		if rrDisplayName == "" {
			return fmt.Errorf("--display-name is required")
		}
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{"DisplayName": rrDisplayName}
		if rrRouteAction != "" {
			fields["RouteAction"] = rrRouteAction
		}
		if rrEnabled != "" {
			fields["Enabled"] = rrEnabled
		}
		if err := client.CreateRoutingRule(srv.Host, srv.Port, cupiUser, cupiPass, fields); err != nil {
			return err
		}
		fmt.Printf("Added routing rule %s\n", rrDisplayName)
		return nil
	},
}

var routingRulesUpdateCmd = &cobra.Command{
	Use:   "update <name-or-id>",
	Short: "Update a routing rule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{}
		if rrDisplayName != "" {
			fields["DisplayName"] = rrDisplayName
		}
		if rrRouteAction != "" {
			fields["RouteAction"] = rrRouteAction
		}
		if rrEnabled != "" {
			fields["Enabled"] = rrEnabled
		}
		if len(fields) == 0 {
			return fmt.Errorf("no fields to update; use --help to see available flags")
		}
		if err := client.UpdateRoutingRule(srv.Host, srv.Port, cupiUser, cupiPass, args[0], fields); err != nil {
			return err
		}
		fmt.Printf("Updated routing rule %s\n", args[0])
		return nil
	},
}

var routingRulesRemoveCmd = &cobra.Command{
	Use:   "remove <name-or-id>",
	Short: "Remove a routing rule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		if err := client.DeleteRoutingRule(srv.Host, srv.Port, cupiUser, cupiPass, args[0]); err != nil {
			return err
		}
		fmt.Printf("Removed routing rule %s\n", args[0])
		return nil
	},
}

var rrConditionsCmd = &cobra.Command{
	Use:   "conditions",
	Short: "Manage routing rule conditions",
}

var rrConditionsListCmd = &cobra.Command{
	Use:   "list <rule-id>",
	Short: "List conditions for a routing rule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		items, err := client.ListRoutingRuleConditions(srv.Host, srv.Port, cupiUser, cupiPass, args[0])
		if err != nil {
			return err
		}
		var rows []map[string]string
		for _, v := range items {
			rows = append(rows, map[string]string{
				"objectId":     v.ObjectId,
				"operatorType": v.OperatorType,
				"parameter":    v.Parameter,
				"operandTwo":   v.OperandTwo,
			})
		}
		return output.Print(rows, outputFlag)
	},
}

var rrConditionsAddCmd = &cobra.Command{
	Use:   "add <rule-id>",
	Short: "Add a condition to a routing rule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if rrCondParameter == "" {
			return fmt.Errorf("--parameter is required")
		}
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		fields := map[string]interface{}{"Parameter": rrCondParameter}
		if rrCondOperator != "" {
			fields["OperatorType"] = rrCondOperator
		}
		if rrCondOperand != "" {
			fields["OperandTwo"] = rrCondOperand
		}
		if err := client.CreateRoutingRuleCondition(srv.Host, srv.Port, cupiUser, cupiPass, args[0], fields); err != nil {
			return err
		}
		fmt.Printf("Added condition to routing rule %s\n", args[0])
		return nil
	},
}

var rrConditionsRemoveCmd = &cobra.Command{
	Use:   "remove <rule-id> <condition-id>",
	Short: "Remove a condition from a routing rule",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, cupiUser, cupiPass, err := resolveCredentials(cmd, auth.CredTypeCUPI)
		if err != nil {
			return err
		}
		if err := client.DeleteRoutingRuleCondition(srv.Host, srv.Port, cupiUser, cupiPass, args[0], args[1]); err != nil {
			return err
		}
		fmt.Printf("Removed condition %s\n", args[1])
		return nil
	},
}

func init() {
	routingRulesListCmd.Flags().StringVar(&rrQuery, "query", "", "Filter query")

	routingRulesAddCmd.Flags().StringVar(&rrDisplayName, "display-name", "", "Display name (required)")
	routingRulesAddCmd.Flags().StringVar(&rrRouteAction, "route-action", "", "Route action")
	routingRulesAddCmd.Flags().StringVar(&rrEnabled, "enabled", "", "Enabled (true/false)")

	routingRulesUpdateCmd.Flags().StringVar(&rrDisplayName, "display-name", "", "Display name")
	routingRulesUpdateCmd.Flags().StringVar(&rrRouteAction, "route-action", "", "Route action")
	routingRulesUpdateCmd.Flags().StringVar(&rrEnabled, "enabled", "", "Enabled (true/false)")

	rrConditionsAddCmd.Flags().StringVar(&rrCondParameter, "parameter", "", "Condition parameter (required)")
	rrConditionsAddCmd.Flags().StringVar(&rrCondOperator, "operator", "", "Operator type")
	rrConditionsAddCmd.Flags().StringVar(&rrCondOperand, "operand", "", "Second operand value")

	rrConditionsCmd.AddCommand(rrConditionsListCmd)
	rrConditionsCmd.AddCommand(rrConditionsAddCmd)
	rrConditionsCmd.AddCommand(rrConditionsRemoveCmd)

	routingRulesCmd.AddCommand(routingRulesListCmd)
	routingRulesCmd.AddCommand(routingRulesGetCmd)
	routingRulesCmd.AddCommand(routingRulesAddCmd)
	routingRulesCmd.AddCommand(routingRulesUpdateCmd)
	routingRulesCmd.AddCommand(routingRulesRemoveCmd)
	routingRulesCmd.AddCommand(rrConditionsCmd)
}
