// Refusal route tests for the gate checkpoint: each checkpoint refusal prints the route of
// its face in the shared refusal-route registry, and that route clears the refusal. The walk
// runs in the external test package, because the commit, the preflight read, and the doctor
// that a route names each import the gate.
package gate_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/adopt"
	"github.com/gibbonmi/bench/internal/commit"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/refusalroute/routetest"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// The caller's own checkpoint arguments: the first chunk, or the completion.
var (
	chunkArgs    = []string{"--checkpoint", recordtest.Spec, "--chunk", "1"}
	completeArgs = []string{"--checkpoint", recordtest.Spec, "--complete"}
)

// gateSet is one checkpoint route fixture: the review fixture whose checkout the checkpoint
// grades, and the label of the active assignment that owns that checkout.
type gateSet struct {
	f     *recordtest.Fixture
	label string
}

// gateFaceFixture produces exactly one gate face. build makes the set and breaks it so the
// checkpoint prints that face. The registry walk requires a fixture for each gate face, so
// a face added with no fixture turns TestCheckpointFacesFollowTheirRoutes red.
type gateFaceFixture struct {
	face string
	// args are the caller's checkpoint arguments.
	args  []string
	build func(t *testing.T) gateSet
	// refuse runs the refused gate, where the caller's checkpoint alone does not reach the
	// face. It answers the exit.
	refuse func(t *testing.T, set gateSet, stderr io.Writer) int
	// contains states what the printed route must hold.
	contains string
	// clear removes the fixture's own fault scaffolding once the face printed, which no
	// operator's tree holds, so the route repairs only the cause the face names.
	clear func(t *testing.T, set gateSet)
	// carry carries out a printed instruction by the fixture's own means, keyed by the
	// step's index in the face's route.
	carry map[int]func(t *testing.T, set gateSet)
	// slots are the replacement pairs for the operator slots of a printed command step.
	slots []string
	// after carries out what the route's last step leaves to the operator.
	after func(t *testing.T, set gateSet)
}

// TestCheckpointFacesFollowTheirRoutes is RR44 and RR68. The registry is the source of the
// gate's face set, so each gate face needs a producing fixture. Each fixture's printed route
// is carried out step by step: an instruction by the fixture's own means, and each command
// step verbatim through the verb's own entry. The checkpoint then reruns out of the face:
// as the route's own last step, or after a route that ends elsewhere.
func TestCheckpointFacesFollowTheirRoutes(t *testing.T) {
	var keys [][2]string
	for _, fixture := range gateFaceFixtures() {
		keys = append(keys, [2]string{fixture.face, ""})
	}
	faces := routetest.Fixtures(t, refusalroute.Gate, keys)
	for _, fixture := range gateFaceFixtures() {
		t.Run(fixture.face, func(t *testing.T) { followGateFace(t, faces[fixture.face], fixture) })
	}
}

// TestCheckpointRouteRendersTheLedgerLabel grades the label slot over the values a ledger can
// hold. A label with a space prints shell-quoted, and a label that is not line-safe prints
// the slot, so no control byte reaches the route line.
func TestCheckpointRouteRendersTheLedgerLabel(t *testing.T) {
	for _, tc := range []struct{ name, label, want string }{
		{"space", "gate route", sanitize.ShellQuote("gate route")},
		{"not line-safe", "gate\x1broute", "<" + refusalroute.FactLabel + ">"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			set := assigned(t, tc.label, true)
			set.f.Write("tracked.txt", "edited\n")
			var stderr bytes.Buffer
			gate.RunCommand(append([]string{set.f.Root}, completeArgs...), io.Discard, &stderr)
			next, printed := printedRoute(stderr.String())
			if !printed || !strings.HasPrefix(next, "bench commit --in "+tc.want+" ") || strings.ContainsRune(next, '\x1b') {
				t.Fatalf("stderr = %q, want a next= route that addresses the worktree as %q", stderr.String(), tc.want)
			}
		})
	}
}

