package projectmetadata_test

import (
	"path/filepath"
	"testing"

	"github.com/Flaque/filet"
	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/pkg/helpers/projectmetadata"
)

func Test(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")

	path, err := projectmetadata.New(tmpdir).Path()
	require.NoError(t, err)
	require.Equal(t, filepath.Join(tmpdir, ".devbuddy"), path)
	require.DirExists(t, path)

	gitignorePath := filepath.Join(path, ".gitignore")
	require.FileExists(t, gitignorePath)
	require.True(t, filet.FileSays(t, gitignorePath, []byte("*")))
}
