package treetarget

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/worktree"
)

// The usage lines, the refusal pairs, and the child facts here are authored apart from Run:
// the spec fixes each line, the child argv, and the child directory. A changed word, a
// started child, or a child that keeps the recorded spelling reds these rows. Each test
// passes its executables and its home as arguments, so no test binds the environment.

// targetRepo is a primary checkout and the Bench home that its assignments live under.
type targetRepo struct{ root, home string }

func newTargetRepo(t *testing.T, home string) targetRepo {
	t.Helper()
	return targetRepo{root: committedRepo(t), home: home}
}

// create makes one active assignment labeled label through the worktree create verb, and
// answers its ledger record.
func (r targetRepo) create(t *testing.T, request, label string) intent.Assignment {
	t.Helper()
	var stdout, stderr bytes.Buffer
	if code := worktree.CreateCommand(r.root, r.home, []string{"--request", request, "--label", label}, &stdout, &stderr); code != 0 {
		t.Fatalf("worktree create %q = (%d, %q, %q)", label, code, stdout.String(), stderr.String())
	}
	assignment, ok, err := intent.FindAssignmentByRequest(r.root, intent.RequestDigest(request))
	if err != nil || !ok {
		t.Fatalf("the ledger holds no assignment for the request %q: %v", request, err)
	}
	return assignment
}

// run runs the verb name with the flag arguments args through Run from the primary
// checkout. wrapper and running are the two executables that the child can run.
func (r targetRepo) run(wrapper, running, name string, args ...string) (stdout, stderr string, code int) {
	var out, errOut bytes.Buffer
	code = Run(Call{Name: name, Flag: "--in", Args: args, Root: r.root, Home: r.home, Wrapper: wrapper, Running: running,
		Stdin: strings.NewReader(""), Stdout: &out, Stderr: &errOut})
	return out.String(), errOut.String(), code
}

// writeScript writes text as an executable file at path, with its parent directories.
func writeScript(t *testing.T, path, text string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(text), 0o755); err != nil {
		t.Fatal(err)
	}
}

// markerScript writes an executable script that records its argv, its physical directory,
// and its environment in a marker file, and then runs body. It answers the script and the
// marker.
func markerScript(t *testing.T, body string) (script, marker string) {
	t.Helper()
	dir := t.TempDir()
	script, marker = filepath.Join(dir, "child"), filepath.Join(dir, "marker")
	writeScript(t, script, "#!/bin/sh\n{ printf 'argv'; printf ' %s' \"$@\"; printf '\\ndir %s\\n' \"$(pwd -P)\"; env; } > '"+marker+"'\n"+body+"\n")
	return script, marker
}

// childRecord is what one marker script recorded: its argv line, its physical directory,
// and its environment.
type childRecord struct {
	argv, dir string
	env       map[string]string
}

func readMarker(t *testing.T, marker string) childRecord {
	t.Helper()
	data, err := os.ReadFile(marker)
	if err != nil {
		t.Fatalf("no child started: %v", err)
	}
	lines := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	record := childRecord{argv: lines[0], dir: strings.TrimPrefix(lines[1], "dir "), env: map[string]string{}}
	for _, line := range lines[2:] {
		if key, value, ok := strings.Cut(line, "="); ok {
			record.env[key] = value
		}
	}
	return record
}

func requireNoChild(t *testing.T, marker string) {
	t.Helper()
	if _, err := os.Stat(marker); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("a child started: %v", err)
	}
}

// physicalPath answers the spelling of path with every link resolved.
func physicalPath(t *testing.T, path string) string {
	t.Helper()
	physical, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatal(err)
	}
	return physical
}

