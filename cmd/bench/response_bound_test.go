package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/poolkey"
	"github.com/gibbonmi/bench/internal/responsebound/responseboundtest"
)

// The expectations here are authored apart from the owner: a bounded response over 10
// lines prints exactly 10, and its fifth line is the spill line. The spec fixes both
// counts, and only the owner package may read the line value from the policy registry.

// boundRun is one dispatch through the production dispatcher under a private Bench home.
type boundRun struct {
	stdout, stderr, home string
	code                 int
}

// runRegistry runs argv through the production dispatcher against registry, under a
// private Bench home.
func runRegistry(t *testing.T, registry []commandDefinition, argv ...string) boundRun {
	t.Helper()
	home := t.TempDir()
	t.Setenv(benchhome.Env, home)
	old := commandRegistry
	defer func() { commandRegistry = old }()
	commandRegistry = registry
	var out, errOut bytes.Buffer
	code := Command{Stdout: &out, Stderr: &errOut, Executable: "bench"}.Run(argv)
	return boundRun{stdout: out.String(), stderr: errOut.String(), home: home, code: code}
}

// linesHandler answers a handler that prints count numbered lines on the stream that
// stream selects and exits 0.
func linesHandler(stream func(Command) io.Writer, count int) commandHandler {
	return func(c Command, _ []string) int {
		fmt.Fprint(stream(c), numberedLines(count))
		return 0
	}
}

func numberedLines(count int) string {
	var all strings.Builder
	for i := 1; i <= count; i++ {
		fmt.Fprintf(&all, "line %02d\n", i)
	}
	return all.String()
}

func stdoutOf(c Command) io.Writer { return c.Stdout }

func stderrOf(c Command) io.Writer { return c.Stderr }

// runPlanted registers one public command declared with bound and scope that runs run, and
// runs it through the production dispatcher under a private Bench home. The zero scope
// declares none, so the response carries no tree row.
func runPlanted(t *testing.T, bound boundDisposition, scope treeScope, run commandHandler) boundRun {
	t.Helper()
	return runRegistry(t, []commandDefinition{{
		Name:      "planted",
		Inventory: publicInventory(helpRow{Order: 1, Description: "print numbered lines"}),
		Bound:     bound,
		Scope:     scope,
		Run:       run,
	}}, "planted")
}

// exitingHandler answers a handler that prints lines numbered lines on stdout and exits
// exit.
func exitingHandler(lines, exit int) commandHandler {
	return func(c Command, args []string) int {
		linesHandler(stdoutOf, lines)(c, args)
		return exit
	}
}

// BO1: a bounded public command that prints 25 lines produces exactly 10 stdout lines.
func TestDispatcherBoundsPublicResponse(t *testing.T) {
	run := runPlanted(t, boundResponse, 0, exitingHandler(25, 0))
	lines := strings.Split(strings.TrimSuffix(run.stdout, "\n"), "\n")
	if run.code != 0 || run.stderr != "" || len(lines) != 10 {
		t.Fatalf("bounded fixture = (%d, %q, %q), want exactly 10 stdout lines", run.code, run.stdout, run.stderr)
	}
	if !strings.HasPrefix(lines[4], "spilled{lines=25,") {
		t.Fatalf("fifth line = %q, want the spill line", lines[4])
	}
}

// BO2: a bounded command that prints exactly 10 lines produces its exact bytes and no
// spill file.
func TestDispatcherPassesBoundaryResponse(t *testing.T) {
	run := runPlanted(t, boundResponse, 0, exitingHandler(10, 0))
	if want := numberedLines(10); run.code != 0 || run.stderr != "" || run.stdout != want {
		t.Fatalf("boundary fixture = (%d, %q, %q), want (0, %q, \"\")", run.code, run.stdout, run.stderr, want)
	}
	if entries, err := os.ReadDir(run.home); err != nil || len(entries) != 0 {
		t.Fatalf("a 10-line response wrote the Bench home: %v (%v)", entries, err)
	}
}

// BO7: the bound keeps the command's own exit code.
func TestDispatcherKeepsExitCode(t *testing.T) {
	run := runPlanted(t, boundResponse, 0, exitingHandler(30, 3))
	if run.code != 3 {
		t.Fatalf("bounded fixture exit = %d, want 3", run.code)
	}
	if lines := strings.Count(run.stdout, "\n"); lines != 10 {
		t.Fatalf("bounded fixture printed %d lines, want 10: %q", lines, run.stdout)
	}
}

