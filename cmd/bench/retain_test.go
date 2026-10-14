package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/poolkey"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/responsebound/responseboundtest"
)

// retainFixture is one assignment checkout of a new repository under a private Bench home.
type retainFixture struct {
	home, repo, checkout, id string
}

// useRetainFixture registers definitions, makes an assignment checkout under a private Bench
// home, and runs the test from that checkout.
func useRetainFixture(t *testing.T, definitions ...commandDefinition) retainFixture {
	t.Helper()
	home := t.TempDir()
	t.Setenv(benchhome.Env, home)
	repo, checkout, id := responseboundtest.AssignmentCheckout(t, home)
	t.Chdir(checkout)
	old := commandRegistry
	t.Cleanup(func() { commandRegistry = old })
	commandRegistry = definitions
	return retainFixture{home: home, repo: repo, checkout: checkout, id: id}
}

// boundedCommand is a bounded, tree-scoped command whose handler prints lines numbered lines
// and exits with exit, so its response carries the identity row.
func boundedCommand(name string, lines, exit int) commandDefinition {
	return commandDefinition{Name: name, Inventory: publicInventory(), Bound: boundResponse, Scope: scopeTree, Run: exitingHandler(lines, exit)}
}

func (f retainFixture) retained(t *testing.T) responsebound.Retained {
	t.Helper()
	retained, err := responsebound.ReadRetained(f.home, f.checkout)
	if err != nil {
		t.Fatalf("read the retained response: %v", err)
	}
	return retained
}

// RT1, RT2, RT3: a call within the bound and a spilled call each leave their complete output
// without the identity row or a spill frame, with the typed verb, the arguments after it, the
// exit code, and the worktree's commit and dirty paths.
func TestRetainKeepsTheCompleteResponse(t *testing.T) {
	for _, row := range []struct {
		name        string
		lines, exit int
	}{
		{"within the bound", 3, 3},
		{"spilled", 25, 0},
	} {
		t.Run(row.name, func(t *testing.T) {
			f := useRetainFixture(t, boundedCommand("test", row.lines, row.exit))
			if err := os.WriteFile(filepath.Join(f.checkout, "note.txt"), []byte("dirty\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, code := dispatch("test", "--package", "./x"); code != row.exit {
				t.Fatalf("test exit = %d, want %d", code, row.exit)
			}
			want := responsebound.Retained{
				Verb: "test", Args: []string{"--package", "./x"}, Exit: row.exit,
				Head: gittest.Output(t, f.checkout, "rev-parse", "HEAD"), Dirty: []string{"note.txt"},
				Output: []byte(numberedLines(row.lines)),
			}
			if got := f.retained(t); !reflect.DeepEqual(got, want) {
				t.Fatalf("retained = %+v\nwant %+v", got, want)
			}
		})
	}
}

// RT4: a record call, a worktree call, and a retiring leaf each leave the earlier retained
// response in place.
func TestRetainSkipsRecordWorktreeAndRetiringCalls(t *testing.T) {
	retire := func(c Command, _ string, args []string) int { return exitingHandler(2, 0)(c, args) }
	family := commandDefinition{
		Name: "family", Inventory: publicInventory(), Bound: boundResponse, LeafUsage: func() string { return "usage\n" },
		Leaves: []commandLeaf{{Name: "retire", Root: rootNone, Bound: boundResponse, Retires: true, Run: retire}},
	}
	for _, argv := range [][]string{{"record", "chunk"}, {"worktree", "exec"}, {"family", "retire"}} {
		t.Run(argv[0], func(t *testing.T) {
			f := useRetainFixture(t, boundedCommand("test", 1, 0), boundedCommand("record", 2, 0), boundedCommand("worktree", 2, 0), family)
			if _, code := dispatch("test"); code != 0 {
				t.Fatalf("test exit = %d", code)
			}
			before := f.retained(t)
			if _, code := dispatch(argv...); code != 0 {
				t.Fatalf("%q exit = %d", argv, code)
			}
			if got := f.retained(t); !reflect.DeepEqual(got, before) {
				t.Fatalf("retained after %q = %+v, want the test response %+v", argv, got, before)
			}
		})
	}
}

// RT6: a bounded call in the primary checkout writes no response store.
func TestRetainSkipsThePrimaryCheckout(t *testing.T) {
	f := useRetainFixture(t, boundedCommand("test", 3, 0))
	t.Chdir(f.repo)
	if _, code := dispatch("test"); code != 0 {
		t.Fatalf("test exit = %d", code)
	}
	if _, err := os.Stat(filepath.Join(f.home, "responses")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("response store after a primary call = %v, want none", err)
	}
}

// RT8: a symlink at the assignment scope refuses the retention write, and the response and
// the exit code stay the command's own. The same call without the symlink retains, so the
// refusal is the retention write failing and not a path the dispatcher never reached.
func TestRetainFailureKeepsTheResponse(t *testing.T) {
	f := useRetainFixture(t, boundedCommand("test", 3, 3))
	elsewhere := t.TempDir()
	scope := filepath.Join(f.home, "responses", poolkey.Key(f.repo), f.id)
	if err := os.MkdirAll(filepath.Dir(scope), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(elsewhere, scope); err != nil {
		t.Fatal(err)
	}
	if out, code := dispatch("test"); code != 3 || out != numberedLines(3) {
		t.Fatalf("test with a symlinked scope = (%d, %q), want (3, %q)", code, out, numberedLines(3))
	}
	if entries, err := os.ReadDir(elsewhere); err != nil || len(entries) != 0 {
		t.Fatalf("the retention followed the symlink: %v (%v)", entries, err)
	}
	if err := os.Remove(scope); err != nil {
		t.Fatal(err)
	}
	if _, code := dispatch("test"); code != 3 || f.retained(t).Exit != 3 {
		t.Fatalf("test without the symlink exit = %d, want 3 and a retained response", code)
	}
}
