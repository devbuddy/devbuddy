package taskapi

import "fmt"

type taskParser func(*TaskConfig, *Task) error

type taskDefinition struct {
	name          string
	requiredTask  string
	parser        taskParser
	osRequirement string // The platform this task can run on. "debian", "macos"
}

var taskDefinitions = make(map[string]*taskDefinition)

func RegisterTaskDefinition(key string, name string, parser taskParser) *taskDefinition {
	if _, ok := taskDefinitions[key]; ok {
		panic(fmt.Sprint("Can't re-register a TaskDefinition:", name))
	}
	td := &taskDefinition{name: name}
	taskDefinitions[name] = td
	return td
}

func (td *taskDefinition) AddRequiredTask(name string) *taskDefinition {
	if td.requiredTask != "" {
		panic("only one required task supported")
	}
	td.requiredTask = name
	return td
}

func (td *taskDefinition) SetOsRequirement(requirement string) *taskDefinition {
	td.osRequirement = requirement
	return td
}
