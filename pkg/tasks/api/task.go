package api

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/context"
)

// Task represents a task created by a taskDefinition.parser and specified by the TaskInfo
type Task struct {
	*TaskDefinition
	Info    string
	Actions []TaskAction
}

func (t *Task) AddAction(action TaskAction) {
	t.Actions = append(t.Actions, action)
}

func (t *Task) AddActionBuilder(description string, runFunc func(*context.Context) error) *taskActionBuilder {
	if runFunc == nil {
		panic("runFunc cannot be nil")
	}
	action := &taskAction{desc: description, runFunc: runFunc}
	t.Actions = append(t.Actions, action)

	return &taskActionBuilder{action}
}

func (t *Task) Describe() string {
	description := fmt.Sprintf("Task %s (%s)", t.Name, t.Info)

	if t.RequiredTask != "" {
		description += fmt.Sprintf(" required_task=%s", t.RequiredTask)
	}

	for _, action := range t.Actions {
		f := action.Feature()
		if action.Feature() != nil {
			feature := *f
			description += fmt.Sprintf(" feature=%s:%s", feature.Name, feature.Param)
		}
	}

	description += fmt.Sprintf(" actions=%d", len(t.Actions))

	return description
}
