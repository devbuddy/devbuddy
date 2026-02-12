package autoenv

import (
	"testing"

	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/stretchr/testify/require"
)

func TestReadFeatureCacheNotSet(t *testing.T) {
	e := env.New([]string{})
	cache := ReadFeatureCache(e)
	require.Nil(t, cache)
}

func TestReadFeatureCacheCorruptJSON(t *testing.T) {
	e := env.New([]string{"__BUD_FEATURE_CACHE=not-valid-json"})
	cache := ReadFeatureCache(e)
	require.Nil(t, cache)
}

func TestFeatureCacheRoundTrip(t *testing.T) {
	e := env.New([]string{})

	features := NewFeatureSet().
		With(NewFeatureInfo("golang", "1.21")).
		With(NewFeatureInfo("env", `{"FOO":"bar"}`))

	original := NewFeatureCache("myproject-12345", "98765", features)
	WriteFeatureCache(e, original)

	got := ReadFeatureCache(e)
	require.NotNil(t, got)
	require.Equal(t, original.ProjectSlug, got.ProjectSlug)
	require.Equal(t, original.Checksum, got.Checksum)
	require.Equal(t, original.Features.String(), got.Features.String())
}

func TestFeatureCacheRoundTripEmptyFeatures(t *testing.T) {
	e := env.New([]string{})

	original := NewFeatureCache("proj-111", "55555", NewFeatureSet())
	WriteFeatureCache(e, original)

	got := ReadFeatureCache(e)
	require.NotNil(t, got)
	require.Equal(t, "proj-111", got.ProjectSlug)
	require.Equal(t, "55555", got.Checksum)
	require.Empty(t, got.Features)
}

func TestFeatureCacheOverwrite(t *testing.T) {
	e := env.New([]string{})

	first := NewFeatureCache("proj-a", "111", NewFeatureSet().With(NewFeatureInfo("rust", "1.0")))
	WriteFeatureCache(e, first)

	second := NewFeatureCache("proj-b", "222", NewFeatureSet().With(NewFeatureInfo("python", "3.9")))
	WriteFeatureCache(e, second)

	got := ReadFeatureCache(e)
	require.NotNil(t, got)
	require.Equal(t, "proj-b", got.ProjectSlug)
	require.Equal(t, "222", got.Checksum)
	require.Equal(t, "python:3.9", got.Features.String())
}
