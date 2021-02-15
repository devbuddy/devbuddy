package context

import "time"

func Timeout(dur time.Duration) runOptionsFn {
	return func(opt *runOptions) {
		opt.timeout = dur
	}
}

func ExitCode(exitCode int) runOptionsFn {
	return func(opt *runOptions) {
		opt.exitCode = exitCode
	}
}

type runOptions struct {
	timeout  time.Duration
	exitCode int
}

type runOptionsFn func(*runOptions)

func buildRunOptions(fns []runOptionsFn) *runOptions {
	options := &runOptions{
		timeout: 5 * time.Second,
	}
	for _, fn := range fns {
		fn(options)
	}
	return options
}
