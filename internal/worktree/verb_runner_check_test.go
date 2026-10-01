package worktree

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// verbRecorder is the testing.TB a must-form test passes. Its embedded TB is nil, so a must
// form that reaches past the four overrides panics instead of passing silently.
type verbRecorder struct {
	testing.TB
	failures []string
}

func (r *verbRecorder) Helper() {}

func (r *verbRecorder) Fatal(args ...any) { r.failures = append(r.failures, fmt.Sprint(args...)) }

func (r *verbRecorder) Fatalf(format string, args ...any) {
	r.failures = append(r.failures, fmt.Sprintf(format, args...))
}

func (r *verbRecorder) FailNow() { r.failures = append(r.failures, "FailNow") }

// usageCommand is the command name that leads a usage constant. A grammar that refuses an
// extra argument names only that command, not the whole usage line.
func usageCommand(line string) string {
	return strings.Join(strings.Fields(line)[:3], " ")
}

// runnerRepo is a call on a new repository with a private home, so no test binds the
// process environment.
func runnerRepo(t *testing.T) verbCall {
	t.Helper()
	root := newWorktreeRepo(t)
	return verbCall{root: root, home: filepath.Join(root, ".bench-home")}
}

// landedRunnerPlan is a real `clean --landed` plan over two landed assignments.
func landedRunnerPlan(t *testing.T) verbResult {
	t.Helper()
	call := runnerRepo(t)
	for _, name := range []string{"first", "second"} {
		landAssignment(t, call.root, mustCreate(t, call.root, call.home, "runner-landed-"+name, name), name+".txt")
	}
	call.args = []string{"--landed"}
	result := runVerb(t, verbClean, call)
	requireTest(t, result.exit == 0, "clean --landed = %d %q %q", result.exit, result.stdout, result.stderr)
	return result
}

// explicitErrorPlan is a real explicit-set `clean` plan whose one target does not resolve.
func explicitErrorPlan(t *testing.T) verbResult {
	t.Helper()
	call := runnerRepo(t)
	call.args = []string{"--target", "no-such-target"}
	result := runVerb(t, verbClean, call)
	requireTest(t, result.exit == 1, "clean --target = %d %q %q", result.exit, result.stdout, result.stderr)
	return result
}

func TestVerbRunnerKeyReachesItsOwnVerb(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		key  verbKey
		args []string
		want string
	}{
		{verbCreate, nil, "usage: " + usage.WorktreeCreate},
		{verbRelease, nil, "usage: " + usage.WorktreeRelease},
		{verbClean, nil, "invalid invocation; run " + usage.WorktreeClean},
		{verbReclaim, []string{"--bogus"}, toon.Usage(usage.WorktreeReclaim, "--bogus")},
		{verbReauthorize, []string{"p"}, "usage: " + usage.WorktreeReauthorize},
		{verbMerge, []string{"t"}, "usage: " + usage.WorktreeMerge},
		{verbReset, []string{"t"}, "usage: " + usage.WorktreeReset},
		{verbLand, []string{"p"}, "usage: " + usage.WorktreeLand},
		{verbLandResume, []string{"p"}, "usage: " + usage.WorktreeLandResume},
		{verbList, []string{"x"}, toon.Usage(usage.WorktreeList, "x")},
		{verbPath, nil, "usage: " + usage.WorktreePath},
		{verbShow, []string{"a", "b", "c"}, toon.Usage(usageCommand(usage.WorktreeShow), "c")},
		{verbBuild, []string{"a", "b"}, toon.Usage(usageCommand(usage.WorktreeBuild), "b")},
		{verbExec, []string{"t", "x"}, "usage: " + usage.WorktreeExec},
	} {
		t.Run(string(row.key), func(t *testing.T) {
			result := runVerb(t, row.key, verbCall{root: t.TempDir(), home: t.TempDir(), args: row.args})
			if result.exit != 2 || !strings.Contains(result.stdout+result.stderr, row.want) {
				t.Fatalf("%s = %d %q %q, want exit 2 with %q", row.key, result.exit, result.stdout, result.stderr, row.want)
			}
		})
	}
}

