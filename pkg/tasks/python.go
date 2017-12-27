package tasks

import (
	"fmt"
	"strings"

	color "github.com/logrusorgru/aurora"

	"github.com/pior/dad/pkg/executor"
)

type Python struct {
	version string
}

func (p *Python) Load(definition interface{}) (bool, error) {
	def, ok := definition.(map[interface{}]interface{})
	if !ok {
		return false, nil
	}
	if version, ok := def["python"]; ok {
		p.version, ok = version.(string)
		if !ok {
			return false, nil
		}
		return true, nil
	}

	return false, nil
}

func (p *Python) Perform() (err error) {
	fmt.Printf("%s Python: %s\n", color.Brown("★"), color.Cyan(p.version))

	output, code, err := executor.Capture("pyenv", "versions", "--bare", "--skip-aliases")
	if err != nil {
		return
	}
	if code != 0 {
		return fmt.Errorf("failed to run pyenv versions. exit code: %d", code)
	}

	installedVersions := strings.Split(strings.TrimSpace(output), "\n")

	if stringInSlice(p.version, installedVersions) {
		fmt.Println(color.Green("  Already good!"))
		return nil
	}

	code, err = executor.Run("pyenv", "install", p.version)
	if err != nil {
		return
	}
	if code != 0 {
		return fmt.Errorf("failed to install the required python version. exit code: %d", code)
	}

	return nil
}
