package recordtest

import (
	"io"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/responsebound/responseboundtest"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
)

// stagePlan is the chunk plan of a Staging fixture: a probe requirement, and a requirement
// whose planned command holds a single-quoted operand.
var stagePlan = []rr.Requirement{
	{ID: "tests", Command: "bench test --package ./internal/example", Probe: "omit source check"},
	{ID: "quoted", Command: "bench test --run 'TestA|TestB'"},
}

// Staging prepares a delegated one-chunk fixture in an assignment worktree of the pool below a
// new Bench home, with stagePlan as its chunk plan. BENCH_HOME names that home, so a stage
// finds the assignment scope. It answers the fixture and the home.
func Staging(t *testing.T) (*Fixture, string) {
	t.Helper()
	home := t.TempDir()
	t.Setenv(benchhome.Env, home)
	_, checkout, _ := responseboundtest.AssignmentCheckout(t, home)
	f := Prepare(t, checkout, 1, Spec, Body, func(plan *rr.Plan) {
		plan.Chunks[0].Verification = append([]rr.Requirement{}, stagePlan...)
	}, Delegate)
	return f, home
}

// Retain makes output the retained response of the fixture's worktree, as the call
// `bench <verb> <args>…` printed it with exit at the current tip.
func (f *Fixture) Retain(home, verb string, args []string, exit int, output string) {
	f.T.Helper()
	owner := responsebound.New(home, io.Discard, io.Discard, func() string { return f.Root }, false)
	_, _ = io.WriteString(owner.Stdout(), output)
	owner.Finish()
	err := owner.Retain(home, f.Root, func() (responsebound.Retained, error) {
		return responsebound.Retained{Verb: verb, Args: args, Exit: exit, Head: f.Tip(), Dirty: []string{}}, nil
	})
	if err != nil {
		f.T.Fatal(err)
	}
}
