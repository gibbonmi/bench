package runtime

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode"

	"github.com/gibbonmi/bench/internal/contract"
)

// runtimeMutationCase is one of the six post-start mutating operations an operator can
// reach through `bench spec build`, paired with the exact CLI args that drive it far
// enough — past every early, operation-specific gate the service checks before the
// shared precondition call — to reach the working-subject refusal under test.
type runtimeMutationCase struct {
	op   string
	args []string
}

// composeRuntimeMutationCases starts one run on f and prepares every piece of live state
// each of the six mutating operations needs to clear its own early gate (a real
// assignment for checkpoint, a checkpointed assignment for integrate, a receipt bound
// to the current candidate for review) so the caller's applicable working-subject
// condition is the first thing the command can still refuse on.
func composeRuntimeMutationCases(t *testing.T, f contract.Fixture) []runtimeMutationCase {
	t.Helper()
	f.WriteFile("specs/demo/spec.md", "# demo\n\nStatus: staged\n")
	f.WriteFile("specs/demo/tickets/one.md", "# One\n\nOwnership fence: internal/t1\n\n- [ ] [R61] checkpoint target\n")
	f.WriteFile("specs/demo/tickets/two.md", "# Two\n\nOwnership fence: internal/t2\n\n- [ ] [R62] integrate target\n")
	f.CommitAll("staged spec")
	f.Bench("gate").RequireExit(0)
	f.Bench("spec", "build", "start", "demo").RequireExit(0)

	checkpointTarget := assignRuntimeBuild(t, f, "one.md", "checkpoint-target-request")
	integrateTarget := assignRuntimeBuild(t, f, "two.md", "integrate-target-request")
	writeAssignmentChange(t, integrateTarget.Path, "internal/t2/change.go", "package t2\n")
	integrateReceipt := runtimeCheckpointReceipt(t, f, integrateTarget, []string{"internal/t2/change.go"})
	f.Bench("spec", "build", "checkpoint", "demo", "--assignment", integrateTarget.ID, "--evidence", integrateReceipt).RequireExit(0)

	reviewReceipt := runtimeReviewReceipt(t, f)

	return []runtimeMutationCase{
		{"assign", []string{"spec", "build", "assign", "demo", "--ticket", "one.md", "--request", "operation-named-assign"}},
		// checkpoint's own early gate only requires the assignment to exist, not to already
		// carry a receipt, so an unvalidated placeholder evidence path is enough to reach the
		// shared precondition call unblocked.
		{"checkpoint", []string{"spec", "build", "checkpoint", "demo", "--assignment", checkpointTarget.ID, "--evidence", "irrelevant-receipt-path"}},
		{"integrate", []string{"spec", "build", "integrate", "demo", "--assignment", integrateTarget.ID}},
		{"review", []string{"spec", "build", "review", "demo", "--evidence", reviewReceipt}},
		{"promote", []string{"spec", "build", "promote", "demo"}},
		// ApplyAbandon reaches the shared precondition call before it ever compares the
		// fingerprint against a real plan, so a placeholder value reaches the same refusal.
		{"abandon", []string{"spec", "build", "abandon", "demo", "--apply", "placeholder-fingerprint"}},
	}
}

// requireOperationNamed asserts message names op (and not "start", unless op is
// "start" itself), between the two literal fragments every workingSubject refusal
// shares — the way the sentence is built independent of which drift condition fired.
func requireOperationNamed(t *testing.T, op, want, message string) {
	t.Helper()
	if !strings.Contains(message, want) {
		t.Fatalf("%s refusal = %q, want it to contain %q", op, message, want)
	}
	if op != "start" && strings.Contains(message, "spec build start") {
		t.Fatalf("%s refusal borrows start's wording: %q", op, message)
	}
}

