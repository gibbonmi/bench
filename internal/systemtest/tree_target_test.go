//go:build system

package systemtest

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/treetarget/kittest"
	"github.com/gibbonmi/bench/internal/treetarget/treetargettest"
	"github.com/gibbonmi/bench/internal/worktree"
)

// The row cells here are authored apart from the renderer: the spec fixes one row of the
// target, the full HEAD commit, and the dirty cell. The header and the TOON quoting of a
// cell come from their owners, so this file restates neither.

// onlyInAlphaSpec is a repo-relative spec path that only the alpha branch holds.
const onlyInAlphaSpec = "specs/only-in-alpha/spec.md"

// treeTargetFixture answers a scaffolded repository and one active assignment labeled alpha.
// The alpha branch holds one commit that adds onlyInAlphaSpec, so the two trees differ in
// both their HEAD and their content.
func treeTargetFixture(t *testing.T) (recordedPublicationRepo, systemLandingWorktree) {
	t.Helper()
	fixture := scaffoldRecordedPublicationRepo(t)
	alpha := systemCreateLandingWorktree(t, fixture.root, fixture.home, "alpha", "tree target")
	if err := os.MkdirAll(filepath.Join(alpha.path, filepath.Dir(onlyInAlphaSpec)), 0o755); err != nil {
		t.Fatal(err)
	}
	systemCommit(t, alpha.path, onlyInAlphaSpec, "# Only in alpha\n", "spec only in alpha")
	return fixture, alpha
}

// runTreeTarget runs the selected executable in dir. The environment removes the invoking
// wrapper, so a tree-target child is the running executable. It also removes the command
// observation, so stderr holds only the lines of the verb.
func runTreeTarget(t *testing.T, fixture recordedPublicationRepo, dir string, args ...string) processResult {
	t.Helper()
	unobserved := capability.WithoutEnvironment(fixture.environment(fixture.home), "BENCH_COMMAND_OBSERVE")
	return systemSelected(t, dir, append(unobserved, worktree.WrapperEnv), args...)
}

// treeRowCells answers the cells of the identity row that leads stdout, unquoted and joined
// by commas, or the empty string when stdout does not start with the row block.
func treeRowCells(stdout string) string {
	block := strings.TrimSuffix(stdout, treetargettest.WithoutRow(stdout))
	_, row, _ := strings.Cut(strings.TrimSuffix(block, "\n"), "\n")
	cells := strings.Split(strings.TrimSpace(row), ",")
	for i, cell := range cells {
		if value, err := systemTOONCell(cell); err == nil {
			cells[i] = value
		}
	}
	return strings.Join(cells, ",")
}

// TT26, TT28, TT41, TT54, TT55: a tree-scoped verb names its tree in its first block. A tree
// target runs the verb in the named tree, from the primary checkout and from a worktree,
// and a relative operand resolves in that tree. The exec route and directory inference
// still name the tree that they run in.
func TestTreeTargetRunsInNamedWorktree(t *testing.T) {
	fixture, alpha := treeTargetFixture(t)
	primaryRow := "primary," + systemGitOutput(t, fixture.root, "rev-parse", "HEAD") + ",false"
	alphaRow := "alpha," + systemGitOutput(t, alpha.path, "rev-parse", "HEAD") + ",false"
	for _, row := range []struct {
		name, dir, row, body string
		args                 []string
	}{
		{name: "TT26 a label from the primary checkout", dir: fixture.root, row: alphaRow, args: []string{"status", "--in", "alpha"}},
		{name: "TT28 the primary keyword from a worktree", dir: alpha.path, row: primaryRow, args: []string{"status", "--in", "primary"}},
		{name: "TT41 a relative operand in the target", dir: fixture.root, row: alphaRow, body: "\nspec: " + onlyInAlphaSpec + "\n", args: []string{"coverage", "--in", "alpha", onlyInAlphaSpec}},
		{name: "TT54 the exec route", dir: fixture.root, row: alphaRow, args: []string{"worktree", "exec", "alpha", "--", owner.selected.path, "status"}},
		{name: "TT55 directory inference", dir: alpha.path, row: alphaRow, args: []string{"status"}},
	} {
		t.Run(row.name, func(t *testing.T) {
			result := runTreeTarget(t, fixture, row.dir, row.args...)
			if result.code != 0 || treeRowCells(result.stdout) != row.row || !strings.Contains(result.stdout, row.body) {
				t.Fatalf("bench %q in %s = (%d, %q, %q), want exit 0, the row %q first, and %q", row.args, row.dir, result.code, result.stdout, result.stderr, row.row, row.body)
			}
		})
	}
}

// TT50: the real executable, with BENCH_KIT set to the kit root, refuses a kit worktree
// target whose build no longer matches the tree, and it starts no child. The build is a
// script that freshness.Publish seals and that writes a marker when it runs.
func TestTreeTargetRefusesStaleKitBuild(t *testing.T) {
	fixture, alpha := treeTargetFixture(t)
	kittest.WriteTree(t, alpha.path)
	marker := filepath.Join(t.TempDir(), "marker")
	script := filepath.Join(t.TempDir(), "bench")
	if err := os.WriteFile(script, []byte("#!/bin/sh\n: > '"+marker+"'\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	build := freshness.PublishedExecutable(alpha.path)
	if err := freshness.Publish(alpha.path, script, build, filepath.Dir(build), "tree-target"); err != nil {
		t.Fatal(err)
	}
	kittest.EditBuildInput(t, alpha.path)
	result := runTreeTarget(t, fixture, fixture.root, "status", "--in", "alpha")
	want := "bench status --in: worktree build does not match the tree\nnext=bench worktree build alpha\n"
	if result.code != 1 || result.stdout != "" || result.stderr != want {
		t.Fatalf("bench status --in alpha = (%d, %q, %q), want (1, \"\", %q)", result.code, result.stdout, result.stderr, want)
	}
	if _, err := os.Stat(marker); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("a child started: %v", err)
	}
}
