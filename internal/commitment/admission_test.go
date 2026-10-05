package commitment_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitcmd"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/intent"
)

func admissionRepo(t *testing.T, dependent, parallel bool) (string, string, string) {
	t.Helper()
	p := policy([]commitment.Milestone{milestone("M1", "A", "B", "C")}, "M1")
	for i := range p.Milestones[0].Outcomes {
		outcome := &p.Milestones[0].Outcomes[i]
		outcome.Deliverables = []commitment.DeliveryBinding{{Source: commitment.SourceBinding{ID: outcome.ID + ".spec", Path: deliverable(outcome.ID), Identity: commitment.Identity([]byte(deliverableBody))}}}
	}
	if dependent {
		p.Milestones[0].Outcomes[1].Dependencies = []string{"A"}
	}
	if parallel {
		p.ParallelGrants = []commitment.ParallelGrant{{Outcomes: []string{"A", "B"}}}
	}
	root := commitmenttest.Repo(t, p)
	for _, id := range []string{"A", "B", "C"} {
		commitmenttest.Write(t, root, deliverable(id), deliverableBody)
	}
	commitmenttest.Commit(t, root, "approve deliverable sources")
	return root, commitmenttest.Assignment(t, root, "first"), commitmenttest.Assignment(t, root, "second")
}

const deliverableBody = "# Delivery\n\nStatus: staged\n"

func deliverable(outcome string) string { return "specs/" + outcome + "/spec.md" }

func startOutcome(root, outcome, request string) (string, int) {
	return commitcmd.Command(root, []string{"start", "--outcome", outcome, "--request", request, "--deliverable", deliverable(outcome)})
}

func mustCommand(t *testing.T, root string, args ...string) string {
	t.Helper()
	out, code := commitcmd.Command(root, args)
	if code != 0 {
		t.Fatalf("%v = (%s,%d)", args, out, code)
	}
	return out
}

func mustStart(t *testing.T, root, outcome, request string) {
	t.Helper()
	out, code := startOutcome(root, outcome, request)
	if code != 0 {
		t.Fatalf("start %s = (%s,%d)", outcome, out, code)
	}
}

func TestCommitmentOrderedStart(t *testing.T) {
	root, first, _ := admissionRepo(t, false, false)
	path, err := intent.Address(root)
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	out, code := startOutcome(first, "B", "first")
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if code != 1 || !strings.Contains(out, "A") || string(before) != string(after) {
		t.Fatalf("premature B=(%s,%d), unchanged=%v", out, code, string(before) == string(after))
	}
	mustStart(t, first, "A", "first")
}

func TestCommitmentUncommittedStart(t *testing.T) {
	root, first, _ := admissionRepo(t, false, false)
	out, code := startOutcome(first, "uncommitted", "first")
	state, err := intent.Read(root)
	if err != nil {
		t.Fatal(err)
	}
	if code != 1 || !strings.Contains(out, "bench commitment plan") || state.Commitment != nil {
		t.Fatalf("uncommitted start=(%s,%d), state=%+v", out, code, state.Commitment)
	}
}

func TestCommitmentBlockedSuccessor(t *testing.T) {
	root, first, second := admissionRepo(t, false, false)
	mustStart(t, first, "A", "first")
	mustCommand(t, first, "block", "--outcome", "A", "--reason", "Waiting for an input")
	mustStart(t, second, "B", "second")
	state, err := intent.Read(root)
	if err != nil {
		t.Fatal(err)
	}
	if state.Commitment == nil || len(state.Commitment.Claims) != 1 || state.Commitment.Claims[0].Outcome != "B" || len(state.Commitment.Blockers) != 1 || state.Commitment.Blockers[0].Reason != "Waiting for an input" {
		t.Fatalf("blocked successor state=%+v", state.Commitment)
	}
}

func TestCommitmentBlockedDependency(t *testing.T) {
	_, first, second := admissionRepo(t, true, false)
	mustCommand(t, first, "block", "--outcome", "A", "--reason", "Waiting")
	out, code := startOutcome(second, "B", "second")
	if code != 1 || !strings.Contains(out, "dependency") {
		t.Fatalf("dependent start=(%s,%d)", out, code)
	}
}

func TestCommitmentAllBlocked(t *testing.T) {
	_, first, _ := admissionRepo(t, false, false)
	for _, id := range []string{"A", "B", "C"} {
		mustCommand(t, first, "block", "--outcome", id, "--reason", "Waiting")
	}
	out, code := startOutcome(first, "uncommitted", "first")
	if code != 1 || !strings.Contains(out, "bench commitment plan") {
		t.Fatalf("all blocked start=(%s,%d)", out, code)
	}
}

func TestCommitmentSingleOutcome(t *testing.T) {
	_, first, second := admissionRepo(t, false, false)
	mustStart(t, first, "A", "first")
	out, code := startOutcome(second, "B", "second")
	if code != 1 || !strings.Contains(out, "active") {
		t.Fatalf("second outcome=(%s,%d)", out, code)
	}
}

func TestCommitmentExactParallelGrant(t *testing.T) {
	root, first, second := admissionRepo(t, false, true)
	mustStart(t, first, "A", "first")
	mustStart(t, second, "B", "second")
	third := commitmenttest.Assignment(t, root, "third")
	out, code := startOutcome(third, "C", "third")
	if code != 1 {
		t.Fatalf("ungranted C=(%s,%d)", out, code)
	}
}

