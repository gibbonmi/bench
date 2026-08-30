// Package preflight is the start-oracle over artifacts-vs-reality. The
// gate answers "is the tree done?"; this package answers "do the declared
// artifacts still describe reality?" at the moment a phase starts.
//
// It follows the decision-domain split. A thin gatherer (git, the exported
// diff-base resolution, the spec resolver, the coverage parser, tickets/
// enumeration) collects immutable Facts. Decide is a pure function over
// those facts with no I/O of its own, so the five checks are table-tested
// without a repository.
package preflight

import "strings"

// Facts is the immutable evidence Decide classifies. Every field is
// gathered once, up front, by Gather, and never re-read mid-verdict. So a
// rerun over the same Facts value always answers the same way.
type Facts struct {
	Mode string

	// SpecPath is the resolved spec's repo-relative path, printed on the `spec:` line.
	SpecPath string

	// DefaultBranchResolved and DefaultBranchCurrent back base-current. An
	// unresolved default branch is red on its own; a resolved one is current
	// only when its tip is an ancestor of HEAD.
	DefaultBranchResolved bool
	DefaultBranchCurrent  bool

	// ReviewBaseResolved and ReviewBaseHint back both paths-authorized and
	// diff-nonempty: neither check can answer without a resolved review base.
	ReviewBaseResolved bool
	ReviewBaseHint     string

	// SourceBase is the base an explicit-base invocation pinned; it alone gates the
	// `source` presentation, so an implicit-base invocation still prints no source
	// table. SourceTip is the derived source tip either way — the snapshot head an
	// explicit range captured, or HEAD.
	SourceBase string
	SourceTip  string

	// PinnedSourceTip is the commit --source-tip named, already resolved to its full
	// identity by the gatherer, and empty when the flag is omitted. Resolving it is
	// the gatherer's job; whether it agrees with SourceTip is Decide's.
	PinnedSourceTip string

	// ExplicitSourceRange records that the invocation supplied --base. Its validity
	// comes from the same resolved source range that supplies ReviewBaseResolved,
	// rather than from destination default-branch ancestry.
	ExplicitSourceRange bool

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

	// DefaultBranch is the resolved default branch name, empty when it does not
	// resolve. The stale-base remedy names it, so the printed command can never
	// name a branch this repository does not have.
	DefaultBranch string

	// AssignmentTarget is the id of the active assignment whose worktree is this
	// preflight root, empty in the primary checkout and in any tree no active
	// assignment owns. The stale-base remedy addresses that id.
	AssignmentTarget string

	// TicketsDirExists reports whether specs/<slug>/tickets/ exists at all;
	// present-but-empty counts as existing. Build mode uses it to tell an
	// absent tickets/ (row checks not-applicable) from a present one. A
	// present, empty directory runs row checks for real, reading as
	// unowned rows rather than a pass.
	TicketsDirExists bool
}

// CheckResult is one verdict row: the check's name, its verdict ("green" or
// "red"), a detail string, and the remedy that answers it. The detail is empty
// for green, and names the offending path/ID(s) for red. Next is empty on every
// row but the one red a single command repairs, because a remedy on a row that
// no command answers sends the reader down a false path.
type CheckResult struct {
	Check, Verdict, Detail, Next string
}

const (
	verdictGreen = "green"
	verdictRed   = "red"
	verdictNA    = "not-applicable"

	modeBuild = "build"
)

// Verdict is Decide's complete answer: the check rows, in fixed order, and
// whether any of them is red, the caller's exit-code source. A
// not-applicable row never contributes to Red. It is a printed,
// definitive verdict in its own right, not a soft pass standing in for a
// real one.
type Verdict struct {
	Checks []CheckResult
	Red    bool
}

