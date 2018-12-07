package osidentity

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetVariant(t *testing.T) {
	identity := Identity{"darwin", "18.0.2"}

	variant, err := identity.GetVariant()

	require.NoError(t, err, "GetVariant() failed")
	require.Equal(t, "mojave", variant)

	identity = Identity{"darwin", "6.0.2"}

	variant, err = identity.GetVariant()

	require.NoError(t, err, "GetVariant() failed")
	require.Equal(t, "jaguar", variant)

}

func TestInvalidGetVariant(t *testing.T) {
	identity := Identity{"darwin", "1.0.2"}

	variant, err := identity.GetVariant()

	require.Error(t, err, fmt.Errorf("Cannot identify variant '1' for darwin"))
	require.Equal(t, "", variant)
}
