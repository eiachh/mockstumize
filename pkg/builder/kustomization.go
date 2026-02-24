package builder

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const kustomizationFileName = "kustomization.yaml"

type Kustomization struct {
	Resources []string `yaml:"resources"`
	Patches   []string `yaml:"patchesStrategicMerge,omitempty"`
}

func LoadKustomization(dir string) (Kustomization, error) {
	path := filepath.Join(dir, kustomizationFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		return Kustomization{}, fmt.Errorf("read kustomization: %w", err)
	}

	var k Kustomization
	if err := yaml.Unmarshal(data, &k); err != nil {
		return Kustomization{}, fmt.Errorf("parse kustomization: %w", err)
	}

	return k, nil
}
