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
	"sort"
	"strings"
	"testing"
)

// [OR1] The optimistic exact-green answer is a generation-owned decision: a reuse hit
// derives its subject from one evaluation-owned pre generation — one parsed listing, at
// most one materialization, no post generation, no gate child — and says so.
func TestExecuteReusingFreshGreenAnswersFromEvaluationOwnedGeneration(t *testing.T) {
	// The recorder joins PATH, and PATH is part of the oracle identity, so it must be
	// present for the seed too or the reuse it is counting would refuse on oracle drift.
	recorder := installGateGitRecorder(t)
	root := reusableEvidenceRepo(t, 0)
	t.Setenv("BENCH_KIT", root)
	if got := Execute(context.Background(), root, io.Discard, io.Discard); got.ActionExit != 0 {
		t.Fatalf("seed execution = %+v, want green", got)
	}
	if got := gateRunCount(t, root); got != 1 {
		t.Fatalf("seed gate runs = %d, want 1", got)
	}
	seeded := len(recorder.operations(t))
	var stdout bytes.Buffer
	got := ExecuteReusingFreshGreen(context.Background(), root, &stdout, io.Discard)
	if got.ActionExit != 0 || !got.Inspection.ReusableGreen || !strings.Contains(stdout.String(), "fresh verdict reused") {
		t.Fatalf("reuse result = %+v stdout=%q, want announced reused green", got, stdout.String())
	}
	if got := gateRunCount(t, root); got != 1 {
		t.Fatalf("gate runs after reuse = %d, want the seed run only", got)
	}
	operations := recorder.operations(t)[seeded:]
	if got := countGateGitOperation(operations, "ls-tree"); got != 1 {
		t.Fatalf("reuse-hit parsed listings = %d, want the evaluation-owned pre generation only; operations=%v", got, operations)
	}
	if got := countGateGitOperation(operations, "write-tree"); got > 1 {
		t.Fatalf("reuse-hit working-tree materializations = %d, want at most one; operations=%v", got, operations)
	}
}

// [OR2] A subject with nothing to reuse still pays the real execution: the gate child
// runs and records the verdict it earned.
func TestExecuteReusingFreshGreenFallsThroughToRealExecution(t *testing.T) {
	t.Parallel()
	root := reusableEvidenceRepo(t, 0)
	got := executeReusingFreshGreenAtKit(context.Background(), root, root, io.Discard, io.Discard)
	if got.ActionExit != 0 || got.GateExit != 0 {
		t.Fatalf("fall-through result = %+v, want green real run", got)
	}
	if got := gateRunCount(t, root); got != 1 {
		t.Fatalf("gate runs = %d, want one real execution", got)
	}
	if inspection := Inspect(root); !inspection.ReusableGreen {
		t.Fatalf("fall-through verdict = %+v, want recorded reusable green", inspection)
	}
}

func TestGateEvaluationKeepsAcceptedGenerationForStrippedIdentity(t *testing.T) {
	t.Parallel()
	root := gateTestRepo(t, "#!/usr/bin/env bash\nexit 0\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	evaluation := newWorkingTreeEvaluation(root)
	if _, err := evaluation.acceptPre(); err != nil {
		t.Fatal(err)
	}
	want, err := buildStrippedSubjectForGeneration(root, evaluation.pre)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "work.txt"), []byte("moved\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := evaluation.capturePost(); err != nil {
		t.Fatal(err)
	}
	got, err := evaluation.acceptedStrippedSubject()
	if err != nil {
		t.Fatal(err)
	}
	if got.Tree != want.Tree {
		t.Fatalf("stripped tree after post capture = %s, want accepted pre generation %s", got.Tree, want.Tree)
	}
}

