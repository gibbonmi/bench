// Package preflight is the start-oracle over artifacts-vs-reality: the gate answers
// "is the tree done?", and this package answers "do the declared artifacts still
// describe reality?" at the moment a phase starts. It follows the decision-domain
// split — a thin gatherer (git, the exported diff-base resolution, the spec resolver,
// the coverage parser, tickets/ enumeration) collects immutable Facts, and Decide is a
// pure function over those facts with no I/O of its own, so the five checks are
// table-tested without a repository.
package preflight

import "strings"

// Facts is the immutable evidence Decide classifies. Every field is gathered once,
// up front, by Gather — never re-read mid-verdict, so a rerun over the same Facts
// value always answers the same way.
type Facts struct {
	Mode string

	// SpecPath is the resolved spec's repo-relative path, printed on the `spec:` line.
	SpecPath string

	// DefaultBranchResolved and DefaultBranchCurrent back base-current: an unresolved
	// default branch is red on its own; a resolved one is current only when its tip is
	// an ancestor of HEAD.
	DefaultBranchResolved bool
	DefaultBranchCurrent  bool

	// ReviewBaseResolved and ReviewBaseHint back both paths-authorized and
	// diff-nonempty: neither check can answer without a resolved review base.
	ReviewBaseResolved bool
	ReviewBaseHint     string

	// SourceBase and SourceTip are the full identities pinned by an explicit-base
	// invocation. Empty values preserve the existing implicit-base presentation.
	SourceBase string
	SourceTip  string

	// ChangedPaths is the changed-file set since the resolved review base. Explicit
	// source builds include committed, index, tracked-worktree, and untracked paths. An
	// empty set is a legitimate answer, not an unresolved one.
	ChangedPaths []string

	// FenceEntries is the spec's declared `## Ownership fences` tokens: backticked,
	// outside parentheses. paths-authorized checks every changed path against these.
	FenceEntries []string

	// DeclaredRowIDs is the spec's coverage map row IDs, in map order.
	DeclaredRowIDs []string

	// TicketTokens is every uppercase-tag token (`[A-Z]+[0-9]+`, word-boundary
	// matched) found across tickets/ file content, regardless of tag.
	TicketTokens []string

	// SpecTag is the alphabetic prefix shared by the spec's own declared row IDs
	// (e.g. "PF" for PF1..n) — the tag rows-membership scopes its check to, so a
	// foreign-tag token (FT93) is ignored rather than flagged.
	SpecTag string

	// TicketsDirExists reports whether specs/<slug>/tickets/ exists at all —
	// present-but-empty counts as existing. Build mode uses it to tell an absent
	// tickets/ (row checks not-applicable) from a present one (row checks run for
	// real, so an empty directory reads as unowned rows rather than a pass).
	TicketsDirExists bool
}

// CheckResult is one verdict row: the check's name, its verdict ("green" or "red"),
// and a detail string — empty for green, naming the offending path/ID(s) for red.
type CheckResult struct {
	Check, Verdict, Detail string
}

const (
	verdictGreen = "green"
	verdictRed   = "red"
	verdictNA    = "not-applicable"

	modeBuild = "build"
)

// Verdict is Decide's complete answer: the five check rows, in fixed order, and
// whether any of them is red — the caller's exit-code source. A not-applicable row
// never contributes to Red: it is a printed, definitive verdict in its own right,
// not a soft pass standing in for a real one.
type Verdict struct {
	Checks []CheckResult
	Red    bool
}

// Decide classifies immutable Facts into the five-check verdict. It performs no I/O
// and consults nothing but its argument, so the same Facts value always yields the
// same Verdict — the byte-identical-rerun guarantee lives here. Mode applicability
// lives here rather than in the gatherer: build mode always runs base-current and
// paths-authorized, runs rows-owned and rows-membership for real only when
// specs/<slug>/tickets/ exists, and never runs diff-nonempty.
func Decide(f Facts) Verdict {
	checks := []CheckResult{
		baseCurrentCheck(f),
		pathsAuthorizedCheck(f),
		rowsOwnedRow(f),
		rowsMembershipRow(f),
		diffNonemptyRow(f),
	}
	red := false
	for _, c := range checks {
		if c.Verdict == verdictRed {
			red = true
		}
	}
	return Verdict{Checks: checks, Red: red}
}

func green(name string) CheckResult { return CheckResult{Check: name, Verdict: verdictGreen} }
func red(name, detail string) CheckResult {
	return CheckResult{Check: name, Verdict: verdictRed, Detail: detail}
}
func notApplicable(name string) CheckResult {
	return CheckResult{Check: name, Verdict: verdictNA}
}

