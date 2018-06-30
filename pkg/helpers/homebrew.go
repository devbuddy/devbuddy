package helpers

import (
	"path/filepath"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/utils"
)

type caskroom struct {
	prefix string
}

type cellar struct {
	prefix string
}

// Homebrew represent an homebrew installation
type Homebrew struct {
	caskroom *caskroom
	cellar   *cellar
}

// NewHomebrew is returning a new Cellar
func NewHomebrew() *Homebrew {
	return NewHomebrewWithPrefix("/usr/local")
}

// NewHomebrewWithPrefix is returning a new Cellar at prefix
func NewHomebrewWithPrefix(prefix string) *Homebrew {
	return &Homebrew{
		cellar:   &cellar{prefix: prefix},
		caskroom: &caskroom{prefix: prefix},
	}
}

// IsInstalled returns true if `pkg` is installed in cellar or in caskroom
func (h *Homebrew) IsInstalled(formula string) (installed bool) {
	path := buildFormulaPath(formula)

	return h.cellar.isInstalled(path) || h.caskroom.isInstalled(path)
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

// isInstalled returns true if formulua was installed in Caskrook
func (c *caskroom) isInstalled(formula string) bool {
	legacyPrefix := "/opt/homebrew-cask"

	return utils.PathExists(filepath.Join(legacyPrefix, "Caskroom", formula)) ||
		utils.PathExists(filepath.Join(c.prefix, "Caskroom", formula))
}

// isInstalled returns true if formulua was installed in cellar
func (c *cellar) isInstalled(formula string) bool {
	return utils.PathExists(filepath.Join(c.prefix, "Cellar", formula))
}
