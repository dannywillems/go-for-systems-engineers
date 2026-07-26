package exercises

import "testing"

// (1 * 2) + 3
func sample() Expr {
	return Add{L: Mul{L: Lit{1}, R: Lit{2}}, R: Lit{3}}
}

func TestHeight(t *testing.T) {
	if got := Height(sample()); got != 3 {
		t.Fatalf("Height = %d, want 3", got)
	}
	if got := Height(Lit{5}); got != 1 {
		t.Fatalf("Height(Lit) = %d, want 1", got)
	}
}

func TestContains(t *testing.T) {
	e := sample()
	for _, v := range []int{1, 2, 3} {
		if !Contains(e, v) {
			t.Fatalf("Contains(%d) = false, want true", v)
		}
	}
	if Contains(e, 99) {
		t.Fatal("Contains(99) = true, want false")
	}
}
