package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// mockCmd represents the mock command
var mockCmd = &cobra.Command{
	Use:   "mock",
	Short: "Run mock programs",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mock called")
	},
}

func init() {
	rootCmd.AddCommand(mockCmd)
}
