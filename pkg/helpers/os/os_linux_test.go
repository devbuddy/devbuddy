package os

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOSGetVariant(t *testing.T) {
	os := OS{"linux", "debian"}

	variant, err := os.GetVariant()

	require.NoError(t, err, "GetVariant() failed")
	require.Equal(t, "debian", variant)
}