// An exempt entry prints every line.
func TestDispatcherPassesExemptResponse(t *testing.T) {
	run := runPlanted(t, boundExempt("fixture reason"), 0, exitingHandler(25, 0))
	if want := numberedLines(25); run.stdout != want {
		t.Fatalf("exempt fixture stdout = %q, want %q", run.stdout, want)
	}
	if entries, err := os.ReadDir(run.home); err != nil || len(entries) != 0 {
		t.Fatalf("an exempt response wrote the Bench home: %v (%v)", entries, err)
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

// boundEntry is one public entry or family leaf of a registry and its declared
// disposition. An entry is named by its command, a leaf by its command and its leaf name,
// and a flag-limited disposition adds its flag.
type boundEntry struct {
	name        string
	disposition boundDisposition
}

func boundEntries(registry []commandDefinition) []boundEntry {
	var entries []boundEntry
	named := func(name string, disposition boundDisposition) boundEntry {
		return boundEntry{name: strings.TrimSpace(name + " " + disposition.flag), disposition: disposition}
	}
	for _, definition := range registry {
		if definition.Inventory.Visibility != inventoryPublic {
			continue
		}
		entries = append(entries, named(definition.Name, definition.Bound))
		for _, leaf := range definition.Leaves {
			entries = append(entries, named(definition.Name+" "+leaf.Name, leaf.Bound))
		}
	}
	return entries
}

// undeclaredBound answers the entries of registry that declare neither the bound nor an
// exemption, sorted.
func undeclaredBound(registry []commandDefinition) []string {
	var missing []string
	for _, entry := range boundEntries(registry) {
		if !entry.disposition.bounded && entry.disposition.exempt == "" {
			missing = append(missing, entry.name)
		}
	}
	sort.Strings(missing)
	return missing
}

// BO21: every public registry entry and every leaf declares a bound disposition. The
// planted entries prove the check refuses a public entry and a leaf with none.
func TestEveryPublicCommandDeclaresBound(t *testing.T) {
	if missing := undeclaredBound(commandRegistry); len(missing) != 0 {
		t.Fatalf("public entries without a bound disposition: %q", missing)
	}
	planted := []commandDefinition{
		{Name: "planted", Inventory: publicInventory()},
		{Name: "internal", Inventory: internalInventory},
		{Name: "flagged", Inventory: publicInventory(), Bound: boundDisposition{flag: "--flag"}},
		{Name: "family", Inventory: publicInventory(), Bound: boundResponse, Leaves: []commandLeaf{{Name: "leaf"}}},
	}
	if missing := undeclaredBound(planted); !reflect.DeepEqual(missing, []string{"family leaf", "flagged --flag", "planted"}) {
		t.Fatalf("planted undeclared entries = %q, want the public entries and the leaf", missing)
	}
}

// spillDirOf answers the directory of the spill file that a bounded response names.
func spillDirOf(t *testing.T, stdout string) string {
	t.Helper()
	lines := strings.SplitAfter(stdout, "\n")
	if len(lines) < 5 {
		t.Fatalf("response = %q, want a spill line", stdout)
	}
	return filepath.Dir(responseboundtest.Path(t, lines[4]))
}

// BO72 at the dispatcher: from inside an assignment worktree, an over-bound leaf whose row
// retires spills under the primary scope, and an over-bound leaf that does not retire
// spills under the assignment scope. The spec fixes the store directory names.
func TestDispatcherSpillsRetiringLeafToPrimary(t *testing.T) {
	home := t.TempDir()
	t.Setenv(benchhome.Env, home)
	repo, checkout, id := responseboundtest.AssignmentCheckout(t, home)
	t.Chdir(checkout)
	old := commandRegistry
	t.Cleanup(func() { commandRegistry = old })
	print11 := func(c Command, _ string, args []string) int { return linesHandler(stdoutOf, 11)(c, args) }
	commandRegistry = []commandDefinition{{
		Name: "family", Inventory: publicInventory(), Bound: boundResponse, LeafUsage: func() string { return "usage\n" },
		Leaves: []commandLeaf{
			{Name: "retire", Root: rootNone, Bound: boundResponse, Retires: true, Run: print11},
			{Name: "keep", Root: rootNone, Bound: boundResponse, Run: print11},
		},
	}}
	store := filepath.Join(home, "responses", poolkey.Key(repo))
	for _, row := range []struct{ leaf, scope string }{{"retire", "primary"}, {"keep", id}} {
		var out bytes.Buffer
		if code := (Command{Stdout: &out, Stderr: &out}).Run([]string{"family", row.leaf}); code != 0 {
			t.Fatalf("family %s exit = %d: %q", row.leaf, code, out.String())
		}
		if got, want := spillDirOf(t, out.String()), filepath.Join(store, row.scope); got != want {
			t.Errorf("family %s spilled under %s, want %s", row.leaf, got, want)
		}
	}
}

// The root lookup is lazy: a verb that removes its own tree before its spill opens finds no
// root, so the spill takes the capped store outside any repository and keys no removed tree.
func TestDispatcherSpillAfterTreeRemovalTakesNoRepository(t *testing.T) {
	home := t.TempDir()
	t.Setenv(benchhome.Env, home)
	_, checkout, _ := responseboundtest.AssignmentCheckout(t, home)
	t.Chdir(checkout)
	old := commandRegistry
	t.Cleanup(func() { commandRegistry = old })
	commandRegistry = []commandDefinition{{
		Name: "remove", Inventory: publicInventory(), Bound: boundResponse,
		Run: func(c Command, args []string) int {
			if err := os.RemoveAll(checkout); err != nil {
				t.Fatalf("remove the tree: %v", err)
			}
			return linesHandler(stdoutOf, 11)(c, args)
		},
	}}
	var out bytes.Buffer
	if code := (Command{Stdout: &out, Stderr: &out}).Run([]string{"remove"}); code != 0 {
		t.Fatalf("remove exit = %d: %q", code, out.String())
	}
	if got, want := spillDirOf(t, out.String()), filepath.Join(home, "responses", "none", "primary"); got != want {
		t.Fatalf("spill after the tree removal went under %s, want %s", got, want)
	}
}

// The retiring verbs of the spec each declare Retires on their own leaf row, and no other
// leaf does.
func TestRetiringLeavesDeclareRetires(t *testing.T) {
	var retiring []string
	for _, leaf := range worktreeLeaves {
		if leaf.Retires {
			retiring = append(retiring, leaf.Name)
		}
	}
	sort.Strings(retiring)
	if want := []string{"clean", "land", "reclaim", "release"}; !reflect.DeepEqual(retiring, want) {
		t.Fatalf("retiring worktree leaves = %q, want %q", retiring, want)
	}
}
