package taskengine

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/helpers/osidentity"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
)

type TaskSelector interface {
	ShouldRun(*taskapi.Context, *taskapi.Task) (bool, error)
}

type TaskSelectorImpl struct {
	osIdent *osidentity.Identity
}

func NewTaskSelector() *TaskSelectorImpl {
	return &TaskSelectorImpl{osIdent: osidentity.Detect()}
}

func (s *TaskSelectorImpl) ShouldRun(ctx *taskapi.Context, task *taskapi.Task) (bool, error) {
	shouldRun, err := s.osRequirementMatch(ctx, task)
	if err != nil {
		return false, err
	}
	if !shouldRun {
		return false, nil
	}

	return true, nil
}

func (s *TaskSelectorImpl) osRequirementMatch(ctx *taskapi.Context, task *taskapi.Task) (bool, error) {
	switch task.OSRequirement {
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
		return false, fmt.Errorf("invalid value for osRequirement: %s", task.OSRequirement)
	}

	return true, nil
}
