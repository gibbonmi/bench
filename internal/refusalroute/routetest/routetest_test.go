package routetest

import (
	"fmt"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/refusalroute"
)

// recorder is the testing.TB a rule's check runs against. Its embedded TB is nil, so a
// check that reaches past the overrides panics instead of passing silently.
type recorder struct {
	testing.TB
	failure string
}

func (r *recorder) Helper() {}

func (r *recorder) Errorf(format string, args ...any) {
	if r.failure == "" {
		r.failure = fmt.Sprintf(format, args...)
	}
}

// Fatalf ends the check as it ends a test.
func (r *recorder) Fatalf(format string, args ...any) {
	r.Errorf(format, args...)
	runtime.Goexit()
}

// failure runs check against a recorder and returns the first failure it reports, or the
// empty string when the check passes.
func failure(check func(t testing.TB)) string {
	r := &recorder{}
	done := make(chan struct{})
	go func() {
		defer close(done)
		check(r)
	}()
	<-done
	return r.failure
}

// requireVerdict requires the check to fail exactly when fails holds.
func requireVerdict(t *testing.T, name string, fails bool, check func(t testing.TB)) {
	t.Helper()
	if got := failure(check); (got != "") != fails {
		t.Errorf("%s: failure = %q, want a failure %v", name, got, fails)
	}
}

var (
	agentFace    = refusalroute.Face{Name: "agent-face", Authority: refusalroute.Agent, Route: []refusalroute.Step{instruction, command}}
	reviewerFace = refusalroute.Face{Name: "reviewer-face", Authority: refusalroute.Reviewer, Route: []refusalroute.Step{instruction, command}}
	instruction  = refusalroute.Instruction(refusalroute.Text("repair the cause"))
	command      = refusalroute.Command(refusalroute.Text("bench gate"))
)

func TestNextReadsTheOneRouteLine(t *testing.T) {
	type read struct {
		route string
		one   bool
	}
	for output, want := range map[string]read{
		"error: refused\nnext=bench gate\n":      {"bench gate", true},
		"next=bench gate":                        {"bench gate", true},
		"error: refused\n":                       {},
		"next=bench doctor\nnext=bench gate\n":   {},
		"refused{reason=x,next=bench gate}\n":    {},
		"error: refused\n next=bench gate\n":     {},
		"error: refused\nnext=\nnext=bench gate": {},
	} {
		if route, one := Next(output); route != want.route || one != want.one {
			t.Errorf("Next(%q) = (%q, %v), want (%q, %v)", output, route, one, want.route, want.one)
		}
	}
}

func TestStepsSplitsTheMarkerAndCountsTheSteps(t *testing.T) {
	steps := []string{"repair the cause", "bench gate"}
	agent, reviewer := steps[0]+"; then "+steps[1], ReviewerMarker+steps[0]+"; then "+steps[1]
	var got []string
	requireVerdict(t, "agent route", false, func(t testing.TB) { got = Steps(t, agentFace, agent, "", nil) })
	if !slices.Equal(got, steps) {
		t.Errorf("agent steps = %q, want %q", got, steps)
	}
	requireVerdict(t, "reviewer route", false, func(t testing.TB) { got = Steps(t, reviewerFace, reviewer, "", nil) })
	if !slices.Equal(got, steps) {
		t.Errorf("reviewer steps = %q, want %q", got, steps)
	}
	requireVerdict(t, "route after its preface", false, func(t testing.TB) { Steps(t, agentFace, "first; "+agent, "first", []string{"bench gate"}) })
	requireVerdict(t, "agent route with the marker", true, func(t testing.TB) { Steps(t, agentFace, reviewer, "", nil) })
	requireVerdict(t, "reviewer route with no marker", true, func(t testing.TB) { Steps(t, reviewerFace, agent, "", nil) })
	requireVerdict(t, "route with a missing step", true, func(t testing.TB) { Steps(t, agentFace, steps[1], "", nil) })
	requireVerdict(t, "route with no named value", true, func(t testing.TB) { Steps(t, agentFace, agent, "", []string{"bench commit"}) })
}

