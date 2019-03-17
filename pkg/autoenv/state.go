package autoenv

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

const autoEnvStateVariableName = "__BUD_AUTOENV"

type savedEnv map[string]*string

type AutoEnvState struct {
	ProjectHash string     `json:"project"`
	Features    FeatureSet `json:"features"`
	SavedEnv    savedEnv   `json:"saved_state"`
}

func (a *AutoEnvState) resetEnv() {
	a.SavedEnv = savedEnv{}
}

// FeatureState remember the current state of the features (whether they are active)
type FeatureState struct {
	env *env.Env
	UI  *termui.UI
}

func (s *FeatureState) read() *AutoEnvState {
	state := &AutoEnvState{}
	state.resetEnv()

	if s.env.Has(autoEnvStateVariableName) {
		err := json.Unmarshal([]byte(s.env.Get(autoEnvStateVariableName)), &state)
		if err != nil {
			panic(fmt.Sprintf("failed to read the state: %s", err))
		}
	}
	return state
}

func (s *FeatureState) write(state *AutoEnvState) {
	serialized, err := json.Marshal(state)
	if err != nil {
		panic(fmt.Sprintf("failed to write the state: %s", err))
	}
	s.env.Set(autoEnvStateVariableName, string(serialized))
}

// GetActiveFeatures returns a Hash of feature name -> param
func (s *FeatureState) GetActiveFeatures() FeatureSet {
	state := s.read()
	return state.Features
}

// SetProject change the state to set the project
func (s *FeatureState) SetProjectSlug(slug string) {
	state := s.read()
	state.ProjectHash = slug
	s.write(state)
}

// GetProjectSlug returns the slug of the project in which DevBuddy was when the state was written
func (s *FeatureState) GetProjectSlug() string {
	state := s.read()
	return state.ProjectHash
}

// SetFeature marks a feature as active
func (s *FeatureState) SetFeature(featureInfo *FeatureInfo) {
	state := s.read()
	state.Features = state.Features.With(featureInfo)
	s.write(state)
}

// UnsetFeature marks a feature as inactive
func (s *FeatureState) UnsetFeature(name string) {
	state := s.read()
	state.Features = state.Features.Without(name)
	s.write(state)
}

func (p savedEnv) String() string {
	elements := []string{}

	for name, value := range p {
		if value == nil {
			elements = append(elements, name+"=nil")
		} else {
			elements = append(elements, fmt.Sprintf("%s=\"%s\"", name, *value))
		}
	}
	return strings.Join(elements, " ")
}

func (s *FeatureState) ForgetEnv() {
	state := s.read()
	state.resetEnv()
	s.write(state)
}

func (s *FeatureState) SaveEnv() {
	state := s.read()

	for _, mutation := range s.env.Mutations() {
		if mutation.Name == autoEnvStateVariableName {
			continue // skip our own variable
		}

		if _, present := state.SavedEnv[mutation.Name]; present {
			continue // skip if we already recorded the initial value for this variable
		}

		if mutation.Previous == nil {
			state.SavedEnv[mutation.Name] = nil
		} else {
			copiedValue := fmt.Sprint(mutation.Previous.Value) // trick to make a copy of the string
			state.SavedEnv[mutation.Name] = &copiedValue
		}
	}

	s.write(state)
}

func (s *FeatureState) RestoreEnv() {
	state := s.read()

	for name, value := range state.SavedEnv {
		if value == nil {
			s.env.Unset(name)
			s.UI.Debug("restoring %s to deleted", name)
		} else {
			s.env.Set(name, *value)
			s.UI.Debug("restoring %s to %s", name, *value)
		}
	}
}