func TestVerbRunnerPoolKeyReturnsThePoolPath(t *testing.T) {
	t.Parallel()
	root, home := t.TempDir(), t.TempDir()
	result := runVerb(t, verbPool, verbCall{home: home, args: []string{root}})
	if want := poolAt(home, root) + "\n"; result.exit != 0 || result.stdout != want || result.stderr != "" {
		t.Fatalf("pool = %d %q %q, want %q", result.exit, result.stdout, result.stderr, want)
	}
}

func TestVerbRunnerLeaseFileKeyMatchesItsEntry(t *testing.T) {
	t.Parallel()
	want, code := LeaseFileCommand(nil)
	result := runVerb(t, verbLeaseFile, verbCall{})
	if result.exit != code || result.stdout != want || result.stderr != "" {
		t.Fatalf("lease-file = %d %q %q, want %d %q", result.exit, result.stdout, result.stderr, code, want)
	}
}

func TestVerbRunnerResumeCleanKeyMatchesItsEntry(t *testing.T) {
	t.Parallel()
	root, home := t.TempDir(), t.TempDir()
	var stdout, stderr bytes.Buffer
	code := ResumeCleanCommand(root, home, []string{"--bogus"}, &stdout, &stderr)
	result := runVerb(t, verbResumeClean, verbCall{root: root, home: home, args: []string{"--bogus"}})
	if result.exit != code || result.stdout != stdout.String() || result.stderr != stderr.String() {
		t.Fatalf("resume-clean = %d %q %q, want %d %q %q", result.exit, result.stdout, result.stderr, code, stdout.String(), stderr.String())
	}
}

func TestVerbRunnerReturnsBothStreamsAndTheExitCode(t *testing.T) {
	t.Parallel()
	call := runnerRepo(t)
	creation := mustCreate(t, call.root, call.home, "runner-streams", "streams")
	call.args = []string{creation.Assignment.Label, "--", "sh", "-c", "printf out; printf err >&2; exit 3"}
	var stdout, stderr bytes.Buffer
	code := ExecCommand(call.root, call.home, call.args, nil, &stdout, &stderr)
	requireTest(t, code == 3 && stdout.Len() > 0 && stderr.Len() > 0, "direct exec = %d %q %q, want exit 3 on both streams", code, stdout.String(), stderr.String())
	result := runVerb(t, verbExec, call)
	if result.exit != code || result.stdout != stdout.String() || result.stderr != stderr.String() {
		t.Fatalf("runner = %d %q %q, want %d %q %q", result.exit, result.stdout, result.stderr, code, stdout.String(), stderr.String())
	}
}

func TestVerbResultRowsDecodeTheWholeDocument(t *testing.T) {
	t.Parallel()
	call := runnerRepo(t)
	mustCreate(t, call.root, call.home, "runner-rows-first", "first")
	mustCreate(t, call.root, call.home, "runner-rows-second", "second")
	result := runVerb(t, verbList, call)
	document, err := axitest.DecodeDocument(result.stdout)
	mustNoError(t, err)
	want, err := document.Rows("worktrees")
	mustNoError(t, err)
	rows, err := readVerbRows(result.stdout, "worktrees")
	if err != nil || len(rows) != 2 || !reflect.DeepEqual(rows, want) {
		t.Fatalf("rows = %#v, %v; want the two decoded rows %#v", rows, err, want)
	}
}

func TestVerbResultRowsRefuseAPartialDocument(t *testing.T) {
	t.Parallel()
	table, err := toon.Table("worktrees", []string{"id"}, [][]string{{"one"}})
	mustNoError(t, err)
	if rows, err := readVerbRows(table+"a line outside the grammar\n", "worktrees"); err == nil {
		t.Fatalf("rows = %#v with no error, want a decode error", rows)
	}
}

func TestVerbResultMustRowsFailsOnAReaderError(t *testing.T) {
	t.Parallel()
	recorder := &verbRecorder{}
	verbResult{stdout: "a line outside the grammar\n"}.mustRows(recorder, "worktrees")
	if len(recorder.failures) != 1 {
		t.Fatalf("recorder failures = %q, want one", recorder.failures)
	}
}

