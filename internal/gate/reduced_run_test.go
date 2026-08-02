package gate

// The reduced execution path: a changeset confined to the declared allowlist runs only
// the phases that can observe it and inherits a full-green ancestor's evidence for the
// rest. The failure class this file guards is a verdict that credits work nobody
// graded, so every test observes execution through durable markers — which gate ran,
// which phases ran, what the slot records — rather than through return values alone.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/canary"
)

// reducedRunFixture builds a repository observable at both execution layers: the
// resolved gate script appends to .git/full-runs, and each declared phase appends its
// name to .git/phase-runs. A full run therefore leaves a gate marker and no phase
// markers; a reduced run leaves phase markers and no gate marker — the two cannot be
// confused. The phase names carry the declaration's meaning, so the fixture guards
// both memberships before relying on them.
func reducedRunFixture(t *testing.T) string {
	t.Helper()
	scope := ReducedScope()
	if !scope.Excludable(canary.PhaseTest) || scope.Excludable(conformancePhaseName) {
		t.Fatal("fixture phase names no longer match the declaration; repoint the fixture")
	}
	if !scope.Member("ROADMAP.md") || !scope.Member("capture/learnings.md") {
		t.Fatal("fixture capture paths are no longer declared; repoint the fixture")
	}
	root := t.TempDir()
	// Reduction is only offered to the kit's own root — the tree that declares the
	// scope being applied — and the wrapper names that root through BENCH_KIT. The
	// fixture claims that identity for itself, exactly as the contract fixtures do,
	// so these tests hold whether or not an enclosing gate run exported the real
	// kit's BENCH_KIT into the test environment.
	t.Setenv("BENCH_KIT", root)
	gitRun(t, root, "init", "-q")
	// The gate script routes through the gate-phases hand-off exactly as the kit's own
	// entry does — a stand-in binary keeps the exec inert, and the marker line above it
	// is the layer this fixture observes. Routing is what makes the root reducible:
	// phaseTableGate proves the resolved gate execs the phase table before any
	// reduction, manifest or not.
	writeGateTestFile(t, root, ".bench/gate.sh",
		"#!/usr/bin/env bash\necho full >> .git/full-runs\nexec true gate-phases \"$PWD\"\n", 0o755)
	writeGateTestFile(t, root, ".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`, 0o644)
	writeGateTestFile(t, root, ".bench/phase-conformance.sh", "echo conformance >> .git/phase-runs\n", 0o644)
	writeGateTestFile(t, root, ".bench/phase-test.sh", "echo test >> .git/phase-runs\n", 0o644)
	writeGateTestFile(t, root, canary.PhaseManifestPath, `{"phases":[`+
		`{"name":"conformance","argv":["bash",".bench/phase-conformance.sh"]},`+
		`{"name":"test","argv":["bash",".bench/phase-test.sh"]}]}`+"\n", 0o644)
	writeGateTestFile(t, root, "ROADMAP.md", "roadmap\n", 0o644)
	writeGateTestFile(t, root, "capture/learnings.md", "learnings\n", 0o644)
	writeGateTestFile(t, root, "graded.txt", "graded content\n", 0o644)
	return root
}

func markerLines(t *testing.T, root, name string) []string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, ".git", name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	return strings.Fields(string(data))
}

func fullRunCount(t *testing.T, root string) int { return len(markerLines(t, root, "full-runs")) }

func phaseRunNames(t *testing.T, root string) []string { return markerLines(t, root, "phase-runs") }

func slotRecord(t *testing.T, root string, now time.Time) verdictRecord {
	t.Helper()
	loaded := loadVerdict(cachePath(t, root), now)
	if loaded.state != Ready {
		t.Fatalf("slot record = %s/%q, want a ready verdict", loaded.state, loaded.reason)
	}
	return loaded.record
}

