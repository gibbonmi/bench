package repository

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/intent"
)

// The fixture modes differ from the staged policy mode, so a restore that drops the
// snapshot mode shows in the read-back.
const (
	oldPolicyMode  = 0o600
	boardMode      = 0o640
	oldBoard       = "# Roadmap\n\n## Recommended sequence\n\n1. old\n"
	projectedBoard = "# Roadmap\n\n## Recommended sequence\n\n1. A\n2. B\n"
)

// approval is one committed repository with one recorded plan that adds outcome B.
type approval struct {
	store                 Store
	plan                  commitment.Plan
	policy, board, ledger string
	oldPolicy             string
}

// newApproval commits a policy, and a board when withBoard holds, and records the plan.
// This internal test builds its own repository because the shared commitment fixtures
// import this package.
func newApproval(t *testing.T, withBoard bool) approval {
	t.Helper()
	root := gittest.RepoOnBranch(t, "main")
	encode := func(outcomes ...string) []byte {
		milestone := commitment.Milestone{ID: "M1"}
		for _, id := range outcomes {
			milestone.Outcomes = append(milestone.Outcomes, commitment.Outcome{ID: id, Criteria: []commitment.Criterion{{ID: id + ".done", Text: "The outcome is delivered."}}})
		}
		data, err := commitment.Bytes(commitment.Policy{Version: 1, ActiveMilestone: "M1", Milestones: []commitment.Milestone{milestone}})
		if err != nil {
			t.Fatal(err)
		}
		return data
	}
	fx := approval{store: Store{Root: root}, policy: filepath.Join(root, filepath.FromSlash(commitment.PolicyPath)), board: filepath.Join(root, "ROADMAP.md"), oldPolicy: string(encode("A"))}
	if err := os.MkdirAll(filepath.Dir(fx.policy), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fx.policy, []byte(fx.oldPolicy), oldPolicyMode); err != nil {
		t.Fatal(err)
	}
	if withBoard {
		if err := os.WriteFile(fx.board, []byte(oldBoard), boardMode); err != nil {
			t.Fatal(err)
		}
	}
	gittest.Output(t, root, "add", "-A")
	gittest.Output(t, root, "commit", "-qm", "initial")
	plan, err := fx.store.Plan(encode("A", "B"))
	if err != nil {
		t.Fatal(err)
	}
	fx.plan = plan
	if fx.ledger, err = intent.Address(root); err != nil {
		t.Fatal(err)
	}
	return fx
}

// faults routes each replacement to the recorder of its destination and that
// destination's ordinal, such as "policy#2" for the policy restore. So a stage fault and
// an independent restore fault fail one approval. It fails t when a replacement runs
// without the ledger lock, because the stage and its restore belong to the transaction.
type faults struct {
	t      *testing.T
	lock   string
	names  map[string]string
	fail   map[string]*durablefiletest.Recorder
	after  map[string]func()
	counts map[string]int
	key    string
	Trace  []string
}

func (fx approval) faults(t *testing.T) *faults {
	return &faults{
		t:     t,
		lock:  fx.ledger + ".lock",
		names: map[string]string{filepath.Dir(fx.policy): "policy", filepath.Dir(fx.board): "board", filepath.Dir(fx.ledger): "ledger"},
		fail:  map[string]*durablefiletest.Recorder{}, after: map[string]func(){}, counts: map[string]int{},
	}
}

// inject fails stage of the replacement that key names.
func (f *faults) inject(key, stage string, cause error) {
	f.fail[key] = &durablefiletest.Recorder{Fail: map[string]error{stage: cause}}
}

func (f *faults) recorder() durablefile.Replacer {
	if rec := f.fail[f.key]; rec != nil {
		return rec.Replacer()
	}
	return durablefile.Native()
}

