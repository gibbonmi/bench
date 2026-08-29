package worktree

import (
	"encoding/json"
	"errors"
	"github.com/gibbonmi/bench/internal/intent"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPlanAutomaticRejectsEveryInvalidMarkerWithoutMutation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		mutate func(*testing.T, string, string, Marker)
		target func(string) string
	}{
		{"absent", func(t *testing.T, _, marker string, _ Marker) { mustRemove(t, marker) }, nil},
		{"empty", func(t *testing.T, _, marker string, _ Marker) { mustWrite(t, marker, nil, 0o600) }, nil},
		{"malformed", func(t *testing.T, _, marker string, _ Marker) { mustWrite(t, marker, []byte("{"), 0o600) }, nil},
		{"strict-duplicate", rewriteMarkerBody(func(s string) string { return strings.Replace(s, "{", `{"schema":"bench-owner/v1",`, 1) + "\n" }), nil},
		{"strict-unknown", rewriteMarkerBody(func(body string) string { return strings.TrimSuffix(body, "}") + `,"unknown":true}` + "\n" }), nil},
		{"strict-trailing-value", rewriteMarkerBody(func(body string) string { return body + " {}\n" }), nil},
		{"strict-missing-newline", rewriteMarkerBody(func(body string) string { return body }), nil},
		{"strict-wrong-type", rewriteMarkerBody(func(s string) string { return strings.Replace(s, `"`+OwnerMarkerSchema+`"`, "1", 1) + "\n" }), nil},
		{"broad-mode", func(t *testing.T, _, marker string, _ Marker) { mustChmod(t, marker, 0o644) }, nil},
		{"symlink", func(t *testing.T, _, marker string, valid Marker) {
			body, _ := json.Marshal(valid)
			target := marker + ".target"
			mustWrite(t, target, append(body, '\n'), 0o600)
			mustRemove(t, marker)
			mustNoError(t, os.Symlink(target, marker))
		}, nil},
		{"wrong-type", func(t *testing.T, _, marker string, _ Marker) {
			mustRemove(t, marker)
			mustNoError(t, os.Mkdir(marker, 0o700))
		}, nil},
		{"unsupported-schema", rewriteMarker(func(m *Marker) { m.Schema = "bench-owner/v0" }), nil},
		{"invalid-id", rewriteMarker(func(m *Marker) { m.OwnerID = "not-hex" }), nil},
		{"id-mismatch", rewriteMarker(func(m *Marker) { m.OwnerID = strings.Repeat("f", 32) }), nil},
		{"path-mismatch", rewriteMarker(func(m *Marker) { m.Path += "-other" }), nil},
		{"registration-mismatch", func(t *testing.T, path, _ string, _ Marker) {
			copyPath := path + "-unregistered"
			mustMkdirAll(t, copyPath, 0o700)
			gitFile, err := os.ReadFile(filepath.Join(path, ".git"))
			mustNoError(t, err)
			mustWrite(t, filepath.Join(copyPath, ".git"), gitFile, 0o600)
		}, func(path string) string { return path + "-unregistered" }},
		{"path-reused", func(t *testing.T, path, _ string, _ Marker) {
			main := gitOutput(t, path, "rev-parse", "--path-format=absolute", "--git-common-dir")
			mainRoot := filepath.Dir(main)
			gitRun(t, mainRoot, "worktree", "unlock", path)
			gitRun(t, mainRoot, "worktree", "remove", "--force", path)
			mustMkdirAll(t, path, 0o700)
		}, nil},
		{"primary-admin", func(*testing.T, string, string, Marker) {}, func(string) string { return "PRIMARY" }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := newWorktreeRepo(t)
			home := filepath.Join(root, ".bench-home")
			creation := mustCreate(t, root, home, "marker-"+tc.name, "marker validation")
			markerFile, err := markerPath(creation.Path)
			mustNoError(t, err)
			valid := Marker{Schema: OwnerMarkerSchema, OwnerID: creation.Assignment.OwnerID, Path: creation.Path}
			tc.mutate(t, creation.Path, markerFile, valid)
			target := creation.Path
			if tc.target != nil {
				target = tc.target(target)
				if target == "PRIMARY" {
					target = root
				}
			}
			before := lifecycleSnapshot(t, root, creation.Path, markerFile)
			plan, err := PlanAutomatic(root, target)
			mustNoError(t, err)
			requireTest(t, plan.Action == ActionRetain, "PlanAutomatic action = %q, want retain", plan.Action)
			if strings.HasPrefix(tc.name, "strict-") {
				requireTest(t, plan.ReasonCode == ReasonMalformed, "PlanAutomatic reason = %q, want malformed", plan.ReasonCode)
			}
			after := lifecycleSnapshot(t, root, creation.Path, markerFile)
			requireTest(t, before == after, "planner mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
		})
	}
}