func TestGateEvaluationBoundsProspectiveSourceWorkAndRunsFullComponentInventory(t *testing.T) {
	fixture := newKitShapedFixture(t)
	writeGateTestFile(t, fixture.root, ".bench/gate.sh", `#!/usr/bin/env bash
set -uo pipefail
gitdir="$(git rev-parse --path-format=absolute --git-common-dir)"
echo full >> "$gitdir/full-runs"
for script in .bench/phase-*.sh; do
  bash "$script" || exit 1
done
exec true gate-phases "$PWD"
`, 0o755)
	for _, phase := range fixture.phases {
		writeGateTestFile(t, fixture.root, ".bench/phase-"+phase.Name+".sh", fmt.Sprintf(`#!/usr/bin/env bash
gitdir="$(git rev-parse --path-format=absolute --git-common-dir)"
echo %s >> "$gitdir/prospective-phase-runs"
`, phase.Name), 0o755)
	}
	if err := os.Remove(filepath.Join(fixture.root, prospectiveGatePath)); err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	gitRun(t, fixture.root, "add", ".")
	gitRun(t, fixture.root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "prospective source")
	tree := gitOutput(t, fixture.root, "write-tree")
	checkout, cleanup, err := prospectiveCheckout(fixture.root, tree)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	recorder := installGateGitRecorder(t)
	evaluation := newProspectiveTreeEvaluation(checkout, fixture.root, tree)
	result := executeSubjectWithEngine(context.Background(), checkout, fixture.root, io.Discard, io.Discard, productionGateEngine{}, nil, forceRun, evaluation)
	if result.ActionExit != 0 {
		t.Fatalf("prospective evaluation result = %+v, want green", result)
	}
	if evaluation.pre == nil || evaluation.post == nil || evaluation.pre == evaluation.post || evaluation.pre.tree != tree {
		t.Fatalf("prospective generations = pre:%p/%s post:%p, want supplied pre and distinct checkout post", evaluation.pre, evaluation.pre.tree, evaluation.post)
	}
	operations := recorder.operations(t)
	if got := countGateGitOperation(operations, "write-tree"); got > 1 {
		t.Fatalf("checkout-based working-tree materializations = %d, want at most 1; operations=%v", got, operations)
	}
	if got := countGateGitOperation(operations, "ls-tree"); got != 2 {
		t.Fatalf("prospective parsed listings = %d, want supplied pre and checkout post only; operations=%v", got, operations)
	}
	gitdir := gitOutput(t, fixture.root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	data, err := os.ReadFile(filepath.Join(gitdir, "prospective-phase-runs"))
	if err != nil {
		t.Fatal(err)
	}
	got := strings.Fields(string(data))
	want := fixture.phaseNames()
	sort.Strings(got)
	sort.Strings(want)
	if !slices.Equal(got, want) {
		t.Fatalf("prospective components = %v, want full inventory %v", got, want)
	}
}

func TestGateEvaluationProspectiveValidationRejectsCheckoutDriftWithoutMaterializing(t *testing.T) {
	for _, tc := range []struct {
		name string
		path string
	}{
		{"tracked", "work.txt"},
		{"untracked", "drift.txt"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := reusableEvidenceRepo(t, 0)
			tree := gitOutput(t, root, "write-tree")
			checkout, cleanup, err := prospectiveCheckout(root, tree)
			if err != nil {
				t.Fatal(err)
			}
			defer cleanup()
			recorder := installGateGitRecorder(t)
			evaluation := newProspectiveTreeEvaluation(checkout, root, tree)
			if _, err := evaluation.acceptPre(); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(checkout, tc.path), []byte("drift\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			if _, err := evaluation.validatePre(); err == nil {
				t.Fatal("prospective checkout drift passed under-lock validation")
			}
			operations := recorder.operations(t)
			if got := countGateGitOperation(operations, "write-tree"); got != 0 {
				t.Fatalf("prospective pre/validation materializations = %d, want zero; operations=%v", got, operations)
			}
			if got := countGateGitOperation(operations, "ls-tree"); got != 1 {
				t.Fatalf("prospective pre/validation listings = %d, want accepted pre only; operations=%v", got, operations)
			}
		})
	}
}

func TestGateEvaluationSourceFaultCannotReuseAuthority(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		source func(string) treeSource
	}{
		{"source", func(string) treeSource { return unavailableTreeSource{} }},
		{"listing", func(root string) treeSource { return failingTreeSource{treeSource: workingTreeSource{root: root}} }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newKitShapedFixture(t)
			if got := executeWithEngine(context.Background(), fixture.root, io.Discard, io.Discard, productionGateEngine{}); got.ActionExit != 0 {
				t.Fatalf("seed evaluation = %+v, want green", got)
			}
			evaluation := newWorkingTreeEvaluation(fixture.root)
			evaluation.preSource = tc.source(fixture.root)
			result := executeSubjectWithEngine(context.Background(), fixture.root, fixture.root, io.Discard, io.Discard, productionGateEngine{}, nil, reuseFreshGreen, evaluation)
			if result.ActionExit != 1 || result.Inspection.ReusableGreen || evaluation.pre != nil || evaluation.scoping.eligible {
				t.Fatalf("faulted evaluation = result:%+v pre:%p scoping:%+v, want refusal without reusable or partial authority", result, evaluation.pre, evaluation.scoping)
			}
		})
	}
}

