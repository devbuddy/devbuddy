package tasks

import (
	"fmt"
	"strings"
)

func init() {
	t := registerTask("pip")
	t.name = "Pip"
	t.requiredFeature = "python"
	t.builder = newPip
}

func newPip(config *taskConfig) (Task, error) {
	var files []string

	for _, value := range config.payload.([]interface{}) {
		if v, ok := value.(string); ok {
			files = append(files, v)
		} else {
			return nil, fmt.Errorf("invalid pip files")
		}
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no pip files specified")
	}

	task := &Task{
		header: strings.Join(files, ", "),
	}
	for _, file := range files {
		task.addAction(&pipInstall{file: file})
	}

	return task, nil
}

type pipInstall struct {
	file    string
	success bool
}

func (p *pipInstall) description() string {
	return fmt.Sprintf("install %s", p.file)
}

func (p *pipInstall) needed(ctx *context) (bool, error) {
	return !p.success, nil
}

func (p *pipInstall) run(ctx *context) error {
	err := command(ctx, "pip", "install", "--require-virtualenv", "-r", p.file).AddOutputFilter("already satisfied").Run()

	if err != nil {
		return fmt.Errorf("Pip failed: %s", err)
	}
	p.success = true
	return nil
}
