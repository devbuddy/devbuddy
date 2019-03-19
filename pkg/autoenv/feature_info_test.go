package autoenv_test

import (
	"testing"

	"github.com/devbuddy/devbuddy/pkg/autoenv"
	"github.com/stretchr/testify/require"
)

func TestFeatureInfo(t *testing.T) {
	fi := autoenv.NewFeatureInfo("NAME", "PARAM")
	require.Equal(t, "NAME", fi.Name)
	require.Equal(t, "PARAM", fi.Param)
	require.Equal(t, "NAME:PARAM", fi.String())
}

func TestFeatureSet(t *testing.T) {
	fs := autoenv.NewFeatureSet()
	require.Nil(t, fs.Get("python"))
	require.Equal(t, "", fs.String())

	fs = fs.With(autoenv.NewFeatureInfo("python", "3.6"))
	require.Equal(t, "3.6", fs.Get("python").Param)
	require.Equal(t, "python:3.6", fs.String())

	fs = fs.With(autoenv.NewFeatureInfo("go", "1.12"))
	require.Equal(t, "python:3.6 go:1.12", fs.String())

	fs = fs.Without("python")
	require.Nil(t, fs.Get("python"))
	require.Equal(t, "go:1.12", fs.String())
}
