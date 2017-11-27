package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/pior/dad/pkg/integration"
)

var rootCmd = &cobra.Command{
	Use: "dad [args] [COMMAND]",
	Run: run,
}

func init() {
	rootCmd.PersistentFlags().Bool("init", false, "Shell initialization")
}

func run(cmd *cobra.Command, args []string) {
	if GetFlagBool(cmd, "init") {
		integration.Print()
		os.Exit(0)
	}

	cmd.Help()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
