package shift

import (
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
)

// TestObjectiveCommitSubjectSharesBannerPolicy pins the one behavior change of the
// objective owner. The durable commit subject carries the sanitizer's escaped, bounded
// preview, the banner's policy, instead of the verbatim objective. The commit subject
// never carries the raw objective.
func TestObjectiveCommitSubjectSharesBannerPolicy(t *testing.T) {
	esc := string(rune(0x1b))
	// The fixture overruns the preview cap so the subject must truncate. "x" is one byte
	// per rune, so the rune count is also the byte count the suffix names.
	const overrun = bounds.PreviewRuneLimit + 30
	long := objective(strings.Repeat("x", overrun))
	want := "shift: iteration 3 — " + strings.Repeat("x", bounds.PreviewRuneLimit) + "… (" + strconv.Itoa(overrun) + " bytes)"
	if got := long.commitSubject(3); !strings.HasPrefix(got, want) || strings.Contains(got, strings.Repeat("x", bounds.PreviewRuneLimit+1)) {
		t.Fatalf("commit subject is not the bounded preview: %q", got)
	}
	hostile := objective("paint it " + esc + "[31mred")
	if got := hostile.commitSubject(1); strings.ContainsRune(got, 0x1b) || !strings.Contains(got, `[31mred`) {
		t.Fatalf("commit subject leaked or dropped the control sequence: %q", got)
	}
	if b, s := hostile.banner("bench/shift-x"), hostile.commitSubject(1); !strings.HasSuffix(b, strings.TrimPrefix(s, "shift: iteration 1 — ")) {
		t.Fatalf("banner %q and commit subject %q render the objective differently", b, s)
	}
}

// TestObjectiveVerbatimProjections pins the surfaces that must not change: prompt,
// predicate argument, and scratch bytes carry the objective exactly.
func TestObjectiveVerbatimProjections(t *testing.T) {
	o := objective("keep the fixtures honest — même en UTF-8")
	if got := o.predicateArgument(); got != string(o) {
		t.Fatalf("predicate argument = %q, want verbatim", got)
	}
	if got := string(o.scratch()); got != string(o)+"\n" {
		t.Fatalf("scratch = %q, want verbatim plus newline", got)
	}
	if got := o.prompt(); !strings.Contains(got, "Objective: "+string(o)) {
		t.Fatalf("prompt does not carry the verbatim objective: %q", got)
	}
}
