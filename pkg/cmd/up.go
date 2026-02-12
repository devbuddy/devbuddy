package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/autoenv"
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/tasks"
	"github.com/devbuddy/devbuddy/pkg/tasks/api"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskengine"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

func init() {
	tasks.RegisterTasks()
}

var upCmd = &cobra.Command{
	Use:          "up",
	Short:        "Ensure the project is up and running",
	RunE:         upRun,
	Args:         noArgs,
	GroupID:      "devbuddy",
	SilenceUsage: true,
}

func upRun(_ *cobra.Command, _ []string) error {
	ctx, err := context.LoadWithProject()
	if err != nil {
		return err
	}

	taskList, err := api.GetTasksFromProject(ctx.Project)
	if err != nil {
		return err
	}

	runner := taskengine.NewTaskRunner(ctx)
	selector := taskengine.NewTaskSelector()

	success, err := taskengine.Run(ctx, runner, selector, taskList)
	if err != nil {
		return err
	}
	if !success {
		return errTasksFailed
	}

	// Update the feature cache so the shell hook can skip re-parsing dev.yml
	if err := writeFeatureCache(ctx, taskList); err != nil {
		ctx.UI.Debug("failed to write feature cache: %s", err)
	}

	return nil
}

func writeFeatureCache(ctx *context.Context, taskList []*api.Task) error {
	featureSet := api.GetFeaturesFromTasks(taskList)
	checksum, err := utils.FileChecksum(filepath.Join(ctx.Project.Path, "dev.yml"))
	if err != nil {
		return fmt.Errorf("computing dev.yml checksum: %w", err)
	}
	cache := autoenv.NewFeatureCache(ctx.Project.Slug(), checksum, featureSet)
	return autoenv.WriteFeatureCacheFinalizer(cache)
}
