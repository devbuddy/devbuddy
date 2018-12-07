package os

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOSGetVariant(t *testing.T) {
	os := OS{"linux", "debian"}

	variant, err := os.GetVariant()

	require.NoError(t, err, "GetVariant() failed")
	require.Equal(t, "debian", variant)
}

func TestOSInvalidGetVariant(t *testing.T) {
	os := OS{"linux", "invalid"}

	variant, err := os.GetVariant()

	require.Error(t, err, fmt.Errorf("Cannot identify variant 'invalid' for linux"))
	require.Equal(t, "", variant)
}
