package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is set at build time via:
//
//	go build -ldflags "-X github.com/Cloverhound/cupi-cli/cmd.Version=v1.2.3" .
//
// When built without ldflags (e.g. go build .) it reports "dev".
var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the cupi version",
	Long:  "Print the current version of the cupi CLI binary.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cupi %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
