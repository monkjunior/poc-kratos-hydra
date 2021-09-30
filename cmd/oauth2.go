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
}