func (f *faults) replacer() durablefile.Replacer {
	return durablefile.Replacer{
		CreateTemp: func(dir, pattern string) (durablefile.File, error) {
			name := f.names[dir]
			f.counts[name]++
			f.key = fmt.Sprintf("%s#%d", name, f.counts[name])
			f.Trace = append(f.Trace, f.key)
			if _, err := os.Stat(f.lock); err != nil {
				f.t.Errorf("%s ran without the ledger lock: %v", f.key, err)
			}
			return f.recorder().CreateTemp(dir, pattern)
		},
		Rename: func(oldpath, newpath string) error {
			if err := f.recorder().Rename(oldpath, newpath); err != nil {
				return err
			}
			if after := f.after[f.key]; after != nil {
				after()
			}
			return nil
		},
		OpenDirectory: func(name string) (durablefile.Directory, error) { return f.recorder().OpenDirectory(name) },
		Remove:        func(name string) error { return f.recorder().Remove(name) },
	}
}

func (fx approval) approve(f *faults) error {
	_, err := fx.store.approve(f.replacer(), fx.plan.ID, "decision-1", fx.plan.Effects.Delayed, fx.plan.Effects.Removed)
	return err
}

// requireBytes fails t unless the file at path holds want, whatever its mode.
func requireBytes(t *testing.T, path, want string) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Fatalf("%s = %q, want %q", path, got, want)
	}
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
		if receipt.Approved != approved || approved && receipt.Decision != "decision-1" {
			t.Fatalf("receipt approved=%t decision=%q, want approved=%t", receipt.Approved, receipt.Decision, approved)
		}
		return
	}
	t.Fatalf("no receipt for plan %s", fx.plan.ID)
}

// requireUnrestored fails t unless err holds each cause and names path as unrestored.
func requireUnrestored(t *testing.T, err error, path string, causes ...error) {
	t.Helper()
	for _, cause := range causes {
		if !errors.Is(err, cause) {
			t.Fatalf("err = %v, want it to hold %v", err, cause)
		}
	}
	if err == nil || !strings.Contains(err.Error(), "restore "+path+":") {
		t.Fatalf("err = %v, want it to name the unrestored %s", err, path)
	}
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
			fx := newApproval(t, true)
			f := fx.faults(t)
			cause := errors.New("injected policy " + testCase.stage + " failure")
			f.inject("policy#1", testCase.stage, cause)

			err := fx.approve(f)

			durablefiletest.RequireFailure(t, err, testCase.stage, testCase.published, cause)
			if !testCase.published {
				requireBytes(t, fx.policy, fx.oldPolicy)
			}
		})
	}
}

// D46: the new policy may be visible after its directory sync fails, so the stage
// restores the old policy and its mode before approval returns, and no receipt approves.
func TestApprovalInitialPolicyFailure(t *testing.T) {
	fx := newApproval(t, true)
	f := fx.faults(t)
	cause := errors.New("injected policy directory-sync failure")
	f.inject("policy#1", "directory-sync", cause)

	err := fx.approve(f)

	if !errors.Is(err, cause) {
		t.Fatalf("err = %v, want %v", err, cause)
	}
	durablefiletest.RequireFile(t, fx.policy, oldPolicyMode, fx.oldPolicy)
	fx.requireApproved(t, false)
}

// D47: a board failure after the policy was staged restores both files.
func TestApprovalBoardFailure(t *testing.T) {
	fx := newApproval(t, true)
	f := fx.faults(t)
	cause := errors.New("injected board directory-sync failure")
	f.inject("board#1", "directory-sync", cause)

	err := fx.approve(f)

	if !errors.Is(err, cause) {
		t.Fatalf("err = %v, want %v", err, cause)
	}
	durablefiletest.RequireFile(t, fx.policy, oldPolicyMode, fx.oldPolicy)
	durablefiletest.RequireFile(t, fx.board, boardMode, oldBoard)
}

