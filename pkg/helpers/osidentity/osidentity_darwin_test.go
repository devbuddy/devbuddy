package osidentity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetectIsMacOS(t *testing.T) {
	i := Detect()
	require.True(t, i.IsMacOS())
	require.False(t, i.IsDebianLike())
}
