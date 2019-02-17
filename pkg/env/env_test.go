package env

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChanged(t *testing.T) {
	env := New([]string{})
	require.Equal(t, []VariableChange{}, env.Changed())

	env.Set("K2", "1")
	require.Equal(t, []VariableChange{{"K2", "1", false}}, env.Changed())

	env.Set("K2", "2")
	require.Equal(t, []VariableChange{{"K2", "2", false}}, env.Changed())

	env = New([]string{"K1=1"})
	env.Unset("K1")
	require.Equal(t, []VariableChange{{"K1", "", true}}, env.Changed())
}

func TestEnviron(t *testing.T) {
	env := New([]string{"K1=V1", "K2=V2", "K3=V3"})
	env.Set("K1", "V1B")
	env.Unset("K2")
	require.ElementsMatch(t, []string{"K1=V1B", "K3=V3"}, env.Environ())
}

func TestHas(t *testing.T) {
	env := New([]string{"K1=V1"})
	require.True(t, env.Has("K1"))
	require.False(t, env.Has("NOPE"))
}

func TestIsInPath(t *testing.T) {
	env := New([]string{"PATH=/usr/bin:/bin"})

	require.True(t, env.IsInPath("/usr/bin"))
	require.True(t, env.IsInPath("/bin"))

	require.False(t, env.IsInPath("/usr/local/bin"))
	require.False(t, env.IsInPath("/usr"))
}
