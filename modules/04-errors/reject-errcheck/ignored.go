// Package ignored COMPILES: dropping an error return value is legal Go, which is
// the whole gap. Go has no #[must_use] on error, so `os.Remove(path)` silently
// discards its error. errcheck (bundled in golangci-lint) is the external tool
// that reintroduces the obligation; capture runs it here to show the finding.
package ignored

import "os"

// Cleanup ignores the error from os.Remove. This compiles cleanly.
func Cleanup(path string) {
	os.Remove(path)
}
