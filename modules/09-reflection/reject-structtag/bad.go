// Package structtagbad COMPILES FINE but `go vet` rejects it. This is the point:
// a malformed struct tag is invisible to the compiler (reflection reads tags as
// opaque strings at run time, so a typo just silently fails to match at run
// time). Go's answer is a STATIC ANALYZER bolted on -- go vet's structtag check
// flags it before it ships:
//
//	struct field tag `json:name` not compatible with reflect.StructTag.Get:
//	bad syntax for struct tag value
//
// Without vet, json.Marshal would just ignore the tag at run time -- the
// error-deferral that separates reflection from compile-time codegen.
package structtagbad

type Person struct {
	Name string `json:name` // missing quotes: should be `json:"name"`
}
