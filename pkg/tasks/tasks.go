package tasks

import (
	"github.com/pior/dad/pkg/project"
)

type taskBuilder func(*taskConfig) (Task, error)

var allTasks = make(map[string]taskBuilder)

type Task interface {
	name() string
	header() string
	perform(*Context) error
}

type TaskWithFeature interface {
	Task
	feature(*project.Project) (string, string)
}

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
			feature, param := t.feature(proj)
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
		task, err = taskBuilder(taskConfig)
		if err != nil {
			return nil, err
		}
	} else {
		task = newInvalid(definition, err)
	}

	return task, nil
}
