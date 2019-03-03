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

// func (i FeatureInfo) String() string {
// 	return fmt.Sprintf("%s (%s)", i.Name, i.Param)
// }

// NewFeatureInfo returns a FeatureInfo
func NewFeatureInfo(name string, param string) FeatureInfo {
	return FeatureInfo{name, param}
}

// FeatureSet represents a set of parameterized features
type FeatureSet map[string]FeatureInfo

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
				set = set.With(FeatureInfo{parts[0], parts[1]})
			}
		}
	}
	return set
}

// With returns a new FeatureSet augmented with the featureInfo provided
func (s FeatureSet) With(featureInfo FeatureInfo) FeatureSet {
	s[featureInfo.Name] = featureInfo
	return s
}

// Without returns a new FeatureSet augmented with the featureInfo provided
func (s FeatureSet) Without(name string) FeatureSet {
	delete(s, name)
	return s
}

// Serialize returns the FeatureSet serialized as a string
func (s FeatureSet) Serialize() string {
	var parts []string
	for _, info := range s {
		parts = append(parts, fmt.Sprintf("%s=%s", info.Name, info.Param))
	}
	return strings.Join(parts, ":")
}
