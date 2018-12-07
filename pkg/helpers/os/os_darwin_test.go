package os

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOSGetVariant(t *testing.T) {
	os := OS{"darwin", "18.0.2"}

	variant, err := os.GetVariant()

	require.NoError(t, err, "GetVariant() failed")
	require.Equal(t, "mojave", variant)

	os = OS{"darwin", "6.0.2"}

	variant, err = os.GetVariant()

	require.NoError(t, err, "GetVariant() failed")
	require.Equal(t, "jaguar", variant)

}

func TestOSInvalidGetVariant(t *testing.T) {
	os := OS{"darwin", "1.0.2"}

	variant, err := os.GetVariant()

	require.Error(t, err, fmt.Errorf("Cannot identify variant '1' for darwin"))
	require.Equal(t, "", variant)
}
