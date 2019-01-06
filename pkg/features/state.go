package features

import (
	"fmt"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/env"
)

const autoEnvVariableName = "BUD_AUTO_ENV_FEATURES"

// FeatureState remember the current state of the features (whether they are active)
type FeatureState struct {
	env *env.Env
}

func serializeFeatureSet(featureSet FeatureSet) string {
	var parts []string
	for _, info := range featureSet {
		parts = append(parts, fmt.Sprintf("%s=%s", info.Name, info.Param))
	}
	return strings.Join(parts, ":")
}

func deserializeFeatureSet(serialized string) FeatureSet {
	features := FeatureSet{}
	for _, feat := range strings.Split(serialized, ":") {
		if feat != "" {
			parts := strings.SplitN(feat, "=", 2)
			if len(parts) == 2 {
				features = features.With(FeatureInfo{parts[0], parts[1]})
			}
		}
	}
	return features
}

// GetActiveFeatures returns a Hash of feature name -> param
func (s *FeatureState) GetActiveFeatures() FeatureSet {
	return deserializeFeatureSet(s.env.Get(autoEnvVariableName))
}

func (s *FeatureState) setActiveFeatures(featureSet FeatureSet) {
	val := serializeFeatureSet(featureSet)
	if len(val) == 0 {
		s.env.Unset(autoEnvVariableName)
	} else {
		s.env.Set(autoEnvVariableName, val)
	}
}

// SetFeature marks a feature as active
func (s *FeatureState) SetFeature(featureInfo FeatureInfo) {
	s.setActiveFeatures(s.GetActiveFeatures().With(featureInfo))
}

// UnsetFeature marks a feature as inactive
func (s *FeatureState) UnsetFeature(name string) {
	s.setActiveFeatures(s.GetActiveFeatures().Without(name))
}
