// The landing's own trace: the gate the authorization runs, and the phases below it,
// sit under the landing span rather than starting a trace of their own.
package worktree

import (
	"bytes"
	"testing"

	"github.com/gibbonmi/bench/internal/otelrecord"
)

// LC39: the lane's phase spans carry the landing span's trace identity. The phases
// inherit the gate run's context, so the authorization's own gate span is the parent
// that decides their trace, and this fixture's shell gate opens no phase of its own.
//
// The landing records under the home this test owns, and the gate resolves the package's
// own process home for itself, so the two spans of one trace are read from two records.
// One repository has one record per home, and this repository is this test's alone.
func TestLandRecordsOneTraceForThePhases(t *testing.T) {
	t.Parallel()
	request := "land-one-trace"
	root, creation, base, tip, _, home := specLessLandingFixture(t, request)

	var stdout, stderr bytes.Buffer
	if code := LandCommand(root, home, "", specLessLandArgs(request, base, tip, creation.Path), &stdout, &stderr); code != 0 {
		t.Fatalf("land = %d, want 0: %q %q", code, stdout.String(), stderr.String())
	}

	landing, ok := otelrecord.NewestLanding(home, root)
	if !ok {
		t.Fatal("the landing's record names no completed landing")
	}
	authorized, err := otelrecord.ReadSpans(privateBenchHome, root)
	if err != nil {
		t.Fatalf("read the authorization's record: %v", err)
	}
	var below int
	for _, span := range authorized {
		if span.Seam != "gate" && span.Seam != otelrecord.SeamGatePhase {
			continue
		}
		below++
		if span.TraceID != landing.TraceID || span.ParentSpanID == "" {
			t.Fatalf("%s %q trace = %s/%s, want the landing trace %s below one of its spans",
				span.Seam, span.Name, span.TraceID, span.ParentSpanID, landing.TraceID)
		}
	}
	if below == 0 {
		t.Fatalf("the authorization's record names no gate: %+v", authorized)
	}
}
