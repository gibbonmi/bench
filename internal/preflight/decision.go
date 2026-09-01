// Package preflight is the start-oracle over artifacts-vs-reality. The
// gate answers "is the tree done?"; this package answers "do the declared
// artifacts still describe reality?" at the moment a phase starts.
//
// It follows the decision-domain split. A thin gatherer (git, the exported
// diff-base resolution, the spec resolver, the coverage parser, tickets/
// enumeration) collects immutable Facts. Decide is a pure function over
// those facts with no I/O of its own, so every check is table-tested
// without a repository.
package preflight

import (
	"strings"

	"github.com/gibbonmi/bench/internal/tickets"
)

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

	// Tickets is every `.md` ticket file under specs/<slug>/tickets/, parsed
	// through the ticket grammar, in enumeration order. Its `Covers:` tokens are
	// the one citation source: a row ID in prose is not evidence.
	Tickets []tickets.Ticket

	// TicketDiagnostics is the ordered grammar fault list across those files, each
	// named by its folder-relative path. The enumeration's own duplicate-identity
	// faults ride here too.
	TicketDiagnostics []string

	// BlockerCycles holds one edge of each blocker cycle across the parsed set.
	BlockerCycles []string

	// WritesPathExists reports, for each `Writes:` entry declared across the parsed
	// tickets, whether the path it names resolves in the tree. The gatherer owns the
	// probe; whether an absent path is a fault is Decide's.
	WritesPathExists map[string]bool

	// WritesFixturePins reports, for each `Writes:` entry, the canary fixture
	// directories that pin the path it names. The gatherer enumerates the live
	// inventory; whether an unnamed fixture is a fault is Decide's.
	WritesFixturePins map[string][]string

	// WritesBoundFiles reports, for each `Writes:` entry, the files the binding
	// registry binds to the package the entry writes into.
	WritesBoundFiles map[string][]string

	// WritesSystemTagged reports, for each `Writes:` entry, whether it names a Go
	// test file whose build constraint carries the system tag.
	WritesSystemTagged map[string]bool

	// TicketPinsKit reports, per ticket basename, whether the body states the
	// literal BENCH_KIT. kit-pin grades that statement against the tagged writes.
	TicketPinsKit map[string]bool

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

