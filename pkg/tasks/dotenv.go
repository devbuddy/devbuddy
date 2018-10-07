package tasks

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/utils"
)

const dotenvFilename = ".env"

func init() {
	t := registerTaskDefinition("dotenv")
	t.name = "Dotenv"
	t.parser = parserDotenv
}

func parserDotenv(config *taskConfig, task *Task) error {
	task.featureName = "dotenv"
	task.featureParam = dotenvFilename

	task.addAction(&dotenvCheck{path: dotenvFilename})

	return nil
}

type dotenvCheck struct {
	path string
}

func (p *dotenvCheck) description() string {
	return fmt.Sprintf("Loading %s", p.path)
}

func (p *dotenvCheck) needed(ctx *context) (bool, error) {
	if !utils.PathExists(p.path) {
		ctx.ui.TaskWarning(
			fmt.Sprintf("The dotenv task expects a \"%s\" file in the root of the project", p.path),
		)
	}
	return false, nil
}

func (p *dotenvCheck) run(ctx *context) error {
	return nil
}
