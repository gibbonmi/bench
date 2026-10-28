package gate

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/testrepo"
)

// TestChunkCheckpointReusesABackgroundGreen is BG1. After a background grade of the tip, the
// chunk checkpoint of that tip prints the reuse line and runs no oracle.
func TestChunkCheckpointReusesABackgroundGreen(t *testing.T) {
	f := backgroundFixture(t)
	requireBackgroundGreen(t, f, nil)
	runs := outcomeRuns(t, f.Root)
	if code, out := runCheckpoint(t, f); code != 0 || !strings.Contains(out, "reused") || outcomeRuns(t, f.Root) != runs {
		t.Fatalf("checkpoint after a background green = %d %s after %d runs; want the reused green and %d runs", code, out, outcomeRuns(t, f.Root), runs)
	}
}

// TestBackgroundChunkGradeKeysTheStrippedTree is BG23. The background green is under the key of
// the HEAD tree without the review record, not under the key of the raw HEAD tree.
func TestBackgroundChunkGradeKeysTheStrippedTree(t *testing.T) {
	f := backgroundFixture(t)
	requireBackgroundGreen(t, f, nil)
	stripped, raw := strippedEvidence(t, f), InspectTree(f.Root, f.Tree())
	if got := backgroundRecord(t, f); !stripped.ReusableGreen || got.Key != stripped.Grade || got.Key == raw.Grade {
		t.Fatalf("evidence of the stripped tree = %+v, of the raw tree = %+v, record = %+v; want the green of the stripped tree under the recorded key", stripped, raw, got)
	}
}

// TestBackgroundGradeRunsWholeWithNoChunkEntry is BG24. A record with no chunk entry resolves no
// chunk id, so the grade runs the whole schedule and records its green.
func TestBackgroundGradeRunsWholeWithNoChunkEntry(t *testing.T) {
	f := recordtest.Attach(t, outcomeFixture(t), 1)
	f.Save()
	f.Commit("retain the record before its first chunk entry")
	seedGateCache(t, f)
	requireBackgroundGreen(t, f, nil)
	if got := strippedEvidence(t, f); !got.ReusableGreen {
		t.Fatalf("evidence of the stripped tree = %+v; want the green of the whole schedule", got)
	}
}

// TestBackgroundGradeSkipsTheEvidencePass is BG26. The chunk checkpoint refuses a frozen chunk
// with no review, and the background grade of the same tip still yields its green.
func TestBackgroundGradeSkipsTheEvidencePass(t *testing.T) {
	f, _ := unreviewedFixture(t)
	seedGateCache(t, f)
	if code, out := runCheckpoint(t, f); code == 0 || !strings.Contains(out, "completion evidence") {
		t.Fatalf("checkpoint of an unreviewed chunk = %d %s; want the evidence refusal", code, out)
	}
	requireBackgroundGreen(t, f, nil)
	if got := strippedEvidence(t, f); !got.ReusableGreen {
		t.Fatalf("evidence of the stripped tree = %+v; want the background green", got)
	}
}

