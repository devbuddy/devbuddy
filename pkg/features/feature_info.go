package features

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
