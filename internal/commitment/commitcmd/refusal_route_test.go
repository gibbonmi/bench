// Refusal route tests for the commitment verb: each refusal prints, in its next table, the
// route of its face in the shared refusal-route registry, and that route clears the refusal.
package commitcmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitcmd"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/refusalroute/routetest"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/worktree"
)

// commitmentSet is one commitment route fixture: the primary checkout, the Bench home that
// holds its assignments, the checkout the verb runs in, the verb's arguments, and the input
// file that the verb reads.
type commitmentSet struct {
	root, home, dir, input string
	args                   []string
}

// commitmentFaceFixture produces exactly one commitment face for one cause. The registry walk
// requires a fixture for each commitment face that the commitment walk proves, so a face
// added with no fixture turns TestCommitmentFacesFollowTheirRoutes red.
type commitmentFaceFixture struct {
	face, cause string
	build       func(t *testing.T) commitmentSet
	// prefix and contains state what the printed route must open with and hold, and names
	// the values of the set that it must name.
	prefix, contains string
	names            func(f commitmentSet) []string
	// carry carries out a printed step by the fixture's own means, keyed by the step's index
	// in the face's route: an instruction, or a step of a reviewer route.
	carry map[int]func(t *testing.T, f commitmentSet)
	// fill does what the operator does to fill the operator slots of a printed command step,
	// and it returns the replacement pairs of those slots.
	fill func(t *testing.T, f commitmentSet) []string
	// after carries out what the route's last step leaves to the operator, and it returns the
	// set that the verb then reruns in.
	after func(t *testing.T, f commitmentSet, last string) commitmentSet
}

// TestCommitmentFacesFollowTheirRoutes is RR45, RR46, and RR48. The registry is the source of
// the commitment face set, so each commitment face that no commit candidate alone raises
// needs a producing fixture. Each fixture's printed route is carried out step by step: an
// instruction or a reviewer step by the fixture's own means, and each other command step
// verbatim through the verb's own entry. The commitment verb then reruns out of the face.
func TestCommitmentFacesFollowTheirRoutes(t *testing.T) {
	var keys [][2]string
	for _, fixture := range commitmentFaceFixtures() {
		keys = append(keys, [2]string{fixture.face, fixture.cause})
	}
	faces := routetest.Fixtures(t, refusalroute.Commitment, keys)
	for _, fixture := range commitmentFaceFixtures() {
		t.Run(fixture.face+"/"+fixture.cause, func(t *testing.T) { followCommitmentFace(t, faces[fixture.face], fixture) })
	}
}

// followCommitmentFace drives one fixture: the verb prints the face, the walk carries out the
// printed route, and the verb reruns out of the face.
func followCommitmentFace(t *testing.T, face refusalroute.Face, fixture commitmentFaceFixture) {
	t.Setenv("BENCH_HOME", t.TempDir())
	f := fixture.build(t)
	f.home = os.Getenv("BENCH_HOME")
	out, code := commitcmd.Command(f.dir, f.args)
	next, printed := routetest.NextCell(out)
	if code != 1 || !printed || !strings.HasPrefix(next, fixture.prefix) || !strings.Contains(next, fixture.contains) {
		t.Fatalf("face %s = (%d, %q), want exit 1 and a next cell that opens with %q and holds %q", face.Name, code, out, fixture.prefix, fixture.contains)
	}
	var names []string
	if fixture.names != nil {
		names = fixture.names(f)
	}
	steps := routetest.Steps(t, face, next, "", names)
	carry := func(index int) (func(), bool) {
		step, carried := fixture.carry[index]
		return func() { step(t, f) }, carried
	}
	last := routetest.Follow(t, face, steps, carry, func(step string) string {
		if fixture.fill != nil {
			step = strings.NewReplacer(fixture.fill(t, f)...).Replace(step)
		}
		return runRouteStep(t, f, step)
	})
	if fixture.after != nil {
		f = fixture.after(t, f, last)
	}
	if out, code := commitcmd.Command(f.dir, f.args); code != 0 {
		t.Fatalf("%s rerun = (%d, %q), want exit 0", face.Name, code, out)
	}
}

// runRouteStep runs one printed command step in process, the way the CLI dispatches it, at
// the fixture's checkout. The shared walk states the rules a printed step and its exit meet.
// It returns the step's stdout.
func runRouteStep(t *testing.T, f commitmentSet, step string) string {
	t.Helper()
	words := routetest.Words(t, step)
	var stdout, stderr bytes.Buffer
	code := 0
	switch strings.Join(words[:min(3, len(words))], " ") {
	case "bench worktree create":
		code = worktree.CreateCommand(f.dir, f.home, words[3:], &stdout, &stderr)
	default:
		if words[1] != "commitment" {
			t.Fatalf("printed step %q runs no verb that a commitment route names", step)
		}
		var out string
		out, code = commitcmd.Command(f.dir, words[2:])
		stdout.WriteString(out)
	}
	if code != 0 {
		t.Fatalf("printed step %q = (%d, %q, %q), want exit 0", step, code, stdout.String(), stderr.String())
	}
	return stdout.String()
}