// TestDirtyCheckoutRefusalNamesEachStrictOperation is CT1: every strict operation names
// itself in the dirty-checkout refusal through the real CLI as a fresh subprocess.
func TestDirtyCheckoutRefusalNamesEachStrictOperation(t *testing.T) {
	t.Run("start", func(t *testing.T) {
		f := setupRuntimeBuildGate(t, "#!/bin/sh\nexit 0\n")
		f.WriteFile("specs/demo/spec.md", "# demo\n\nStatus: staged\n")
		f.WriteFile("specs/demo/tickets/one.md", "# One\n\nOwnership fence: internal/t1\n\n- [ ] [R61] start target\n")
		f.CommitAll("staged spec")
		f.Bench("gate").RequireExit(0)
		f.WriteFile("dirty.txt", "dirty\n")
		probe := f.Bench("spec", "build", "start", "demo")
		probe.RequireExit(1)
		requireOperationNamed(t, "start", "spec build start requires a clean working checkout", probe.Stdout)
		if probe.Stderr != "" {
			t.Fatalf("start dirty-checkout stderr = %q", probe.Stderr)
		}
	})
	f := setupRuntimeBuildGate(t, "#!/bin/sh\nexit 0\n")
	cases := composeRuntimeMutationCases(t, f)
	f.WriteFile("dirty.txt", "dirty\n")
	for _, tc := range cases {
		if tc.op != "promote" && tc.op != "abandon" {
			continue
		}
		t.Run(tc.op, func(t *testing.T) {
			probe := f.Bench(tc.args...)
			probe.RequireExit(1)
			requireOperationNamed(t, tc.op, "spec build "+tc.op+" requires a clean working checkout", probe.Stdout)
			if probe.Stderr != "" {
				t.Fatalf("%s dirty-checkout stderr = %q", tc.op, probe.Stderr)
			}
		})
	}
}

func TestProvisionalOperationsTolerateCoordinatorDirt(t *testing.T) {
	f := setupRuntimeBuildGate(t, "#!/bin/sh\nexit 0\n")
	f.WriteFile("specs/demo/spec.md", "# demo\n\nStatus: staged\n")
	f.WriteFile("specs/demo/tickets/one.md", "# One\n\nOwnership fence: internal/demo\n\n- [ ] [R62] provisional dirt tolerance\n")
	f.WriteFile("coordinator-tracked.txt", "base\n")
	f.WriteFile("coordinator-staged.txt", "base\n")
	f.CommitAll("staged spec")
	f.Bench("gate").RequireExit(0)
	f.Bench("spec", "build", "start", "demo").RequireExit(0)
	f.WriteFile("coordinator-tracked.txt", "changed\n")
	f.WriteFile("coordinator-staged.txt", "changed\n")
	f.Git("add", "coordinator-staged.txt").RequireExit(0)
	f.WriteFile("coordinator-untracked.txt", "changed\n")

	assignment := assignRuntimeBuild(t, f, "one.md", "coordinator-dirt-request")
	writeAssignmentChange(t, assignment.Path, "internal/demo/change.go", "package demo\n")
	f.Bench("spec", "build", "checkpoint", "demo", "--assignment", assignment.ID, "--evidence", runtimeCheckpointReceipt(t, f, assignment, []string{"internal/demo/change.go"})).RequireExit(0)
	f.Bench("spec", "build", "integrate", "demo", "--assignment", assignment.ID).RequireExit(0)
	f.Bench("spec", "build", "review", "demo", "--evidence", runtimeReviewReceipt(t, f)).RequireExit(0)
	status := f.Git("status", "--porcelain=v1", "--untracked-files=all").Stdout
	for _, path := range []string{"coordinator-tracked.txt", "coordinator-staged.txt", "coordinator-untracked.txt"} {
		contract.RequireContains(t, status, path)
	}
}

// TestNoWorkingBranchRefusalNamesEachOperation is CT2: the same six operations name
// themselves in the no-working-branch refusal. The resolver hardcodes two messages —
// fixing only the dirty-checkout one is the cheap half-fix this row exists to catch.
func TestNoWorkingBranchRefusalNamesEachOperation(t *testing.T) {
	f := setupRuntimeBuildGate(t, "#!/bin/sh\nexit 0\n")
	cases := composeRuntimeMutationCases(t, f)
	f.Git("checkout", "--detach", "-q")
	for _, tc := range cases {
		t.Run(tc.op, func(t *testing.T) {
			probe := f.Bench(tc.args...)
			probe.RequireExit(1)
			requireOperationNamed(t, tc.op, "spec build "+tc.op+" requires a checked-out working branch", probe.Stdout)
			if probe.Stderr != "" {
				t.Fatalf("%s no-branch stderr = %q", tc.op, probe.Stderr)
			}
		})
	}
}

