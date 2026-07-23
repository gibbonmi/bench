// Package commit owns `bench commit -m <msg> [--spec <slug>] <path>...`: the thin
// orchestrator that mechanizes "commit on green, never on red" so the invariant lives in
// code, not in prose the agent must remember. It sequences block-check → gate → flip →
// stage → commit: it refuses before gating if any working-tree file outside the named set
// (plus the --spec file) is dirty, runs the project gate through internal/gate and commits
// only on green (reusing a fresh green verdict already recorded for the identical closed
// oracle subject instead of paying the gate twice), flips the spec through internal/spec when --spec is set, and stages
// exactly the named paths via a `:(literal)` pathspec (a named deletion included) —
// never a bare `git add -A` over the whole tree. A named path whose removal is already
// staged (`git rm`, a rename's old half) matches no add-pathspec and is recognized as
// already in the index rather than failed.
// It forms no opinion of the gate's verdict and carries no branch guard: the pre-push hook
// owns default-branch protection, so commit is branch-agnostic.
package commit

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// Command runs the gated commit. `help`, `--help`, and `-h` print the declared help at
// exit 0; usage errors (no paths, no -m, unknown flag) exit 2; operational failures (block, gate-red, flip failure, empty commit, git error) exit 1;
// a green gate that commits exits 0. The gate's live output streams to stdout/stderr, so
// a red gate reports its own first failing phase.
func Command(args []string, stdout, stderr io.Writer) int {
	msg, specSlug, paths, help, usageErr := parseArgs(args)
	if help != "" {
		fmt.Fprintln(stdout, help)
		return 0
	}
	if usageErr != "" {
		fmt.Fprintln(stderr, "usage: bench commit -m <msg> [--spec <slug>] <path>... (--spec marks the spec implemented; "+usageErr+")")
		return 2
	}

	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, toon.NotInRepo())
		return 1
	}

	// The named set the commit will land, root-relative and slash-formed to match
	// git's porcelain output and pathspecs.
	named := make([]string, 0, len(paths)+1)
	for _, p := range paths {
		rel, relErr := rootRel(root, p)
		if relErr != nil {
			fmt.Fprintf(stderr, "error: cannot resolve path %q relative to repo root: %v\n", p, relErr)
			return 1
		}
		named = append(named, rel)
	}

	// Resolve the --spec file up front so it joins the allowed set (it is still clean at
	// block-check; the flip happens only after a green gate) and so a bad slug fails fast.
	if specSlug != "" {
		_, resolved, tried, ok, resErr := spec.Resolve(root, specSlug)
		if resErr != nil {
			fmt.Fprintf(stderr, "error: --spec not readable: %s: %v\n", resolved, resErr)
			return 1
		}
		if !ok {
			fmt.Fprintf(stderr, "error: --spec not found: %s\n", strings.Join(tried, ", "))
			return 1
		}
		rel, relErr := rootRel(root, resolved)
		if relErr != nil {
			fmt.Fprintf(stderr, "error: cannot resolve --spec %q relative to repo root: %v\n", resolved, relErr)
			return 1
		}
		named = append(named, rel)
		// Fail fast: reject a bad or already-implemented --spec here, before the block-check
		// and the gate, so a spec the post-gate Flip would reject never burns a green gate.
		if _, checkErr := spec.CheckStaged(root, specSlug); checkErr != nil {
			fmt.Fprintf(stderr, "error: %v\n", checkErr)
			return 1
		}
	}

	// Block before gating on any dirty/untracked file outside the named set, so a green
	// verdict describes exactly the diff that lands.
	if offenders := unexplained(root, named); len(offenders) > 0 {
		fmt.Fprintln(stderr, "error: working-tree files outside the named set block the commit — name them, or set them aside:")
		for _, o := range offenders {
			fmt.Fprintf(stderr, "  %s\n", o)
		}
		return 1
	}

	// Classify the named paths before the gate runs, so a naming error never burns a
	// green run. Nothing between here and staging creates or deletes paths (the --spec
	// flip only edits file content), so the plan stays valid.
	toStage, planErr := stagePlan(root, named)
	if planErr != nil {
		fmt.Fprintf(stderr, "error: %v\n", planErr)
		return 1
	}

	// A green verdict recorded for the identical closed oracle subject already proves
	// this diff and the inputs that can affect its oracle (the block-check above pinned
	// the tree to the named set), so re-running the gate buys nothing but its full cost.
	// Anything less — stale, red, untrusted, or absent — pays the real gate run.
	if gv := gate.Inspect(root); gv.ReusableGreen {
		fmt.Fprintln(stdout, "gate: green (fresh verdict reused for this tree)")
	} else if result := gate.Execute(context.Background(), root, stdout, stderr); result.ActionExit != 0 {
		fmt.Fprintln(stderr, "error: gate is red — commit refused (see the failing phase above)")
		return 1
	}

	if specSlug != "" {
		if _, flipErr := spec.Flip(root, specSlug); flipErr != nil {
			fmt.Fprintf(stderr, "error: %v\n", flipErr)
			return 1
		}
	}

	for _, p := range toStage {
		if stageErr := exec.Command("git", "-C", root, "add", "-A", "--", ":(literal)"+p).Run(); stageErr != nil {
			fmt.Fprintf(stderr, "error: staging %q failed: %v\n", p, stageErr)
			return 1
		}
	}

	if nothingStaged(root) {
		fmt.Fprintln(stderr, "error: nothing to commit — the named paths produced no staged change")
		return 1
	}

	if commitErr := exec.Command("git", "-C", root, "commit", "-q", "-m", msg).Run(); commitErr != nil {
		fmt.Fprintf(stderr, "error: git commit failed: %v\n", commitErr)
		return 1
	}
	fmt.Fprintf(stdout, "committed %d path(s)\n", len(named))
	return 0
}

