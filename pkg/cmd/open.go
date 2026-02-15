package cmd

import (
	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers/open"
)

var openCmd = &cobra.Command{
	Use:          "open [pattern]",
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
	ctx, err := context.LoadWithProject()
	if err != nil {
		return err
	}

	linkName := ""
	if len(args) == 1 {
		linkName = args[0]
	}

	if GetFlagBool(cmd, "list") {
		return open.PrintLinks(ctx.Project)
	}

	url, err := open.FindLink(ctx, linkName)
	if err != nil {
		return err
	}

	return open.Open(url)
}
