package mem

import "testing"

func TestPadding(t *testing.T) {
	padded, packed, saved := Sizes()
	if padded != 24 {
		t.Errorf("sizeof(Padded) = %d, want 24 (bool/int64/bool with padding)", padded)
	}
	if packed != 16 {
		t.Errorf("sizeof(Packed) = %d, want 16 (int64/bool/bool)", packed)
	}
	if saved != 8 {
		t.Errorf("saved = %d, want 8", saved)
	}
}

func TestAliasing(t *testing.T) {
	before, after := AliasParent()
	if before[2] == after[2] {
		t.Errorf("expected aliasing to overwrite orig[2]: before=%v after=%v", before, after)
	}
	if after[2] != 99 {
		t.Errorf("orig[2] = %d, want 99 (overwritten by append)", after[2])
	}
}

func TestEscapeBehaviorCompiles(t *testing.T) {
	if NoEscape() != 5 {
		t.Error("NoEscape")
	}
	if *EscapesReturn() != 42 {
		t.Error("EscapesReturn")
	}
	EscapesInterface(1)
}