func TestFollowCarriesOrRunsEachStep(t *testing.T) {
	steps := []string{"repair the cause", "bench gate"}
	carryFirst := func(carried *bool) func(int) (func(), bool) {
		return func(index int) (func(), bool) { return func() { *carried = true }, index == 0 }
	}
	carryAll := func(int) (func(), bool) { return func() {}, true }
	carryNone := func(int) (func(), bool) { return nil, false }
	var carried bool
	var ran []string
	run := func(step string) string { ran = append(ran, step); return "ran " + step }
	var last string
	requireVerdict(t, "carried instruction and run command", false, func(t testing.TB) { last = Follow(t, agentFace, steps, carryFirst(&carried), run) })
	if !carried || !slices.Equal(ran, steps[1:]) || last != "ran bench gate" {
		t.Errorf("walk carried %v, ran %q, and returned %q; want the instruction carried and the command run last", carried, ran, last)
	}
	requireVerdict(t, "carried reviewer command", false, func(t testing.TB) { Follow(t, reviewerFace, steps, carryAll, run) })
	requireVerdict(t, "carried agent command", true, func(t testing.TB) { Follow(t, agentFace, steps, carryAll, run) })
	requireVerdict(t, "instruction no fixture carries", true, func(t testing.TB) { Follow(t, agentFace, steps, carryNone, run) })
}

func TestWordsAcceptOnlyOneFilledBenchCommand(t *testing.T) {
	var words []string
	requireVerdict(t, "one Bench command", false, func(t testing.TB) { words = Words(t, "bench commit --in 'label' -m 'msg' -- 'a.txt'") })
	if want := []string{"bench", "commit", "--in", "label", "-m", "msg", "--", "a.txt"}; !slices.Equal(words, want) {
		t.Errorf("words = %q, want %q", words, want)
	}
	for name, step := range map[string]string{
		"two commands":    "bench gate; bench gate",
		"a pipeline":      "bench gate | cat",
		"an open slot":    "bench commit -m <msg>",
		"a quoted slot":   "bench commit -m '<msg>'",
		"an open quote":   "bench commit -m 'msg",
		"no Bench verb":   "git status",
		"a bare launcher": "bench",
	} {
		requireVerdict(t, name, true, func(t testing.TB) { Words(t, step) })
	}
}

func TestDiagnosticIsTheBareDoctor(t *testing.T) {
	for words, want := range map[string]bool{"bench doctor": true, "bench doctor --fix": false, "bench gate": false} {
		if got := Diagnostic(strings.Fields(words)); got != want {
			t.Errorf("Diagnostic(%q) = %v, want %v", words, got, want)
		}
	}
}

func TestFixturesMatchTheRegistry(t *testing.T) {
	var keys [][2]string
	for _, face := range refusalroute.Faces(refusalroute.Commit) {
		keys = append(keys, [2]string{face.Name, ""})
	}
	var faces map[string]refusalroute.Face
	requireVerdict(t, "a fixture for each face", false, func(t testing.TB) { faces = Fixtures(t, refusalroute.Commit, keys) })
	if len(faces) != len(keys) {
		t.Errorf("faces = %d, want the %d registered commit faces", len(faces), len(keys))
	}
	requireVerdict(t, "a second cause of one face", false, func(t testing.TB) {
		Fixtures(t, refusalroute.Commit, append(slices.Clone(keys), [2]string{keys[0][0], "second"}))
	})
	requireVerdict(t, "a face with no fixture", true, func(t testing.TB) { Fixtures(t, refusalroute.Commit, keys[1:]) })
	requireVerdict(t, "a fixture of no face", true, func(t testing.TB) {
		Fixtures(t, refusalroute.Commit, append(slices.Clone(keys), [2]string{"no-such-face", ""}))
	})
	requireVerdict(t, "two fixtures of one cause", true, func(t testing.TB) {
		Fixtures(t, refusalroute.Commit, append(slices.Clone(keys), keys[0]))
	})
}
