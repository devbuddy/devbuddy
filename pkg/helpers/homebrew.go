package helpers

import (
	"path/filepath"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/utils"
)

type caskroom struct {
	prefixes []string
}

type cellar struct {
	prefixes []string
}

// Homebrew represent an homebrew installation
type Homebrew struct {
	caskroom *caskroom
	cellar   *cellar
}

// NewHomebrew is returning a new Cellar
func NewHomebrew() *Homebrew {
	return NewHomebrewWithPrefix("/usr/local", "/opt/homebrew")
}

// NewHomebrewWithPrefix is returning a new Cellar at prefix
func NewHomebrewWithPrefix(prefixes ...string) *Homebrew {
	return &Homebrew{
		cellar:   &cellar{prefixes: prefixes},
		caskroom: &caskroom{prefixes: prefixes},
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

	if utils.PathExists(filepath.Join(legacyPrefix, "Caskroom", formula)) {
		return true
	}

	for _, prefix := range c.prefixes {
		if utils.PathExists(filepath.Join(prefix, "Caskroom", formula)) {
			return true
		}
	}
	return false
}

// isInstalled returns true if formulua was installed in cellar
func (c *cellar) isInstalled(formula string) bool {
	for _, prefix := range c.prefixes {
		if utils.PathExists(filepath.Join(prefix, "Cellar", formula)) {
			return true
		}
	}
	return false
}
