// Package shift owns the gated loop: the iteration and refactor phases, the touched-path
// staging diff, the adapter preflight, and the `.bench/done.sh` early-completion check.
// It also owns the SIGINT/SIGTERM teardown and the iteration/refactor prompt text, as Go
// string constants. The staging diff parses porcelain `-z` natively, which retires the
// shell's newline misread. It runs acquire, loop, and release inside one process, since
// lease ownership is the recording process's pid. It calls internal/worktree for the
// pool, internal/gate for the oracle in-process, and internal/structure for the
// touched-scope refactor detector.
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
// set of dirty paths, scratch excluded. The NUL framing lets a path with spaces, glob
// characters, or a literal newline survive whole. This is the case the shell's
// `printf | sort` pipeline misread. git.ParsePorcelainZ owns that framing; this layers
// the shift's scratch policy and stable ordering on top.
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
	// Audit #11: tolerate. This runs inside a shift-created, validated worktree, so an
	// empty parse stages nothing, and the loop's own commit/gate is the loud oracle for
	// that worktree.
	raw, _ := git.Raw("-C", root, "status", "--porcelain", "-z", "--no-renames")
	return parseDirtyPaths(raw)
}

// stageTouched stages exactly what the agent touched: the paths dirty after it ran
// (post) minus those dirty before (pre). Snapshotting pre before the gate runs keeps a
// gate byproduct, such as an unignored build artifact, and any pre-existing dirt out of
// the commit. The `:(literal)` pathspec keeps a glob character in a path from being read
// as a pattern. A staging failure, a real `git add` error or the injected test fault,
// returns an error rather than being swallowed. The caller can then snapshot and split
// by evidence instead of committing a partial tree.
func stageTouched(root string, pre, post []string) error {
	before := make(map[string]bool, len(pre))
	for _, p := range pre {
		before[p] = true
	}
	for _, p := range post {
		if p == "" || before[p] {
			continue
		}
		if err := exec.Command("git", "-C", root, "add", "-A", "--", ":(literal)"+p).Run(); err != nil {
			return fmt.Errorf("stage %s: %w", p, err)
		}
	}
	return hitShift(shiftFault, stepStage)
}

// resolveExecutable mirrors bash `type -P`. It returns the path when name resolves to
// an executable file, absolute, relative, or on PATH. It returns "" for a shell keyword,
// builtin, or a name that does not resolve. A BENCH_AGENT that is not a runnable file is
// rejected in the preflight rather than exec-failing silently every iteration.
func resolveExecutable(name string) string {
	p, err := exec.LookPath(name)
	if err != nil {
		return ""
	}
	return p
}

// requireAdapter is the preflight the loop runs before acquiring a worktree or entering
// any iteration. An unset or empty BENCH_AGENT is a configure-your-adapter error. A
// value that does not resolve to an executable file, a missing path or a shell keyword,
// is a not-executable error. Both fail fast, before any agent or gate run.
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
// the objective, the project's per-repo early-completion override. Absent or
// non-executable means "run to the iteration cap". It runs from the worktree root, so
// done.sh sees the shift's in-progress tree.
func objectiveMet(root string, objective objective) bool {
	done := filepath.Join(root, ".bench", "done.sh")
	info, err := os.Stat(done)
	if err != nil || info.IsDir() || info.Mode()&0o111 == 0 {
		return false
	}
	cmd := exec.Command(done, objective.predicateArgument())
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
// argument joined, mirroring the shell's `$*`. An empty objective is not defaulted;
// Loop's validation rejects it with exit 2.
func Command(args []string, stdout, stderr io.Writer) int {
	refresh := len(args) > 0 && args[0] == "--refresh"
	if refresh {
		args = args[1:]
	}
	objective := strings.Join(args, " ")
	return loop(objective, refresh, stdout, stderr)
}
