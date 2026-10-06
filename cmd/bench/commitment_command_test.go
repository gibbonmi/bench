package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
)

func TestCommitmentCommandDisplaysActiveOrder(t *testing.T) {
	root := commitmenttest.Repo(t, commitment.Policy{
		Version:         1,
		ActiveMilestone: "M1",
		Milestones: []commitment.Milestone{{
			ID:       "M1",
			Outcomes: []commitment.Outcome{{ID: "A", Criteria: []commitment.Criterion{{ID: "A.done", Text: "The outcome is delivered."}}}, {ID: "B", Criteria: []commitment.Criterion{{ID: "B.done", Text: "The outcome is delivered."}}}},
		}},
	})
	t.Chdir(root)
	var stdout, stderr bytes.Buffer
	code := (Command{Stdout: &stdout, Stderr: &stderr}).Run([]string{"commitment", "show"})
	if code != 0 {
		t.Fatalf("commitment show code = %d, stderr = %q", code, stderr.String())
	}
	out := spilledResponse(t, stdout.String())
	if a, b := strings.Index(out, ",A"), strings.Index(out, ",B"); a < 0 || b < 0 || a >= b {
		t.Fatalf("commitment show output = %q, want A before B", out)
	}
	const outlook = "commitment_outlook[1]{state,active_milestone,next_outcome,deliverable,active,blocked,waiting,command}:\n  eligible,M1,A,\"\",\"\",\"\",B,"
	if !strings.Contains(out, outlook) {
		t.Fatalf("commitment show output = %q, want the outlook that names A next", out)
	}
}

func TestCommitmentCommandStaleApproval(t *testing.T) {
	current := commitment.Policy{
		Version:         1,
		ActiveMilestone: "M1",
		Milestones: []commitment.Milestone{{
			ID:       "M1",
			Outcomes: []commitment.Outcome{{ID: "A", Criteria: []commitment.Criterion{{ID: "A.done", Text: "The outcome is delivered."}}}},
		}},
	}
	root := commitmenttest.Repo(t, current)
	proposal := current
	proposal.Milestones = []commitment.Milestone{{
		ID:       "M1",
		Outcomes: []commitment.Outcome{{ID: "B", Criteria: []commitment.Criterion{{ID: "B.done", Text: "The outcome is delivered."}}}},
	}}
	payload, err := commitment.Bytes(proposal)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := (commitrepo.Store{Root: root}).Plan(payload)
	if err != nil {
		t.Fatal(err)
	}

	changed := current
	changed.Milestones = []commitment.Milestone{{
		ID:       "M1",
		Outcomes: []commitment.Outcome{{ID: "C", Criteria: []commitment.Criterion{{ID: "C.done", Text: "The outcome is delivered."}}}},
	}}
	commitmenttest.WritePolicy(t, root, changed)
	commitmenttest.Commit(t, root, "advance policy")
	t.Chdir(commitmenttest.Planning(t, root))

	var stdout, stderr bytes.Buffer
	code := (Command{Stdout: &stdout, Stderr: &stderr}).Run([]string{
		"commitment", "approve",
		"--plan", plan.ID,
		"--decision", "decision-1",
		"--delayed", "none",
		"--removed", "A",
	})
	if code == 0 {
		t.Fatalf("stale approval code = 0, stdout = %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "predecessor changed") {
		t.Fatalf("stale approval output = %q, want predecessor refusal", stdout.String())
	}
}