func mustStrippedSubject(t *testing.T, root string) subject {
	t.Helper()
	s, err := buildStrippedSubject(root)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func commonGitDirOf(t *testing.T, root string) string {
	t.Helper()
	return gitOutput(t, root, "rev-parse", "--path-format=absolute", "--git-common-dir")
}

func mustExecuteGreen(t *testing.T, root string, engine gateEngine) {
	t.Helper()
	if got := executeWithEngine(context.Background(), root, io.Discard, io.Discard, engine); got.ActionExit != 0 {
		t.Fatalf("execution = %+v, want green", got)
	}
}

// [R14] A capture-only changeset executes the included phases only, and the recorded
// phase list says so. Today every change runs all phases through the resolved gate, so
// the recorded phase list is complete — the reduction has to be observable both in the
// record and in which layer actually executed.
func TestReducedRunExecutesIncludedPhasesOnly(t *testing.T) {
	root := reducedRunFixture(t)
	mustExecuteGreen(t, root, productionGateEngine{})
	if got := fullRunCount(t, root); got != 1 {
		t.Fatalf("gate runs after the seed = %d, want 1", got)
	}
	if got := phaseRunNames(t, root); len(got) != 0 {
		t.Fatalf("seed run reached the phase table directly: %v", got)
	}

	writeGateTestFile(t, root, "ROADMAP.md", "capture-only edit\n", 0o644)
	var stdout bytes.Buffer
	if got := executeWithEngine(context.Background(), root, &stdout, io.Discard, productionGateEngine{}); got.ActionExit != 0 {
		t.Fatalf("capture-only execution = %+v, want green", got)
	}
	rec := slotRecord(t, root, time.Now().UTC())
	if !rec.Reduced || !reflect.DeepEqual(rec.Phases, []string{conformancePhaseName}) {
		t.Fatalf("recorded verdict = %+v, want a reduced record whose phase list is exactly the included set", rec)
	}
	if got := phaseRunNames(t, root); !reflect.DeepEqual(got, []string{conformancePhaseName}) {
		t.Fatalf("executed phases = %v, want only %q", got, conformancePhaseName)
	}
	if got := fullRunCount(t, root); got != 1 {
		t.Fatalf("gate runs after the reduced run = %d, want 1 — the reduced run paid the resolved gate", got)
	}
	if sp := mustStrippedSubject(t, root); rec.Ancestor != sp.Tree {
		t.Fatalf("recorded ancestor = %q, want the stripped identity %q the full green answered for", rec.Ancestor, sp.Tree)
	}
	if !strings.Contains(stdout.String(), "gate: reduced run") {
		t.Fatalf("reduced run said nothing about itself:\n%s", stdout.String())
	}
}

// [R17] Consecutive capture commits inherit from the same full green, never from each
// other. An implementation reading the ancestor from the previous record's own identity
// finds a reduced verdict there and falls back to a full run — the motivating scenario,
// unserved — so the second run has to stay reduced and carry the first run's ancestor
// identity and time byte-for-byte.
func TestConsecutiveReducedRunsShareAncestor(t *testing.T) {
	root := reducedRunFixture(t)
	mustExecuteGreen(t, root, productionGateEngine{})

	writeGateTestFile(t, root, "capture/learnings.md", "first capture edit\n", 0o644)
	mustExecuteGreen(t, root, productionGateEngine{})
	first := slotRecord(t, root, time.Now().UTC())
	if !first.Reduced {
		t.Fatalf("first capture commit = %+v, want a reduced verdict", first)
	}

	writeGateTestFile(t, root, "capture/learnings.md", "second capture edit\n", 0o644)
	mustExecuteGreen(t, root, productionGateEngine{})
	second := slotRecord(t, root, time.Now().UTC())
	if !second.Reduced {
		t.Fatalf("second consecutive capture commit fell back to a full run: %+v", second)
	}
	if second.Ancestor != first.Ancestor || second.AncestorRecordedAt != first.AncestorRecordedAt {
		t.Fatalf("second run inherited (%s at %s), first inherited (%s at %s); want the same full-green ancestor",
			second.Ancestor, second.AncestorRecordedAt, first.Ancestor, first.AncestorRecordedAt)
	}
	if got := fullRunCount(t, root); got != 1 {
		t.Fatalf("gate runs across the drain = %d, want 1 — a capture commit paid the full gate", got)
	}
}

// [R18] A reduced verdict is never accepted as an ancestor, so no chain of reduced
// records can form: chaining permitted, a fixture chain produces evidence attributed to
// a run that graded nothing. Both halves are pinned — a planted reduced record at the
// ancestor slot forces a full run, and a legitimate reduced run authors no ancestor of
// its own for a later run to chain from.
func TestReducedVerdictIsNotAnAncestor(t *testing.T) {
	t.Run("a reduced record at the ancestor slot is refused", func(t *testing.T) {
		root := reducedRunFixture(t)
		sp := mustStrippedSubject(t, root)
		now := time.Now().UTC().Truncate(time.Second)
		forged := verdictRecord{
			Schema: 1, State: Ready, Status: "green", Tree: sp.Tree, Oracle: sp.Oracle,
			RecordedAt: now.Add(-time.Minute).Format(time.RFC3339),
			Reduced:    true, Phases: []string{conformancePhaseName},
			Ancestor:           strings.Repeat("d", 40),
			AncestorRecordedAt: now.Add(-2 * time.Minute).Format(time.RFC3339),
		}
		gitdir := commonGitDirOf(t, root)
		dir := filepath.Join(gitdir, "bench-gate-evidence")
		if err := ensureEvidenceDir(gitdir, dir); err != nil {
			t.Fatal(err)
		}
		if err := durableReplaceAt(dir, evidenceName(sp), forged); err != nil {
			t.Fatal(err)
		}
		// The forgery is exactly the shape a marker-blind reuse would credit; if this
		// sanity check fails the case below proves nothing.
		if got := inspectEvidence(root, sp, now); !got.ReusableGreen {
			t.Fatalf("forged ancestor is invisible to the fixture: %+v", got)
		}

		mustExecuteGreen(t, root, productionGateEngine{})
		if rec := slotRecord(t, root, time.Now().UTC()); rec.Reduced {
			t.Fatalf("gate inherited from a reduced record: %+v — evidence attributed to a run that graded nothing", rec)
		}
		if got := fullRunCount(t, root); got != 1 {
			t.Fatalf("gate runs = %d, want 1 full run past the forged ancestor", got)
		}
	})

	t.Run("a reduced run authors no ancestor", func(t *testing.T) {
		root := reducedRunFixture(t)
		mustExecuteGreen(t, root, productionGateEngine{})
		gitdir := commonGitDirOf(t, root)
		sp := mustStrippedSubject(t, root)
		before := mustRead(t, evidencePath(gitdir, sp))
		if got := evidenceFiles(t, gitdir); len(got) != 2 {
			t.Fatalf("retained evidence after the full green = %v, want the whole-tree and stripped records", got)
		}

		writeGateTestFile(t, root, "ROADMAP.md", "capture-only edit\n", 0o644)
		mustExecuteGreen(t, root, productionGateEngine{})
		if rec := slotRecord(t, root, time.Now().UTC()); !rec.Reduced {
			t.Fatalf("capture-only execution = %+v, want a reduced verdict", rec)
		}
		if after := mustRead(t, evidencePath(gitdir, sp)); !bytes.Equal(before, after) {
			t.Fatalf("the reduced run re-authored the ancestor:\nbefore %q\nafter  %q", before, after)
		}
		if got := evidenceFiles(t, gitdir); len(got) != 2 {
			t.Fatalf("retained evidence after the reduced run = %v, want it unchanged — a reduced green retains nothing", got)
		}
	})
}

// [R19] The ancestor lookup is content-addressed, not clock-bounded: an ancestor past
// the whole-tree freshness window still serves an allowlist-confined changeset, because
// every phase that can observe the change runs fresh and the inherited phases answer for
// content the stripped identity proves unchanged. The ancestor's own recorded time is
// carried verbatim into each reduced record — never re-stamped — so the record always
// attributes the inherited evidence to the run that produced it.
func TestOldAncestorStillServesReducedRun(t *testing.T) {
	root := reducedRunFixture(t)
	t0 := time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)
	mustExecuteGreen(t, root, &faultEngine{now: t0})

	writeGateTestFile(t, root, "ROADMAP.md", "first capture edit\n", 0o644)
	within := t0.Add(30 * time.Minute)
	mustExecuteGreen(t, root, &faultEngine{now: within})
	rec := slotRecord(t, root, within)
	if !rec.Reduced || rec.AncestorRecordedAt != t0.Format(time.RFC3339) {
		t.Fatalf("in-window record = %+v, want a reduced verdict carrying the ancestor's time %s", rec, t0.Format(time.RFC3339))
	}

	writeGateTestFile(t, root, "ROADMAP.md", "second capture edit\n", 0o644)
	old := t0.Add(freshness + time.Minute)
	mustExecuteGreen(t, root, &faultEngine{now: old})
	rec = slotRecord(t, root, old)
	if !rec.Reduced || rec.AncestorRecordedAt != t0.Format(time.RFC3339) {
		t.Fatalf("old-ancestor record = %+v, want a reduced verdict still carrying the ancestor's own time %s", rec, t0.Format(time.RFC3339))
	}
	if got := fullRunCount(t, root); got != 1 {
		t.Fatalf("gate runs = %d, want 1 — the old ancestor should keep serving reduced runs", got)
	}
}

