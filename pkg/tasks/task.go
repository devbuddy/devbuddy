package tasks

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/features"
)

// Task represents a task created by a taskDefinition.parser and specified by the TaskInfo
type Task struct {
	*taskDefinition
	header  string
	actions []taskAction
	feature features.FeatureInfo
}

func (t *Task) SetInfo(info string) {
	t.header = info
}

func (t *Task) SetFeature(name, param string) {
	t.feature = features.NewFeatureInfo(name, param)
}

func (t *Task) AddAction(action taskAction) {
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
	description := fmt.Sprintf("Task %s (%s)", t.name, t.header)

	if t.feature.Name != "" {
		description += fmt.Sprintf(" has feature %s:%s and", t.feature.Name, t.feature.Param)
	}

	description += fmt.Sprintf(" has %d actions", len(t.actions))

	return description
}
