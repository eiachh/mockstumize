package printer

import (
	"bytes"
	"fmt"
	"sort"

	"gopkg.in/yaml.v3"

	"mockstumize/pkg/builder"
)

// Print encodes every resource in the map into a multi-document YAML string.
// Resources are sorted by their ResID so output is deterministic.
func Print(resources builder.ResourceMap) (string, error) {
	ids := sortedResIDs(resources)

	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(2)
	for _, id := range ids {
		if err := encoder.Encode(resources[id]); err != nil {
			return "", fmt.Errorf("encode resource %s: %w", id, err)
		}
	}
	if err := encoder.Close(); err != nil {
		return "", fmt.Errorf("close encoder: %w", err)
	}
	return buffer.String(), nil
}

func sortedResIDs(resources builder.ResourceMap) []builder.ResID {
	ids := make([]builder.ResID, 0, len(resources))
	for id := range resources {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i].String() < ids[j].String()
	})
	return ids
}