// decisionRoute is the command that the commitment-decision route names, as the spec's face
// table states it.
const decisionRoute = "bench commitment plan --input "

// start is the verb's start of outcome for request with deliverable.
func start(outcome, request, deliverable string) []string {
	return []string{"start", "--outcome", outcome, "--request", request, "--deliverable", deliverable}
}

func commitmentFaceFixtures() []commitmentFaceFixture {
	quoted := sanitize.ShellQuote
	return []commitmentFaceFixture{
		{
			// The caller starts from the primary checkout, which no assignment owns. The
			// operator names its own request and label in the printed create, and the start
			// reruns in the checkout that the create made.
			face: refusalroute.CommitmentNeedsAssignment,
			build: func(t *testing.T) commitmentSet {
				root := commitmenttest.Staged(t, "A")
				return commitmentSet{root: root, dir: root, args: start("A", "owned", "specs/A/spec.md")}
			},
			contains: "bench worktree create --request ",
			fill:     func(*testing.T, commitmentSet) []string { return []string{"<request>", "owned", "<label>", "owned"} },
			after: func(t *testing.T, f commitmentSet, _ string) commitmentSet {
				checkout, err := worktree.TreeTarget(f.root, "owned")
				if err != nil {
					t.Fatalf("created assignment: %v", err)
				}
				f.dir = checkout
				return f
			},
		},
		{
			// One criterion of the evidence is unmet. The operator corrects the evidence at
			// its path, and the printed verification records it.
			face: refusalroute.CommitmentVerifyEvidence,
			build: func(t *testing.T) commitmentSet {
				root := commitmenttest.SeedMilestone(t, commitmenttest.MilestoneSpec, commitmenttest.TicketsFolder)
				evidence := commitmenttest.Evidence(t, root, func(evidence *commitment.MilestoneEvidence) {
					evidence.Results[1].Result = commitment.ResultUnmet
				})
				return commitmentSet{root: root, dir: root, input: evidence, args: []string{"verify", "--milestone", commitmenttest.ClosureMilestone, "--evidence", evidence}}
			},
			contains: "bench commitment verify --milestone " + quoted(commitmenttest.ClosureMilestone) + " --evidence ",
			fill: func(t *testing.T, f commitmentSet) []string {
				copyFile(t, commitmenttest.Evidence(t, f.root, nil), f.input)
				return []string{"<file>", quoted(f.input)}
			},
		},
		{
			// RR45: the outcome is outside the active milestone. The reviewer withdraws the
			// edit that took it out, which is a commitment change.
			face: refusalroute.CommitmentDecision,
			build: func(t *testing.T) commitmentSet {
				root := commitmenttest.Staged(t, "Z", "A")
				commitmenttest.EditPolicy(t, root, func(policy *commitment.Policy) {
					policy.Milestones[0].Outcomes = policy.Milestones[0].Outcomes[1:]
				})
				commitmenttest.Commit(t, root, "take Z out of the milestone")
				return commitmentSet{root: root, dir: commitmenttest.Assignment(t, root, "decision"), args: start("Z", "decision", "specs/Z/spec.md")}
			},
			prefix:   routetest.ReviewerMarker,
			contains: decisionRoute,
			carry: map[int]func(*testing.T, commitmentSet){0: func(t *testing.T, f commitmentSet) {
				gittest.Output(t, f.root, "revert", "--no-edit", "HEAD")
			}},
		},
		{
			// The adoption lists a run that the ledger does not hold. The operator reads the
			// run it means from the printed inventory and lists that run.
			face:     refusalroute.CommitmentRunUnknown,
			build:    unknownRunSet,
			contains: "bench commitment inventory",
			after: func(t *testing.T, f commitmentSet, inventory string) commitmentSet {
				document, err := axitest.DecodeDocument(inventory)
				if err != nil {
					t.Fatal(err)
				}
				rows, err := document.Rows("commitment_inventory")
				if err != nil {
					t.Fatal(err)
				}
				for _, row := range rows {
					if fields := row.(map[string]any); fields["kind"] == "run" && fields["identity"] == intent.RequestDigest(commitmenttest.LegacyRun) {
						writeAdoption(t, f.input, fields["id"].(string))
						return f
					}
				}
				t.Fatalf("inventory = %q, want the legacy run", inventory)
				return f
			},
		},
		planInputFixture("", "adoption plan.json", quoted),
		planInputFixture("input path that is not line-safe", "adoption\x1bplan.json", func(string) string { return "<file>" }),
		{
			// The board names a roadmap row whose detail file the default branch does not
			// hold, a cause with no face of its own. The reviewer adds the detail file.
			face: refusalroute.CommitmentHandback,
			build: func(t *testing.T) commitmentSet {
				root := commitmenttest.Staged(t, "A")
				commitmenttest.Write(t, root, "ROADMAP.md", "# Roadmap\n\n## Parked\n\n**FT1 — A**\n\n## Recommended sequence\n\n1. A\n")
				commitmenttest.Commit(t, root, "board without its detail")
				return commitmentSet{root: root, dir: root, args: []string{"inventory"}}
			},
			prefix:   routetest.ReviewerMarker,
			contains: "; then bench commitment inventory",
			carry: map[int]func(*testing.T, commitmentSet){0: func(t *testing.T, f commitmentSet) {
				commitmenttest.Write(t, f.root, "roadmap/FT1.md", "**FT1 — A**\n\nKeep the A obligation.\n")
				commitmenttest.Commit(t, f.root, "add the detail")
			}},
		},
		{
			// The blocker reason spans two lines, a verb error that raises no face, so the
			// refusal fails closed to the reviewer. The cause is the operand itself, and its
			// re-run prints the reason's placeholder, which the reviewer fills with one line.
			face: refusalroute.CommitmentHandback, cause: "untyped verb error",
			build: func(t *testing.T) commitmentSet {
				root := commitmenttest.Staged(t, "A")
				return commitmentSet{root: root, dir: root, args: block("A", "two\nlines")}
			},
			prefix:   routetest.ReviewerMarker,
			contains: "; then bench commitment block --outcome " + quoted("A") + " --reason <text>",
			carry:    map[int]func(*testing.T, commitmentSet){0: func(*testing.T, commitmentSet) {}},
			fill:     func(*testing.T, commitmentSet) []string { return []string{"<text>", quoted(oneLine)} },
			after: func(_ *testing.T, f commitmentSet, _ string) commitmentSet {
				f.args = block("A", oneLine)
				return f
			},
		},
	}
}

