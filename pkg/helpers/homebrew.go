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
	prefix := "/usr/local"

	return &Homebrew{
		cellar:   &cellar{prefix: prefix},
		caskroom: &caskroom{prefix: prefix},
	}
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

	return h.cellar.IsInstalled(path) || h.caskroom.IsInstalled(path)
}

func buildFormulaPath(filename string) string {
	results := strings.Split(filename, "/")
	formula := results[len(results)-1]
	return strings.TrimSuffix(formula, filepath.Ext(formula))
}

// IsInstalled returns true if formulua was installed in Caskrook
func (c *caskroom) IsInstalled(formula string) bool {
	path := "/opt/homebrew-cask/Caskroom"

	if !utils.PathExists(path) {
		path = filepath.Join(c.prefix, "Caskroom")
	}

	return utils.PathExists(filepath.Join(path, formula))
}

// IsInstalled returns true if formulua was installed in cellar
func (c *cellar) IsInstalled(formula string) bool {
	path := filepath.Join(c.prefix, "Cellar")

	path = filepath.Join(path, formula)
	return utils.PathExists(path)
}
