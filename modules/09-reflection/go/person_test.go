package reflectgen

import (
	"encoding/json"
	"testing"
)

// TestUnexportedFieldSilentlyDropped documents reflection's blind spot: the
// unexported field is set on the value but never appears in the JSON, and no
// error is raised. This is the error-deferral the module is about.
func TestUnexportedFieldSilentlyDropped(t *testing.T) {
	p := NewPerson("Ada", 36, "top-secret")
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got := string(b)
	want := `{"name":"Ada","age":36}`
	if got != want {
		t.Fatalf("Marshal = %s, want %s", got, want)
	}
	if p.Secret() != "top-secret" {
		t.Fatal("secret field should still be set on the value")
	}
}

// TestManualMatchesReflection: the hand-written encoder and reflection agree, so
// the benchmark compares two paths to the same output.
func TestManualMatchesReflection(t *testing.T) {
	p := NewPerson("Ada", 36, "")
	b, _ := json.Marshal(p)
	if ManualMarshal(p) != string(b) {
		t.Fatalf("ManualMarshal = %s, reflection = %s", ManualMarshal(p), b)
	}
}

// TestDescribeAnyStruct: Describe works on a type it was not written for.
func TestDescribeAnyStruct(t *testing.T) {
	type Point struct {
		X, Y int
	}
	got := Describe(Point{X: 1, Y: 2})
	if len(got) != 2 {
		t.Fatalf("Describe(Point) = %v, want 2 fields", got)
	}
}

var benchSink string

func BenchmarkReflectMarshal(b *testing.B) {
	p := NewPerson("Ada", 36, "")
	b.ReportAllocs()
	for b.Loop() {
		out, _ := json.Marshal(p)
		benchSink = string(out)
	}
}

func BenchmarkManualMarshal(b *testing.B) {
	p := NewPerson("Ada", 36, "")
	b.ReportAllocs()
	for b.Loop() {
		benchSink = ManualMarshal(p)
	}
}
