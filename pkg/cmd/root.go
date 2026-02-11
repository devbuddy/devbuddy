package cmd

import (
	"errors"
	"fmt"
	"os"

	color "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/helpers/debug"
	"github.com/devbuddy/devbuddy/pkg/helpers/open"
	"github.com/devbuddy/devbuddy/pkg/hook"
	"github.com/devbuddy/devbuddy/pkg/integration"
)

// errTasksFailed is returned by `bud up` when one or more tasks fail.
var errTasksFailed = errors.New("some tasks failed")

func build(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "bud",
		RunE:              rootRun,
		Version:           version,
		SilenceUsage:      true,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
	}

	rootCmd.AddGroup(
		&cobra.Group{ID: "devbuddy", Title: "DevBuddy Commands:"},
		&cobra.Group{ID: "project", Title: "Project Commands:"},
	)
	rootCmd.SetHelpCommandGroupID("devbuddy")
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.Flags().Bool("shell-init", false, "Shell initialization")
	rootCmd.Flags().Bool("with-completion", false, "Enable completion during initialization")

	rootCmd.Flags().Bool("debug-info", false, "Print some debug information")
	rootCmd.Flags().Bool("report-issue", false, "Create an issue about DevBuddy on Github")

	rootCmd.Flags().Bool("shell-hook", false, "Shell prompt hook")
	if err := rootCmd.Flags().MarkHidden("shell-hook"); err != nil {
		panic(fmt.Sprintf("bug: failed to mark flag as hidden: %s", err))
	}

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

func rootRun(cmd *cobra.Command, _ []string) error {
	if GetFlagBool(cmd, "shell-init") {
		withCompletion := GetFlagBool(cmd, "with-completion")
		integration.Print(withCompletion, cmd)
		return nil
	}

	if GetFlagBool(cmd, "shell-hook") {
		hook.Run()
		return nil
	}

	if GetFlagBool(cmd, "debug-info") {
		fmt.Println(debug.FormatDebugInfo(cmd.Version, os.Environ(), debug.SafeFindCurrentProject()))
		return nil
	}

	if GetFlagBool(cmd, "report-issue") {
		url := debug.NewGithubIssueURL(cmd.Version, os.Environ(), debug.SafeFindCurrentProject())
		return open.Open(url)
	}

	return cmd.Help()
}

func Execute(version string) {
	rootCmd := build(version)
	buildCustomCommands(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		if !errors.Is(err, errTasksFailed) {
			fmt.Fprintln(os.Stderr, color.Red("Error:"), err)
		}
		os.Exit(1)
	}
}
