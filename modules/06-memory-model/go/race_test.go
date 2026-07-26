//go:build racedemo

// This file is excluded from the normal suite (it deliberately races, which
// would fail `make test-race`). The capture tool builds it with
// `-tags racedemo -race` to inject the detector's report into the README.
package conc

import "testing"

func TestDataRace(t *testing.T) {
	_ = RacyInc()
}
