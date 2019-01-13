package taskapi

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/features"
)

// Task represents a task created by a taskDefinition.parser and specified by the TaskInfo
type Task struct {
	*TaskDefinition
	Info    string
	Actions []TaskAction
	Feature features.FeatureInfo
}

func (t *Task) SetFeature(name, param string) {
	t.Feature = features.NewFeatureInfo(name, param)
}

func (t *Task) AddAction(action TaskAction) {
	t.Actions = append(t.Actions, action)
}

func (t *Task) AddActionWithBuilder(description string, runFunc func(*Context) error) *genericTaskActionBuilder {
	if runFunc == nil {
		panic("runFunc cannot be nil")
	}
	action := &genericTaskAction{desc: description, runFunc: runFunc}
	t.Actions = append(t.Actions, action)
	return &genericTaskActionBuilder{action}
}

func (t *Task) Describe() string {
	description := fmt.Sprintf("Task %s (%s)", t.Name, t.Info)

	if t.RequiredTask != "" {
		description += fmt.Sprintf(" required_task=%s", t.RequiredTask)
	}

	if t.Feature.Name != "" {
		description += fmt.Sprintf(" feature=%s:%s", t.Feature.Name, t.Feature.Param)
	}

	description += fmt.Sprintf(" actions=%d", len(t.Actions))

	return description
}
