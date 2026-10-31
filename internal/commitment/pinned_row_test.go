package commitment_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/gittest"
)

// pinnedEdit is a row body that no published row holds.
const pinnedEdit = "**FT1 — A**\n\nRecord the settled A decision.\n"

// pinnedRows commits on main the protected milestone whose outcomes A then B pin FT1 and FT2,
// and returns the root, one owned assignment worktree of it, and the published policy.
func pinnedRows(t *testing.T) (root, worktree string, published commitment.Policy) {
	t.Helper()
	root = gittest.RepoOnBranch(t, "main")
	commitmenttest.SeedProtected(t, root)
	commitmenttest.Commit(t, root, "protect the rows")
	published, _, err := commitrepo.Store{Root: root}.Policy()
	if err != nil {
		t.Fatal(err)
	}
	return root, commitmenttest.Assignment(t, root, "pinned-row"), published
}

// repin returns policy with the first source of outcome index bound at the identity of body.
func repin(policy commitment.Policy, index int, body string) commitment.Policy {
	outcomes := append([]commitment.Outcome(nil), policy.Milestones[0].Outcomes...)
	outcomes[index].Sources = []commitment.SourceBinding{outcomes[index].Sources[0]}
	outcomes[index].Sources[0].Identity = commitment.Identity([]byte(body))
	policy.Milestones = []commitment.Milestone{policy.Milestones[0]}
	policy.Milestones[0].Outcomes = outcomes
	return policy
}

// changedFrom is the refusal of the source of outcome index in proposed whose read bytes are got.
func changedFrom(proposed commitment.Policy, index int, got string) string {
	source := proposed.Milestones[0].Outcomes[index].Sources[0]
	return "commitment source " + strconv.Quote(source.ID) + " changed from " + source.Identity + " to " + commitment.Identity([]byte(got))
}

// Plan accepts the new identity of each pinned row whose working file holds the new bytes
// while HEAD keeps the published bytes. One row, both rows, and a store root through a linked
// parent directory each plan.
func TestPinnedRowPlanReadsTheWorkingFile(t *testing.T) {
	const second = "**FT2 — B**\n\nRecord the settled B decision.\n"
	for _, row := range []struct {
		name   string
		bodies []string
		linked bool
	}{
		{name: "one-row", bodies: []string{pinnedEdit}},
		{name: "two-rows", bodies: []string{pinnedEdit, second}},
		{name: "linked-parent", bodies: []string{pinnedEdit}, linked: true},
	} {
		t.Run(row.name, func(t *testing.T) {
			_, worktree, proposed := pinnedRows(t)
			for i, body := range row.bodies {
				commitmenttest.Write(t, worktree, proposed.Milestones[0].Outcomes[i].Sources[0].Path, body)
				proposed = repin(proposed, i, body)
			}
			storeRoot := worktree
			if row.linked {
				link := filepath.Join(t.TempDir(), "linked")
				if err := os.Symlink(filepath.Dir(worktree), link); err != nil {
					t.Fatal(err)
				}
				storeRoot = filepath.Join(link, filepath.Base(worktree))
			}
			if _, err := (commitrepo.Store{Root: storeRoot}).Plan(mustBytes(t, proposed)); err != nil {
				t.Fatalf("Plan(the new identity of %d working rows) = %v, want acceptance", len(row.bodies), err)
			}
		})
	}
}

// Plan refuses a new identity that the working file does not hold, and the refusal writes
// nothing. The working file holds other bytes, or HEAD holds the new bytes while the working
// file keeps the published bytes.
func TestPinnedRowPlanRefusesBytesTheWorkingFileLacks(t *testing.T) {
	const other = "**FT1 — A**\n\nBytes that the plan does not name.\n"
	for _, row := range []struct {
		name string
		hold func(t *testing.T, worktree, path, published string) string
	}{
		{name: "unwritten", hold: func(t *testing.T, worktree, path, _ string) string {
			commitmenttest.Write(t, worktree, path, other)
			return other
		}},
		{name: "head-only", hold: func(t *testing.T, worktree, path, published string) string {
			commitmenttest.Write(t, worktree, path, pinnedEdit)
			commitmenttest.Commit(t, worktree, "commit the row edit")
			commitmenttest.Write(t, worktree, path, published)
			return published
		}},
	} {
		t.Run(row.name, func(t *testing.T) {
			root, worktree, current := pinnedRows(t)
			path := current.Milestones[0].Outcomes[0].Sources[0].Path
			data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
			if err != nil {
				t.Fatal(err)
			}
			published := string(data)
			got := row.hold(t, worktree, path, published)
			proposed := repin(current, 0, pinnedEdit)
			before := commitmenttest.MilestoneState(t, root, worktree)
			_, err = (commitrepo.Store{Root: worktree}).Plan(mustBytes(t, proposed))
			want := changedFrom(proposed, 0, got)
			if err == nil || err.Error() != want || commitmenttest.MilestoneState(t, root, worktree) != before {
				t.Fatalf("Plan = %v, want %q that writes nothing", err, want)
			}
		})
	}
}

