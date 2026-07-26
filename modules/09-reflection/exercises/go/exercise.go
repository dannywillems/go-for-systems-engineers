// Package exercises: Module 09 reader task. RED until you implement ToMap using
// reflection, so it works on ANY struct without knowing its type in advance.
package exercises

// TODO(reader): using the reflect package, return a map from field name to value
// for every EXPORTED field of the struct v. If a field has a `json:"..."` tag,
// use the tag's name as the key instead of the field name. Skip unexported
// fields. If v is not a struct, return nil.
//
// Hints: reflect.TypeOf, reflect.ValueOf, t.Kind() == reflect.Struct,
// t.NumField(), t.Field(i), f.IsExported(), f.Tag.Get("json"),
// val.Field(i).Interface(). A json tag may be "name,omitempty" -- take the part
// before the first comma.
func ToMap(v any) map[string]any {
	return nil // replace me
}
