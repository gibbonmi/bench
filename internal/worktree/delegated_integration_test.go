// The synthetic delegated run: concurrent ticket authors, one serial integrated
// source chain, chunk acceptance, and the broker-owned landing. It launches no
// model and pays for no comparison; every author is a fixture worktree.
package worktree

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/intent"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/sanitize"
)

const delegatedJourneySpec = "specs/x/spec.md"

// delegatedJourneySpecBody declares the fence the reviewed range must stay inside.
// Each ticket writes its own file, so the fence names one entry per ticket author.
const delegatedJourneySpecBody = "# x\n\nStatus: staged\n\n## User stories\n1. Land a delegated run.\n\n" +
	"### Acceptance coverage map\n| row | story | behavior | seam | why it catches the failure |\n|---|---|---|---|---|\n" +
	"| E1 | 1 | lands | command | catches failure |\n\n## Ownership fences\n\n" +
	"- `source.txt`\n- `ticket-1.txt`\n- `ticket-2.txt`\n- `ticket-3.txt`\n- `ticket-4.txt`\n- `reviews/x.md`\n\n"

// delegatedJourneyChunks is the run the spec describes: two independent chunks,
// one of them with two serial tickets under distinct authors, and one dependent
// chunk. recordtest.Delegate then gives each ticket its own author and maps each
// chunk obligation onto the ticket that owes it.
func delegatedJourneyChunks(plan *rr.Plan) {
	plan.Chunks = []rr.PlannedChunk{
		{ID: "A", Tickets: []string{"1.md", "2.md"}, Verification: []rr.Requirement{
			{ID: "first", Command: "go test ./...", Probe: "omit source check"},
			{ID: "second", Command: "go test -race ./..."},
		}},
		{ID: "B", Tickets: []string{"3.md"}, Verification: []rr.Requirement{
			{ID: "independent", Command: "go test ./..."},
		}},
		{ID: "C", Tickets: []string{"4.md"}, Verification: []rr.Requirement{
			{ID: "dependent", Command: "go test ./..."},
		}},
	}
}

// delegatedJourneyFixture mints the destination, the orchestrator's integration
// assignment, and one assignment per ticket author. The unrelated assignment
// belongs to no chunk of the run, so the landing must leave it alone.
type delegatedJourney struct {
	root, home, base string
	joins            joins
	integration      Creation
	authors          map[string]Creation
	unrelated        Creation
	fixture          *recordtest.Fixture
}

func delegatedJourneyFixture(t *testing.T) *delegatedJourney {
	t.Helper()
	home := filepath.Join(t.TempDir(), "bench-home")
	root := newWorktreeRepo(t)
	mustMkdirAll(t, filepath.Join(root, ".bench"), 0o755)
	// The ordinary gate runs at every chunk checkpoint, before any status
	// transform exists. Only the prospective gate asserts the published
	// transition, so it is the script the landing must run. Both stay POSIX.
	mustWrite(t, filepath.Join(root, ".bench", "gate.sh"), []byte("#!/bin/sh\nset -eu\n[ -f tracked.txt ]\n"), 0o755)
	mustWrite(t, filepath.Join(root, ".bench", "gate-prospective.sh"), []byte("#!/bin/sh\nset -eu\nruntime=$1\ngrep -q '^Status: implemented$' "+delegatedJourneySpec+"\n"), 0o755)
	mustWrite(t, filepath.Join(root, ".bench", "gate-inputs.json"), []byte("{\"schema\":1,\"closure\":\"local\",\"environment\":[],\"paths\":[],\"tools\":[]}\n"), 0o644)
	gitRun(t, root, "add", ".bench")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "declare the journey gate")

	j := delegatedJourneyJoins(t)
	journey := &delegatedJourney{root: root, home: home, joins: j, authors: map[string]Creation{}}
	journey.base = gitOutput(t, root, "rev-parse", "HEAD")
	journey.integration = mustCreate(t, root, home, "delegated-integration", "integration")
	for _, ticket := range []string{"1.md", "2.md", "3.md", "4.md"} {
		label := "author-" + strings.TrimSuffix(ticket, ".md")
		journey.authors[ticket] = mustCreate(t, root, home, "delegated-"+label, label)
	}
	journey.unrelated = mustCreate(t, root, home, "delegated-unrelated", "unrelated")
	commitInWorktree(t, journey.unrelated.Path, "unrelated.txt", "outside the run\n", "unrelated work")
	journey.fixture = recordtest.Prepare(t, journey.integration.Path, 4, delegatedJourneySpec,
		delegatedJourneySpecBody, delegatedJourneyChunks, recordtest.Delegate)
	return journey
}

