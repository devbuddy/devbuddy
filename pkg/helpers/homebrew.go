package helpers

import (
	"path/filepath"
	"strings"

	"github.com/pior/dad/pkg/env"
	"github.com/pior/dad/pkg/utils"
)

type Homebrew struct {
	prefix string
}

// NewHomebrew is returning a new Cellar
func NewHomebrew() *Homebrew {
	return NewHomebrewWithPrefix("/usr/local")
}

func NewHomebrewWithPrefix(prefix string) *Homebrew {
	if !utils.PathExists(filepath.Join(prefix, "Cellar")) {

		path := strings.Split(env.NewFromOS().Get("PATH"), ":")

		if len(path) > 0 && utils.PathExists(path[0]) {
			prefix = filepath.Dir(path[0])
		}
	}

	return &Homebrew{prefix: prefix}
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
