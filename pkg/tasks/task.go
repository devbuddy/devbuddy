package tasks

import "github.com/devbuddy/devbuddy/pkg/features"

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

func buildFromDefinition(definition interface{}) (task *Task, err error) {
	taskConfig, err := parseTaskConfig(definition)
	if err != nil {
		return nil, err
	}

	taskDef := taskDefinitions[taskConfig.name]
	if taskDef == nil {
		taskDef = &taskDefinition{
			name:   "Unknown",
			parser: parseUnknown,
		}
	}

	task = &Task{taskDefinition: taskDef}
	err = taskDef.parser(taskConfig, task)
	return
}
