package repository_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
	"github.com/gibbonmi/bench/internal/intent"
)

// The fixture modes differ from the staged policy mode, so a restore that drops the
// snapshot mode shows in the read-back.
const (
	oldPolicyMode  = 0o600
	boardMode      = 0o640
	projectedBoard = "# Roadmap\n\n## Recommended sequence\n\n1. A\n2. B\n"
	decision       = "decision-1"
)

// approval is one committed repository with one recorded plan that adds outcome B.
type approval struct {
	store                 commitrepo.Store
	plan                  commitment.Plan
	policy, board, ledger string
	oldPolicy, oldBoard   string
}

// policyOf is the policy whose active milestone M1 holds outcomes in order.
func policyOf(outcomes ...string) commitment.Policy {
	milestone := commitment.Milestone{ID: "M1"}
	for _, id := range outcomes {
		milestone.Outcomes = append(milestone.Outcomes, commitment.Outcome{ID: id, Criteria: []commitment.Criterion{{ID: id + ".done", Text: "The outcome is delivered."}}})
	}
	return commitment.Policy{Version: 1, ActiveMilestone: "M1", Milestones: []commitment.Milestone{milestone}}
}

// newApproval commits the policy of outcome A and the board, keeps each file that its flag
// names, and records the plan.
func newApproval(t *testing.T, withPolicy, withBoard bool) approval {
	t.Helper()
	root := commitmenttest.Repo(t, policyOf("A"))
	fx := approval{store: commitrepo.Store{Root: root}, policy: filepath.Join(root, filepath.FromSlash(commitment.PolicyPath)), board: filepath.Join(root, "ROADMAP.md")}
	prepare := func(path string, mode os.FileMode, keep bool) string {
		data, err := os.ReadFile(path)
		if err == nil && keep {
			err = os.Chmod(path, mode)
		} else if err == nil {
			err = os.Remove(path)
		}
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	fx.oldPolicy, fx.oldBoard = prepare(fx.policy, oldPolicyMode, withPolicy), prepare(fx.board, boardMode, withBoard)
	if !withPolicy || !withBoard {
		commitmenttest.Commit(t, root, "remove a fixture file")
	}
	data, err := commitment.Bytes(policyOf("A", "B"))
	if err != nil {
		t.Fatal(err)
	}
	if fx.plan, err = fx.store.Plan(data); err != nil {
		t.Fatal(err)
	}
	if fx.ledger, err = intent.Address(root); err != nil {
		t.Fatal(err)
	}
	return fx
}

// at names the ordinal replacement of destination.
func at(destination string, ordinal int) durablefiletest.Fault {
	return durablefiletest.Fault{Destination: destination, Ordinal: ordinal}
}

// faults routes each replacement of one approval by destination, such as policy#2 for the
// policy restore, and after runs its hook once the rename of its replacement succeeds.
type faults struct {
	durablefiletest.Router
	after map[durablefiletest.Fault]func()
}

func (fx approval) faults() *faults {
	names := map[string]string{fx.policy: "policy", fx.board: "board", fx.ledger: "ledger"}
	return &faults{Router: durablefiletest.Router{Fail: map[durablefiletest.Fault]map[string]error{}, Destination: func(path string) string { return names[path] }}, after: map[durablefiletest.Fault]func(){}}
}

// approve runs the approval of fx with given through f. Each rename fails t without the
// ledger lock, because the stage and its restore belong to the transaction.
func (fx approval) approve(t *testing.T, f *faults, given string) error {
	lock, err := intent.LockPath(fx.store.Root)
	if err != nil {
		t.Fatal(err)
	}
	replacer := f.Replacer()
	rename := replacer.Rename
	replacer.Rename = func(oldpath, newpath string) error {
		err := rename(oldpath, newpath)
		fault := f.Trace[len(f.Trace)-1]
		if _, statErr := os.Stat(lock); statErr != nil {
			t.Errorf("%v ran without the ledger lock: %v", fault, statErr)
		}
		if after := f.after[fault]; err == nil && after != nil {
			after()
		}
		return err
	}
	_, err = commitrepo.ApproveWith(fx.store, replacer, fx.plan.ID, given, fx.plan.Effects.Delayed, fx.plan.Effects.Removed)
	return err
}

// fail runs one approval with stage of fault failed, and it returns the approval error and
// the injected cause.
func (fx approval) fail(t *testing.T, fault durablefiletest.Fault, stage string) (err, cause error) {
	f := fx.faults()
	cause = fmt.Errorf("injected %v %s failure", fault, stage)
	f.Fail[fault] = map[string]error{stage: cause}
	return fx.approve(t, f, decision), cause
}

func requireAbsent(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Lstat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("%s exists after approval: %v", path, err)
	}
}

// requireApproved fails t unless the plan's receipt records approved and, when approved,
// the decision.
func (fx approval) requireApproved(t *testing.T, approved bool) {
	t.Helper()
	ledger, err := intent.Read(fx.store.Root)
	if err != nil {
		t.Fatal(err)
	}
	for _, receipt := range ledger.CommitmentReceipts {
		if receipt.Plan != fx.plan.ID {
			continue
		}
		if receipt.Approved != approved || approved && receipt.Decision != decision {
			t.Fatalf("receipt approved=%t decision=%q, want approved=%t", receipt.Approved, receipt.Decision, approved)
		}
		return
	}
	t.Fatalf("no receipt for plan %s", fx.plan.ID)
}