func TestCommitmentUnblockOrder(t *testing.T) {
	root, first, second := admissionRepo(t, false, false)
	mustStart(t, first, "A", "first")
	mustCommand(t, first, "block", "--outcome", "A", "--reason", "Waiting")
	mustStart(t, second, "B", "second")
	mustCommand(t, first, "unblock", "--outcome", "A")
	out, code := startOutcome(first, "A", "first")
	state, err := intent.Read(root)
	if err != nil {
		t.Fatal(err)
	}
	if code != 1 || !strings.Contains(out, "active") || len(state.Commitment.Claims) != 1 || state.Commitment.Claims[0].Outcome != "B" {
		t.Fatalf("unblock displaced B: (%s,%d),%+v", out, code, state.Commitment)
	}
}

func TestCommitmentSameOutcomeTickets(t *testing.T) {
	root, first, second := admissionRepo(t, false, false)
	mustStart(t, first, "A", "first")
	mustStart(t, second, "A", "second")
	state, err := intent.Read(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(state.Commitment.Bindings) != 2 || len(state.Commitment.Claims) != 1 {
		t.Fatalf("ticket claims=%+v", state.Commitment)
	}
}

func TestCommitmentAdoptionRequired(t *testing.T) {
	root := commitmenttest.Repo(t, policy([]commitment.Milestone{milestone("M1", "A")}, "M1"))
	if err := os.Remove(filepath.Join(root, ".bench", "commitment.json")); err != nil {
		t.Fatal(err)
	}
	commitmenttest.Commit(t, root, "without adoption")
	first := commitmenttest.Assignment(t, root, "first")
	out, code := startOutcome(first, "A", "first")
	if code != 1 || !strings.Contains(out, "adoption required") || !strings.Contains(out, "bench commitment plan --input") {
		t.Fatalf("adoption=(%s,%d)", out, code)
	}
}

func TestCommitmentRetainedClaim(t *testing.T) {
	if root := os.Getenv("DC_START_PROCESS_ROOT"); root != "" {
		out, code := startOutcome(root, "A", "first")
		if code != 0 {
			t.Fatalf("child start=(%s,%d)", out, code)
		}
		return
	}
	root, first, second := admissionRepo(t, false, false)
	child := exec.Command(os.Args[0], "-test.run=^TestCommitmentRetainedClaim$")
	child.Env = append(os.Environ(), "DC_START_PROCESS_ROOT="+first)
	if out, err := child.CombinedOutput(); err != nil {
		t.Fatalf("start process: %v: %s", err, out)
	}
	if _, err := intent.PurgeAssignments(root, func(assignment intent.Assignment) bool { return assignment.Request != intent.RequestDigest("first") }); err != nil {
		t.Fatal(err)
	}
	out, code := startOutcome(second, "B", "second")
	state, err := intent.Read(root)
	if err != nil {
		t.Fatal(err)
	}
	if code != 1 || !strings.Contains(out, "active") || len(state.Commitment.Claims) != 1 || state.Commitment.Claims[0].Outcome != "A" || len(state.Commitment.Bindings) != 1 {
		t.Fatalf("lost retained claim: (%s,%d),%+v", out, code, state.Commitment)
	}
}

func TestCommitmentStartIdentityAndAtomicRefusal(t *testing.T) {
	root, first, _ := admissionRepo(t, false, false)
	path, err := intent.Address(root)
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{
		{"start", "--outcome", "A", "--request", "wrong", "--deliverable", deliverable("A")},
		{"start", "--outcome", "A", "--request", "first", "--deliverable", deliverable("B")},
	} {
		out, code := commitcmd.Command(first, args)
		after, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if code != 1 || string(before) != string(after) {
			t.Fatalf("non-atomic refusal: (%s,%d), unchanged=%v", out, code, string(before) == string(after))
		}
	}
	mustStart(t, first, "A", "first")
	before, err = os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	mustStart(t, first, "A", "first")
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("identical start rewrote the ledger")
	}
}

func TestCommitmentMalformedAdmission(t *testing.T) {
	root, first, _ := admissionRepo(t, false, false)
	commitmenttest.Write(t, root, ".bench/commitment.json", "{")
	commitmenttest.Commit(t, root, "malformed policy")
	out, code := startOutcome(first, "A", "first")
	if code != 1 || !strings.Contains(out, "commitment policy:") || strings.Contains(out, "adoption required") {
		t.Fatalf("corruption=(%s,%d)", out, code)
	}
}

func TestCommitmentStartPublishedIdentity(t *testing.T) {
	for _, kind := range []string{"spec", "tickets-only"} {
		for _, change := range []string{"changed", "deleted"} {
			t.Run(kind+"/"+change, func(t *testing.T) {
				source := "specs/light/spec.md"
				member := source
				files := map[string]string{source: deliverableBody}
				folder := kind == "tickets-only"
				if folder {
					source = "specs/light"
					member = source + "/tickets/01-deliver.md"
					files = map[string]string{member: "# First ticket\n", source + "/tickets/02-deliver.md": "# Second ticket\n"}
				}
				root := boundDeliverable(t, source, files, folder)
				first := commitmenttest.Assignment(t, root, "first")
				path, err := intent.Address(root)
				if err != nil {
					t.Fatal(err)
				}
				before, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				if change == "deleted" {
					if err := os.Remove(filepath.Join(root, member)); err != nil {
						t.Fatal(err)
					}
				} else {
					commitmenttest.Write(t, root, member, files[member]+"Changed acceptance.\n")
				}
				commitmenttest.Commit(t, root, "change approved source")
				out, code := commitcmd.Command(first, []string{"start", "--outcome", "A", "--request", "first", "--deliverable", source})
				after, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				if code != 1 || string(before) != string(after) {
					t.Fatalf("stale deliverable=(%s,%d), unchanged=%v", out, code, string(before) == string(after))
				}
			})
		}
	}
}