// followGateFace drives one fixture: the checkpoint prints the face, the walk carries out the
// printed route, and the checkpoint reruns out of the face.
func followGateFace(t *testing.T, face refusalroute.Face, fixture gateFaceFixture) {
	set := fixture.build(t)
	var stdout, stderr bytes.Buffer
	var code int
	if fixture.refuse != nil {
		code = fixture.refuse(t, set, &stderr)
	} else {
		code = gate.RunCommand(append([]string{set.f.Root}, fixture.args...), &stdout, &stderr)
	}
	next, printed := printedRoute(stderr.String())
	if code == 0 || !printed || !strings.Contains(next, fixture.contains) {
		t.Fatalf("face %s = (%d, %q), want a refusal whose next= route holds %q", face.Name, code, stderr.String(), fixture.contains)
	}
	steps := routetest.Steps(t, face, next, "", nil)
	if fixture.clear != nil {
		fixture.clear(t, set)
	}
	carry := func(index int) (func(), bool) {
		step, carried := fixture.carry[index]
		return func() { step(t, set) }, carried
	}
	routetest.Follow(t, face, steps, carry, func(step string) int {
		return runGateRouteStep(t, set, strings.NewReplacer(fixture.slots...).Replace(step))
	})
	// A route that ends with the caller's own gate has rerun it, and the step exited 0.
	if strings.HasPrefix(steps[len(steps)-1], "bench gate ") {
		return
	}
	if fixture.after != nil {
		fixture.after(t, set)
	}
	stdout.Reset()
	stderr.Reset()
	if code := gate.RunCommand(append([]string{set.f.Root}, fixture.args...), &stdout, &stderr); code != 0 {
		t.Fatalf("%s rerun = (%d, %q, %q), want exit 0", face.Name, code, stdout.String(), stderr.String())
	}
}

// runGateRouteStep runs one printed command step in process, the way the CLI dispatches it, in
// the fixture's checkout. A tree target must name the fixture's own assignment, which owns
// that checkout. The shared walk states the rules a printed step and its exit meet.
func runGateRouteStep(t *testing.T, set gateSet, step string) int {
	t.Helper()
	words := routetest.Words(t, step)
	t.Chdir(set.f.Root)
	var stdout, stderr bytes.Buffer
	code := 0
	switch strings.Join(words[:min(3, len(words))], " ") {
	case "bench gate --in", "bench commit --in":
		if words[3] != set.label {
			t.Fatalf("printed step %q names the worktree %q, want the fixture's assignment %q", step, words[3], set.label)
		}
		if words[1] == "gate" {
			code = gate.Command(words[4:], nil, &stdout, &stderr)
		} else {
			_, code = commit.Run(words[4:], &stdout, &stderr)
		}
	case "bench preflight review":
		var out string
		out, code = preflight.Command(words[2:])
		stdout.WriteString(out)
	case "bench doctor":
		code = adopt.Doctor(words[2:], io.Discard, io.Discard, "fixture")
	default:
		t.Fatalf("printed step %q runs no verb that a checkpoint route names", step)
	}
	if code != 0 && !routetest.Diagnostic(words) {
		t.Fatalf("printed step %q = (%d, %q, %q), want exit 0", step, code, stdout.String(), stderr.String())
	}
	return code
}

// assigned is a checkpoint fixture whose checkout an active assignment of label owns. With
// complete, the fixture also retains its completion evidence.
func assigned(t *testing.T, label string, complete bool) gateSet {
	t.Helper()
	f := gate.AttachedCheckpointFixture(t, recordtest.Attach)
	if complete {
		gate.RetainCompletion(f)
	}
	commitmenttest.Register(t, f.Root, label)
	return gateSet{f: f, label: label}
}

func gateFaceFixtures() []gateFaceFixture {
	// record is the path of the fixtures' review record, which the spec path names.
	record, _ := rr.RecordPath(recordtest.Spec)
	return []gateFaceFixture{
		evidenceFixture(record),
		dirtyFixture(record),
		{
			// The spec drops its staged status, which only the reviewer restores.
			face:     "checkpoint-composition",
			args:     completeArgs,
			contains: routetest.ReviewerMarker,
			build: func(t *testing.T) gateSet {
				set := assigned(t, "gate-composition", true)
				staged := string(readFile(t, set.f.Root, recordtest.Spec))
				set.f.Write(recordtest.Spec, strings.Replace(staged, "Status: staged\n", "", 1))
				set.f.Commit("drop the staged status")
				return set
			},
			carry: map[int]func(*testing.T, gateSet){0: func(t *testing.T, set gateSet) {
				staged := string(readFile(t, set.f.Root, recordtest.Spec))
				set.f.Write(recordtest.Spec, strings.Replace(staged, "# Example\n\n", "# Example\n\nStatus: staged\n", 1))
				set.f.Commit("restore the staged status")
			}},
		},
		{
			// The tip moves under the run's lock, after the run accepted its subject.
			face:     "checkpoint-tip-moved",
			args:     chunkArgs,
			contains: "bench gate --in " + sanitize.ShellQuote("gate-tip") + " ",
			build:    func(t *testing.T) gateSet { return assigned(t, "gate-tip", false) },
			refuse: func(t *testing.T, set gateSet, stderr io.Writer) int {
				move := func() { set.f.Git("commit", "-q", "--allow-empty", "-m", "move the tip") }
				return gate.RunCheckpointArmed(t, set.f.Root, chunkArgs, move, stderr)
			},
		},
		{
			// The branch that HEAD names is gone, so the checkpoint has no source until the
			// fixture restores it.
			face:     "checkpoint-subject-unavailable",
			args:     completeArgs,
			contains: "bench doctor; then ",
			build: func(t *testing.T) gateSet {
				set := assigned(t, "gate-subject", true)
				set.f.Git("update-ref", "-m", "fixture", "refs/fixture/tip", "HEAD")
				set.f.Git("update-ref", "-d", "refs/heads/main")
				return set
			},
			clear: func(_ *testing.T, set gateSet) { set.f.Git("update-ref", "refs/heads/main", "refs/fixture/tip") },
		},
		{
			// The completion run grades the source tree itself, which is not the published
			// transform of that source. The checkpoint grades the published tree, so the
			// reviewer has nothing left to clear.
			face:     "gate-handback",
			args:     completeArgs,
			contains: routetest.ReviewerMarker,
			build:    func(t *testing.T) gateSet { return assigned(t, "gate-handback", true) },
			refuse: func(_ *testing.T, set gateSet, stderr io.Writer) int {
				ctx := gate.WithCompletion(context.Background(), recordtest.Spec, set.f.Tip())
				return gate.ExecuteTree(ctx, set.f.Root, set.f.Tree(), io.Discard, stderr).ActionExit
			},
			carry: map[int]func(*testing.T, gateSet){0: func(*testing.T, gateSet) {}},
		},
	}
}

