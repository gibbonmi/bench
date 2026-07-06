// Package commit owns `bench commit [-m <msg>] [--spec <slug>] <path>...`: the thin
// orchestrator that mechanizes "commit on green, never on red" so the invariant lives in
// code, not in prose the agent must remember. It sequences block-check → gate → flip →
// stage → commit: it refuses before gating if any working-tree file outside the named set
// (plus the --spec file) is dirty, runs the project gate through internal/gate and commits
// only on green, flips the spec through internal/spec when --spec is set, and stages
// exactly the named paths via a `:(literal)` pathspec (a named deletion included) —
// never a bare `git add -A` over the whole tree.
// It forms no opinion of the gate's verdict and carries no branch guard: the pre-push hook
// owns default-branch protection, so commit is branch-agnostic.
package commit

import (
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/spec"
)

// Command runs the gated commit. Usage errors (no paths, no -m, unknown flag) exit 2;
// operational failures (block, gate-red, flip failure, empty commit, git error) exit 1;
// a green gate that commits exits 0. The gate's live output streams to stdout/stderr, so
// a red gate reports its own first failing phase.
func Command(args []string, stdout, stderr io.Writer) int {
	msg, specSlug, paths, usageErr := parseArgs(args)
	if usageErr != "" {
		fmt.Fprintln(stderr, "usage: bench commit -m <msg> [--spec <slug>] <path>... ("+usageErr+")")
		return 2
	}

	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, "error: not in a git repository — run inside a Bench-linked repo")
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

	if rc := gate.RunAndRecord(root, stdout, stderr); rc != 0 {
		fmt.Fprintln(stderr, "error: gate is red — commit refused (see the failing phase above)")
		return 1
	}

	if specSlug != "" {
		if _, flipErr := spec.Flip(root, specSlug); flipErr != nil {
			fmt.Fprintf(stderr, "error: %v\n", flipErr)
			return 1
		}
	}

	for _, p := range named {
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

// parseArgs pulls -m <msg>, --spec <slug>, and the positional paths out of args. usageErr
// is non-empty on any misuse (no -m, no paths, unknown flag, a flag missing its value).
func parseArgs(args []string) (msg string, specSlug string, paths []string, usageErr string) {
	msgSet := false
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-m":
			if i+1 >= len(args) {
				return "", "", nil, "-m needs a message"
			}
			i++
			msg, msgSet = args[i], true
		case a == "--spec":
			if i+1 >= len(args) {
				return "", "", nil, "--spec needs a slug"
			}
			i++
			specSlug = args[i]
		case a == "-h" || a == "--help":
			return "", "", nil, "help"
		case strings.HasPrefix(a, "-"):
			return "", "", nil, "unknown flag: " + a
		default:
			paths = append(paths, a)
		}
	}
	if !msgSet {
		return "", "", nil, "-m <msg> is required"
	}
	if len(paths) == 0 {
		return "", "", nil, "at least one <path> is required"
	}
	return msg, specSlug, paths, ""
}

// unexplained lists the working-tree paths (tracked-modified or untracked) that are not in
// the allowed set, sorted. An empty result means the tree equals the named diff.
func unexplained(root string, allowed []string) []string {
	allow := make(map[string]bool, len(allowed))
	for _, p := range allowed {
		allow[p] = true
	}
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
