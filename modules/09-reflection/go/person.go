// Package reflectgen is Module 09: how each language turns a typed value into a
// generic operation (serialize, compare, print) it did not hand-write. Go's
// answer is RUNTIME REFLECTION: encoding/json and text/template walk a value's
// type descriptor at run time via the reflect package. There is no compile-time
// step and no generated code; the type information travels inside the value and
// is inspected dynamically. The cost, and the deferral of errors to run time, is
// the subject of this module.
package reflectgen

// region:tags:start

// Person carries struct TAGS that reflection reads at run time. encoding/json
// never sees this type at compile time; it discovers Name/Age and their tags by
// walking reflect.Type when Marshal is called. secret is unexported, so
// reflection from another package cannot see it at all -- it is silently omitted
// from the JSON, with no error.
//
//nolint:govet // field order follows the demo/JSON narrative, not byte alignment
type Person struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	secret string // unexported: invisible to reflection, silently dropped
}

// region:tags:end

// NewPerson builds a Person including the unexported field, so a test can prove
// the field exists yet never appears in the JSON.
func NewPerson(name string, age int, secret string) Person {
	return Person{Name: name, Age: age, secret: secret}
}

// Secret exposes the unexported field for tests.
func (p Person) Secret() string { return p.secret }
