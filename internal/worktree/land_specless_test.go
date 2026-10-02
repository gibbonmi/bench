// Spec-less landing tests: publication, gate refusal, base resolution, and resume without a spec.
package worktree

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// WL1, WL2, WL5, and WL10: the spec-less landing publishes and releases like a spec
// landing, keeps both reviewed parents and the marker, transitions no spec, and prints
// the same record.
func TestLandCommandSpecLessLandsPublishesAndReleases(t *testing.T) {
	t.Parallel()
	request := "spec-less-land"
	f := specLessLandingFixture(t, request)
	specsBefore := gitOutput(t, f.creation.Path, "rev-parse", f.tip+":specs")
	r := runVerb(t, verbLand, f.call(specLessLandArgs(request, f.base, f.tip, f.creation.Path)...))
	published := gitOutput(t, f.root, "rev-parse", "main")
	tree := gitOutput(t, f.root, "rev-parse", published+"^{tree}")
	want := wantEffects("not-applicable") + "landed{source_base=" + f.base + ",source_tip=" + f.tip + ",destination_base=" + f.base + ",published_commit=" + published + ",tree=" + tree + ",worktree=released,census=0}\n"
	if r.exit != 0 || r.stdout != want {
		t.Fatalf("spec-less land = (%d, %q, %q), want (0, %q)", r.exit, r.stdout, r.stderr, want)
	}
	parents := strings.Fields(gitOutput(t, f.root, "rev-list", "--parents", "-n", "1", published))
	if len(parents) != 3 || parents[1] != f.base || parents[2] != f.tip {
		t.Fatalf("published parents = %q, want destination %s and source %s", parents, f.base, f.tip)
	}
	if got := gitOutput(t, f.root, "rev-parse", "refs/bench/green/main"); got != published {
		t.Fatalf("project-green = %s, want %s", got, published)
	}
	if got := gitOutput(t, f.root, "rev-parse", published+":specs"); got != specsBefore {
		t.Fatalf("published specs tree = %s, want the composed %s", got, specsBefore)
	}
	if got := gitOutput(t, f.root, "show", published+":specs/x/spec.md"); strings.Contains(got, "Status: implemented") {
		t.Fatalf("spec-less landing transitioned a spec: %q", got)
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("gate tally = %q, %v", got, err)
	}
	if _, err := os.Stat(f.creation.Path); !os.IsNotExist(err) {
		t.Fatalf("spec-less landing retained the worktree: %v", err)
	}
}

// WL3 and WL22: the gate still owns the spec-less path, and its refusal exits 1 with
// nothing published.
func TestLandCommandSpecLessGateRefusalPublishesNothing(t *testing.T) {
	t.Parallel()
	request := "spec-less-gate-red"
	f := specLessLandingFixture(t, request)
	mustWrite(t, filepath.Join(f.root, ".bench", "gate-prospective.sh"), []byte("#!/bin/sh\nset -eu\nruntime=$1\nprintf g >> '"+f.tally+"'\nexit 1\n"), 0o755)
	gitRun(t, f.root, "add", ".bench/gate-prospective.sh")
	gitRun(t, f.root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "red prospective gate")
	base := gitOutput(t, f.root, "rev-parse", "HEAD")
	gitRun(t, f.creation.Path, "rebase", "main")
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")

	r := runVerb(t, verbLand, f.call(specLessLandArgs(request, base, tip, f.creation.Path)...))
	if r.exit != 1 || !strings.HasPrefix(r.stdout, "refused{detail=prospective authorization refused") {
		t.Fatalf("spec-less gate refusal = (%d, %q, %q), want exit 1 and an authorization refusal", r.exit, r.stdout, r.stderr)
	}
	if got := gitOutput(t, f.root, "rev-parse", "main"); got != base {
		t.Fatalf("gate refusal published main=%s, want %s", got, base)
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("gate tally = %q, %v", got, err)
	}
}

// WL6: the identity proofs still run without a spec, and the refusal names both sides.
func TestLandCommandSpecLessRefusesSourceTipMismatch(t *testing.T) {
	t.Parallel()
	request := "spec-less-tip-mismatch"
	f := specLessLandingFixture(t, request)
	r := runVerb(t, verbLand, f.call(specLessLandArgs(request, f.base, f.base, f.creation.Path)...))
	// The route re-points the caller's own command at the tip the worktree holds, and a
	// spec-less landing re-runs spec-less.
	want := "refused{detail=worktree source tip mismatch,observed=" + f.base + ",wanted=" + f.tip +
		",next=" + landingRerun(request, f.base, f.tip, "", f.creation.Path, f.creation.Assignment.ID) + "}\n"
	if r.exit != 1 || r.stdout != want || len(r.stderr) != 0 {
		t.Fatalf("spec-less tip mismatch = (%d, %q, %q), want (1, %q, empty)", r.exit, r.stdout, r.stderr, want)
	}
	if _, err := os.Stat(f.tally); !os.IsNotExist(err) {
		t.Fatalf("identity refusal ran the gate: %v", err)
	}
}

