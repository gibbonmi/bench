package gate

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// reuseMarkerRepo is a throwaway repo whose gate appends one line per run to
// .git/runs, so "the oracle did not run again" is a counted fact rather than an
// inference from the returned exit code.
func reuseMarkerRepo(t *testing.T, exit int, manifest string) string {
	t.Helper()
	return gateTestRepo(t, fmt.Sprintf("#!/usr/bin/env bash\necho run >> .git/runs\nexit %d\n", exit), manifest)
}

func gateRunCount(t *testing.T, root string) int {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, ".git", "runs"))
	if errors.Is(err, os.ErrNotExist) {
		return 0
	}
	if err != nil {
		t.Fatal(err)
	}
	return len(strings.Fields(string(data)))
}

// TestReusableGreenIsReusedWithoutRunningOrWriting grades the three halves of a reuse
// that no single assertion covers: the oracle is not re-run, the durable record is not
// rewritten (a rewritten RecordedAt would slide the freshness window forward on every
// read), and the returned tuple projects the green rather than a zero Inspection that
// would read as `absent` on a green tree.
func TestReusableGreenIsReusedWithoutRunningOrWriting(t *testing.T) {
	root := reuseMarkerRepo(t, 0, `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	first := Execute(context.Background(), root, io.Discard, io.Discard)
	if first.ActionExit != 0 || !first.Inspection.ReusableGreen {
		t.Fatalf("first execution = %+v, want a reusable green", first)
	}
	before := mustRead(t, cachePath(t, root))

	var stdout bytes.Buffer
	second := Execute(context.Background(), root, &stdout, io.Discard)
	if got := gateRunCount(t, root); got != 1 {
		t.Fatalf("gate runs = %d, want 1 — the reusable green paid a second run", got)
	}
	if got, want := stdout.String(), "gate: green (fresh verdict reused for this tree)\n"; got != want {
		t.Fatalf("reuse stdout = %q, want %q", got, want)
	}
	if after := mustRead(t, cachePath(t, root)); !bytes.Equal(before, after) {
		t.Fatalf("reuse rewrote the verdict record:\nbefore %q\nafter  %q", before, after)
	}
	if second.GateExit != 0 || second.ActionExit != 0 || second.Inspection.State != Ready || second.Inspection.Status != "green" || !second.Inspection.ReusableGreen {
		t.Fatalf("reused result = %+v, want 0/0 with an inspection projecting the reusable green", second)
	}
}

// TestFreshFlagForcesARealRunPastAReusableGreen grades the operator's escape from a
// green the oracle would no longer stand behind. The flag has to survive RunCommand's
// argument parsing, where it shares the argument list with the optional root positional,
// and it has to stay opt-in: the same tree without it still answers from the record.
func TestFreshFlagForcesARealRunPastAReusableGreen(t *testing.T) {
	const closed = `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`
	for _, order := range []struct {
		name string
		args func(root string) []string
	}{
		{"flag first", func(root string) []string { return []string{"--fresh", root} }},
		{"root first", func(root string) []string { return []string{root, "--fresh"} }},
	} {
		t.Run(order.name, func(t *testing.T) {
			root := reuseMarkerRepo(t, 0, closed)
			if got := Execute(context.Background(), root, io.Discard, io.Discard); got.ActionExit != 0 || !got.Inspection.ReusableGreen {
				t.Fatalf("seed execution = %+v, want a reusable green", got)
			}
			if got := RunCommand(order.args(root), io.Discard, io.Discard); got != 0 {
				t.Fatalf("forced run exit = %d, want 0", got)
			}
			if got := gateRunCount(t, root); got != 2 {
				t.Fatalf("gate runs after --fresh = %d, want 2 — the flag did not force a real run", got)
			}
			if got := RunCommand([]string{root}, io.Discard, io.Discard); got != 0 {
				t.Fatalf("following run exit = %d, want 0", got)
			}
			if got := gateRunCount(t, root); got != 2 {
				t.Fatalf("gate runs after a plain run = %d, want 2 — the force outlived the flag", got)
			}
		})
	}
}

// TestNonReusableSubjectsPayARealRun pins the short-circuit to the ReusableGreen
// predicate. The cheapest wrong implementation short-circuits on any cached green; each
// case here leaves exactly one run behind and requires the next execution to reach the
// oracle anyway.
func TestNonReusableSubjectsPayARealRun(t *testing.T) {
	const closed = `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`
	seedGreen := func(t *testing.T, root string) {
		t.Helper()
		if got := Execute(context.Background(), root, io.Discard, io.Discard); got.ActionExit != 0 {
			t.Fatalf("seed execution = %+v, want green", got)
		}
	}
	cases := []struct {
		name    string
		arrange func(t *testing.T) (string, gateEngine)
	}{
		{"recorded red", func(t *testing.T) (string, gateEngine) {
			root := reuseMarkerRepo(t, 1, closed)
			if got := Execute(context.Background(), root, io.Discard, io.Discard); got.ActionExit == 0 {
				t.Fatalf("red gate execution = %+v, want a red verdict", got)
			}
			return root, productionGateEngine{}
		}},
		{"expired verdict", func(t *testing.T) (string, gateEngine) {
			root := reuseMarkerRepo(t, 0, closed)
			seedGreen(t, root)
			return root, &faultEngine{now: time.Now().UTC().Truncate(time.Second).Add(freshness + time.Minute)}
		}},
		{"pending record", func(t *testing.T) (string, gateEngine) {
			root := reuseMarkerRepo(t, 0, closed)
			seedGreen(t, root)
			plan := mustSubject(t, root)
			pending := verdictRecord{Schema: 1, State: Pending, Tree: plan.Tree, Oracle: plan.Oracle, StartedAt: time.Now().UTC().Add(-time.Hour).Truncate(time.Second).Format(time.RFC3339), OwnerPID: 99999999}
			if err := durableReplace(filepath.Dir(cachePath(t, root)), pending); err != nil {
				t.Fatal(err)
			}
			return root, productionGateEngine{}
		}},
		{"open subject", func(t *testing.T) (string, gateEngine) {
			root := reuseMarkerRepo(t, 0, "")
			seedGreen(t, root)
			return root, productionGateEngine{}
		}},
		{"changed tree", func(t *testing.T) (string, gateEngine) {
			root := reuseMarkerRepo(t, 0, closed)
			seedGreen(t, root)
			writeGateTestFile(t, root, "changed.txt", "changed\n", 0o644)
			return root, productionGateEngine{}
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, engine := tc.arrange(t)
			if got := gateRunCount(t, root); got != 1 {
				t.Fatalf("arranged gate runs = %d, want 1", got)
			}
			if got := inspectAt(root, engine.Now()); got.ReusableGreen {
				t.Fatalf("arranged subject is reusable: %+v", got)
			}
			executeWithEngine(context.Background(), root, io.Discard, io.Discard, engine)
			if got := gateRunCount(t, root); got != 2 {
				t.Fatalf("gate runs = %d, want 2 — a non-reusable subject skipped the oracle", got)
			}
		})
	}
}

func TestInspectTreeNeverRunsGate(t *testing.T) {
	root := reusableEvidenceRepo(t, 0)
	tree := gitOutput(t, root, "write-tree")
	if got := Execute(context.Background(), root, io.Discard, io.Discard); got.ActionExit != 0 {
		t.Fatalf("seed execution = %+v, want green", got)
	}
	if got := InspectTree(root, tree); !got.ReusableGreen {
		t.Fatalf("bootstrap inspection = %+v, want retained exact green", got)
	}
	if got := gateRunCount(t, root); got != 1 {
		t.Fatalf("bootstrap gate runs = %d, want 1", got)
	}
}

func reusableEvidenceRepo(t *testing.T, exit int) string {
	t.Helper()
	root := reuseMarkerRepo(t, exit, `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	gitRun(t, root, "add", ".")
	gitRun(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "subject")
	return root
}

func TestExecuteTreeRunsProspectiveWrapper(t *testing.T) {
	root := reusableEvidenceRepo(t, 0)
	writeGateTestFile(t, root, ".bench/gate.sh", "#!/usr/bin/env bash\ngitdir=\"$(git rev-parse --git-common-dir)\"\nprintf prospective > \"$gitdir/prospective-run\"\nprintf run >> \"$gitdir/prospective-runs\"\n", 0o755)
	gitRun(t, root, "add", ".bench/gate.sh")
	tree := gitOutput(t, root, "write-tree")
	gitRun(t, root, "reset", "--hard", "HEAD")
	if got := ExecuteTree(context.Background(), root, tree, io.Discard, io.Discard); got.ActionExit != 0 {
		t.Fatalf("prospective execution = %+v, want green", got)
	}
	worktreeGitDir := gitOutput(t, root, "rev-parse", "--absolute-git-dir")
	if got, err := os.ReadFile(filepath.Join(worktreeGitDir, "prospective-run")); err != nil || string(got) != "prospective" {
		t.Fatalf("prospective wrapper output = %q, %v; want prospective", got, err)
	}
	if got := InspectTree(root, tree); !got.ReusableGreen {
		t.Fatalf("prospective evidence = %+v, want reusable green", got)
	}
	gitdir := gitOutput(t, root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if got := evidenceFiles(t, gitdir); len(got) != 1 {
		t.Fatalf("retained evidence after one prospective execution = %v, want one record", got)
	}
	if got := ExecuteTree(context.Background(), root, tree, io.Discard, io.Discard); got.ActionExit != 0 {
		t.Fatalf("second prospective execution = %+v, want reuse", got)
	}
	if got := evidenceFiles(t, gitdir); len(got) != 1 {
		t.Fatalf("retained evidence after repeated prospective execution = %v, want one stable record", got)
	}
	if got, err := os.ReadFile(filepath.Join(gitdir, "prospective-runs")); err != nil || string(got) != "run" {
		t.Fatalf("prospective wrapper runs = %q, %v; want one run", got, err)
	}
}

func TestValidateProjectGreenRequiresTipMarkerAndClosedSubject(t *testing.T) {
	root := reusableEvidenceRepo(t, 0)
	if got := Execute(context.Background(), root, io.Discard, io.Discard); got.ActionExit != 0 {
		t.Fatalf("seed execution = %+v, want green", got)
	}
	branch := gitOutput(t, root, "branch", "--show-current")
	gitRun(t, root, "update-ref", "refs/bench/green/"+branch, "HEAD")
	if got := ValidateProjectGreen(root, branch); !got.ReusableGreen {
		t.Fatalf("matching project-green = %+v, want reusable", got)
	}
	gitRun(t, root, "update-ref", "refs/bench/green/not-the-working-branch", "HEAD")
	if got := ValidateProjectGreen(root, "not-the-working-branch"); got.ReusableGreen || got.Reason != "working branch changed" {
		t.Fatalf("wrong-branch project-green = %+v, want branch refusal", got)
	}
	gitRun(t, root, "commit", "--allow-empty", "-q", "-m", "advance")
	if got := ValidateProjectGreen(root, branch); got.ReusableGreen {
		t.Fatalf("advanced tip project-green = %+v, want refusal", got)
	}
	gitRun(t, root, "update-ref", "refs/bench/green/"+branch, "HEAD")
	if err := os.Remove(filepath.Join(root, ".bench", "gate-inputs.json")); err != nil {
		t.Fatal(err)
	}
	if got := ValidateProjectGreen(root, branch); got.ReusableGreen {
		t.Fatalf("open subject project-green = %+v, want refusal", got)
	}
}

func TestEvidenceDoesNotReuseLatestProjectionOrRed(t *testing.T) {
	root := reusableEvidenceRepo(t, 0)
	if got := Execute(context.Background(), root, io.Discard, io.Discard); got.ActionExit != 0 {
		t.Fatalf("seed execution = %+v, want green", got)
	}
	plan := mustSubject(t, root)
	gitdir := gitOutput(t, root, "rev-parse", "--absolute-git-dir")
	red := verdictRecord{Schema: 1, State: Ready, Status: "red", Tree: plan.Tree, Oracle: plan.Oracle, RecordedAt: time.Now().UTC().Truncate(time.Second).Format(time.RFC3339)}
	if err := durableReplace(gitdir, red); err != nil {
		t.Fatal(err)
	}
	if got := InspectTree(root, plan.Tree); !got.ReusableGreen {
		t.Fatalf("retained green became dependent on latest projection: %+v", got)
	}
	if got := Execute(context.Background(), root, io.Discard, io.Discard); got.ActionExit != 0 {
		t.Fatalf("retained green execution = %+v, want reuse", got)
	}
	if got := gateRunCount(t, root); got != 1 {
		t.Fatalf("gate runs = %d, want 1 after latest projection red", got)
	}
}

func TestRetainedEvidenceSurvivesHistoryOnlyAdvance(t *testing.T) {
	root := reusableEvidenceRepo(t, 0)
	tree := gitOutput(t, root, "write-tree")
	if got := Execute(context.Background(), root, io.Discard, io.Discard); got.ActionExit != 0 {
		t.Fatalf("seed execution = %+v, want green", got)
	}
	gitRun(t, root, "commit", "--allow-empty", "-q", "-m", "history only")
	if got := InspectTree(root, tree); !got.ReusableGreen {
		t.Fatalf("history-only advance invalidated retained evidence: %+v", got)
	}
	if got := gateRunCount(t, root); got != 1 {
		t.Fatalf("history-only inspection gate runs = %d, want 1", got)
	}
}

func TestForcedRedInvalidatesRetainedGreen(t *testing.T) {
	root := gateTestRepo(t, "#!/usr/bin/env bash\necho run >> .git/runs\ntest ! -f .git/force-red\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	gitRun(t, root, "add", ".")
	gitRun(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "subject")
	if got := Execute(context.Background(), root, io.Discard, io.Discard); got.ActionExit != 0 {
		t.Fatalf("green execution = %+v, want green", got)
	}
	writeGateTestFile(t, root, ".git/force-red", "\n", 0o600)
	if got := RunCommand([]string{"--fresh", root}, io.Discard, io.Discard); got == 0 {
		t.Fatal("forced red execution unexpectedly passed")
	}
	if got := Execute(context.Background(), root, io.Discard, io.Discard); got.ActionExit == 0 {
		t.Fatalf("normal execution reused the invalidated green: %+v", got)
	}
	if got := gateRunCount(t, root); got != 3 {
		t.Fatalf("gate runs = %d, want 3 after green, forced red, and normal red", got)
	}
}

func evidenceFiles(t *testing.T, gitdir string) []string {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(gitdir, "bench-gate-evidence"))
	if err != nil {
		t.Fatal(err)
	}
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.Type().IsRegular() {
			files = append(files, entry.Name())
		}
	}
	return files
}

// [PC15] A partial verdict graded only the components whose inputs moved, so it answers for
// its own tree and never for the whole one. The refusal is what keeps every reuse path safe
// without a change at each of them — `bench commit` reads reusability through the
// authorization package, which reads it here — so the record stays readable, the inspection
// names the partition and carries it for the consumers that render and refuse against it,
// and only ReusableGreen is withheld.
func TestPartialVerdictIsNotReusable(t *testing.T) {
	root := reuseMarkerRepo(t, 0, `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	if got := Execute(context.Background(), root, io.Discard, io.Discard); got.ActionExit != 0 {
		t.Fatalf("seed execution = %+v, want green", got)
	}
	plan := mustSubject(t, root)
	seeded := time.Now().UTC().Truncate(time.Second)
	partial := partialTestRecord(seeded)
	partial.Tree, partial.Oracle = plan.Tree, plan.Oracle
	if err := durableReplace(filepath.Dir(cachePath(t, root)), partial); err != nil {
		t.Fatal(err)
	}

	got := inspectAt(root, seeded)
	if got.State != Ready || got.Status != "green" || got.ReusableGreen || got.Reason != "partial verdict" {
		t.Fatalf("partial inspection = %+v, want a readable green that is not a reusable whole-tree green", got)
	}
	if got.Partition == nil {
		t.Fatalf("partial inspection carried no partition: %+v", got)
	}
	skipped := make([]string, 0, len(got.Partition.Skipped))
	for _, skip := range got.Partition.Skipped {
		skipped = append(skipped, skip.Component)
	}
	if !reflect.DeepEqual(skipped, partial.Skipped) || !reflect.DeepEqual(got.Partition.Executed, partial.Executed) {
		t.Fatalf("inspection partition = %+v, want executed %q and skipped %q", got.Partition, partial.Executed, partial.Skipped)
	}

	// The same tree and oracle under a full record is reusable, so the refusal above is the
	// partition's and not some other property of the seeded record.
	full := fullTestRecord(seeded)
	full.Tree, full.Oracle = plan.Tree, plan.Oracle
	if err := durableReplace(filepath.Dir(cachePath(t, root)), full); err != nil {
		t.Fatal(err)
	}
	if got := inspectAt(root, seeded); !got.ReusableGreen || got.Partition != nil {
		t.Fatalf("full inspection over the same tree = %+v, want a reusable green carrying no partition", got)
	}
}
