package autoenv

import (
	"encoding/json"
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/env"
)

const autoEnvVariableName = "__BUD_AUTOENV"

type StateData struct {
	ProjectSlug string     `json:"project"`
	Features    FeatureSet `json:"features"`
}

// StateManager remember the current state of the features (whether they are active)
type StateManager struct {
	env *env.Env
}

func (s *StateManager) read() *StateData {
	state := &StateData{}

	if s.env.Has(autoEnvVariableName) {
		err := json.Unmarshal([]byte(s.env.Get(autoEnvVariableName)), &state)
		if err != nil {
			panic(fmt.Sprintf("failed to read the state: %s", err))
		}
	}
	return state
}

func (s *StateManager) write(state *StateData) {
	serialized, err := json.Marshal(state)
	if err != nil {
		panic(fmt.Sprintf("failed to write the state: %s", err))
	}
	s.env.Set(autoEnvVariableName, string(serialized))
}

// GetActiveFeatures returns the FeatureSet recorded in the state
func (s *StateManager) GetActiveFeatures() FeatureSet {
	state := s.read()
	return state.Features
}

// SetProjectSlug records the project slug in the state
func (s *StateManager) SetProjectSlug(slug string) {
	state := s.read()
	state.ProjectSlug = slug
	s.write(state)
}

// GetProjectSlug returns the slug of the project in which DevBuddy was when the state was written
func (s *StateManager) GetProjectSlug() string {
	state := s.read()
	return state.ProjectSlug
}

// SetFeature marks a feature as active
func (s *StateManager) SetFeature(featureInfo *FeatureInfo) {
	state := s.read()
	state.Features = state.Features.With(featureInfo)
	s.write(state)
}

// UnsetFeature marks a feature as inactive
func (s *StateManager) UnsetFeature(name string) {
	state := s.read()
	state.Features = state.Features.Without(name)
	s.write(state)
}
