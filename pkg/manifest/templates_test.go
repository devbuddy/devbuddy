package manifest

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTemplatesPointToBudDoc(t *testing.T) {
	for _, name := range ListTemplates() {
		t.Run(name, func(t *testing.T) {
			template, err := LoadTemplate(name)
			require.NoError(t, err)

			lines := strings.Split(string(template), "\n")
			require.GreaterOrEqual(t, len(lines), 3)
			require.Contains(t, lines[0], "DevBuddy config file")
			require.Contains(t, lines[1], "bud --doc")
		})
	}
}
