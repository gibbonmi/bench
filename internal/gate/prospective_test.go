package gate

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/conformance/registry"
)

func TestExecuteTreeBuildsExactUnpublishedBenchkitSource(t *testing.T) {
	kit := kitRootForTest(t)
	root := gateTestRepo(t, string(mustReadGateTestFile(t, filepath.Join(kit, ".bench", "gate.sh"))), `{"schema":1,"closure":"local","environment":["HOME"],"paths":[],"tools":[]}`)
	for _, rel := range []string{
		".bench/gate-prospective.sh", "scripts/go-build.sh", "scripts/go-build.inputs",
		"package.json", "internal/releaseevidence/requirements.json", "internal/freshness/freshness.go",
		"internal/freshness/check/main.go", "internal/freshness/cmd/main.go",
	} {
		mode := os.FileMode(0o644)
		if strings.HasSuffix(rel, ".sh") {
			mode = 0o755
		}
		writeGateTestFile(t, root, rel, string(mustReadGateTestFile(t, filepath.Join(kit, filepath.FromSlash(rel)))), mode)
	}
	writeGateTestFile(t, root, ".gitignore", "dist/\n", 0o644)
	writeGateTestFile(t, root, "go.mod", "module github.com/gibbonmi/bench\n\ngo 1.25\n", 0o644)
	writeGateTestFile(t, root, "cmd/bench/main.go", prospectiveBenchMain("source A"), 0o644)
	gitRun(t, root, "add", ".")
	gitRun(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "source A")
	writeGateTestFile(t, root, "cmd/bench/main.go", prospectiveBenchMain("source B"), 0o644)
	gitRun(t, root, "add", "cmd/bench/main.go")
	tree := gitOutput(t, root, "write-tree")
	gitRun(t, root, "reset", "--hard", "HEAD")

	direct := exec.Command(filepath.Join(root, ".bench", "gate.sh"))
	direct.Dir = root
	direct.Env = append(os.Environ(), "BENCH_GATE_PROSPECTIVE=1", "GOCACHE="+t.TempDir())
	directOutput, directErr := direct.CombinedOutput()
	if directErr == nil || !strings.Contains(string(directOutput), "rebuild with") {
		t.Fatalf("ordinary real wrapper with ambient marker = %v, output=%q; want missing-artifact freshness refusal", directErr, directOutput)
	}

	var stdout, stderr bytes.Buffer
	if got := ExecuteTree(context.Background(), root, tree, &stdout, &stderr); got.ActionExit != 0 {
		t.Fatalf("prospective benchkit execution = %+v, want green; stdout=%q stderr=%q", got, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "source B") || strings.Contains(stdout.String(), "source A") {
		t.Fatalf("prospective output = %q, want only unpublished source B", stdout.String())
	}
	if _, err := os.Lstat(filepath.Join(root, "dist")); !os.IsNotExist(err) {
		t.Fatalf("prospective execution populated ordinary checkout dist: %v", err)
	}
	writeGateTestFile(t, root, "cmd/bench/main.go", prospectiveBenchMain("source B"), 0o644)
	gitRun(t, root, "add", "cmd/bench/main.go")
	gitRun(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "source B")
	if committed := gitOutput(t, root, "show", "-s", "--format=%T", "HEAD"); committed != tree {
		t.Fatalf("committed source B tree = %s, want gated tree %s", committed, tree)
	}
	branch := gitOutput(t, root, "branch", "--show-current")
	gitRun(t, root, "update-ref", "refs/bench/green/"+branch, "HEAD")
	if got := ValidateProjectGreen(root, branch); !got.ReusableGreen {
		t.Fatalf("prospective green did not validate committed project-green: %+v", got)
	}
}

func TestExecuteTreeIgnoresStoredCheckSlotsAndRunsFullConformanceInventory(t *testing.T) {
	fixture := newKitShapedFixture(t)
	testBinary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	gitdir := gitOutput(t, fixture.root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	selectionPath := filepath.Join(gitdir, "prospective-conformance-selection")
	fakeBin := filepath.Join(gitdir, "prospective-fake-bin")
	if err := os.MkdirAll(fakeBin, 0o755); err != nil {
		t.Fatal(err)
	}
	writeGateTestFile(t, fakeBin, "go", fmt.Sprintf(`#!/usr/bin/env bash
case " $* " in
  *" ./internal/conformance "*) printf '%%s|%%s\n' "${%s-}" "${%s-}" >> %q ;;
esac
exit 0
`, registry.ConformanceChecksEnv, registry.ConformanceInheritedEnv, selectionPath), 0o755)
	writeGateTestFile(t, fakeBin, "shellcheck", "#!/usr/bin/env bash\nexit 0\n", 0o755)
	writeGateTestFile(t, fixture.root, prospectiveGatePath, fmt.Sprintf(`#!/usr/bin/env bash
set -euo pipefail
root="${1:?missing prospective root}"
export BENCH_KIT="$root"
export BENCH_PROSPECTIVE_PHASE_HELPER=1
export BENCH_PROSPECTIVE_PHASE_ROOT="$root"
export PATH=%q:"$PATH"
exec %q -test.run '^TestProspectiveFullInventoryHelper$'
`, fakeBin, testBinary), 0o755)

	mustExecuteGreen(t, fixture.root, productionGateEngine{})
	slots, valid := loadConformanceCheckSlots(fixture.root)
	if got, want := len(slots), len(ordinaryConformanceChecks(registry.Dev)); !valid || got != want {
		t.Fatalf("seeded ordinary check slots = %d valid=%v, want %d valid slots", got, valid, want)
	}

	gitRun(t, fixture.root, "add", ".")
	gitRun(t, fixture.root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "seed prospective fixture")
	if err := os.Remove(filepath.Join(fixture.root, filepath.FromSlash(canary.PhaseManifestPath))); err != nil {
		t.Fatal(err)
	}
	gitRun(t, fixture.root, "add", "-u", canary.PhaseManifestPath)
	tree := gitOutput(t, fixture.root, "write-tree")
	gitRun(t, fixture.root, "reset", "--hard", "HEAD")

	var stdout, stderr bytes.Buffer
	if got := ExecuteTree(context.Background(), fixture.root, tree, &stdout, &stderr); got.ActionExit != 0 {
		t.Fatalf("prospective execution = %+v, want green; stdout=%q stderr=%q", got, stdout.String(), stderr.String())
	}
	data, err := os.ReadFile(selectionPath)
	if err != nil {
		t.Fatalf("read prospective conformance selection: %v", err)
	}
	lines := strings.Fields(string(data))
	if got, want := lines, []string{strings.Join(registry.OrdinaryNames(registry.Dev), ",") + "|"}; !slices.Equal(got, want) {
		t.Fatalf("prospective conformance selections = %v, want exactly the full dev inventory %v", got, want)
	}
}

func TestProspectiveFullInventoryHelper(t *testing.T) {
	if os.Getenv("BENCH_PROSPECTIVE_PHASE_HELPER") != "1" {
		return
	}
	root := os.Getenv("BENCH_PROSPECTIVE_PHASE_ROOT")
	if code := PhasesCommand([]string{root}, io.Discard, io.Discard); code != 0 {
		t.Fatalf("prospective PhasesCommand = %d, want green", code)
	}
}

func TestOrdinaryGreenRemainsProspectiveBootstrapEvidence(t *testing.T) {
	root := reusableEvidenceRepo(t, 0)
	writeGateTestFile(t, root, prospectiveGatePath, "#!/usr/bin/env bash\nexit 97\n", 0o755)
	gitRun(t, root, "add", prospectiveGatePath)
	gitRun(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "prospective hook")
	tree := gitOutput(t, root, "write-tree")

	if got := Execute(context.Background(), root, io.Discard, io.Discard); got.ActionExit != 0 {
		t.Fatalf("ordinary execution = %+v, want green", got)
	}
	if got := InspectTree(root, tree); !got.ReusableGreen {
		t.Fatalf("ordinary green did not remain prospective bootstrap evidence: %+v", got)
	}
	if got := ExecuteTree(context.Background(), root, tree, io.Discard, io.Discard); got.ActionExit != 0 {
		t.Fatalf("prospective bootstrap did not reuse ordinary green: %+v", got)
	}
}

func TestExecuteTreeReusesExactGreenBeforeGateLock(t *testing.T) {
	root := reusableEvidenceRepo(t, 0)
	tree := gitOutput(t, root, "write-tree")
	if got := Execute(context.Background(), root, io.Discard, io.Discard); got.ActionExit != 0 {
		t.Fatalf("ordinary execution = %+v, want green", got)
	}
	before := mustRead(t, cachePath(t, root))

	holdGateLock(t, root)

	var stdout bytes.Buffer
	got := ExecuteTree(context.Background(), root, tree, &stdout, io.Discard)
	if got.ActionExit != 0 || !got.Inspection.ReusableGreen {
		t.Fatalf("held-lock prospective reuse = %+v, want reusable green", got)
	}
	if got, want := stdout.String(), "gate: green (fresh verdict reused for this tree)\n"; got != want {
		t.Fatalf("reuse stdout = %q, want %q", got, want)
	}
	if got := gateRunCount(t, root); got != 1 {
		t.Fatalf("gate runs = %d, want 1", got)
	}
	if after := mustRead(t, cachePath(t, root)); !bytes.Equal(before, after) {
		t.Fatalf("held-lock reuse rewrote the verdict record:\nbefore %q\nafter  %q", before, after)
	}
}

func TestNonReusableEvidenceReachesGateLock(t *testing.T) {
	for _, tc := range []struct {
		name    string
		arrange func(t *testing.T, root string)
	}{
		{"stale", func(t *testing.T, root string) {
			plan := mustSubject(t, root)
			replaceRetainedEvidence(t, root, plan, verdictRecord{Schema: 1, State: Ready, Status: "green", Tree: plan.Tree, Oracle: plan.Oracle, RecordedAt: time.Now().UTC().Add(-freshness - time.Minute).Format(time.RFC3339)})
		}},
		{"tree mismatched", func(t *testing.T, root string) {
			plan := mustSubject(t, root)
			replaceRetainedEvidence(t, root, plan, verdictRecord{Schema: 1, State: Ready, Status: "green", Tree: plan.Tree, Oracle: plan.Oracle, RecordedAt: time.Now().UTC().Format(time.RFC3339)})
			writeGateTestFile(t, root, "changed.txt", "changed\n", 0o644)
		}},
		{"oracle mismatched", func(t *testing.T, root string) {
			old, err := buildSubjectForPolicy(root, root, "oracle-v1/freshness-v1")
			if err != nil {
				t.Fatal(err)
			}
			replaceRetainedEvidence(t, root, old, verdictRecord{Schema: 1, State: Ready, Status: "green", Tree: old.Tree, Oracle: old.Oracle, RecordedAt: time.Now().UTC().Format(time.RFC3339)})
		}},
		{"unavailable", func(t *testing.T, root string) {}},
		{"red", func(t *testing.T, root string) {
			plan := mustSubject(t, root)
			replaceRetainedEvidence(t, root, plan, verdictRecord{Schema: 1, State: Ready, Status: "red", Tree: plan.Tree, Oracle: plan.Oracle, RecordedAt: time.Now().UTC().Format(time.RFC3339)})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := reusableEvidenceRepo(t, 0)
			tc.arrange(t, root)
			holdGateLock(t, root)

			var stderr bytes.Buffer
			got := Execute(context.Background(), root, io.Discard, &stderr)
			if got.ActionExit != 1 || got.Inspection.ReusableGreen {
				t.Fatalf("non-reusable execution = %+v, want lock refusal", got)
			}
			if !strings.Contains(stderr.String(), "gate execution already in progress") {
				t.Fatalf("stderr = %q, want gate-lock refusal", stderr.String())
			}
			if got := gateRunCount(t, root); got != 0 {
				t.Fatalf("gate runs = %d, want lock refusal before execution", got)
			}
		})
	}
}

func replaceRetainedEvidence(t *testing.T, root string, plan subject, record verdictRecord) {
	t.Helper()
	gitdir := gitOutput(t, root, "rev-parse", "--absolute-git-dir")
	dir := filepath.Join(gitdir, "bench-gate-evidence")
	if err := os.Mkdir(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := durableReplaceAt(dir, evidenceName(plan), record); err != nil {
		t.Fatal(err)
	}
}

func holdGateLock(t *testing.T, root string) {
	t.Helper()
	gitdir := gitOutput(t, root, "rev-parse", "--absolute-git-dir")
	engine := productionGateEngine{}
	lock, err := engine.OpenLock(filepath.Join(gitdir, "bench-gate.lock"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = engine.Unlock(lock)
		_ = lock.Close()
	})
	if err := engine.Acquire(lock); err != nil {
		t.Fatal(err)
	}
}

type evidenceAfterPreEngine struct {
	productionGateEngine
	root   string
	plan   subject
	record []byte
}

func (e *evidenceAfterPreEngine) OpenLock(path string) (gateFile, error) {
	if err := retainGreen(e.root, e.plan, e.Now()); err != nil {
		return nil, err
	}
	gitdir, err := commonGitDir(e.root)
	if err != nil {
		return nil, err
	}
	e.record, err = os.ReadFile(evidencePath(gitdir, e.plan))
	if err != nil {
		return nil, err
	}
	return e.productionGateEngine.OpenLock(path)
}

func TestEvidenceAppearingAfterPrecheckReusesUnderLock(t *testing.T) {
	root := reusableEvidenceRepo(t, 0)
	plan := mustSubject(t, root)
	engine := &evidenceAfterPreEngine{root: root, plan: plan}

	var stdout bytes.Buffer
	got := executeWithEngine(context.Background(), root, &stdout, io.Discard, engine)
	if got.ActionExit != 0 || !got.Inspection.ReusableGreen {
		t.Fatalf("late evidence execution = %+v, want reusable green", got)
	}
	if got, want := stdout.String(), "gate: green (fresh verdict reused for this tree)\n"; got != want {
		t.Fatalf("reuse stdout = %q, want %q", got, want)
	}
	if got := gateRunCount(t, root); got != 0 {
		t.Fatalf("gate runs = %d, want no execution", got)
	}
	gitdir, err := commonGitDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if after := mustRead(t, evidencePath(gitdir, plan)); !bytes.Equal(engine.record, after) {
		t.Fatalf("late evidence reuse rewrote the verdict record:\nbefore %q\nafter  %q", engine.record, after)
	}
}

func TestExecuteTreeRefusesUnavailableSuppliedTreeWithoutAuthority(t *testing.T) {
	root := reusableEvidenceRepo(t, 0)
	if got := Execute(context.Background(), root, io.Discard, io.Discard); got.ActionExit != 0 {
		t.Fatalf("ordinary seed = %+v, want green", got)
	}
	beforeRuns := gateRunCount(t, root)
	var stderr bytes.Buffer
	result := ExecuteTree(context.Background(), root, strings.Repeat("f", 40), io.Discard, &stderr)
	if result.ActionExit != 1 || result.Inspection.ReusableGreen {
		t.Fatalf("unavailable prospective tree = %+v, want refusal without authority", result)
	}
	if got := gateRunCount(t, root); got != beforeRuns {
		t.Fatalf("gate runs = %d, want %d after unavailable prospective tree", got, beforeRuns)
	}
	if !strings.Contains(stderr.String(), "prospective gate subject unavailable") {
		t.Fatalf("stderr = %q, want unavailable subject refusal", stderr.String())
	}
}

func TestPolicyVersionMismatchInvalidatesGreen(t *testing.T) {
	root := reusableEvidenceRepo(t, 0)
	old, err := buildSubjectForPolicy(root, root, "oracle-v1/freshness-v1")
	if err != nil {
		t.Fatal(err)
	}
	current, err := buildSubject(root)
	if err != nil {
		t.Fatal(err)
	}
	if old.Oracle == current.Oracle {
		t.Fatal("policy version did not change oracle identity")
	}
	gitdir := gitOutput(t, root, "rev-parse", "--absolute-git-dir")
	if err := durableReplace(gitdir, verdictRecord{Schema: 1, State: Ready, Status: "green", Tree: old.Tree, Oracle: old.Oracle, RecordedAt: time.Now().UTC().Format(time.RFC3339)}); err != nil {
		t.Fatal(err)
	}
	if got := Inspect(root); got.ReusableGreen || got.Reason != "oracle changed" {
		t.Fatalf("policy-mismatched green = %+v, want oracle-changed refusal", got)
	}
}

func prospectiveBenchMain(sentinel string) string {
	return fmt.Sprintf(`package main

import (
	"fmt"
	"os"
)

var version string

func main() {
	if len(os.Args) != 3 {
		os.Exit(90)
	}
	switch os.Args[1] {
	case "freshness-check":
		return
	case "gate-phases":
		fmt.Println(%q)
	default:
		os.Exit(91)
	}
}
`, sentinel)
}

func mustReadGateTestFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
