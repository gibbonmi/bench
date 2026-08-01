package worktree

import (
	"errors"
	"fmt"
	"github.com/gibbonmi/bench/internal/bounds"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestReleaseProvisionalRemovesExactRetainedLiveCheckpoint(t *testing.T) {
	fixture := newLiveProvisionalFixture(t)
	created := fixture.created

	if head := gitOutput(t, created.Path, "rev-parse", "HEAD"); head != created.Assignment.Start {
		t.Fatalf("assignment HEAD = %s, want base %s", head, created.Assignment.Start)
	}
	if index := gitOutput(t, created.Path, "write-tree"); index != gitOutput(t, fixture.root, "rev-parse", created.Assignment.Start+"^{tree}") {
		t.Fatalf("assignment index = %s, want base tree", index)
	}
	if err := ReleaseProvisional(fixture.root, fixture.request, created.Path, fixture.evidence); err != nil {
		t.Fatalf("ReleaseProvisional: %v", err)
	}
	if _, err := os.Stat(created.Path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("released worktree remains: %v", err)
	}
	if refs := gitOutput(t, fixture.root, "for-each-ref", "--format=%(refname)", "refs/bench/recovery/"); refs != "" {
		t.Fatalf("exact retained payload created redundant recovery refs: %s", refs)
	}
}

type liveProvisionalFixture struct {
	root, request string
	created       Creation
	evidence      ProvisionalEvidence
}

func newLiveProvisionalFixture(t *testing.T) liveProvisionalFixture {
	t.Helper()
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	mustWrite(t, filepath.Join(root, ".gitignore"), []byte("ignored-release.txt\n"), 0o644)
	gitRun(t, root, "add", ".gitignore")
	gitRun(t, root, "commit", "-qm", "release ignore fixture")
	request := "provisional-live"
	created := mustCreate(t, root, request, "provisional live")
	mustWrite(t, filepath.Join(created.Path, "provisional.txt"), []byte("checkpoint\n"), 0o644)
	tree := benchgit.TreeHash(created.Path)
	checkpoint := gitOutput(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit-tree", tree, "-p", created.Assignment.Start, "-m", "attributed checkpoint")
	integrated := gitOutput(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit-tree", tree, "-p", created.Assignment.Start, "-m", "integrated candidate")
	ref := "refs/bench/test-checkpoint/" + created.Assignment.ID
	integratedRef := "refs/bench/test-candidate/" + created.Assignment.ID
	gitRun(t, root, "update-ref", ref, checkpoint)
	gitRun(t, root, "update-ref", integratedRef, integrated)
	evidence := ProvisionalEvidence{Base: created.Assignment.Start, CheckpointRef: ref, Checkpoint: checkpoint, IntegratedRef: integratedRef, Integrated: integrated}
	return liveProvisionalFixture{root: root, request: request, created: created, evidence: evidence}
}

func TestReleaseProvisionalRefusesLiveCheckpointDrift(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(*testing.T, *liveProvisionalFixture)
	}{
		{"unexplained file", func(t *testing.T, f *liveProvisionalFixture) {
			mustWrite(t, filepath.Join(f.created.Path, "later.txt"), []byte("later\n"), 0o644)
		}},
		{"staged index", func(t *testing.T, f *liveProvisionalFixture) {
			gitRun(t, f.created.Path, "add", "provisional.txt")
		}},
		{"retargeted checkpoint ref", func(t *testing.T, f *liveProvisionalFixture) {
			gitRun(t, f.root, "update-ref", f.evidence.CheckpointRef, f.evidence.Base, f.evidence.Checkpoint)
		}},
		{"retargeted integrated ref", func(t *testing.T, f *liveProvisionalFixture) {
			gitRun(t, f.root, "update-ref", f.evidence.IntegratedRef, f.evidence.Base, f.evidence.Integrated)
		}},
		{"deleted checkpoint ref", func(t *testing.T, f *liveProvisionalFixture) {
			gitRun(t, f.root, "update-ref", "-d", f.evidence.CheckpointRef, f.evidence.Checkpoint)
		}},
		{"checkpoint not based at the assignment base", func(t *testing.T, f *liveProvisionalFixture) {
			parent := gitOutput(t, f.root, "rev-parse", f.evidence.Base+"^")
			tree := gitOutput(t, f.root, "rev-parse", f.evidence.Checkpoint+"^{tree}")
			checkpoint := gitOutput(t, f.root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit-tree", tree, "-p", parent, "-m", "reparented checkpoint")
			gitRun(t, f.root, "update-ref", f.evidence.CheckpointRef, checkpoint, f.evidence.Checkpoint)
			f.evidence.Checkpoint = checkpoint
		}},
		{"nondeclared ignored file", func(t *testing.T, f *liveProvisionalFixture) {
			mustWrite(t, filepath.Join(f.created.Path, "ignored-release.txt"), []byte("ignored\n"), 0o644)
		}},
		{"control byte path", func(t *testing.T, f *liveProvisionalFixture) {
			mustWrite(t, filepath.Join(f.created.Path, "unsafe\npath.txt"), []byte("unsafe\n"), 0o644)
			tree := benchgit.TreeHash(f.created.Path)
			checkpoint := gitOutput(t, f.root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit-tree", tree, "-p", f.evidence.Base, "-m", "hostile checkpoint")
			integrated := gitOutput(t, f.root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit-tree", tree, "-p", f.evidence.Base, "-m", "hostile integration")
			gitRun(t, f.root, "update-ref", f.evidence.CheckpointRef, checkpoint, f.evidence.Checkpoint)
			gitRun(t, f.root, "update-ref", f.evidence.IntegratedRef, integrated, f.evidence.Integrated)
			f.evidence.Checkpoint, f.evidence.Integrated = checkpoint, integrated
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newLiveProvisionalFixture(t)
			tc.mutate(t, &fixture)
			if err := ReleaseProvisional(fixture.root, fixture.request, fixture.created.Path, fixture.evidence); err == nil {
				t.Fatal("drifted provisional release succeeded")
			}
			if _, err := os.Stat(fixture.created.Path); err != nil {
				t.Fatalf("refused provisional release removed checkout: %v", err)
			}
			stored, ok, err := intent.FindAssignmentByRequest(fixture.root, requestDigest(fixture.request))
			if err != nil || !ok || stored.State != intent.StateActive {
				t.Fatalf("refused provisional release mutated assignment: %#v, %v, %v", stored, ok, err)
			}
			if refs := gitOutput(t, fixture.root, "for-each-ref", "--format=%(refname)", "refs/bench/recovery/"); refs != "" {
				t.Fatalf("refused release created recovery refs: %s", refs)
			}
		})
	}
}

func TestReleaseProvisionalDiscardsDeclaredIgnoredOutputWithoutRecovery(t *testing.T) {
	fixture := newLiveProvisionalFixture(t)
	mustMkdirAll(t, filepath.Join(fixture.root, ".bench"), 0o755)
	mustWrite(t, filepath.Join(fixture.root, ".bench", "build-outputs.json"), []byte("{\"schema\":1,\"paths\":[\"ignored-release.txt\"]}\n"), 0o644)
	mustWrite(t, filepath.Join(fixture.created.Path, "ignored-release.txt"), []byte("disposable\n"), 0o644)

	if err := ReleaseProvisional(fixture.root, fixture.request, fixture.created.Path, fixture.evidence); err != nil {
		t.Fatalf("ReleaseProvisional: %v", err)
	}
	if _, err := os.Stat(fixture.created.Path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("released worktree remains: %v", err)
	}
	if refs := gitOutput(t, fixture.root, "for-each-ref", "--format=%(refname)", "refs/bench/recovery/"); refs != "" {
		t.Fatalf("declared output release created redundant recovery refs: %s", refs)
	}
}

func TestCreateStampsAssignment(t *testing.T) {
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	before := time.Now().UTC().Truncate(time.Second)
	creation := mustCreate(t, root, "stamped", "creation stamp")
	after := time.Now().UTC()
	stored, err := assignmentByID(root, creation.Assignment.ID)
	mustNoError(t, err)
	requireTest(t, creation.Assignment.CreatedAt != nil && stored.CreatedAt != nil && *stored.CreatedAt == *creation.Assignment.CreatedAt,
		"created assignment stamp = %v, ledger stamp = %v", creation.Assignment.CreatedAt, stored.CreatedAt)
	stamp, err := time.Parse(time.RFC3339, *stored.CreatedAt)
	mustNoError(t, err)
	requireTest(t, !stamp.Before(before) && !stamp.After(after), "stamp %s outside [%s, %s]", stamp, before, after)
}

// orphanNow is the fixed instant every predicate case is judged against, so no case
// depends on the wall clock.
var orphanNow = time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)

func stampedAt(when time.Time) *string {
	stamp := when.UTC().Format(time.RFC3339)
	return &stamp
}

func TestOrphanedOnAgedRecord(t *testing.T) {
	aged := intent.Assignment{State: intent.StateActive, CreatedAt: stampedAt(orphanNow.Add(-8 * 24 * time.Hour))}
	requireTest(t, orphaned(aged, orphanNow), "orphaned(active, stamped 8 days ago) = false, want true")
}

func TestOrphanedTreatsAbsentStampAsAged(t *testing.T) {
	unstamped := intent.Assignment{State: intent.StateActive}
	requireTest(t, orphaned(unstamped, orphanNow), "orphaned(active, unstamped) = false, want true")
}

func TestOrphanedRequiresAge(t *testing.T) {
	young := intent.Assignment{State: intent.StateActive, CreatedAt: stampedAt(orphanNow.Add(-6 * 24 * time.Hour))}
	requireTest(t, !orphaned(young, orphanNow), "orphaned(active, stamped 6 days ago) = true, want false")
}

// The window's edge is a decision rather than an accident of the comparison: a record
// aged exactly AssignmentStale is still inside the window, and only a strictly older one
// is abandoned.
func TestOrphanedExcludesTheExactWindowEdge(t *testing.T) {
	edge := intent.Assignment{State: intent.StateActive, CreatedAt: stampedAt(orphanNow.Add(-bounds.AssignmentStale))}
	requireTest(t, !orphaned(edge, orphanNow), "orphaned(active, stamped exactly AssignmentStale ago) = true, want false")
}

func TestOrphanedRejectsFutureStamp(t *testing.T) {
	skewed := intent.Assignment{State: intent.StateActive, CreatedAt: stampedAt(orphanNow.Add(30 * 24 * time.Hour))}
	requireTest(t, !orphaned(skewed, orphanNow), "orphaned(active, stamped 30 days ahead) = true, want false")
}

func TestOrphanedOnlyActiveState(t *testing.T) {
	for _, state := range []intent.AssignmentState{intent.StateCleanupPending, intent.StateRecovered, intent.StateComplete} {
		t.Run(string(state), func(t *testing.T) {
			aged := intent.Assignment{State: state, CreatedAt: stampedAt(orphanNow.Add(-8 * 24 * time.Hour))}
			requireTest(t, !orphaned(aged, orphanNow), "orphaned(%s, stamped 8 days ago) = true, want false", state)
		})
	}
}

// An unparseable stamp is unknown age, and unknown must not read as abandoned:
// ValidateAssignment rejects such a record on every ledger read, so reaching here at
// all means the caller bypassed the ledger.
func TestOrphanedRejectsUnparseableStamp(t *testing.T) {
	for _, stamp := range []string{"", "yesterday", "2026-07-\x0027T12:00:00Z"} {
		t.Run(fmt.Sprintf("%q", stamp), func(t *testing.T) {
			bad := intent.Assignment{State: intent.StateActive, CreatedAt: &stamp}
			requireTest(t, !orphaned(bad, orphanNow), "orphaned(active, created_at=%q) = true, want false", stamp)
		})
	}
}

// backdate ages a live assignment by rewriting its ledger stamp, which is how a test
// reaches an aged record through PlanAutomatic — the planner reads the ledger and the
// clock itself, so neither is an argument a caller can set.
func backdate(t *testing.T, root string, assignment intent.Assignment, age time.Duration) {
	t.Helper()
	assignment.CreatedAt = stampedAt(time.Now().Add(-age))
	mustNoError(t, intent.PutAssignment(root, assignment))
}

func TestPlanAutomaticLabelsOrphaned(t *testing.T) {
	root, creation := newOwnedAssignment(t, "orphan-label")
	backdate(t, root, creation.Assignment, 8*24*time.Hour)
	plan, err := PlanAutomatic(root, creation.Path)
	requireTest(t, err == nil && plan.Action == ActionRetain && plan.ReasonCode == ReasonOrphaned,
		"PlanAutomatic over an aged clean assignment = action %q reason %q, %v", plan.Action, plan.ReasonCode, err)
}

func TestPlanAutomaticKeepsEarlierRetainReason(t *testing.T) {
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	mustWrite(t, filepath.Join(root, ".gitignore"), []byte("ignored.txt\n"), 0o644)
	gitRun(t, root, "add", ".gitignore")
	gitRun(t, root, "commit", "-qm", "ignore")
	creation := mustCreate(t, root, "orphan-residue", "ignored residue")
	mustWrite(t, filepath.Join(creation.Path, "ignored.txt"), []byte("residue\n"), 0o644)
	backdate(t, root, creation.Assignment, 8*24*time.Hour)
	plan, err := PlanAutomatic(root, creation.Path)
	requireTest(t, err == nil && plan.Action == ActionRetain && plan.ReasonCode == ReasonIgnored,
		"PlanAutomatic over an aged assignment holding ignored residue = action %q reason %q, %v", plan.Action, plan.ReasonCode, err)
}

// A record carrying no creation stamp is what every worktree cut before the field
// existed looks like, and its git lock was written from the same fields. Release and
// explicit cleanup both re-derive that lock string, so this pins that neither path
// reads the stamp.
func TestReleaseAndPlanExplicitAcceptUnstampedAssignment(t *testing.T) {
	root, creation := newOwnedAssignment(t, "unstamped-lock")
	unstamped := creation.Assignment
	unstamped.CreatedAt = nil
	mustNoError(t, intent.PutAssignment(root, unstamped))
	plan, err := PlanExplicit(root, creation.Path)
	requireTest(t, err == nil && plan.Action == ActionRemove && plan.ReasonCode == "",
		"PlanExplicit over an unstamped assignment = %#v, %v", plan, err)
	code := ReleaseCommand(root, []string{"--request", "landed-unstamped-lock", creation.Path}, io.Discard, io.Discard)
	requireTest(t, code == 0, "ReleaseCommand over an unstamped assignment exit=%d", code)
}

func TestReleaseProvisionalRequiresExactRetainedCleanCheckpoint(t *testing.T) {
	setup := func(t *testing.T) (string, Creation, string, ProvisionalEvidence) {
		root := newWorktreeRepo(t)
		t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
		request := "provisional-" + strings.ReplaceAll(t.Name(), "/", "-")
		created := mustCreate(t, root, request, "provisional")
		mustWrite(t, filepath.Join(created.Path, "provisional.txt"), []byte("checkpoint\n"), 0o644)
		gitRun(t, created.Path, "add", "provisional.txt")
		gitRun(t, created.Path, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "checkpoint")
		tree := gitOutput(t, created.Path, "rev-parse", "HEAD^{tree}")
		checkpoint := gitOutput(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit-tree", tree, "-p", created.Assignment.Start, "-m", "attributed checkpoint")
		integrated := gitOutput(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit-tree", tree, "-p", created.Assignment.Start, "-m", "integrated candidate")
		ref := "refs/bench/test-checkpoint/" + created.Assignment.ID
		integratedRef := "refs/bench/test-candidate/" + created.Assignment.ID
		gitRun(t, root, "update-ref", ref, checkpoint)
		gitRun(t, root, "update-ref", integratedRef, integrated)
		return root, created, request, ProvisionalEvidence{Base: created.Assignment.Start, CheckpointRef: ref, Checkpoint: checkpoint, IntegratedRef: integratedRef, Integrated: integrated}
	}

	t.Run("success and replay", func(t *testing.T) {
		root, created, request, evidence := setup(t)
		if head := gitOutput(t, created.Path, "rev-parse", "HEAD"); head == evidence.Checkpoint {
			t.Fatal("positive control must use a different assignment commit with the checkpoint tree")
		}
		mustNoError(t, ReleaseProvisional(root, request, created.Path, evidence))
		if _, err := os.Stat(created.Path); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("released worktree remains: %v", err)
		}
		mustNoError(t, ReleaseProvisional(root, request, created.Path, evidence))
	})

	for _, tc := range []struct {
		name   string
		mutate func(*testing.T, string, Creation, ProvisionalEvidence)
	}{
		{"retargeted ref", func(t *testing.T, root string, _ Creation, evidence ProvisionalEvidence) {
			gitRun(t, root, "update-ref", evidence.CheckpointRef, "HEAD", evidence.Checkpoint)
		}},
		{"dirty worktree", func(t *testing.T, _ string, created Creation, _ ProvisionalEvidence) {
			mustWrite(t, filepath.Join(created.Path, "dirty.txt"), []byte("dirty\n"), 0o644)
		}},
		{"index only", func(t *testing.T, _ string, created Creation, _ ProvisionalEvidence) {
			mustWrite(t, filepath.Join(created.Path, "index.txt"), []byte("index\n"), 0o644)
			gitRun(t, created.Path, "add", "index.txt")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, created, request, evidence := setup(t)
			tc.mutate(t, root, created, evidence)
			if err := ReleaseProvisional(root, request, created.Path, evidence); err == nil {
				t.Fatal("mismatched provisional release succeeded")
			}
			if _, err := os.Stat(created.Path); err != nil {
				t.Fatalf("refused provisional release removed checkout: %v", err)
			}
			stored, ok, err := intent.FindAssignmentByRequest(root, requestDigest(request))
			if err != nil || !ok || stored.State != intent.StateActive {
				t.Fatalf("refused provisional release mutated assignment: %#v, %v, %v", stored, ok, err)
			}
		})
	}
}

func TestCompactProvisionalAssignmentPropagatesMalformedLedger(t *testing.T) {
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	created := mustCreate(t, root, "provisional-malformed-ledger", "provisional")
	assignment := created.Assignment
	assignment.State = intent.StateComplete
	mustNoError(t, intent.PutAssignment(root, assignment))
	ledger, err := intent.Address(root)
	mustNoError(t, err)
	original, err := os.ReadFile(ledger)
	mustNoError(t, err)
	mustWrite(t, ledger, []byte("{"), 0o600)
	if err := compactProvisionalAssignment(root, assignment.ID); err == nil {
		t.Fatal("malformed assignment ledger was treated as absent")
	}
	mustWrite(t, ledger, original, 0o600)
	stored, ok, err := intent.FindAssignmentByRequest(root, assignment.Request)
	if err != nil || !ok || stored.ID != assignment.ID || stored.State != intent.StateComplete {
		t.Fatalf("failed compaction destroyed recoverable assignment: %#v, %v, %v", stored, ok, err)
	}
}
