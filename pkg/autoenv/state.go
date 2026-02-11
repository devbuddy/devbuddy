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
			elements = append(elements, fmt.Sprintf("%s=%q", name, *value))
		}
	}
	return strings.Join(elements, " ")
}

type StateData struct {
	ProjectSlug   string            `json:"project"`
	Features      FeatureSet        `json:"features"`
	SavedEnv      savedEnv          `json:"saved_env"`
	FileChecksums map[string]string `json:"file_checksums,omitempty"`
}

// StateManager remember the current state of the features (whether they are active)
type StateManager struct {
	env *env.Env
	UI  *termui.UI
}

func (s *StateManager) read() (*StateData, error) {
	state := &StateData{SavedEnv: savedEnv{}}

	if s.env.Has(autoEnvVariableName) {
		err := json.Unmarshal([]byte(s.env.Get(autoEnvVariableName)), &state)
		if err != nil {
			return nil, fmt.Errorf("failed to read the state: %s", err)
		}
	}
	return state, nil
}

func (s *StateManager) write(state *StateData) error {
	serialized, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to write the state: %s", err)
	}
	s.env.Set(autoEnvVariableName, string(serialized))
	return nil
}

// GetActiveFeatures returns the FeatureSet recorded in the state
func (s *StateManager) GetActiveFeatures() (FeatureSet, error) {
	state, err := s.read()
	if err != nil {
		return nil, err
	}
	return state.Features, nil
}

// SetProjectSlug records the project slug in the state
func (s *StateManager) SetProjectSlug(slug string) error {
	state, err := s.read()
	if err != nil {
		return err
	}
	state.ProjectSlug = slug
	return s.write(state)
}

// GetProjectSlug returns the slug of the project in which DevBuddy was when the state was written
func (s *StateManager) GetProjectSlug() (string, error) {
	state, err := s.read()
	if err != nil {
		return "", err
	}
	return state.ProjectSlug, nil
}

// SetFeature marks a feature as active
func (s *StateManager) SetFeature(featureInfo *FeatureInfo) error {
	state, err := s.read()
	if err != nil {
		return err
	}
	state.Features = state.Features.With(featureInfo)
	return s.write(state)
}

// UnsetFeature marks a feature as inactive
func (s *StateManager) UnsetFeature(name string) error {
	state, err := s.read()
	if err != nil {
		return err
	}
	state.Features = state.Features.Without(name)
	return s.write(state)
}

// SaveEnv records the environment mutations in the state
func (s *StateManager) SaveEnv() error {
	state, err := s.read()
	if err != nil {
		return err
	}

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

	return s.write(state)
}

// RestoreEnv reverts the environment as recorded in the state
func (s *StateManager) RestoreEnv() error {
	state, err := s.read()
	if err != nil {
		return err
	}

	for name, value := range state.SavedEnv {
		if value == nil {
			s.env.Unset(name)
			s.UI.Debug("restoring %s to deleted", name)
		} else {
			s.env.Set(name, *value)
			s.UI.Debug("restoring %s to %s", name, *value)
		}
	}
	return nil
}

// ForgetEnv clears the environment mutations previously recorded in the state
func (s *StateManager) ForgetEnv() error {
	state, err := s.read()
	if err != nil {
		return err
	}
	state.SavedEnv = savedEnv{}
	return s.write(state)
}

// GetFileChecksums returns the file checksums recorded in the state
func (s *StateManager) GetFileChecksums() (map[string]string, error) {
	state, err := s.read()
	if err != nil {
		return nil, err
	}
	return state.FileChecksums, nil
}

// SetFileChecksums records file checksums in the state
func (s *StateManager) SetFileChecksums(checksums map[string]string) error {
	state, err := s.read()
	if err != nil {
		return err
	}
	state.FileChecksums = checksums
	return s.write(state)
}
