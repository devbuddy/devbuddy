package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/project"
)

var openCmd = &cobra.Command{
	Use:   "open [github|pullrequest]",
	Short: "Open a link about your project",
	Run:   openRun,
	Args:  onlyOneArg,
}

func openRun(cmd *cobra.Command, args []string) {
	proj, err := project.FindCurrent()
	checkError(err)

	var url string

	switch args[0] {
	case "github", "gh":
		url, err = helpers.NewGitRepo(proj.Path).BuildGithubProjectURL()
	case "pullrequest", "pr":
		url, err = helpers.NewGitRepo(proj.Path).BuildGithubPullrequestURL()
	default:
		url = proj.Manifest.Open[args[0]]
		if url != "" {
			break
		}
		err = fmt.Errorf("no link for '%s'", args[0])
	}
	checkError(err)
	if url != "" {
		err = helpers.Open(url)
	}
	checkError(err)
}
