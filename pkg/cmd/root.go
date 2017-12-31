package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/pior/dad/pkg/hook"
	"github.com/pior/dad/pkg/integration"
)

var rootCmd *cobra.Command

func build(version string) {
	rootCmd = &cobra.Command{
		Use:     "dad",
		Run:     rootRun,
		Version: version,
	}

	rootCmd.PersistentFlags().Bool("shell-init", false, "Shell initialization")
	rootCmd.PersistentFlags().Bool("with-completion", false, "Enable completion during initialization")

	rootCmd.PersistentFlags().Bool("shell-hook", false, "Shell prompt hook")
	err := rootCmd.PersistentFlags().MarkHidden("shell-hook")
	checkError(err)

	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(cdCmd)
	rootCmd.AddCommand(upCmd)
}

func rootRun(cmd *cobra.Command, args []string) {
	var err error

	if GetFlagBool(cmd, "shell-init") {
		if GetFlagBool(cmd, "with-completion") {
			err = rootCmd.GenBashCompletion(os.Stdout)
			checkError(err)
		}
		integration.Print()
		os.Exit(0)
	}

	if GetFlagBool(cmd, "shell-hook") {
		hook.Hook()
		os.Exit(0)
	}

	err = cmd.Help()
	checkError(err)
}

func Execute(version string) {
	build(version)
	buildCustomCommands()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
