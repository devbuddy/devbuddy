package project

import (
	"fmt"
	"hash/adler32"
	"os"
	"path/filepath"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/executor"
)

// Project represents a project whether it exists locally or not
type Project struct {
	hosting *hostingInfo
	Path    string // Local path of this project on disk
}

// NewFromID creates an instance of Project from a short id like "org/name" or "name"
func NewFromID(id string, conf *config.Config) (p *Project, err error) {
	if !strings.Contains(id, "/") {
		if conf.DefaultOrg != "" {
			id = conf.DefaultOrg + "/" + id
		}
	}

	hosting, err := newHostingInfoByURL(id)
	if err != nil {
		return nil, err
	}

	path := filepath.Join(conf.SourceDir, hosting.platform, hosting.organisation, hosting.repository)

	p = &Project{
		hosting: hosting,
		Path:    path,
	}
	return
}

// NewFromPath creates an instance of Project from the local path
func NewFromPath(path string) *Project {
	return &Project{
		hosting: newHostingInfoByPath(path),
		Path:    path,
	}
}

// Name returns a logical id like platform:org/project
func (p *Project) Name() string {
	return p.hosting.repository
}

// FullName returns a logical id like platform:org/project
func (p *Project) FullName() string {
	return fmt.Sprintf("%s:%s/%s", p.hosting.platform, p.hosting.organisation, p.hosting.repository)
}

// Slug returns a short, unique but humanly recognizable id based on the path
func (p *Project) Slug() string {
	locationToken := adler32.Checksum([]byte(filepath.Clean(p.Path)))
	return fmt.Sprintf("%s-%d", p.hosting.repository, locationToken)
}

// GetRemoteURL builds the Git remote url for the project
func (p *Project) GetRemoteURL() (url string, err error) {
	return p.hosting.remoteURL, nil
}

// Exists checks whether the project exists locally
func (p *Project) Exists() bool {
	if p.Path == "" {
		return false
	}
	if _, err := os.Stat(p.Path); err == nil {
		return true
	}
	return false
}

// Clone runs the Git command needed to clone the project
func (p *Project) Clone(exec *executor.Executor) (err error) {
	err = os.MkdirAll(filepath.Dir(p.Path), 0755)
	if err != nil {
		return
	}

	return exec.Run(executor.New("git", "clone", p.hosting.remoteURL, p.Path)).Error
}

// Create creates the project directory locally
func (p *Project) Create() (err error) {
	err = os.MkdirAll(p.Path, 0755)
	return
}
