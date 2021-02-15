package context

import (
	"os"
	"path"
)

type Config struct {
	ShellName   string
	BinaryPath  string
	DockerImage string
}

func LoadConfig() Config {
	return Config{
		ShellName:   env("TEST_SHELL", "bash"),
		BinaryPath:  path.Join(cwd(), "bud"),
		DockerImage: env("TEST_DOCKER_IMAGE_NAME", "devbuddy-test-env-linux"),
	}
}

func cwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic("Error when getting the current working directory: " + err.Error())
	}
	return cwd
}

func env(name string, defaultValue string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	return defaultValue
}
