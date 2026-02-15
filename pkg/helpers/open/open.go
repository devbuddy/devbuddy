package open

import (
	"fmt"
	"os/exec"
	"runtime"

	color "github.com/logrusorgru/aurora"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/manifest"
	"github.com/devbuddy/devbuddy/pkg/project"
)

// Open a file or URL with the default application, return immediately.
// Use `xdg-open` or `open` depending on the platform.
func Open(location string) error {
	openCommand := "xdg-open"
	if runtime.GOOS == "darwin" {
		openCommand = "open"
	}

	return exec.Command(openCommand, location).Start()
}

// FindLink returns the url of a link about the project.
// Possible links are github/pullrequest pages and arbitrary links declared in dev.yml. In case of collision, links
// declared in dev.yml have precedence over Github links.
func FindLink(ctx *context.Context, linkName string) (url string, err error) {
	man, err := manifest.Load(ctx.Project.Path)
	if err != nil {
		return "", err
	}

	if linkName == "" {
		if len(man.Open) == 1 {
			for _, url = range man.Open {
				return url, nil
			}
		}
		return "", fmt.Errorf("which link should I open?")
	}

	link := project.FindBestLinkMatch(linkName, BuildIndex(man.Open))
	url = man.Open[link]

	if url != "" {
		return
	}

	switch linkName {
	case "github", "gh":
		url, err = helpers.NewGitRepo(ctx, ctx.Project.Path).BuildGithubProjectURL()
		return
	case "pullrequest", "pr":
		url, err = helpers.NewGitRepo(ctx, ctx.Project.Path).BuildGithubPullrequestURL()
		return
	default:
		err = fmt.Errorf("no link for '%s'", linkName)
	}

	return
}

func PrintLinks(proj *project.Project) (err error) {
	man, err := manifest.Load(proj.Path)
	if err != nil {
		return err
	}

	if len(man.Open) == 0 {
		return fmt.Errorf("no links found in the project")
	}
	for title, url := range man.Open {
		fmt.Println(color.Green(title), "\t", url)
	}

	return nil
}

func BuildIndex(links map[string]string) []string {
	var index []string

	for key := range links {
		index = append(index, key)
	}

	return index
}
