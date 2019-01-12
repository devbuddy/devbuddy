package tasks

import (
	"fmt"
)

type taskParser func(*TaskConfig, *Task) error

type taskDefinition struct {
	key           string     // the value used in dev.yml to identify a task
	name          string     // the displayed name of the task
	requiredTask  string     // another task that should be declared before this task
	parser        taskParser // the config parser for this task
	osRequirement string     // The platform this task can run on. "debian", "macos"
}

var taskDefinitions = make(map[string]*taskDefinition)

func Register(key, name string, parserFunc taskParser) *taskDefinition {
	if _, ok := taskDefinitions[key]; ok {
		panic(fmt.Sprint("Can't re-register a taskDefinition:", name))
	}
	if key == "" || name == "" {
		panic("key and name cannot be empty")
	}

	td := &taskDefinition{key: key, name: name, parser: parserFunc}
	taskDefinitions[key] = td
	return td
}

func (t *taskDefinition) SetRequiredTask(name string) *taskDefinition {
	if name == "" {
		panic("name cannot be empty")
	}
	t.requiredTask = name
	return t
}

func (t *taskDefinition) SetOsRequirement(requirement string) *taskDefinition {
	if requirement == "" {
		panic("requirement cannot be empty")
	}
	t.osRequirement = requirement
	return t
}
