package preflight

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/preflight/ownership"
	"github.com/gibbonmi/bench/internal/preflight/staleness"
	specref "github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/tickets"
)

// BootstrapFailure is the fail-closed bootstrap answer. An artifact
// preflight cannot load or parse never becomes a green-by-omission; it
// becomes exactly one structured error naming what failed.
//
// Named failures include the missing or dangling spec path, and the
// coverage validator's message or the row-ID opt-in hint. They also
// include the found (non-staged) status. Every case fail-closed.
type BootstrapFailure struct {
	Kind, Hint string
}

func gather(root, mode, slug string, source *diff.SourceRange, sourcePaths []string, sourceTipPin string, noFollow, audit bool) (Facts, *BootstrapFailure) {
	pinnedTip, pinErr := resolvePin(root, sourceTipPin)
	if pinErr != nil {
		return Facts{}, pinErr
	}
	content, resolved, tried, ok, err := resolveSpecInput(root, slug, noFollow)
	if err != nil {
		return Facts{}, &BootstrapFailure{"spec not readable", "spec " + slug + ": " + err.Error()}
	}
	if !ok {
		return Facts{}, &BootstrapFailure{"spec not found", "no spec resolved for " + slug + " (tried: " + strings.Join(tried, ", ") + ")"}
	}

	status, statusErr := specStatus(root, slug, resolved)
	if statusErr != nil {
		return Facts{}, statusErr
	}
	if status != "staged" {
		return Facts{}, &BootstrapFailure{"spec not staged", "spec " + slug + " has Status: " + status + " (want staged)"}
	}

	optIn, ids, violations, cerr := coverage.ParseSpec(resolved)
	if cerr != nil {
		return Facts{}, &BootstrapFailure{"coverage map not readable", cerr.Error()}
	}
	if !optIn {
		return Facts{}, &BootstrapFailure{"coverage map not opted into row IDs", "add a leading `row` column to the map header"}
	}
	if len(violations) > 0 {
		return Facts{}, &BootstrapFailure{"coverage map invalid", strings.Join(violations, "; ")}
	}

	defaultBranch := ""
	defaultBranchResolved, defaultBranchCurrent := false, false
	reviewBase, reviewBaseResolved, reviewBaseHint := "", false, ""
	changedPaths := append([]string(nil), sourcePaths...)
	var resolvedSource diff.SourceRange
	var landingBase string
	if source != nil {
		resolvedSource = *source
		reviewBase, reviewBaseResolved = resolvedSource.Base, true
		landingBase, _, _ = reviewBaseFacts(root)
	} else {
		defaultBranch, defaultBranchResolved, defaultBranchCurrent = baseCurrentFacts(root)
		resolvedSource.Tip = headTip(root)
		reviewBase, reviewBaseResolved, reviewBaseHint = reviewBaseFacts(root)
		landingBase = reviewBase
		if reviewBaseResolved {
			changedPaths, err = diff.ChangedFilePathsAt(root, reviewBase)
			if err != nil {
				return Facts{}, &BootstrapFailure{"changed files not readable", err.Error()}
			}
		}
	}
	// Build mode grades the checkout tickets, because an uncommitted ticket edit counts there.
	// Review mode grades a clean committed source, so it grades the tickets that the source
	// tip tracks, and a checkout file that the tip does not hold authorizes nothing.
	specPath := filepath.ToSlash(specref.RelTo(root, resolved))
	var tracked func() ([]tickets.Entry, []string, error)
	if mode == modeReview {
		tracked = func() ([]tickets.Entry, []string, error) {
			return staleness.TicketEntriesAt(root, resolvedSource.Tip, specPath)
		}
	}
	ticketFacts, ticketErr := gatherTicketFolder(root, filepath.Join(filepath.Dir(resolved), "tickets"), mode, specTag(ids), noFollow, tracked)
	if ticketErr != nil {
		return Facts{}, ticketErr
	}
	defaultTipPaths, defaultTipFailure := defaultTipEqualPaths(root, resolvedSource.Tip, changedPaths)
	if defaultTipFailure != nil {
		return Facts{}, defaultTipFailure
	}
	sealPresent, sealRefusal := binarySealFacts(root)
	var audited staleness.AuditFacts
	if audit {
		source := ownership.Source{Base: resolvedSource.Base, Tip: resolvedSource.Tip, SpecPath: specPath}
		var auditFailure *BootstrapFailure
		if audited, auditFailure = collectAudit(root, content, source, ticketFacts.parsed); auditFailure != nil {
			return Facts{}, auditFailure
		}
	}

	return withFoldResolution(root, withCompletionPlan(root, Facts{
		Mode:                  mode,
		SpecPath:              specPath,
		DefaultBranch:         defaultBranch,
		DefaultBranchResolved: defaultBranchResolved,
		DefaultBranchCurrent:  defaultBranchCurrent,
		AssignmentTarget:      assignmentTarget(root),
		ReviewBaseResolved:    reviewBaseResolved,
		ReviewBaseHint:        reviewBaseHint,
		SourceBase:            resolvedSource.Base,
		SourceTip:             resolvedSource.Tip,
		PinnedSourceTip:       pinnedTip,
		ExplicitSourceRange:   source != nil,
		ChangedPaths:          changedPaths,
		DeletedPaths:          resolvedSource.DeletedPaths,
		DeletedBefore:         deletedBefore(root, reviewBase, ticketFacts.parsed),
		BuildDeletedPaths:     buildDeletedPaths(root, specPath, landingBase, resolvedSource.Tip),
		DefaultTipPaths:       defaultTipPaths,
		DeclaredRowIDs:        ids,
		ChunkTable:            coverage.ReadChunkTable(content),
		SpecTag:               specTag(ids),
		Tickets:               ticketFacts.parsed,
		TicketDiagnostics:     ticketFacts.diagnostics,
		BlockerCycles:         ticketFacts.cycles,
		WritesPathExists:      ticketFacts.writes,
		LandedDeletions:       landedDeletions(root, mode, ticketFacts.parsed, ticketFacts.writes),
		MapSourceDiagnostics:  mapSourceDiagnostics(root, mode, resolvedSource.Tip, specPath),
		Audit:                 audited,
		WritesSystemTagged:    ticketFacts.systemTag,
		TicketPinsKit:         ticketFacts.kitPinned,
		TicketsDirExists:      ticketFacts.dirExists,
		BinarySealPresent:     sealPresent,
		BinarySealRefusal:     sealRefusal,
	})), nil
}