// WL23: the range proof still runs without a spec, and it refuses before the gate.
func TestLandCommandSpecLessRefusesNonAncestorBaseBeforeTheGate(t *testing.T) {
	t.Parallel()
	request := "spec-less-nonancestor-base"
	f := specLessLandingFixture(t, request)
	commitInWorktree(t, f.root, "destination-only", "destination\n", "destination movement")
	unrelated := gitOutput(t, f.root, "rev-parse", "HEAD")
	r := runVerb(t, verbLand, f.call(specLessLandArgs(request, unrelated, f.tip, f.creation.Path)...))
	if r.exit != 1 || !strings.Contains(r.stdout, "reviewed source range is invalid") || !strings.Contains(r.stdout, "not an ancestor") {
		t.Fatalf("spec-less non-ancestor base = (%d, %q, %q), want an invalid-range refusal", r.exit, r.stdout, r.stderr)
	}
	if _, err := os.Stat(f.tally); !os.IsNotExist(err) {
		t.Fatalf("range refusal ran the gate: %v", err)
	}
}

// WL24: the landed record carries the resolved review base, never the flag's spelling.
func TestLandCommandSpecLessLandedSourceBaseIsTheResolvedBase(t *testing.T) {
	t.Parallel()
	request := "spec-less-resolved-base"
	f := specLessLandingFixture(t, request)
	r := runVerb(t, verbLand, f.call(specLessLandArgs(request, f.base[:12], f.tip, f.creation.Path)...))
	if r.exit != 0 || !strings.Contains(r.stdout, "landed{source_base="+f.base+",") || strings.Contains(r.stdout, f.base[:12]+",") {
		t.Fatalf("abbreviated spec-less base = (%d, %q, %q), want the resolved base in the record", r.exit, r.stdout, r.stderr)
	}
}

// WL7: a spec-less landing interrupted after publication resumes without a spec.
func TestResumeLandCommandSpecLessCompletesAnInterruptedLanding(t *testing.T) {
	t.Parallel()
	request := "spec-less-resume"
	f := specLessLandingFixture(t, request)
	working := defaultJoins()
	broken := working
	broken.advanceLandingMarker = func(context.Context, string, string, string, string) error {
		return errors.New("injected marker interruption")
	}
	r := runVerb(t, verbLand, f.callWith(broken, specLessLandArgs(request, f.base, f.tip, f.creation.Path)...))
	if r.exit != 3 || !strings.Contains(r.stdout, "worktree=incomplete:marker") {
		t.Fatalf("interrupted spec-less landing = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	if strings.Contains(r.stdout, "--spec") {
		t.Fatalf("spec-less resume instruction named a spec: %q", r.stdout)
	}
	published := gitOutput(t, f.root, "rev-parse", "main")
	args := []string{"--resume", published, "--request", request, "--base", f.base, "--source-tip", f.tip, f.creation.Path}
	r = runVerb(t, verbLand, f.callWith(working, args...))
	if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released,census=0}") || len(r.stderr) != 0 {
		t.Fatalf("spec-less resume = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	if got := gitOutput(t, f.root, "rev-parse", "refs/bench/green/main"); got != published {
		t.Fatalf("project-green = %s, want %s", got, published)
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("spec-less resume reran the gate: tally=%q error=%v", got, err)
	}
}

// WL25: a resume without a spec completes a published spec-backed landing's marker and
// release, and publishes nothing a second time.
func TestResumeLandCommandWithoutSpecCompletesASpecBackedLanding(t *testing.T) {
	t.Parallel()
	request := "spec-backed-spec-less-resume"
	f := publicLandingFixture(t, request, "", "")
	working := defaultJoins()
	broken := working
	broken.advanceLandingMarker = func(context.Context, string, string, string, string) error {
		return errors.New("injected marker interruption")
	}
	r := runVerb(t, verbLand, f.callWith(broken, landArgs(request, f.base, f.tip, f.creation.Path)...))
	if r.exit != 3 || !strings.Contains(r.stdout, "worktree=incomplete:marker") {
		t.Fatalf("interrupted landing = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	published := gitOutput(t, f.root, "rev-parse", "main")
	args := []string{"--resume", published, "--request", request, "--base", f.base, "--source-tip", f.tip, f.creation.Path}
	r = runVerb(t, verbLand, f.callWith(working, args...))
	if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released,census=0}") || len(r.stderr) != 0 {
		t.Fatalf("spec-less resume of a spec landing = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	if got := gitOutput(t, f.root, "rev-parse", "main"); got != published {
		t.Fatalf("resume republished: main=%s, want %s", got, published)
	}
	if got := gitOutput(t, f.root, "rev-parse", "refs/bench/green/main"); got != published {
		t.Fatalf("project-green = %s, want %s", got, published)
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("resume reran the gate: tally=%q error=%v", got, err)
	}
}

// The edge under WL1: the flag is optional, but an empty value stays a usage error.
func TestLandCommandRefusesAnEmptySpecValue(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"first run", []string{"--request", "r", "--base", "b", "--source-tip", "s", "--spec", "", "-m", "m", "path"}},
		{"resume", []string{"--resume", "p", "--request", "r", "--base", "b", "--source-tip", "s", "--spec", "", "path"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := runVerb(t, verbLand, verbCall{home: Home(), args: tc.args})
			if r.exit != 2 || len(r.stdout) != 0 || !strings.Contains(r.stderr, `--spec ""`) {
				t.Fatalf("empty --spec = (%d, %q, %q), want (2, empty, a usage line naming the empty value)", r.exit, r.stdout, r.stderr)
			}
		})
	}
}
