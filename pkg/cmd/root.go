package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/pior/dad/pkg/integration"
)

var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{
		Use: "dad",
		Run: rootRun,
	}

	rootCmd.PersistentFlags().Bool("shell-init", false, "Shell initialization")
	// rootCmd.PersistentFlags().MarkHidden("shell-init")
	rootCmd.PersistentFlags().Bool("with-completion", false, "Enable completion during initialization")
	// rootCmd.PersistentFlags().MarkHidden("shell-completion")

	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(cdCmd)
}

func rootRun(cmd *cobra.Command, args []string) {
	if GetFlagBool(cmd, "shell-init") {
		if GetFlagBool(cmd, "with-completion") {
			rootCmd.GenBashCompletion(os.Stdout)
		}
		integration.Print()
		os.Exit(0)
	}

	cmd.Help()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
