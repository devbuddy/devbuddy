package os

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOSGetVariant(t *testing.T) {
	os := OS{"darwin", "18.0.2"}

	variant, err := os.GetVariant()

	require.NoError(t, err, "GetVariant() failed")
	require.Equal(t, "mojave", variant)
}
