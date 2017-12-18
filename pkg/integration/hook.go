package integration

import (
	"fmt"
	"os"

	color "github.com/logrusorgru/aurora"

	"github.com/pior/dad/pkg/project"
)

func Hook() {
	proj, err := project.FindCurrent()
	if err != nil {
		// We didnt find a project, just quit silently
		return
	}

	if false {
		notify(fmt.Sprintf("Active project: %s", proj.Path))
	}
}

func notify(msg string) {
	fmt.Fprintf(os.Stderr, "👴  %s\n", color.Cyan(msg))
}
