package solutions

import "testing"

type sample struct {
	Name   string `json:"name"`
	Age    int    `json:"age,omitempty"`
	Nick   string // no tag: key is the field name
	secret string // unexported: skipped
}

func TestToMapExportedFields(t *testing.T) {
	m := ToMap(sample{Name: "Ada", Age: 36, Nick: "A", secret: "x"})
	if m["name"] != "Ada" {
		t.Fatalf("m[name] = %v, want Ada", m["name"])
	}
	if m["age"] != 36 {
		t.Fatalf("m[age] = %v, want 36 (tag before comma)", m["age"])
	}
	if m["Nick"] != "A" {
		t.Fatalf("m[Nick] = %v, want A (no tag -> field name)", m["Nick"])
	}
	if _, ok := m["secret"]; ok {
		t.Fatal("unexported field secret must be skipped")
	}
	if len(m) != 3 {
		t.Fatalf("len(m) = %d, want 3", len(m))
	}
}

func TestToMapNonStruct(t *testing.T) {
	if m := ToMap(42); m != nil {
		t.Fatalf("ToMap(42) = %v, want nil", m)
	}
}
