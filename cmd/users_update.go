package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/spf13/cobra"
)

var (
	updateUserDtmf        string
	updateUserFirstName   string
	updateUserLastName    string
	updateUserDisplayName string
	updateUserDepartment  string
)

var usersUpdateCmd = &cobra.Command{
	Use:   "update <alias-or-objectId>",
	Short: "Update a CUC user",
	Long: `Update fields on a CUC voicemail user.

Examples:
  cupi users update jsmith --display-name "John Smith"
  cupi users update jsmith --dtmf 1002 --department "Engineering"`,
	Args: cobra.ExactArgs(1),
	RunE: runUsersUpdate,
}

func runUsersUpdate(cmd *cobra.Command, args []string) error {
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

	fields := map[string]interface{}{}
	if updateUserDtmf != "" {
		fields["DtmfAccessId"] = updateUserDtmf
	}
	if updateUserFirstName != "" {
		fields["FirstName"] = updateUserFirstName
	}
	if updateUserLastName != "" {
		fields["LastName"] = updateUserLastName
	}
	if updateUserDisplayName != "" {
		fields["DisplayName"] = updateUserDisplayName
	}
	if updateUserDepartment != "" {
		fields["Department"] = updateUserDepartment
	}

	if len(fields) == 0 {
		return fmt.Errorf("at least one field must be specified to update")
	}

	if err := client.UpdateUser(serverCfg.Host, serverCfg.Port, user, pass, aliasOrID, fields); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	fmt.Printf("Updated user: %s\n", aliasOrID)
	return nil
}

func init() {
	usersUpdateCmd.Flags().StringVar(&updateUserDtmf, "dtmf", "", "DTMF access ID / extension")
	usersUpdateCmd.Flags().StringVar(&updateUserFirstName, "first-name", "", "First name")
	usersUpdateCmd.Flags().StringVar(&updateUserLastName, "last-name", "", "Last name")
	usersUpdateCmd.Flags().StringVar(&updateUserDisplayName, "display-name", "", "Display name")
	usersUpdateCmd.Flags().StringVar(&updateUserDepartment, "department", "", "Department")
}
