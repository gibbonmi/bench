// Stack route walk: each stack fixture prints a stack face, the walk carries out the printed
// route, and the rerun at its end lands. TestLandingFacesFollowTheirRoutes runs this loop.
package worktree

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/refusalroute/routetest"
)

// stackRouteFixture produces one stack face for the route walk. cause keys the fixture apart
// from another fixture of its face and names its subtest. build answers the fixture, the
// overrides of the call, and the 0-based member that the face refuses. spacedDir puts the
// stack file under a directory whose name holds a space. carry is the fixture's own means
// for a printed step, keyed by the step's index in the face's route.
type stackRouteFixture struct {
	face, cause string
	spacedDir   bool
	build       func(t *testing.T, prefix string) (f stackFixture, overrides []string, refused int)
	carry       map[int]func(t *testing.T, w stackWalk)
}

// stackWalk is what a carried step reads: the fixture, the stack file, the refused member,
// and the run that printed the face.
type stackWalk struct {
	f       stackFixture
	file    string
	refused int
	first   verbResult
}

// member is the refused member of the walk.
func (w stackWalk) member() stackFixtureMember { return w.f.members[w.refused] }

// stackRouteFixtures are the stack fixtures of the route walk. The registry census reads the
// face of each one too.
func stackRouteFixtures() []stackRouteFixture {
	red := func(n, member int) func(t *testing.T, prefix string) (stackFixture, []string, int) {
		return func(t *testing.T, prefix string) (stackFixture, []string, int) {
			return redStackFixture(t, prefix, n, member), nil, member
		}
	}
	// overridden is three members with member 2 red, called with an override on the 0-based
	// member.
	overridden := func(member int) func(t *testing.T, prefix string) (stackFixture, []string, int) {
		return func(t *testing.T, prefix string) (stackFixture, []string, int) {
			f := redStackFixture(t, prefix, 3, 1)
			return f, []string{f.members[member].request + "=" + advanceMember(t, f.members[member])}, 1
		}
	}
	return []stackRouteFixture{
		{face: faceStackRed, cause: "red-middle-of-three", spacedDir: true, build: red(3, 1), carry: stackRedCarry},
		{face: faceStackRed, cause: "red-last-of-two", build: red(2, 1), carry: stackRedCarry},
		{face: faceStackRed, cause: "red-first-of-two", build: red(2, 0), carry: stackRedCarry},
		{face: faceStackRed, cause: "red-after-an-override", build: overridden(0), carry: stackRedCarry},
		// The fold moves the tip of member 3, so the rerun drops its override.
		{face: faceStackRed, cause: "red-before-an-override", build: overridden(2), carry: stackRedCarry},
	}
}

// stackRedCarry carries the instructions of the stack-red route: the repair of the red bytes
// under a planted commit lane, the review that takes the lane away, and the fold.
var stackRedCarry = map[int]func(t *testing.T, w stackWalk){
	0: func(t *testing.T, w stackWalk) {
		m := w.member()
		mustWrite(t, filepath.Join(m.creation.Path, m.file), []byte("repaired\n"), 0o644)
		plantSourceCommitLane(t, m.creation)
	},
	2: func(t *testing.T, w stackWalk) {
		mustRemove(t, filepath.Join(w.member().creation.Path, ".bench", "phases.json"))
	},
	3: foldLandedStack,
}

// foldLandedStack carries out the fold instruction: when members landed, it removes each
// landed member from the stack file, folds the destination into each remaining member, and
// sets the new source_tip of each folded member in the stack file.
func foldLandedStack(t *testing.T, w stackWalk) {
	t.Helper()
	landed, _ := recordField(w.first.stdout, "stack{", "landed", "worktree")
	k, err := strconv.Atoi(landed)
	mustNoError(t, err)
	if k == 0 {
		return
	}
	data, err := os.ReadFile(w.file)
	mustNoError(t, err)
	var doc stackFile
	mustNoError(t, json.Unmarshal(data, &doc))
	doc.Members = doc.Members[k:]
	for i, m := range doc.Members {
		gitRun(t, m.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "merge", "-q", "--no-ff", "-m", "fold the destination", "main")
		doc.Members[i].SourceTip = gitOutput(t, m.Path, "rev-parse", "HEAD")
	}
	doc.write(t, w.file)
}

// stackRouteKeys are the face and cause of each stack fixture, for routetest.Fixtures.
func stackRouteKeys() [][2]string {
	var keys [][2]string
	for _, fixture := range stackRouteFixtures() {
		keys = append(keys, [2]string{fixture.face, fixture.cause})
	}
	return keys
}

// followStackRoutes runs each stack fixture through the route walk.
func followStackRoutes(t *testing.T, wrapper string, faces map[string]refusalroute.Face) {
	for _, fixture := range stackRouteFixtures() {
		t.Run(fixture.cause, func(t *testing.T) {
			t.Parallel()
			followStackRoute(t, wrapper, faces[fixture.face], fixture)
		})
	}
}

// followStackRoute drives one stack fixture: the stack prints the face for the refused
// member, the walk carries out the printed route, and the rerun lands. A stack rerun ends
// with the stack record of every remaining member, and a single rerun with its landed record.
func followStackRoute(t *testing.T, wrapper string, face refusalroute.Face, fixture stackRouteFixture) {
	f, overrides, refused := fixture.build(t, "stack-route-"+fixture.cause)
	file, args := f.stack(t, overrides...)
	if fixture.spacedDir {
		spaced := filepath.Join(t.TempDir(), "stack dir")
		mustMkdirAll(t, spaced, 0o755)
		args[1] = filepath.Join(spaced, "stack.json")
		mustNoError(t, os.Rename(file, args[1]))
		file = args[1]
	}
	w := stackWalk{f: f, file: file, refused: refused, first: runVerb(t, verbLand, f.call(args...))}
	m := w.member()
	next := stackRefusalNext(t, w.first.stdout, strconv.Itoa(refused+1), m)
	steps := routetest.Steps(t, face, next, "", []string{m.creation.Assignment.Label})
	fill := func(step string) string {
		return strings.NewReplacer(operatorSlots(t, m.creation.Path, stackRepairedTip, repairedTipArg)...).Replace(step)
	}
	carry := func(index int) (func(), bool) {
		step, carried := fixture.carry[index]
		return func() { step(t, w) }, carried
	}
	last := routetest.Follow(t, face, steps, carry, func(step string) verbResult {
		return runPrintedStep(t, wrapper, f.repoHome, fill(step))
	})
	if strings.HasPrefix(steps[len(steps)-1], "bench worktree land --stack ") {
		if _, printed := recordField(last.stdout, "stack{", "landed"); last.exit != 0 || !printed {
			t.Fatalf("%s stack rerun = (%d, %q, %q), want exit 0 and a stack record", fixture.cause, last.exit, last.stdout, last.stderr)
		}
		return
	}
	if _, landed := landedField(last.stdout, "worktree"); last.exit != 0 || !landed {
		t.Fatalf("%s rerun = (%d, %q, %q), want exit 0 and a landed record", fixture.cause, last.exit, last.stdout, last.stderr)
	}
}
