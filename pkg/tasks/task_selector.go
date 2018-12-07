package tasks

import (
	"fmt"

	osidentity "github.com/devbuddy/devbuddy/pkg/helpers/os"
)

type TaskSelector interface {
	ShouldRun(*Context, *Task) error
}

type TaskSelectorImpl struct {
	osIdent osidentity.Identity
}

func NewTaskSelector() (*TaskSelectorImpl, error) {
	osIdent, err := osidentity.Detect()
	if err != nil {
		return nil, err
	}
	return nil, &TaskSelectorImpl{osIdent: osIdent}
}

func (s *TaskSelectorImpl) ShouldRun(ctx *Context, task *Task) (bool, error) {
	shouldRun, err := s.osRequirementMatch(ctx, task)
	if err != nil {
		return false, err
	}
	if !shouldRun {
		return false, nil
	}

	return true, nil
}

func (s *TaskSelectorImpl) osRequirementMatch(ctx *Context, task *Task) (bool, error) {
	switch task.osRequirement {
	case "":
		return true, nil
	case "debian":
		if !ident.isDebianLike() {
			return false, nil
		}
	case "macos":
		if !ident.isMacOS() {
			return false, nil
		}
	default:
		return false, fmt.Errorf("invalid value for osRequirement: %s", task.osRequirement)
	}

	return true, nil
}