// Decide classifies a complete immutable Facts snapshot into a Verdict. It
// performs no I/O and accepts no state beyond f. The same Facts value always
// yields the same rows and Red value. Gather owns the snapshot and failure;
// Decide does not gather, validate inputs, or retry an interrupted read.
//
// Mode applicability lives here rather than in the gatherer. An explicit
// source range makes base-current grade that range's validity, while a
// bare invocation grades default branch ancestry. Build mode always runs
// paths-authorized, runs every ticket row for real only when
// specs/<slug>/tickets/ exists, and never runs diff-nonempty.
//
// tip-current is the one conditional row. It appears only when
// --source-tip pinned a tip, directly after base-current. So the two
// halves of the source identity are graded together, ahead of every
// check that presupposes that identity.
func Decide(f Facts) Verdict {
	checks := []CheckResult{baseCurrentCheck(f)}
	if f.PinnedSourceTip != "" {
		checks = append(checks, tipCurrentCheck(f))
	}
	checks = append(checks,
		pathsAuthorizedCheck(f),
		ticketRow(f, "tickets-parse", ticketsParseCheck),
		ticketRow(f, "blockers-resolve", blockersResolveCheck),
		ticketRow(f, "writes-resolve", writesResolveCheck),
		ticketRow(f, "fixture-closure", fixtureClosureCheck),
		ticketRow(f, "registry-closure", registryClosureCheck),
		ticketRow(f, "kit-pin", kitPinCheck),
		ticketRow(f, "rows-owned", rowsOwnedCheck),
		ticketRow(f, "rows-membership", rowsMembershipCheck),
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

// ticketRow gates one ticket-reading check by build-mode ticket-directory
// applicability: not-applicable when build mode has no tickets/ directory at
// all, real otherwise. Every row that reads a parsed ticket shares this one
// gate, so no two of them can drift on when a fresh build is graded.
func ticketRow(f Facts, name string, check func(Facts) CheckResult) CheckResult {
	if f.Mode == modeBuild && !f.TicketsDirExists {
		return notApplicable(name)
	}
	return check(f)
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
	if unauthorized := unauthorizedPaths(f); len(unauthorized) > 0 {
		return red("paths-authorized", "not authorized by any ownership fence: "+strings.Join(unauthorized, ", "))
	}
	return green("paths-authorized")
}

// unauthorizedPaths is every changed path no authorizing entry covers, in changed-path
// order. The row's detail sentence and the landing's typed path list both read this one
// answer, so a report and a refusal can never name a different set.
func unauthorizedPaths(f Facts) []string {
	var unauthorized []string
	entries := authorizingEntries(f)
	for _, p := range f.ChangedPaths {
		if !fenceAuthorizes(p, entries) {
			unauthorized = append(unauthorized, p)
		}
	}
	return unauthorized
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

// newMarker declares a `Writes:` entry as a path the ticket creates. An entry
// carrying it is green whether or not the tree already holds the path, because
// a blocker ticket may land the file first.
const newMarker = "(new)"

// splitWritesEntry separates one `Writes:` entry into the tree path it names
// and whether it carries the (new) marker. The gatherer's existence probe and
// the writes-resolve row read this one split, so the path probed and the path
// graded can never disagree.
func splitWritesEntry(entry string) (path string, isNew bool) {
	path = strings.TrimSpace(entry)
	if !strings.HasSuffix(path, newMarker) {
		return path, false
	}
	return strings.TrimSpace(strings.TrimSuffix(path, newMarker)), true
}

// ticketsParseCheck reports the ticket grammar itself: an absent required
// field, a duplicate field, an unterminated fence, a citation fault, a
// declared blocker that resolves against no sibling, or one basename claimed
// at two depths. Its detail is its own, so a grammar fault never hides behind
// the ownership rows below it.
func ticketsParseCheck(f Facts) CheckResult {
	if len(f.TicketDiagnostics) > 0 {
		return red("tickets-parse", "ticket grammar fault(s): "+strings.Join(f.TicketDiagnostics, "; "))
	}
	return green("tickets-parse")
}

// blockersResolveCheck reports the dependency faults only the whole parsed set
// can show: a blocker cycle leaves the coordinator an empty frontier with no
// signal. The per-ticket edge faults belong to the grammar row above.
func blockersResolveCheck(f Facts) CheckResult {
	if len(f.BlockerCycles) > 0 {
		return red("blockers-resolve", "blocker cycle(s): "+strings.Join(f.BlockerCycles, "; "))
	}
	return green("blockers-resolve")
}

// writesResolveCheck grades every declared ownership entry against the tree.
// An entry that names no path and claims no (new) file is a typo that would
// charge a delegate against nothing, so it reds with the entry named.
func writesResolveCheck(f Facts) CheckResult {
	var unresolved []string
	seen := map[string]bool{}
	for _, ticket := range f.Tickets {
		for _, entry := range ticket.Writes {
			if _, isNew := splitWritesEntry(entry); isNew {
				continue
			}
			named := ticket.Name + ": " + entry
			if seen[named] || f.WritesPathExists[entry] {
				continue
			}
			seen[named] = true
			unresolved = append(unresolved, named)
		}
	}
	if len(unresolved) > 0 {
		return red("writes-resolve", "Writes: entry names no tree path and carries no (new) marker: "+strings.Join(unresolved, ", "))
	}
	return green("writes-resolve")
}

// ownedPaths is every tree path one ticket declares, with the (new) marker
// stripped. The three closures below grade their required names against this one
// set, so no two of them can disagree about what a ticket already owns.
func ownedPaths(ticket tickets.Ticket) []string {
	paths := make([]string, 0, len(ticket.Writes))
	for _, entry := range ticket.Writes {
		path, _ := splitWritesEntry(entry)
		paths = append(paths, path)
	}
	return paths
}

// pathCovered reports whether one required path is already named by the ticket,
// either exactly or through a directory entry that contains it.
func pathCovered(required string, owned []string) bool {
	for _, path := range owned {
		if required == path || strings.HasPrefix(required, path+"/") {
			return true
		}
	}
	return false
}

// fixtureClosureCheck grades the red-capable fixture into the ticket. A ticket
// that edits a fixture-pinned line without naming the owning fixture directory
// leaves the proof outside the charge, and the bite breaks unnoticed.
func fixtureClosureCheck(f Facts) CheckResult {
	var unnamed []string
	seen := map[string]bool{}
	for _, ticket := range f.Tickets {
		owned := ownedPaths(ticket)
		for _, entry := range ticket.Writes {
			for _, fixture := range f.WritesFixturePins[entry] {
				if pathCovered(fixture, owned) {
					continue
				}
				named := ticket.Name + ": " + entry + " is pinned by " + fixture
				if seen[named] {
					continue
				}
				seen[named] = true
				unnamed = append(unnamed, named)
			}
		}
	}
	if len(unnamed) > 0 {
		return red("fixture-closure", "Writes: entry names a fixture-pinned path without naming the fixture: "+strings.Join(unnamed, ", "))
	}
	return green("fixture-closure")
}

// registryClosureCheck grades the declared binding into the ticket. A ticket that
// writes a bound package and omits a bound registry finds that registry mid-build
// and pays a repair round.
func registryClosureCheck(f Facts) CheckResult {
	var omitted []string
	seen := map[string]bool{}
	for _, ticket := range f.Tickets {
		owned := ownedPaths(ticket)
		for _, entry := range ticket.Writes {
			for _, file := range f.WritesBoundFiles[entry] {
				if pathCovered(file, owned) {
					continue
				}
				named := ticket.Name + ": " + entry + " requires " + file
				if seen[named] {
					continue
				}
				seen[named] = true
				omitted = append(omitted, named)
			}
		}
	}
	if len(omitted) > 0 {
		return red("registry-closure", "Writes: entry names a bound package without naming every bound file: "+strings.Join(omitted, ", "))
	}
	return green("registry-closure")
}

// kitPinCheck grades the kit pin into the ticket. A system-tagged test file reads
// BENCH_KIT, so an ambient value flips the fixture verdict under composition
// unless the ticket states the variable.
func kitPinCheck(f Facts) CheckResult {
	var unpinned []string
	for _, ticket := range f.Tickets {
		if f.TicketPinsKit[ticket.Name] {
			continue
		}
		for _, entry := range ticket.Writes {
			if f.WritesSystemTagged[entry] {
				unpinned = append(unpinned, ticket.Name+": "+entry)
			}
		}
	}
	if len(unpinned) > 0 {
		return red("kit-pin", "ticket writes a system-tagged test file without stating BENCH_KIT: "+strings.Join(unpinned, ", "))
	}
	return green("kit-pin")
}

// coversTokens is every row ID the parsed tickets cite, in enumeration order.
// It is the one ownership evidence both row checks read.
func coversTokens(f Facts) []string {
	var tokens []string
	for _, ticket := range f.Tickets {
		tokens = append(tokens, ticket.Covers...)
	}
	return tokens
}

func rowsOwnedCheck(f Facts) CheckResult {
	cited := coversTokens(f)
	var uncited []string
	for _, id := range f.DeclaredRowIDs {
		if !containsStr(cited, id) {
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
	for _, tok := range coversTokens(f) {
		if tickets.TagOf(tok) != f.SpecTag {
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