// [R20] An allowlist-confined change with no ancestor runs the full gate: a first
// commit, a fresh clone, or a pruned cache has nothing sound to inherit, and treating
// the absence as reusable would emit a reduced verdict with an empty ancestor field.
func TestNoAncestorFallsBackToFullRun(t *testing.T) {
	t.Run("first run with no evidence at all", func(t *testing.T) {
		root := reducedRunFixture(t)
		mustExecuteGreen(t, root, productionGateEngine{})
		rec := slotRecord(t, root, time.Now().UTC())
		if rec.Reduced || rec.Ancestor != "" {
			t.Fatalf("first-run record = %+v, want a full verdict with no ancestor field", rec)
		}
		if got := fullRunCount(t, root); got != 1 {
			t.Fatalf("gate runs = %d, want 1 full run", got)
		}
	})

	t.Run("pruned evidence cache", func(t *testing.T) {
		root := reducedRunFixture(t)
		mustExecuteGreen(t, root, productionGateEngine{})
		if err := os.RemoveAll(filepath.Join(commonGitDirOf(t, root), "bench-gate-evidence")); err != nil {
			t.Fatal(err)
		}
		writeGateTestFile(t, root, "ROADMAP.md", "capture-only edit\n", 0o644)
		mustExecuteGreen(t, root, productionGateEngine{})
		if rec := slotRecord(t, root, time.Now().UTC()); rec.Reduced {
			t.Fatalf("pruned-cache record = %+v — a missing ancestor was treated as reusable", rec)
		}
		if got := fullRunCount(t, root); got != 2 {
			t.Fatalf("gate runs = %d, want 2 — the pruned cache did not force a full run", got)
		}
	})
}

