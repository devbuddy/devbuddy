package hook

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/autoenv"
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/tasks/api"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

func Run() {
	// In the shell hook, the stdout is evaluated by the shell
	// stderr is used to display messages to the user

	// Also, we can't annoy the user here, so we always just quit silently

	ctx, err := context.Load(true)
	if err != nil {
		return
	}

	err = run(ctx)
	if err != nil {
		ctx.UI.Debug("%s", err)
	}
}

func run(ctx *context.Context) error {
	allFeatures, err := getDesiredFeatures(ctx)
	if err != nil {
		return err
	}
	ctx.UI.Debug("Desired features: %+v", allFeatures)

	autoenv.Sync(ctx, allFeatures)
	emitEnvironmentChangeAsShellCommands(ctx)
	emitShellHashResetCommand(ctx)

	return nil
}

func getDesiredFeatures(ctx *context.Context) (autoenv.FeatureSet, error) {
	if ctx.Project == nil {
		return autoenv.NewFeatureSet(), nil
	}

	cache := autoenv.ReadFeatureCache(ctx.Env)

	if cache != nil && cache.ProjectSlug == ctx.Project.Slug() {
		// Cache hit: use cached features, but check if dev.yml changed
		checksum, err := utils.FileChecksum(filepath.Join(ctx.Project.Path, "dev.yml"))
		if err == nil && checksum != cache.Checksum && cache.ShouldWarnDevYmlChanged() {
			ctx.UI.HookDevYmlChanged()
			cache.MarkWarned()
			autoenv.WriteFeatureCache(ctx.Env, cache)
		}
		return cache.Features, nil
	}

	// No cache for this project (first visit in session): parse once and cache
	features, err := loadFeaturesFromProject(ctx.Project)
	if err != nil {
		return nil, err
	}
	checksum, _ := utils.FileChecksum(filepath.Join(ctx.Project.Path, "dev.yml"))
	autoenv.WriteFeatureCache(ctx.Env, autoenv.NewFeatureCache(ctx.Project.Slug(), checksum, features))
	return features, nil
}

func loadFeaturesFromProject(proj *project.Project) (autoenv.FeatureSet, error) {
	allTasks, err := api.GetTasksFromProject(proj)
	if err != nil {
		return nil, err
	}
	return api.GetFeaturesFromTasks(allTasks), nil
}

func emitEnvironmentChangeAsShellCommands(ctx *context.Context) {
	for _, mutation := range ctx.Env.Mutations() {
		ctx.UI.Debug("Apply change: %s\n%s", mutation.Name, mutation.DiffString())

		if mutation.Current == nil {
			fmt.Printf("unset %s\n", mutation.Name)
		} else {
			fmt.Printf("export %s=\"%s\"\n", mutation.Name, shellEscapeDoubleQuoted(mutation.Current.Value))
		}
	}
}

// shellEscapeDoubleQuoted escapes a string for safe inclusion in a bash
// double-quoted context. In double quotes, \, ", $, and ` are special.
func shellEscapeDoubleQuoted(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, `$`, `\$`)
	s = strings.ReplaceAll(s, "`", "\\`")
	return s
}

func emitShellHashResetCommand(ctx *context.Context) {
	for _, mutation := range ctx.Env.Mutations() {
		if mutation.Name == "PATH" {
			fmt.Printf("hash -r")
			return
		}
	}
}
