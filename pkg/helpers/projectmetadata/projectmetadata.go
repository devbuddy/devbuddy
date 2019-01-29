package projectmetadata

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/devbuddy/devbuddy/pkg/utils"
)

var dirName = ".devbuddy"

// ProjectMetadata is the place to store metadata files about a project
type ProjectMetadata struct {
	projectPath string
}

// New returns an instance of ProjectMetadata
func New(projectPath string) *ProjectMetadata {
	return &ProjectMetadata{projectPath: projectPath}
}

func (p *ProjectMetadata) Path() string {
	return filepath.Join(p.projectPath, dirName)
}

// Prepare makes sure the metadata directory is ready
func (p *ProjectMetadata) Prepare() (err error) {
	if !utils.PathExists(p.projectPath) {
		return fmt.Errorf("failed to initialize the project metadata dir: path does not exist: %s", p.projectPath)
	}

	if !utils.PathExists(p.Path()) {
		err = os.MkdirAll(p.Path(), 0755)
		if err != nil {
			return
		}
	}

	gitignore := filepath.Join(p.Path(), ".gitignore")
	if !utils.PathExists(gitignore) {
		err = ioutil.WriteFile(gitignore, []byte("*"), 0644)
		if err != nil {
			return
		}
	}
	return nil
}
