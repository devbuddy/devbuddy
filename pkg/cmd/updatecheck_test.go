package cmd

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

type updateCheckFinisherFunc func()

func (f updateCheckFinisherFunc) Finish() {
	f()
}

func TestMaybeRunUpdateCheckSkipsShellHook(t *testing.T) {
	var called bool

	maybeStartUpdateCheck([]string{"bud", "--shell-hook"}, "v0.16.1", &bytes.Buffer{}, func(string, io.Writer) updateCheckFinisher {
		called = true
		return nil
	})

	require.False(t, called)
}

func TestMaybeStartUpdateCheckOnlyRunsForBudUp(t *testing.T) {
	var called bool

	finisher := maybeStartUpdateCheck([]string{"bud", "up"}, "v0.16.1", &bytes.Buffer{}, func(version string, out io.Writer) updateCheckFinisher {
		called = true
		require.Equal(t, "v0.16.1", version)
		require.NotNil(t, out)
		return updateCheckFinisherFunc(func() {})
	})

	require.True(t, called)
	require.NotNil(t, finisher)

	called = false
	finisher = maybeStartUpdateCheck([]string{"bud", "inspect"}, "v0.16.1", &bytes.Buffer{}, func(string, io.Writer) updateCheckFinisher {
		called = true
		return updateCheckFinisherFunc(func() {})
	})
	require.False(t, called)
	require.Nil(t, finisher)
}
