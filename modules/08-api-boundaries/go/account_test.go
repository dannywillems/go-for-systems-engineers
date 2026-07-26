package api_test

import (
	"errors"
	"testing"

	"api"
)

// The test lives in the EXTERNAL test package api_test, so it sees only the
// exported surface -- exactly what a real consumer sees. It cannot touch
// a.balance; that is the point.
func TestAccountThroughPublicAPI(t *testing.T) {
	a, err := api.Open(100)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	a.Deposit(50)
	if got := a.Balance(); got != 150 {
		t.Fatalf("Balance = %d, want 150", got)
	}
	if err := a.Withdraw(200); !errors.Is(err, api.ErrOverdraft) {
		t.Fatalf("Withdraw over balance = %v, want ErrOverdraft", err)
	}
	if got := a.Balance(); got != 150 {
		t.Fatalf("Balance after failed withdraw = %d, want 150", got)
	}
}

func TestOpenRejectsNegative(t *testing.T) {
	if _, err := api.Open(-1); !errors.Is(err, api.ErrOverdraft) {
		t.Fatalf("Open(-1) = %v, want ErrOverdraft", err)
	}
}
