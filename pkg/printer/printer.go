package printer

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"

	"mockstumize/pkg/builder"
)

func Print(resources []builder.Resource) (string, error) {
	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(2)
	for index, resource := range resources {
		if index > 0 {
			buffer.WriteString("---\n")
		}
		if err := encoder.Encode(resource); err != nil {
			return "", fmt.Errorf("encode resource: %w", err)
		}
	}
	if err := encoder.Close(); err != nil {
		return "", fmt.Errorf("close encoder: %w", err)
	}
	return buffer.String(), nil
}
