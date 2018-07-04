package tasks

import (
	"fmt"
)

type taskDefinition struct {
	key             string
	name            string
	requiredFeature string
	builder         taskBuilder
}

var taskDefinitions = make(map[string]*taskDefinition)

func registerTask(name string) *taskDefinition {
	if _, ok := taskDefinitions[name]; ok {
		panic(fmt.Sprint("Can't re-register a taskDefinition:", name))
	}
	td := &taskDefinition{name: name}
	taskDefinitions[name] = td
	return td
}
