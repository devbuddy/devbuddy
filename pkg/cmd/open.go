package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pior/dad/pkg/helpers"
	"github.com/pior/dad/pkg/project"
)

var openCmd = &cobra.Command{
	Use:   "open [github]",
	Short: "Open a link about your project",
	Run:   openRun,
	Args:  onlyOneArg,
}

func openRun(cmd *cobra.Command, args []string) {
	proj, err := project.FindCurrent()
	checkError(err)

	switch args[0] {
	case "github", "gh":
		err = openGithub(proj)
	default:
		err = fmt.Errorf("I don't know how to open %s", args[0])
	}
	checkError(err)
}

func openGithub(proj *project.Project) error {
	gitRepo := helpers.NewGitRepo(proj.Path)
	branch, err := gitRepo.GetCurrentBranch()
	if err != nil {
		return err
	}
	remoteURL, err := gitRepo.GetRemoteURL()
	if err != nil {
		return err
	}
	webURL, err := helpers.WebURLFromGitURL(remoteURL, branch)
	if err != nil {
		return err
	}
	return helpers.Open(webURL)
}
