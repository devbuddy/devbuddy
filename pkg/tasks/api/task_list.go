package api

import (
	"encoding/json"
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
		task, err = newEnvTask(manifest.Env)
		if err != nil {
			return nil, err
		}
		taskList = append(taskList, task)
	}

	for _, payload := range manifest.Up {
		task, err = NewTaskFromPayload(payload)
		if err != nil {
			return nil, err
		}
		taskList = append(taskList, task)
	}

	return taskList, nil
}

func NewTaskFromPayload(payload any) (*Task, error) {
	taskConfig, err := NewTaskConfig(payload)
	if err != nil {
		return nil, fmt.Errorf("parsing task: %w", err)
	}

	taskDef := GetDefinitionOrUnknown(taskConfig.name)

	task := &Task{TaskDefinition: taskDef}

	err = taskDef.Parser(taskConfig, task)
	if err != nil {
		return nil, fmt.Errorf(`task "%s": %w`, task.Key, err)
	}

	return task, nil
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

// newEnvTask builds an env task directly, encoding the env vars as JSON in the
// feature param so the env feature can activate without re-loading dev.yml.
func newEnvTask(envVars map[string]string) (*Task, error) {
	encoded, err := json.Marshal(envVars)
	if err != nil {
		return nil, fmt.Errorf("encoding env vars: %w", err)
	}
	taskDef := GetDefinitionOrUnknown("env")
	task := &Task{TaskDefinition: taskDef}
	noop := func(ctx *context.Context) error { return nil }
	task.AddActionBuilder("", noop).SetFeature("env", string(encoded))
	return task, nil
}
