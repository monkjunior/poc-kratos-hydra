package cmd

import (
	"fmt"

	"github.com/monkjunior/poc-kratos-hydra/pkg/config"
	"github.com/spf13/cobra"
)

// testConfigCmd represents the testConfig command
var testConfigCmd = &cobra.Command{
	Use:   "test-config",
	Short: "Show the current config",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Config: %+v\n", config.Cfg)
	},
}

func init() {
	rootCmd.AddCommand(testConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testConfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