// grammar is the declared argument shape usage.Parse enforces — arity, flag recognition,
// `--`, and help all come from there rather than a local switch, which is what makes a
// path beginning with a dash expressible. MinArgs stays 0 so the no-paths case keeps its
// own named reason below rather than the generic missing-positional one.
var grammar = usage.Grammar{
	Cmd:     "bench commit",
	Help:    "usage: bench commit -m <msg> [--spec <slug>] [--] <path>...",
	Flags:   []usage.Flag{{Name: "-m", HasValue: true}, {Name: "--spec", HasValue: true}},
	MaxArgs: -1,
}

// parseArgs routes args through the one argument grammar and applies the two requirements
// arity alone cannot state: -m is mandatory, and at least one path must be named. help is
// non-empty when the invocation asked for help, which is a success the caller prints;
// usageErr is non-empty on any misuse.
func parseArgs(args []string) (msg string, specSlug string, paths []string, help string, usageErr string) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		if code == 0 {
			return "", "", nil, line, ""
		}
		return "", "", nil, "", line
	}
	msg, msgSet := parsed.Flags["-m"]
	if !msgSet {
		return "", "", nil, "", "-m <msg> is required"
	}
	if len(parsed.Positionals) == 0 {
		return "", "", nil, "", "at least one <path> is required"
	}
	return msg, parsed.Flags["--spec"], parsed.Positionals, "", ""
}

// stagePlan classifies each named path into what the staging loop can act on. `git add -A
// -- :(literal)p` is fatal (exit 128) when p matches nothing in the worktree or the index —
// exactly the state a staged removal leaves behind (`git rm`, or a `git mv` rename's old
// half). Such a path needs no staging: absent from worktree and index but present in HEAD
// means its deletion is already in the index. A path absent from all three is a naming
// error reported with a real message instead of git's raw exit status.
func stagePlan(root string, named []string) ([]string, error) {
	var stage []string
	for _, p := range named {
		if inWorktree(root, p) || inIndex(root, p) {
			stage = append(stage, p)
			continue
		}
		if !inHead(root, p) {
			return nil, fmt.Errorf("named path %q not found in worktree, index, or HEAD", p)
		}
	}
	return stage, nil
}

func inWorktree(root, p string) bool {
	_, err := os.Lstat(filepath.Join(root, filepath.FromSlash(p)))
	return err == nil
}

func inIndex(root, p string) bool {
	out, err := git.Raw("-C", root, "ls-files", "-z", "--", ":(literal)"+p)
	return err == nil && len(out) > 0
}

func inHead(root, p string) bool {
	return exec.Command("git", "-C", root, "cat-file", "-e", "HEAD:"+p).Run() == nil
}

// unexplained lists the working-tree paths (tracked-modified or untracked) that are not in
// the allowed set, sorted. An empty result means the tree equals the named diff.
func unexplained(root string, allowed []string) []string {
	allow := make(map[string]bool, len(allowed))
	for _, p := range allowed {
		allow[p] = true
	}
	// Audit #12 — tolerate (flagged for reviewer veto): an empty parse relaxes this
	// attribution guard, but the subsequent real `git commit`/`git add` fails loudly on a
	// broken repo, so no false-clean output escapes; hardening it into a loud error is
	// deferred (see the spec's Out of scope).
	// --untracked-files=all lists untracked files individually rather than collapsing a
	// new directory to `dir/`, so a named path inside a fresh directory matches instead
	// of reading as an unexplained offender.
	raw, _ := git.Raw("-C", root, "status", "--porcelain", "-z", "--no-renames", "--untracked-files=all")
	var offenders []string
	for _, e := range git.ParsePorcelainZ(raw) {
		if e.Path == "" || allow[e.Path] {
			continue
		}
		offenders = append(offenders, e.Path)
	}
	return offenders
}

// rootRel converts a path argument (as given, cwd-relative) to its slash-formed,
// repo-root-relative form for pathspec staging and porcelain comparison.
func rootRel(root, arg string) (string, error) {
	abs, err := filepath.Abs(arg)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(root, abs)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(rel), nil
}

// nothingStaged reports whether the index has no staged changes — the empty-commit guard.
func nothingStaged(root string) bool {
	return exec.Command("git", "-C", root, "diff", "--cached", "--quiet").Run() == nil
}