// Each start or approval cause that only a commitment change clears raises the decision face
// at its source. An unparsed approval operand or an unknown plan id hands back. The walk
// proves each route, so this test reads only which face each cause prints.
func TestCommitmentCausesRaiseTheirFaces(t *testing.T) {
	faces := map[string]refusalroute.Face{}
	for _, face := range refusalroute.Faces(refusalroute.Commitment) {
		faces[face.Name] = face
	}
	decision, approval := []string{decisionRoute}, []string{"; then bench commitment approve "}
	for _, row := range []struct {
		name, face string
		names      []string
		set        func(t *testing.T, root string) (dir string, args []string)
	}{
		{"unapproved deliverable", refusalroute.CommitmentDecision, decision, func(t *testing.T, root string) (string, []string) {
			return commitmenttest.Assignment(t, root, "decision"), start("A", "decision", "specs/B/spec.md")
		}},
		{"precedes", refusalroute.CommitmentDecision, decision, func(t *testing.T, root string) (string, []string) {
			return commitmenttest.Assignment(t, root, "decision"), start("B", "decision", "specs/B/spec.md")
		}},
		{"no active milestone", refusalroute.CommitmentDecision, decision, func(t *testing.T, root string) (string, []string) {
			commitmenttest.EditPolicy(t, root, func(policy *commitment.Policy) { policy.ActiveMilestone = "" })
			commitmenttest.Commit(t, root, "no active milestone")
			return commitmenttest.Assignment(t, root, "decision"), start("A", "decision", "specs/A/spec.md")
		}},
		{"approval after another approval", refusalroute.CommitmentDecision, decision, func(t *testing.T, root string) (string, []string) {
			dir, first, args := commitmenttest.Planning(t, root), approvalOf(t, root, commitmenttest.Reworded(t, root, "A is read.")), approvalOf(t, root, commitmenttest.Reworded(t, root, "A is kept."))
			if out, code := commitcmd.Command(dir, first); code != 0 {
				t.Fatalf("first approval = (%d, %q), want exit 0", code, out)
			}
			return dir, args
		}},
		{"approval after the policy moved", refusalroute.CommitmentDecision, decision, func(t *testing.T, root string) (string, []string) {
			args := approvalOf(t, root, commitmenttest.Reworded(t, root, "A is read."))
			commitmenttest.EditPolicy(t, root, func(policy *commitment.Policy) { policy.Milestones[0].Outcomes[1].Criteria[0].Text = "B moved." })
			commitmenttest.Commit(t, root, "move the policy past the plan")
			return commitmenttest.Planning(t, root), args
		}},
		{"approval after a bound source changed", refusalroute.CommitmentDecision, decision, func(t *testing.T, root string) (string, []string) {
			args := approvalOf(t, root, commitmenttest.Reworded(t, root, "A is read."))
			commitmenttest.Write(t, root, "specs/B/spec.md", commitmenttest.StagedBody+"\nB changed.\n")
			commitmenttest.Commit(t, root, "change a bound source past the plan")
			return commitmenttest.Planning(t, root), args
		}},
		{"approval after a listed run changed", refusalroute.CommitmentDecision, decision, func(t *testing.T, _ string) (string, []string) {
			root, run := commitmenttest.Legacy(t)
			owner, _ := intent.AssignmentForWorktree(run)
			args := approvalOf(t, root, commitmenttest.LegacyAdoption(owner.ID).Encode(t))
			gittest.Output(t, run, "rm", "-q", commitmenttest.LegacyScope)
			commitmenttest.Commit(t, run, "drop the listed scope past the plan")
			return run, args
		}},
		{"approval operand", refusalroute.CommitmentHandback, approval, func(t *testing.T, root string) (string, []string) {
			return commitmenttest.Planning(t, root), []string{"approve", "--plan", "p", "--decision", "d", "--delayed", " A", "--removed", "none"}
		}},
		{"unknown plan", refusalroute.CommitmentHandback, approval, func(t *testing.T, root string) (string, []string) {
			return commitmenttest.Planning(t, root), []string{"approve", "--plan", "p", "--decision", "d", "--delayed", "none", "--removed", "none"}
		}},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Setenv("BENCH_HOME", t.TempDir())
			dir, args := row.set(t, commitmenttest.Staged(t, "A", "B"))
			out, code := commitcmd.Command(dir, args)
			next, printed := routetest.NextCell(out)
			if code != 1 || !printed {
				t.Fatalf("%s = (%d, %q), want exit 1 and one next cell", row.name, code, out)
			}
			routetest.Steps(t, faces[row.face], next, "", row.names)
		})
	}
}