// collectAudit reads the premise snapshot of the spec content against the parsed tickets,
// then the staleness facts of the gathered source pair. No verdict row reads the snapshot,
// so only the audit form collects it. A refused collection refuses the gather, so no
// partial audit reaches the verdict.
func collectAudit(root string, content []byte, source ownership.Source, parsed []tickets.Ticket) (staleness.AuditFacts, *BootstrapFailure) {
	snapshot := ownership.Collect(root, source, ownership.ParseManifest(content), parsed)
	audited, err := staleness.Collect(root, source, snapshot)
	var refusal *staleness.Error
	if errors.As(err, &refusal) {
		return staleness.AuditFacts{}, &BootstrapFailure{refusal.Stage + " refused", refusal.Error()}
	}
	if err != nil {
		return staleness.AuditFacts{}, &BootstrapFailure{"audit not readable", err.Error()}
	}
	return audited, nil
}

// resolvePin resolves a --source-tip value to its full commit identity. An
// empty pin is the flag's absence, not a failure. An unresolvable one is
// refused in the same shape ResolveSourceRange refuses an unreachable
// --base, so a typo never reaches the verdict table as a drift.
func resolvePin(root, pin string) (string, *BootstrapFailure) {
	if pin == "" {
		return "", nil
	}
	resolved, err := git.ResolveCommit(root, pin)
	if err != nil {
		return "", &BootstrapFailure{"cannot resolve --source-tip", "'" + pin + "' does not name a commit reachable in this repository"}
	}
	return resolved, nil
}

// headTip is the derived source tip of a bare invocation, the counterpart
// to the snapshot head an explicit source range captures. An unreadable
// HEAD answers empty. tip-current renders that as its own red rather than
// a bootstrap failure. A bare invocation with no pin has no use for the
// value at all.
func headTip(root string) string {
	tip, err := git.Output("-C", root, "rev-parse", "HEAD")
	if err != nil {
		return ""
	}
	return tip
}

