package osidentity

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetVariant(t *testing.T) {
	identity := DetectFromReleaseId("debian")

	variant, err := identity.GetVariant()

	require.NoError(t, err, "GetVariant() failed")
	require.Equal(t, "debian", variant)
}

func TestInvalidGetVariant(t *testing.T) {
	identity := DetectFromReleaseId("invalid")

	variant, err := identity.GetVariant()

	require.Error(t, err, fmt.Errorf("Cannot identify variant 'invalid' for linux"))
	require.Equal(t, "", variant)
}
