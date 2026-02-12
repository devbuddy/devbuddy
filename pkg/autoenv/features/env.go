package features

import (
	"encoding/json"
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/context"
)

func init() {
	register.Register(env{})
}

type env struct{}

func (env) Name() string {
	return "env"
}

// Activate sets the env vars encoded as JSON in param.
func (env) Activate(ctx *context.Context, param string) (bool, error) {
	var envVars map[string]string
	if err := json.Unmarshal([]byte(param), &envVars); err != nil {
		return false, fmt.Errorf("env feature: invalid param: %w", err)
	}
	for name, value := range envVars {
		ctx.Env.Set(name, value)
	}
	return false, nil
}

func (env) Deactivate(ctx *context.Context, param string) {}
