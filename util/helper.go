package util

import (
	"io"
	"reflect"
)

func GetOrDefault(m map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if value, ok := m[key]; ok {
		return value
	} else {
		return defaultValue
	}
}

func GetTypeName(i interface{}) string {
	t := reflect.TypeOf(i)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

func CloseQuietly(closer io.Closer) {
	_ = closer.Close()
}