// TestReviewReachesPreconditionRefusal is CT3: `review` validates its receipt against
// the run and the current candidate before it ever reaches the shared precondition
// call, so a stub or mismatched receipt returns the receipt refusal. A real three-axis
// receipt makes the no-working-branch refusal the one actually observed.
func TestReviewReachesPreconditionRefusal(t *testing.T) {
	f := setupRuntimeBuildGate(t, "#!/bin/sh\nexit 0\n")
	f.WriteFile("specs/demo/spec.md", "# demo\n\nStatus: staged\n")
	f.WriteFile("specs/demo/tickets/one.md", "# One\n\nOwnership fence: internal/demo\n\n- [ ] [R63] review target\n")
	f.CommitAll("staged spec")
	f.Bench("gate").RequireExit(0)
	f.Bench("spec", "build", "start", "demo").RequireExit(0)
	receipt := runtimeReviewReceipt(t, f)
	f.Git("checkout", "--detach", "-q")
	probe := f.Bench("spec", "build", "review", "demo", "--evidence", receipt)
	probe.RequireExit(1)
	if strings.Contains(probe.Stdout, "invalid spec build review receipt") {
		t.Fatalf("review with a valid candidate-bound receipt still hit the receipt refusal: %q", probe.Stdout)
	}
	requireOperationNamed(t, "review", "spec build review requires a checked-out working branch", probe.Stdout)
}

// TestStartMarkerAncestryEndToEnd is CT4: `start` fast-forwards a strict-ancestor
// `refs/bench/green/<branch>` marker and refuses a divergent one, exercised end to end
// through the real CLI. Both subtests plant the marker with an earlier real `start`,
// then probe a second, unrelated slug on the same branch so the fresh-start path —
// the one that reads the marker unconditionally — is the one under test.
func TestStartMarkerAncestryEndToEnd(t *testing.T) {
	t.Run("ancestor marker fast-forwards and start succeeds", func(t *testing.T) {
		f := setupRuntimeBuildGate(t, "#!/bin/sh\nexit 0\n")
		f.WriteFile("specs/demo/spec.md", "# demo\n\nStatus: staged\n")
		f.WriteFile("specs/demo/tickets/one.md", "# One\n\nOwnership fence: internal/demo\n\n- [ ] [R64] ancestor marker\n")
		f.CommitAll("staged spec")
		f.Bench("gate").RequireExit(0)
		f.Bench("spec", "build", "start", "demo").RequireExit(0)
		branch := strings.TrimSpace(f.Git("symbolic-ref", "--short", "HEAD").Stdout)
		markerAfterFirstStart := strings.TrimSpace(f.Git("rev-parse", "refs/bench/green/"+branch).Stdout)

		f.WriteFile("specs/ancestor-demo/spec.md", "# ancestor-demo\n\nStatus: staged\n")
		f.WriteFile("specs/ancestor-demo/tickets/one.md", "# One\n\nOwnership fence: internal/ancestor\n\n- [ ] [R65] descendant tip\n")
		f.CommitAll("descendant commit")
		newTip := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
		f.Bench("gate").RequireExit(0)

		started := f.Bench("spec", "build", "start", "ancestor-demo")
		started.RequireExit(0)
		if started.Stderr != "" {
			t.Fatalf("ancestor start stderr = %q", started.Stderr)
		}
		markerAfterSecondStart := strings.TrimSpace(f.Git("rev-parse", "refs/bench/green/"+branch).Stdout)
		if markerAfterFirstStart == newTip {
			t.Fatalf("fixture did not advance the tip past the planted marker")
		}
		if markerAfterSecondStart != newTip {
			t.Fatalf("ancestor marker did not fast-forward to the tip: got %s, want %s", markerAfterSecondStart, newTip)
		}
	})

	t.Run("divergent marker refuses start", func(t *testing.T) {
		f := setupRuntimeBuildGate(t, "#!/bin/sh\nexit 0\n")
		f.WriteFile(".root-keep", "root\n")
		f.CommitAll("root commit")
		root := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
		f.WriteFile("specs/demo/spec.md", "# demo\n\nStatus: staged\n")
		f.WriteFile("specs/demo/tickets/one.md", "# One\n\nOwnership fence: internal/demo\n\n- [ ] [R66] divergent marker\n")
		f.CommitAll("staged spec")
		f.Bench("gate").RequireExit(0)
		f.Bench("spec", "build", "start", "demo").RequireExit(0)
		branch := strings.TrimSpace(f.Git("symbolic-ref", "--short", "HEAD").Stdout)
		plantedMarker := strings.TrimSpace(f.Git("rev-parse", "refs/bench/green/"+branch).Stdout)

		f.Git("reset", "--hard", root).RequireExit(0)
		f.WriteFile("specs/sibling-demo/spec.md", "# sibling-demo\n\nStatus: staged\n")
		f.WriteFile("specs/sibling-demo/tickets/one.md", "# One\n\nOwnership fence: internal/sibling\n\n- [ ] [R67] sibling tip\n")
		f.CommitAll("sibling commit")
		siblingTip := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
		if ok := f.GitAllow("merge-base", "--is-ancestor", plantedMarker, siblingTip); ok.ExitCode == 0 {
			t.Fatalf("fixture built an ancestor tip, not a divergent one")
		}
		f.Bench("gate").RequireExit(0)

		probe := f.Bench("spec", "build", "start", "sibling-demo")
		probe.RequireExit(1)
		if !strings.Contains(probe.Stdout, "is not an ancestor of the tip") {
			t.Fatalf("divergent-marker refusal = %q, want the owner's ancestor wording", probe.Stdout)
		}
		if stillMarker := strings.TrimSpace(f.Git("rev-parse", "refs/bench/green/"+branch).Stdout); stillMarker != plantedMarker {
			t.Fatalf("refused start moved the marker: got %s, want unchanged %s", stillMarker, plantedMarker)
		}
	})
}

