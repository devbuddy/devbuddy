package cmd

import (
	"github.com/spf13/cobra"

	"github.com/pior/dad/pkg/helpers/open"
	"github.com/pior/dad/pkg/project"
)

var openCmd = &cobra.Command{
	Use:   "open [github|pullrequest]",
	Short: "Open a link about your project",
	Run:   openRun,
	Args:  zeroOrOneArg,
}

func openRun(cmd *cobra.Command, args []string) {
	linkName := ""
	if len(args) == 1 {
		linkName = args[0]
	}

	proj, err := project.FindCurrent()
	checkError(err)

	url, err := open.FindLink(proj, linkName)
	checkError(err)

	err = open.Open(url)
	checkError(err)
}