// delegatedJourneyJoins gives the merge verb a declared lane that costs nothing.
// The lane arrives through the seam, so this journey binds no process environment.
func delegatedJourneyJoins(t *testing.T) joins {
	t.Helper()
	tally := filepath.Join(t.TempDir(), "lane-tally")
	j := defaultJoins()
	j.mergeLane = func(string) (*gate.Lane, error) {
		return &gate.Lane{Checks: []gate.Phase{
			{Name: "unit", Argv: []string{"sh", "-c", "printf g >> " + sanitize.ShellQuote(tally)}},
		}}, nil
	}
	return j
}

// fold carries one committed contribution into the target assignment through the
// existing merge owner. It returns the target's new tip.
func (d *delegatedJourney) fold(t *testing.T, target Creation, from string) string {
	t.Helper()
	code, stdout, stderr := runMerge(t, d.joins, d.root, d.home, "--from", from, target.Assignment.ID)
	if code != 0 {
		t.Fatalf("fold %s into %s = %d; stdout=%q stderr=%q", from, target.Assignment.Label, code, stdout, stderr)
	}
	return gitOutput(t, d.root, "rev-parse", target.Assignment.Branch)
}

// record appends one chunk's evidence at the integration tip and commits it.
func (d *delegatedJourney) record(t *testing.T, base string) {
	t.Helper()
	d.fixture.RecordChunk(base)
	d.fixture.Save()
	d.fixture.Commit("retain delegated chunk evidence")
}

func (d *delegatedJourney) checkpoint(t *testing.T, chunk string) (int, string) {
	t.Helper()
	var out, err bytes.Buffer
	code := gate.RunCommand([]string{d.integration.Path, "--checkpoint", delegatedJourneySpec, "--chunk", chunk}, &out, &err)
	return code, out.String() + err.String()
}

