package cmd

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
	"github.com/Cloverhound/cupi-cli/internal/auth"
	"github.com/Cloverhound/cupi-cli/internal/client"
	"github.com/spf13/cobra"
)

var (
	addUserAlias       string
	addUserDtmf        string
	addUserFirstName   string
	addUserLastName    string
	addUserDisplayName string
	addUserTemplate    string
)

var usersAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a CUC user",
	Long: `Create a new CUC voicemail user.

The --template flag specifies the user template alias to base the new user on.
Defaults to 'voicemailusertemplate' if not specified.

Examples:
  cupi users add --alias jsmith --dtmf 1001 --first-name John --last-name Smith
  cupi users add --alias jsmith --dtmf 1001 --template voicemailusertemplate
  cupi --dry-run users add --alias testuser --dtmf 9999`,
	RunE: runUsersAdd,
}

func runUsersAdd(cmd *cobra.Command, args []string) error {
	if addUserAlias == "" {
		return fmt.Errorf("--alias is required")
	}
	if addUserDtmf == "" {
		return fmt.Errorf("--dtmf is required")
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

	templateAlias := addUserTemplate
	if templateAlias == "" {
		templateAlias = "voicemailusertemplate"
	}

	fields := map[string]interface{}{
		"Alias":        addUserAlias,
		"DtmfAccessId": addUserDtmf,
	}
	if addUserFirstName != "" {
		fields["FirstName"] = addUserFirstName
	}
	if addUserLastName != "" {
		fields["LastName"] = addUserLastName
	}
	if addUserDisplayName != "" {
		fields["DisplayName"] = addUserDisplayName
	}

	newUser, err := client.CreateUser(serverCfg.Host, serverCfg.Port, user, pass, templateAlias, fields)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	fmt.Printf("Created user: alias=%s objectId=%s\n", newUser.Alias, newUser.ObjectId)
	return nil
}

func init() {
	usersAddCmd.Flags().StringVar(&addUserAlias, "alias", "", "User alias/login name (required)")
	usersAddCmd.Flags().StringVar(&addUserDtmf, "dtmf", "", "DTMF access ID / extension (required)")
	usersAddCmd.Flags().StringVar(&addUserFirstName, "first-name", "", "First name")
	usersAddCmd.Flags().StringVar(&addUserLastName, "last-name", "", "Last name")
	usersAddCmd.Flags().StringVar(&addUserDisplayName, "display-name", "", "Display name")
	usersAddCmd.Flags().StringVar(&addUserTemplate, "template", "", "User template alias (default: voicemailusertemplate)")
}
