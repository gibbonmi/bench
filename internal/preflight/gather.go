package preflight

import (
	"bytes"
	"errors"
	"fmt"
	"go/build/constraint"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/canonicalpath"
	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	specref "github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/tickets"
)

// BootstrapFailure is the fail-closed bootstrap answer. An artifact
// preflight cannot load or parse never becomes a green-by-omission; it
// becomes exactly one structured error naming what failed.
//
// Named failures include the missing or dangling spec path, and the
// coverage validator's message or the row-ID opt-in hint. They also
// include the empty fences section, and the found (non-staged) status.
// Every case fail-closed.
type BootstrapFailure struct {
	Kind, Hint string
}

// Gather gathers preflight facts for root, mode, slug, and the first optional
// explicit base. It supplies no source-tip pin. It returns one immutable
// Facts snapshot for Decide or one BootstrapFailure; it never classifies a
// check. Exactly one result is non-zero.
func Gather(root, mode, slug string, explicitBase ...string) (Facts, *BootstrapFailure) {
	base := ""
	if len(explicitBase) > 0 {
		base = explicitBase[0]
	}
	return GatherPinned(root, mode, slug, base, "")
}

// GatherPinned gathers preflight facts for root, mode, slug, an explicit
// base, and a source-tip pin. With an explicit base, it reads one
// movement-checked source snapshot and retries one snapshot drift. It returns
// snapshot drift after a second movement, or one BootstrapFailure for a read
// or pin failure. A resolved pin that names the wrong commit is a Decide row.
func GatherPinned(root, mode, slug, explicitBase, sourceTipPin string) (Facts, *BootstrapFailure) {
	if explicitBase != "" {
		var gathered Facts
		var gatherFailure *BootstrapFailure
		result := diff.MovementCheckedRetry(root, func(snapshot diff.MovementSnapshot) (string, string) {
			var err error
			var resolveKind, resolveHint string
			source, resolveKind, resolveHint := snapshot.ResolveSourceRange(explicitBase)
			if resolveKind != "" {
				return resolveKind, resolveHint
			}
			paths, err := snapshot.SourceSnapshotPaths(source)
			if err != nil {
				return "changed files not readable", err.Error()
			}
			if mode == "review" {
				dirty, statusErr := git.Output("-C", root, "status", "--porcelain")
				if statusErr != nil {
					return "source status unreadable", statusErr.Error()
				}
				if dirty != "" {
					return "source not clean", "review source has uncommitted changes"
				}
			}
			gathered, gatherFailure = gather(root, mode, slug, &source, paths, sourceTipPin)
			if gatherFailure != nil {
				return gatherFailure.Kind, gatherFailure.Hint
			}
			return "", ""
		})
		if result.Kind != "" {
			if gatherFailure != nil {
				return Facts{}, gatherFailure
			}
			return Facts{}, &BootstrapFailure{result.Kind, result.Hint}
		}
		if result.DriftKind != "" {
			return Facts{}, &BootstrapFailure{"snapshot drift", result.DriftHint}
		}
		return gathered, nil
	}
	return gather(root, mode, slug, nil, nil, sourceTipPin)
}

