package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/hook"
	"github.com/devbuddy/devbuddy/pkg/integration"
)

var rootCmd *cobra.Command

func build(version string) {
	rootCmd = &cobra.Command{
		Use:     "bud",
		Run:     rootRun,
		Version: version,
	}

	rootCmd.Flags().Bool("shell-init", false, "Shell initialization")
	rootCmd.Flags().Bool("with-completion", false, "Enable completion during initialization")

	rootCmd.Flags().Bool("shell-hook", false, "Shell prompt hook")
	err := rootCmd.Flags().MarkHidden("shell-hook")
	checkError(err)

	rootCmd.AddCommand(cdCmd)
	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(inspectCmd)
	rootCmd.AddCommand(openCmd)
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(upgradeCmd)
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