// D23 and D36: the approval returns the policy replacement's classified failure, and a
// failure before the rename leaves the old policy.
func TestCommitmentReplacement(t *testing.T) {
	type row struct {
		stage     string
		published bool
	}
	cases := []row{{stage: "rename"}}
	for _, stage := range durablefiletest.PostRenameStages {
		cases = append(cases, row{stage: stage, published: true})
	}
	for _, testCase := range cases {
		t.Run(testCase.stage, func(t *testing.T) {
			fx := newApproval(t, true, true)

			err, cause := fx.fail(t, at("policy", 1), testCase.stage)

			durablefiletest.RequireFailure(t, err, testCase.stage, testCase.published, cause)
			if !testCase.published {
				durablefiletest.RequireBytes(t, fx.policy, fx.oldPolicy)
			}
		})
	}
}

// D46: the new policy may be visible after its directory sync fails, so the stage
// restores the old policy and its mode before approval returns, and no receipt approves.
func TestApprovalInitialPolicyFailure(t *testing.T) {
	fx := newApproval(t, true, true)

	err, cause := fx.fail(t, at("policy", 1), "directory-sync")

	durablefiletest.RequireFailure(t, err, "directory-sync", true, cause)
	durablefiletest.RequireFile(t, fx.policy, oldPolicyMode, fx.oldPolicy)
	fx.requireApproved(t, false)
}

// D47: a board failure after the policy was staged restores both files, whether or not the
// board rename completed.
func TestApprovalBoardFailure(t *testing.T) {
	for _, testCase := range []struct {
		stage     string
		published bool
	}{{stage: "directory-sync", published: true}, {stage: "rename"}} {
		t.Run(testCase.stage, func(t *testing.T) {
			fx := newApproval(t, true, true)

			err, cause := fx.fail(t, at("board", 1), testCase.stage)

			durablefiletest.RequireFailure(t, err, testCase.stage, testCase.published, cause)
			durablefiletest.RequireFile(t, fx.policy, oldPolicyMode, fx.oldPolicy)
			durablefiletest.RequireFile(t, fx.board, boardMode, fx.oldBoard)
		})
	}
}

// D48: a receipt write that fails before its rename leaves the old ledger, so the
// transaction restores both staged files once under the ledger lock.
func TestApprovalUnpublishedLedgerFailure(t *testing.T) {
	fx := newApproval(t, true, true)
	f := fx.faults()
	cause := errors.New("injected ledger rename failure")
	f.Fail[at("ledger", 1)] = map[string]error{"rename": cause}

	err := fx.approve(t, f, decision)

	durablefiletest.RequireFailure(t, err, "rename", false, cause)
	if got, want := fmt.Sprint(f.Trace), "[policy#1 board#1 ledger#1 policy#2 board#2]"; got != want {
		t.Fatalf("replacements = %s, want %s", got, want)
	}
	durablefiletest.RequireBytes(t, fx.policy, fx.oldPolicy)
	durablefiletest.RequireBytes(t, fx.board, fx.oldBoard)
}

// D30: a refusal before the first write returns the caller's own error, so no
// replacement starts and each regular file keeps its bytes.
func TestApprovalRefusal(t *testing.T) {
	cases := []struct {
		name, decision string
		setup          func(*testing.T, approval)
	}{
		{name: "decision", setup: func(*testing.T, approval) {}},
		{name: "policy-directory", decision: decision, setup: func(t *testing.T, fx approval) { replaceWithDirectory(t, fx.policy) }},
		{name: "board-directory", decision: decision, setup: func(t *testing.T, fx approval) { replaceWithDirectory(t, fx.board) }},
		{name: "board-without-sequence", decision: decision, setup: func(t *testing.T, fx approval) {
			commitmenttest.Write(t, fx.store.Root, "ROADMAP.md", "# Roadmap\n")
		}},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			fx := newApproval(t, true, true)
			testCase.setup(t, fx)
			kept := map[string]string{}
			for _, path := range []string{fx.policy, fx.board} {
				// A directory has no bytes to keep.
				if data, err := os.ReadFile(path); err == nil {
					kept[path] = string(data)
				}
			}
			f := fx.faults()

			err := fx.approve(t, f, testCase.decision)

			var replaceErr *durablefile.Error
			if err == nil || errors.As(err, &replaceErr) {
				t.Fatalf("err = %v, want the caller's refusal", err)
			}
			if len(f.Trace) != 0 {
				t.Fatalf("replacements = %v after a refusal, want none", f.Trace)
			}
			for path, want := range kept {
				durablefiletest.RequireBytes(t, path, want)
			}
		})
	}
}

// replaceWithDirectory turns the regular file at path into a directory.
func replaceWithDirectory(t *testing.T, path string) {
	t.Helper()
	if err := errors.Join(os.Remove(path), os.Mkdir(path, 0o755)); err != nil {
		t.Fatal(err)
	}
}