func gather(root, mode, slug string, source *diff.SourceRange, sourcePaths []string, sourceTipPin string) (Facts, *BootstrapFailure) {
	pinnedTip, pinErr := resolvePin(root, sourceTipPin)
	if pinErr != nil {
		return Facts{}, pinErr
	}
	content, resolved, tried, ok, err := specref.Resolve(root, slug)
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

	fenceEntries, _ := specref.FenceTokens(content)
	if len(fenceEntries) == 0 {
		return Facts{}, &BootstrapFailure{"ownership fences empty", "## Ownership fences declares no backticked entry outside parentheses"}
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

	ticketFacts, ticketErr := gatherTickets(root, filepath.Join(filepath.Dir(resolved), "tickets"), mode, specTag(ids))
	if ticketErr != nil {
		return Facts{}, ticketErr
	}

	defaultBranch := ""
	defaultBranchResolved, defaultBranchCurrent := false, false
	reviewBase, reviewBaseResolved, reviewBaseHint := "", false, ""
	changedPaths := append([]string(nil), sourcePaths...)
	var resolvedSource diff.SourceRange
	if source != nil {
		resolvedSource = *source
		reviewBase, reviewBaseResolved = resolvedSource.Base, true
	} else {
		defaultBranch, defaultBranchResolved, defaultBranchCurrent = baseCurrentFacts(root)
		resolvedSource.Tip = headTip(root)
		reviewBase, reviewBaseResolved, reviewBaseHint = reviewBaseFacts(root)
		if reviewBaseResolved {
			changedPaths, err = diff.ChangedFilePathsAt(root, reviewBase)
			if err != nil {
				return Facts{}, &BootstrapFailure{"changed files not readable", err.Error()}
			}
		}
	}

	sealPresent, sealRefusal := binarySealFacts(root)

	return Facts{
		Mode:                  mode,
		SpecPath:              filepath.ToSlash(specref.RelTo(root, resolved)),
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
		FenceEntries:          fenceEntries,
		DeclaredRowIDs:        ids,
		SpecTag:               specTag(ids),
		Tickets:               ticketFacts.parsed,
		TicketDiagnostics:     ticketFacts.diagnostics,
		BlockerCycles:         ticketFacts.cycles,
		WritesPathExists:      ticketFacts.writes,
		WritesFixturePins:     ticketFacts.pins,
		WritesBoundFiles:      ticketFacts.bound,
		WritesSystemTagged:    ticketFacts.systemTag,
		TicketPinsKit:         ticketFacts.kitPinned,
		TicketsDirExists:      ticketFacts.dirExists,
		BinarySealPresent:     sealPresent,
		BinarySealRefusal:     sealRefusal,
	}, nil
}

// destinationBinary is the published Bench executable a hand run, a hook, the
// wrapper, and the landing all execute. It is the one path the seal row grades.
const destinationBinary = "dist/bench"

// binarySealFacts grades root's published binary through the seal verifier,
// the one primitive that decides staleness by source digest. An absent binary
// is reported as absent rather than graded, because a linked consumer repo
// publishes none. The verifier's refusal is carried whole, so the row prints
// the rebuild sentence the verifier composed rather than a second copy.
func binarySealFacts(root string) (present bool, refusal string) {
	executable := filepath.Join(root, filepath.FromSlash(destinationBinary))
	if _, err := os.Stat(executable); err != nil {
		return false, ""
	}
	if err := freshness.Verify(root, executable); err != nil {
		return true, err.Error()
	}
	return true, ""
}

// resolvePin resolves a --source-tip value to its full commit identity. An
// empty pin is the flag's absence, not a failure. An unresolvable one is
// refused in the same shape ResolveSourceRange refuses an unreachable
// --base, so a typo never reaches the verdict table as a drift.
func resolvePin(root, pin string) (string, *BootstrapFailure) {
	if pin == "" {
		return "", nil
	}
	resolved, err := git.Output("-C", root, "rev-parse", "--verify", pin+"^{commit}")
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
// staged spec's ownership fences. It returns an error for gather failure, a
// red authorization row, or an unresolved range. A caller's earlier Facts or
// range cannot replace this final authorization.
func AuthorizeReviewedSource(root, slug, base string) (diff.SourceRange, error) {
	facts, failure := Gather(root, "review", slug, base)
	if failure != nil {
		return diff.SourceRange{}, fmt.Errorf("%s: %s", failure.Kind, failure.Hint)
	}
	check := pathsAuthorizedCheck(facts)
	if check.Verdict == verdictRed {
		if unauthorized := unauthorizedPaths(facts); len(unauthorized) > 0 {
			return diff.SourceRange{}, UnauthorizedPathsError{Detail: check.Detail, Paths: unauthorized}
		}
		return diff.SourceRange{}, errors.New(check.Detail)
	}
	if facts.SourceBase == "" || facts.SourceTip == "" {
		return diff.SourceRange{}, errors.New("reviewed source range is unresolved")
	}
	return diff.SourceRange{Base: facts.SourceBase, Tip: facts.SourceTip, CommittedPaths: facts.ChangedPaths}, nil
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
	pins        map[string][]string
	bound       map[string][]string
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
func gatherTickets(root, dir, mode, tag string) (ticketFacts, *BootstrapFailure) {
	d := bounds.ClassifyDir(dir)
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

	files, duplicates, refusal := tickets.Enumerate(dir, d.Entries)
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
// probe is the gatherer's whole I/O contribution to the ownership row;
// the policy over the resulting bit belongs to Decide. The duplicate-identity
// diagnostics the enumeration reports open the diagnostic list, so a duplicate
// basename is named before the grammar faults below it.
func gradeTickets(root string, files []tickets.Entry, duplicates []string, tag string) (ticketFacts, *BootstrapFailure) {
	pins, pinErr := canary.FixturePins(root)
	if pinErr != nil {
		return ticketFacts{}, &BootstrapFailure{"fixture inventory not readable", pinErr.Error()}
	}
	facts := ticketFacts{
		writes:    map[string]bool{},
		pins:      map[string][]string{},
		bound:     map[string][]string{},
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
			path, _ := splitWritesEntry(entry)
			facts.writes[entry] = treeHolds(root, path)
			if pinning := pins[path]; len(pinning) > 0 {
				facts.pins[entry] = pinning
			}
			if bound := tickets.BoundFiles(path); len(bound) > 0 {
				facts.bound[entry] = bound
			}
			if systemTagged(root, path) {
				facts.systemTag[entry] = true
			}
		}
	}
	facts.cycles = tickets.Cycles(facts.parsed)
	return facts, nil
}

// treeHolds reports whether one `Writes:` path names something in the tree.
// Any file type answers yes: the row grades whether the path exists, not
// what sits at it. The lstat never follows a link, so a dangling symlink
// reads as the absent path it is.
func treeHolds(root, path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Lstat(filepath.Join(root, filepath.FromSlash(path)))
	return err == nil
}

// kitEnvMarker is the literal a ticket states to pin the kit a system-tagged
// fixture reads. An ambient value flips such a fixture under composition, so the
// ticket names the variable rather than inheriting whatever the session holds.
const kitEnvMarker = "BENCH_KIT"

// systemTagged reports whether one `Writes:` path names a Go test file whose
// //go:build expression carries the system tag. The constraint parse is the same
// one the coverage citations grade with, so the two never disagree about which
// file is system-tagged. A path that is no test file, that does not read, or that
// carries no constraint answers no.
func systemTagged(root, path string) bool {
	if !strings.HasSuffix(path, "_test.go") {
		return false
	}
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "//") {
			return false
		}
		if !constraint.IsGoBuild(trimmed) {
			continue
		}
		expr, parseErr := constraint.Parse(trimmed)
		return parseErr == nil && exprHoldsTag(expr, "system")
	}
	return false
}

// exprHoldsTag reports whether a build expression mentions tag anywhere. A
// negated mention still counts: the file's verdict depends on the tag either way.
func exprHoldsTag(expr constraint.Expr, tag string) bool {
	switch typed := expr.(type) {
	case *constraint.TagExpr:
		return typed.Tag == tag
	case *constraint.NotExpr:
		return exprHoldsTag(typed.X, tag)
	case *constraint.AndExpr:
		return exprHoldsTag(typed.X, tag) || exprHoldsTag(typed.Y, tag)
	case *constraint.OrExpr:
		return exprHoldsTag(typed.X, tag) || exprHoldsTag(typed.Y, tag)
	}
	return false
}

// relTo is path expressed against base, falling back to the full path when
// the two share no relation. A diagnostic then still names something real.
func relTo(base, path string) string {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return path
	}
	return rel
}