func TestCreatePersistsCanonicalPathThroughSymlinkedHome(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	realHome := t.TempDir()
	linkedHome := filepath.Join(t.TempDir(), "bench-home")
	if err := os.Symlink(realHome, linkedHome); err != nil {
		t.Fatal(err)
	}

	creation, err := createAt(defaultJoins(), root, linkedHome, "symlinked-home", "symlinked home", nil, currentTime())
	mustNoError(t, err)
	requireTest(t, creation.Path == creation.Assignment.Worktree,
		"creation path %q differs from assignment path %q", creation.Path, creation.Assignment.Worktree)
}

func TestPlanAutomaticRequiresCompleteAssignmentJoin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		mutate func(*testing.T, string, *intent.Assignment)
	}{
		{"owner", func(_ *testing.T, _ string, a *intent.Assignment) { a.OwnerID = strings.Repeat("e", 32) }},
		{"assignment", func(_ *testing.T, _ string, a *intent.Assignment) { a.ID = strings.Repeat("d", 32) }},
		{"request", func(_ *testing.T, _ string, a *intent.Assignment) { a.Request = strings.Repeat("a", 64) }},
		{"label", func(_ *testing.T, _ string, a *intent.Assignment) { a.Label = "different label" }},
		{"start", func(_ *testing.T, _ string, a *intent.Assignment) { a.Start = strings.Repeat("f", 40) }},
		{"full-branch", func(_ *testing.T, _ string, a *intent.Assignment) { a.Branch += "-other" }},
		{"active-state", func(_ *testing.T, _ string, a *intent.Assignment) { a.State = intent.StateActive }},
		{"recovery", func(_ *testing.T, _ string, a *intent.Assignment) {
			a.Recovery = []intent.Recovery{{Ref: "refs/bench/recovery/" + a.OwnerID + "/" + a.ID + "/1", Root: a.Start, Payloads: []string{a.Start}}}
		}},
		{"current-registration", func(t *testing.T, root string, a *intent.Assignment) {
			gitRun(t, root, "-C", a.Worktree, "switch", "-q", "--detach")
		}},
		{"legacy-only", func(t *testing.T, root string, _ *intent.Assignment) {
			path, err := intent.Address(root)
			mustNoError(t, err)
			mustWrite(t, path, []byte(`{"schema":1,"entries":[]}`+"\n"), 0o600)
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := newWorktreeRepo(t)
			home := filepath.Join(root, ".bench-home")
			creation := mustCreate(t, root, home, "assignment-"+tc.name, "assignment join")
			assignment := creation.Assignment
			assignment.State = intent.StateCleanupPending
			tc.mutate(t, root, &assignment)
			if tc.name != "legacy-only" {
				writeAssignmentLedger(t, root, assignment)
			}
			plan, err := PlanAutomatic(root, creation.Path)
			mustNoError(t, err)
			requireTest(t, plan.Action == ActionRetain, "PlanAutomatic action = %q, want retain", plan.Action)
		})
	}
	t.Run("complete", func(t *testing.T) {
		root := newWorktreeRepo(t)
		home := filepath.Join(root, ".bench-home")
		creation := mustCreate(t, root, home, "assignment-complete", "assignment join")
		assignment := creation.Assignment
		assignment.State = intent.StateCleanupPending
		mustNoError(t, intent.PutAssignment(root, assignment))
		plan, err := PlanAutomatic(root, creation.Path)
		mustNoError(t, err)
		requireTest(t, plan.Action != ActionRetain && plan.Assignment == assignment.ID,
			"complete join plan = %#v, want actionable assignment %s", plan, assignment.ID)
	})
}
func writeAssignmentLedger(t *testing.T, root string, assignment intent.Assignment) {
	t.Helper()
	ledger, err := intent.Read(root)
	mustNoError(t, err)
	ledger.Assignments = []intent.Assignment{assignment}
	body, err := json.Marshal(ledger)
	mustNoError(t, err)
	path, err := intent.Address(root)
	mustNoError(t, err)
	mustWrite(t, path, append(body, '\n'), 0o600)
}
func TestConcurrentCreateSerializesByRequest(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	attempted, registered, proceed := make(chan string, 8), make(chan struct{}), make(chan struct{})
	j := defaultJoins()
	j.creationLockAttempt = func(request string) { attempted <- request }
	type result struct {
		creation Creation
		err      error
	}
	results := make(chan result, 2)
	go func() {
		created, err := createAt(j, root, home, "same-request", "same label", func(step LifecycleStep) error {
			if step == StepRegistration {
				close(registered)
				<-proceed
			}
			return nil
		}, currentTime())
		results <- result{created, err}
	}()
	<-attempted
	<-registered
	go func() {
		created, err := createAt(j, root, home, "same-request", "same label", nil, currentTime())
		results <- result{created, err}
	}()
	<-attempted
	close(proceed)
	first, second := <-results, <-results
	requireTest(t, first.err == nil && second.err == nil && first.creation.Path == second.creation.Path && first.creation.Assignment.ID == second.creation.Assignment.ID,
		"concurrent create = %#v/%v and %#v/%v", first.creation, first.err, second.creation, second.err)
	_, err := createAt(j, root, home, "same-request", "changed label", nil, currentTime())
	requireTest(t, err != nil, "changed-label replay did not conflict")
	assignments, err := intent.Assignments(root)
	requireTest(t, err == nil && len(assignments) == 1 && strings.Count(gitOutput(t, root, "worktree", "list", "--porcelain"), "worktree ") == 2,
		"serialized bundle = %#v, %v", assignments, err)
	start := make(chan struct{})
	for _, request := range []string{"distinct-a", "distinct-b"} {
		request := request
		go func() {
			<-start
			created, err := createAt(j, root, home, request, request, nil, currentTime())
			results <- result{created, err}
		}()
	}
	close(start)
	distinctA, distinctB := <-results, <-results
	requireTest(t, distinctA.err == nil && distinctB.err == nil && distinctA.creation.Path != distinctB.creation.Path,
		"distinct concurrent create = %#v / %#v", distinctA, distinctB)
	assignments, err = intent.Assignments(root)
	requireTest(t, err == nil && len(assignments) == 3, "distinct request assignments = %#v, %v", assignments, err)
}
func TestActualSIGINTAfterUnlockRestoresExactLockAndReplays(t *testing.T) {
	t.Parallel()
	if os.Getenv("BENCH_SIGINT_HELPER") == "1" {
		_, err := ApplyAutomatic(os.Getenv("BENCH_SIGINT_ROOT"), os.Getenv("BENCH_SIGINT_PATH"), nil)
		requireTest(t, errors.Is(err, ErrCleanupInterrupted), "cleanup error = %v, want distinct interruption", err)
		return
	}
	root, creation, _ := newPendingAssignment(t, "actual-sigint")
	shimDir := t.TempDir()
	realGit, err := exec.LookPath("git")
	mustNoError(t, err)
	shim := `#!/bin/sh
prev=
for arg do
  if [ "$prev" = worktree ] && [ "$arg" = remove ]; then
    kill -INT "$PPID"
    while :; do sleep 1; done
  fi
  prev=$arg
done
exec "$REAL_GIT" "$@"
`
	mustWrite(t, filepath.Join(shimDir, "git"), []byte(shim), 0o755)
	cmd := descendant(t, os.Args[0], "-test.run=^TestActualSIGINTAfterUnlockRestoresExactLockAndReplays$")
	cmd.Env = append(os.Environ(), "BENCH_SIGINT_HELPER=1", "BENCH_SIGINT_ROOT="+root, "BENCH_SIGINT_PATH="+creation.Path, "REAL_GIT="+realGit, "PATH="+shimDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	out, err := cmd.CombinedOutput()
	requireTest(t, err == nil, "SIGINT helper: %v\n%s", err, out)
	err = validateCreationBundle(root, creation.Assignment)
	requireTest(t, err == nil, "SIGINT residual is not exact-locked: %v", err)
	replay, err := ApplyAutomatic(root, creation.Path, nil)
	requireTest(t, err == nil && replay.Action == ActionRemoved, "SIGINT replay = %#v, %v", replay, err)
}
func TestLifecycleFaultBoundariesRemainLockedOrAbsent(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name                               string
		step                               LifecycleStep
		removeFault, relockFault, survives bool
	}{
		{"registration", StepRegistration, false, false, false},
		{"marker", StepMarker, false, false, false},
		{"record", StepRecord, false, false, false},
		{"rollback-remove", StepRegistration, true, false, true},
		{"rollback-relock", StepRegistration, true, true, true},
	} {
		t.Run("create-"+tc.name, func(t *testing.T) {
			root := newWorktreeRepo(t)
			home := filepath.Join(root, ".bench-home")
			fault := errors.New("fault after " + string(tc.step))
			creation, err := createAt(defaultJoins(), root, home, "fault-"+tc.name, "fault creation", func(got LifecycleStep) error {
				if got == tc.step {
					return fault
				}
				if got == StepRollbackRemove && tc.removeFault {
					return errors.New("injected rollback removal failure")
				}
				if got == StepRelock && tc.relockFault {
					return errors.New("injected relock verification failure")
				}
				return nil
			}, currentTime())
			requireTest(t, errors.Is(err, fault) && creation.Path == "",
				"Create = %#v, %v; want no returned path and injected error", creation, err)
			registrations := gitOutput(t, root, "worktree", "list", "--porcelain")
			assignments, readErr := intent.Assignments(root)
			requireTest(t, tc.survives || strings.Count(registrations, "worktree ") == 1 && len(assignments) == 0,
				"rollback residue registrations=%s assignments=%#v error=%v", registrations, assignments, readErr)
			if tc.survives {
				valid := readErr == nil && len(assignments) == 1
				if valid {
					valid = validateCreationBundle(root, assignments[0]) == nil
				}
				requireTest(t, valid, "surviving rollback bundle registrations=%s assignments=%#v error=%v", registrations, assignments, readErr)
				requireTest(t, !tc.relockFault || strings.Contains(err.Error(), "HIGH-SEVERITY"), "relock residual lacks severity: %v", err)
			}
		})
	}
	// Every step reachable from the automatic path. StepRecoveryRef is not among
	// them: the automatic planner retains a checkout it could only remove by
	// preserving first. The explicit-retry cases below are where that boundary
	// is graded.
	for _, tc := range []struct {
		name       string
		step       LifecycleStep
		wantExists bool
		wantBranch bool
	}{
		{"unlock", StepUnlock, true, true},
		{"sigint-at-unlock", StepUnlock, true, true},
		{"removal-command", StepRemovalAttempt, true, true},
		{"after-removal", StepRemoval, false, true},
		{"after-branch-removal", StepBranch, false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, creation, _ := newOwnedAssignment(t, "fault-"+tc.name)
			markPending(t, root, creation.Assignment)
			fault := errors.New("fault " + tc.name)
			_, err := ApplyAutomatic(root, creation.Path, failLifecycleStep(tc.step, fault))
			requireTest(t, errors.Is(err, fault), "ApplyAutomatic error = %v, want %v", err, fault)
			_, statErr := os.Stat(creation.Path)
			requireTest(t, (statErr == nil) == tc.wantExists, "checkout existence = %v, want %v", statErr == nil, tc.wantExists)
			branchExists := descendant(t, "git", "-C", root, "show-ref", "--verify", "--quiet", creation.Assignment.Branch).Run() == nil
			requireTest(t, branchExists == tc.wantBranch, "assignment branch existence = %v, want %v", branchExists, tc.wantBranch)
			assignments, readErr := intent.Assignments(root)
			requireTest(t, readErr == nil && len(assignments) == 1 && assignments[0].State == intent.StateCleanupPending,
				"fault assignment = %#v, %v", assignments, readErr)
			if tc.wantExists {
				registration := gitOutput(t, root, "worktree", "list", "--porcelain")
				requireTest(t, strings.Contains(registration, "locked "+lockReason(assignments[0])), "surviving checkout is not re-locked:\n%s", registration)
			}
		})
	}
	for _, step := range []LifecycleStep{StepReceipt, StepRecoveryRef, StepUnlock, StepRemoval, StepBranch} {
		t.Run("explicit-retry-"+string(step), func(t *testing.T) {
			root, creation, _ := newOwnedAssignment(t, "explicit-retry-"+string(step))
			if step != StepRecoveryRef {
				markPending(t, root, creation.Assignment)
			}
			mustWrite(t, filepath.Join(creation.Path, "dirty.txt"), []byte("recover once\n"), 0o644)
			plan, err := PlanExplicit(root, creation.Path)
			mustNoError(t, err)
			fault := errors.New("interrupt at " + string(step))
			faulted := defaultJoins()
			faulted.cleanupBoundary = failLifecycleStep(step, fault)
			_, err = applyExplicitWith(faulted, root, creation.Path, plan.Fingerprint, CleanupOptions{})
			requireTest(t, errors.Is(err, fault), "first apply error = %v, want %v", err, fault)
			if step == StepRecoveryRef {
				registration := gitOutput(t, root, "worktree", "list", "--porcelain")
				requireTest(t, strings.Contains(registration, "worktree "+creation.Path) && strings.Contains(registration, "locked "+lockReason(creation.Assignment)),
					"recovery-ref failure lost exact locked checkout:\n%s", registration)
			}
			replay, err := ApplyExplicit(root, creation.Path, plan.Fingerprint)
			requireTest(t, err == nil && replay.Action == ActionRemoved, "interrupted retry = %#v, %v", replay, err)
			refs := strings.Fields(gitOutput(t, root, "for-each-ref", "--format=%(refname)", "refs/bench/recovery/"))
			requireTest(t, len(refs) == 1, "interrupted retry refs = %#v", refs)
		})
	}
}
func rewriteMarker(edit func(*Marker)) func(*testing.T, string, string, Marker) {
	return func(t *testing.T, _ string, markerFile string, marker Marker) {
		edit(&marker)
		body, err := json.Marshal(marker)
		mustNoError(t, err)
		mustWrite(t, markerFile, append(body, '\n'), 0o600)
	}
}
func rewriteMarkerBody(edit func(string) string) func(*testing.T, string, string, Marker) {
	return func(t *testing.T, _ string, markerFile string, marker Marker) {
		body, err := json.Marshal(marker)
		mustNoError(t, err)
		mustWrite(t, markerFile, []byte(edit(string(body))), 0o600)
	}
}
func lifecycleSnapshot(t *testing.T, root, path, marker string) string {
	t.Helper()
	parts := []string{
		gitOutput(t, root, "worktree", "list", "--porcelain"),
		gitOutput(t, root, "for-each-ref", "--format=%(refname) %(objectname)", "refs/heads/", "refs/bench/"),
	}
	ledger, _ := os.ReadFile(filepath.Join(root, ".git", intent.Filename))
	parts = append(parts, string(ledger))
	if info, err := os.Lstat(marker); err == nil {
		parts = append(parts, info.Mode().String())
		if info.Mode()&os.ModeSymlink != 0 {
			target, _ := os.Readlink(marker)
			parts = append(parts, target)
		} else if info.Mode().IsRegular() {
			body, _ := os.ReadFile(marker)
			parts = append(parts, string(body))
		}
	} else {
		parts = append(parts, err.Error())
	}
	if _, err := os.Lstat(path); err == nil {
		parts = append(parts, "path-present")
	} else {
		parts = append(parts, "path-absent")
	}
	return strings.Join(parts, "\n--\n")
}
func mustWrite(t *testing.T, path string, data []byte, mode os.FileMode) {
	t.Helper()
	mustNoError(t, os.WriteFile(path, data, mode))
	mustNoError(t, os.Chmod(path, mode))
}
func mustMkdirAll(t *testing.T, path string, mode os.FileMode) {
	t.Helper()
	mustNoError(t, os.MkdirAll(path, mode))
}

func failLifecycleStep(want LifecycleStep, failure error) func(LifecycleStep) error {
	return func(got LifecycleStep) error {
		if got == want {
			return failure
		}
		return nil
	}
}
