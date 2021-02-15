package expect

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUnknownCommand(t *testing.T) {
	ep := NewExpect("/nopenopenope")
	err := ep.Start()
	require.EqualError(t, err, "fork/exec /nopenopenope: no such file or directory")
}

func TestFailingCommand(t *testing.T) {
	ep := NewExpect("false")
	err := ep.Start()
	require.NoError(t, err)

	_, err = ep.Expect("")
	require.EqualError(t, err, "process stopped: exit status 1")
}

func TestExpectFunc(t *testing.T) {
	ep := NewExpect("/bin/echo", "hello world")
	err := ep.Start()
	require.NoError(t, err)

	wstr := "hello world\n"
	l, eerr := ep.ExpectFunc(func(a string) bool { return len(a) > 10 })
	require.NoError(t, eerr)
	require.Equal(t, wstr, l)

	cerr := ep.Close()
	require.NoError(t, cerr)
}

func TestEcho(t *testing.T) {
	ep := NewExpect("/bin/echo", "hello world")
	err := ep.Start()
	require.NoError(t, err)

	l, eerr := ep.Expect("world")
	require.NoError(t, eerr)

	wstr := "hello world"
	require.Equal(t, wstr, l[:len(wstr)])

	cerr := ep.Close()
	require.NoError(t, cerr)
}

func TestExited(t *testing.T) {
	ep := NewExpect("/bin/echo", "")
	err := ep.Start()
	require.NoError(t, err)

	l, err := ep.Expect("foobar")
	require.EqualError(t, err, "process stopped")
	require.Equal(t, "", l)

	cerr := ep.Close()
	require.NoError(t, cerr)
}

func TestClose(t *testing.T) {
	ep := NewExpect("/bin/echo", "")
	err := ep.Start()
	require.NoError(t, err)

	err = ep.Close()
	require.NoError(t, err)

	err = ep.Close()
	require.EqualError(t, err, "already closed")
}

func TestStop(t *testing.T) {
	ep := NewExpect("/bin/sleep", "100")
	err := ep.Start()
	require.NoError(t, err)

	err = ep.Stop()
	require.NoError(t, err)

	err = ep.Stop()
	require.EqualError(t, err, "already closed")
}

func TestExpectOnClosed(t *testing.T) {
	ep := NewExpect("/bin/echo", "hello world")
	err := ep.Start()
	require.NoError(t, err)

	err = ep.Close()
	require.NoError(t, err)

	_, err = ep.Expect("...")
	require.EqualError(t, err, "already closed")
}

func TestSend(t *testing.T) {
	ep := NewExpect("/usr/bin/tr", "a", "b")
	err := ep.Start()
	require.NoError(t, err)

	err = ep.Send("a\n")
	require.NoError(t, err)

	_, err = ep.Expect("b")
	require.NoError(t, err)

	err = ep.Stop()
	require.NoError(t, err)
}

func TestSignal(t *testing.T) {
	ep := NewExpect("/bin/sleep", "100")
	err := ep.Start()
	require.NoError(t, err)

	ep.Signal(os.Interrupt)

	done := make(chan struct{})
	go func() {
		defer close(done)

		cerr := ep.Close()
		require.EqualError(t, cerr, "process stopped: signal: interrupt")
	}()

	select {
	case <-time.After(5 * time.Second):
		t.Fatalf("signal test timed out")
	case <-done:
	}
}
