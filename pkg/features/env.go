package features

import (
	"fmt"
	"strings"
)

type Env struct {
	env         map[string]string
	verbatimEnv map[string]string
}

func NewEnv(env []string) (e *Env) {
	e = &Env{
		env:         make(map[string]string),
		verbatimEnv: make(map[string]string),
	}

	for _, pair := range env {
		parts := strings.SplitN(pair, "=", 2)
		e.env[parts[0]] = parts[1]
		e.verbatimEnv[parts[0]] = parts[1]
	}

	return
}

type EnvVarChange struct {
	Name    string
	Value   string
	Deleted bool
}

func (e *Env) Set(name string, value string) {
	e.env[name] = value
}

func (e *Env) Unset(name string) {
	delete(e.env, name)
}

func (e *Env) Get(name string) string {
	return e.env[name]
}

func (e *Env) Changed() []EnvVarChange {
	c := []EnvVarChange{}

	for k, v := range e.env {
		if v != e.verbatimEnv[k] {
			c = append(c, EnvVarChange{Name: k, Value: v})
		}
	}
	for k := range e.verbatimEnv {
		if _, ok := e.env[k]; !ok {
			c = append(c, EnvVarChange{Name: k, Deleted: true})
		}
	}
	return c
}

const AutoEnvVariableName = "DEV_AUTO_ENV_FEATURES"

func (e *Env) GetActiveFeatures() map[string]string {
	features := map[string]string{}

	val := e.env[AutoEnvVariableName]
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

func (e *Env) SetActiveFeatures(features map[string]string) {
	var parts []string

	for feat, param := range features {
		parts = append(parts, fmt.Sprintf("%s=%s", feat, param))
	}

	val := strings.Join(parts, ":")

	if len(val) == 0 {
		delete(e.env, AutoEnvVariableName)
	} else {
		e.env[AutoEnvVariableName] = val
	}
}
