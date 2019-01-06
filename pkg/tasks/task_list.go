package tasks

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
		task, err = buildFromDefinition(taskdef)
		if err != nil {
			return nil, err
		}
		taskList = append(taskList, task)
	}

	return taskList, nil
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
		s += fmt.Sprintf("Task %s (%s)\n", task.name, task.header)
		if task.feature.Name != "" {
			s += fmt.Sprintf("  Provides: %s\n", task.feature)
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
