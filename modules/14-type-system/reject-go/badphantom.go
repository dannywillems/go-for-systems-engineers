// Package badphantom DOES NOT COMPILE: it passes an OrderID where a UserID is
// required. The phantom Tag makes them distinct types even though both are
// int64, so the mix is a compile error -- the whole point of phantom typing.
package badphantom

type ID[Tag any] int64

type User struct{}

type Order struct{}

type UserID = ID[User]

type OrderID = ID[Order]

func LookupUser(id UserID) int64 { return int64(id) }

var _ = LookupUser(OrderID(7))
