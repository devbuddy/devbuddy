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
	env := New([]string{"PATH=/bin"})

	env.PrependToPath("/sbin")
	require.Equal(t, "/sbin:/bin", env.Get("PATH"))

	env.PrependToPath("/usr/bin")
	require.Equal(t, "/usr/bin:/sbin:/bin", env.Get("PATH"))
}

func TestMutations(t *testing.T) {
	env := New([]string{"PATH=/bin", "K1=V1"})

	require.Equal(t, []VariableMutation{}, env.Mutations())

	env.Set("K2", "V1")
	require.Equal(t, []VariableMutation{
		VariableMutation{"K2", nil, &variable{"K2", "V1"}},
	}, env.Mutations())

	env.Set("K2", "V2")
	require.Equal(t, []VariableMutation{
		VariableMutation{"K2", nil, &variable{"K2", "V2"}},
	}, env.Mutations())

	env.Unset("K2")
	require.Equal(t, []VariableMutation{}, env.Mutations())

	env.Set("K1", "V2")
	require.Equal(t, []VariableMutation{
		VariableMutation{"K1", &variable{"K1", "V1"}, &variable{"K1", "V2"}},
	}, env.Mutations())

	env.Unset("K1")
	require.Equal(t, []VariableMutation{
		VariableMutation{"K1", &variable{"K1", "V1"}, nil},
	}, env.Mutations())

	env.Set("K1", "V1")
	require.Equal(t, []VariableMutation{}, env.Mutations())

	env.PrependToPath("/foo")
	require.Equal(t, []VariableMutation{
		VariableMutation{"PATH", &variable{"PATH", "/bin"}, &variable{"PATH", "/foo:/bin"}},
	}, env.Mutations())
}

func TestMutationsDiffString(t *testing.T) {
	m := VariableMutation{"NAME", &variable{"NAME", "V1"}, &variable{"NAME", "V2"}}
	require.Equal(t, "  - V1\n  + V2\n", m.DiffString())

	m = VariableMutation{"NAME", nil, &variable{"NAME", "V2"}}
	require.Equal(t, "  + V2\n", m.DiffString())

	m = VariableMutation{"NAME", &variable{"NAME", "V1"}, nil}
	require.Equal(t, "  - V1\n", m.DiffString())
}

func TestMergeEnviron(t *testing.T) {
	merged := MergeEnviron(
		[]string{"BASE=base", "OVERRIDE=default"},
		[]string{"CUSTOM=custom", "OVERRIDE=override", "MALFORMED"},
	)

	env := New(merged)
	require.Equal(t, "base", env.Get("BASE"))
	require.Equal(t, "custom", env.Get("CUSTOM"))
	require.Equal(t, "override", env.Get("OVERRIDE"))
}
