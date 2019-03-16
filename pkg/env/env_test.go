package env

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnviron(t *testing.T) {
	values := []string{"K1=V1", "K2=V2", "K3=V3"}
	require.ElementsMatch(t, values, New(values).Environ())
}

func TestSetUnset(t *testing.T) {
	env := New([]string{})

	env.Unset("NOOP")
	require.ElementsMatch(t, []string{}, env.Environ())

	env.Set("K1", "V1")
	require.ElementsMatch(t, []string{"K1=V1"}, env.Environ())

	env.Set("K1", "V2")
	require.ElementsMatch(t, []string{"K1=V2"}, env.Environ())

	env.Set("K2", "V1")
	require.ElementsMatch(t, []string{"K1=V2", "K2=V1"}, env.Environ())

	env.Unset("K1")
	require.ElementsMatch(t, []string{"K2=V1"}, env.Environ())

	env.Unset("K2")
	require.ElementsMatch(t, []string{}, env.Environ())
}

func TestHas(t *testing.T) {
	env := New([]string{})

	env.Set("K1", "V1")
	require.True(t, env.Has("K1"))

	env.Unset("K1")
	require.False(t, env.Has("K1"))
}

func TestPATH(t *testing.T) {
	env := New([]string{})

	env.PrependToPath("/bin")
	require.Equal(t, "/bin", env.Get("PATH"))

	env.PrependToPath("/sbin")
	require.Equal(t, "/sbin:/bin", env.Get("PATH"))
}
