package solutions

import "testing"

func TestCToF(t *testing.T) {
	if CToF(100) != 212 || CToF(0) != 32 {
		t.Fatal("CToF wrong")
	}
}

func TestStack(t *testing.T) {
	var s Stack[int]
	s.Push(1)
	s.Push(2)
	if s.Len() != 2 {
		t.Fatal("len")
	}
	if v, ok := s.Pop(); !ok || v != 2 {
		t.Fatal("pop")
	}
	if _, ok := (&Stack[int]{}).Pop(); ok {
		t.Fatal("empty pop should be false")
	}
}
