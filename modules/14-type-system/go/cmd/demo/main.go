package main

import (
	"fmt"

	ts "github.com/dannywillems/go-for-systems-engineers/modules/14-type-system/go"
)

func main() {
	u := ts.UserID(7)
	o := ts.OrderID(7)
	fmt.Println(ts.LookupUser(u))
	// LookupUser(o) would not compile: OrderID is not a UserID. We can still
	// see they are the same underlying value if we ask explicitly:
	fmt.Printf("UserID(7) and OrderID(7) share underlying int64: %v\n",
		int64(u) == int64(o))
}
