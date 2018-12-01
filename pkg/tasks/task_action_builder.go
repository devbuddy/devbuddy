package tasks

type genericTaskActionCondition struct {
	pre  func(*context) *actionResult
	post func(*context) *actionResult
}

type genericTaskActionBuilder struct {
	desc string

	conditions     []*genericTaskActionCondition
	monitoredFiles []string

	runFunc func(*context) error
}

func actionBuilder(description string, runFunc func(*context) error) *genericTaskActionBuilder {
	return &genericTaskActionBuilder{desc: description, runFunc: runFunc}
}

func (a *genericTaskActionBuilder) On(condition *genericTaskActionCondition) *genericTaskActionBuilder {
	a.conditions = append(a.conditions, condition)
	return a
}

func (a *genericTaskActionBuilder) OnFunc(condFunc func(*context) *actionResult) *genericTaskActionBuilder {
	a.On(&genericTaskActionCondition{pre: condFunc, post: condFunc})
	return a
}

// OnFileChange specifies that the action will run when a file changes or does not exist.
// The action will NOT fail if the file is not created.
func (a *genericTaskActionBuilder) OnFileChange(path string) *genericTaskActionBuilder {
	a.monitoredFiles = append(a.monitoredFiles, path)
	return a
}

//
func (a *genericTaskActionBuilder) Build() *genericTaskAction {
	return &genericTaskAction{
		builder: &genericTaskActionBuilder{ // Hand made copy
			desc:           a.desc,
			conditions:     append(a.conditions[:0:0], a.conditions...),
			monitoredFiles: append(a.monitoredFiles[:0:0], a.monitoredFiles...),
			runFunc:        a.runFunc,
		},
	}
}
