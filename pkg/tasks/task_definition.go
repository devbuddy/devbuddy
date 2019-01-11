package tasks

import (
	"fmt"
)

type taskParser func(*TaskConfig, *Task) error

type taskDefinition struct {
	name          string
	requiredTask  string
	parser        taskParser
	osRequirement string // The platform this task can run on. "debian", "macos"
}

var taskDefinitions = make(map[string]*taskDefinition)

func registerTaskDefinition(name string) *taskDefinition {
	if _, ok := taskDefinitions[name]; ok {
		panic(fmt.Sprint("Can't re-register a taskDefinition:", name))
	}
	td := &taskDefinition{name: name}
	taskDefinitions[name] = td
	return td
}
