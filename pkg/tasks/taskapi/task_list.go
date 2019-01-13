package taskapi

import (
	"fmt"

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
		task, err = NewTaskFromDefinition(taskdef)
		if err != nil {
			return nil, err
		}
		taskList = append(taskList, task)
	}

	return taskList, nil
}

func NewTaskFromDefinition(definition interface{}) (task *Task, err error) {
	taskConfig, err := NewTaskConfig(definition)
	if err != nil {
		return nil, err
	}

	taskDef := taskDefinitions[taskConfig.name]
	if taskDef == nil {
		taskDef = newUnknownTaskDefinition()
	}

	task = &Task{TaskDefinition: taskDef}
	err = taskDef.Parser(taskConfig, task)
	return
}

func newUnknownTaskDefinition() *TaskDefinition {
	parser := func(config *TaskConfig, task *Task) error {
		task.AddActionWithBuilder("", func(ctx *Context) error {
			ctx.UI.TaskWarning(fmt.Sprintf("Unknown task: \"%s\"", config.name))
			return nil
		})
		return nil
	}
	return &TaskDefinition{Name: "Unknown", Parser: parser}
}

func GetFeaturesFromTasks(tasks []*Task) features.FeatureSet {
	featureSet := features.FeatureSet{}

	for _, task := range tasks {
		if task.Feature.Name != "" {
			featureSet = featureSet.With(task.Feature)
		}
	}

	return featureSet
}
