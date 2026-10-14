package repository_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/refusalroute"
)

// Each row commits production paths in an unbound assignment with no light-path ticket. The
// refusal keeps the face of the binding fault and names the paths that need a binding: all
// of them up to five, else the first five and the total. A planning path is never named. The
// printed route of the face states the three ways out ahead of the bind command.
func TestUnboundRefusalNamesProductionPaths(t *testing.T) {
	root, _ := lightPathRepo(t)
	many := []string{"internal/x/a.go", "internal/x/b.go", "internal/x/c.go", "internal/x/d.go", "internal/x/e.go", "internal/x/f.go", "internal/x/g.go"}
	for _, row := range []struct {
		name         string
		steps        []func(*testing.T, string)
		want, absent []string
	}{
		{name: "one-path", steps: steps(write("internal/x/y.go")), want: []string{unbound + ` for production path "internal/x/y.go"`}},
		{name: "planning-path-omitted", steps: steps(write("internal/x/y.go", "capture/IDEAS.md")), want: []string{`"internal/x/y.go"`}, absent: []string{"IDEAS.md"}},
		{name: "five-paths", steps: steps(write(many[:5]...)), want: []string{`"internal/x/a.go", "internal/x/b.go", "internal/x/c.go", "internal/x/d.go", "internal/x/e.go"`}, absent: []string{"more"}},
		{name: "over-the-bound", steps: steps(write(many...)), want: []string{`"internal/x/e.go", and 2 more (7 total)`}, absent: []string{"internal/x/f.go", "internal/x/g.go"}},
	} {
		t.Run(row.name, func(t *testing.T) {
			worktree := commitmenttest.Assignment(t, root, row.name)
			for _, step := range row.steps {
				step(t, worktree)
			}
			commitmenttest.Commit(t, worktree, row.name)
			err := (commitrepo.Store{Root: worktree}).AuthorizeCandidate(gittest.Output(t, worktree, "rev-parse", "HEAD^{tree}"))
			if err == nil {
				t.Fatal("AuthorizeCandidate = nil, want the unbound refusal")
			}
			var raised refusalroute.Raised
			if !errors.As(err, &raised) || raised.Name != refusalroute.CommitmentUnbound {
				t.Fatalf("AuthorizeCandidate = %v, want a refusal that raises %q", err, refusalroute.CommitmentUnbound)
			}
			for _, want := range row.want {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("refusal = %q, want it to hold %q", err, want)
				}
			}
			for _, absent := range row.absent {
				if strings.Contains(err.Error(), absent) {
					t.Errorf("refusal = %q, want no %q", err, absent)
				}
			}
			route := refusalroute.Printed("", err, nil).Route
			for _, want := range []string{"remove the production paths from the change", "one-ticket light-path folder", "; bench commitment start --outcome "} {
				if !strings.Contains(route, want) {
					t.Errorf("route = %q, want it to hold %q", route, want)
				}
			}
		})
	}
}