func TestVerbResultFingerprintReadsTheTableCell(t *testing.T) {
	t.Parallel()
	result := landedRunnerPlan(t)
	value, err := readVerbFingerprint(result.stdout)
	if err != nil || len(value) != 64 || !strings.Contains(result.stdout, "--landed --apply "+value) {
		t.Fatalf("fingerprint = %q, %v; want the set fingerprint that the apply action names in %q", value, err, result.stdout)
	}
	for _, row := range result.mustRows(t, "worktree_cleanup") {
		if cell := row.(map[string]any)["fingerprint"]; cell != value {
			t.Fatalf("row fingerprint = %#v, want %q", cell, value)
		}
	}
}

func TestVerbResultFingerprintReadsTheRecordCell(t *testing.T) {
	t.Parallel()
	call := runnerRepo(t)
	creation := mustCreate(t, call.root, call.home, "runner-record", "record")
	commitInWorktree(t, creation.Path, "ahead", "ahead\n", "ahead")
	call.args = []string{"--to", creation.Assignment.Start, creation.Assignment.ID}
	result := runVerb(t, verbReset, call)
	requireTest(t, result.exit == 0, "reset --to = %d %q %q", result.exit, result.stdout, result.stderr)
	record, _, _ := strings.Cut(result.stdout, "\n")
	value, err := readVerbFingerprint(result.stdout)
	if err != nil || value == "" || !strings.Contains(record, " --apply "+value+",") || !strings.HasSuffix(record, ",fingerprint="+value+"}") {
		t.Fatalf("fingerprint = %q, %v; want the record cell that the next cell applies in %q", value, err, record)
	}
}

func TestVerbResultFingerprintRefusesAnAbsentValue(t *testing.T) {
	t.Parallel()
	table, err := toon.Table("worktree_cleanup", []string{"target"}, [][]string{{"one"}})
	mustNoError(t, err)
	if value, err := readVerbFingerprint(table); !errors.Is(err, errNoVerbFingerprint) {
		t.Fatalf("fingerprint = %q, %v; want the no-fingerprint error", value, err)
	}
}

func TestVerbResultFingerprintRefusesConflictingCells(t *testing.T) {
	t.Parallel()
	table, err := toon.Table("worktree_cleanup", []string{"target", "fingerprint"}, [][]string{{"one", "aaaa"}, {"two", "bbbb"}})
	mustNoError(t, err)
	if value, err := readVerbFingerprint(table); err == nil {
		t.Fatalf("fingerprint = %q with no error, want a conflict error", value)
	}
}

func TestVerbResultMustFingerprintFailsOnAReaderError(t *testing.T) {
	t.Parallel()
	recorder := &verbRecorder{}
	verbResult{stdout: "no fingerprint here\n"}.mustFingerprint(recorder)
	if len(recorder.failures) != 1 {
		t.Fatalf("recorder failures = %q, want one", recorder.failures)
	}
}

func TestVerbResultFingerprintKeepsANumericLookingCell(t *testing.T) {
	t.Parallel()
	const numeric = "1.50"
	table, err := toon.Table("worktree_cleanup", []string{"target", "fingerprint"}, [][]string{{"one", numeric}, {"two", numeric}})
	mustNoError(t, err)
	requireTest(t, strings.Contains(table, `"`+numeric+`"`), "the producer did not quote %q: %q", numeric, table)
	if value, err := readVerbFingerprint(table); err != nil || value != numeric {
		t.Fatalf("fingerprint = %q, %v; want %q", value, err, numeric)
	}
}

func TestVerbResultMustNoFingerprintAcceptsAnErrorPlan(t *testing.T) {
	t.Parallel()
	result := explicitErrorPlan(t)
	recorder := &verbRecorder{}
	result.mustNoFingerprint(recorder)
	if len(recorder.failures) != 0 {
		t.Fatalf("recorder failures = %q for %q, want none", recorder.failures, result.stdout)
	}
}

func TestVerbResultMustNoFingerprintRefusesAPlan(t *testing.T) {
	t.Parallel()
	recorder := &verbRecorder{}
	landedRunnerPlan(t).mustNoFingerprint(recorder)
	if len(recorder.failures) != 1 {
		t.Fatalf("recorder failures = %q, want one", recorder.failures)
	}
}

