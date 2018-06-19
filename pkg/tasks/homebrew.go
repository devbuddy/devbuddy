package tasks

import (
	"fmt"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	allTasks["homebrew"] = newHomebrew
}

type Homebrew struct {
	formulas []string
}

func newHomebrew(config *taskConfig) (Task, error) {
	task := &Homebrew{}

	for _, value := range config.payload.([]interface{}) {
		if v, ok := value.(string); ok {
			task.formulas = append(task.formulas, v)
		} else {
			return nil, fmt.Errorf("invalid homebrew formulas")
		}
	}

	if len(task.formulas) == 0 {
		return nil, fmt.Errorf("no homebrew formulas specified")
	}

	return task, nil
}

func (h *Homebrew) name() string {
	return "Homebrew"
}

func (h *Homebrew) header() string {
	return strings.Join(h.formulas, ", ")
}

func (h *Homebrew) actions(ctx *context) (actions []taskAction) {
	for _, f := range h.formulas {
		actions = append(actions, &brewInstall{formula: f})
	}
	return
}

type brewInstall struct {
	formula string
	success bool
}

func (b *brewInstall) description() string {
	return fmt.Sprintf("installing %s", b.formula)
}

func (b *brewInstall) needed(ctx *context) (bool, error) {
	brew := helpers.NewHomebrew()

	installed := brew.IsInstalled(b.formula)

	return !installed, nil
}

func (b *brewInstall) run(ctx *context) error {
	err := command(ctx, "brew", "install", b.formula).Run()

	if err != nil {
		return fmt.Errorf("Homebrew failed: %s", err)
	}

	b.success = true
	return nil
}
