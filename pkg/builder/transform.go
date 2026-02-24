package builder

import "fmt"

// ApplyNameTransformers updates resource names with the given prefix/suffix.
// The ResID keys are updated to match the new metadata.name value.
func ApplyNameTransformers(resources ResourceMap, prefix, suffix string) (ResourceMap, error) {
	if prefix == "" && suffix == "" {
		return resources, nil
	}

	transformed := make(ResourceMap, len(resources))
	for id, res := range resources {
		meta, ok := res["metadata"].(Resource)
		if !ok {
			if rawMeta, ok := res["metadata"].(map[string]any); ok {
				meta = Resource(rawMeta)
			} else {
				return nil, fmt.Errorf("resource %s is missing metadata", id)
			}
		}
		nameValue, _ := meta["name"].(string)
		if nameValue == "" {
			return nil, fmt.Errorf("resource %s is missing metadata.name", id)
		}

		newName := prefix + nameValue + suffix
		meta["name"] = newName
		res["metadata"] = meta

		newID := ResID{
			APIVersion: id.APIVersion,
			Kind:       id.Kind,
			Name:       newName,
		}
		transformed[newID] = res
	}

	return transformed, nil
}
