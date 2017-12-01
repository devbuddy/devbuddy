package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/pior/dad/pkg/integration"
)

var rootCmd = &cobra.Command{
	Use: "dad",
	Run: rootRun,
}

func init() {
	rootCmd.PersistentFlags().Bool("shell-init", false, "Shell initialization")
	// rootCmd.PersistentFlags().MarkHidden("shell-init")
	rootCmd.PersistentFlags().Bool("shell-completion", false, "Shell completion")
	// rootCmd.PersistentFlags().MarkHidden("shell-completion")

	rootCmd.AddCommand(cloneCmd)
}

func rootRun(cmd *cobra.Command, args []string) {
	if GetFlagBool(cmd, "shell-init") {
		integration.Print()
		os.Exit(0)
	}

	if GetFlagBool(cmd, "shell-completion") {
		fmt.Println("complete-me complete-me-again")
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