// approvalOf plans input in root and returns the verb's approval of that plan.
func approvalOf(t *testing.T, root string, input []byte) []string {
	t.Helper()
	plan, err := (commitrepo.Store{Root: root}).Plan(input)
	if err != nil {
		t.Fatal(err)
	}
	return []string{"approve", "--plan", plan.ID, "--decision", "d", "--delayed", "none", "--removed", "none"}
}

// oneLine is the blocker reason that the reviewer fills in.
const oneLine = "waits on review"

// block is the verb's block of outcome for reason.
func block(outcome, reason string) []string {
	return []string{"block", "--outcome", outcome, "--reason", reason}
}

// planInputFixture is a commitment-plan-input fixture: the input file at name does not
// parse. The route re-runs the plan with the input as named. The walk corrects the file,
// and the route's re-run then plans it. An input that prints its placeholder is the path
// the operator holds.
func planInputFixture(cause, name string, named func(string) string) commitmentFaceFixture {
	return commitmentFaceFixture{
		face: refusalroute.CommitmentPlanInput, cause: cause,
		build: func(t *testing.T) commitmentSet {
			root := commitmenttest.Staged(t, "A")
			input := filepath.Join(t.TempDir(), name)
			commitmenttest.Write(t, filepath.Dir(input), name, "not a proposal")
			return commitmentSet{root: root, dir: root, input: input, args: []string{"plan", "--input", input}}
		},
		names: func(f commitmentSet) []string {
			return []string{"; then bench commitment plan --input " + named(f.input)}
		},
		carry: map[int]func(*testing.T, commitmentSet){0: func(t *testing.T, f commitmentSet) {
			commitmenttest.Write(t, filepath.Dir(f.input), filepath.Base(f.input), string(commitmenttest.Reworded(t, f.root, "The outcome is delivered and read.")))
		}},
		fill: func(_ *testing.T, f commitmentSet) []string { return []string{"<file>", sanitize.ShellQuote(f.input)} },
	}
}

// unknownRunSet is the run-unknown fixture: the legacy repository, and an adoption that lists
// its run under an id the ledger does not hold.
func unknownRunSet(t *testing.T) commitmentSet {
	root, _ := commitmenttest.Legacy(t)
	input := filepath.Join(t.TempDir(), "adoption.json")
	writeAdoption(t, input, strings.Repeat("d", 32))
	return commitmentSet{root: root, dir: root, input: input, args: []string{"plan", "--input", input}}
}

// writeAdoption writes at path the legacy adoption that lists the legacy run under the id run.
func writeAdoption(t *testing.T, path, run string) {
	t.Helper()
	commitmenttest.Write(t, filepath.Dir(path), filepath.Base(path), string(commitmenttest.LegacyAdoption(run).Encode(t)))
}

// copyFile replaces the file at to with the content of the file at from.
func copyFile(t *testing.T, from, to string) {
	t.Helper()
	data, err := os.ReadFile(from)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(to, data, 0o644); err != nil {
		t.Fatal(err)
	}
}
