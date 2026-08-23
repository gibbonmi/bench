package shift

import (
	"strings"
	"testing"
)

// TestObjectiveCommitSubjectSharesBannerPolicy pins the one behavior change of the
// objective owner. The durable commit subject carries the sanitizer's escaped, bounded
// preview, the banner's policy, instead of the verbatim objective. The commit subject
// never carries the raw objective.
func TestObjectiveCommitSubjectSharesBannerPolicy(t *testing.T) {
	esc := string(rune(0x1b))
	long := objective(strings.Repeat("x", 150))
	if got := long.commitSubject(3); !strings.HasPrefix(got, "shift: iteration 3 — "+strings.Repeat("x", 120)+"… (150 bytes)") || strings.Contains(got, strings.Repeat("x", 121)) {
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
