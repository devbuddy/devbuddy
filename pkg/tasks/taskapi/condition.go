package taskapi

import (
	"path/filepath"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers/store"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

type Condition interface {
	Before(*context.Context) *ActionResult
	After(*context.Context) *ActionResult
}

type funcCondition struct {
	fn func(*context.Context) *ActionResult
}

func FuncCondition(fn func(*context.Context) *ActionResult) Condition {
	return &funcCondition{fn}
}

func (c funcCondition) Before(ctx *context.Context) *ActionResult {
	return c.fn(ctx)
}

func (c funcCondition) After(ctx *context.Context) *ActionResult {
	return c.fn(ctx)
}

type fileCondition struct {
	path string
}

func FileCondition(path string) Condition {
	return fileCondition{path}
}

func (c fileCondition) Before(ctx *context.Context) *ActionResult {
	fullPath := filepath.Join(ctx.Project.Path, c.path)

	if !utils.PathExists(fullPath) {
		return NotNeeded()
	}

	fileChecksum, err := utils.FileChecksum(fullPath)
	if err != nil {
		return Failed("failed to get the file checksum: %s", err)
	}

	checksumStore, err := store.Open(ctx.Project.Path, "checksum")
	if err != nil {
		return Failed("failed to open the internal project state: %s", err)
	}

	storedChecksum, err := checksumStore.GetString(c.path)
	if err != nil {
		return Failed("failed to read the previous file checksum: %s", err)
	}

	if fileChecksum != storedChecksum {
		return Needed("file %s has changed", c.path)
	}
	return NotNeeded()
}

func (c fileCondition) After(ctx *context.Context) *ActionResult {
	fullPath := filepath.Join(ctx.Project.Path, c.path)

	if !utils.PathExists(fullPath) {
		return NotNeeded()
	}

	fileChecksum, err := utils.FileChecksum(fullPath)
	if err != nil {
		return Failed("failed to get the file checksum: %s", err)
	}

	checksumStore, err := store.Open(ctx.Project.Path, "checksum")
	if err != nil {
		return Failed("failed to open the internal project state: %s", err)
	}

	err = checksumStore.SetString(c.path, fileChecksum)
	if err != nil {
		return Failed("failed to store the current file checksum: %s", err)
	}

	return NotNeeded()
}
