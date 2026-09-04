package git

import "strings"

// detachedHead is the name CheckedOutBranch reports for a detached head. The literal is
// a state, not a branch, so no destination resolves from it.
const detachedHead = "HEAD"

// BarePushDestination is the single owner of the destination a bare `git push` targets
// from root. It reads the same config git reads, because git's own `@{push}` fails for a
// topic branch that has no remote-tracking ref. Under `simple`, `current`, or an unset
// push.default, the destination is the checked-out branch. Under `upstream` or
// `tracking`, it is the upstream branch's name on the remote. Under `matching` or
// `nothing`, on a detached HEAD, outside a repository, or on any probe error, it returns
// ("", false), so a caller that guards a push fails closed rather than guessing a branch.
func BarePushDestination(root string) (string, bool) {
	if !OK("-C", root, "rev-parse", "--git-dir") {
		return "", false
	}
	branch, err := CheckedOutBranch(root)
	if err != nil || branch == "" || branch == detachedHead {
		return "", false
	}
	// A missing push.default exits nonzero, and that unset value is the git default.
	mode, err := Output("-C", root, "config", "--get", "push.default")
	if err != nil {
		mode = ""
	}
	switch mode {
	case "", "simple", "current":
		return branch, true
	case "upstream", "tracking":
		return upstreamBranch(root, branch)
	}
	return "", false
}

// upstreamBranch returns the short name the upstream ref carries on its remote. The
// probe is `rev-parse --abbrev-ref @{upstream}`, which is one call and reports the
// upstream as `<remote>/<branch>`; `for-each-ref` would need the branch's own ref
// operand for the same answer. The remote name comes from branch.<name>.remote, so a
// branch name that itself contains a slash still strips correctly. An upstream that does
// not carry the configured remote's prefix reports no destination.
func upstreamBranch(root, branch string) (string, bool) {
	upstream, err := Output("-C", root, "rev-parse", "--abbrev-ref", "@{upstream}")
	if err != nil || upstream == "" {
		return "", false
	}
	remote, err := Output("-C", root, "config", "--get", "branch."+branch+".remote")
	if err != nil || remote == "" {
		return "", false
	}
	name := strings.TrimPrefix(upstream, remote+"/")
	if name == upstream || name == "" {
		return "", false
	}
	return name, true
}
