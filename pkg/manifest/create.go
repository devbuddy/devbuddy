package manifest

import (
	"path"

	"github.com/devbuddy/devbuddy/pkg/utils"
)

const defaultManifestContent = `# DevBuddy config file
# See https://github.com/devbuddy/devbuddy/blob/master/docs/Config.md

up:
  # MacOS:
  # - homebrew:
  #   - curl
  #   - golangci/tap/golangci-lint

  # Linux:
  # - apt:
  #   - python3-dev

  # Go:
  # - go:
  #     version: 1.12.4
  #     modules: true
  # - golang_dep

  # Python:
  # - python: 3.7.0
  # - pipfile
  # - pip:
  #   - requirements.txt
  #   - tests/requirements.txt
  # - python_develop

  # Custom task:
  # - custom:
  #     name: Create the local config file
  #     met?: test -f config/local.yml
  #     meet: cp config/local.yml.tmpl config/local.yml

  - custom:
      name: Edit dev.yml then remove me
      met?: 'false'
      meet: $EDITOR dev.yml

commands:
  test:
    desc: Run all tests
    run: go test ./... -cover

open:
  devbuddy: https://github.com/devbuddy/devbuddy/blob/master/docs/Config.md#config-devyml
`

// Create writes a default project manifest in the specified path
func Create(projectPath string) error {
	manifestPath := path.Join(projectPath, manifestFilename)
	return utils.WriteNewFile(manifestPath, []byte(defaultManifestContent), 0666)
}
