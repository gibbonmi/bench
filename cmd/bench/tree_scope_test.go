package main

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/env"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/treetarget/treetargettest"
)

// TT6 and TT7: each repository-scoped verb refuses a tree target as its first argument,
// before the verb runs. The verbs come from the registry, so no test restates the scope
// table. The flag spelling is authored apart from the dispatcher, so a refusal that reads
// another argument fails. A wrapper-only definition never reaches the dispatcher.
func TestRepositoryVerbRefusesTreeTarget(t *testing.T) {
	t.Setenv(benchhome.Env, t.TempDir())
	// A verb that runs past a missing refusal must fail fast, so `bench shift` finds no
	// harness adapter and starts no loop.
	t.Setenv("BENCH_AGENT", "")
	root := newAXIEnvelopeRepo(t)
	// An ignored inbox lets `bench idea` park its line in this primary checkout, so a verb
	// that runs before the refusal leaves the file behind.
	writeAXIFixture(t, filepath.Join(root, ".gitignore"), "capture/\n")

	refused := 0
	for _, definition := range commandRegistry {
		if definition.Scope != scopeRepository || definition.WrapperOnly {
			continue
		}
		refused++
		want := toon.Usage("bench "+definition.Name, "--in") + "\n"
		for _, args := range [][]string{{"--in", "primary"}, {"--in"}} {
			argv := append([]string{definition.Name}, args...)
			t.Run(strings.Join(argv, " "), func(t *testing.T) {
				result := runAXICommandAt(t, root, argv)
				if result.code != 2 || result.stdout != want || result.stderr != "" {
					t.Fatalf("bench %q = (%d, %q, %q), want (2, %q, \"\")", argv, result.code, result.stdout, result.stderr, want)
				}
			})
		}
	}
	if refused == 0 {
		t.Fatal("the command registry holds no repository-scoped verb")
	}
	t.Run("idea", func(t *testing.T) {
		result := runAXICommandAt(t, root, []string{"idea", "--in", "primary", "x"})
		if result.code != 2 {
			t.Fatalf("idea --in primary x = (%d, %q, %q), want exit 2", result.code, result.stdout, result.stderr)
		}
		if _, err := os.Stat(filepath.Join(root, "capture", "IDEAS.md")); !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("idea --in primary x left capture/IDEAS.md: %v", err)
		}
	})
}

// The row expectations below are authored apart from the renderer: the spec fixes the block
// `tree[1]{target,head,dirty}:`, its one row, and its place as the first stdout block. A
// renamed field, a moved row, or a quoted cell reds these rows.

const treeRowHeader = "tree[1]{target,head,dirty}:\n"

// withoutTreeRow answers a dispatcher response without its leading identity row block. The
// dispatch helpers call it, so a fixture that grades the verb output keeps its expectation.
func withoutTreeRow(out string) string { return treetargettest.WithoutRow(out) }

// treeRowRepo makes the current directory a new committed primary checkout. It answers the
// row block that a call in that clean checkout prints. The commit dates are fixed, so HEAD
// is one known commit: TOON quotes a commit that starts with a zero and a digit, and this
// commit does not.
func treeRowRepo(t *testing.T) string {
	t.Helper()
	root := gittest.RepoOnBranch(t, "main")
	t.Setenv("GIT_AUTHOR_DATE", "2026-01-01T00:00:00Z")
	t.Setenv("GIT_COMMITTER_DATE", "2026-01-01T00:00:00Z")
	runAXIGit(t, "-C", root, "commit", "--allow-empty", "-qm", "base")
	t.Chdir(root)
	return treeRowHeader + "  primary," + strings.TrimSpace(runAXIGit(t, "-C", root, "rev-parse", "HEAD")) + ",false\n"
}

// runTreeCall runs argv through the production dispatcher and answers both streams.
func runTreeCall(argv ...string) (stdout, stderr string, code int) {
	var out, errOut bytes.Buffer
	code = Command{Stdout: &out, Stderr: &errOut, Executable: "bench"}.Run(argv)
	return out.String(), errOut.String(), code
}

// TT10, TT11, TT15: in a primary checkout, the first stdout block of `bench roadmap` is the
// row block, and its row names the primary target, the full HEAD commit, and a clean tree.
func TestTreeRowLeadsTreeResponse(t *testing.T) {
	want := treeRowRepo(t)
	stdout, stderr, code := runTreeCall("roadmap")
	if code == 2 || !strings.HasPrefix(stdout, want) {
		t.Fatalf("roadmap = (%d, %q, %q), want the row block %q first", code, stdout, stderr, want)
	}
	if _, cells, _ := strings.Cut(want, "  primary,"); len(strings.Split(cells, ",")[0]) != 40 {
		t.Fatalf("row %q does not carry the 40-character HEAD commit", want)
	}
}

// TT37: the tree-target flag names a tree target only as the first argument after the verb.
// A late flag reaches the verb's own grammar: the gate refuses it at exit 2, and `bench
// version` ignores its arguments. The invoking wrapper is a marker script, so a started
// child leaves the marker and never runs this test binary. The first-position call starts
// that child, so the absent marker of a late flag is not silent by accident.
func TestTreeTargetOnlyAsFirstArgument(t *testing.T) {
	treeRowRepo(t)
	marker := filepath.Join(t.TempDir(), "child-started")
	wrapper := filepath.Join(t.TempDir(), "wrapper")
	writeExecutable(t, wrapper, "#!/bin/sh\n: > '"+marker+"'\n")
	t.Setenv(env.WrapperEnv, wrapper)
	if stdout, stderr, code := runTreeCall("gate", "--fresh", "--in", "primary"); code != 2 || stdout != "" || stderr != gate.CommandUsage+"\n" {
		t.Fatalf("gate --fresh --in primary = (%d, %q, %q), want (2, \"\", %q)", code, stdout, stderr, gate.CommandUsage+"\n")
	}
	version, _, _ := runTreeCall("version")
	if stdout, stderr, code := runTreeCall("version", "x", "--in", "primary"); code != 0 || stdout != version || stderr != "" {
		t.Fatalf("version x --in primary = (%d, %q, %q), want (0, %q, \"\")", code, stdout, stderr, version)
	}
	if _, err := os.Stat(marker); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("a late tree-target flag started a child: %v", err)
	}
	if stdout, stderr, code := runTreeCall("status", "--in", "primary"); code != 0 || stdout != "" {
		t.Fatalf("status --in primary = (%d, %q, %q), want the child's exit 0 and no parent row", code, stdout, stderr)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("a first-position tree-target flag started no child: %v", err)
	}
}

