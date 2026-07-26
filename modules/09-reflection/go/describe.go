package reflectgen

import (
	"fmt"
	"reflect"
	"strconv"
)

// region:reflect:start

// Describe is what reflection buys: one function that inspects the fields of ANY
// struct, including types it has never heard of, with no per-type code. It reads
// the reflect.Type at run time. This is impossible with compile-time codegen
// alone -- the derive/synthesis approaches must be asked for each type in
// advance -- but every call pays the reflection walk, and a non-struct argument
// is a run-time error, not a compile-time one.
func Describe(v any) []string {
	t := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	if t.Kind() != reflect.Struct {
		return []string{fmt.Sprintf("<not a struct: %s>", t.Kind())}
	}
	var out []string
	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() { // reflection sees the field but cannot read its value
			out = append(out, fmt.Sprintf("%s %s (unexported)", f.Name, f.Type))
			continue
		}
		tag := f.Tag.Get("json")
		out = append(out, fmt.Sprintf("%s %s json=%q = %v",
			f.Name, f.Type, tag, val.Field(i).Interface()))
	}
	return out
}

// region:reflect:end

// ManualMarshal is a hand-written JSON encoder for Person: no reflection, no
// generated code, allocation-lean (one buffer, one final string). It is the
// baseline the reflection path is measured against.
func ManualMarshal(p Person) string {
	buf := make([]byte, 0, 32)
	buf = append(buf, `{"name":`...)
	buf = appendQuoted(buf, p.Name)
	buf = append(buf, `,"age":`...)
	buf = strconv.AppendInt(buf, int64(p.Age), 10)
	buf = append(buf, '}')
	return string(buf)
}

func appendQuoted(buf []byte, s string) []byte {
	buf = append(buf, '"')
	for i := range len(s) {
		if c := s[i]; c == '"' || c == '\\' {
			buf = append(buf, '\\', c)
		} else {
			buf = append(buf, c)
		}
	}
	return append(buf, '"')
}
