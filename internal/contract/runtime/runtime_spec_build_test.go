package runtime

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
	benchgit "github.com/gibbonmi/bench/internal/git"
)

type runtimeBuildState struct {
	Run          string                            `json:"run"`
	Candidate    string                            `json:"candidate"`
	CandidateTip string                            `json:"candidate_tip"`
	History      []json.RawMessage                 `json:"history"`
	Assignments  map[string]runtimeBuildAssignment `json:"assignments"`
}

type runtimeBuildAssignment struct {
	ID, Path, Base, TicketDigest, Created string
	Rows, Assumptions                     []string
}

func TestRuntimeSpecBuildPorcelainRoutesRealAndLinkedWrappers(t *testing.T) {
	f := setupRuntimeBuildGate(t, "#!/bin/sh\nexit 0\n", contract.WithSpacePath())
	f.WriteFile("specs/demo/spec.md", "# demo\n\nStatus: staged\n")
	f.WriteFile("specs/demo/tickets/one.md", "# One\n\nOwnership fence: internal/demo\n\n- [ ] [R1] first attempt\n")
	f.CommitAll("staged spec")
	f.Bench("gate").RequireExit(0)

	started := f.Bench("spec", "build", "start", "demo")
	started.RequireExit(0)
	contract.RequireContains(t, started.Stdout, "spec_build[1]{slug,state,subject,next}:")
	contract.RequireContains(t, started.Stdout, "demo,active")
	if started.Stderr != "" {
		t.Fatalf("start stderr = %q", started.Stderr)
	}
	firstAssignment := assignRuntimeBuild(t, f, "one.md", "first request")
	first := readRuntimeBuildState(t, f)
	f.Git("update-ref", "refs/bench/recovery/retained", "HEAD").RequireExit(0)
	oldRefs := strings.Fields(f.Git("for-each-ref", "--format=%(refname) %(objectname)", "refs/bench/specbuild/", "refs/bench/recovery/").Stdout)
	plan := f.Bench("spec", "build", "abandon", "demo")
	plan.RequireExit(0)
	lines := strings.Split(plan.Stdout, "\n")
	if len(lines) < 2 {
		t.Fatalf("abandon plan = %q", plan.Stdout)
	}
	fingerprint := strings.Trim(strings.Split(strings.TrimSpace(lines[1]), ",")[0], `"`)
	f.Bench("spec", "build", "abandon", "demo", "--apply", fingerprint).RequireExit(0)
	f.WriteFile("specs/demo/tickets/two.md", "# Two\n\nOwnership fence: internal/demo\n\n- [ ] [R2] second attempt\n")
	f.CommitAll("ticket-only descendant")
	f.Bench("gate").RequireExit(0)
	restarted := f.Bench("spec", "build", "start", "demo")
	restarted.RequireExit(0)
	contract.RequireContains(t, restarted.Stdout, "demo,active")
	secondAssignment := assignRuntimeBuild(t, f, "two.md", "second request")
	second := readRuntimeBuildState(t, f)
	if first.Run == second.Run || first.Candidate == second.Candidate || firstAssignment.ID == secondAssignment.ID || len(second.History) != 1 {
		t.Fatalf("fresh attempt collided: first=%#v second=%#v history=%d", first, second, len(second.History))
	}
	for i := 0; i < len(oldRefs); i += 2 {
		if got := strings.TrimSpace(f.Git("rev-parse", oldRefs[i]).Stdout); got != oldRefs[i+1] {
			t.Fatalf("retained ref %s = %s, want %s", oldRefs[i], got, oldRefs[i+1])
		}
	}

	f.Bench("link").RequireExit(0)
	deep := filepath.Join(f.Root, "nested", "cwd")
	contract.Mkdir(t, deep)
	linked := filepath.Join(f.Root, ".bench", "bin", "bench.sh")
	status := contract.RunAt(t, f, deep, nil, "bash", linked, "spec", "build", "status", "demo", "--full")
	status.RequireExit(0)
	for _, want := range []string{"demo,active", secondAssignment.ID, "review[0]"} {
		contract.RequireContains(t, status.Stdout, want)
	}
	if status.Stderr != "" {
		t.Fatalf("linked status stderr = %q", status.Stderr)
	}
}

