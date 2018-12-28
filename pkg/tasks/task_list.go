package tasks

import (
	"fmt"

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
