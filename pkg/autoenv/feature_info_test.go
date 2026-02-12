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

func TestFeatureInfoDisplayString(t *testing.T) {
	tests := []struct {
		name  string
		param string
		want  string
	}{
		{"python", "3.6.5", "python 3.6.5"},
		{"golang", "1.13.1+mod", "golang 1.13.1+mod"},
		{"env", `{"VAR":"val"}`, "env"},
		{"env", "", "env"},
	}
	for _, tt := range tests {
		t.Run(tt.name+"_"+tt.param, func(t *testing.T) {
			fi := autoenv.NewFeatureInfo(tt.name, tt.param)
			require.Equal(t, tt.want, fi.DisplayString())
		})
	}
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
