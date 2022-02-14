package debug

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/pkg/test"
)

func TestFormatDebugInfo(t *testing.T) {
	text := FormatDebugInfo("versionONE", []string{"SHELL=/bin/bash"}, "")
	require.Contains(t, text, "`versionONE`")
	require.Contains(t, text, "SHELL=\"/bin/bash\"\n")
	require.Contains(t, text, "Project not found")

	tmpdir := t.TempDir()

	text = FormatDebugInfo("", []string{}, tmpdir)
	require.Contains(t, text, "Failed to read manifest: no manifest at")

	writer := test.Project(tmpdir)
	writer.Manifest().WriteString("up: [{go: 1.2.3}]")

	text = FormatDebugInfo("", []string{}, tmpdir)
	require.Contains(t, text, "0. `map[go:1.2.3]`")
}

func TestNewGithubIssueURL(t *testing.T) {
	url := NewGithubIssueURL("", []string{}, "")
	require.Contains(t, url, "https://github.com/devbuddy/devbuddy/issues/new?")
	require.Contains(t, url, "Project+not+found")
}
