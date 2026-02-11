package context

import (
	"fmt"
	"os"
	"path"
)

type Config struct {
	ShellName   string
	BinaryPath  string
	DockerImage string
}

func LoadConfig() (Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return Config{}, fmt.Errorf("getting current working directory: %w", err)
	}

	dockerImage, ok := os.LookupEnv("TEST_DOCKER_IMAGE")
	if !ok {
		return Config{}, fmt.Errorf("missing env var TEST_DOCKER_IMAGE")
	}

	shellName := "bash"
	if v, ok := os.LookupEnv("TEST_SHELL"); ok {
		shellName = v
	}

	return Config{
		ShellName:   shellName,
		BinaryPath:  path.Join(cwd, "bud"),
		DockerImage: dockerImage,
	}, nil
}