// [RB1] A root that is not the kit runs unreduced. ReducedScope() is the kit's own
// declaration, compiled into every binary the kit ships, so a linked repo gating through
// that binary — BENCH_KIT naming the kit checkout, the graded root a different tree —
// must pay the full run even when a valid ancestor sits in its cache: the allowlist was
// never its own, and inheriting evidence against it is a scope the repository never
// declared.
func TestForeignRootNeverReduces(t *testing.T) {
	root := reducedRunFixture(t)
	// The graded root stops being the kit: BENCH_KIT now names a different checkout,
	// which is exactly a linked repo's shape.
	t.Setenv("BENCH_KIT", t.TempDir())
	mustExecuteGreen(t, root, productionGateEngine{})

	writeGateTestFile(t, root, "ROADMAP.md", "capture-only edit\n", 0o644)
	mustExecuteGreen(t, root, productionGateEngine{})
	if rec := slotRecord(t, root, time.Now().UTC()); rec.Reduced {
		t.Fatalf("foreign-root record = %+v — the root inherited a reduced scope it never declared", rec)
	}
	if got := phaseRunNames(t, root); len(got) != 0 {
		t.Fatalf("foreign root reached the phase table directly: %v", got)
	}
	if got := fullRunCount(t, root); got != 2 {
		t.Fatalf("gate runs = %d, want 2 — the capture-only edit on a foreign root did not pay a full run", got)
	}
}

// [RB3] The gate-phases hand-off branch of phaseTableGate is the one governing this
// repository's own reductions — the kit declares no phase manifest, so its eligibility
// rests entirely on its gate script carrying the hand-off — yet every reduced fixture
// above writes a phases.json and returns through the manifest branch, so this branch is
// otherwise reached by no test at all. The expectation is authored against the ticket,
// not the implementation; the recorded red is the substring-branch mutation probe.
func TestHandOffBranchAdmitsManifestlessRoots(t *testing.T) {
	// The kit's own shape, minus everything incidental: no .bench/phases.json, and a
	// gate script whose exec line hands off to gate-phases as the real one does.
	routed := "#!/usr/bin/env bash\nexec env BENCH_KIT=\"$kit\" \"$bench\" gate-phases \"$root\"\n"
	handWritten := "#!/usr/bin/env bash\ngo test ./...\n"
	cases := []struct {
		name   string
		script string
		res    Resolution
		want   bool
	}{
		{"manifestless gate script handing off to gate-phases", routed, Resolution{Kind: GateSh}, true},
		{"manifestless hand-written gate script", handWritten, Resolution{Kind: GateSh}, false},
		{"hand-off script but the resolution is not the script", routed, Resolution{Kind: BenchGate, Command: "true"}, false},
		{"missing gate script", "", Resolution{Kind: GateSh}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			if tc.script != "" {
				writeGateTestFile(t, root, ".bench/gate.sh", tc.script, 0o755)
			}
			if got := phaseTableGate(root, tc.res); got != tc.want {
				t.Fatalf("phaseTableGate = %v, want %v", got, tc.want)
			}
		})
	}
}

