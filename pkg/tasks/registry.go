package tasks

import (
	"fmt"
	"reflect"

	"github.com/pior/dad/pkg/project"
)

var allTasks = make(map[string]TaskBuilder)

type TaskBuilder func() Task

func GetTasksFromProject(proj *project.Project) (taskList []Task, err error) {
	var task Task

	for _, taskdef := range proj.Manifest.Up {
		task, err = buildFromDefinition(taskdef)
		if err != nil {
			return nil, err
		}
		taskList = append(taskList, task)
	}

	return taskList, nil
}

func GetFeaturesFromTasks(proj *project.Project, tasks []Task) map[string]string {
	features := map[string]string{}

	for _, task := range tasks {
		if t, ok := task.(TaskWithFeature); ok {
			feature, param := t.Feature(proj)
			features[feature] = param
		}
	}

	return features
}

func buildFromDefinition(definition interface{}) (task Task, err error) {
	taskBuilder, err := findTaskBuilder(definition)
	if err != nil {
		return nil, err
	}
	task = taskBuilder()
	ok, err := task.Load(definition)
	if err != nil {
		return nil, err
	}
	if ok {
		return task, nil
	}

	return nil, fmt.Errorf("error parsing tasks: %+v", definition)
}

func findTaskBuilder(definition interface{}) (TaskBuilder, error) {
	name, err := extractTaskName(definition)
	if err != nil {
		return nil, fmt.Errorf("%s (%+v)", err, definition)
	}
	taskBuilder, found := allTasks[name]
	if found {
		return taskBuilder, nil
	}
	return NewUnknown, nil
}

func extractTaskName(definition interface{}) (string, error) {
	val := reflect.ValueOf(definition)
	if val.Kind() == reflect.Map {
		keys := val.MapKeys()
		if len(keys) != 1 {
			return "", fmt.Errorf("invalid map length")
		}
		definition = keys[0].Interface()
	}

	if name, ok := definition.(string); ok {
		return name, nil
	}

	return "", fmt.Errorf("invalid structure")
}
