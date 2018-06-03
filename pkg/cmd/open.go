package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pior/dad/pkg/helpers"
	"github.com/pior/dad/pkg/project"
)

var openCmd = &cobra.Command{
	Use:   "open [github|pullrequest]",
	Short: "Open a link about your project",
	Run:   openRun,
}

func openRun(cmd *cobra.Command, args []string) {
	proj, err := project.FindCurrent()
	checkError(err)

	url, err := findOpenURL(proj, args)
	checkError(err)

	err = helpers.Open(url)
	checkError(err)
}

func findOpenURL(proj *project.Project, args []string) (url string, err error) {
	if len(args) == 0 {
		if len(proj.Manifest.Open) == 1 {
			for _, url = range proj.Manifest.Open {
				return url, nil
			}
		}
		return "", fmt.Errorf("expecting one argument")
	}

	if len(args) > 1 {
		return "", fmt.Errorf("expecting one argument")
	}

	switch args[0] {
	case "github", "gh":
		url, err = helpers.NewGitRepo(proj.Path).BuildGithubProjectURL()
	case "pullrequest", "pr":
		url, err = helpers.NewGitRepo(proj.Path).BuildGithubPullrequestURL()
	default:
		url = proj.Manifest.Open[args[0]]
		if url == "" {
			err = fmt.Errorf("no link for '%s'", args[0])
		}
	}
	return url, err
}