// D48: a receipt write that fails before its rename leaves the old ledger, so the
// transaction restores both staged files once under the ledger lock.
func TestApprovalUnpublishedLedgerFailure(t *testing.T) {
	fx := newApproval(t, true)
	f := fx.faults(t)
	cause := errors.New("injected ledger rename failure")
	f.inject("ledger#1", "rename", cause)

	err := fx.approve(f)

	durablefiletest.RequireFailure(t, err, "rename", false, cause)
	if want := []string{"policy#1", "board#1", "ledger#1", "policy#2", "board#2"}; !slices.Equal(f.Trace, want) {
		t.Fatalf("replacements = %v, want %v", f.Trace, want)
	}
	requireBytes(t, fx.policy, fx.oldPolicy)
	requireBytes(t, fx.board, oldBoard)
}

// D49: the approved receipt may be visible after a failure past the ledger rename, so
// the new policy and board stay with it.
func TestApprovalPublishedLedgerFailure(t *testing.T) {
	for _, stage := range durablefiletest.PostRenameStages {
		t.Run(stage, func(t *testing.T) {
			fx := newApproval(t, true)
			f := fx.faults(t)
			cause := errors.New("injected ledger " + stage + " failure")
			f.inject("ledger#1", stage, cause)

			err := fx.approve(f)

			durablefiletest.RequireFailure(t, err, stage, true, cause)
			durablefiletest.RequireFile(t, fx.policy, 0o644, string(fx.plan.Proposed))
			durablefiletest.RequireFile(t, fx.board, boardMode, projectedBoard)
			fx.requireApproved(t, true)
		})
	}
}

// D50: a restore that fails, or whose read-back differs from its snapshot, joins its
// cause to the initiating failure and names the file it left unrestored.
func TestApprovalRestoreFailure(t *testing.T) {
	cases := []struct {
		name, initial, initialStage, restore, restoreStage string
		tamper                                             bool
	}{
		{name: "stage", initial: "policy#1", initialStage: "directory-sync", restore: "policy#2", restoreStage: "rename"},
		{name: "compensation", initial: "ledger#1", initialStage: "rename", restore: "board#2", restoreStage: "directory-sync"},
		{name: "read-back", initial: "board#1", initialStage: "directory-sync", restore: "board#2", tamper: true},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			fx := newApproval(t, true)
			unrestored := map[string]string{"policy": fx.policy, "board": fx.board}[strings.Split(testCase.restore, "#")[0]]
			f := fx.faults(t)
			causes := []error{errors.New("injected " + testCase.initial + " failure")}
			f.inject(testCase.initial, testCase.initialStage, causes[0])
			if testCase.tamper {
				f.after[testCase.restore] = func() {
					if err := os.Chmod(unrestored, 0o666); err != nil {
						t.Fatal(err)
					}
				}
			} else {
				causes = append(causes, errors.New("injected "+testCase.restore+" failure"))
				f.inject(testCase.restore, testCase.restoreStage, causes[1])
			}

			err := fx.approve(f)

			requireUnrestored(t, err, unrestored, causes...)
		})
	}
}

// D51: without a board, each failure leaves no board, and the policy and receipt follow
// the ledger outcome.
func TestApprovalFailureWithoutBoard(t *testing.T) {
	cases := []struct {
		key, stage string
		approved   bool
	}{
		{key: "policy#1", stage: "directory-sync"},
		{key: "ledger#1", stage: "rename"},
		{key: "ledger#1", stage: "directory-sync", approved: true},
	}
	for _, testCase := range cases {
		t.Run(testCase.key+"/"+testCase.stage, func(t *testing.T) {
			fx := newApproval(t, false)
			f := fx.faults(t)
			cause := errors.New("injected " + testCase.key + " failure")
			f.inject(testCase.key, testCase.stage, cause)

			err := fx.approve(f)

			if !errors.Is(err, cause) {
				t.Fatalf("err = %v, want %v", err, cause)
			}
			want := fx.oldPolicy
			if testCase.approved {
				want = string(fx.plan.Proposed)
			}
			requireBytes(t, fx.policy, want)
			requireAbsent(t, fx.board)
			fx.requireApproved(t, testCase.approved)
		})
	}
}
