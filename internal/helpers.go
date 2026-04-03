package internal

func getModuleName(config map[string]any) string {
	if v, ok := config["module"].(string); ok && v != "" {
		return v
	}
	return "crm"
}

func resolveValue(key string, current, config map[string]any) string {
	if v, ok := current[key].(string); ok && v != "" {
		return v
	}
	if v, ok := config[key].(string); ok && v != "" {
		return v
	}
	return ""
}

func resolveMap(key string, current, config map[string]any) map[string]any {
	if v, ok := current[key].(map[string]any); ok {
		return v
	}
	if v, ok := config[key].(map[string]any); ok {
		return v
	}
	return nil
}

func resolveAnySlice(key string, current, config map[string]any) []map[string]any {
	for _, m := range []map[string]any{current, config} {
		switch v := m[key].(type) {
		case []map[string]any:
			return v
		case []any:
			result := make([]map[string]any, 0, len(v))
			for _, item := range v {
				if m, ok := item.(map[string]any); ok {
					result = append(result, m)
				}
			}
			return result
		}
	}
	return nil
}