// Approve stages the proposal while the row edit stays uncommitted in the working file, and
// it refuses a plan whose working row changed after the plan.
func TestPinnedRowApproveRereadsTheWorkingFile(t *testing.T) {
	const later = "**FT1 — A**\n\nAn edit after the plan.\n"
	for _, row := range []struct {
		name  string
		after string
	}{
		{name: "kept-after-plan"},
		{name: "changed-after-plan", after: later},
	} {
		t.Run(row.name, func(t *testing.T) {
			_, worktree, current := pinnedRows(t)
			path := current.Milestones[0].Outcomes[0].Sources[0].Path
			commitmenttest.Write(t, worktree, path, pinnedEdit)
			proposed := repin(current, 0, pinnedEdit)
			input := mustBytes(t, proposed)
			store := commitrepo.Store{Root: worktree}
			plan, err := store.Plan(input)
			if err != nil {
				t.Fatal(err)
			}
			if row.after != "" {
				commitmenttest.Write(t, worktree, path, row.after)
			}
			policyFile := filepath.Join(worktree, filepath.FromSlash(commitment.PolicyPath))
			before, err := os.ReadFile(policyFile)
			if err != nil {
				t.Fatal(err)
			}
			_, err = store.Approve(plan.ID, "reviewer sign-off", plan.Effects.Delayed, plan.Effects.Removed)
			staged, readErr := os.ReadFile(policyFile)
			if readErr != nil {
				t.Fatal(readErr)
			}
			if row.after == "" {
				if err != nil || !bytes.Equal(staged, input) {
					t.Fatalf("Approve(the working row edit) = %v with policy %q, want the proposal %q staged", err, staged, input)
				}
				if status := gittest.Output(t, worktree, "status", "--porcelain", "--", path); !strings.HasPrefix(status, "M") && !strings.HasPrefix(status, " M") {
					t.Fatalf("status of %s = %q, want the row edit uncommitted", path, status)
				}
				return
			}
			want := changedFrom(proposed, 0, row.after)
			if err == nil || !strings.Contains(err.Error(), want) || !bytes.Equal(staged, before) {
				t.Fatalf("Approve(a row that changed after the plan) = %v, want %q that leaves the policy file unchanged", err, want)
			}
		})
	}
}

// The working identity drops the one occurrence ledger line, as the tree read does, and a
// malformed ledger refuses.
func TestPinnedRowPlanNormalizesTheOccurrenceLedger(t *testing.T) {
	const kept = "**FT1 — A**\n\nRecord the settled A decision.\n\n"
	for _, row := range []struct {
		name, body, identity, want string
	}{
		{name: "one-ledger", body: kept + "Occurrences: alpha-1\n", identity: kept},
		{name: "two-ledgers", body: kept + "Occurrences: alpha-1\nOccurrences: beta-2\n", identity: kept, want: "has a malformed occurrence ledger"},
	} {
		t.Run(row.name, func(t *testing.T) {
			_, worktree, current := pinnedRows(t)
			commitmenttest.Write(t, worktree, current.Milestones[0].Outcomes[0].Sources[0].Path, row.body)
			_, err := (commitrepo.Store{Root: worktree}).Plan(mustBytes(t, repin(current, 0, row.identity)))
			if row.want == "" && err != nil {
				t.Fatalf("Plan(the identity of the row without its ledger line) = %v, want acceptance", err)
			}
			if row.want != "" && (err == nil || !strings.Contains(err.Error(), row.want)) {
				t.Fatalf("Plan(a row with two ledger lines) = %v, want a refusal naming %q", err, row.want)
			}
		})
	}
}
