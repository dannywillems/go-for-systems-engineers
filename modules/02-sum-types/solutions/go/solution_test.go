package solutions

import "testing"

func sample() Expr {
	return Add{L: Mul{L: Lit{1}, R: Lit{2}}, R: Lit{3}}
}

func TestHeight(t *testing.T) {
	if Height(sample()) != 3 || Height(Lit{5}) != 1 {
		t.Fatal("Height wrong")
	}
}

func TestContains(t *testing.T) {
	e := sample()
	if !Contains(e, 1) || !Contains(e, 3) || Contains(e, 99) {
		t.Fatal("Contains wrong")
	}
}
