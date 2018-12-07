package osidentity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsDebianLike(t *testing.T) {
	i := Identity{"linux", "debian"}

	require.True(t, i.IsDebianLike())

	i = Identity{"other", "debian"}

	require.False(t, i.IsDebianLike())

	i = Identity{"other", "redhat"}

	require.False(t, i.IsDebianLike())
}

func TestIsMacOS(t *testing.T) {
	i := Identity{"darwin", ""}

	require.False(t, i.IsDebianLike())

	require.True(t, i.IsMacOS())

	i = Identity{"linux", "debian"}

	require.False(t, i.IsMacOS())
}
