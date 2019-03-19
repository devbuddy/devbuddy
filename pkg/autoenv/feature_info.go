package autoenv

import (
	"fmt"
	"strings"
)

// FeatureInfo represents a parameterized feature
type FeatureInfo struct {
	Name  string
	Param string
}

// NewFeatureInfo returns a FeatureInfo
func NewFeatureInfo(name string, param string) *FeatureInfo {
	return &FeatureInfo{name, param}
}

func (f FeatureInfo) String() string {
	return fmt.Sprintf("%s:%s", f.Name, f.Param)
}

// FeatureSet represents a set of parameterized features
type FeatureSet []*FeatureInfo

// NewFeatureSet returns a FeatureSet
func NewFeatureSet() FeatureSet {
	return FeatureSet{}
}

// NewFeatureSetFromString returns a new FeatureSet from a serialized string
func NewFeatureSetFromString(data string) FeatureSet {
	set := FeatureSet{}
	for _, feat := range strings.Split(data, ":") {
		if feat != "" {
			parts := strings.SplitN(feat, "=", 2)
			if len(parts) == 2 {
				set = set.With(NewFeatureInfo(parts[0], parts[1]))
			}
		}
	}
	return set
}

// With returns a new FeatureSet augmented with the featureInfo provided
func (s FeatureSet) With(featureInfo *FeatureInfo) FeatureSet {
	return append(s.Without(featureInfo.Name), featureInfo)
}

// Without returns a new FeatureSet augmented with the featureInfo provided
func (s FeatureSet) Without(name string) FeatureSet {
	newSet := FeatureSet{}
	for _, element := range s {
		if element.Name != name {
			newSet = append(newSet, element)
		}
	}
	return newSet
}

// Get returns a new FeatureSet augmented with the featureInfo provided
func (s FeatureSet) Get(name string) *FeatureInfo {
	for _, element := range s {
		if element.Name == name {
			return element
		}
	}
	return nil
}

func (s FeatureSet) String() string {
	elements := []string{}
	for _, element := range s {
		elements = append(elements, element.String())
	}
	return strings.Join(elements, " ")
}

// Serialize returns the FeatureSet serialized as a string
func (s FeatureSet) Serialize() string {
	var parts []string
	for _, info := range s {
		parts = append(parts, fmt.Sprintf("%s=%s", info.Name, info.Param))
	}
	return strings.Join(parts, ":")
}
