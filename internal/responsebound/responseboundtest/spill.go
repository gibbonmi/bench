// Package responseboundtest serves consumer tests of the response bound. It reads the
// spill path of an over-bound response, and it makes the assignment checkout whose tree
// sets the spill scope. The spill line has one reader here, so a test in any package
// finds the spill file the same way. The frame owner in package responsebound classifies
// the line, so this package holds no second marker parser.
package responseboundtest

import (
	"testing"

	"github.com/gibbonmi/bench/internal/responsebound"
)

// Spill is the parsed spill line of one over-bound response.
type Spill struct {
	// Path is the absolute path of the spill file. It is the last field, so it can hold a
	// comma.
	Path string
}

// ParseLine parses one spill line, with or without its newline. It projects the frame
// owner's classifier, so a line without the spill frame or a path answers false.
func ParseLine(line string) (Spill, bool) {
	return spillOf(responsebound.ProjectionLine(line))
}

// Find answers the spill line of response when it is the first reserved frame that the
// frame owner finds. A response within the bound has none.
func Find(response string) (Spill, bool) {
	projection, line, err := responsebound.First([]byte(response))
	return spillOf(projection, line > 0, err)
}

// spillOf projects one classified frame onto the spill line, so a frame of another kind, a
// malformed frame, or no frame answers false.
func spillOf(projection responsebound.Projection, ok bool, err error) (Spill, bool) {
	if !ok || err != nil || projection.Kind != responsebound.Spilled {
		return Spill{}, false
	}
	return Spill{Path: projection.Path}, true
}

// Path answers the spill path that response names, and fails the test without one.
func Path(t testing.TB, response string) string {
	t.Helper()
	spill, ok := Find(response)
	if !ok {
		t.Fatalf("response = %q, want a spill line", response)
	}
	return spill.Path
}
