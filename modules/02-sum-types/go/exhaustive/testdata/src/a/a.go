package a

//sumtype:decl
type Shape interface{ shape() }

type Circle struct{}
type Square struct{}
type Triangle struct{}

func (Circle) shape()   {}
func (Square) shape()   {}
func (Triangle) shape() {}

func missing(s Shape) int {
	switch s.(type) { // want `non-exhaustive type switch on Shape: missing cases Square, Triangle`
	case Circle:
		return 1
	}
	return 0
}

func withDefault(s Shape) int {
	switch s.(type) { // ok: a default clause covers the rest
	case Circle:
		return 1
	default:
		return 0
	}
}

func exhaustive(s Shape) int {
	switch s.(type) { // ok: every variant handled
	case Circle:
		return 1
	case Square:
		return 2
	case Triangle:
		return 3
	}
	return 0
}
