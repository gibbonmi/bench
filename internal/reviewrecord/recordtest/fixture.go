// Package recordtest supplies committed completion-evidence fixtures to consumer tests.
package recordtest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/testrepo"
)

const Spec = "specs/example/spec.md"

type Fixture struct {
	T      testing.TB
	Root   string
	Record rr.Record
	Plan   rr.Plan
	// spec and body retain what Prepare wrote, so RewritePlan can replace the
	// plan fence without the consumer repeating the surrounding document.
	spec string
	body string
}

func New(t testing.TB, count int) *Fixture { return Attach(t, gittest.RepoOnBranch(t, "main"), count) }

func Attach(t testing.TB, root string, count int) *Fixture {
	return AttachAt(t, root, count, Spec)
}

func AttachAt(t testing.TB, root string, count int, spec string) *Fixture {
	return Prepare(t, root, count, spec, "# Example\n\nStatus: staged\n\n")
}

// AttachDelegated is Attach in the version 2 delegated form.
func AttachDelegated(t testing.TB, root string, count int) *Fixture {
	return Prepare(t, root, count, Spec, "# Example\n\nStatus: staged\n\n", Delegate)
}

// NewDelegated is New in the version 2 delegated form.
func NewDelegated(t testing.TB, count int) *Fixture {
	return AttachDelegated(t, gittest.RepoOnBranch(t, "main"), count)
}

// WritePlan writes the fixture's current plan back into its spec without
// committing it. An amendment that must land inside the next chunk's own delta
// stages the plan here, so the chunk's base keeps its predecessor's source.
func (f *Fixture) WritePlan() {
	f.T.Helper()
	data, err := json.Marshal(f.Plan)
	if err != nil {
		f.T.Fatal(err)
	}
	f.Write(f.spec, f.body+"```bench-completion-plan\n"+string(data)+"\n```\n")
}

// RewritePlan writes the fixture's current plan back into its spec, commits it,
// and reloads it. A consumer mutates one field of f.Plan, calls this, and
// grades the plan the reader then sees.
func (f *Fixture) RewritePlan() {
	f.T.Helper()
	f.WritePlan()
	f.Commit("rewrite fixture plan")
	f.Reload()
}

// Reload reads the committed plan again, keeping the retained record. A plan
// amendment changes the digest that later evidence must name.
func (f *Fixture) Reload() {
	f.T.Helper()
	plan, err := rr.ReadPlan(f.Root, f.Tree(), f.Record.Spec)
	if err != nil {
		f.T.Fatal(err)
	}
	f.Plan = plan
}

// Option reshapes the planned fixture before it is written and committed.
type Option func(*rr.Plan)

// Orchestrator is the delegated fixture's orchestrator session.
const Orchestrator = "fixture-orchestrator"

// Author is the delegated fixture's author session for one ticket basename.
func Author(ticket string) string { return "fixture-author-" + ticket }

// Delegate turns a planned fixture into the version 2 delegated form. It gives
// every ticket its own author and maps each chunk obligation onto the ticket
// that owes it, so a consumer starts from a valid delegated plan and mutates
// one field to express the case under test.
func Delegate(plan *rr.Plan) {
	plan.Version = 2
	execution := &rr.Execution{Mode: "delegate", RunID: "fixture-run", OrchestratorSession: Orchestrator, AuthorLimit: 2, Assignments: map[string][]rr.Assignment{}}
	for i := range plan.Chunks {
		chunk := &plan.Chunks[i]
		for j := range chunk.Verification {
			chunk.Verification[j].Ticket = chunk.Tickets[j%len(chunk.Tickets)]
		}
		for _, ticket := range chunk.Tickets {
			execution.Assignments[ticket] = []rr.Assignment{Assign(ticket)}
		}
	}
	plan.Execution = execution
}

// Assign is one first-dispatch assignment for a ticket.
func Assign(ticket string) rr.Assignment {
	return rr.Assignment{Session: Author(ticket), Assignment: "fixture-assignment-" + ticket, Model: "unknown", Effort: "unknown", Source: "fixture-source", NativeRef: "fixture:dispatch-" + ticket}
}

