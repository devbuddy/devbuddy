package tasks

import (
	"fmt"
	"strings"

	"github.com/pior/dad/pkg/helpers"
)

func init() {
	allTasks["homebrew"] = newHomebrew
}

type Homebrew struct {
	files []string
}

func newHomebrew(config *taskConfig) (Task, error) {
	task := &Homebrew{}

	for _, value := range config.payload.([]interface{}) {
		if v, ok := value.(string); ok {
			task.files = append(task.files, v)
		} else {
			return nil, fmt.Errorf("invalid homebrew packages")
		}
	}

	if len(task.files) == 0 {
		return nil, fmt.Errorf("no homebrew packages specified")
	}

	return task, nil
}

func (h *Homebrew) name() string {
	return "Homebrew"
}

func (h *Homebrew) header() string {
	return strings.Join(h.files, ", ")
}

func (h *Homebrew) perform(ctx *Context) error {
	packageHelper := helpers.NewHomebrew()

	if ctx.env.Os() != "darwin" {
		return fmt.Errorf("homebrew is only supported on darwin operating system")
	}

	for _, file := range h.files {
		if packageHelper.PackageIsInstalled(file) {
			continue
		}

		code, err := runCommand(ctx, "brew", "install", file)
		if err != nil {
			return err
		}

		if code != 0 {
			return fmt.Errorf("Homebrew failed with code %d", code)
		}
	}
	return nil
}
