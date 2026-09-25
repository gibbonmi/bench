package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
)

// The expectations here are authored apart from the owner: a bounded response over 10
// lines prints exactly 10. A shared constant would let a changed bound move both sides.

// runBoundFixture registers one public command that prints lines numbered lines and
// exits exit, declared with bound, and runs it through the production dispatcher under a
// private Bench home.
func runBoundFixture(t *testing.T, bound boundDisposition, lines, exit int) (stdout, stderr, home string, code int) {
	t.Helper()
	home = t.TempDir()
	t.Setenv(benchhome.Env, home)
	old := commandRegistry
	t.Cleanup(func() { commandRegistry = old })
	commandRegistry = []commandDefinition{{
		Name:      "fixture",
		Inventory: publicInventory(helpRow{Order: 1, Description: "print numbered lines"}),
		Bound:     bound,
		Run: func(c Command, _ []string) int {
			for i := 1; i <= lines; i++ {
				fmt.Fprintf(c.Stdout, "line %02d\n", i)
			}
			return exit
		},
	}}
	var out, errOut bytes.Buffer
	code = Command{Stdout: &out, Stderr: &errOut, Executable: "bench"}.Run([]string{"fixture"})
	return out.String(), errOut.String(), home, code
}

func numberedLines(count int) string {
	var all strings.Builder
	for i := 1; i <= count; i++ {
		fmt.Fprintf(&all, "line %02d\n", i)
	}
	return all.String()
}

// BO1: a bounded public command that prints 25 lines produces exactly 10 stdout lines.
func TestDispatcherBoundsPublicResponse(t *testing.T) {
	stdout, stderr, _, code := runBoundFixture(t, boundResponse, 25, 0)
	lines := strings.Split(strings.TrimSuffix(stdout, "\n"), "\n")
	if code != 0 || stderr != "" || len(lines) != 10 {
		t.Fatalf("bounded fixture = (%d, %q, %q), want exactly 10 stdout lines", code, stdout, stderr)
	}
	if !strings.HasPrefix(lines[4], "spilled{lines=25,") {
		t.Fatalf("fifth line = %q, want the spill line", lines[4])
	}
}

// BO2: a bounded command that prints exactly 10 lines produces its exact bytes and no
// spill file.
func TestDispatcherPassesBoundaryResponse(t *testing.T) {
	stdout, stderr, home, code := runBoundFixture(t, boundResponse, 10, 0)
	if want := numberedLines(10); code != 0 || stderr != "" || stdout != want {
		t.Fatalf("boundary fixture = (%d, %q, %q), want (0, %q, \"\")", code, stdout, stderr, want)
	}
	if _, err := os.Stat(filepath.Join(home, "responses")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("a 10-line response created the spill store: %v", err)
	}
}

// BO7: the bound keeps the command's own exit code.
func TestDispatcherKeepsExitCode(t *testing.T) {
	stdout, _, _, code := runBoundFixture(t, boundResponse, 30, 3)
	if code != 3 {
		t.Fatalf("bounded fixture exit = %d, want 3", code)
	}
	if lines := strings.Count(stdout, "\n"); lines != 10 {
		t.Fatalf("bounded fixture printed %d lines, want 10: %q", lines, stdout)
	}
}

// A pending entry is not bounded yet, so it prints every line.
func TestDispatcherPassesPendingResponse(t *testing.T) {
	stdout, _, _, _ := runBoundFixture(t, boundPending, 25, 0)
	if want := numberedLines(25); stdout != want {
		t.Fatalf("pending fixture stdout = %q, want %q", stdout, want)
	}
}

// BO31: an exec grammar refusal passes through the bound unchanged, so its first stderr
// line still marks a refusal and not the child's own exit 2.
func TestExecGrammarRefusalKeepsUsageLine(t *testing.T) {
	t.Setenv(benchhome.Env, t.TempDir())
	for _, argv := range [][]string{{"worktree", "exec", "target"}, {"worktree", "exec", "--unknown", "target", "--", "true"}} {
		result := runAXICommandAt(t, newAXIEnvelopeRepo(t), argv)
		if result.code != 2 || result.stdout != "" || !strings.HasPrefix(result.stderr, "usage: bench worktree exec") {
			t.Errorf("%q = (%d, %q, %q), want exit 2 and a stderr usage line", argv, result.code, result.stdout, result.stderr)
		}
	}
}

// boundDeclarations answers each public entry and each family leaf of registry under its
// declared disposition. An entry is named by its command, and a leaf by its command and
// its leaf name.
func boundDeclarations(registry []commandDefinition) map[boundDisposition][]string {
	declared := map[boundDisposition][]string{}
	for _, definition := range registry {
		if definition.Inventory.Visibility != inventoryPublic {
			continue
		}
		declared[definition.Bound] = append(declared[definition.Bound], definition.Name)
		for _, leaf := range definition.Leaves {
			declared[leaf.Bound] = append(declared[leaf.Bound], definition.Name+" "+leaf.Name)
		}
	}
	for disposition := range declared {
		sort.Strings(declared[disposition])
	}
	return declared
}

// BO21: every public registry entry and every leaf declares a bound disposition. The
// planted entry proves the check refuses a public entry with none.
func TestEveryPublicCommandDeclaresBound(t *testing.T) {
	declared := boundDeclarations(commandRegistry)
	if missing := declared[boundUndeclared]; len(missing) != 0 {
		t.Fatalf("public entries without a bound disposition: %q", missing)
	}
	if bounded := declared[boundResponse]; !reflect.DeepEqual(bounded, []string{"worktree exec"}) {
		t.Fatalf("bounded entries = %q, want only the worktree exec leaf", bounded)
	}
	planted := []commandDefinition{
		{Name: "planted", Inventory: publicInventory()},
		{Name: "internal", Inventory: internalInventory},
		{Name: "family", Inventory: publicInventory(), Bound: boundPending, Leaves: []commandLeaf{{Name: "leaf"}}},
	}
	if missing := boundDeclarations(planted)[boundUndeclared]; !reflect.DeepEqual(missing, []string{"family leaf", "planted"}) {
		t.Fatalf("planted undeclared entries = %q, want the public entry and the leaf", missing)
	}
}
