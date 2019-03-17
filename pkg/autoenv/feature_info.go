package autoenv

import (
	"fmt"
	"strings"
)

// FeatureInfo represents a parameterized feature
type FeatureInfo struct {
	Name  string `json:"name"`
	Param string `json:"param"`
}

// NewFeatureInfo returns a FeatureInfo
func NewFeatureInfo(name string, param string) *FeatureInfo {
	return &FeatureInfo{name, param}
}

// FeatureSet represents a set of parameterized features
type FeatureSet []*FeatureInfo

// NewFeatureSet returns a FeatureSet
func NewFeatureSet() FeatureSet {
	return FeatureSet{}
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
		elements = append(elements, fmt.Sprintf("%s:%s", element.Name, element.Param))
	}
	return strings.Join(elements, " ")
}
