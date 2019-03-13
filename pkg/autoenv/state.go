package autoenv

import (
	"fmt"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/termui"

	"github.com/devbuddy/devbuddy/pkg/env"
)

const autoEnvVariableName = "BUD_AUTO_ENV_FEATURES"

// FeatureState remember the current state of the features (whether they are active)
type FeatureState struct {
	env *env.Env
	UI  *termui.UI
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

const autoEnvStateVariableName = "__BUD_AUTOENV_STATE"

type previousState map[string]*string

func (p previousState) String() string {
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

type AutoEnvState struct {
	Previous map[string]*string `json:previous`
}

func (s *FeatureState) loadState() (previousState, error) {
	state := make(previousState)

	if s.env.Has(autoEnvStateVariableName) {
		err := Unmarshal(s.env.Get(autoEnvStateVariableName), &state)
		if err != nil {
			return nil, err
		}
	}

	return state, nil
}

func (s *FeatureState) Forget() {
	s.env.Unset(autoEnvStateVariableName)
}

func (s *FeatureState) Save() error {
	state, err := s.loadState()
	if err != nil {
		return err
	}
	s.UI.Debug("Loaded state: %s", state)

	for _, mutation := range s.env.Mutations() {
		if mutation.Name == autoEnvStateVariableName || mutation.Name == autoEnvVariableName {
			continue
		}
		if _, present := state[mutation.Name]; present {
			continue // only the first mutation should be recorded to keep the original state
		}
		if mutation.Previous == nil {
			state[mutation.Name] = nil
		} else {
			copiedValue := fmt.Sprint(mutation.Previous.Value)
			state[mutation.Name] = &copiedValue
		}
	}

	serialized, err := Marshal(state)
	if err != nil {
		return err
	}

	s.env.Set(autoEnvStateVariableName, string(serialized))

	return nil
}

func (s *FeatureState) Restore() error {
	state, err := s.loadState()
	if err != nil {
		return err
	}

	for name, value := range state {
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
