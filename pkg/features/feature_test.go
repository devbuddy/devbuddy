package features

import (
	"testing"

	"github.com/devbuddy/devbuddy/pkg/features/definitions"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	names := definitions.Names()
	require.ElementsMatch(t, []string{"python", "golang"}, names)

	for _, name := range names {
		require.NotNil(t, definitions.Get(name))
	}
}
