package features

import (
	"testing"

	"github.com/devbuddy/devbuddy/pkg/features/definitions"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	def := definitions.Get("python")
	require.NotNil(t, def)

	def = definitions.Get("golang")
	require.NotNil(t, def)

	names := definitions.Names()
	require.ElementsMatch(t, []string{"python", "golang"}, names)
}