// D49: the approved receipt may be visible after a failure past the ledger rename, so
// the new policy and board stay with it.
func TestApprovalPublishedLedgerFailure(t *testing.T) {
	for _, stage := range durablefiletest.PostRenameStages {
		t.Run(stage, func(t *testing.T) {
			fx := newApproval(t, true, true)

			err, cause := fx.fail(t, at("ledger", 1), stage)

			durablefiletest.RequireFailure(t, err, stage, true, cause)
			durablefiletest.RequireFile(t, fx.policy, 0o644, string(fx.plan.Proposed))
			durablefiletest.RequireFile(t, fx.board, boardMode, projectedBoard)
			fx.requireApproved(t, true)
		})
	}
}

// D50: a restore that fails, or whose read-back differs from its snapshot, joins its
// cause to the initiating failure and names the file it left unrestored. The restore of
// the other file still runs.
func TestApprovalRestoreFailure(t *testing.T) {
	cases := []struct {
		name                       string
		initial, restore           durablefiletest.Fault
		initialStage, restoreStage string
		tamper                     func(path string) error
	}{
		{name: "stage", initial: at("policy", 1), initialStage: "directory-sync", restore: at("policy", 2), restoreStage: "rename"},
		{name: "compensation", initial: at("ledger", 1), initialStage: "rename", restore: at("board", 2), restoreStage: "directory-sync"},
		{name: "first-of-two", initial: at("ledger", 1), initialStage: "rename", restore: at("policy", 2), restoreStage: "rename"},
		{name: "read-back-mode", initial: at("board", 1), initialStage: "directory-sync", restore: at("board", 2),
			tamper: func(path string) error { return os.Chmod(path, 0o666) }},
		{name: "read-back-bytes", initial: at("board", 1), initialStage: "directory-sync", restore: at("board", 2),
			tamper: func(path string) error { return os.WriteFile(path, []byte("# Tampered\n"), boardMode) }},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			fx := newApproval(t, true, true)
			unrestored, other, otherMode, otherOld := fx.policy, fx.board, os.FileMode(boardMode), fx.oldBoard
			if testCase.restore.Destination == "board" {
				unrestored, other, otherMode, otherOld = fx.board, fx.policy, oldPolicyMode, fx.oldPolicy
			}
			f := fx.faults()
			causes := []error{fmt.Errorf("injected %v failure", testCase.initial)}
			f.Fail[testCase.initial] = map[string]error{testCase.initialStage: causes[0]}
			if testCase.tamper != nil {
				f.after[testCase.restore] = func() {
					if err := testCase.tamper(unrestored); err != nil {
						t.Fatal(err)
					}
				}
			} else {
				causes = append(causes, fmt.Errorf("injected %v failure", testCase.restore))
				f.Fail[testCase.restore] = map[string]error{testCase.restoreStage: causes[1]}
			}

			err := fx.approve(t, f, decision)

			if err == nil || !strings.Contains(err.Error(), "restore "+unrestored+":") || !errors.Is(err, causes[0]) || !errors.Is(err, causes[len(causes)-1]) {
				t.Fatalf("err = %v, want it to hold %v and name the unrestored %s", err, causes, unrestored)
			}
			durablefiletest.RequireFile(t, other, otherMode, otherOld)
		})
	}
}

// D51: without a board, each failure leaves no board, and the policy and receipt follow
// the ledger outcome.
func TestApprovalFailureWithoutBoard(t *testing.T) {
	cases := []struct {
		fault    durablefiletest.Fault
		stage    string
		approved bool
	}{
		{fault: at("policy", 1), stage: "directory-sync"},
		{fault: at("ledger", 1), stage: "rename"},
		{fault: at("ledger", 1), stage: "directory-sync", approved: true},
	}
	for _, testCase := range cases {
		t.Run(testCase.fault.String()+"/"+testCase.stage, func(t *testing.T) {
			fx := newApproval(t, true, false)

			err, cause := fx.fail(t, testCase.fault, testCase.stage)

			if !errors.Is(err, cause) {
				t.Fatalf("err = %v, want %v", err, cause)
			}
			want := fx.oldPolicy
			if testCase.approved {
				want = string(fx.plan.Proposed)
			}
			durablefiletest.RequireBytes(t, fx.policy, want)
			requireAbsent(t, fx.board)
			fx.requireApproved(t, testCase.approved)
		})
	}
}

// A first approval has no policy to restore, so each failure that leaves the old ledger
// removes the staged policy.
func TestApprovalFailureWithoutPolicy(t *testing.T) {
	for _, testCase := range []struct {
		fault durablefiletest.Fault
		stage string
	}{{fault: at("policy", 1), stage: "directory-sync"}, {fault: at("ledger", 1), stage: "rename"}} {
		t.Run(testCase.fault.String(), func(t *testing.T) {
			fx := newApproval(t, false, true)

			err, cause := fx.fail(t, testCase.fault, testCase.stage)

			if !errors.Is(err, cause) {
				t.Fatalf("err = %v, want %v", err, cause)
			}
			requireAbsent(t, fx.policy)
		})
	}
}
