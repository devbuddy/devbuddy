package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/helpers/open"
	"github.com/devbuddy/devbuddy/pkg/project"
)

var openCmd = &cobra.Command{
	Use:   "open [github|pullrequest]",
	Short: "Open a link about your project",
	Run:   openRun,
	Args:  zeroOrOneArg,
}

func init() {
	openCmd.Flags().Bool("list", false, "List available project's URLs")
}

func openRun(cmd *cobra.Command, args []string) {
	linkName := ""
	if len(args) == 1 {
		linkName = args[0]
	}

	proj, err := project.FindCurrent()
	checkError(err)

	if GetFlagBool(cmd, "list") {
		err = open.PrintLinks(proj)
		checkError(err)
		os.Exit(0)
	}

	url, err := open.FindLink(proj, linkName)
	checkError(err)

	err = open.Open(url)
	checkError(err)
}
