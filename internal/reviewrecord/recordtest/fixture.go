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
}

func New(t testing.TB, count int) *Fixture { return Attach(t, gittest.RepoOnBranch(t, "main"), count) }

func Attach(t testing.TB, root string, count int) *Fixture {
	t.Helper()
	f := &Fixture{T: t, Root: root}
	plan := rr.Plan{Version: 1, FinalVerification: []rr.Requirement{{ID: "acceptance", Command: "go test ./..."}, {ID: "integration", Command: "go test -tags=system ./..."}}}
	for i := 1; i <= count; i++ {
		id, ticket := fmt.Sprint(i), fmt.Sprintf("%d.md", i)
		plan.Chunks = append(plan.Chunks, rr.PlannedChunk{ID: id, Tickets: []string{ticket}, Verification: []rr.Requirement{{ID: "tests", Command: "go test ./...", Probe: "omit source check"}}})
		blocker := "none"
		if i > 1 {
			blocker = fmt.Sprintf("%d.md", i-1)
		}
		f.Write("specs/example/tickets/"+ticket, fmt.Sprintf("# Chunk %d\n\nBlocked by: %s\nWrites: source.txt\nCovers: E%d\n\n## What to build\n\nImplement the behavior.\n\n## Acceptance\n\n- [ ] E%d: The behavior works.\n", i, blocker, i, i))
	}
	data, err := json.Marshal(plan)
	if err != nil {
		t.Fatal(err)
	}
	f.Write(Spec, "# Example\n\nStatus: staged\n\n```bench-completion-plan\n"+string(data)+"\n```\n")
	f.Write("source.txt", "base\n")
	f.Commit("fixture plan")
	f.Plan, err = rr.ReadPlan(root, f.Tree(), Spec)
	if err != nil {
		t.Fatal(err)
	}
	f.Record = rr.Record{Version: 1, Spec: Spec, PlanDigest: f.Plan.Digest, ImplementationSession: "fixture-author"}
	return f
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
	return rr.Evidence{ID: id, Performer: performer, Role: role, Model: "unknown", Effort: "unknown", SourceDigest: source, State: "completed", Outcome: "pass", NativeRef: Native("completed with no findings")}
}

func (f *Fixture) Verification(prefix, source string, requirements []rr.Requirement) []rr.Verification {
	result := []rr.Verification{}
	for _, requirement := range requirements {
		zero := 0
		item := rr.Verification{Evidence: f.Evidence(prefix+"-"+requirement.ID, source, "author-verification"), Requirement: requirement.ID, Command: requirement.Command, ExitCode: &zero}
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
	digest, err := rr.SourceDigest(f.Root, f.Tree(), Spec)
	if err != nil {
		f.T.Fatal(err)
	}
	chunk := rr.Chunk{ID: planned.ID, Base: base, Tip: f.Tip(), PlanDigest: f.Plan.Digest, SourceDigest: digest, AcceptanceRows: planned.Rows}
	chunk.Verification = f.Verification("chunk-"+planned.ID, digest, planned.Verification)
	for _, axis := range []string{"Standards", "Spec", "Coverage"} {
		chunk.Reviews = append(chunk.Reviews, rr.Review{Evidence: f.Evidence("chunk-"+planned.ID+"-"+axis, digest, "independent-review"), Axis: axis, Base: base, Tip: chunk.Tip})
	}
	f.Record.Chunks = append(f.Record.Chunks, chunk)
}

func (f *Fixture) Complete() {
	source := f.Record.Chunks[len(f.Record.Chunks)-1].SourceDigest
	completion := rr.Completion{State: "completed", SourceDigest: source, Performer: "fixture-author", Reconciliation: map[string]string{}}
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
	f.Write("reviews/example.md", "# Review outcomes\n\n```bench-review-record\n"+string(data)+"\n```\n")
}
