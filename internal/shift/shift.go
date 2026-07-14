// Package shift owns the gated loop: the iteration and refactor phases, the
// touched-path staging diff (porcelain `-z` parsed natively, which retires the shell's
// newline misread), the adapter preflight, the `.bench/done.sh` early-completion
// check, the SIGINT/SIGTERM teardown, and the iteration/refactor prompt text as Go
// string constants. It runs acquire → loop → release inside one process (lease
// ownership is the recording process's pid), calling internal/worktree for the pool,
// internal/gate for the oracle in-process, and internal/structure for the touched-scope
// refactor detector.
package shift

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
)

// scratchFiles are the shift's per-run scratch, carried between iterations and excluded
// from the touched-path staging diff so they never ride into an iteration commit.
var scratchFiles = map[string]bool{
	".bench-objective": true,
	".bench-notes.md":  true,
}

// parseDirtyPaths turns `git status --porcelain -z --no-renames` output into the sorted
// set of dirty paths, scratch excluded. The NUL framing (a path with spaces, glob
// characters, or a literal newline survives whole — the case the shell's `printf | sort`
// pipeline misread) is git.ParsePorcelainZ's; this layers the shift's scratch policy and
// stable ordering on top.
func parseDirtyPaths(raw []byte) []string {
	var paths []string
	for _, entry := range git.ParsePorcelainZ(raw) {
		if entry.Path == "" || scratchFiles[entry.Path] {
			continue
		}
		paths = append(paths, entry.Path)
	}
	sort.Strings(paths)
	return paths
}

// dirtyPaths reads the current dirty set of the worktree at root.
func dirtyPaths(root string) []string {
	// Audit #11 — tolerate: runs inside a shift-created, validated worktree; an empty parse
	// stages nothing, and the loop's own commit/gate is the loud oracle for that worktree.
	raw, _ := git.Raw("-C", root, "status", "--porcelain", "-z", "--no-renames")
	return parseDirtyPaths(raw)
}

// stageTouched stages exactly what the agent touched: the paths dirty after it ran
// (post) minus those dirty before (pre). Snapshotting pre before the gate runs keeps a
// gate byproduct (an unignored build artifact) and any pre-existing dirt out of the
// commit; the `:(literal)` pathspec keeps a glob character in a path from being read as
// a pattern.
func stageTouched(root string, pre, post []string) {
	before := make(map[string]bool, len(pre))
	for _, p := range pre {
		before[p] = true
	}
	for _, p := range post {
		if p == "" || before[p] {
			continue
		}
		_ = exec.Command("git", "-C", root, "add", "-A", "--", ":(literal)"+p).Run()
	}
}

// resolveExecutable mirrors bash `type -P`: it returns the path when name resolves to
// an executable file (absolute, relative, or on PATH) and "" for a shell keyword,
// builtin, or a name that does not resolve — so a BENCH_AGENT that is not a runnable
// file is rejected in the preflight rather than exec-failing silently every iteration.
func resolveExecutable(name string) string {
	p, err := exec.LookPath(name)
	if err != nil {
		return ""
	}
	return p
}

// requireAdapter is the preflight the loop runs before acquiring a worktree or entering
// any iteration: an unset/empty BENCH_AGENT is a configure-your-adapter error, and a
// value that does not resolve to an executable file (a missing path, or a shell
// keyword) is a not-executable error. Both fail fast, before any agent or gate run.
func requireAdapter(agent string) error {
	if agent == "" {
		return fmt.Errorf("no harness adapter configured: set BENCH_AGENT to an adapter executable (references in .bench/adapters/)")
	}
	if resolveExecutable(agent) == "" {
		return fmt.Errorf("adapter not executable: BENCH_AGENT='%s' is neither an executable file nor a command on PATH", agent)
	}
	return nil
}

// objectiveMet reports whether an executable `.bench/done.sh` exists and exits zero for
// the objective — the project's per-repo early-completion override. Absent or
// non-executable means "run to the iteration cap". Run from the worktree root, so
// done.sh sees the shift's in-progress tree.
func objectiveMet(root, objective string) bool {
	done := filepath.Join(root, ".bench", "done.sh")
	info, err := os.Stat(done)
	if err != nil || info.IsDir() || info.Mode()&0o111 == 0 {
		return false
	}
	cmd := exec.Command(done, objective)
	cmd.Dir = root
	return cmd.Run() == nil
}

// cleanupScratch removes the shift's scratch files from the worktree root.
func cleanupScratch(root string) {
	for name := range scratchFiles {
		os.Remove(filepath.Join(root, name))
	}
}

// Command is the `bench shift [objective...]` entry: the objective is every positional
// argument joined (mirroring the shell's `$*`). An empty objective is not defaulted —
// Loop's validation rejects it with exit 2.
func Command(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	objective := strings.Join(args, " ")
	return Loop(objective, stdin, stdout, stderr)
}
