package tasks

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/project"
)

type taskParser func(*taskConfig, *Task) error

type taskDefinition struct {
	key             string
	name            string
	requiredFeature string
	parser          taskParser
}

var taskDefinitions = make(map[string]*taskDefinition)

func registerTask(name string) *taskDefinition {
	if _, ok := taskDefinitions[name]; ok {
		panic(fmt.Sprint("Can't re-register a taskDefinition:", name))
	}
	td := &taskDefinition{name: name}
	taskDefinitions[name] = td
	return td
}

type Task struct {
	// *taskDefinition
	header       string
	actions      []taskAction
	perform      func(*context) error
	featureName  string
	featureParam string
}

func (t *Task) addAction(action taskAction) {
	t.actions = append(t.actions, action)
}

type taskAction interface {
	description() string
	needed(*context) (bool, error)
	run(*context) error
}

func GetTasksFromProject(proj *project.Project) (taskList []Task, err error) {
	var task *Task

	for _, taskdef := range proj.Manifest.Up {
		task, err = buildFromDefinition(taskdef)
		if err != nil {
			return nil, err
		}
		taskList = append(taskList, task)
	}

	return taskList, nil
}

func GetFeaturesFromTasks(proj *project.Project, tasks []*Task) map[string]string {
	features := map[string]string{}

	for _, task := range tasks {
		if task.featureName != "" {
			feature, param := t.feature(proj)
			features[task.featureName] = task.featureParam
		}
	}

	return features
}

func InspectTasks(taskList []*Task, proj *project.Project) (s string) {
	for _, task := range taskList {
		s += fmt.Sprintf("Task %s\n", task.name())
		s += fmt.Sprintf("  Internal: %+v\n", task)
		if task.featureName != "" {
			s += fmt.Sprintf("  Feature: %s=%s\n", task.featureName, task.featureParam)
		}
	}
	return s
}

func buildFromDefinition(definition interface{}) (task *Task, err error) {
	taskConfig, err := parseTaskConfig(definition)
	if err != nil {
		return newInvalid(definition, err), nil
	}

	taskDef := taskDefinitions[taskConfig.name]
	if taskDef == nil {
		taskDef = taskDefinitions["unknown"]
	}

	// task.taskDefinition = taskDef

	err = taskDef.parser(taskConfig, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}
