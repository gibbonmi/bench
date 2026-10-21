package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/canonicalpath"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// wrapperKitShape is one invocation of bin/bench.sh: the tree whose wrapper runs, the
// working directory, extra environment, and the kit the wrapper must export.
type wrapperKitShape struct {
	name string
	tree string
	cwd  string
	env  []string
	want string
}

// wrapperKitFixture holds the repository shapes of the wrapper's kit rule. Each wrapper
// tree runs a stub dist/bench that writes the BENCH_KIT it receives to argvFile.
type wrapperKitFixture struct {
	kit, worktree string
	shapes        []wrapperKitShape
	argvFile      string
	env           []string
}

// newWrapperKitFixture builds the wrapper's repository, a linked worktree of it, a
// different repository, an adopted repository with its kit at <repo>/.bench, and a plain
// directory. Its shapes set no environment override.
func newWrapperKitFixture(t *testing.T) *wrapperKitFixture {
	t.Helper()
	base := t.TempDir()
	kit := filepath.Join(base, "kit")
	argvFile := filepath.Join(base, "argv")
	copyExecutable(t, filepath.Join("..", "..", "bin", "bench.sh"), filepath.Join(kit, "bin", "bench.sh"))
	gitInitRepo(t, kit)
	// dist/bench stays untracked, exactly as in a real checkout. The worktree case reaches
	// it only because main_tree_kit re-anchors the binary lookup at the main tree.
	writeExecutable(t, filepath.Join(kit, "dist", "bench"), `#!/usr/bin/env bash
printf '%s\n' "$BENCH_KIT" > "$BENCH_TEST_ARGV"
`)
	worktree := filepath.Join(base, "wt")
	add := exec.Command("git", "worktree", "add", "-q", "-b", "probe", worktree)
	add.Dir = kit
	if out, err := add.CombinedOutput(); err != nil {
		t.Fatalf("git worktree add: %v\n%s", err, out)
	}
	project := filepath.Join(base, "project")
	gitInitRepo(t, project)
	// In the adopted-repo layout, the kit is a subdirectory of the project repo. The
	// project's own CWD shares its common dir with the wrapper's tree, so only the
	// "kit is its tree's top level" condition keeps the wrapper naming <repo>/.bench.
	adopted := filepath.Join(base, "adopted")
	adoptedKit := filepath.Join(adopted, ".bench")
	copyExecutable(t, filepath.Join("..", "..", "bin", "bench.sh"), filepath.Join(adoptedKit, "bin", "bench.sh"))
	writeExecutable(t, filepath.Join(adoptedKit, "dist", "bench"), `#!/usr/bin/env bash
printf '%s\n' "$BENCH_KIT" > "$BENCH_TEST_ARGV"
`)
	gitInitRepo(t, adopted)
	plain := filepath.Join(base, "plain")
	if err := os.MkdirAll(plain, 0o755); err != nil {
		t.Fatal(err)
	}

	env := capability.WithoutEnvironment(os.Environ(), runbinary.Env)
	env = capability.WithoutEnvironment(capability.WithoutEnvironment(env, "BENCH_KIT"), "BENCH_WRAPPER")
	env = append(env, "BENCH_TEST_ARGV="+argvFile, "BENCH_HOME="+filepath.Join(base, "home"))
	return &wrapperKitFixture{
		kit:      kit,
		worktree: worktree,
		argvFile: argvFile,
		env:      env,
		shapes: []wrapperKitShape{
			{name: "worktree of the wrapper's repository", tree: kit, cwd: worktree, want: worktree},
			{name: "an adopted repo's own kit subdirectory", tree: adoptedKit, cwd: adopted, want: adoptedKit},
			{name: "the wrapper's own main checkout", tree: kit, cwd: kit, want: kit},
			{name: "a different repository (linked project)", tree: kit, cwd: project, want: kit},
			{name: "outside any repository", tree: kit, cwd: plain, want: kit},
		},
	}
}

// run invokes the shape's wrapper from the shape's working directory and answers the kit
// the wrapper exports.
func (f *wrapperKitFixture) run(t *testing.T, shape wrapperKitShape) string {
	t.Helper()
	if err := os.Remove(f.argvFile); err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	cmd := exec.Command("bash", filepath.Join(shape.tree, "bin", "bench.sh"), "gate-phases", "/tmp/root")
	cmd.Dir = shape.cwd
	cmd.Env = append(append([]string(nil), f.env...), shape.env...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("bench.sh from %s: %v\n%s", shape.cwd, err, out)
	}
	return strings.TrimSpace(readPath(t, f.argvFile))
}

// TestShellWrapperKitMatchesKitForRoot holds the Bash kit rule and gate.KitForRoot to one
// answer: the wrapper's export is the oracle for each shape. The wrapper exports physical
// paths, so the two answers compare under their canonical spellings.
func TestShellWrapperKitMatchesKitForRoot(t *testing.T) {
	f := newWrapperKitFixture(t)
	for _, shape := range f.shapes {
		t.Run(shape.name, func(t *testing.T) {
			wrapperKit := f.run(t, shape)
			goKit := gate.KitForRoot(shape.tree, shape.cwd)
			if canonicalSpelling(t, goKit) != canonicalSpelling(t, wrapperKit) {
				t.Fatalf("gate.KitForRoot(%q, %q) = %q, want the wrapper's kit %q", shape.tree, shape.cwd, goKit, wrapperKit)
			}
		})
	}
}

func canonicalSpelling(t *testing.T, path string) string {
	t.Helper()
	resolved, err := canonicalpath.Resolve(path)
	if err != nil {
		t.Fatal(err)
	}
	return resolved
}