func TestRuntimeSpecBuildUsageAndActionableEnvironmentErrors(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	contract.RequireFreshBench(t)
	f := contract.NewFixture(t)
	f.WriteFile("specs/demo/spec.md", "Status: staged\n")
	for _, args := range [][]string{
		{"spec", "build"}, {"spec", "build", "ninth", "demo"},
		{"spec", "build", "commit", "demo"}, {"spec", "build", "update-ref", "demo"},
		{"spec", "build", "worktree", "demo"}, {"spec", "build", "cherry-pick", "demo"},
		{"spec", "build", "assign", "demo", "--ticket", "promote", "--request", "--assignment", "--", "literal"},
		{"spec", "build", "status", "--", "demo"},
	} {
		probe := f.Bench(args...)
		probe.RequireExit(2)
		if !strings.HasPrefix(probe.Stdout, "usage:") || probe.Stderr != "" {
			t.Fatalf("usage %v = stdout %q stderr %q", args, probe.Stdout, probe.Stderr)
		}
	}

	outside := contract.NewFixture(t, contract.WithNoRepo())
	probe := outside.Bench("spec", "build", "status", "demo")
	probe.RequireExit(1)
	contract.RequireContains(t, probe.Stdout, "Git repository unavailable")
	contract.RequireContains(t, probe.Stdout, "run this command inside the working checkout")
	if probe.Stderr != "" {
		t.Fatalf("outside-repo stderr = %q", probe.Stderr)
	}

	missingGate := contract.NewFixture(t)
	missingGate.WriteFile(".gitignore", ".bench-contract-env/\n")
	missingGate.WriteFile("specs/demo/spec.md", "Status: staged\n")
	missingGate.CommitAll("staged spec")
	gateProbe := missingGate.Bench("spec", "build", "start", "demo")
	gateProbe.RequireExit(1)
	contract.RequireContains(t, gateProbe.Stdout, "run bench gate --fresh, then retry start")

	control := contract.NewFixture(t)
	control.WriteFile("specs/control\nslug/spec.md", "Status: staged\n")
	controlProbe := control.Bench("spec", "build", "status", "control\nslug")
	controlProbe.RequireExit(0)
	if strings.Count(controlProbe.Stdout, "\n") != 2 || !strings.Contains(controlProbe.Stdout, `control\nslug`) {
		t.Fatalf("control-bearing output split: %q", controlProbe.Stdout)
	}

	missingKit := t.TempDir()
	offline := contract.RunAt(t, f, f.Root, map[string]string{"BENCH_KIT": missingKit, "BENCH_OFFLINE": "1"}, "bash", filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh"), "spec", "build", "status", "demo")
	if offline.ExitCode == 0 || !strings.Contains(offline.Stdout+offline.Stderr, "bench repair") {
		t.Fatalf("missing binary route = exit %d stdout %q stderr %q", offline.ExitCode, offline.Stdout, offline.Stderr)
	}
}

func TestRuntimeSpecBuildDogfoodsThreeTicketFrontierAndFourthTicketRefill(t *testing.T) {
	f := setupRuntimeBuildGate(t, "#!/bin/sh\nexit 0\n")
	f.WriteFile("specs/demo/spec.md", "Status: staged\n")
	for i := 1; i <= 4; i++ {
		blocked := []string{"", "", "", "Blocked by: one.md\n\n"}[i-1]
		f.WriteFile("specs/demo/tickets/"+[]string{"one", "two", "three", "four"}[i-1]+".md", "# Ticket\n\n"+blocked+"Ownership fence: internal/t"+strconv.Itoa(i)+"\n\n- [ ] [R4"+strconv.Itoa(i)+"] trace\n")
	}
	f.CommitAll("staged frontier")
	f.Bench("gate").RequireExit(0)
	f.Bench("spec", "build", "start", "demo").RequireExit(0)
	first, second, third := assignRuntimeBuild(t, f, "one.md", "frontier-one"), assignRuntimeBuild(t, f, "two.md", "frontier-two"), assignRuntimeBuild(t, f, "three.md", "frontier-three")
	if first.Base != second.Base || first.Base != third.Base || first.Path == second.Path || second.Path == third.Path {
		t.Fatalf("three-ticket frontier did not fill independent sibling slots: %+v %+v %+v", first, second, third)
	}
	writeAssignmentChange(t, first.Path, "internal/t1/change.go", "package t1\n")
	f.Bench("spec", "build", "checkpoint", "demo", "--assignment", first.ID, "--evidence", runtimeCheckpointReceipt(t, f, first, []string{"internal/t1/change.go"})).RequireExit(0)
	f.Bench("spec", "build", "integrate", "demo", "--assignment", first.ID).RequireExit(0)
	fourth := assignRuntimeBuild(t, f, "four.md", "frontier-four")
	for _, active := range []runtimeBuildAssignment{second, third} {
		if _, err := os.Stat(active.Path); err != nil {
			t.Fatalf("fourth ticket waited for active sibling %s to drain: %v", active.ID, err)
		}
	}
	if fourth.Base == first.Base {
		t.Fatalf("newly eligible fourth ticket was not assigned from integrated candidate: %+v", fourth)
	}
}

func TestRuntimeSpecBuildRedRepairGreenRetainsComposedEvidence(t *testing.T) {
	f := setupRuntimeBuildGate(t, "#!/bin/sh\ngitdir=\"$(git rev-parse --path-format=absolute --git-common-dir)\"\nprintf 'run\\n' >> \"$gitdir/spec-gate-count\"\ntest ! -f internal/demo/red.go\n")
	f.WriteFile("specs/demo/spec.md", "# demo\n\nStatus: staged\n")
	f.WriteFile("specs/demo/tickets/red.md", "# Red implementation\n\nOwnership fence: internal/demo\n\n- [ ] [R24] red implementation\n")
	f.WriteFile("specs/demo/tickets/repair.md", "# Repair implementation\n\nOwnership fence: internal/demo\n\n- [ ] [R40] repair implementation\n")
	f.CommitAll("staged spec")
	f.Bench("gate").RequireExit(0)
	gateCount := filepath.Join(gitDir(t, f), "spec-gate-count")
	requireGateRuns(t, gateCount, 1)
	if dirty := f.Git("status", "--porcelain=v1", "--untracked-files=all").Stdout; dirty != "" {
		t.Fatalf("baseline gate dirtied fixture: %q", dirty)
	}
	f.Bench("spec", "build", "start", "demo").RequireExit(0)

	first := assignRuntimeBuild(t, f, "red.md", "red-request")
	writeAssignmentChange(t, first.Path, "internal/demo/red.go", "package demo\n")
	firstReceipt := runtimeCheckpointReceipt(t, f, first, []string{"internal/demo/red.go"})
	f.Bench("spec", "build", "checkpoint", "demo", "--assignment", first.ID, "--evidence", firstReceipt).RequireExit(0)
	requireGateRuns(t, gateCount, 1)
	f.Bench("spec", "build", "integrate", "demo", "--assignment", first.ID).RequireExit(0)
	requireGateRuns(t, gateCount, 1)
	firstReview := runtimeReviewReceipt(t, f)
	f.Bench("spec", "build", "review", "demo", "--evidence", firstReview).RequireExit(0)
	red := f.Bench("spec", "build", "promote", "demo")
	red.RequireExit(1)
	contract.RequireContains(t, red.Stdout, "candidate")
	requireGateRuns(t, gateCount, 2)

	repair := assignRuntimeBuild(t, f, "repair.md", "repair-request")
	if err := os.Remove(filepath.Join(repair.Path, "internal", "demo", "red.go")); err != nil {
		t.Fatal(err)
	}
	writeAssignmentChange(t, repair.Path, "internal/demo/green.go", "package demo\n")
	repairReceipt := runtimeCheckpointReceipt(t, f, repair, []string{"internal/demo/green.go", "internal/demo/red.go"})
	f.Bench("spec", "build", "checkpoint", "demo", "--assignment", repair.ID, "--evidence", repairReceipt).RequireExit(0)
	f.Bench("spec", "build", "integrate", "demo", "--assignment", repair.ID).RequireExit(0)
	requireGateRuns(t, gateCount, 2)
	if stale := f.Bench("spec", "build", "promote", "demo"); stale.ExitCode == 0 || !strings.Contains(stale.Stdout, "review") {
		t.Fatalf("promotion reused stale review: %+v", stale)
	}
	repairReview := runtimeReviewReceipt(t, f)
	f.Bench("spec", "build", "review", "demo", "--evidence", repairReview).RequireExit(0)
	f.Bench("spec", "build", "promote", "demo").RequireExit(0)
	requireGateRuns(t, gateCount, 3)

	full := f.Bench("spec", "build", "status", "demo", "--full")
	full.RequireExit(0)
	for _, want := range []string{"demo,terminal", first.ID, repair.ID, "released", "Standards", "Spec", "Coverage"} {
		contract.RequireContains(t, full.Stdout, want)
	}
	if got := strings.TrimSpace(f.Git("log", "-1", "--format=%P").Stdout); strings.Contains(got, " ") {
		t.Fatalf("promotion exposed provisional merge ancestry: %q", got)
	}
}

// A promoted run's provisional refs are working state, not evidence: the reviewed
// content reaches the branch through the promotion commit, so every assignment branch
// and every candidate and checkpoint ref under refs/bench/specbuild/ is dead the moment
// promotion succeeds. Nothing else enumerates them — release compacts the intent row
// that named the branch — so a promotion that leaves them behind strands them for good.
func TestRuntimeSpecBuildPromotionReclaimsProvisionalRefs(t *testing.T) {
	f := setupRuntimeBuildGate(t, "#!/bin/sh\nexit 0\n")
	f.WriteFile("specs/demo/spec.md", "# demo\n\nStatus: staged\n")
	f.WriteFile("specs/demo/tickets/one.md", "# One\n\nOwnership fence: internal/demo\n\n- [ ] [R61] reclamation\n")
	f.CommitAll("staged spec")
	f.Bench("gate").RequireExit(0)
	f.Bench("spec", "build", "start", "demo").RequireExit(0)

	assigned := assignRuntimeBuild(t, f, "one.md", "reclaim-request")
	writeAssignmentChange(t, assigned.Path, "internal/demo/change.go", "package demo\n")
	receipt := runtimeCheckpointReceipt(t, f, assigned, []string{"internal/demo/change.go"})
	f.Bench("spec", "build", "checkpoint", "demo", "--assignment", assigned.ID, "--evidence", receipt).RequireExit(0)
	f.Bench("spec", "build", "integrate", "demo", "--assignment", assigned.ID).RequireExit(0)
	f.Bench("spec", "build", "review", "demo", "--evidence", runtimeReviewReceipt(t, f)).RequireExit(0)
	f.Bench("spec", "build", "promote", "demo").RequireExit(0)

	for _, namespace := range []string{"refs/heads/bench/assign/", "refs/bench/specbuild/"} {
		if residue := strings.TrimSpace(f.Git("for-each-ref", "--format=%(refname)", namespace).Stdout); residue != "" {
			t.Errorf("promotion stranded provisional refs under %s:\n%s", namespace, residue)
		}
	}
}

func TestRuntimeSpecBuildInterruptStopsProspectiveGateChildren(t *testing.T) {
	f := setupRuntimeBuildGate(t, "#!/bin/sh\ngitdir=\"$(git rev-parse --path-format=absolute --git-common-dir)\"\nif test -f \"$gitdir/spec-block\"; then sleep 30 & echo $! > \"$gitdir/spec-child\"; wait; fi\n")
	f.WriteFile("specs/demo/spec.md", "Status: staged\n")
	f.WriteFile("specs/demo/tickets/one.md", "# One\n\nOwnership fence: internal/demo\n\n- [ ] [R59] interrupt\n")
	f.CommitAll("staged spec")
	f.Bench("gate").RequireExit(0)
	f.Bench("spec", "build", "start", "demo").RequireExit(0)
	assignment := assignRuntimeBuild(t, f, "one.md", "interrupt-request")
	writeAssignmentChange(t, assignment.Path, "internal/demo/change.go", "package demo\n")
	receipt := runtimeCheckpointReceipt(t, f, assignment, []string{"internal/demo/change.go"})
	f.Bench("spec", "build", "checkpoint", "demo", "--assignment", assignment.ID, "--evidence", receipt).RequireExit(0)
	f.Bench("spec", "build", "integrate", "demo", "--assignment", assignment.ID).RequireExit(0)
	f.Bench("spec", "build", "review", "demo", "--evidence", runtimeReviewReceipt(t, f)).RequireExit(0)

	gitdir := gitDir(t, f)
	contract.WriteFileAbs(t, filepath.Join(gitdir, "spec-block"), "block\n")
	launcher := filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh")
	cmd := exec.Command("bash", launcher, "spec", "build", "promote", "demo")
	cmd.Dir, cmd.Env = f.Root, surfaceEnv(f, nil)
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) })
	childText := waitForSurfacePath(t, filepath.Join(gitdir, "spec-child"), cmd)
	child, err := strconv.Atoi(childText)
	if err != nil {
		t.Fatalf("child pid = %q", childText)
	}
	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		if err == nil {
			t.Fatalf("interrupted promote exited zero: %s", output.String())
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("interrupted promote did not exit: %s", output.String())
	}
	deadline := time.Now().Add(3 * time.Second)
	for syscall.Kill(child, 0) == nil && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if syscall.Kill(child, 0) == nil {
		t.Fatalf("prospective gate child %d survived interrupt", child)
	}
}

