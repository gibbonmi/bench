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

	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/preflight/evidencecmd"
	"github.com/gibbonmi/bench/internal/preflight/staleness"
	"github.com/gibbonmi/bench/internal/tickets"
)

// Facts is the immutable evidence Decide classifies. Every field is
// gathered once, up front, by Gather, and never re-read mid-verdict. So a
// rerun over the same Facts value always answers the same way.
type Facts struct {
	Mode      string
	Operation evidencecmd.Kind

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

	// ChangedPaths is the complete changed-file set since the resolved review base.
	ChangedPaths []string
	// DeletedPaths is the explicit committed range's exact D-status subset.
	DeletedPaths []string
	// DeletedBefore reports, for each (deleted) `Writes:` entry, whether its path is a
	// tracked file at the review base. The gatherer owns the probe; writes-resolve grades it.
	DeletedBefore map[string]bool
	// BuildDeletedPaths is the exact D-status subset of the build's landing range, from the
	// landing base to SourceTip. It is empty when the review record at SourceTip does not
	// read or holds no chunk, and when the landing range does not resolve.
	BuildDeletedPaths []string
	// DefaultTipPaths is the ChangedPaths subset whose bytes and mode at SourceTip equal
	// those at the default-branch tip. Such a path came in unchanged through a fold of the
	// default branch, so it is no write of the build.
	DefaultTipPaths []string
	// FoldResolutionPaths is every path that a fold entry of the review record at SourceTip
	// lists in its resolution. Each path keeps its bytes and mode from the newest such fold.
	// The fold review grades such a path, so no fence names it.
	FoldResolutionPaths []string

	// DeclaredRowIDs is the spec's coverage map row IDs, in map order.
	DeclaredRowIDs []string

	// ChunkTable is the acceptance column of the spec's Implementation chunks table: the
	// listed row IDs, the token faults, and the detail of an absent table or column.
	// chunk-rows grades it against DeclaredRowIDs.
	ChunkTable coverage.ChunkTable

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

	// CompletionPlanDigest is the digest the checkpoint's plan reader derived from
	// the spec's bench-completion-plan fence at SourceTip, and CompletionPlanError
	// is that reader's message, empty when the plan parses. The gatherer reads the
	// fence once; whether an unreadable plan is a fault is Decide's. The review
	// preparation's completion metadata records the same pair, so prepared evidence
	// and a verdict can never disagree about the plan.
	CompletionPlanDigest string
	CompletionPlanError  string

	// WritesPathExists reports, for each `Writes:` entry declared across the parsed
	// tickets, whether the path it names resolves in the tree. The gatherer owns the
	// probe; whether an absent path is a fault is Decide's.
	WritesPathExists map[string]bool

	// LandedDeletions maps each plain `Writes:` entry that the checkout does not hold to the
	// full identity of the default-branch commit that deleted its path, when the
	// default-branch tip lacks the path too. Only build mode fills it. writes-landed grades it.
	LandedDeletions map[string]string

	// MapSourceDiagnostics is the ordered list of Sources faults of every candidate decision
	// map at SourceTip, graded against the published tree and the predicted delivery closure
	// of SpecPath. Only build mode fills it. map-sources grades it.
	MapSourceDiagnostics []string

	// Audit is the staleness collection of the explicit pair. Only the audit form gathers it.
	// staleness.Decide grades it, and Decide reads none of it.
	Audit staleness.AuditFacts

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

	// BinarySealPresent reports whether the root publishes a dist/bench at all.
	// A root that publishes none is not graded, because a linked consumer repo
	// builds no Bench binary of its own.
	BinarySealPresent bool

	// BinarySealRefusal is the seal verifier's refusal for that published
	// binary, nil when the seal matches the sources. The refusal already
	// carries the reason and the one rebuild sentence, so the row renders
	// them rather than composing a second copy.
	BinarySealRefusal *freshness.Refusal

	// BinarySealRebuild is the owner's rebuild remedy and root kind for that root.
	BinarySealRebuild freshness.Rebuild
}

// CheckResult is one verdict row: the check's name, its verdict ("green" or
// "red"), a detail string, and the remedy that answers it. The detail is empty
// for green, and names the offending path/ID(s) for red. Next is empty on every
// row but the one red a single command repairs, because a remedy on a row that
// no command answers sends the reader down a false path.
type CheckResult struct {
	Check, Verdict, Detail, Next string
	// exempt marks a row whose own applicability rule can exclude it from this mode or
	// source. A row that a missing input holds not-applicable is not exempt.
	exempt bool
}

const (
	verdictGreen = "green"
	verdictRed   = "red"
	verdictNA    = "not-applicable"
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
// Four rows are conditional. tip-current appears only when --source-tip
// pinned a tip, directly after base-current, so the two halves of the
// source identity are graded together, ahead of every check that
// presupposes that identity. writes-landed appears only in build mode,
// directly after writes-resolve, because only the build entry tells a
// `Writes:` entry that landed work deleted from a typing error.
// binary-seal appears only in build mode, directly after kit-pin, because
// it guards the transaction a build is about to open. map-sources appears
// only in build mode, directly after binary-seal, and reads no ticket, so
// it is never not applicable: a decision-map source that the published tree
// does not keep stops the build before its first chunk.
func Decide(f Facts) Verdict {
	checks := []CheckResult{baseCurrentCheck(f)}
	if f.PinnedSourceTip != "" {
		checks = append(checks, tipCurrentCheck(f))
	}
	checks = append(checks,
		pathsAuthorizedCheck(f),
		ticketRow(f, "tickets-parse", ticketsParseCheck),
		ticketRow(f, "completion-plan", completionPlanCheck),
		ticketRow(f, "blockers-resolve", blockersResolveCheck),
		ticketRow(f, "writes-resolve", writesResolveCheck),
	)
	if f.Mode == modeBuild {
		checks = append(checks, ticketRow(f, "writes-landed", writesLandedCheck))
	}
	checks = append(checks, ticketRow(f, "kit-pin", kitPinCheck))
	if f.Mode == modeBuild {
		// A root that publishes no binary is outside the seal's own scope.
		checks = append(checks, exempt(binarySealCheck(f)), mapSourcesCheck(f))
	}
	checks = append(checks,
		ticketRow(f, "rows-owned", rowsOwnedCheck),
		ticketRow(f, "rows-membership", rowsMembershipCheck),
		ticketRow(f, chunkRowsName, chunkRowsCheck),
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

// exempt returns row with its CheckResult.exempt flag set.
func exempt(row CheckResult) CheckResult {
	row.exempt = true
	return row
}

// requiredGreen reports whether every row of the verdict is green or is an exempt
// not-applicable row; CheckResult.exempt defines exempt.
func (v Verdict) requiredGreen() bool {
	for _, c := range v.Checks {
		if c.Verdict != verdictGreen && (c.Verdict != verdictNA || !c.exempt) {
			return false
		}
	}
	return true
}

// ticketRow gates one ticket-reading check by build-mode ticket-directory
// applicability: not-applicable when build mode has no tickets/ directory at
// all, real otherwise. Every row that reads a parsed ticket shares this one
// gate, so no two of them can drift on when a fresh build is graded. The
// completion plan names tickets, so completion-plan shares the gate too.
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
		return exempt(notApplicable("diff-nonempty"))
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

// mapSourcesCheck reds each Sources fault of the candidate decision maps at the source tip.
func mapSourcesCheck(f Facts) CheckResult {
	if len(f.MapSourceDiagnostics) > 0 {
		return red("map-sources", strings.Join(f.MapSourceDiagnostics, "; "))
	}
	return green("map-sources")
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
