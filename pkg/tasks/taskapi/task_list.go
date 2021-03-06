package taskapi

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/autoenv"
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/manifest"
	"github.com/devbuddy/devbuddy/pkg/project"
)

func GetTasksFromProject(proj *project.Project) (taskList []*Task, err error) {
	var task *Task

	manifest, err := manifest.Load(proj.Path)
	if err != nil {
		return nil, err
	}

	if len(manifest.Env) != 0 {
		task, err = NewTaskFromDefinition("env")
		if err != nil {
			return nil, err
		}
		taskList = append(taskList, task)
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
		task.AddActionWithBuilder("", func(ctx *context.Context) error {
			ctx.UI.TaskWarning(fmt.Sprintf("Unknown task: \"%s\"", config.name))
			return nil
		})
		return nil
	}
	return &TaskDefinition{Name: "Unknown", Parser: parser}
}

func GetFeaturesFromTasks(tasks []*Task) autoenv.FeatureSet {
	featureSet := autoenv.NewFeatureSet()

	for _, task := range tasks {
		for _, action := range task.Actions {
			if action.Feature() != nil {
				featureSet = featureSet.With(action.Feature())
			}
		}
	}

	return featureSet
}