// TT27, TT53: Run starts the given executable with the verb and its remaining arguments,
// in the physical root of the target. The Bench home is a link, and the child's PWD holds
// the physical spelling of the recorded worktree path.
func TestRunStartsChildInTarget(t *testing.T) {
	home := filepath.Join(t.TempDir(), "home-link")
	if err := os.Symlink(t.TempDir(), home); err != nil {
		t.Fatal(err)
	}
	r := newTargetRepo(t, home)
	alpha := r.create(t, "tree-target-alpha", "alpha")
	want := physicalPath(t, alpha.Worktree)
	if strings.HasPrefix(want, home) {
		t.Fatalf("the physical worktree path %q carries the home link %q", want, home)
	}
	script, marker := markerScript(t, "exit 0")
	if stdout, stderr, code := r.run("", script, "status", "alpha", "--route"); code != 0 || stdout != "" || stderr != "" {
		t.Fatalf("status --in alpha --route = (%d, %q, %q), want the child's exit 0 and no parent output", code, stdout, stderr)
	}
	child := readMarker(t, marker)
	if child.argv != "argv status --route" || child.dir != want {
		t.Fatalf("child = %q in %q, want argv status --route in %q", child.argv, child.dir, want)
	}
	if child.env["PWD"] != want {
		t.Fatalf("child PWD = %q, want the physical worktree path %q", child.env["PWD"], want)
	}
}

// TT29, TT30, TT31, TT33, TT34, TT35, TT36, TT57, TT58, TT59: a value that names no active
// assignment by its exact label refuses before any child starts. A missing, empty, dash, or
// path value is a grammar refusal at exit 2 on stdout, and its usage line escapes a control
// character. Every other refusal is the worktree target refusal at exit 1 on stderr. An
// active label whose tree is gone, or whose tree no longer holds the assignment branch,
// fails the missing-tree or the creation-bundle check of the lookup. Each wanted line is
// exact and holds no raw control byte.
func TestRunRefusesBeforeChild(t *testing.T) {
	r := newTargetRepo(t, t.TempDir())
	alpha := r.create(t, "tree-target-alpha", "alpha")
	retired := r.create(t, "tree-target-retired", "retired")
	retired.State = intent.StateComplete
	if err := intent.PutAssignment(r.root, retired); err != nil {
		t.Fatal(err)
	}
	twinA := r.create(t, "tree-target-twin-a", "twin")
	twinB := r.create(t, "tree-target-twin-b", "twin")
	// The gone branch is still the tip of main, so it has landed, and the refusal routes to
	// the batch clean.
	gone := r.create(t, "tree-target-gone", "gone")
	if err := os.RemoveAll(gone.Worktree); err != nil {
		t.Fatal(err)
	}
	detached := r.create(t, "tree-target-detached", "detached")
	gitIn(t, detached.Worktree, "checkout", "-q", "--detach")
	missing := "usage: bench gate --in (missing argument: <label|primary>)\n"
	usage := func(value string) string { return "usage: bench gate --in (unknown argument: " + value + ")\n" }
	refusal := func(reason string) string { return "bench gate --in: " + reason + "\nnext=bench worktree list\n" }
	// The ambiguity refusal names both ids in the ledger's order, which the spec leaves open.
	ambiguous := func(first, second string) string { return refusal("target is ambiguous: " + first + ", " + second) }
	for _, row := range []struct {
		name                string
		args                []string
		code                int
		stdout, stderr, alt string
	}{
		{name: "TT29 missing value", code: 2, stdout: missing},
		{name: "TT57 empty value", args: []string{""}, code: 2, stdout: missing},
		{name: "TT58 dash value", args: []string{"--help"}, code: 2, stdout: usage("--help")},
		{name: "dash value with a control character", args: []string{"-\x01"}, code: 2, stdout: usage(`-\u0001`)},
		{name: "TT30 absolute pool path", args: []string{alpha.Worktree}, code: 2, stdout: usage(alpha.Worktree)},
		{name: "TT31 relative path", args: []string{"./alpha"}, code: 2, stdout: usage("./alpha")},
		{name: "TT59 home form", args: []string{"~"}, code: 2, stdout: usage("~")},
		{name: "home form of another user", args: []string{"~nobody"}, code: 2, stdout: usage("~nobody")},
		{name: "TT33 unknown label", args: []string{"beta"}, code: 1, stderr: refusal("target is unassigned")},
		{name: "TT34 released label", args: []string{"retired"}, code: 1, stderr: refusal("assignment " + retired.ID + " is not active")},
		{name: "TT35 colliding label", args: []string{"twin"}, code: 1, stderr: ambiguous(twinA.ID, twinB.ID), alt: ambiguous(twinB.ID, twinA.ID)},
		{name: "TT36 control character", args: []string{"al\x01pha"}, code: 1, stderr: refusal("target contains control characters")},
		{name: "missing tree", args: []string{"gone"}, code: 1, stderr: "bench gate --in: worktree tree is missing\nnext=bench worktree clean --landed\n"},
		{name: "creation bundle", args: []string{"detached"}, code: 1, stderr: refusal("assignment branch is not checked out")},
	} {
		t.Run(row.name, func(t *testing.T) {
			script, marker := markerScript(t, "exit 0")
			stdout, stderr, code := r.run("", script, "gate", row.args...)
			if code != row.code || stdout != row.stdout || (stderr != row.stderr && (row.alt == "" || stderr != row.alt)) {
				t.Fatalf("gate --in %q = (%d, %q, %q), want (%d, %q, %q)", row.args, code, stdout, stderr, row.code, row.stdout, row.stderr)
			}
			requireNoChild(t, marker)
		})
	}
}

