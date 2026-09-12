package reviewrecord_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

const fixtureBody = "# Example\n\nStatus: staged\n\n"

// delegated builds a committed delegated fixture with every chunk recorded.
func delegated(t *testing.T, count int) *recordtest.Fixture {
	t.Helper()
	f := recordtest.NewDelegated(t, count)
	for i := 0; i < count; i++ {
		f.AddChunk()
	}
	f.Complete()
	f.Save()
	f.Commit("retain delegated evidence")
	return f
}

// accept runs the complete source-bound check the checkpoint and landing share.
// It reads the committed record, so a test grades what the tree actually holds.
func accept(f *recordtest.Fixture) error {
	return rr.Check(f.Root, f.Tree(), f.Tip(), f.Record.Spec, "", true)
}

// planFixture commits a delegated plan carrying one deliberate defect. It
// writes without reloading, because the reader's refusal is what the test
// grades and an invalid plan cannot be read back.
func planFixture(t *testing.T, count int, mutate func(*rr.Plan)) *recordtest.Fixture {
	t.Helper()
	f := recordtest.NewDelegated(t, count)
	mutate(&f.Plan)
	f.WritePlan()
	f.Commit("stage the plan under test")
	return f
}

// mergeTickets puts every planned ticket in one chunk, so one chunk owes
// verification to two distinct authors.
func mergeTickets(p *rr.Plan) {
	p.Chunks[0].Tickets = append(p.Chunks[0].Tickets, p.Chunks[1].Tickets...)
	p.Chunks = p.Chunks[:1]
}

// twoTicketChunk is a recorded single chunk owned by two ticket authors.
func twoTicketChunk(t *testing.T) *recordtest.Fixture {
	t.Helper()
	f := recordtest.Prepare(t, gittest.RepoOnBranch(t, "main"), 2, recordtest.Spec, fixtureBody, mergeTickets, recordtest.Delegate)
	f.AddChunk()
	f.Complete()
	f.Save()
	f.Commit("retain two-author chunk evidence")
	return f
}

// DI2: valid delegated evidence is accepted, and verification ownership comes
// from the source-bound plan rather than from the record.
func TestDelegatedIdentityAuthority(t *testing.T) {
	f := delegated(t, 1)
	if err := accept(f); err != nil {
		t.Fatalf("valid delegated evidence refused: %v", err)
	}
	f.Record.Chunks[0].Verification[0].Performer = recordtest.Orchestrator
	f.Save()
	f.Commit("rewrite the recorded performer")
	if err := accept(f); err == nil || !strings.Contains(err.Error(), "verification tests") {
		t.Fatalf("an evidence-only performer rewrite changed verification ownership: %v", err)
	}
}

// DI3: malformed, duplicate, conflicting, and mixed-version declarations refuse.
func TestDelegatedIdentityRefusals(t *testing.T) {
	cases := []struct {
		name, want string
		mutate     func(*rr.Plan)
	}{
		{"missing execution", "missing execution declaration", func(p *rr.Plan) { p.Execution = nil }},
		{"wrong mode", "invalid execution mode", func(p *rr.Plan) { p.Execution.Mode = "retain" }},
		{"missing run id", "missing run id", func(p *rr.Plan) { p.Execution.RunID = "" }},
		{"missing orchestrator", "missing run id", func(p *rr.Plan) { p.Execution.OrchestratorSession = "" }},
		{"zero author limit", "invalid author limit", func(p *rr.Plan) { p.Execution.AuthorLimit = 0 }},
		{"foreign assignment ticket", "invalid assignment ticket", func(p *rr.Plan) {
			p.Execution.Assignments["9.md"] = []rr.Assignment{recordtest.Assign("9.md")}
		}},
		{"absent assignment history", "missing assignment history", func(p *rr.Plan) {
			delete(p.Execution.Assignments, "1.md")
		}},
		{"incomplete assignment", "missing session", func(p *rr.Plan) {
			p.Execution.Assignments["1.md"][0].NativeRef = ""
		}},
		{"author is the orchestrator", "identities must differ", func(p *rr.Plan) {
			p.Execution.Assignments["1.md"][0].Session = recordtest.Orchestrator
		}},
		{"version 1 plan with execution", "invalid execution declaration", func(p *rr.Plan) { p.Version = 1 }},
		{"replacement without a trigger", "invalid replacement trigger", func(p *rr.Plan) {
			p.Execution.Assignments["1.md"] = replaced(p.Execution.Assignments["1.md"][0], "")
		}},
		{"replacement without stop proof", "stopped-writer evidence", func(p *rr.Plan) {
			history := replaced(p.Execution.Assignments["1.md"][0], "session-lost")
			history[1].Stopped = ""
			p.Execution.Assignments["1.md"] = history
		}},
		{"no-progress without reassessment", "requires a reassessment", func(p *rr.Plan) {
			history := replaced(p.Execution.Assignments["1.md"][0], "no-progress")
			history[1].Reassessment = ""
			p.Execution.Assignments["1.md"] = history
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := planFixture(t, 1, tc.mutate)
			if _, err := rr.ReadPlan(f.Root, f.Tree(), recordtest.Spec); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("an invalid delegated plan was accepted or lost its reason: %v", err)
			}
		})
	}
}

