package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/helpers/debug"
	"github.com/devbuddy/devbuddy/pkg/helpers/open"
	"github.com/devbuddy/devbuddy/pkg/hook"
	"github.com/devbuddy/devbuddy/pkg/integration"
)

func build(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "bud",
		Run:     rootRun,
		Version: version,
	}

	rootCmd.Flags().Bool("shell-init", false, "Shell initialization")
	rootCmd.Flags().Bool("with-completion", false, "Enable completion during initialization")

	rootCmd.Flags().Bool("debug-info", false, "Print some debug information")
	rootCmd.Flags().Bool("report-issue", false, "Create an issue about DevBuddy on Github")

	rootCmd.Flags().Bool("shell-hook", false, "Shell prompt hook")
	err := rootCmd.Flags().MarkHidden("shell-hook")
	checkError(err)

	rootCmd.AddCommand(cdCmd)
	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(inspectCmd)
	rootCmd.AddCommand(openCmd)
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(upgradeCmd)

	return rootCmd
}

func rootRun(cmd *cobra.Command, args []string) {
	var err error

	if GetFlagBool(cmd, "shell-init") {
		if GetFlagBool(cmd, "with-completion") {
			err = cmd.GenBashCompletion(os.Stdout)
			checkError(err)
		}
		integration.Print()
		os.Exit(0)
	}

	if GetFlagBool(cmd, "shell-hook") {
		hook.Run()
		os.Exit(0)
	}

	if GetFlagBool(cmd, "debug-info") {
		fmt.Println(debug.FormatDebugInfo(cmd.Version, os.Environ(), debug.SafeFindCurrentProject()))
		os.Exit(0)
	}

	if GetFlagBool(cmd, "report-issue") {
		url := debug.NewGithubIssueURL(cmd.Version, os.Environ(), debug.SafeFindCurrentProject())
		err := open.Open(url)
		checkError(err)
		os.Exit(0)
	}

	err = cmd.Help()
	checkError(err)
}

func Execute(version string) {
	rootCmd := build(version)
	buildCustomCommands(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