// UnauthorizedPathsError is the fence refusal AuthorizeReviewedSource returns when the
// reviewed range holds paths no ownership fence covers. Detail is the report's own
// sentence, and Paths is the same set unjoined. A caller that renders a path table reads
// Paths, so it never splits a sentence back apart.
type UnauthorizedPathsError struct {
	Detail string
	Paths  []string
}

func (e UnauthorizedPathsError) Error() string { return e.Detail }

// AuthorizeReviewedSource gathers a new review snapshot for root, slug, and
// base. It returns the source range only when its committed paths pass the
// ownership fence that the staged spec's tickets derive. The review gather reads
// the tickets of the source tip from committed objects. It returns an error for
// gather failure, a red authorization row, or an unresolved range. A caller's
// earlier Facts or range cannot replace this final authorization.
func AuthorizeReviewedSource(root, slug, base string) (diff.SourceRange, error) {
	facts, failure := Gather(root, modeReview, slug, base)
	if failure != nil {
		return diff.SourceRange{}, fmt.Errorf("%s: %s", failure.Kind, failure.Hint)
	}
	if facts.SourceBase == "" || facts.SourceTip == "" {
		return diff.SourceRange{}, errors.New("reviewed source range is unresolved")
	}
	check := pathsAuthorizedCheck(facts)
	if check.Verdict == verdictRed {
		if unauthorized := unauthorizedPaths(facts); len(unauthorized) > 0 {
			return diff.SourceRange{}, UnauthorizedPathsError{Detail: check.Detail, Paths: unauthorized}
		}
		return diff.SourceRange{}, errors.New(check.Detail)
	}
	return diff.SourceRange{Base: facts.SourceBase, Tip: facts.SourceTip, CommittedPaths: facts.ChangedPaths, DeletedPaths: facts.DeletedPaths}, nil
}

// specStatus resolves the typed Status: value for slug via the spec
// package's Facts, the one source of typed spec status. It matches by
// slug rather than re-parsing the content this package already holds.
func specStatus(root, slug, resolved string) (string, *BootstrapFailure) {
	facts, err := specref.Facts(root)
	if err != nil {
		return "", &BootstrapFailure{"spec status not readable", err.Error()}
	}
	for _, f := range facts {
		if f.Slug == slug {
			return f.Status, nil
		}
	}
	// resolved may point at a path Facts' folder-spec enumeration does not
	// cover (e.g. a non-standard argument). Either way, no typed status is
	// available to trust.
	return "", &BootstrapFailure{"spec status not readable", resolved + " did not resolve through folder-spec enumeration"}
}

// specTag is the alphabetic prefix shared by a spec's declared row IDs, e.g. "PF" for
// PF1..n. A spec with no declared IDs (already refused upstream by the coverage
// bootstrap check) has no tag; callers past bootstrap never observe that case.
func specTag(ids []string) string {
	if len(ids) == 0 {
		return ""
	}
	return tickets.TagOf(ids[0])
}

// ticketFacts is everything the gatherer learns from specs/<slug>/tickets/.
// Decide grades these values and reads no file of its own.
type ticketFacts struct {
	parsed      []tickets.Ticket
	diagnostics []string
	cycles      []string
	writes      map[string]bool
	systemTag   map[string]bool
	kitPinned   map[string]bool
	dirExists   bool
}

// gatherTickets enumerates specs/<slug>/tickets/ and parses every `.md`
// file through the ticket grammar, recursing into subdirectories so a
// ticket filed under tickets/sub/ is graded the same as one at the top
// level. Every entry — file or subdirectory, at every depth — is
// lstat-classified before it is opened or descended into. So a FIFO, a
// dangling symlink, or another special file is refused rather than
// blocking, no matter how deep it sits.
//
// Review mode requires the top-level directory to exist. Build mode
// instead reports whether it exists at all. This lets the verdict core
// tell an absent directory (row checks not-applicable) from a
// present-but-empty one (row checks run for real, reading as unowned
// rows).
func gatherTicketsWithPolicy(root, dir, mode, tag string, noFollow bool) (ticketFacts, *BootstrapFailure) {
	return gatherTicketFolder(root, dir, mode, tag, noFollow, nil)
}