// Decide classifies immutable Facts into the verdict. It performs no I/O
// and consults nothing but its argument, so the same Facts value always
// yields the same Verdict; the byte-identical-rerun guarantee lives here.
//
// Mode applicability lives here rather than in the gatherer. An explicit
// source range makes base-current grade that range's validity, while a
// bare invocation grades default branch ancestry. Build mode always runs
// paths-authorized, runs rows-owned and rows-membership for real only when
// specs/<slug>/tickets/ exists, and never runs diff-nonempty.
//
// tip-current is the one conditional row. It appears only when
// --source-tip pinned a tip, directly after base-current. So the two
// halves of the source identity are graded together, ahead of every
// check that presupposes that identity. An invocation with no pin renders
// exactly the five rows it always has.
func Decide(f Facts) Verdict {
	checks := []CheckResult{baseCurrentCheck(f)}
	if f.PinnedSourceTip != "" {
		checks = append(checks, tipCurrentCheck(f))
	}
	checks = append(checks,
		pathsAuthorizedCheck(f),
		rowsOwnedRow(f),
		rowsMembershipRow(f),
		diffNonemptyRow(f),
	)
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
	if f.ExplicitSourceRange {
		if !f.ReviewBaseResolved {
			return red("base-current", "source base does not resolve: "+f.ReviewBaseHint)
		}
		return green("base-current")
	}
	if !f.DefaultBranchResolved {
		return red("base-current", "default branch does not resolve")
	}
	if !f.DefaultBranchCurrent {
		return staleBaseRed(f)
	}
	return green("base-current")
}

// staleBaseRed is the one red that carries its own remedy: the default branch has
// moved ahead of this tree, and one merge repairs it. A root no active assignment
// owns renders the placeholder id, so the operator reads where their own target
// goes rather than a command with a hole in it.
func staleBaseRed(f Facts) CheckResult {
	target := f.AssignmentTarget
	if target == "" {
		target = "<target>"
	}
	row := red("base-current", "default branch tip is not an ancestor of HEAD")
	row.Next = "bench worktree merge --from " + f.DefaultBranch + " " + target
	return row
}

// tipCurrentCheck verifies the reviewer's frozen tip against the one
// preflight derived. Both values are already-resolved full identities.
// So a pin spelled as a branch or as HEAD is green whenever it names
// the same commit; the check grades agreement, not spelling.
func tipCurrentCheck(f Facts) CheckResult {
	if f.SourceTip == "" {
		return red("tip-current", "source tip does not resolve, so --source-tip "+f.PinnedSourceTip+" cannot be verified")
	}
	if f.PinnedSourceTip != f.SourceTip {
		return red("tip-current", "--source-tip "+f.PinnedSourceTip+" is not the derived source tip "+f.SourceTip)
	}
	return green("tip-current")
}

func pathsAuthorizedCheck(f Facts) CheckResult {
	if !f.ReviewBaseResolved {
		return red("paths-authorized", "review base does not resolve: "+f.ReviewBaseHint)
	}
	var unauthorized []string
	entries := authorizingEntries(f)
	for _, p := range f.ChangedPaths {
		if !fenceAuthorizes(p, entries) {
			unauthorized = append(unauthorized, p)
		}
	}
	if len(unauthorized) > 0 {
		return red("paths-authorized", "not authorized by any ownership fence: "+strings.Join(unauthorized, ", "))
	}
	return green("paths-authorized")
}

// captureEntry is the phase-owned capture folder. Every phase close writes the
// handoff, the learnings, or the ideas there, so a reviewed range always carries it,
// and no spec fences it. It is authorized for every range.
const captureEntry = "capture"

// authorizingEntries is every entry paths-authorized consults. This is
// the spec's declared fence entries, the phase-owned capture folder, plus
// the active spec's own folder, which authorizes an in-range amendment of
// the spec without a self-fence entry.
//
// The implicit entries are derived rather than carried as their own facts,
// so the printed spec path and the authorized folder cannot disagree. They
// are appended to a copy so the gathered FenceEntries slice is never
// mutated. Mode is deliberately not consulted: build preflight, review
// preflight, and the landing's final source authorization all get the same
// answer.
func authorizingEntries(f Facts) []string {
	entries := append(append([]string{}, f.FenceEntries...), captureEntry)
	if folder := specFolder(f.SpecPath); folder != "" {
		entries = append(entries, folder)
	}
	return entries
}

// specFolder is the directory containing the resolved spec, empty when the
// spec path carries no directory at all. The result is an ordinary fence
// entry, so the segment-boundary rule below is the one that grades it,
// never a second prefix rule.
func specFolder(specPath string) string {
	i := strings.LastIndex(specPath, "/")
	if i < 0 {
		return ""
	}
	return specPath[:i]
}

// fenceAuthorizes reports whether path is covered by one of the spec's
// declared fence entries: an exact match, or a `/`-separated prefix.
// `internal/git` never authorizes `internal/git2`, only `internal/git`
// itself or anything under `internal/git/`.
//
// A fence entry conventionally spelled with its own trailing slash (a
// directory marker, e.g. `internal/preflight/`) is normalized before
// comparison, so the trailing slash is never itself an extra path segment.
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
