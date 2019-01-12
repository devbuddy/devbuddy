package taskapi

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/features"
)

// Task represents a task created by a taskDefinition.parser and specified by the TaskInfo
type Task struct {
	*TaskDefinition
	header  string
	actions []TaskAction
	feature features.FeatureInfo
}

func (t *Task) SetInfo(info string) {
	t.header = info
}

func (t *Task) SetFeature(name, param string) {
	t.feature = features.NewFeatureInfo(name, param)
}

func (t *Task) AddAction(action TaskAction) {
	t.actions = append(t.actions, action)
}

func (t *Task) AddActionWithBuilder(description string, runFunc func(*Context) error) *genericTaskActionBuilder {
	if runFunc == nil {
		panic("runFunc cannot be nil")
	}
	action := &genericTaskAction{desc: description, runFunc: runFunc}
	t.actions = append(t.actions, action)
	return &genericTaskActionBuilder{action}
}

func (t *Task) Describe() string {
	description := fmt.Sprintf("Task %s (%s)", t.Name, t.header)

	if t.RequiredTask != "" {
		description += fmt.Sprintf(" required_task=%s", t.RequiredTask)
	}

	if t.feature.Name != "" {
		description += fmt.Sprintf(" feature=%s:%s", t.feature.Name, t.feature.Param)
	}

	description += fmt.Sprintf(" actions=%d", len(t.actions))

	return description
}
