package gitguard

import (
	"path"
	"strings"
)

// classify maps a git subcommand and its post-subcommand args to a deny label, or "" to
// allow. Each destructive class has its own verb rule, including the
// carve-outs for harness-delegate scratch (worktree-* branches, .claude/worktrees/
// paths). The two verdicts that need repo truth (checkout ref-ness, forced-creation
// clobber) and the push destination rule call through chk; every other verdict is pure.
// Only the push rule reads redirected, because only the push facts come from the
// process working directory; every other verdict grades the tokens alone.
func classify(sub string, args []string, viaXargs, redirected bool, chk Checker) string {
	switch sub {
	case "push":
		if key := pushVerdict(args, viaXargs, redirected, chk); key != "" {
			return denyLabels[key]
		}
	case "reset":
		if contains(args, "--hard") {
			return denyLabels["reset"]
		}
	case "clean":
		if forced(args) {
			return denyLabels["clean"]
		}
	case "branch":
		if key := branchVerdict(args); key != "" {
			return denyLabels[key]
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
	case "filter-branch":
		return denyLabels["filter-branch"]
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
	case "stash":
		if key := stashVerdict(args); key != "" {
			return denyLabels[key]
		}
	case "rm":
		if rmVerdict(args) {
			return denyLabels["rm-force"]
		}
	}
	return ""
}

// pushVerdict returns the deny key for a push, or "" to allow. A push publishing a topic
// branch is ordinary agent work; the merge into the default branch, a rewrite, a
// deletion, and the three broadcast forms stay the reviewer's. The lexical classes are
// read first, then the destination of every refspec. No option exempts a push, so
// `--dry-run` reaches the same verdict as the real run. The rule fails closed: an
// unresolvable default branch or destination denies with the unresolved key. Two states
// leave the destination unreadable before any token is graded: a redirect names a
// repository other than the one the facts describe, and an xargs prefix appends the
// destination from stdin. Both deny first.
func pushVerdict(args []string, viaXargs, redirected bool, chk Checker) string {
	if viaXargs || redirected {
		return "push-unresolved"
	}
	free := pushFreeArgs(args)
	var refspecs []string
	if len(free) > 1 {
		refspecs = free[1:]
	}
	if pushForced(args, refspecs) {
		return "push-force"
	}
	if pushDeletes(args, refspecs) {
		return "push-delete"
	}
	for _, broadcast := range []struct{ opt, key string }{
		{"--all", "push-all"},
		{"--mirror", "push-mirror"},
		{"--tags", "push-tags"},
	} {
		if contains(args, broadcast.opt) {
			return broadcast.key
		}
	}

	def, ok := chk.defaultBranch()
	if !ok {
		return "push-unresolved"
	}
	if len(refspecs) == 0 {
		dest, ok := chk.bareDestination()
		if !ok {
			return "push-unresolved"
		}
		refspecs = []string{dest}
	}
	for _, spec := range refspecs {
		dest, ok := pushDestination(spec, chk)
		if !ok {
			return "push-unresolved"
		}
		if dest == def {
			return "push-default"
		}
	}
	return ""
}

// pushValueOptions are the push options that take a separate value. Each is skipped with
// its value, so `git push -o ci.skip origin topic` reads `origin` as the remote rather
// than `ci.skip`. The `=value` spelling is one token and needs no entry.
var pushValueOptions = []string{"-o", "--push-option", "--repo", "--receive-pack", "--exec"}

// pushFreeArgs returns push's non-option args in order: the remote first, then every
// refspec. Tokens after `--` are free args, the way checkoutFreeArgs reads its own.
func pushFreeArgs(args []string) []string {
	var free []string
	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "--" {
			return append(free, args[i+1:]...)
		}
		if hasAny(pushValueOptions, arg) {
			i += 2
			continue
		}
		if strings.HasPrefix(arg, "-") && arg != "-" {
			i++
			continue
		}
		free = append(free, arg)
		i++
	}
	return free
}

// pushForced reads force in each of its spellings: the two whole options, a short
// cluster carrying `f` (`-fu`), the lease option with or without a value, and the `+`
// refspec prefix. `--force-if-includes` only adds a safety check and is not force.
func pushForced(args, refspecs []string) bool {
	for _, arg := range args {
		if arg == "--force-with-lease" || strings.HasPrefix(arg, "--force-with-lease=") {
			return true
		}
	}
	if hasAny(args, "-f", "--force") || anyShortFlagHas(args, "f") {
		return true
	}
	for _, spec := range refspecs {
		if strings.HasPrefix(spec, "+") {
			return true
		}
	}
	return false
}

