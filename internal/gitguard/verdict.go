package gitguard

import (
	"path"
	"strings"
)

// classify maps a git subcommand and its post-subcommand args to a deny label, or "" to
// allow. Each destructive class has its own verb rule, including the
// carve-outs for harness-delegate scratch (worktree-* branches, .claude/worktrees/
// paths). The two verdicts that need repo truth (checkout ref-ness, forced-creation
// clobber) call through chk; every other verdict is pure.
func classify(sub string, args []string, viaXargs bool, chk Checker) string {
	switch sub {
	case "push":
		return denyLabels["push"]
	case "reset":
		if contains(args, "--hard") {
			return denyLabels["reset"]
		}
	case "clean":
		if contains(args, "--force") || anyShortFlagHas(args, "f") {
			return denyLabels["clean"]
		}
	case "branch":
		if branchVerdict(args) {
			if hasAny(args, "-f", "--force") && !hasAny(args, "-D", "-d", "--delete") {
				return denyLabels["branch-force"]
			}
			return denyLabels["branch-delete"]
		}
	case "checkout":
		if checkoutVerdict(args, viaXargs, chk) {
			return denyLabels["checkout"]
		}
	case "switch":
		if switchVerdict(args, chk) {
			return denyLabels["switch"]
		}
	case "restore":
		if restoreVerdict(args, viaXargs) {
			return denyLabels["restore"]
		}
	case "rebase":
		return denyLabels["rebase"]
	case "stash":
		if key := stashVerdict(args); key != "" {
			return denyLabels[key]
		}
	case "commit":
		if contains(args, "--amend") {
			return denyLabels["amend"]
		}
	case "update-ref":
		if contains(args, "-d") {
			return denyLabels["update-ref"]
		}
	case "tag":
		if contains(args, "-d") || contains(args, "--delete") {
			return denyLabels["tag"]
		}
	case "reflog":
		if reflogVerdict(args) {
			return denyLabels["reflog"]
		}
	case "worktree":
		if worktreeVerdict(args) {
			return denyLabels["worktree"]
		}
	}
	return ""
}

func branchVerdict(args []string) bool {
	if !hasAny(args, "-D", "-d", "--delete", "-f", "--force") {
		return false
	}
	// Force-move (-f/--force without a delete flag) always blocks.
	if hasAny(args, "-f", "--force") && !hasAny(args, "-D", "-d", "--delete") {
		return true
	}
	// Carve-out: deleting harness-delegate branches (worktree-*) is cleanup of
	// agent-created scratch, not reviewer history.
	names := freeArgs(args)
	if len(names) == 0 {
		return true
	}
	for _, name := range names {
		if !strings.HasPrefix(name, "worktree-") {
			return true
		}
	}
	return false
}

func checkoutVerdict(args []string, viaXargs bool, chk Checker) bool {
	if viaXargs {
		return true
	}
	if hasAny(args, "-f", "--force") {
		return true
	}
	if hasPathspecFile(args) {
		return true
	}
	if target, ok := forcedCreationTarget(args, "-B"); ok && chk.BranchExists(target) {
		return true
	}
	if contains(args, "--") {
		return hasExplicitPathspec(args)
	}
	free := checkoutFreeArgs(args)
	if len(free) == 0 {
		return false
	}
	if len(free) > 1 {
		return true
	}
	return !chk.RefResolves(free[0])
}

func switchVerdict(args []string, chk Checker) bool {
	if hasAny(args, "-f", "--force", "--discard-changes") {
		return true
	}
	target, ok := forcedCreationTarget(args, "-C")
	return ok && chk.BranchExists(target)
}

func restoreVerdict(args []string, viaXargs bool) bool {
	if viaXargs {
		return true
	}
	staged := hasAny(args, "--staged", "-S")
	worktree := hasAny(args, "--worktree", "-W")
	if staged && !worktree {
		return false
	}
	return contains(args, ".") || restoreHasPathspec(args)
}

var stashReadOnly = map[string]bool{"list": true, "show": true}

// stashVerdict splits stash into its two deny classes, returning the deny key or "" to
// allow: "stash" for the verbs that destroy stash history, "stash-worktree" for the verbs
// that cross-apply working-tree state between worktrees. It fails closed — the allow set
// is exactly the read-only verbs, so a bare `git stash` (which is `push`) and any verb not
// enumerated here block. The subcommand is args[0] and nothing else, because git accepts
// one only there: scanning for the first *free* argument instead lets `git stash -m list`
// read as the allowed `list` while it pushes a stash and reverts the working tree. There
// is no viaXargs parameter, unlike checkoutVerdict and restoreVerdict: those turn
// destructive when paths arrive from stdin, whereas the argument-less stash form already
// blocks, so nothing xargs can supply loosens this verdict.
func stashVerdict(args []string) string {
	if len(args) == 0 {
		return "stash-worktree"
	}
	switch sub := args[0]; {
	case stashReadOnly[sub]:
		return ""
	case sub == "drop" || sub == "clear":
		return "stash"
	default:
		return "stash-worktree"
	}
}