// DI23: concurrent authors reach acceptance through a serial integrated source
// chain, and a branch-only pass cannot close the second integrated chunk.
// DI30: the existing landing owner alone publishes the implemented status, and it
// releases no unrelated assignment.
func TestDelegatedIntegrationJourney(t *testing.T) {
	t.Parallel()
	d := delegatedJourneyFixture(t)
	first, second := d.authors["1.md"], d.authors["2.md"]

	// A same-chunk successor that starts before its predecessor commits carries
	// nothing: the predecessor's branch tip is still the destination tip.
	if tip := d.fold(t, second, first.Assignment.Label); tip != gitOutput(t, d.root, "rev-parse", d.base) {
		t.Fatalf("a pending predecessor moved the successor source to %s", tip)
	}
	if _, err := os.Stat(filepath.Join(second.Path, "ticket-1.txt")); !os.IsNotExist(err) {
		t.Fatalf("the successor holds the predecessor's work before its green commit: %v", err)
	}

	// The predecessor's green commit is what permits the successor's dispatch. No
	// intermediate chunk review stands between the two tickets.
	commitInWorktree(t, first.Path, "ticket-1.txt", "ticket one\n", "ticket 1 work")
	d.fold(t, second, first.Assignment.Label)
	if got, err := os.ReadFile(filepath.Join(second.Path, "ticket-1.txt")); err != nil || string(got) != "ticket one\n" {
		t.Fatalf("the successor source = %q, %v; want the predecessor's committed bytes", got, err)
	}
	commitInWorktree(t, second.Path, "ticket-2.txt", "ticket two\n", "ticket 2 work")
	commitInWorktree(t, d.authors["3.md"].Path, "ticket-3.txt", "ticket three\n", "ticket 3 work")

	// Chunk A enters the integration source as one contribution, and its review
	// reads the tip that holds both tickets.
	chunkABase := gitOutput(t, d.root, "rev-parse", d.integration.Assignment.Branch)
	chunkATip := d.fold(t, d.integration, second.Assignment.Label)
	for _, name := range []string{"ticket-1.txt", "ticket-2.txt"} {
		if _, err := os.Stat(filepath.Join(d.integration.Path, name)); err != nil {
			t.Fatalf("the chunk A tip is missing %s: %v", name, err)
		}
	}
	d.record(t, chunkABase)
	if code, out := d.checkpoint(t, "A"); code != 0 {
		t.Fatalf("the integrated chunk A checkpoint refused: %d %s", code, out)
	}

	// Chunk B folds next, so only one chunk contribution enters the source at a
	// time. Its author verification must name chunk B's own integrated source.
	chunkBTip := d.fold(t, d.integration, d.authors["3.md"].Assignment.Label)
	if chunkBTip == chunkATip {
		t.Fatal("the independent chunk did not reach the integration source")
	}
	d.record(t, chunkATip)
	branchOnly := d.fixture.Record.Chunks[1].Verification[0].SourceDigest
	d.fixture.Record.Chunks[1].Verification[0].SourceDigest = d.fixture.Record.Chunks[0].SourceDigest
	d.fixture.Save()
	d.fixture.Commit("relabel chunk B verification onto the earlier source")
	if code, out := d.checkpoint(t, "B"); code == 0 || !strings.Contains(out, "verification independent") {
		t.Fatalf("a source that predates the chunk B fold closed its checkpoint: %d %s", code, out)
	}
	d.fixture.Record.Chunks[1].Verification[0].SourceDigest = branchOnly
	d.fixture.Save()
	d.fixture.Commit("restore the integrated chunk B verification")
	if code, out := d.checkpoint(t, "B"); code != 0 {
		t.Fatalf("the integrated chunk B checkpoint refused: %d %s", code, out)
	}

	// The dependent chunk starts only after both prerequisite checkpoints pass.
	commitInWorktree(t, d.authors["4.md"].Path, "ticket-4.txt", "ticket four\n", "ticket 4 work")
	chunkCBase := d.fixture.Record.Chunks[1].Tip
	d.fold(t, d.integration, d.authors["4.md"].Assignment.Label)
	d.record(t, chunkCBase)
	if code, out := d.checkpoint(t, "C"); code != 0 {
		t.Fatalf("the integrated chunk C checkpoint refused: %d %s", code, out)
	}

	// The orchestrator's final verification is retained evidence, not publication.
	d.fixture.Complete()
	d.fixture.Save()
	d.fixture.Commit("retain the delegated final verification")
	if gitOutput(t, d.root, "rev-parse", "main") != d.base {
		t.Fatal("final verification alone moved the destination branch")
	}

	tip := gitOutput(t, d.root, "rev-parse", d.integration.Assignment.Branch)
	var stdout, stderr bytes.Buffer
	code := LandCommand(d.root, d.home, "", landArgs("delegated-integration", d.base, tip, d.integration.Path), &stdout, &stderr)
	if code != 0 || !strings.Contains(stdout.String(), "worktree=released") {
		t.Fatalf("delegated landing = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	published := gitOutput(t, d.root, "rev-parse", "main")
	if got := gitOutput(t, d.root, "show", published+":"+delegatedJourneySpec); !strings.Contains(got, "Status: implemented") {
		t.Fatalf("published spec = %q, want the broker-owned status transform", got)
	}

	// The landing releases the source it carried and leaves the unrelated
	// assignment registered on disk.
	assignments, err := intent.Assignments(d.root)
	if err != nil {
		t.Fatal(err)
	}
	var retained []string
	for _, assignment := range assignments {
		retained = append(retained, assignment.Label)
	}
	if !contains(retained, d.unrelated.Assignment.Label) {
		t.Fatalf("retained assignments = %q, want the unrelated assignment kept", retained)
	}
	if _, err := os.Stat(filepath.Join(d.unrelated.Path, "unrelated.txt")); err != nil {
		t.Fatalf("the landing disturbed the unrelated assignment: %v", err)
	}
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
