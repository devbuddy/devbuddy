package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/pkg/ui"
)

func TestCreateManifestPromptsForTemplate(t *testing.T) {
	projectPath := t.TempDir()
	_, cmdUI := ui.NewBufferedTesting(false)
	prompts := &ui.FakePrompts{SelectValue: "go"}
	cmdUI.SetPrompts(prompts)

	err := createManifest(cmdUI, projectPath, "")

	require.NoError(t, err)
	require.FileExists(t, filepath.Join(projectPath, "dev.yml"))
	require.Len(t, prompts.SelectRequests, 1)
	require.Equal(t, "Select a template", prompts.SelectRequests[0].Label)
	require.Contains(t, prompts.SelectRequests[0].Options, ui.SelectOption{Value: "go", Label: "go"})
}

func TestCreateManifestPromptCancellationSkipsManifest(t *testing.T) {
	projectPath := t.TempDir()
	_, cmdUI := ui.NewBufferedTesting(false)
	cmdUI.SetPrompts(&ui.FakePrompts{SelectErr: ui.ErrPromptCancelled})

	err := createManifest(cmdUI, projectPath, "")

	require.NoError(t, err)
	_, statErr := os.Stat(filepath.Join(projectPath, "dev.yml"))
	require.ErrorIs(t, statErr, os.ErrNotExist)
}