func reflogVerdict(args []string) bool {
	sub, ok := firstFreeArg(args)
	return ok && sub == "expire"
}

func worktreeVerdict(args []string) bool {
	free := freeArgs(args)
	if len(free) == 0 || free[0] != "remove" {
		return false
	}
	if !hasAny(args, "-f", "--force") {
		return false
	}
	// Carve-out: force-removing harness-delegate worktrees under .claude/worktrees/ is
	// cleanup of agent-created scratch. Paths are normalized so traversal out of that
	// directory still blocks.
	paths := free[1:]
	if len(paths) == 0 {
		return true
	}
	for _, p := range paths {
		if !isDelegateWorktreePath(p) {
			return true
		}
	}
	return false
}

// --- shared helpers -----------------------------------------------------------

func hasExplicitPathspec(args []string) bool {
	idx := indexOf(args, "--")
	if idx < 0 {
		return false
	}
	for _, arg := range args[idx+1:] {
		if arg != "" {
			return true
		}
	}
	return false
}

func isPathspecFileArg(arg string) bool {
	return arg == "--pathspec-from-file" || strings.HasPrefix(arg, "--pathspec-from-file=")
}

func hasPathspecFile(args []string) bool {
	for _, arg := range args {
		if isPathspecFileArg(arg) {
			return true
		}
	}
	return false
}

func restoreHasPathspec(args []string) bool {
	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "--" {
			for _, rest := range args[i+1:] {
				if rest != "" {
					return true
				}
			}
			return false
		}
		if isPathspecFileArg(arg) {
			return true
		}
		if arg == "-s" || arg == "--source" {
			i += 2
			continue
		}
		if strings.HasPrefix(arg, "--source=") {
			i++
			continue
		}
		if strings.HasPrefix(arg, "-") {
			i++
			continue
		}
		return true
	}
	return false
}

// forcedCreationTarget returns the branch named by opt (`-B main` or `-Bmain`), and
// whether one was present.
func forcedCreationTarget(args []string, opt string) (string, bool) {
	for i, arg := range args {
		if arg == opt {
			if i+1 < len(args) {
				return args[i+1], true
			}
			return "", false
		}
		if strings.HasPrefix(arg, opt) && len(arg) > len(opt) {
			return arg[len(opt):], true
		}
	}
	return "", false
}

func checkoutFreeArgs(args []string) []string {
	var free []string
	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "-b" || arg == "-B" || arg == "--orphan" {
			i += 2
			continue
		}
		if strings.HasPrefix(arg, "-") {
			i++
			continue
		}
		free = append(free, arg)
		i++
	}
	return free
}

// isDelegateWorktreePath reports whether arg normalizes to a path under
// .claude/worktrees/ (traversal out of it via `..` normalizes away and no longer
// matches, so it still blocks).
func isDelegateWorktreePath(arg string) bool {
	norm := path.Clean(arg)
	return strings.HasPrefix(norm, ".claude/worktrees/") || strings.Contains(norm, "/.claude/worktrees/")
}

// --- tiny slice utilities ----------------------------------------------------

func contains(args []string, s string) bool { return indexOf(args, s) >= 0 }

func indexOf(args []string, s string) int {
	for i, a := range args {
		if a == s {
			return i
		}
	}
	return -1
}

func hasAny(args []string, opts ...string) bool {
	for _, arg := range args {
		for _, opt := range opts {
			if arg == opt {
				return true
			}
		}
	}
	return false
}

// anyShortFlagHas reports whether any `-…` short-flag cluster contains one of the given
// letters (the `clean -f`/`-fd` test: an arg starting with `-` whose remainder contains
// the letter).
func anyShortFlagHas(args []string, letters string) bool {
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && strings.ContainsAny(arg[1:], letters) {
			return true
		}
	}
	return false
}

// freeArgs returns the non-option args, in order.
func freeArgs(args []string) []string {
	var free []string
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			free = append(free, arg)
		}
	}
	return free
}

// firstFreeArg returns the first non-option arg and whether one exists.
func firstFreeArg(args []string) (string, bool) {
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			return arg, true
		}
	}
	return "", false
}
