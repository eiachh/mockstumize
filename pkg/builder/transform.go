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
		meta, err := ensureMetadata(res, id)
		if err != nil {
			return nil, err
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

// ApplyCommonMetadata attaches labels and annotations to each resource.
func ApplyCommonMetadata(resources ResourceMap, labels, annotations map[string]string) (ResourceMap, error) {
	if len(labels) == 0 && len(annotations) == 0 {
		return resources, nil
	}

	transformed := make(ResourceMap, len(resources))
	for id, res := range resources {
		meta, err := ensureMetadata(res, id)
		if err != nil {
			return nil, err
		}

		if len(labels) > 0 {
			currentLabels := ensureStringMap(meta, "labels")
			for key, value := range labels {
				currentLabels[key] = value
			}
			meta["labels"] = currentLabels
		}

		if len(annotations) > 0 {
			currentAnnotations := ensureStringMap(meta, "annotations")
			for key, value := range annotations {
				currentAnnotations[key] = value
			}
			meta["annotations"] = currentAnnotations
		}

		res["metadata"] = meta
		transformed[id] = res
	}

	return transformed, nil
}

func ensureMetadata(res Resource, id ResID) (Resource, error) {
	meta, ok := res["metadata"].(Resource)
	if ok {
		return meta, nil
	}
	if rawMeta, ok := res["metadata"].(map[string]any); ok {
		return Resource(rawMeta), nil
	}
	return nil, fmt.Errorf("resource %s is missing metadata", id)
}

func ensureStringMap(meta Resource, key string) map[string]string {
	if current, ok := meta[key].(map[string]any); ok {
		result := make(map[string]string, len(current))
		for k, v := range current {
			if value, ok := v.(string); ok {
				result[k] = value
			}
		}
		return result
	}
	if current, ok := meta[key].(map[string]string); ok {
		result := make(map[string]string, len(current))
		for k, v := range current {
			result[k] = v
		}
		return result
	}
	return make(map[string]string)
}
