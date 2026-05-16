package harness

import (
	"fmt"
	"os"
)

type Config struct {
	ShellName              string
	BinaryPath             string // set by the caller (e.g. via budbuild.TempBinaryPath)
	DockerImage            string
	WorkspaceHostPath      string
	WorkspaceContainerPath string
	UsePTY                 bool
}

func LoadConfig() (Config, error) {
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
		DockerImage: dockerImage,
	}, nil
}
