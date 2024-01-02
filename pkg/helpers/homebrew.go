package helpers

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/utils"
)

const EnvHomebrewPrefix = "HOMEBREW_PREFIX"

func HomeBrewPackageIsInstalled(formula string) bool {
	formulaPath := buildFormulaPath(formula)

	var searchPaths []string

	prefix, ok := os.LookupEnv(EnvHomebrewPrefix)
	if ok {
		searchPaths = []string{
			filepath.Join(prefix, "Caskroom", formulaPath),
			filepath.Join(prefix, "Cellar", formulaPath),
		}
	} else {
		searchPaths = []string{
			filepath.Join("/usr/local", "Caskroom", formulaPath),         // on Intel
			filepath.Join("/opt/homebrew", "Caskroom", formulaPath),      // on Apple Silicon
			filepath.Join("/opt/homebrew-cask", "Caskroom", formulaPath), // legacy prefix
			filepath.Join("/usr/local", "Cellar", formulaPath),           // on Intel
			filepath.Join("/opt/homebrew", "Cellar", formulaPath),        // on Apple Silicon
		}
	}

	for _, searchPath := range searchPaths {
		if utils.PathExists(searchPath) {
			return true
		}
	}
	return false
}

// buildFormulaPath building a formula name from a full path by doing the following operations:
// 1. split the path by `/`
// 2. returns the filename
// 3. removes the file extension
// 4. returns the resulting formula name
func buildFormulaPath(path string) string {
	results := strings.Split(path, "/")
	formula := results[len(results)-1]
	return strings.TrimSuffix(formula, filepath.Ext(formula))
}