// TestVerbResultFingerprintTreatsAPlaceholderAsAbsent builds the faulted unclaimed set as
// TestCleanUnclaimedErrorRowRefusesTheSet does: the loose ref file names a blob, because git
// update-ref can refuse a ref that names no commit.
func TestVerbResultFingerprintTreatsAPlaceholderAsAbsent(t *testing.T) {
	t.Parallel()
	explicit := explicitErrorPlan(t)
	root, home := unclaimedBranchFixture(t)
	addUnclaimedBranch(t, root, "a")
	loose := filepath.Join(root, ".git", filepath.FromSlash(unclaimedBranchRef("b")))
	mustNoError(t, os.MkdirAll(filepath.Dir(loose), 0o755))
	mustWrite(t, loose, []byte(gitOutput(t, root, "hash-object", "-w", "tracked.txt")+"\n"), 0o644)
	faulted := runVerb(t, verbClean, verbCall{root: root, home: home, args: []string{"--discard-branch", "--unclaimed"}})
	requireTest(t, faulted.exit == 1, "faulted plan = %d %q %q", faulted.exit, faulted.stdout, faulted.stderr)
	for _, plan := range []struct {
		result verbResult
		cell   string
	}{{explicit, unapplicableFingerprint}, {faulted, ""}} {
		for _, row := range plan.result.mustRows(t, "worktree_cleanup") {
			requireTest(t, row.(map[string]any)["fingerprint"] == plan.cell, "row = %#v, want fingerprint %q", row, plan.cell)
		}
		if value, err := readVerbFingerprint(plan.result.stdout); !errors.Is(err, errNoVerbFingerprint) {
			t.Fatalf("fingerprint = %q, %v for %q; want the no-fingerprint error", value, err, plan.result.stdout)
		}
	}
}

func TestVerbRunnerPassesTheJoinsValue(t *testing.T) {
	t.Parallel()
	call := runnerRepo(t)
	creation := mustCreate(t, call.root, call.home, "runner-joins", "joins")
	probe := 0
	j := defaultJoins()
	j.buildSubject = func(context.Context, string, string) error { probe++; return nil }
	call.joins, call.args = &j, []string{creation.Assignment.Label}
	result := runVerb(t, verbBuild, call)
	if result.exit != 0 || probe != 1 {
		t.Fatalf("build = %d %q %q with probe %d, want exit 0 and one stub call", result.exit, result.stdout, result.stderr, probe)
	}
}

func TestVerbRunnerFeedsStdinToExec(t *testing.T) {
	t.Parallel()
	call := runnerRepo(t)
	creation := mustCreate(t, call.root, call.home, "runner-stdin", "stdin")
	fed := "first\nsecond\x00"
	call.stdin, call.args = strings.NewReader(fed), []string{creation.Assignment.Label, "--", "cat"}
	result := runVerb(t, verbExec, call)
	if result.exit != 0 || result.stdout != fed {
		t.Fatalf("exec = %d %q %q, want the fed bytes %q", result.exit, result.stdout, result.stderr, fed)
	}
}

func TestVerbRunnerReturnsTheExecAssignment(t *testing.T) {
	t.Parallel()
	call := runnerRepo(t)
	creation := mustCreate(t, call.root, call.home, "runner-assignment", "assignment")
	call.args = []string{creation.Assignment.Label, "--", "true"}
	result := runVerb(t, verbExec, call)
	if result.exit != 0 || result.assignment != creation.Assignment.ID {
		t.Fatalf("exec = %d %q with assignment %q, want %q", result.exit, result.stderr, result.assignment, creation.Assignment.ID)
	}
}

func TestVerbCallRefusesAKitValue(t *testing.T) {
	t.Parallel()
	err := checkVerbCall(verbCall{kit: t.TempDir()})
	if want := "verb runner: a kit value waits for the seam reduction spec"; err == nil || err.Error() != want {
		t.Fatalf("checkVerbCall = %v, want %q", err, want)
	}
}

func TestVerbCallRefusesAClockValue(t *testing.T) {
	t.Parallel()
	err := checkVerbCall(verbCall{clock: time.Now})
	if want := "verb runner: a clock value waits for the seam reduction spec"; err == nil || err.Error() != want {
		t.Fatalf("checkVerbCall = %v, want %q", err, want)
	}
}
