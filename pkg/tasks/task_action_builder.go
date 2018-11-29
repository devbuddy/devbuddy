package tasks

import (
	"path/filepath"

	"github.com/devbuddy/devbuddy/pkg/helpers/store"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

type actionRunFunc func(*context) error

type actionCondition struct {
	pre  func(*context) *actionResult
	post func(*context) *actionResult
}

type actionWithBuilder struct {
	desc       string
	conditions []*actionCondition
	runFunc    actionRunFunc
	ran        bool
}

// action API

func (a *actionWithBuilder) description() string {
	return a.desc
}

func (a *actionWithBuilder) needed(ctx *context) (result *actionResult) {
	if a.ran {
		return a.post(ctx)
	}
	return a.pre(ctx)
}

func (a *actionWithBuilder) pre(ctx *context) (result *actionResult) {
	if len(a.conditions) == 0 {
		return actionNeeded("")
	}
	for _, condition := range a.conditions {
		result = condition.pre(ctx)
		if result.Error != nil || result.Needed {
			return result
		}
	}
	return actionNotNeeded()
}

func (a *actionWithBuilder) post(ctx *context) (result *actionResult) {
	for _, condition := range a.conditions {
		result = condition.post(ctx)
		if result.Error != nil || result.Needed {
			return result
		}
	}
	return actionNotNeeded()
}

func (a *actionWithBuilder) run(ctx *context) error {
	a.ran = true
	return a.runFunc(ctx)
}

// Builder API

func newAction(description string, runFunc actionRunFunc) *actionWithBuilder {
	return &actionWithBuilder{desc: description, runFunc: runFunc}
}

func (a *actionWithBuilder) on(condition *actionCondition) *actionWithBuilder {
	a.conditions = append(a.conditions, condition)
	return a
}

func (a *actionWithBuilder) onFunc(condFunc func(*context) *actionResult) *actionWithBuilder {
	a.on(&actionCondition{pre: condFunc, post: condFunc})
	return a
}

func (a *actionWithBuilder) onFileChange(path string) *actionWithBuilder {
	a.on(&actionCondition{
		pre: func(ctx *context) *actionResult {
			fileChecksum, err := utils.FileChecksum(filepath.Join(ctx.proj.Path, path))
			if err != nil {
				return actionFailed("failed to get the file checksum: %s", err)
			}

			storedChecksum, err := ctx.store.GetString("checksum", store.KeyFromPath(path))
			if err != nil {
				return actionFailed("failed to read the previous file checksum: %s", err)
			}

			if fileChecksum != storedChecksum {
				return actionNeeded("file %s changed", path)
			}
			return actionNotNeeded()
		},

		post: func(ctx *context) *actionResult {
			fileChecksum, err := utils.FileChecksum(filepath.Join(ctx.proj.Path, path))
			if err != nil {
				return actionFailed("failed to get the file checksum: %s", err)
			}

			err = ctx.store.SetString("checksum", store.KeyFromPath(path), fileChecksum)
			if err != nil {
				return actionFailed("failed to store the current file checksum: %s", err)
			}

			return actionNotNeeded()
		},
	})
	return a
}

func (a *actionWithBuilder) onFeatureChange(name string) *actionWithBuilder {
	a.on(&actionCondition{
		pre: func(ctx *context) *actionResult {

			return actionNotNeeded()
		},
		post: func(ctx *context) *actionResult {

			return actionNotNeeded()
		},
	})
	return a
}