// TestBackgroundGradeRunsAsTheTargetExecutable is BG27. The entry runs the grade as the
// executable that the wrapper of the target worktree resolves: it clears the inherited run
// binary and kit, so the wrapper selects the dist/bench of the worktree, not the launcher.
func TestBackgroundGradeRunsAsTheTargetExecutable(t *testing.T) {
	root := gittest.RepoOnBranch(t, "main")
	wrapper, err := os.ReadFile(filepath.Join("..", "..", "bin", "bench.sh"))
	if err != nil {
		t.Fatal(err)
	}
	outcomeWrite(t, root, "bin/bench.sh", string(wrapper), 0o755)
	argv := filepath.Join(t.TempDir(), "argv")
	recorder := "#!/bin/sh\nprintf '%s\\n' \"$0\" \"$BENCH_KIT\" \"${BENCH_RUN_BINARY-unset}\" \"$" + backgroundResolvedEnv + "\" \"$@\" > \"$BENCH_TEST_ARGV\"\n"
	outcomeWrite(t, root, "dist/bench", recorder, 0o755)
	launcher := filepath.Join(t.TempDir(), "launcher")
	outcomeWrite(t, filepath.Dir(launcher), filepath.Base(launcher), recorder, 0o755)
	launcher, err = filepath.EvalSymlinks(launcher)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_TEST_ARGV", argv)
	t.Setenv(runbinary.Env, launcher)
	t.Setenv("BENCH_KIT", filepath.Dir(launcher))
	t.Setenv(backgroundResolvedEnv, "")
	previous := reenterImage
	t.Cleanup(func() { reenterImage = previous })
	reentered := errors.New("the test ran the replacement image as a child")
	reenterImage = func(dir, name string, args, env []string) error {
		command := sanitize.ShellQuote(name)
		for _, arg := range args {
			command += " " + sanitize.ShellQuote(arg)
		}
		var out bytes.Buffer
		if run := runResolved(context.Background(), dir, Resolution{Kind: BenchGate, Command: command}, env, &out, &out); run.Code != 0 {
			t.Fatalf("replacement image %s = %d\n%s", command, run.Code, out.String())
		}
		return reentered
	}
	var stderr bytes.Buffer
	if code := BackgroundCommand([]string{root, "HEAD"}, io.Discard, &stderr); code != 1 || !strings.Contains(stderr.String(), reentered.Error()) {
		t.Fatalf("entry = %d %s; want the replacement of its image", code, stderr.String())
	}
	physical, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{filepath.Join(physical, "dist", "bench"), physical, "unset", "1", backgroundVerb, root, "HEAD"}
	if got := strings.Split(strings.TrimSuffix(string(outcomeRead(t, argv)), "\n"), "\n"); !slices.Equal(got, want) {
		t.Fatalf("resolved executable and argv = %q; want %q", got, want)
	}
}

// TestBackgroundGradeWritesNoGateCache is BG28. The grade leaves the gate cache of its worktree
// byte for byte unchanged, before its oracle and after it.
func TestBackgroundGradeWritesNoGateCache(t *testing.T) {
	f := backgroundFixture(t)
	before := outcomeRecord(t, f.Root)
	requireBackgroundGreen(t, f, nil)
	during := outcomeRead(t, recordDuringWitness.path(t, f.Root))
	if after := outcomeRecord(t, f.Root); !bytes.Equal(during, before) || !bytes.Equal(after, before) {
		t.Fatalf("gate cache before = %s, during = %s, after = %s; want one unchanged cache", before, during, after)
	}
}

// TestForegroundTakesTheLockDuringAGrade is BG29. While the grade holds its live claim, the
// checkout lock is free, and a fresh foreground gate in the worktree takes it and runs.
func TestForegroundTakesTheLockDuringAGrade(t *testing.T) {
	f := backgroundFixture(t)
	key := strippedEvidence(t, f).Grade
	var claim backgroundstate.Claim
	var live, held bool
	var code int
	var stderr bytes.Buffer
	ran, runs := false, outcomeRuns(t, f.Root)
	grade, out := gradeBackground(t, f, func() {
		ran = true
		claim, _ = backgroundstate.ReadClaim(commonDir(t, f.Root), key)
		live, _ = claim.Live()
		held, _ = ExecutionInProgress(f.Root)
		code = RunCommand([]string{flagFresh, f.Root}, io.Discard, &stderr)
	})
	if !ran || !live || claim.Worktree != f.Root || held || code != 0 {
		t.Fatalf("during the grade: ran = %v, claim = %+v (live %v), lock held = %v, foreground = %d %s; want a live claim, a free lock, and a green foreground run", ran, claim, live, held, code, stderr.String())
	}
	if grade != 0 || outcomeRuns(t, f.Root) != runs+2 {
		t.Fatalf("background grade = %d %s after %d runs; want exit 0 after the foreground run and its own", grade, out, outcomeRuns(t, f.Root))
	}
}