// rowsOwnedRow gates rowsOwnedCheck by build-mode ticket-directory applicability:
// not-applicable when build mode has no tickets/ directory at all, real otherwise.
func rowsOwnedRow(f Facts) CheckResult {
	if f.Mode == modeBuild && !f.TicketsDirExists {
		return notApplicable("rows-owned")
	}
	return rowsOwnedCheck(f)
}

// rowsMembershipRow mirrors rowsOwnedRow for rows-membership.
func rowsMembershipRow(f Facts) CheckResult {
	if f.Mode == modeBuild && !f.TicketsDirExists {
		return notApplicable("rows-membership")
	}
	return rowsMembershipCheck(f)
}

// diffNonemptyRow is not-applicable in build mode unconditionally: a build has no
// review base to require a nonempty diff against.
func diffNonemptyRow(f Facts) CheckResult {
	if f.Mode == modeBuild {
		return notApplicable("diff-nonempty")
	}
	return diffNonemptyCheck(f)
}

func baseCurrentCheck(f Facts) CheckResult {
	if !f.DefaultBranchResolved {
		return red("base-current", "default branch does not resolve")
	}
	if !f.DefaultBranchCurrent {
		return red("base-current", "default branch tip is not an ancestor of HEAD")
	}
	return green("base-current")
}

func pathsAuthorizedCheck(f Facts) CheckResult {
	if !f.ReviewBaseResolved {
		return red("paths-authorized", "review base does not resolve: "+f.ReviewBaseHint)
	}
	var unauthorized []string
	for _, p := range f.ChangedPaths {
		if !fenceAuthorizes(p, f.FenceEntries) {
			unauthorized = append(unauthorized, p)
		}
	}
	if len(unauthorized) > 0 {
		return red("paths-authorized", "not authorized by any ownership fence: "+strings.Join(unauthorized, ", "))
	}
	return green("paths-authorized")
}

// fenceAuthorizes reports whether path is covered by one of the spec's declared
// fence entries: an exact match, or a `/`-separated prefix — `internal/git` never
// authorizes `internal/git2`, only `internal/git` itself or anything under
// `internal/git/`. A fence entry conventionally spelled with its own trailing slash
// (a directory marker, e.g. `internal/preflight/`) is normalized before comparison so
// the trailing slash is never itself an extra path segment.
func fenceAuthorizes(path string, fences []string) bool {
	for _, fence := range fences {
		trimmed := strings.TrimSuffix(fence, "/")
		if path == trimmed || strings.HasPrefix(path, trimmed+"/") {
			return true
		}
	}
	return false
}

func rowsOwnedCheck(f Facts) CheckResult {
	var uncited []string
	for _, id := range f.DeclaredRowIDs {
		if !containsStr(f.TicketTokens, id) {
			uncited = append(uncited, id)
		}
	}
	if len(uncited) > 0 {
		return red("rows-owned", "declared row(s) cited by no ticket file: "+strings.Join(uncited, ", "))
	}
	return green("rows-owned")
}

func rowsMembershipCheck(f Facts) CheckResult {
	seen := map[string]bool{}
	var phantom []string
	for _, tok := range f.TicketTokens {
		if tagOf(tok) != f.SpecTag {
			continue
		}
		if containsStr(f.DeclaredRowIDs, tok) {
			continue
		}
		if !seen[tok] {
			seen[tok] = true
			phantom = append(phantom, tok)
		}
	}
	if len(phantom) > 0 {
		return red("rows-membership", "ticket token(s) under this spec's tag name no declared row: "+strings.Join(phantom, ", "))
	}
	return green("rows-membership")
}

func diffNonemptyCheck(f Facts) CheckResult {
	if !f.ReviewBaseResolved {
		return red("diff-nonempty", "review base does not resolve: "+f.ReviewBaseHint)
	}
	if len(f.ChangedPaths) == 0 {
		return red("diff-nonempty", "no changed files since the resolved review base")
	}
	return green("diff-nonempty")
}

func containsStr(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

// tagOf is a token's alphabetic prefix — the row-ID tag rows-membership scopes its
// comparison to, so a foreign tag (FT93) is ignored rather than flagged.
func tagOf(token string) string {
	i := 0
	for i < len(token) && token[i] >= 'A' && token[i] <= 'Z' {
		i++
	}
	return token[:i]
}