func TestGateEvaluationBlobFaultWidensConformanceInventory(t *testing.T) {
	t.Parallel()
	fixture := seededScopingFixture(t)
	writeGateTestFile(t, fixture.root, "ROADMAP.md", "capture-only edit\n", 0o644)
	evaluation := newWorkingTreeEvaluation(fixture.root)
	evaluation.preSource = &recordingTreeSource{treeSource: workingTreeSource{root: fixture.root}, err: errBlobUnavailable}
	result := executeSubjectWithEngine(context.Background(), fixture.root, fixture.root, io.Discard, io.Discard, productionGateEngine{}, nil, reuseFreshGreen, evaluation)
	if result.ActionExit != 0 {
		t.Fatalf("blob-fault evaluation = %+v, want complete-inventory green", result)
	}
	if len(evaluation.scoping.checks.Inherited) != 0 || evaluation.scoping.checks.Identities != nil || !evaluation.scoping.checks.CanaryFull {
		t.Fatalf("blob-fault conformance authority = %+v, want no inherited or shortened identity map and a full canary", evaluation.scoping.checks)
	}
}

func TestGateEvaluationBoundsOrdinarySourceWorkAcrossAllIdentityFamilies(t *testing.T) {
	fixture := newKitShapedFixture(t)
	recorder := installGateGitRecorder(t)
	evaluation := newWorkingTreeEvaluation(fixture.root)
	result := executeSubjectWithEngine(context.Background(), fixture.root, fixture.root, io.Discard, io.Discard, productionGateEngine{}, nil, forceRun, evaluation)
	if result.ActionExit != 0 {
		t.Fatalf("evaluation result = %+v, want green", result)
	}
	if evaluation.pre == nil || evaluation.post == nil || evaluation.pre == evaluation.post {
		t.Fatalf("generations = pre:%p post:%p, want distinct accepted and post generations", evaluation.pre, evaluation.post)
	}
	postStripped, err := evaluation.postStrippedSubject()
	if err != nil {
		t.Fatalf("post-stripped identity unavailable: %v", err)
	}
	wantPostStripped, err := buildStrippedSubjectForGeneration(fixture.root, evaluation.post)
	if err != nil {
		t.Fatal(err)
	}
	if !postStripped.Closed || !sameSubject(postStripped, wantPostStripped) {
		t.Fatalf("post-stripped identity = %+v, want closed identity %+v from distinct post generation", postStripped, wantPostStripped)
	}
	if evaluation.acceptedSubject.Tree != evaluation.pre.tree || evaluation.acceptedStripped.Tree == "" {
		t.Fatalf("whole/stripped pre identities = whole:%+v stripped:%+v generation:%s", evaluation.acceptedSubject, evaluation.acceptedStripped, evaluation.pre.tree)
	}
	if !evaluation.scoping.eligible || len(evaluation.scoping.identities) == 0 || len(evaluation.scoping.checks.Identities) == 0 || evaluation.scoping.checks.Canary.Shared == "" {
		t.Fatalf("pre identity families incomplete: components=%d checks=%d canary=%+v", len(evaluation.scoping.identities), len(evaluation.scoping.checks.Identities), evaluation.scoping.checks.Canary)
	}
	operations := recorder.operations(t)
	if got := countGateGitOperation(operations, "write-tree"); got > 3 {
		t.Fatalf("working-tree materializations = %d, want at most 3; operations=%v", got, operations)
	}
	if got := countGateGitOperation(operations, "ls-tree"); got != 2 {
		t.Fatalf("parsed listings = %d, want accepted pre and distinct post only; operations=%v", got, operations)
	}
	for object, reads := range gateGitBlobReads(operations) {
		if reads > 1 {
			t.Fatalf("blob %s reads = %d, want at most one in the identity-bearing pre generation", object, reads)
		}
	}
}

type gateGitRecorder struct{ path string }

func installGateGitRecorder(t *testing.T) gateGitRecorder {
	t.Helper()
	realGit, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "operations")
	writeGateTestFile(t, dir, "git", fmt.Sprintf("#!/usr/bin/env bash\nprintf '%%s\\n' \"$*\" >> %q\nexec %q \"$@\"\n", path, realGit), 0o755)
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return gateGitRecorder{path: path}
}

func (r gateGitRecorder) operations(t *testing.T) []string {
	t.Helper()
	data, err := os.ReadFile(r.path)
	if err != nil {
		t.Fatal(err)
	}
	return strings.Split(strings.TrimSpace(string(data)), "\n")
}

func countGateGitOperation(operations []string, operation string) int {
	count := 0
	for _, line := range operations {
		if strings.Contains(" "+line+" ", " "+operation+" ") {
			count++
		}
	}
	return count
}

func gateGitBlobReads(operations []string) map[string]int {
	reads := map[string]int{}
	for _, line := range operations {
		fields := strings.Fields(line)
		for i := range fields {
			if fields[i] == "cat-file" && i+2 < len(fields) && fields[i+1] == "blob" {
				reads[fields[i+2]]++
			}
		}
	}
	return reads
}
