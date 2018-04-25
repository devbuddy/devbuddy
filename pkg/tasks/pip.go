package tasks

import (
	"fmt"
)

func init() {
	allTasks["pip"] = newPip
}

type Pip struct {
	files []string
}

func newPip(config *taskConfig) (Task, error) {
	task := &Pip{}

	for _, value := range config.payload.([]interface{}) {
		if v, ok := value.(string); ok {
			task.files = append(task.files, v)
		} else {
			return nil, fmt.Errorf("invalid pip files")
		}
	}
	if len(task.files) == 0 {
		return nil, fmt.Errorf("no pip files specified")
	}

	return task, nil
}

func (p *Pip) name() string {
	return "Pip"
}

func (p *Pip) header() string {
	return "" //strings.Join(p.files, ", ")
}

func (p *Pip) preRunValidation(ctx *Context) (err error) {
	_, hasPythonFeature := ctx.features["python"]
	if !hasPythonFeature {
		return fmt.Errorf("You must specify a Python environment to use this task")
	}
	return nil
}

func (p *Pip) actions(ctx *Context) (actions []taskAction) {
	for _, file := range p.files {
		actions = append(actions, &pipInstall{file: file})
	}
	return
}

type pipInstall struct {
	file    string
	success bool
}

func (p *pipInstall) description() string {
	return fmt.Sprintf("install %s", p.file)
}

func (p *pipInstall) needed(ctx *Context) (bool, error) {
	return !p.success, nil
}

func (p *pipInstall) run(ctx *Context) error {
	code, err := runCommand(ctx, "pip", "install", "--require-virtualenv", "-r", p.file)
	if err != nil {
		return err
	}
	if code != 0 {
		return fmt.Errorf("Pip failed with code %d", code)
	}
	p.success = true
	return nil
}