// TT20: an exempt call keeps its stdout artifact byte-clean and prints the row on stderr.
func TestExemptTreeCallPrintsRowOnStderr(t *testing.T) {
	want := treeRowRepo(t)
	stdout, stderr, code := runTreeCall("dashboard", "--stdout")
	if code == 2 || strings.Contains(stdout, "tree[") || !strings.Contains(stderr, want) {
		t.Fatalf("dashboard --stdout = (%d, stdout %q, stderr %q), want no row on stdout and %q on stderr", code, stdout, stderr, want)
	}
}

// TT21, TT22: a help form of a tree-scoped verb and a repository-scoped verb print no row.
func TestHelpFormPrintsNoTreeRow(t *testing.T) {
	treeRowRepo(t)
	for _, argv := range [][]string{{"gate", "--help"}, {"version"}} {
		if stdout, stderr, code := runTreeCall(argv...); code != 0 || strings.Contains(stdout+stderr, "tree[") {
			t.Errorf("%q = (%d, %q, %q), want exit 0 and no row", argv, code, stdout, stderr)
		}
	}
}

// TT24: outside a repository a tree-scoped verb prints no row and keeps its own answer,
// which names the spec path that it cannot find.
func TestTreeRowOutsideRepository(t *testing.T) {
	t.Chdir(t.TempDir())
	stdout, stderr, code := runTreeCall("coverage", "x")
	if code != 1 || stderr != "" || !strings.HasPrefix(stdout, "error: spec not found: x,") || strings.Contains(stdout, "tree[") {
		t.Fatalf("coverage x outside a repository = (%d, %q, %q), want the not-in-repo answer and no row", code, stdout, stderr)
	}
}

// plantedTreeVerb runs one planted tree-scoped verb with bound and run in a new primary
// checkout. It answers the row block of the clean checkout and the run.
func plantedTreeVerb(t *testing.T, bound boundDisposition, run commandHandler) (string, boundRun) {
	t.Helper()
	row := treeRowRepo(t)
	return row, runPlanted(t, bound, scopeTree, run)
}

// TT25: a spilled tree-scoped response keeps the row header as its first inline line.
func TestTreeRowSurvivesSpill(t *testing.T) {
	_, run := plantedTreeVerb(t, boundResponse, linesHandler(stdoutOf, 40))
	if first, _, _ := strings.Cut(run.stdout, "\n"); first+"\n" != treeRowHeader || !strings.Contains(run.stdout, "\nspilled{lines=40,") {
		t.Fatalf("spilled tree response = %q, want the row header first and a spill line", run.stdout)
	}
}

// TT60: the two row lines do not count toward the bound, so 10 verb lines do not spill.
func TestTreeRowOutsideResponseBound(t *testing.T) {
	row, run := plantedTreeVerb(t, boundResponse, linesHandler(stdoutOf, 10))
	if want := row + numberedLines(10); run.code != 0 || run.stdout != want {
		t.Fatalf("10-line tree response = (%d, %q), want (0, %q)", run.code, run.stdout, want)
	}
}

// TT20, TT56: an exempt call prints the row on stderr only when its exit is not 2. The exit
// 1 row shows that the checkout prints a row, so the exit 2 row is not silent by accident.
func TestExemptTreeRowFollowsExitRule(t *testing.T) {
	exempt := boundExempt(boundReasonArtifact)
	row, run := plantedTreeVerb(t, exempt, exitingHandler(1, 1))
	if run.code != 1 || run.stdout != numberedLines(1) || run.stderr != row {
		t.Fatalf("exempt exit 1 = (%d, %q, %q), want the row %q on stderr", run.code, run.stdout, run.stderr, row)
	}
	_, run = plantedTreeVerb(t, exempt, exitingHandler(1, 2))
	if run.code != 2 || strings.Contains(run.stdout+run.stderr, "tree[") {
		t.Fatalf("exempt exit 2 = (%d, %q, %q), want no row on either stream", run.code, run.stdout, run.stderr)
	}
}

// The dispatcher computes the row before the verb runs, so the row names the tree that the
// verb read. A verb that makes the tree dirty still prints the clean row, on both paths.
func TestTreeRowPrecedesVerb(t *testing.T) {
	dirty := func(c Command, args []string) int {
		writeAXIFixture(t, "untracked.txt", "new\n")
		return exitingHandler(1, 0)(c, args)
	}
	row, run := plantedTreeVerb(t, boundResponse, dirty)
	if want := row + numberedLines(1); run.code != 0 || run.stdout != want {
		t.Fatalf("bounded dirtying verb = (%d, %q, %q), want (0, %q)", run.code, run.stdout, run.stderr, want)
	}
	row, run = plantedTreeVerb(t, boundExempt(boundReasonArtifact), dirty)
	if run.code != 0 || run.stderr != row {
		t.Fatalf("exempt dirtying verb = (%d, %q, %q), want the row %q on stderr", run.code, run.stdout, run.stderr, row)
	}
}
