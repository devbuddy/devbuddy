package helpers

import (
	"path/filepath"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

type Homebrew struct {
	env    *env.Env
	prefix string
}

// NewHomebrew is returning a new Cellar
func NewHomebrew(env *env.Env) *Homebrew {
	return NewHomebrewWithPrefix(env, "/usr/local")
}

// NewHomebrewWithPrefix is returning a new Cellar at prefix
func NewHomebrewWithPrefix(env *env.Env, prefix string) *Homebrew {
	if !utils.PathExists(filepath.Join(prefix, "Cellar")) {

		paths := env.GetPathParts()

		if len(paths) > 0 && utils.PathExists(paths[0]) {
			prefix = filepath.Dir(paths[0])
		}
	}

	return &Homebrew{env: env, prefix: prefix}
}

func pathToFormula(filename string) string {
	results := strings.Split(filename, "/")
	formula := results[len(results)-1]
	return strings.TrimSuffix(formula, filepath.Ext(formula))
}

func (h *Homebrew) IsInCaskroom(formula string) bool {
	path := "/opt/homebrew-cask/Caskroom"

	if !utils.PathExists(path) {
		path = filepath.Join(h.prefix, "Caskroom")
	}

	return utils.PathExists(filepath.Join(path, formula))
}

func (h *Homebrew) IsInCellar(formula string) bool {
	path := filepath.Join(h.prefix, "Cellar")

	path = filepath.Join(path, formula)
	return utils.PathExists(path)
}

// IsInstalled returns true if `pkg` is installed in cellar or in caskroom
func (h *Homebrew) IsInstalled(formula string) bool {
	formula = pathToFormula(formula)
	return h.IsInCellar(formula) || h.IsInCaskroom(formula)
}
