package ts

import "testing"

func TestPhantomIDsCompileAndWork(t *testing.T) {
	// Compile-time assertions that the distinct phantom types exist.
	var (
		_ UserID
		_ OrderID
	)
	if LookupUser(UserID(42)) != "user #42" {
		t.Fatal("LookupUser wrong")
	}
}