// TestWrongChunkOnlyMisses is BG41. The latest chunk entry names a chunk id that the plan does
// not hold, and the grade writes its one green only under the key that it ran.
func TestWrongChunkOnlyMisses(t *testing.T) {
	f := backgroundFixture(t)
	f.Record.Chunks[0].ID = "2"
	resave(t, f, "name a wrong chunk id")
	before := storeEntries(t, f.Root)
	requireBackgroundGreen(t, f, nil)
	added := slices.DeleteFunc(storeEntries(t, f.Root), func(name string) bool { return slices.Contains(before, name) })
	if key := strippedEvidence(t, f).Grade; !slices.Equal(added, []string{key}) || backgroundRecord(t, f).Key != key {
		t.Fatalf("new store entries = %v, record = %+v; want the one key %s that the grade ran", added, backgroundRecord(t, f), key)
	}
}

// TestBackgroundSubjectIgnoresAWorktreeEdit is BG42. The worktree edits a Go file before the
// grade and during it, and the grade of the committed tip still records its green.
func TestBackgroundSubjectIgnoresAWorktreeEdit(t *testing.T) {
	f := attachedCheckpointFixture(t, recordtest.Attach, func(f *recordtest.Fixture) { f.Write("source.go", "package source\n") })
	seedGateCache(t, f)
	f.Write("source.go", "package source\n\n// An edit before the grade.\n")
	requireBackgroundGreen(t, f, func() { f.Write("source.go", "package source\n\n// An edit during the grade.\n") })
	if got := backgroundRecord(t, f); got.Outcome != backgroundstate.Green || !strippedEvidence(t, f).ReusableGreen {
		t.Fatalf("record = %+v; want a green of the committed tip", got)
	}
}

// TestBackgroundGreenRecordsGreen is BG106. A green grade records green with its key, its tip,
// and its log, and it releases its claim.
func TestBackgroundGreenRecordsGreen(t *testing.T) {
	f := backgroundFixture(t)
	requireBackgroundGreen(t, f, nil)
	key := strippedEvidence(t, f).Grade
	want := backgroundstate.Record{Outcome: backgroundstate.Green, Key: key, Tip: f.Tip(), Log: backgroundstate.LogPath(adminDir(t, f.Root))}
	if got := backgroundRecord(t, f); got != want {
		t.Fatalf("background record = %+v; want %+v", got, want)
	}
	if claim, err := backgroundstate.ReadClaim(commonDir(t, f.Root), key); err != nil || claim != (backgroundstate.Claim{}) {
		t.Fatalf("claim after the grade = %+v, %v; want it released", claim, err)
	}
}

// TestBackgroundRedRecordsRed is BG107. A grade whose oracle refuses the tree records red with
// its key and its tip.
func TestBackgroundRedRecordsRed(t *testing.T) {
	f := backgroundFixture(t)
	f.Write(".gate-red", "")
	f.Git("add", "-f", ".gate-red")
	f.Commit("commit the red marker")
	code, out := gradeBackground(t, f, nil)
	got := backgroundRecord(t, f)
	if code != 7 || got.Outcome != backgroundstate.Red || got.Key != strippedEvidence(t, f).Grade || got.Tip != f.Tip() {
		t.Fatalf("red grade = %d %s, record = %+v; want exit 7 and a red record of the tip", code, out, got)
	}
}

