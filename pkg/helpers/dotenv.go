package helpers

import (
	"fmt"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/utils"
)

func DotenvRead(path string) (map[string]string, error) {
	lines, err := utils.ReadLines(path)
	if err != nil {
		return nil, err
	}

	vars := map[string]string{}
	for _, line := range lines {
		elems := strings.SplitN(line, "=", 2)
		if len(elems) != 2 {
			return vars, fmt.Errorf("invalid dotenv file: \"%s\"", line)
		}
		vars[elems[0]] = elems[1]
	}
	return vars, nil
}
