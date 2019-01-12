package taskapi

type genericTaskActionCondition struct {
	pre  func(*Context) *ActionResult
	post func(*Context) *ActionResult
}

type genericTaskActionBuilder struct {
	*genericTaskAction
}

func actionBuilder(description string, runFunc func(*Context) error) *genericTaskActionBuilder {
	return &genericTaskActionBuilder{&genericTaskAction{desc: description, runFunc: runFunc}}
}

// On registers a new condition
func (a *genericTaskActionBuilder) On(condition *genericTaskActionCondition) *genericTaskActionBuilder {
	a.conditions = append(a.conditions, condition)
	return a
}

// OnFunc registers a condition defined as a single function
func (a *genericTaskActionBuilder) OnFunc(condFunc func(*Context) *ActionResult) *genericTaskActionBuilder {
	a.On(&genericTaskActionCondition{pre: condFunc, post: condFunc})
	return a
}

// OnFileChange specifies that the action will run when a file changes.
// The action will NOT run if the file does not exist.
// The action will NOT fail if the file is not created.
func (a *genericTaskActionBuilder) OnFileChange(path string) *genericTaskActionBuilder {
	a.monitoredFiles = append(a.monitoredFiles, path)
	return a
}

// Build returns a new task action with the behaviour specified by the builder
func (a *genericTaskActionBuilder) Build() *genericTaskAction {
	return a.genericTaskAction
}