// [RC6] A reduced verdict past the freshness window has expired like any other verdict:
// reporting it as a current reduced verdict would let the board render the reduced row
// for evidence the window has already retired. The reduced marker still travels with the
// expiry — reducedness and staleness are independent facts, and a consumer needs both.
func TestExpiredReducedVerdictReportsExpiry(t *testing.T) {
	root := gateTestRepo(t, "#!/usr/bin/env bash\nexit 0\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	if got := Execute(context.Background(), root, io.Discard, io.Discard); got.ActionExit != 0 {
		t.Fatalf("seed execution = %+v, want green", got)
	}
	plan := mustSubject(t, root)
	seeded := time.Now().UTC().Truncate(time.Second)
	reduced := reducedTestRecord(seeded)
	reduced.Tree, reduced.Oracle = plan.Tree, plan.Oracle
	if err := durableReplace(filepath.Dir(cachePath(t, root)), reduced); err != nil {
		t.Fatal(err)
	}
	// Inside the window the record reads as what it is; if this half fails, the half
	// below proves nothing about ordering.
	if fresh := inspectAt(root, seeded); !fresh.Reduced || fresh.Reason != "reduced verdict" {
		t.Fatalf("in-window inspection = %+v, want the reduced reason", fresh)
	}
	expired := inspectAt(root, seeded.Add(freshness+time.Minute))
	if expired.State != Ready || !expired.Reduced || expired.ReusableGreen || expired.Reason != "verdict expired" {
		t.Fatalf("past-window inspection = %+v, want an expired ready verdict that keeps the reduced marker", expired)
	}
}

// [RC8] A phase manifest chooses which table a routed gate resolves, never whether the
// gate routes at all: a root whose declared gate script never hands off to gate-phases
// pays the full resolved run even with a manifest beside it, because reducing there
// would swap a hand-written oracle for a phase table it never execs.
func TestManifestWithoutHandOffNeverReduces(t *testing.T) {
	root := reducedRunFixture(t)
	writeGateTestFile(t, root, ".bench/gate.sh",
		"#!/usr/bin/env bash\necho full >> .git/full-runs\nexit 0\n", 0o755)
	mustExecuteGreen(t, root, productionGateEngine{})

	writeGateTestFile(t, root, "ROADMAP.md", "capture-only edit\n", 0o644)
	mustExecuteGreen(t, root, productionGateEngine{})
	if rec := slotRecord(t, root, time.Now().UTC()); rec.Reduced {
		t.Fatalf("non-routing gate reduced: %+v — the manifest chose whether the gate routes, not which table", rec)
	}
	if got := phaseRunNames(t, root); len(got) != 0 {
		t.Fatalf("non-routing gate reached the phase table directly: %v", got)
	}
	if got := fullRunCount(t, root); got != 2 {
		t.Fatalf("gate runs = %d, want 2 — the capture-only edit did not pay the resolved gate", got)
	}
}

// componentSlotFilesIn returns the name of every evidence-store file that decodes as a
// component slot record — the "component" key is the one componentSlotRecord carries and
// no verdict class does, so its presence in the raw bytes is what tells a slot apart from
// the whole-tree and stripped verdicts the same store also holds.
func componentSlotFilesIn(t *testing.T, gitdir string) []string {
	t.Helper()
	var slots []string
	for _, name := range evidenceFiles(t, gitdir) {
		data := mustRead(t, filepath.Join(gitdir, "bench-gate-evidence", name))
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(data, &fields); err != nil {
			t.Fatalf("decode evidence file %s: %v", name, err)
		}
		if _, isSlot := fields["component"]; isSlot {
			slots = append(slots, name)
		}
	}
	return slots
}