func setupRuntimeBuildGate(t *testing.T, script string, opts ...contract.FixtureOption) contract.Fixture {
	t.Helper()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RequireFreshBench(t)
	f := contract.NewFixture(t, opts...)
	f.WriteFile(".bench/gate.sh", script)
	if err := os.Chmod(filepath.Join(f.Root, ".bench", "gate.sh"), 0o755); err != nil {
		t.Fatal(err)
	}
	f.WriteFile(".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	f.WriteFile(".gitignore", ".bench-contract-env/\n")
	f.Git("config", "user.email", "bench@local")
	f.Git("config", "user.name", "bench")
	return f
}

func assignRuntimeBuild(t *testing.T, f contract.Fixture, ticket, request string) runtimeBuildAssignment {
	t.Helper()
	before, probe := readRuntimeBuildState(t, f), f.Bench("spec", "build", "assign", "demo", "--ticket", ticket, "--request", request)
	probe.RequireExit(0)
	state := readRuntimeBuildState(t, f)
	for id, assignment := range state.Assignments {
		if before.Assignments[id].ID == "" && assignment.Path != "" && assignment.ID != "" {
			if _, err := os.Stat(assignment.Path); err == nil {
				return assignment
			}
		}
	}
	t.Fatalf("active assignment request %q missing from durable state: %+v", request, state.Assignments)
	return runtimeBuildAssignment{}
}

func writeAssignmentChange(t *testing.T, root, name, body string) {
	t.Helper()
	contract.WriteFileAbs(t, filepath.Join(root, filepath.FromSlash(name)), body)
	command := func(args ...string) {
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
	}
	command("add", "-A")
	command("-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "delegate checkpoint")
}

func runtimeCheckpointReceipt(t *testing.T, f contract.Fixture, assignment runtimeBuildAssignment, ownership []string) string {
	t.Helper()
	tree := benchgit.TreeHash(assignment.Path)
	if tree == "none" {
		t.Fatal("live assignment tree is unavailable")
	}
	rows := make([]map[string]any, len(assignment.Rows))
	for i, row := range assignment.Rows {
		rows[i] = map[string]any{"row": row, "outcome": "passed"}
	}
	assumptions := make([]string, len(assignment.Assumptions))
	for i, assumption := range assignment.Assumptions {
		assumptions[i] = runtimeDigest(assumption)
	}
	sort.Strings(assumptions)
	receipt := map[string]any{
		"version": 1, "run": readRuntimeBuildState(t, f).Run, "assignment": assignment.ID, "base": assignment.Base, "tree": tree,
		"ticket_digest": assignment.TicketDigest, "rows": rows, "checks": []map[string]any{{"name": "focused runtime check", "passed": true}},
		"probe":     map[string]any{"producer": "coordinator", "assignment": assignment.ID, "tree": tree, "command": "focused runtime check", "exit": 0, "output_digest": runtimeDigest("pass"), "produced": time.Now().UTC().Add(time.Second).Format(time.RFC3339Nano)},
		"ownership": ownership, "assumptions": assumptions,
	}
	return writeRuntimeJSON(t, receipt)
}

func requireDirtyAssignmentAtBase(t *testing.T, f contract.Fixture, assignment runtimeBuildAssignment) {
	t.Helper()
	head := strings.TrimSpace(contract.RunAt(t, f, assignment.Path, nil, "git", "rev-parse", "HEAD").Stdout)
	status := contract.RunAt(t, f, assignment.Path, nil, "git", "status", "--porcelain", "--untracked-files=all").Stdout
	if head != assignment.Base || status == "" {
		t.Fatalf("assignment state before checkpoint: head=%s base=%s status=%q", head, assignment.Base, status)
	}
}

func runtimeReviewReceipt(t *testing.T, f contract.Fixture) string {
	t.Helper()
	state := readRuntimeBuildState(t, f)
	return writeRuntimeJSON(t, map[string]any{"version": 1, "run": state.Run, "candidate": state.CandidateTip, "axes": []map[string]any{{"axis": "Standards"}, {"axis": "Spec"}, {"axis": "Coverage"}}})
}

func readRuntimeBuildState(t *testing.T, f contract.Fixture) runtimeBuildState {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(gitDir(t, f), "bench", "specbuild", "*.json"))
	if err != nil || len(paths) != 1 {
		t.Fatalf("spec build state paths = %v, %v", paths, err)
	}
	var state runtimeBuildState
	data, err := os.ReadFile(paths[0])
	if err != nil || json.Unmarshal(data, &state) != nil {
		t.Fatalf("read spec build state: %v", err)
	}
	return state
}

func writeRuntimeJSON(t *testing.T, value any) string {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "receipt.json")
	if err := os.WriteFile(path, append(data, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func runtimeDigest(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func requireGateRuns(t *testing.T, path string, want int) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil || strings.Count(string(data), "run\n") != want {
		t.Fatalf("gate runs = %q, %v; want %d", data, err, want)
	}
}
