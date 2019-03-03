package autoenv

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

// GetActiveFeatures returns a Hash of feature name -> param
func (s *FeatureState) GetActiveFeatures() FeatureSet {
	_, set := s.get()
	return set
}

// SetProject change the state to set the project
func (s *FeatureState) SetProjectSlug(slug string) {
	_, set := s.get()
	s.set(slug, set)
}

// IsProject returns whether the state is for this project
func (s *FeatureState) GetProjectSlug() string {
	slug, _ := s.get()
	return slug
}

func (s *FeatureState) get() (string, FeatureSet) {
	data := s.env.Get(autoEnvVariableName)

	if strings.HasPrefix(data, "1:") {
		parts := strings.SplitN(data, ":", 3)
		return parts[1], NewFeatureSetFromString(parts[2])
	}
	return "", NewFeatureSetFromString(data)
}

func (s *FeatureState) set(slug string, featureSet FeatureSet) {
	val := fmt.Sprintf("1:%s:%s", slug, featureSet.Serialize())
	if len(val) == 0 {
		s.env.Unset(autoEnvVariableName)
	} else {
		s.env.Set(autoEnvVariableName, val)
	}
}

// SetFeature marks a feature as active
func (s *FeatureState) SetFeature(featureInfo FeatureInfo) {
	pKey, set := s.get()
	s.set(pKey, set.With(featureInfo))
}

// UnsetFeature marks a feature as inactive
func (s *FeatureState) UnsetFeature(name string) {
	pKey, set := s.get()
	s.set(pKey, set.Without(name))
}