// evidenceFixture is the completion-evidence fixture: a delivery that the preflight read
// grades green, with its review record removed. The read reports the gap, and the operator
// writes the record back. The chunk is the one the conformant plan declares.
func evidenceFixture(record string) gateFaceFixture {
	var saved []byte
	return gateFaceFixture{
		face:     "checkpoint-completion-evidence",
		args:     []string{"--checkpoint", recordtest.Spec, "--chunk", "c1"},
		contains: "bench preflight review ",
		build: func(t *testing.T) gateSet {
			t.Setenv("BENCH_GATE", "true")
			root, _ := preflighttest.SeedConformant(t)
			recordtest.RetainSingleChunk(t, root, recordtest.Spec, preflighttest.RunGit(t, "rev-parse", "main"))
			saved = readFile(t, root, record)
			if err := os.Remove(filepath.Join(root, record)); err != nil {
				t.Fatal(err)
			}
			return gateSet{f: &recordtest.Fixture{T: t, Root: root}}
		},
		after: func(_ *testing.T, set gateSet) { set.f.Write(record, string(saved)) },
	}
}

// dirtyFixture is the dirty-checkout fixture: the completion record waits uncommitted in a
// linked worktree that an assignment owns. The label holds a space, so the route quotes it.
func dirtyFixture(record string) gateFaceFixture {
	const label = "gate route"
	return gateFaceFixture{
		face:     "checkpoint-dirty-checkout",
		args:     completeArgs,
		contains: "bench commit --in " + sanitize.ShellQuote(label) + " ",
		slots:    []string{"<msg>", "record", "<path>...", record},
		build: func(t *testing.T) gateSet {
			t.Setenv("BENCH_HOME", t.TempDir())
			t.Setenv("BENCH_KIT", t.TempDir())
			f := gate.AttachedCheckpointFixture(t, linkedAttach)
			commitmenttest.Register(t, f.Root, label)
			commitmenttest.Admit(t, f.Root, label, recordtest.Spec)
			f.Complete()
			f.Save()
			return gateSet{f: f, label: label}
		},
	}
}

// linkedAttach attaches the review fixture to the primary checkout, approves its delivery
// there, and answers the fixture in a linked worktree branched from it, because the commit
// refuses a primary checkout. The linked gate exits 0: the outcome gate copies a run record
// that only the primary checkout's Git directory keeps.
func linkedAttach(t testing.TB, root string, count int) *recordtest.Fixture {
	f := recordtest.Attach(t, root, count)
	f.Write(".bench/gate.sh", "#!/bin/sh\nexit 0\n")
	if err := os.Chmod(filepath.Join(root, ".bench", "gate.sh"), 0o755); err != nil {
		t.Fatal(err)
	}
	commitmenttest.SeedAdmission(t, root, recordtest.Spec)
	f.Commit("approve the delivery")
	linked := *f
	linked.Root = filepath.Join(t.TempDir(), "linked")
	f.Git("worktree", "add", "-q", "-b", "assignment", linked.Root)
	return &linked
}

// printedRoute is the route of a refusal's next= line on stderr.
func printedRoute(stderr string) (string, bool) {
	_, next, printed := strings.Cut(stderr, "\n"+refusalroute.NextField+"=")
	next, _, _ = strings.Cut(next, "\n")
	return next, printed
}

func readFile(t testing.TB, root, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, path))
	if err != nil {
		t.Fatal(err)
	}
	return data
}
