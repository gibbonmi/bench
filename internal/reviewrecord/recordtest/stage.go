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

// Staging prepares a delegated one-chunk fixture in an assignment worktree with stagePlan as
// its chunk plan, through Pooled. It answers the fixture and the home.
func Staging(t *testing.T) (*Fixture, string) {
	t.Helper()
	return Pooled(t, 1, func(plan *rr.Plan) {
		plan.Chunks[0].Verification = append([]rr.Requirement{}, stagePlan...)
	}, Delegate)
}

// Pooled prepares a fixture with count planned chunks in an assignment worktree of the pool
// below a new Bench home. BENCH_HOME names that home, so a stage finds the assignment scope.
// The repository of the pool holds the default branch, so a fold can follow a chunk. It
// answers the fixture and the home.
func Pooled(t *testing.T, count int, options ...Option) (*Fixture, string) {
	t.Helper()
	home := t.TempDir()
	t.Setenv(benchhome.Env, home)
	repo, checkout, _ := responseboundtest.AssignmentCheckout(t, home)
	f := Prepare(t, checkout, count, Spec, Body, options...)
	f.primary = repo
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

// Run is a staged run of `bench test` that printed output with exit at the current tip in a
// clean checkout.
func (f *Fixture) Run(exit int, output string) responsebound.Staged {
	f.T.Helper()
	return responsebound.Staged{Retained: responsebound.Retained{Verb: "test", Args: []string{}, Exit: exit, Head: f.Tip(), Dirty: []string{}, Output: []byte(output)}}
}

// Stage writes staged as the staged excerpt of requirement of chunk in the fixture's
// assignment scope below home.
func (f *Fixture) Stage(home, chunk, requirement string, staged responsebound.Staged) {
	f.T.Helper()
	slug, err := rr.Slug(f.spec)
	if err == nil {
		err = responsebound.Stage(home, f.Root, slug, chunk, requirement, staged)
	}
	if err != nil {
		f.T.Fatal(err)
	}
}
