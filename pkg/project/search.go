package project

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/sahilm/fuzzy"

	"github.com/pior/dad/pkg/config"
)

func FindBestMatch(expr string, conf *config.Config) (found *Project, err error) {
	projects, err := GetAllProjects(conf.SourceDir)
	if err != nil {
		return
	}

	if len(projects) == 0 {
		err = fmt.Errorf("no projects found at all! Try cloning one first")
		return
	}

	found = projectMatch(expr, projects)
	if found == nil {
		err = fmt.Errorf("no project found for %s", expr)
	}
	return
}

func projectMatch(expr string, projects []*Project) *Project {
	// First, try to match on project name only
	names := []string{}
	for _, p := range projects {
		names = append(names, p.RepositoryName)
	}
	matches := fuzzy.Find(expr, names)
	if matches.Len() >= 1 {
		return projects[matches[0].Index]
	}

	// Then, extend match to the organisation name as well
	names = []string{}
	for _, p := range projects {
		names = append(names, p.id)
	}
	matches = fuzzy.Find(expr, names)
	if matches.Len() >= 1 {
		return projects[matches[0].Index]
	}

	return nil
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
				id:               filepath.Join(org, repo),
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
