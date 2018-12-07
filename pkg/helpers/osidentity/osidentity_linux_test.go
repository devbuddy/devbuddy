package osidentity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetectIsDebianLike(t *testing.T) {
	i := Detect()
	require.True(t, i.IsDebianLike())
}

func TestDetectIsMacOS(t *testing.T) {
	i := Detect()
	require.False(t, i.IsMacOS())
}
