package tasks

import (
	"fmt"
	"reflect"

	"github.com/devbuddy/devbuddy/pkg/features"
	"github.com/devbuddy/devbuddy/pkg/manifest"
	"github.com/devbuddy/devbuddy/pkg/project"
)

func GetTasksFromProject(proj *project.Project) (taskList []*Task, err error) {
	var task *Task

	manifest, err := manifest.Load(proj.Path)
	if err != nil {
		return nil, err
	}

	for _, taskdef := range manifest.Up {
		task, err = buildFromDefinition(taskdef)
		if err != nil {
			return nil, err
		}
		taskList = append(taskList, task)
	}

	return taskList, nil
}

func buildFromDefinition(definition interface{}) (task *Task, err error) {
	taskConfig, err := parseTaskConfig(definition)
	if err != nil {
		return nil, err
	}

	taskDef := taskDefinitions[taskConfig.name]
	if taskDef == nil {
		taskDef = &TaskDefinition{
			Name:   "Unknown",
			Parser: parseUnknown,
		}
	}

	task = &Task{TaskDefinition: taskDef}
	err = taskDef.Parser(taskConfig, task)
	return
}

func parseTaskConfig(definition interface{}) (*TaskConfig, error) {
	val := reflect.ValueOf(definition)

	if val.Kind() == reflect.Map {
		keys := val.MapKeys()
		if len(keys) != 1 {
			return nil, fmt.Errorf("invalid map length")
		}
		name, ok := keys[0].Interface().(string)
		if !ok {
			return nil, fmt.Errorf("task name should be a string")
		}
		payload := val.MapIndex(keys[0]).Interface()
		return &TaskConfig{name: name, payload: payload}, nil
	}

	if val.Kind() == reflect.String {
		return &TaskConfig{name: definition.(string), payload: nil}, nil
	}

	return nil, fmt.Errorf("invalid task: \"%+v\"", definition)
}

func GetFeaturesFromTasks(tasks []*Task) features.FeatureSet {
	featureSet := features.FeatureSet{}

	for _, task := range tasks {
		if task.feature.Name != "" {
			featureSet = featureSet.With(task.feature)
		}
	}

	return featureSet
}

func InspectTasks(taskList []*Task, proj *project.Project) (s string) {
	for _, task := range taskList {
		s += fmt.Sprintf("Task %s (%s)\n", task.Name, task.header)
		if task.feature.Name != "" {
			s += fmt.Sprintf("  Provides: %s\n", task.feature)
		}
		if task.RequiredTask != "" {
			s += fmt.Sprintf("  Requires: %s\n", task.RequiredTask)
		}
		for _, action := range task.actions {
			s += fmt.Sprintf("  Action: %T %+v\n", action, action)
		}
	}
	return s
}
