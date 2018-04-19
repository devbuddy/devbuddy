package tasks

import (
	"fmt"

	"github.com/pior/dad/pkg/project"
)

var allTasks = make(map[string]taskBuilder)

type taskBuilder func() Task

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
	taskConfig, err := parseTaskConfig(definition)
	if err == nil {
		taskBuilder := allTasks[taskConfig.name]
		if taskBuilder == nil {
			taskBuilder = newUnknown
		}
		task = taskBuilder()
	} else {
		task = newInvalid(definition, err)
	}

	ok, err := task.load(taskConfig)
	if err != nil {
		return nil, err
	}
	if ok {
		return task, nil
	}

	return nil, fmt.Errorf("error parsing tasks: %+v", definition)
}
