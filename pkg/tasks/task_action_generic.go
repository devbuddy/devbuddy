package tasks

import (
	"path/filepath"

	"github.com/devbuddy/devbuddy/pkg/helpers/store"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

type genericTaskAction struct {
	desc           string
	conditions     []*genericTaskActionCondition
	monitoredFiles []string
	runFunc        func(*Context) error

	runCalled bool
}

func (a *genericTaskAction) Description() string {
	return a.desc
}

func (a *genericTaskAction) Needed(ctx *Context) (result *ActionResult) {
	if a.runCalled {
		return a.post(ctx)
	}
	return a.pre(ctx)
}

func (a *genericTaskAction) Run(ctx *Context) error {
	a.runCalled = true
	return a.runFunc(ctx)
}

// internals

func (a *genericTaskAction) pre(ctx *Context) (result *ActionResult) {
	hasConditions := false

	for _, condition := range a.conditions {
		hasConditions = true

		result = condition.pre(ctx)
		if result.Error != nil || result.Needed {
			return result
		}
	}

	for _, filePath := range a.monitoredFiles {
		hasConditions = true

		result = genericTaskActionPreConditionForFile(ctx, filePath)
		if result.Error != nil || result.Needed {
			return result
		}
	}

	if hasConditions {
		return actionNotNeeded()
	}
	return actionNeeded("action without conditions")
}

func (a *genericTaskAction) post(ctx *Context) (result *ActionResult) {
	for _, condition := range a.conditions {
		result = condition.post(ctx)
		if result.Error != nil || result.Needed {
			return result
		}
	}
	for _, filePath := range a.monitoredFiles {
		result = genericTaskActionPostConditionForFile(ctx, filePath)
		if result.Error != nil || result.Needed {
			return result
		}
	}
	return actionNotNeeded()
}

func genericTaskActionPreConditionForFile(ctx *Context, path string) *ActionResult {
	fullPath := filepath.Join(ctx.proj.Path, path)

	if !utils.PathExists(fullPath) {
		return actionNotNeeded()
	}

	fileChecksum, err := utils.FileChecksum(fullPath)
	if err != nil {
		return actionFailed("failed to get the file checksum: %s", err)
	}

	storedChecksum, err := store.New(ctx.proj.Path).GetString("checksum" + path)
	if err != nil {
		return actionFailed("failed to read the previous file checksum: %s", err)
	}

	if fileChecksum != storedChecksum {
		return actionNeeded("file %s has changed", path)
	}
	return actionNotNeeded()
}

func genericTaskActionPostConditionForFile(ctx *Context, path string) *ActionResult {
	fullPath := filepath.Join(ctx.proj.Path, path)

	if !utils.PathExists(fullPath) {
		return actionNotNeeded()
	}

	fileChecksum, err := utils.FileChecksum(fullPath)
	if err != nil {
		return actionFailed("failed to get the file checksum: %s", err)
	}

	err = store.New(ctx.proj.Path).SetString("checksum"+path, fileChecksum)
	if err != nil {
		return actionFailed("failed to store the current file checksum: %s", err)
	}

	return actionNotNeeded()
}
