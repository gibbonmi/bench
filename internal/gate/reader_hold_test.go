package gate

// The reader hold. Each test builds the landing fixture of the landing plan. A red whole
// landing run commits an uncovered change after the destination, so the plan refuses the
// reader set, and plants the fixture gate's red.

import (
	"bytes"
	"context"
	"slices"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// hold writes the reader hold of commit in the store of the fixture.
func (f *landingFixture) hold(commit string) {
	f.t.Helper()
	gitdir, err := benchgit.CommonDir(f.root)
	if err != nil {
		f.t.Fatal(err)
	}
	if err := writeReaderHold(durablefile.Native(), gitdir, commit); err != nil {
		f.t.Fatal(err)
	}
}

// redWholeLanding runs the planted red whole gate over an uncovered change after the
// destination, under the fixture's landing value, and answers its result.
func (f *landingFixture) redWholeLanding() Result {
	f.t.Helper()
	f.put("AGENTS.md", "change\n")
	f.commit("change")
	outcomeWrite(f.t, f.root, ".gate-red", "\n", 0o644)
	owner := func(context.Context, string) (*runbinary.Selection, error) { return fixtureSelection(f.root), nil }
	var stdout, stderr bytes.Buffer
	result := executeSubjectReplacing(f.landing(), f.root, f.root, &stdout, &stderr, nil, forceRun, newGateEvaluation(f.root), owner, "", durablefile.Native())
	if result.ActionExit != 7 || result.Schedule != ScheduleWhole {
		f.t.Fatalf("landing = %+v\nstdout:\n%s\nstderr:\n%s\nwant the planted red of the whole gate", result, stdout.String(), stderr.String())
	}
	return result
}

// TestRedWholeLandingWritesTheReaderHold proves that a red whole run with a landing value
// writes a reader hold that names the run's green marker.
func TestRedWholeLandingWritesTheReaderHold(t *testing.T) {
	f := newLandingFixture(t, nil)
	f.redWholeLanding()
	if held, err := readerHeld(f.root, f.marker); err != nil || !held {
		t.Fatalf("reader hold of %s after a red whole landing = %t %v; want the hold", f.marker, held, err)
	}
}

// TestReaderPlanRunsWholeUnderAHold proves that a hold that names the current green marker
// sends a covered landing to the whole gate.
func TestReaderPlanRunsWholeUnderAHold(t *testing.T) {
	f := newLandingFixture(t, nil)
	f.hold(f.marker)
	f.refuses(docsChange, nil, refusedHold)
}

// TestReaderPlanIgnoresAHoldOfAnOlderMarker proves that a hold that names an older green
// marker leaves the reader set selectable.
func TestReaderPlanIgnoresAHoldOfAnOlderMarker(t *testing.T) {
	f := newLandingFixture(t, nil)
	f.hold(gittest.Output(t, f.root, "rev-parse", f.marker+"~1"))
	f.selects(docsChange)
}

// TestRedWholeLandingListsTheSuspects proves that a red whole landing run after two
// reader-record commits lists both commits as suspects, and not the commit of the chain that
// holds no reader record.
func TestRedWholeLandingListsTheSuspects(t *testing.T) {
	f := newLandingFixture(t, nil)
	f.chain(3, 0)
	want := []string{gittest.Output(t, f.root, "rev-parse", f.destination+"~1"), f.destination}
	if got := f.redWholeLanding().Suspects; !slices.Equal(got, want) {
		t.Fatalf("suspects = %q; want the reader landings %q", got, want)
	}
}