// Prepare supplies a plan and tickets beside the consumer fixture's own spec body.
func Prepare(t testing.TB, root string, count int, spec, body string, options ...Option) *Fixture {
	t.Helper()
	f := &Fixture{T: t, Root: root, spec: spec, body: body}
	plan := rr.Plan{Version: 1, FinalVerification: []rr.Requirement{{ID: "acceptance", Command: "go test ./..."}, {ID: "integration", Command: "go test -tags=system ./..."}}}
	for i := 1; i <= count; i++ {
		id, ticket := fmt.Sprint(i), fmt.Sprintf("%d.md", i)
		plan.Chunks = append(plan.Chunks, rr.PlannedChunk{ID: id, Tickets: []string{ticket}, Verification: []rr.Requirement{{ID: "tests", Command: "go test ./...", Probe: "omit source check"}, {ID: "additional", Command: "go test -race ./..."}}})
		blocker := "none"
		if i > 1 {
			blocker = fmt.Sprintf("%d.md", i-1)
		}
		f.Write(filepath.ToSlash(filepath.Join(filepath.Dir(spec), "tickets", ticket)), fmt.Sprintf("# Chunk %d\n\nBlocked by: %s\nWrites: source.txt\nCovers: E%d\n\n## What to build\n\nImplement the behavior.\n\n## Acceptance\n\n- [ ] E%d: The behavior works.\n", i, blocker, i, i))
	}
	for _, option := range options {
		option(&plan)
	}
	data, err := json.Marshal(plan)
	if err != nil {
		t.Fatal(err)
	}
	f.Write(spec, body+"```bench-completion-plan\n"+string(data)+"\n```\n")
	f.Write("source.txt", "base\n")
	f.Commit("fixture plan")
	f.loadPlan(spec)
	return f
}

func (f *Fixture) loadPlan(spec string) {
	f.T.Helper()
	var err error
	f.Plan, err = rr.ReadPlan(f.Root, f.Tree(), spec)
	if err != nil {
		f.T.Fatal(err)
	}
	f.Record = rr.Record{Version: 1, Spec: spec, PlanDigest: f.Plan.Digest, ImplementationSession: "fixture-author"}
	if f.Plan.Delegated() {
		// A version 2 record names no identity of its own; the plan owns it.
		f.Record.Version, f.Record.ImplementationSession = 2, ""
	}
}

// RetainSingleChunk adds complete fixture evidence for an already committed source.
func RetainSingleChunk(t testing.TB, root, spec, base string) {
	t.Helper()
	f := &Fixture{T: t, Root: root}
	f.loadPlan(spec)
	if len(f.Plan.Chunks) != 1 {
		t.Fatal("single-chunk fixture requires one planned chunk")
	}
	f.RecordChunk(base)
	f.Complete()
	f.Save()
	path, err := rr.RecordPath(spec)
	if err != nil {
		t.Fatal(err)
	}
	f.Git("add", "--", ":(literal)"+path)
	f.Git("-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "retain fixture completion")
}

func (f *Fixture) Write(path, data string) {
	f.T.Helper()
	full := filepath.Join(f.Root, path)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		f.T.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(data), 0o644); err != nil {
		f.T.Fatal(err)
	}
}

func (f *Fixture) Git(args ...string) string {
	f.T.Helper()
	out, err := benchgit.Output(append([]string{"-C", f.Root}, args...)...)
	if err != nil {
		f.T.Fatalf("git %v: %v", args, err)
	}
	return out
}
func (f *Fixture) Tip() string  { return f.Git("rev-parse", "HEAD") }
func (f *Fixture) Tree() string { return f.Git("rev-parse", "HEAD^{tree}") }
func (f *Fixture) Commit(message string) {
	f.T.Helper()
	if err := testrepo.CommitAll(f.Root, message); err != nil {
		f.T.Fatal(err)
	}
}

func Native(text string) rr.NativeRef {
	return rr.NativeRef{Ref: "fixture:terminal-result", Excerpt: text, Digest: rr.Digest([]byte(text))}
}

