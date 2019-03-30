package helpers

import (
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/joho/godotenv"
)

// LoadEnvfile sets all the environment variables defined in an envfile
func LoadEnvfile(env *env.Env, path string) error {
	loadedEnv, err := godotenv.Read(path)
	if err != nil {
		return err
	}

	for name, value := range loadedEnv {
		env.Set(name, value)
	}

	return nil
}
