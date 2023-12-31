package projectmetadata

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/devbuddy/devbuddy/pkg/utils"
)

var dirName = ".devbuddy"

// ProjectMetadata is the place to store metadata files about a project
type ProjectMetadata struct {
	path string
}

// New returns an instance of ProjectMetadata
func New(projectPath string) *ProjectMetadata {
	if !utils.PathExists(projectPath) {
		panic("project path does not exist: " + projectPath)
	}

	return &ProjectMetadata{
		path: filepath.Join(projectPath, dirName),
	}
}

// Path returns the path of the project metadata directory.
func (p *ProjectMetadata) Path() (string, error) {
	err := p.prepare()
	if err != nil {
		return "", fmt.Errorf("failed to initialize the project metadata dir: %w", err)
	}

	return p.path, nil
}

func (p *ProjectMetadata) prepare() error {
	if !utils.PathExists(p.path) {
		err := os.MkdirAll(p.path, 0755)
		if err != nil {
			return err
		}
	}

	gitignore := filepath.Join(p.path, ".gitignore")
	if !utils.PathExists(gitignore) {
		err := os.WriteFile(gitignore, []byte("*"), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
