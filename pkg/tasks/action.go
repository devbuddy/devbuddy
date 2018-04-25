package tasks

type taskAction interface {
	description() string
	needed(*Context) (bool, error)
	run(*Context) error
}
