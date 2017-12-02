package project

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/executor"
)

type Project struct {
	HostingPlatform  string
	OrganisationName string
	RepositoryName   string
	Path             string
}

func NewFromIdentifier(id string) (p *Project, err error) {
	p = &Project{HostingPlatform: "github.com"}

	if match := regexp.MustCompile(`([^/]+)/([^/]+)`).FindStringSubmatch(id); match != nil {
		p.OrganisationName = match[1]
		p.RepositoryName = match[2]
		return
	}

	return
}

func (p *Project) GetRemoteUrl() (url string, err error) {
	if p.HostingPlatform == "github.com" {
		url = fmt.Sprintf("git@github.com:%s/%s.git", p.OrganisationName, p.RepositoryName)
		return
	}
	err = fmt.Errorf("Unknown project hosting platform: %s", p.HostingPlatform)
	return
}

func (p *Project) InferPath(conf *config.Config) {
	p.Path = filepath.Join(conf.SourceDir, p.HostingPlatform, p.OrganisationName, p.RepositoryName)
}

func (p *Project) Exists() bool {
	if p.Path == "" {
		panic("Project path can't be null")
	}
	if _, err := os.Stat(p.Path); err == nil {
		return true
	}
	return false
}

func (p *Project) Clone() (err error) {
	err = os.MkdirAll(filepath.Dir(p.Path), 0755)
	if err != nil {
		return
	}

	url, err := p.GetRemoteUrl()
	if err != nil {
		return
	}
	err = executor.Run("git", "clone", url, p.Path)
	return
}
