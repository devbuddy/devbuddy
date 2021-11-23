package projectmetadata_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/pkg/helpers/projectmetadata"
	"github.com/devbuddy/devbuddy/pkg/test"
)

func Test(t *testing.T) {
	tmpdir, gitignorePath := test.File(t, ".devbuddy/.gitignore")

	path, err := projectmetadata.New(tmpdir).Path()
	require.NoError(t, err)

	require.Equal(t, filepath.Join(tmpdir, ".devbuddy"), path)
	require.DirExists(t, path)

	data := test.ReadFile(gitignorePath)
	require.Equal(t, data, []byte("*"))
}
