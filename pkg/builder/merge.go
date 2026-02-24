package builder

// MergeResMaps merges patches into base resources. Each patch is matched to a
// base resource by its ResID (apiVersion + kind + name). Matching resources are
// deep-merged so that the patch values override the base. Patches that do not
// match an existing resource are added as new resources.
func MergeResMaps(base, patches ResourceMap) (ResourceMap, error) {
	result := make(ResourceMap, len(base))

	// Start with a copy of every base resource.
	for id, res := range base {
		result[id] = copyResource(res)
	}

	// Apply each patch on top of its matching base resource.
	for id, patch := range patches {
		baseRes, exists := result[id]
		if !exists {
			result[id] = copyResource(patch)
			continue
		}
		result[id] = deepMerge(baseRes, patch)
	}

	return result, nil
}

// deepMerge recursively merges src into dst. Values in src take precedence.
// Maps are merged key-by-key; all other types are replaced outright.
func deepMerge(dst, src Resource) Resource {
	merged := make(Resource, len(dst))
	for k, v := range dst {
		merged[k] = v
	}
	for k, srcVal := range src {
		dstVal, exists := merged[k]
		if !exists {
			merged[k] = srcVal
			continue
		}

		srcMap := toResource(srcVal)
		dstMap := toResource(dstVal)
		if srcMap != nil && dstMap != nil {
			merged[k] = deepMerge(dstMap, srcMap)
		} else {
			merged[k] = srcVal
		}
	}
	return merged
}

// toResource attempts to convert a value to a Resource (map[string]any).
// Returns nil if the value is not a map type.
func toResource(v any) Resource {
	switch m := v.(type) {
	case Resource:
		return m
	case map[string]any:
		return Resource(m)
	default:
		return nil
	}
}

// copyResource creates a shallow copy of a Resource map.
func copyResource(r Resource) Resource {
	cp := make(Resource, len(r))
	for k, v := range r {
		cp[k] = v
	}
	return cp
}
