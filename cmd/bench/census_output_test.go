package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/census"
	"github.com/gibbonmi/bench/internal/worktree"
)

// useCensusFamily registers one bounded family whose `print` leaf prints lines numbered
// lines and exits exit, and whose `retire` leaf does the same and retires an assignment.
func useCensusFamily(t *testing.T, lines, exit int) {
	t.Helper()
	old := commandRegistry
	t.Cleanup(func() { commandRegistry = old })
	printLines := func(c Command, _ string, args []string) int {
		linesHandler(stdoutOf, lines)(c, args)
		return exit
	}
	commandRegistry = []commandDefinition{{
		Name: "family", Inventory: publicInventory(), Bound: boundResponse, LeafUsage: func() string { return "usage\n" },
		Leaves: []commandLeaf{
			{Name: "print", Root: rootNone, Bound: boundResponse, Run: printLines},
			{Name: "retire", Root: rootNone, Bound: boundResponse, Retires: true, Run: printLines},
		},
	}}
}

// dispatch runs argv through the production dispatcher and answers its response.
func dispatch(argv ...string) (string, int) {
	var out bytes.Buffer
	code := Command{Stdout: &out, Stderr: &out}.Run(argv)
	return out.String(), code
}

// outputRecordFields answers the fields after the time of each output record line that
// an assignment holds, or none when the file is absent.
func outputRecordFields(t *testing.T, home, root, id string) [][]string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(census.Dir(home, root), id+".output"))
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	var records [][]string
	for _, line := range strings.Split(strings.TrimSuffix(string(data), "\n"), "\n") {
		records = append(records, strings.Split(line, "\t")[1:])
	}
	return records
}

// BO57 at the dispatcher: a bounded leaf run inside an assignment worktree appends one
// record under that assignment with the verb head, the complete output's size, and the
// disposition. A retiring leaf appends none, because its retirement drops the file.
func TestDispatcherRecordsOutputInAssignment(t *testing.T) {
	home := t.TempDir()
	t.Setenv(benchhome.Env, home)
	repo, checkout, id := assignmentCheckout(t, home)
	t.Chdir(checkout)
	useCensusFamily(t, 11, 0)
	for _, leaf := range []string{"print", "retire"} {
		if out, code := dispatch("family", leaf); code != 0 {
			t.Fatalf("family %s = (%d, %q)", leaf, code, out)
		}
	}
	want := [][]string{{"bench family print", "11", fmt.Sprint(len(numberedLines(11))), "spilled"}}
	if got := outputRecordFields(t, home, repo, id); fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("output records = %q, want %q", got, want)
	}
}

// BO58: `bench worktree exec <target>` from the primary checkout records its response
// under the target assignment, with the exec head.
func TestExecOutputRecordUsesTarget(t *testing.T) {
	home := t.TempDir()
	t.Setenv(benchhome.Env, home)
	repo := newAXIEnvelopeRepo(t)
	creation, err := worktree.Create(repo, "census-output-exec", "census-output-exec", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Chdir(repo)
	if out, code := dispatch("worktree", "exec", creation.Assignment.Label, "--", "echo", "child"); code != 0 {
		t.Fatalf("exec = (%d, %q)", code, out)
	}
	want := [][]string{{"bench worktree exec", "1", fmt.Sprint(len("child\n")), "inline"}}
	if got := outputRecordFields(t, home, repo, creation.Assignment.ID); fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("output records under the target = %q, want %q", got, want)
	}
}

// BO62: a symlink at the census directory refuses the record, and the verb's response and
// exit code stay the command's own. The same run without the symlink records, so the
// refusal is the record path failing and not a path the dispatcher never reached.
func TestOutputRecordFailureKeepsVerdict(t *testing.T) {
	home := t.TempDir()
	t.Setenv(benchhome.Env, home)
	repo, checkout, id := assignmentCheckout(t, home)
	t.Chdir(checkout)
	useCensusFamily(t, 3, 3)
	elsewhere := t.TempDir()
	dir := census.Dir(home, repo)
	if err := os.MkdirAll(filepath.Dir(dir), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(elsewhere, dir); err != nil {
		t.Fatal(err)
	}
	if out, code := dispatch("family", "print"); code != 3 || out != numberedLines(3) {
		t.Fatalf("family print with a symlinked census = (%d, %q), want (3, %q)", code, out, numberedLines(3))
	}
	if entries, err := os.ReadDir(elsewhere); err != nil || len(entries) != 0 {
		t.Fatalf("the record followed the symlink: %v (%v)", entries, err)
	}
	if err := os.Remove(dir); err != nil {
		t.Fatal(err)
	}
	if out, code := dispatch("family", "print"); code != 3 || len(outputRecordFields(t, home, repo, id)) != 1 {
		t.Fatalf("family print without the symlink = (%d, %q), want exit 3 and one record", code, out)
	}
}
