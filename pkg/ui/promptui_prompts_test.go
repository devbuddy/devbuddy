package ui

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPromptUISelectTemplatesAlignInactiveRowsWithPandaMarker(t *testing.T) {
	templates := promptUISelectTemplates()

	require.Equal(t, "🐼 {{ .Label | cyan }}", templates.Active)
	require.Equal(t, "   {{ .Label }}", templates.Inactive)
}
