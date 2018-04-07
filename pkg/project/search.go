package project

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/pior/dad/pkg/config"
)

func FindBestMatch(expr string, conf *config.Config) (proj *Project, err error) {
	projects, err := GetAllProjects(conf.SourceDir)
	if err != nil {
		return
	}

	if len(projects) == 0 {
		err = fmt.Errorf("no projects found at all! Try cloning one first")
		return
	}

	// Exact match on ID
	for _, p := range projects {
		if p.ID == expr {
			return p, nil
		}
	}

	// Exact match on RepositoryName
	for _, p := range projects {
		if p.RepositoryName == expr {
			return p, nil
		}
	}

	// Prefix match on ID
	for _, p := range projects {
		if strings.HasPrefix(p.ID, expr) {
			return p, nil
		}
	}

	// Prefix match on RepositoryName
	for _, p := range projects {
		if strings.HasPrefix(p.RepositoryName, expr) {
			return p, nil
		}
	}

	// Other substring match on ID
	for _, p := range projects {
		if strings.Contains(p.ID, expr) {
			return p, nil
		}
	}

	// Other substring match on RepositoryName
	for _, p := range projects {
		if strings.Contains(p.RepositoryName, expr) {
			return p, nil
		}
	}

	err = fmt.Errorf("no project found for %s", expr)
	return
}

func GetAllProjects(sourceDir string) ([]*Project, error) {
	var projects []*Project

	host := "github.com"

	hostPath := filepath.Join(sourceDir, host)
	var orgPath string
	var projPath string

	orgs, err := listChildDir(hostPath)
	if err != nil {
		return nil, err
	}

	for _, org := range orgs {
		orgPath = filepath.Join(hostPath, org)

		repos, err := listChildDir(orgPath)
		if err != nil {
			return nil, err
		}

		for _, repo := range repos {
			projPath = filepath.Join(orgPath, repo)

			projects = append(projects, &Project{
				HostingPlatform:  host,
				OrganisationName: org,
				RepositoryName:   repo,
				ID:               filepath.Join(org, repo),
				Path:             projPath,
			})
		}
	}

	return projects, nil
}

func listChildDir(path string) (paths []string, err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		err = fmt.Errorf("error listing files in %s: %s", path, err)
		return
	}

	for _, f := range files {
		if f.IsDir() {
			paths = append(paths, f.Name())
		}
	}
	return
}
