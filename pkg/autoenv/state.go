package autoenv

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

const autoEnvVariableName = "__BUD_AUTOENV"

type savedEnv map[string]*string

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

type StateData struct {
	ProjectSlug string     `json:"project"`
	Features    FeatureSet `json:"features"`
	SavedEnv    savedEnv   `json:"saved_env"`
}

// StateManager remember the current state of the features (whether they are active)
type StateManager struct {
	env *env.Env
	UI  *termui.UI
}

func (s *StateManager) read() *StateData {
	state := &StateData{SavedEnv: savedEnv{}}

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
	return s.read().ProjectSlug
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

// SaveEnv records the environment mutations in the state
func (s *StateManager) SaveEnv() {
	state := s.read()

	for _, mutation := range s.env.Mutations() {
		if mutation.Name == autoEnvVariableName {
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

// RestoreEnv reverts the environment as recorded in the state
func (s *StateManager) RestoreEnv() {
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

// ForgetEnv clears the environment mutations previously recorded in the state
func (s *StateManager) ForgetEnv() {
	state := s.read()
	state.SavedEnv = savedEnv{}
	s.write(state)
}
