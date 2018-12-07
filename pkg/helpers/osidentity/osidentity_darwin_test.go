package osidentity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetectIsMacOS(t *testing.T) {
	i := Detect()
	require.True(t, i.IsMacOS())
}

func TestDetectIsDebianLike(t *testing.T) {
	i := Detect()
	require.False(t, i.IsDebianLike())
}
