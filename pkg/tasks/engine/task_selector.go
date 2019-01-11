package engine

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/helpers/osidentity"
	"github.com/devbuddy/devbuddy/pkg/tasks"
)

type TaskSelector interface {
	ShouldRun(*tasks.Context, *tasks.Task) (bool, error)
}

type TaskSelectorImpl struct {
	osIdent *osidentity.Identity
}

func NewTaskSelector() *TaskSelectorImpl {
	return &TaskSelectorImpl{osIdent: osidentity.Detect()}
}

func (s *TaskSelectorImpl) ShouldRun(ctx *tasks.Context, task *tasks.Task) (bool, error) {
	shouldRun, err := s.osRequirementMatch(ctx, task)
	if err != nil {
		return false, err
	}
	if !shouldRun {
		return false, nil
	}

	return true, nil
}

func (s *TaskSelectorImpl) osRequirementMatch(ctx *tasks.Context, task *tasks.Task) (bool, error) {
	switch task.osRequirement {
	case "":
		break
	case "debian":
		if !s.osIdent.IsDebianLike() {
			return false, nil
		}
	case "macos":
		if !s.osIdent.IsMacOS() {
			return false, nil
		}
	default:
		return false, fmt.Errorf("invalid value for osRequirement: %s", task.osRequirement)
	}

	return true, nil
}