func (f *Fixture) Evidence(id, source, role string) rr.Evidence {
	performer := "fixture-author"
	if role == "independent-review" {
		performer = "fixture-reviewer-" + id
	}
	return f.EvidenceAs(id, source, role, performer)
}

// EvidenceAs is one occurrence attributed to an explicit performer.
func (f *Fixture) EvidenceAs(id, source, role, performer string) rr.Evidence {
	return rr.Evidence{ID: id, Performer: performer, Role: role, Model: "unknown", Effort: "unknown", SourceDigest: source, State: "completed", Outcome: "pass", NativeRef: Native("completed with no findings")}
}

// owes is the performer and role one requirement expects. A version 1 fixture
// keeps its single author. A delegated fixture reads the chunk obligation's
// owning ticket, and gives the final obligations to the orchestrator.
func (f *Fixture) owes(requirement rr.Requirement) (string, string) {
	if !f.Plan.Delegated() {
		return "fixture-author", "author-verification"
	}
	if requirement.Ticket == "" {
		return f.Plan.Execution.OrchestratorSession, "integration-verification"
	}
	session, _ := f.Plan.Author(requirement.Ticket)
	return session, "author-verification"
}

func (f *Fixture) Verification(prefix, source string, requirements []rr.Requirement) []rr.Verification {
	result := []rr.Verification{}
	for _, requirement := range requirements {
		zero := 0
		performer, role := f.owes(requirement)
		item := rr.Verification{Evidence: f.EvidenceAs(prefix+"-"+requirement.ID, source, role, performer), Requirement: requirement.ID, Command: requirement.Command, ExitCode: &zero}
		if requirement.Probe != "" {
			item.Probe = &rr.Probe{Mutation: requirement.Probe, Outcome: "bit", ExitCode: 1, Restore: "pass", NativeRef: Native("mutation failed at the required assertion; restore passed")}
		}
		result = append(result, item)
	}
	return result
}

func (f *Fixture) AddChunk() {
	f.T.Helper()
	planned := f.Plan.Chunks[len(f.Record.Chunks)]
	base := f.Tip()
	f.Write("source.txt", "implemented chunk "+planned.ID+"\n")
	f.Commit("chunk " + planned.ID)
	f.RecordChunk(base)
}

func (f *Fixture) RecordChunk(base string) {
	f.T.Helper()
	planned := f.Plan.Chunks[len(f.Record.Chunks)]
	digest, err := rr.SourceDigest(f.Root, f.Tree(), f.Record.Spec)
	if err != nil {
		f.T.Fatal(err)
	}
	chunk := rr.Chunk{ID: planned.ID, Base: base, Tip: f.Tip(), PlanDigest: f.Plan.Digest, SourceDigest: digest, AcceptanceRows: planned.Rows}
	chunk.Verification = f.Verification("chunk-"+planned.ID, digest, planned.Verification)
	for _, axis := range rr.Axes() {
		chunk.Reviews = append(chunk.Reviews, rr.Review{Evidence: f.Evidence("chunk-"+planned.ID+"-"+axis, digest, "independent-review"), Axis: axis, Base: base, Tip: chunk.Tip})
	}
	f.Record.Chunks = append(f.Record.Chunks, chunk)
}

func (f *Fixture) Complete() {
	source := f.Record.Chunks[len(f.Record.Chunks)-1].SourceDigest
	reconciler, _ := f.owes(rr.Requirement{})
	completion := rr.Completion{State: "completed", SourceDigest: source, Performer: reconciler, Reconciliation: map[string]string{}}
	for _, chunk := range f.Plan.Chunks {
		for _, row := range chunk.Rows {
			completion.Reconciliation[row] = "covered"
		}
	}
	completion.Verification = f.Verification("final", source, f.Plan.FinalVerification)
	f.Record.Completion = completion
}

func (f *Fixture) Save() {
	f.T.Helper()
	data, err := json.MarshalIndent(f.Record, "", "  ")
	if err != nil {
		f.T.Fatal(err)
	}
	path, err := rr.RecordPath(f.Record.Spec)
	if err != nil {
		f.T.Fatal(err)
	}
	f.Write(path, "# Review outcomes\n\n```bench-review-record\n"+string(data)+"\n```\n")
}
