package tasks

type taskAction interface {
	description() string
	needed(*context) (bool, error)
	run(*context) error
}