// replaced appends a successor assignment to one first dispatch.
func replaced(first rr.Assignment, trigger string) []rr.Assignment {
	successor := recordtest.Assign("successor")
	successor.Predecessor = first.Session
	successor.Trigger = trigger
	successor.Stopped = "fixture:writer-stopped"
	successor.Preserved = "fixture:preserved-source"
	if trigger == "no-progress" {
		successor.Reassessment = "the seam is wider than the ticket assumed"
	}
	return []rr.Assignment{first, successor}
}

// DI41: one native author session cannot own two ticket histories.
func TestDelegatedDistinctTicketAuthors(t *testing.T) {
	f := planFixture(t, 2, func(p *rr.Plan) {
		p.Execution.Assignments["2.md"][0].Session = recordtest.Author("1.md")
	})
	if _, err := rr.ReadPlan(f.Root, f.Tree(), recordtest.Spec); err == nil || !strings.Contains(err.Error(), "one session owns one ticket") {
		t.Fatalf("one session owned two ticket histories: %v", err)
	}
}

// DI31: an undispatched future ticket accepts an empty assignment history.
func TestDelegatedPendingAssignments(t *testing.T) {
	f := planFixture(t, 2, func(p *rr.Plan) { p.Execution.Assignments["2.md"] = nil })
	if _, err := rr.ReadPlan(f.Root, f.Tree(), recordtest.Spec); err != nil {
		t.Fatalf("a bounded-author plan with one undispatched ticket refused: %v", err)
	}
}

// DI37: one chunk accepts verification from each of its distinct ticket authors.
func TestDelegatedTicketOwners(t *testing.T) {
	f := twoTicketChunk(t)
	if err := accept(f); err != nil {
		t.Fatalf("a two-ticket chunk with two distinct authors refused: %v", err)
	}
	if a, b := f.Record.Chunks[0].Verification[0].Performer, f.Record.Chunks[0].Verification[1].Performer; a == b {
		t.Fatalf("the two-author fixture used one performer %q for both obligations", a)
	}
}

// DI38: a foreign, uncovered, or wrongly performed obligation refuses the chunk.
func TestDelegatedTicketObligations(t *testing.T) {
	cases := []struct {
		name, want string
		mutate     func(*rr.Plan)
	}{
		{"foreign ticket", "outside this chunk", func(p *rr.Plan) { p.Chunks[0].Verification[0].Ticket = "9.md" }},
		{"uncovered ticket", "no verification requirement", func(p *rr.Plan) {
			mergeTickets(p)
			for i := range p.Chunks[0].Verification {
				p.Chunks[0].Verification[i].Ticket = "1.md"
			}
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := planFixture(t, 2, tc.mutate)
			if _, err := rr.ReadPlan(f.Root, f.Tree(), recordtest.Spec); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("an invalid ticket obligation map was accepted: %v", err)
			}
		})
	}
	t.Run("another ticket's performer", func(t *testing.T) {
		f := twoTicketChunk(t)
		f.Record.Chunks[0].Verification[0].Performer = f.Record.Chunks[0].Verification[1].Performer
		f.Save()
		f.Commit("give one author both obligations")
		if err := accept(f); err == nil || !strings.Contains(err.Error(), "verification tests") {
			t.Fatalf("another ticket's author satisfied an obligation: %v", err)
		}
	})
}

// DI1: the version 1 fixture matrix keeps its pre-change outcomes.
func TestDelegatedLegacyParity(t *testing.T) {
	f := recordtest.New(t, 1)
	f.AddChunk()
	f.Complete()
	f.Save()
	f.Commit("retain version 1 evidence")
	if f.Record.Version != 1 || f.Record.ImplementationSession == "" {
		t.Fatal("the version 1 fixture lost its own identity")
	}
	if err := accept(f); err != nil {
		t.Fatalf("valid version 1 evidence refused: %v", err)
	}
	f.Record.Chunks[0].Verification[0].Performer = "other"
	f.Save()
	f.Commit("rewrite the version 1 performer")
	if err := accept(f); err == nil || !strings.Contains(err.Error(), "verification tests") {
		t.Fatalf("version 1 verification ownership drifted: %v", err)
	}
}

// A version 2 record carries no identity of its own, and a version 1 record
// cannot omit one. An unknown version refuses before evidence validation.
func TestDelegatedRecordVersions(t *testing.T) {
	cases := []struct {
		name, want, session string
		version             int
	}{
		{"version 2 names an identity", "invalid implementation session", "smuggled", 2},
		{"version 1 omits its identity", "invalid implementation session", "", 1},
		{"unknown version", "unsupported version", "any", 3},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			record := rr.Record{Version: tc.version, Spec: recordtest.Spec, PlanDigest: "sha256:x", ImplementationSession: tc.session}
			data, _ := json.Marshal(record)
			if _, err := rr.Parse(data); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("an invalid record version was accepted: %v", err)
			}
		})
	}
}
