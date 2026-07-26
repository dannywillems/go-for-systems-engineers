package exercises

import "testing"

func TestCToF(t *testing.T) {
	if got := CToF(100); got != 212 {
		t.Fatalf("CToF(100) = %v, want 212", got)
	}
	if got := CToF(0); got != 32 {
		t.Fatalf("CToF(0) = %v, want 32", got)
	}
}

func TestStack(t *testing.T) {
	var s Stack[string]
	if s.Len() != 0 {
		t.Fatalf("empty Len = %d, want 0", s.Len())
	}
	s.Push("a")
	s.Push("b")
	if s.Len() != 2 {
		t.Fatalf("Len = %d, want 2", s.Len())
	}
	if v, ok := s.Pop(); !ok || v != "b" {
		t.Fatalf("Pop = %q,%v, want b,true", v, ok)
	}
	if _, _ = s.Pop(); s.Len() != 0 {
		t.Fatal("stack should be empty")
	}
	if _, ok := s.Pop(); ok {
		t.Fatal("Pop on empty should be ok=false")
	}
}
