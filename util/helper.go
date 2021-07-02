package util

func GetOrDefault(m map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if value, ok := m[key]; ok {
		return value
	} else {
		return defaultValue
	}
}
