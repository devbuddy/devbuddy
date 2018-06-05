package project

import (
	"fmt"
	"hash/adler32"
	"os"
	"path/filepath"
	"regexp"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/manifest"
)

// Project represents a project whether it exists locally or not
type Project struct {
	HostingPlatform  string // Name and directory name of the hosting platform like "github.com"
	OrganisationName string // Name and directory name of the organisation owning this project
	RepositoryName   string // Name and directory name of this project
	id               string // Short id like "org/name"
	Path             string // Full path of this project

	Manifest *manifest.Manifest // Manifest of this project
}

// NewFromID creates an instance of Project from a short id like "org/name"
func NewFromID(id string, conf *config.Config) (p *Project, err error) {
	reGithubFull := regexp.MustCompile(`([^/]+)/([^/]+)`)

	if match := reGithubFull.FindStringSubmatch(id); match != nil {
		p = &Project{
			HostingPlatform:  "github.com",
			OrganisationName: match[1],
			RepositoryName:   match[2],
			id:               id,
		}
	} else {
		err = fmt.Errorf("Unrecognized remote project: %s", id)
		return
	}

	p.Path = filepath.Join(conf.SourceDir, p.HostingPlatform, p.OrganisationName, p.RepositoryName)
	return
}

// FullName returns a logical id like platform:org/project
func (p *Project) FullName() string {
	return fmt.Sprintf("%s:%s/%s", p.HostingPlatform, p.OrganisationName, p.RepositoryName)
}

// Slug returns a short, unique but humanly recognizable id based on the path
func (p *Project) Slug() string {
	locationToken := adler32.Checksum([]byte(filepath.Clean(p.Path)))
	return fmt.Sprintf("%s-%d", p.RepositoryName, locationToken)
}

// GetRemoteURL builds the Git remote url for the project
func (p *Project) GetRemoteURL() (url string, err error) {
	if p.HostingPlatform == "github.com" {
		url = fmt.Sprintf("git@github.com:%s/%s.git", p.OrganisationName, p.RepositoryName)
		return
	}
	err = fmt.Errorf("Unknown project hosting platform: %s", p.HostingPlatform)
	return
}

// Exists checks whether the project exists locally
func (p *Project) Exists() bool {
	if p.Path == "" {
		panic("Project path can't be null")
	}
	if _, err := os.Stat(p.Path); err == nil {
		return true
	}
	return false
}

// Clone runs the Git command needed to clone the project
func (p *Project) Clone() (err error) {
	err = os.MkdirAll(filepath.Dir(p.Path), 0755)
	if err != nil {
		return
	}

	url, err := p.GetRemoteURL()
	if err != nil {
		return
	}

	return executor.New("git", "clone", url, p.Path).Run()
}

// Create creates the project directory locally
func (p *Project) Create() (err error) {
	err = os.MkdirAll(p.Path, 0755)
	return
}
