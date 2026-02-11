package cmd

import (
	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/helpers/open"
	"github.com/devbuddy/devbuddy/pkg/project"
)

var openCmd = &cobra.Command{
	Use:          "open [github|pullrequest]",
	Short:        "Open a link about your project",
	RunE:         openRun,
	Args:         zeroOrOneArg,
	GroupID:      "devbuddy",
	SilenceUsage: true,
}

func init() {
	openCmd.Flags().Bool("list", false, "List available project's URLs")
	_ = openCmd.MarkZshCompPositionalArgumentWords(1, "github", "gh", "pullrequest", "pr")
}

func openRun(cmd *cobra.Command, args []string) error {
	linkName := ""
	if len(args) == 1 {
		linkName = args[0]
	}

	proj, err := project.FindCurrent()
	if err != nil {
		return err
	}

	if GetFlagBool(cmd, "list") {
		return open.PrintLinks(proj)
	}

	url, err := open.FindLink(proj, linkName)
	if err != nil {
		return err
	}

	return open.Open(url)
}
