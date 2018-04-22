package env

import (
	"fmt"
	"strings"
)

const autoEnvVariableName = "DAD_AUTO_ENV_FEATURES"

func (e *Env) GetActiveFeatures() map[string]string {
	features := map[string]string{}

	val := e.env[autoEnvVariableName]
	if val != "" {
		for _, feat := range strings.Split(val, ":") {
			if feat != "" {
				parts := strings.SplitN(feat, "=", 2)
				if len(parts) == 2 {
					features[parts[0]] = parts[1]
				}
			}
		}
	}

	return features
}

func (e *Env) setActiveFeatures(features map[string]string) {
	var parts []string

	for feat, param := range features {
		parts = append(parts, fmt.Sprintf("%s=%s", feat, param))
	}

	val := strings.Join(parts, ":")

	if len(val) == 0 {
		delete(e.env, autoEnvVariableName)
	} else {
		e.env[autoEnvVariableName] = val
	}
}

func (e *Env) SetFeature(name, param string) {
	features := e.GetActiveFeatures()
	features[name] = param
	e.setActiveFeatures(features)
}

func (e *Env) UnsetFeature(name string) {
	features := e.GetActiveFeatures()
	delete(features, name)
	e.setActiveFeatures(features)
}
