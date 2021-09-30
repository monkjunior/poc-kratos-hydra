package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// oauth2Cmd represents the oauth2 command
var oauth2Cmd = &cobra.Command{
	Use:   "oauth2",
	Short: "Collection of tools to interact with Hydra",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("oauth2 called")
	},
}

func init() {
	rootCmd.AddCommand(oauth2Cmd)

	oauth2Cmd.PersistentFlags().Bool("fake-tls-termination", false, "Fake tls termination by adding to http headers \"X-Forwarded-Proto: https\"")
	oauth2Cmd.PersistentFlags().Bool("skip-tls-verify", false, "Foolishly accept TLS certificates signed by unknown certificate authorities")
}
