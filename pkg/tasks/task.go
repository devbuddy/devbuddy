package tasks

import "github.com/devbuddy/devbuddy/pkg/features"

// Task represents a task created by a taskDefinition.parser and specified by the TaskInfo
type Task struct {
	*taskDefinition
	header  string
	actions []taskAction
	feature features.FeatureInfo
}

func (t *Task) addAction(action taskAction) {
	t.actions = append(t.actions, action)
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
