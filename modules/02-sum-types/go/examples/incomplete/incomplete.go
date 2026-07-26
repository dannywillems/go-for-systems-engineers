// Package incomplete contains a deliberately non-exhaustive type switch over a
// sealed interface. It COMPILES (that is the whole problem); the capture tool
// runs the exhaustive analyzer on it to show the diagnostic Go itself omits.
package incomplete

//sumtype:decl
type Color interface{ color() }

type Red struct{}
type Green struct{}
type Blue struct{}

func (Red) color()   {}
func (Green) color() {}
func (Blue) color()  {}

// Name is missing the Blue case and has no default: valid Go, wrong program.
func Name(c Color) string {
	switch c.(type) {
	case Red:
		return "red"
	case Green:
		return "green"
	}
	return "?"
}
