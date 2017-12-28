package hook

import (
	"fmt"
	"os"

	color "github.com/logrusorgru/aurora"

	"github.com/pior/dad/pkg/project"
)

func Hook() {
	// In the shell hook, the stdout is evaluated by the shell
	// stderr is used to display messages to the user

	proj, err := project.FindCurrent()

	if err != nil && err != project.ErrProjectNotFound {
		// We can't annoy the user here, just quit silently
		return
	}

	if false {
		notify(fmt.Sprintf("Project: %s", proj.Path))
	}
}

func notify(msg string) {
	fmt.Fprintf(os.Stderr, "👴  %s\n", color.Cyan(msg))
}
