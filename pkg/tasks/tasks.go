package tasks

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/manifest"
	"github.com/devbuddy/devbuddy/pkg/project"
)

type taskParser func(*taskConfig, *Task) error

type taskDefinition struct {
	// key          string
	name          string
	requiredTask  string
	parser        taskParser
	osRequirement string // The platform this task can run on. "debian", "macos"
}

var taskDefinitions = make(map[string]*taskDefinition)

func registerTaskDefinition(name string) *taskDefinition {
	if _, ok := taskDefinitions[name]; ok {
		panic(fmt.Sprint("Can't re-register a taskDefinition:", name))
	}
	td := &taskDefinition{name: name}
	taskDefinitions[name] = td
	return td
}

// Task represents a task created by a taskDefinition.parser and specified by the TaskInfo
type Task struct {
	*taskDefinition
	header       string
	actions      []taskAction
	perform      func(*Context) error
	featureName  string
	featureParam string
}

func (t *Task) addAction(action taskAction) {
	t.actions = append(t.actions, action)
}

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

func GetFeaturesFromTasks(tasks []*Task) map[string]string {
	features := map[string]string{}

	for _, task := range tasks {
		if task.featureName != "" {
			features[task.featureName] = task.featureParam
		}
	}

	return features
}

func InspectTasks(taskList []*Task, proj *project.Project) (s string) {
	for _, task := range taskList {
		s += fmt.Sprintf("Task %s (%s)\n", task.name, task.header)
		if task.featureName != "" {
			s += fmt.Sprintf("  Provides: %s=%s\n", task.featureName, task.featureParam)
		}
		if task.requiredTask != "" {
			s += fmt.Sprintf("  Requires: %s\n", task.requiredTask)
		}
		for _, action := range task.actions {
			s += fmt.Sprintf("  Action: %T %+v\n", action, action)
		}
	}
	return s
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
