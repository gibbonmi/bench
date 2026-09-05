package gitguard

import "strings"

// The push family: the destination rule and its lexical classes, split from the verb
// switch in verdict.go along the one verdict that reads repository facts and the
// process working directory.

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
