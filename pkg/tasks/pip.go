package tasks

import (
	"fmt"
	"strings"
)

func init() {
	t := registerTask("pip")
	t.name = "Pip"
	t.requiredFeature = "python"
	t.parser = parserPip
}

func parserPip(config *taskConfig, task *Task) error {
	var files []string

	for _, value := range config.payload.([]interface{}) {
		if v, ok := value.(string); ok {
			files = append(files, v)
		} else {
			return fmt.Errorf("invalid pip files")
		}
	}
	if len(files) == 0 {
		return fmt.Errorf("no pip files specified")
	}

	task.header = strings.Join(files, ", ")

	for _, file := range files {
		task.addAction(&pipInstall{file: file})
	}

	return nil
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
