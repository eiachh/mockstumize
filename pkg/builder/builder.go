package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Resource is a generic YAML document represented as an unstructured map.
type Resource map[string]any

// ResID uniquely identifies a Kubernetes-style resource by its
// apiVersion, kind, and name.
type ResID struct {
	APIVersion string
	Kind       string
	Name       string
}

func (id ResID) String() string {
	return id.APIVersion + "/" + id.Kind + "/" + id.Name
}

// ResourceMap maps each ResID to its full resource body.
type ResourceMap map[ResID]Resource

// Build reads a kustomization directory and returns a ResourceMap containing
// every resource referenced (directly or transitively), with any
// patchesStrategicMerge applied on top.
func Build(dir string) (ResourceMap, error) {
	kustomization, err := LoadKustomization(dir)
	if err != nil {
		return nil, err
	}

	resources := make(ResourceMap)
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
			for id, res := range nested {
				resources[id] = res
			}
			continue
		}

		items, err := loadResources(resolvedPath)
		if err != nil {
			return nil, err
		}
		for id, res := range items {
			resources[id] = res
		}
	}

	// Load and apply strategic merge patches.
	if len(kustomization.Patches) > 0 {
		patches := make(ResourceMap)
		for _, patchPath := range kustomization.Patches {
			resolvedPath := filepath.Join(dir, patchPath)
			items, err := loadResources(resolvedPath)
			if err != nil {
				return nil, fmt.Errorf("load patch %s: %w", patchPath, err)
			}
			for id, res := range items {
				patches[id] = res
			}
		}

		resources, err = MergeResMaps(resources, patches)
		if err != nil {
			return nil, err
		}
	}

	return resources, nil
}

// resIDFromResource extracts a ResID from a generic YAML map by reading
// the apiVersion, kind, and metadata.name fields.
func resIDFromResource(r Resource) (ResID, error) {
	apiVersion, _ := r["apiVersion"].(string)
	kind, _ := r["kind"].(string)

	var name string
	switch meta := r["metadata"].(type) {
	case Resource:
		name, _ = meta["name"].(string)
	case map[string]any:
		name, _ = meta["name"].(string)
	}

	if kind == "" || name == "" {
		return ResID{}, fmt.Errorf("resource is missing kind or metadata.name")
	}

	return ResID{
		APIVersion: apiVersion,
		Kind:       kind,
		Name:       name,
	}, nil
}

func loadResources(path string) (ResourceMap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read resource %s: %w", path, err)
	}

	decoder := yaml.NewDecoder(strings.NewReader(string(data)))
	resources := make(ResourceMap)
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

		id, err := resIDFromResource(doc)
		if err != nil {
			return nil, fmt.Errorf("identify resource in %s: %w", path, err)
		}
		resources[id] = doc
	}

	return resources, nil
}