// [PC18a] A root that is not the kit executes every component: the whole table runs through
// the resolved gate, no per-component slot is authored, and no skip is announced — the same
// refusal the whole-changeset guard already makes, restated at the per-component site. Red
// mutation: drop the kit-identity guard from scopeComponents so a foreign root scopes
// against declarations that were never its own.
func TestForeignRootScopesNoComponent(t *testing.T) {
	fixture := newKitShapedFixture(t)
	// The graded root stops being the kit: BENCH_KIT now names a different checkout, which
	// is exactly a linked repo's shape.
	t.Setenv("BENCH_KIT", t.TempDir())
	var stdout bytes.Buffer
	if got := executeWithEngine(context.Background(), fixture.root, &stdout, io.Discard, productionGateEngine{}); got.ActionExit != 0 {
		t.Fatalf("foreign-root execution = %+v, want green", got)
	}
	if got := fullRunCount(t, fixture.root); got != 1 {
		t.Fatalf("gate runs = %d, want 1 full run", got)
	}
	executed, want := phaseRunNames(t, fixture.root), fixture.phaseNames()
	sort.Strings(executed)
	sort.Strings(want)
	if !reflect.DeepEqual(executed, want) {
		t.Fatalf("executed phases = %v, want every table phase %v", executed, want)
	}
	gitdir := commonGitDirOf(t, fixture.root)
	if slots := componentSlotFilesIn(t, gitdir); len(slots) != 0 {
		t.Fatalf("foreign root authored component slots: %v", slots)
	}
	if strings.Contains(stdout.String(), "gate: skipping") {
		t.Fatalf("foreign root announced a per-component skip:\n%s", stdout.String())
	}
}

// [PC18b] A root with no Go module executes every component its table carries and scopes
// none: the per-component declarations never reach it, so the whole-changeset guard this
// root already relies on ([R14] et al.) keeps deciding for it undisturbed, and the
// per-component evidence store stays empty throughout. Red mutation: resolve component
// identities before checking for a Go module.
func TestNoGoRootScopesNoComponent(t *testing.T) {
	root := reducedRunFixture(t)
	mustExecuteGreen(t, root, productionGateEngine{})
	if got := fullRunCount(t, root); got != 1 {
		t.Fatalf("gate runs after the seed = %d, want 1", got)
	}
	gitdir := commonGitDirOf(t, root)
	if slots := componentSlotFilesIn(t, gitdir); len(slots) != 0 {
		t.Fatalf("no-Go-module root authored component slots on its seed run: %v", slots)
	}

	writeGateTestFile(t, root, "ROADMAP.md", "capture-only edit\n", 0o644)
	var stdout bytes.Buffer
	if got := executeWithEngine(context.Background(), root, &stdout, io.Discard, productionGateEngine{}); got.ActionExit != 0 {
		t.Fatalf("capture-only execution = %+v, want green", got)
	}
	rec := slotRecord(t, root, time.Now().UTC())
	if !rec.Reduced || !reflect.DeepEqual(rec.Phases, []string{conformancePhaseName}) {
		t.Fatalf("capture-only record = %+v, want the whole-changeset reduction still serving this root", rec)
	}
	if slots := componentSlotFilesIn(t, gitdir); len(slots) != 0 {
		t.Fatalf("no-Go-module root authored component slots: %v", slots)
	}
	if strings.Contains(stdout.String(), "gate: skipping") {
		t.Fatalf("no-Go-module root announced a per-component skip:\n%s", stdout.String())
	}
}

