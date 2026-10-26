// The refusal line of a red whole landing run that names reader landings as its suspects.
package landing

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate/authorization"
)

// TestRefusalMessageNamesTheReaderLandings proves that a refusal whose result carries two
// suspects appends both to the current line, and that a result without suspects keeps that
// line.
func TestRefusalMessageNamesTheReaderLandings(t *testing.T) {
	plain := authorization.Result{Kind: authorization.Candidate}
	suspected := plain
	suspected.Suspects = []string{strings.Repeat("a", 40), strings.Repeat("b", 40)}
	want := refusalMessage(plain) + "; reader landings since the green marker: " + suspected.Suspects[0] + ", " + suspected.Suspects[1]
	if got := refusalMessage(suspected); got != want {
		t.Fatalf("refusalMessage = %q, want %q", got, want)
	}
}