// pushDeletes reads deletion in each of its spellings: the two options and a refspec
// whose source is empty (`:topic`, or `:` alone).
func pushDeletes(args, refspecs []string) bool {
	if hasAny(args, "--delete", "-d") {
		return true
	}
	for _, spec := range refspecs {
		if strings.HasPrefix(spec, ":") {
			return true
		}
	}
	return false
}

// pushDestination resolves the branch a refspec updates, and whether it resolved. A
// `src:dst` refspec targets `dst` and a bare `src` targets itself, so the source name
// never decides the verdict. `HEAD`, and its `@` synonym, target the checked-out branch,
// and a detached HEAD leaves it unresolved. Git reads a bare `heads/` prefix as
// `refs/heads/`, so both spellings of the full ref strip to the branch name; a
// destination under any other `refs/` namespace names no branch and is never protected.
func pushDestination(spec string, chk Checker) (string, bool) {
	dest := strings.TrimPrefix(spec, "+")
	if i := strings.Index(dest, ":"); i >= 0 {
		dest = dest[i+1:]
	}
	if dest == "HEAD" || dest == "@" {
		return chk.checkedOut()
	}
	if strings.HasPrefix(dest, "refs/heads/") {
		return strings.TrimPrefix(dest, "refs/heads/"), true
	}
	if strings.HasPrefix(dest, "refs/") {
		return "", true
	}
	if strings.HasPrefix(dest, "heads/") {
		return strings.TrimPrefix(dest, "heads/"), true
	}
	return dest, true
}

func branchVerdict(args []string) string {
	if !hasAny(args, "-D", "-d", "--delete", "-f", "--force") {
		return ""
	}
	// Force-move (-f/--force without a delete flag) always blocks.
	if hasAny(args, "-f", "--force") && !hasAny(args, "-D", "-d", "--delete") {
		return "branch-force"
	}
	// Carve-out: deleting harness-delegate branches (worktree-*) is cleanup of
	// agent-created scratch, not reviewer history.
	names := freeArgs(args)
	if len(names) == 0 {
		return branchDeleteKey(args)
	}
	for _, name := range names {
		if !strings.HasPrefix(name, "worktree-") {
			return branchDeleteKey(args)
		}
	}
	return ""
}

func branchDeleteKey(args []string) string {
	if hasAny(args, "-D", "-f", "--force") {
		return "branch-delete-force"
	}
	return "branch-delete-safe"
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

// stashVerdict blocks the two stash operations that discard work, leaving the rest of
// the verb (bare, list, pop, push, apply, …) usable. The operation is the first free
// arg, resolved the way reflog resolves `expire`. Options are skipped but their
// values are not, and neither is anything after `--`. So `git stash -m drop` and
// `git stash -- drop` reach the same verdict as the positional spelling.
func stashVerdict(args []string) string {
	op, ok := firstFreeArg(args)
	if !ok {
		return ""
	}
	switch op {
	case "drop":
		return "stash-drop"
	case "clear":
		return "stash-clear"
	}
	return ""
}

// rmVerdict blocks only the recursive-and-forced combination, reading short-flag
// clusters the way clean's force test does (so `-rf` counts as both). An ordinary
// `git rm <path>`, a `git rm -r <dir>`, and a `git rm --cached <path>` are ordinary
// file work and stay allowed. git rm spells recursion only as `-r`/`-R`, so there is no
// long form to test beside the cluster.
func rmVerdict(args []string) bool {
	return forced(args) && anyShortFlagHas(args, "rR")
}

// --- shared helpers -----------------------------------------------------------

// forced reports whether the invocation carries git's force option in either spelling.
// clean and rm both deny on it, so the two rules read it from one place.
func forced(args []string) bool {
	return contains(args, "--force") || anyShortFlagHas(args, "f")
}

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
// .claude/worktrees/. Traversal out of it via `..` normalizes away and no longer
// matches, so it still blocks.
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

// anyShortFlagHas reports whether any short-flag cluster contains one of the given
// letters. Such a cluster is an arg starting with a single `-` whose remainder
// contains the letter (the `clean -f`/`-fd` test). A `--long` option is not a cluster
// — the `r` in `--force` is part of a name, not a flag. So each rule spells its long
// forms out separately.
func anyShortFlagHas(args []string, letters string) bool {
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") && strings.ContainsAny(arg[1:], letters) {
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
