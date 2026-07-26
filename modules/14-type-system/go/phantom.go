// Package ts collects type-system quirks: what Go's type system can and cannot
// express, next to the same idea in Rust, OCaml, Swift, and Kotlin.
//
// Go's system is deliberately small: nominal named types, structural
// interfaces, parametric generics (Module 03). It has NO higher-kinded types,
// NO generic methods (Module 04), NO GADTs, and NO declaration-site variance.
// But it can do more than most Go code uses -- phantom type parameters, below.
package ts

import "fmt"

// region:phantom:start

// ID[Tag] is a PHANTOM-typed identifier: the Tag type parameter appears only in
// the type, never in a field, so it costs nothing at runtime (an int64) yet
// makes ID[User] and ID[Order] distinct types the compiler refuses to mix.
type ID[Tag any] int64

type User struct{}

type Order struct{}

type UserID = ID[User]

type OrderID = ID[Order]

// LookupUser accepts only a UserID; passing an OrderID is a compile error, even
// though both are int64. See reject-go for the rejection.
func LookupUser(id UserID) string {
	return fmt.Sprintf("user #%d", int64(id))
}

// region:phantom:end