// TestBackgroundGradeSetsTheWidthInput is BG114. The fixture gate records the width input of its
// environment. The foreground run sees that input unset, and the background chunk grade sees it
// set.
func TestBackgroundGradeSetsTheWidthInput(t *testing.T) {
	const witness outcomeWitness = ".gate-width"
	f := recordtest.Attach(t, outcomeFixture(t, func(_ *testrepo.GateFixture, body string) string {
		return body + `printf '%s' "${` + widthEnv + `-}" > "$common/` + string(witness) + "\"\n"
	}), 1)
	f.AddChunk()
	f.Save()
	f.Commit("retain review evidence")
	seedGateCache(t, f)
	if got := outcomeRead(t, witness.path(t, f.Root)); len(got) != 0 {
		t.Fatalf("width input of the foreground oracle = %q; want it unset", got)
	}
	requireBackgroundGreen(t, f, nil)
	if got := outcomeRead(t, witness.path(t, f.Root)); len(got) == 0 {
		t.Fatalf("width input of the background oracle is unset; want it set")
	}
}

// backgroundFixture is a checkpoint fixture whose worktree holds a gate cache, as a worktree
// does after its first foreground gate. The fixture gate copies that cache on each run.
func backgroundFixture(t *testing.T) *recordtest.Fixture {
	t.Helper()
	f := checkpointFixture(t)
	seedGateCache(t, f)
	return f
}

// seedGateCache runs the ordinary gate of f once, so its worktree holds a gate cache.
func seedGateCache(t *testing.T, f *recordtest.Fixture) {
	t.Helper()
	if result := Execute(context.Background(), f.Root, io.Discard, io.Discard); result.ActionExit != 0 {
		t.Fatalf("seed gate = %+v; want a green", result)
	}
}

// gradeBackground runs one background grade of the HEAD tip of f, and calls arm once the grade
// holds its claim. It answers the exit and the output: the outcome line and the log of the run.
func gradeBackground(t *testing.T, f *recordtest.Fixture, arm func()) (int, string) {
	t.Helper()
	var hook postAcquireContextArm
	if arm != nil {
		hook = func(ctx context.Context) (context.Context, func()) { arm(); return ctx, func() {} }
	}
	var out bytes.Buffer
	code := gradeInBackground(context.Background(), f.Root, f.Tip(), hook, &out, &out)
	log, _ := os.ReadFile(backgroundstate.LogPath(adminDir(t, f.Root)))
	return code, out.String() + string(log)
}

// requireBackgroundGreen requires that one background grade of f exits 0 after one more oracle
// run.
func requireBackgroundGreen(t *testing.T, f *recordtest.Fixture, arm func()) {
	t.Helper()
	runs := outcomeRuns(t, f.Root)
	if code, out := gradeBackground(t, f, arm); code != 0 || outcomeRuns(t, f.Root) != runs+1 {
		t.Fatalf("background grade = %d %s after %d runs; want exit 0 after %d runs", code, out, outcomeRuns(t, f.Root), runs+1)
	}
}

// strippedEvidence answers the evidence of the HEAD tree of f without its review record.
func strippedEvidence(t *testing.T, f *recordtest.Fixture) EvidenceInspection {
	t.Helper()
	stripped, err := rr.SourceDigest(f.Root, f.Tree(), recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	return InspectTree(f.Root, stripped)
}

// backgroundRecord answers the background record of the worktree of f.
func backgroundRecord(t *testing.T, f *recordtest.Fixture) backgroundstate.Record {
	t.Helper()
	record, err := backgroundstate.ReadRecord(adminDir(t, f.Root))
	if err != nil {
		t.Fatal(err)
	}
	return record
}

// storeEntries answers the names of the evidence store of root.
func storeEntries(t *testing.T, root string) []string {
	t.Helper()
	entries, err := os.ReadDir(filepath.Dir(evidencePath(commonDir(t, root), subject{})))
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names
}

func adminDir(t *testing.T, root string) string {
	t.Helper()
	dir, err := benchgit.AdminDir(root)
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func commonDir(t *testing.T, root string) string {
	t.Helper()
	dir, err := benchgit.CommonDir(root)
	if err != nil {
		t.Fatal(err)
	}
	return dir
}
