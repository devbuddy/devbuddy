package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGOVersion(t *testing.T) {
	parse := func(s string) GOVersion {
		v, err := ParseGOVersion(s)
		require.NoError(t, err)
		return v
	}

	require.Equal(t, GOVersion{1, 5}, parse("1.5"))
	require.Equal(t, GOVersion{1, 5}, parse("1.5.0"))

	require.True(t, parse("0.10").LessThan(GOVersion{1, 6}))
	require.True(t, parse("1.5").LessThan(GOVersion{1, 6}))
	require.True(t, parse("1.5.6").LessThan(GOVersion{1, 6}))

	require.False(t, parse("1.6").LessThan(GOVersion{1, 6}))
	require.False(t, parse("1.7").LessThan(GOVersion{1, 6}))
	require.False(t, parse("2.0").LessThan(GOVersion{1, 6}))
}