// TestAbandonRetiresRunOnMovedTipEndToEnd is CT5: `abandon --apply` retires a run whose
// branch tip advanced past its recorded base, end to end through the real CLI, and the
// recovery refs the plan promised survive the cleanup untouched.
func TestAbandonRetiresRunOnMovedTipEndToEnd(t *testing.T) {
	f := setupRuntimeBuildGate(t, "#!/bin/sh\nexit 0\n")
	f.WriteFile("specs/demo/spec.md", "# demo\n\nStatus: staged\n")
	f.WriteFile("specs/demo/tickets/one.md", "# One\n\nOwnership fence: internal/demo\n\n- [ ] [R68] moved-tip abandon\n")
	f.CommitAll("staged spec")
	f.Bench("gate").RequireExit(0)
	f.Bench("spec", "build", "start", "demo").RequireExit(0)
	assignment := assignRuntimeBuild(t, f, "one.md", "moved-tip-request")

	f.Git("update-ref", "refs/bench/recovery/moved-tip", "HEAD").RequireExit(0)
	recoveryObject := strings.TrimSpace(f.Git("rev-parse", "refs/bench/recovery/moved-tip").Stdout)

	f.WriteFile("advance.txt", "advance\n")
	f.CommitAll("advance the working branch")

	plan := f.Bench("spec", "build", "abandon", "demo")
	plan.RequireExit(0)
	lines := strings.Split(plan.Stdout, "\n")
	if len(lines) < 2 {
		t.Fatalf("abandon plan = %q", plan.Stdout)
	}
	fingerprint := strings.Trim(strings.Split(strings.TrimSpace(lines[1]), ",")[0], `"`)

	apply := f.Bench("spec", "build", "abandon", "demo", "--apply", fingerprint)
	apply.RequireExit(0)
	if strings.Contains(apply.Stdout, "bench spec build promote") {
		t.Fatalf("abandon --apply on a moved tip returned the recomposition refusal: %q", apply.Stdout)
	}

	status := f.Bench("spec", "build", "status", "demo")
	status.RequireExit(0)
	contract.RequireContains(t, status.Stdout, "demo,terminal")

	if _, err := os.Stat(assignment.Path); err == nil {
		t.Fatalf("abandon --apply left the assignment worktree behind: %s", assignment.Path)
	}
	if got := strings.TrimSpace(f.Git("rev-parse", "refs/bench/recovery/moved-tip").Stdout); got != recoveryObject {
		t.Fatalf("recovery ref did not survive abandon: got %s, want %s", got, recoveryObject)
	}
}

