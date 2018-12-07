package osidentity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetectIsDebianLike(t *testing.T) {
	i := Detect()
	require.True(t, i.IsDebianLike())
	require.False(t, i.IsMacOS())
}
