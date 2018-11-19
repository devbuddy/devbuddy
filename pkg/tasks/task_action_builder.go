package tasks

import "github.com/devbuddy/devbuddy/pkg/helpers/store"

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

func (s *actionWithBuilder) description() string {
	return s.desc
}

func (s *actionWithBuilder) needed(ctx *context) (result *actionResult) {
	if s.ran {
		return s.post(ctx)
	}
	return s.pre(ctx)
}

func (s *actionWithBuilder) pre(ctx *context) (result *actionResult) {
	if len(s.conditions) == 0 {
		return actionNeeded("")
	}
	for _, condition := range s.conditions {
		result = condition.pre(ctx)
		if result.Error != nil || result.Needed {
			return result
		}
	}
	return actionNotNeeded()
}

func (s *actionWithBuilder) post(ctx *context) (result *actionResult) {
	for _, condition := range s.conditions {
		result = condition.post(ctx)
		if result.Error != nil || result.Needed {
			return result
		}
	}
	return actionNotNeeded()
}

func (s *actionWithBuilder) run(ctx *context) error {
	s.ran = true
	return s.runFunc(ctx)
}

func (s *actionWithBuilder) addFileChangeCondition(path string) *actionWithBuilder {

	pre := func(ctx *context) *actionResult {
		changed, err := store.New(ctx.proj.Path).HasFileChanged("setup.py")
		if err != nil {
			return actionFailed("failed to check if setup.py has changed: %s", err)
		}
		if changed {
			return actionNeeded("setup.py was modified")
		}
		return actionNotNeeded()
	}
	post := func(ctx *context) *actionResult {
		err := store.New(ctx.proj.Path).RecordFileChange("setup.py")
		// TODO ...
	}

	s.conditions = append(s.conditions, &actionCondition{pre: pre, post: post})
	return s
}

func (s *actionWithBuilder) addFeatureCondition(name string) *actionWithBuilder {
	return s
}

func (s *actionWithBuilder) addCustomCondition(condFunc func(*context) *actionResult) *actionWithBuilder {
	s.conditions = append(s.conditions, &actionCondition{pre: condFunc, post: condFunc})
	return s
}
