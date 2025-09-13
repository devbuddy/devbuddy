package manifest

import (
	"errors"
	"os"
	"path"

	"github.com/devbuddy/devbuddy/pkg/utils"
)

// Create writes a default project manifest in the specified path
func Create(projectPath string, templateName string) error {
	template, err := LoadTemplate(templateName)
	if err != nil {
		return err
	}

	manifestPath := path.Join(projectPath, manifestFilename)
	err = utils.WriteNewFile(manifestPath, template, 0666)
	if os.IsExist(err) {
		return errors.New("the manifest already exists")
	}
	return err
}
