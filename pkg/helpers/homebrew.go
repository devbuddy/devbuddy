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

func pathToPackage(filename string) string {
	results := strings.Split(filename, "/")
	pkg := results[len(results)-1]
	return strings.TrimSuffix(pkg, filepath.Ext(pkg))
}

func (h *Homebrew) PackageIsInCaskroom(pkg string) bool {
	path := "/opt/homebrew-cask/Caskroom"

	if !utils.PathExists("/opt/homebrew-cask/Caskroom") {
		path = filepath.Join(h.prefix, "Caskroom")
	}

	return utils.PathExists(filepath.Join(path, pkg))
}

func (h *Homebrew) PackageIsInCellar(pkg string) bool {
	path := filepath.Join(h.prefix, "Cellar")

	path = filepath.Join(path, pkg)
	return utils.PathExists(path)
}

// PackageIsInstalled returns true if `pkg` is installed in cellar or in caskroom
func (h *Homebrew) PackageIsInstalled(pkg string) bool {
	pkg = pathToPackage(pkg)
	return h.PackageIsInCellar(pkg) || h.PackageIsInCaskroom(pkg)
}
