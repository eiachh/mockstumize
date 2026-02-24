package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Resource map[string]any

func Build(dir string) ([]Resource, error) {
	kustomization, err := LoadKustomization(dir)
	if err != nil {
		return nil, err
	}

	var resources []Resource
	for _, resourcePath := range kustomization.Resources {
		resolvedPath := filepath.Join(dir, resourcePath)
		info, err := os.Stat(resolvedPath)
		if err != nil {
			return nil, fmt.Errorf("stat resource %s: %w", resolvedPath, err)
		}

		if info.IsDir() {
			nested, err := Build(resolvedPath)
			if err != nil {
				return nil, err
			}
			resources = append(resources, nested...)
			continue
		}

		items, err := loadResources(resolvedPath)
		if err != nil {
			return nil, err
		}
		resources = append(resources, items...)
	}

	return resources, nil
}

func loadResources(path string) ([]Resource, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read resource %s: %w", path, err)
	}

	decoder := yaml.NewDecoder(strings.NewReader(string(data)))
	var resources []Resource
	for {
		var doc Resource
		if err := decoder.Decode(&doc); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("parse resource %s: %w", path, err)
		}
		if len(doc) == 0 {
			continue
		}
		resources = append(resources, doc)
	}

	return resources, nil
}