// TT32: the label lookup runs before the path-shape step, so a label with `/` still names
// its worktree.
func TestRunLabelWinsOverPathShape(t *testing.T) {
	r := newTargetRepo(t, t.TempDir())
	team := r.create(t, "tree-target-team", "team/alpha")
	script, marker := markerScript(t, "exit 0")
	if stdout, stderr, code := r.run("", script, "gate", "team/alpha"); code != 0 {
		t.Fatalf("gate --in team/alpha = (%d, %q, %q), want the child's exit 0", code, stdout, stderr)
	}
	if child, want := readMarker(t, marker), physicalPath(t, team.Worktree); child.dir != want {
		t.Fatalf("child ran in %q, want the team/alpha worktree %q", child.dir, want)
	}
}

// TT38: the parent passes the child's output and exit through, and it prints no worktree
// line for a failed child.
func TestRunPassesChildResult(t *testing.T) {
	r := newTargetRepo(t, t.TempDir())
	r.create(t, "tree-target-alpha", "alpha")
	script, _ := markerScript(t, "printf x; exit 3")
	if stdout, stderr, code := r.run("", script, "gate", "alpha"); code != 3 || stdout != "x" || stderr != "" {
		t.Fatalf("gate --in alpha = (%d, %q, %q), want (3, \"x\", \"\") with no worktree line", code, stdout, stderr)
	}
}

// TT39, TT40, TT49: a target without build inputs runs the invoking wrapper, or else the
// running executable. A dist/bench file in the target never runs.
func TestRunSelectsChildExecutable(t *testing.T) {
	r := newTargetRepo(t, t.TempDir())
	alpha := r.create(t, "tree-target-alpha", "alpha")
	if freshness.DeclaresBuildInputs(alpha.Worktree) {
		t.Fatal("the alpha fixture declares build inputs, so it is not a linked-repository target")
	}
	strayMarker := filepath.Join(t.TempDir(), "stray")
	writeScript(t, filepath.Join(alpha.Worktree, "dist", "bench"), "#!/bin/sh\n: > '"+strayMarker+"'\n")
	t.Run("TT39 TT49 the wrapper", func(t *testing.T) {
		wrapper, wrapperMarker := markerScript(t, "exit 0")
		running, runningMarker := markerScript(t, "exit 0")
		stdout, stderr, code := r.run(wrapper, running, "gate", "alpha")
		requireNoChild(t, strayMarker)
		if code != 0 {
			t.Fatalf("gate --in alpha = (%d, %q, %q), want the wrapper's exit 0", code, stdout, stderr)
		}
		if child := readMarker(t, wrapperMarker); child.argv != "argv gate" {
			t.Fatalf("wrapper argv = %q, want argv gate", child.argv)
		}
		requireNoChild(t, runningMarker)
	})
	t.Run("TT40 the running executable", func(t *testing.T) {
		running, runningMarker := markerScript(t, "exit 0")
		stdout, stderr, code := r.run("", running, "gate", "alpha")
		requireNoChild(t, strayMarker)
		if code != 0 {
			t.Fatalf("gate --in alpha = (%d, %q, %q), want the running executable's exit 0", code, stdout, stderr)
		}
		if child := readMarker(t, runningMarker); child.argv != "argv gate" {
			t.Fatalf("running executable argv = %q, want argv gate", child.argv)
		}
	})
}
