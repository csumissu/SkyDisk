package util

import "reflect"

func GetOrDefault(m map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if value, ok := m[key]; ok {
		return value
	} else {
		return defaultValue
	}
}

func Convert(i ...interface{}) map[string]interface{} {
	if i == nil || len(i) == 0 {
		return nil
	}

	m := make(map[string]interface{})
	for _, v := range i {
		m[GetTypeName(v)] = v
	}
	return m
}

func GetTypeName(i interface{}) string {
	t := reflect.TypeOf(i)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}
