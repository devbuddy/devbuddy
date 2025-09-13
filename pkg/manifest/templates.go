package manifest

import (
	"embed"
	"io/fs"
	"strings"
)

//go:embed templates/*
var templates embed.FS

func ListTemplates() []string {
	entries, err := fs.ReadDir(templates, "templates")
	if err != nil {
		return nil
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() {
			name := strings.TrimSuffix(entry.Name(), ".yml")
			names = append(names, name)
		}
	}
	return names
}

func LoadTemplate(name string) ([]byte, error) {
	// Read the template file from the embedded filesystem
	template, err := fs.ReadFile(templates, "templates/"+name+".yml")
	if err != nil {
		return nil, err
	}
	return template, nil
}
