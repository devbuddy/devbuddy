package manifest

import (
	"path"

	"github.com/devbuddy/devbuddy/pkg/utils"
)

const defaultManifestContent = `# Created by "bud init"

env:
  ENV: development

up:
  - go: 1.10.1
  - golang_dep
  - python: 3.6.5
  - pip: [tests/requirements.txt]
  - homebrew:
    - curl
  - custom:
      name: Install gometalinter
      met?: which gometalinter.v2 > /dev/null
      meet: go get gopkg.in/alecthomas/gometalinter.v2

commands:
  test:
    desc: Run the unittests
    run: script/test

  lint:
    desc: Lint the project
    run: script/lint

  release:
    desc: Create a new release (bud release [VERSION])
    run: script/release

  godoc:
    desc: Starting GoDoc server on http://0.0.0.0:6060
    run: (sleep 1; open http://0.0.0.0:6060)& godoc -http=:6060

open:
  devbuddy: https://github.com/devbuddy/devbuddy/blob/master/docs/Config.md#config-devyml
`

// Create writes a default project manifest in the specified path
func Create(projectPath string) error {
	manifestPath := path.Join(projectPath, manifestFilename)
	return utils.WriteNewFile(manifestPath, []byte(defaultManifestContent), 0666)
}