// canonicalRoot is a path's absolute, symlink-resolved, cleaned form — the one
// identity two spellings of the same tree share. internal/worktree derives the
// same form for its own targets, but that package imports this one, so a shared
// call would close an import cycle.
func canonicalRoot(path string) (string, error) {
	return canonicalpath.Resolve(path)
}

// baseCurrentFacts backs the base-current check: it resolves the default
// branch and reports whether its tip is an ancestor of HEAD:
// merge-base(default, HEAD) equal to rev-parse(default).
//
// The resolved name comes back with the two predicates, because the stale-base
// remedy names that branch and git.ResolvedDefault is its one source.
//
// An unresolved default branch answers no name and both predicates false. The
// check itself renders that as red without a separate bootstrap failure, since
// map #7 names this a per-check red rather than a bootstrap precondition.
func baseCurrentFacts(root string) (branch string, resolved, current bool) {
	def, ok := git.ResolvedDefault(root)
	if !ok {
		return "", false, false
	}
	mergeBase, err1 := git.Output("merge-base", def, "HEAD")
	tip, err2 := git.Output("rev-parse", def)
	if err1 != nil || err2 != nil {
		return def, true, false
	}
	return def, true, mergeBase == tip
}

// assignmentTarget is the id of the active assignment that owns this preflight
// root, empty otherwise. The tree is matched by canonical path, because the
// ledger records a resolved path while a root may arrive through a symlink. An
// unreadable ledger answers empty: the remedy then prints its placeholder, which
// is the same answer a root outside the pool gets.
func assignmentTarget(root string) string {
	canonical, err := canonicalRoot(root)
	if err != nil {
		return ""
	}
	assignments, err := intent.Assignments(root)
	if err != nil {
		return ""
	}
	for _, a := range assignments {
		if a.State != intent.StateActive {
			continue
		}
		if owned, ownedErr := canonicalRoot(a.Worktree); ownedErr == nil && owned == canonical {
			return a.ID
		}
	}
	return ""
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
