package debug

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/manifest"
	"github.com/devbuddy/devbuddy/pkg/project"
)

func FormatDebugInfo(version string, environ []string, projectPath string) string {
	tmpl := `
**DevBuddy version**

` + "`%s`" + `

**Environment**

%s
**dev.yml**

%s
`

	man := "Project not found."
	if projectPath != "" {
		man = fmtManifest(projectPath)
	}

	return fmt.Sprintf(tmpl, version, fmtEnvVar(environ), man)
}

func NewGithubIssueURL(version string, environ []string, projectPath string) string {
	body := `
## Describe the issue

What happened? What did you expect?

## Steps to reproduce

1. ...
2. ...

## Debug Information
`
	body += FormatDebugInfo(version, environ, projectPath)
	bodyParam := url.QueryEscape(body)
	return fmt.Sprintf("https://github.com/devbuddy/devbuddy/issues/new?labels=user-bug&body=%s", bodyParam)
}

// SafeFindCurrentProject returns the path of the project if found, and an empty string otherwise
func SafeFindCurrentProject() string {
	proj, err := project.FindCurrent()
	if err != nil {
		return ""
	}
	return proj.Path
}

func fmtEnvVar(environ []string) string {
	names := [][]string{
		{
			"SHELL",
			"SHLVL",
			"TERM",
		},
		{
			"PATH",
			"USER",
			"PWD",
		},
		{
			"BUD_DEBUG",
			"BUD_FINALIZER_FILE",
			"BUD_AUTO_ENV_FEATURES",
		},
		{
			"GOROOT",
			"GOPATH",
			"GO111MODULE",
			"VIRTUAL_ENV",
		},
	}

	vars := environAsMap(environ)

	out := "```shell\n"
	for idx, section := range names {
		if idx != 0 {
			out += "\n"
		}
		for _, name := range section {
			out += fmt.Sprintf("%s=\"%s\"\n", name, vars[name])
		}
	}
	out += "```\n"
	return out
}

func environAsMap(environ []string) map[string]string {
	vars := map[string]string{}
	for _, entry := range environ {
		splitted := strings.SplitN(entry, "=", 2)
		vars[splitted[0]] = splitted[1]
	}
	return vars
}

func fmtManifest(path string) string {
	man, err := manifest.Load(path)
	if err != nil {
		return fmt.Sprintf("Failed to read manifest: %s", err)
	}

	out := fmt.Sprintf("Project path: `%s`\n\n", path)

	for idx, task := range man.Up {
		out += fmt.Sprintf("%d. `%+v`\n", idx, task)
	}
	return out
}