// gatherTicketFolder is gatherTicketsWithPolicy with an optional tracked read. The checkout
// walk keeps its folder policy and its refusals, and a non-nil tracked read then supplies the
// ticket files that the facts grade. A failed tracked read refuses the gather, and no
// checkout ticket replaces it.
func gatherTicketFolder(root, dir, mode, tag string, noFollow bool, tracked func() ([]tickets.Entry, []string, error)) (ticketFacts, *BootstrapFailure) {
	var d bounds.ClassifiedDir
	if noFollow {
		d = bounds.ClassifyDirNoFollow(dir)
	} else {
		d = bounds.ClassifyDir(dir)
	}
	switch d.State {
	case bounds.StateAbsent:
		if mode == "review" {
			return ticketFacts{}, &BootstrapFailure{"tickets directory absent", dir + " does not exist"}
		}
		return ticketFacts{}, nil
	case bounds.StateEmpty:
		return ticketFacts{dirExists: true}, nil
	case bounds.StateParsed:
		// fall through to enumeration below
	default:
		return ticketFacts{}, &BootstrapFailure{"tickets directory not readable", dir + " is " + string(d.State) + ": " + d.Reason}
	}

	var files []tickets.Entry
	var duplicates []string
	var refusal *tickets.Refusal
	if noFollow {
		files, duplicates, refusal = tickets.EnumerateNoFollow(dir, d.Entries)
	} else {
		files, duplicates, refusal = tickets.Enumerate(dir, d.Entries)
	}
	if refusal == nil && tracked != nil {
		var err error
		if files, duplicates, err = tracked(); err != nil && !errors.As(err, &refusal) {
			return ticketFacts{}, &BootstrapFailure{"tickets not readable at the source tip", err.Error()}
		}
	}
	if refusal != nil {
		return ticketFacts{}, &BootstrapFailure{refusal.Kind, refusal.Message(dir)}
	}
	facts, gradeErr := gradeTickets(root, files, duplicates, tag)
	if gradeErr != nil {
		return ticketFacts{}, gradeErr
	}
	facts.dirExists = true
	return facts, nil
}

// gradeTickets parses each enumerated ticket against its siblings and the
// spec tag, and probes the tree for every declared `Writes:` entry. The
// duplicate-identity diagnostics the enumeration reports open the diagnostic
// list, so a duplicate basename is named before the grammar faults below it.
func gradeTickets(root string, files []tickets.Entry, duplicates []string, tag string) (ticketFacts, *BootstrapFailure) {
	facts := ticketFacts{
		writes:    map[string]bool{},
		systemTag: map[string]bool{},
		kitPinned: map[string]bool{},
	}
	facts.diagnostics = append(facts.diagnostics, duplicates...)

	siblings := make([]string, 0, len(files))
	for _, file := range files {
		siblings = append(siblings, file.Name)
	}

	for _, file := range files {
		parsed, diagnostics := tickets.ParseTicket(file.Name, file.Data, siblings, tag)
		if field, value, unrepresentable := tickets.UnrepresentableValue(parsed); unrepresentable {
			return ticketFacts{}, &BootstrapFailure{"ticket path not representable", fmt.Sprintf("%s declares a %s: entry %q with a byte spec-TOON cannot represent", file.Rel, field, value)}
		}
		for _, diagnostic := range diagnostics {
			facts.diagnostics = append(facts.diagnostics, file.Rel+": "+diagnostic)
		}
		facts.parsed = append(facts.parsed, parsed)
		facts.kitPinned[file.Name] = bytes.Contains(file.Data, []byte(kitEnvMarker))
		for _, entry := range parsed.Writes {
			facts.probe(root, entry)
		}
	}
	facts.cycles = tickets.Cycles(facts.parsed)
	return facts, nil
}

// assignmentTarget is the id of the active assignment that owns this preflight
// root, empty otherwise. intent.AssignmentForWorktree owns the match, so the
// assignment this remedy names and the section `bench handoff` rewrites resolve
// from one lookup. An unowned tree answers empty, and the remedy then prints its
// placeholder.
func assignmentTarget(root string) string {
	assignment, owned := intent.AssignmentForWorktree(root)
	if !owned {
		return ""
	}
	return assignment.ID
}

// reviewBaseFacts wraps the exported diff-base resolution, the single
// source `bench diff` itself consumes. So preflight and `bench diff` can
// never disagree about the base.
func reviewBaseFacts(root string) (base string, resolved bool, hint string) {
	base, _, errKind, errHint := diff.ResolveReviewBase(root)
	if errKind != "" {
		return "", false, errKind + ": " + errHint
	}
	return base, true, ""
}
