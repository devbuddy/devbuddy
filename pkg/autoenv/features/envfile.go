package features

import (
	"slices"
	"sort"
	"strings"

	"github.com/joho/godotenv"

	"github.com/devbuddy/devbuddy/pkg/context"
)

const envfileTrackedVarsKey = "__BUD_ENVFILE_VARS"

func init() {
	register.Register(envfile{})
}

type envfile struct{}

func (envfile) Name() string {
	return "envfile"
}

func (envfile) Activate(ctx *context.Context, param string) (bool, error) {
	newVars, err := godotenv.Read(param)
	if err != nil {
		return true, err
	}

	previousNames := parseTrackedVars(ctx.Env.Get(envfileTrackedVarsKey))

	// Set all variables from the file
	newNames := make([]string, 0, len(newVars))
	for name, value := range newVars {
		ctx.Env.Set(name, value)
		newNames = append(newNames, name)
	}
	sort.Strings(newNames)

	// Unset variables that were previously set but are no longer in the file
	for _, name := range previousNames {
		if !slices.Contains(newNames, name) {
			ctx.Env.Unset(name)
		}
	}

	ctx.Env.Set(envfileTrackedVarsKey, strings.Join(newNames, ","))

	return false, nil
}

func (envfile) Deactivate(ctx *context.Context, param string) {
	for _, name := range parseTrackedVars(ctx.Env.Get(envfileTrackedVarsKey)) {
		ctx.Env.Unset(name)
	}
	ctx.Env.Unset(envfileTrackedVarsKey)
}

func (envfile) WatchedFiles(param string) []string {
	return []string{param}
}

func parseTrackedVars(raw string) []string {
	if raw == "" {
		return nil
	}
	return strings.Split(raw, ",")
}
