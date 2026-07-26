// Package solutions is the corrigé for Module 09. Run via `make solutions M=09`.
package solutions

import (
	"reflect"
	"strings"
)

func ToMap(v any) map[string]any {
	t := reflect.TypeOf(v)
	if t == nil || t.Kind() != reflect.Struct {
		return nil
	}
	val := reflect.ValueOf(v)
	m := make(map[string]any, t.NumField())
	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		key := f.Name
		if tag := f.Tag.Get("json"); tag != "" {
			if name, _, _ := strings.Cut(tag, ","); name != "" {
				key = name
			}
		}
		m[key] = val.Field(i).Interface()
	}
	return m
}
