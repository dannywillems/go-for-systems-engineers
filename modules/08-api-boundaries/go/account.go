// Package api is Module 08: how each language enforces API boundaries at COMPILE
// time. Go's unit of encapsulation is the PACKAGE, and visibility is by
// identifier case: a Capitalized name is exported, a lowercase name is package-
// private and the compiler refuses cross-package access to it.
//
// This gives Go an "abstract type": an exported struct whose FIELDS are all
// unexported is opaque to other packages -- they can hold an Account and call
// its methods, but cannot read, write, or construct its representation. The type
// NAME is exported; the representation is not.
package api

import "errors"

// ErrOverdraft is a sentinel the caller can match with errors.Is.
var ErrOverdraft = errors.New("api: overdraft")

// region:opaque:start

// Account is opaque: every field is unexported, so no other package can read or
// write balance, nor build an Account with a struct literal. The only way to get
// one is Open, and the only way to change it is the methods.
type Account struct {
	balance int64
}

// Open is the sole constructor: it enforces the invariant (balance >= 0) that
// the hidden field guarantees from then on.
func Open(initial int64) (*Account, error) {
	if initial < 0 {
		return nil, ErrOverdraft
	}
	return &Account{balance: initial}, nil
}

// Deposit and Withdraw are the only mutators; Balance is the only reader.
func (a *Account) Deposit(amount int64) { a.balance += amount }

func (a *Account) Withdraw(amount int64) error {
	if amount > a.balance {
		return ErrOverdraft
	}
	a.balance -= amount
	return nil
}

func (a *Account) Balance() int64 { return a.balance }

// region:opaque:end