// TestRefusalSurvivesControlByteBranchName is CT6: a refusal reached while the checked-out
// working branch carries a control byte in its own name never leaks that byte raw. Git's
// own ref-format validator rejects every ASCII control byte (0x00-0x1F, 0x7F) in a branch
// name outright, so the hostile branch here carries U+0085 (NEL), a Unicode control rune
// in the C1 range git's validator does not police — proving the branch is real, not merely
// asserted. Dirtying the checkout with a file name carrying the same control byte is what
// actually exercises escaping: git's own porcelain status quoting is the first line of
// defense, and sanitize.Controls at the CLI boundary is the second: this test pins that the
// composition of the two never regresses to a raw byte reaching the operator.
func TestRefusalSurvivesControlByteBranchName(t *testing.T) {
	f := setupRuntimeBuildGate(t, "#!/bin/sh\nexit 0\n")
	f.WriteFile("specs/demo/spec.md", "# demo\n\nStatus: staged\n")
	f.WriteFile("specs/demo/tickets/one.md", "# One\n\nOwnership fence: internal/demo\n\n- [ ] [R69] control byte branch\n")
	f.CommitAll("staged spec")
	f.Bench("gate").RequireExit(0)
	control := string([]byte{0xc2, 0x85})
	branch := "hostile" + control + "branch"
	f.Git("checkout", "-b", branch, "-q").RequireExit(0)
	f.Bench("gate").RequireExit(0)
	f.Bench("spec", "build", "start", "demo").RequireExit(0)
	f.WriteFile("dirty"+control+"file.txt", "dirty\n")

	probe := f.Bench("spec", "build", "promote", "demo")
	probe.RequireExit(1)
	requireOperationNamed(t, "promote", "spec build promote requires a clean working checkout", probe.Stdout)
	if probe.Stderr != "" {
		t.Fatalf("control-byte-branch refusal leaked to stderr: %q", probe.Stderr)
	}
	for _, r := range probe.Stdout {
		if unicode.IsControl(r) && r != '\n' {
			t.Fatalf("refusal carries a raw control rune %U: %q", r, probe.Stdout)
		}
	}
	if !strings.Contains(probe.Stdout, "dirty") || !strings.Contains(probe.Stdout, "file.txt") {
		t.Fatalf("refusal dropped the hostile entry instead of escaping it: %q", probe.Stdout)
	}
}

// TestSpecBuildRefusalsCrossProcessBoundary is CT7: a fresh `bash bin/bench.sh`
// subprocess — with zero in-memory state of its own — reloads the run record from disk.
// A refusing subprocess that could only produce
// the right operation name and the right run identity by having actually read the
// persisted record proves the record survived the boundary, and a fourth, independent
// subprocess reading it back afterward proves the refusal did not corrupt it.
func TestSpecBuildRefusalsCrossProcessBoundary(t *testing.T) {
	f := setupRuntimeBuildGate(t, "#!/bin/sh\nexit 0\n")
	f.WriteFile("specs/demo/spec.md", "# demo\n\nStatus: staged\n")
	f.WriteFile("specs/demo/tickets/one.md", "# One\n\nOwnership fence: internal/demo\n\n- [ ] [R70] cross process\n")
	f.CommitAll("staged spec")
	f.Bench("gate").RequireExit(0)
	f.Bench("spec", "build", "start", "demo").RequireExit(0)                  // subprocess 1
	assignment := assignRuntimeBuild(t, f, "one.md", "cross-process-request") // subprocess 2
	before := readRuntimeBuildState(t, f)

	f.WriteFile("dirty.txt", "dirty\n")
	refusal := f.Bench("spec", "build", "promote", "demo") // subprocess 3
	refusal.RequireExit(1)
	requireOperationNamed(t, "promote", "spec build promote requires a clean working checkout", refusal.Stdout)

	after := readRuntimeBuildState(t, f) // read directly from disk, independent of subprocess 3's exit
	if after.Run != before.Run || after.Candidate != before.Candidate {
		t.Fatalf("refusing subprocess mutated the persisted run identity: before=%+v after=%+v", before, after)
	}
	if _, ok := after.Assignments[assignment.ID]; !ok {
		if got := len(after.Assignments); got != len(before.Assignments) {
			t.Fatalf("refusing subprocess altered assignment count: before=%d after=%d", len(before.Assignments), got)
		}
	}

	if err := os.Remove(filepath.Join(f.Root, "dirty.txt")); err != nil {
		t.Fatal(err)
	}
	status := f.Bench("spec", "build", "status", "demo", "--full") // subprocess 4, independent of subprocess 3
	status.RequireExit(0)
	contract.RequireContains(t, status.Stdout, assignment.ID)
	contract.RequireContains(t, status.Stdout, "demo,active")
}
