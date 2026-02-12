package autoenv

import (
	"encoding/json"
	"time"

	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/integration"
)

const featureCacheVariableName = "__BUD_FEATURE_CACHE"

const devYmlChangedWarningInterval = 60 * time.Second

// FeatureCache stores the parsed feature set for a project so we can skip
// re-parsing dev.yml on every shell prompt.
type FeatureCache struct {
	ProjectSlug string     `json:"project_slug"`
	Checksum    string     `json:"checksum"`
	Features    FeatureSet `json:"features"`
	WarnedAt    int64      `json:"warned_at,omitempty"` // unix timestamp of last dev.yml-changed warning
}

func NewFeatureCache(slug, checksum string, features FeatureSet) *FeatureCache {
	return &FeatureCache{
		ProjectSlug: slug,
		Checksum:    checksum,
		Features:    features,
	}
}

// ShouldWarnDevYmlChanged returns true if enough time has passed since the
// last warning. This avoids spamming the user on every prompt.
func (c *FeatureCache) ShouldWarnDevYmlChanged() bool {
	if c.WarnedAt == 0 {
		return true
	}
	return time.Since(time.Unix(c.WarnedAt, 0)) >= devYmlChangedWarningInterval
}

// MarkWarned records the current time as the last warning timestamp.
func (c *FeatureCache) MarkWarned() {
	c.WarnedAt = time.Now().Unix()
}

// ReadFeatureCache reads the cached feature set from the env var.
// Returns nil if the env var is not set or contains invalid JSON.
func ReadFeatureCache(e *env.Env) *FeatureCache {
	if !e.Has(featureCacheVariableName) {
		return nil
	}
	var cache FeatureCache
	if err := json.Unmarshal([]byte(e.Get(featureCacheVariableName)), &cache); err != nil {
		return nil
	}
	return &cache
}

// WriteFeatureCache writes the feature cache into the env var so it will be
// exported by the hook's shell output.
func WriteFeatureCache(e *env.Env, cache *FeatureCache) {
	data, err := json.Marshal(cache)
	if err != nil {
		return
	}
	e.Set(featureCacheVariableName, string(data))
}

// WriteFeatureCacheFinalizer writes the feature cache via a setenv finalizer
// so the calling shell picks it up after `bud up` exits.
func WriteFeatureCacheFinalizer(cache *FeatureCache) error {
	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}
	return integration.AddFinalizerSetEnv(featureCacheVariableName, string(data))
}