// [PC18c] A symlinked path to the kit root still counts as the kit — a capture-only
// changeset through the symlink narrows exactly as it would through the literal path — and a
// stat failure on either side of the identity check runs every component. Red mutation:
// treat a stat failure as a match.
func TestSymlinkedKitCountsAndStatFailureRunsAll(t *testing.T) {
	fixture := newKitShapedFixture(t)
	// BENCH_KIT is set to the symlink from the seed run onward, never to the literal path:
	// a phase's own argv carries the kit spelling it was resolved with, so comparing a run
	// seeded under one spelling against a narrowed run under another would move every
	// toolchain component's identity on the spelling alone and prove nothing about the
	// symlink admission this case exists to check.
	link := filepath.Join(t.TempDir(), "kit-symlink")
	if err := os.Symlink(fixture.root, link); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_KIT", link)
	mustExecuteGreen(t, fixture.root, productionGateEngine{})
	seeded := phaseRunNames(t, fixture.root)

	for _, path := range captureSurfacePaths(ReducedScope()) {
		writeGateTestFile(t, fixture.root, path, "capture-only edit\n", 0o644)
	}
	var stdout bytes.Buffer
	if got := executeWithEngine(context.Background(), fixture.root, &stdout, io.Discard, productionGateEngine{}); got.ActionExit != 0 {
		t.Fatalf("symlinked narrowed execution = %+v, want green", got)
	}
	if rec := slotRecord(t, fixture.root, time.Now().UTC()); rec.partition() == nil {
		t.Fatalf("symlinked capture-only record = %+v, want a partial verdict — the symlink did not admit the kit's declarations", rec)
	}
	if got := fullRunCount(t, fixture.root); got != 1 {
		t.Fatalf("resolved gate runs through the symlink = %d, want 1 — the capture-only edit narrowed instead of paying the full gate", got)
	}
	want, _ := unconditionalPhaseNames(fixture.phases)
	narrowed := append([]string(nil), phaseRunNames(t, fixture.root)[len(seeded):]...)
	sort.Strings(narrowed)
	sort.Strings(want)
	if !reflect.DeepEqual(narrowed, want) {
		t.Fatalf("phases executed through the symlink = %v, want the unconditional set %v", narrowed, want)
	}
	afterSymlink := phaseRunNames(t, fixture.root)

	// A stat failure on either path must run everything: point BENCH_KIT at a path that
	// stats to nothing, mid-run's shape — created, then removed before the gate reads it.
	vanished := filepath.Join(t.TempDir(), "kit-vanishes")
	if err := os.Mkdir(vanished, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(vanished); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_KIT", vanished)
	for _, path := range captureSurfacePaths(ReducedScope()) {
		writeGateTestFile(t, fixture.root, path, "second capture-only edit\n", 0o644)
	}
	stdout.Reset()
	if got := executeWithEngine(context.Background(), fixture.root, &stdout, io.Discard, productionGateEngine{}); got.ActionExit != 0 {
		t.Fatalf("stat-failure execution = %+v, want green", got)
	}
	if rec := slotRecord(t, fixture.root, time.Now().UTC()); rec.Reduced || rec.partition() != nil {
		t.Fatalf("stat-failure record = %+v, want a full verdict", rec)
	}
	if got := fullRunCount(t, fixture.root); got != 2 {
		t.Fatalf("resolved gate runs after the stat failure = %d, want 2 — the vanished kit path did not pay the full gate", got)
	}
	fullTable := append([]string(nil), phaseRunNames(t, fixture.root)[len(afterSymlink):]...)
	sort.Strings(fullTable)
	want = append([]string(nil), fixture.phaseNames()...)
	sort.Strings(want)
	if !reflect.DeepEqual(fullTable, want) {
		t.Fatalf("phases executed after the stat failure = %v, want every table phase %v", fullTable, want)
	}
	if strings.Contains(stdout.String(), "gate: skipping") {
		t.Fatalf("stat-failure root announced a per-component skip:\n%s", stdout.String())
	}
}

// `bench gate --fresh` stays the operator's escape to a real whole-tree run: forceRun
// never consults the inheritance path, and the force does not outlive the flag.
func TestFreshFlagForcesFullRunPastAValidAncestor(t *testing.T) {
	root := reducedRunFixture(t)
	mustExecuteGreen(t, root, productionGateEngine{})

	writeGateTestFile(t, root, "ROADMAP.md", "capture-only edit\n", 0o644)
	if got := RunCommand([]string{"--fresh", root}, io.Discard, io.Discard); got != 0 {
		t.Fatalf("forced run exit = %d, want 0", got)
	}
	if rec := slotRecord(t, root, time.Now().UTC()); rec.Reduced {
		t.Fatalf("--fresh record = %+v — the forced run inherited instead of running the whole tree", rec)
	}
	if got := fullRunCount(t, root); got != 2 {
		t.Fatalf("gate runs after --fresh = %d, want 2", got)
	}

	writeGateTestFile(t, root, "ROADMAP.md", "another capture-only edit\n", 0o644)
	mustExecuteGreen(t, root, productionGateEngine{})
	if rec := slotRecord(t, root, time.Now().UTC()); !rec.Reduced {
		t.Fatalf("post-fresh record = %+v — the force outlived the flag", rec)
	}
	if got := fullRunCount(t, root); got != 2 {
		t.Fatalf("gate runs after the plain run = %d, want 2", got)
	}
}
