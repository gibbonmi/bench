package preflight

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/git"
	specref "github.com/gibbonmi/bench/internal/spec"
)

// BootstrapFailure is the fail-closed bootstrap answer: an artifact preflight cannot
// load or parse never becomes a green-by-omission, it becomes exactly one structured
// error naming what failed — the missing or dangling spec path, the coverage
// validator's own message or the row-ID opt-in hint, the empty fences section, or the
// found (non-staged) status — every case fail-closed.
type BootstrapFailure struct {
	Kind, Hint string
}

// tokenRe is the one grammar for a row-ID-shaped token in ticket-file prose: an
// uppercase tag plus digits, word-boundary matched so "PF1a" never mistakes itself
// for "PF1".
var tokenRe = regexp.MustCompile(`\b[A-Z]+[0-9]+\b`)

// fencesEndRe bounds the `## Ownership fences` section the same way coverage.go
// bounds `### Acceptance coverage map`: a level-2-or-deeper heading ends it. The
// section itself is opened by an exact `## Ownership fences` line match.
var fencesEndRe = regexp.MustCompile(`^#{2,} `)

// Gather is the thin gatherer: it reads git, the exported diff-base resolution, the
// spec resolver, the coverage parser, and tickets/ enumeration, and returns either an
// immutable Facts value ready for Decide, or the one BootstrapFailure that explains
// why no Facts value can be trusted. Exactly one of the two return values is non-zero.
func Gather(root, mode, slug string, explicitBase ...string) (Facts, *BootstrapFailure) {
	if len(explicitBase) > 0 && explicitBase[0] != "" {
		for attempt := 0; attempt < 2; attempt++ {
			var gathered Facts
			var gatherFailure *BootstrapFailure
			result := diff.MovementChecked(root, func(snapshot diff.MovementSnapshot) (string, string) {
				var err error
				var resolveKind, resolveHint string
				source, resolveKind, resolveHint := snapshot.ResolveSourceRange(explicitBase[0])
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
				gathered, gatherFailure = gather(root, mode, slug, &source, paths)
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
			if result.DriftKind == "" {
				return gathered, nil
			}
			if attempt == 1 {
				return Facts{}, &BootstrapFailure{"snapshot drift", result.DriftHint}
			}
		}
		return Facts{}, &BootstrapFailure{"snapshot drift", "the repository changed while reading; retry the exact invocation"}
	}
	return gather(root, mode, slug, nil, nil)
}

func gather(root, mode, slug string, source *diff.SourceRange, sourcePaths []string) (Facts, *BootstrapFailure) {
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

	fenceEntries := fenceTokens(content)
	if len(fenceEntries) == 0 {
		return Facts{}, &BootstrapFailure{"ownership fences empty", "## Ownership fences declares no backticked entry outside parentheses"}
	}

	tokens, ticketsDirExists, ticketErr := gatherTicketTokens(filepath.Join(filepath.Dir(resolved), "tickets"), mode)
	if ticketErr != nil {
		return Facts{}, ticketErr
	}

	defaultBranchResolved, defaultBranchCurrent := false, false
	reviewBase, reviewBaseResolved, reviewBaseHint := "", false, ""
	changedPaths := append([]string(nil), sourcePaths...)
	var resolvedSource diff.SourceRange
	if source != nil {
		resolvedSource = *source
		reviewBase, reviewBaseResolved = resolvedSource.Base, true
	} else {
		defaultBranchResolved, defaultBranchCurrent = baseCurrentFacts(root)
		reviewBase, reviewBaseResolved, reviewBaseHint = reviewBaseFacts(root)
		if reviewBaseResolved {
			changedPaths, err = diff.ChangedFilePathsAt(root, reviewBase)
			if err != nil {
				return Facts{}, &BootstrapFailure{"changed files not readable", err.Error()}
			}
		}
	}

	return Facts{
		Mode:                  mode,
		SpecPath:              filepath.ToSlash(specref.RelTo(root, resolved)),
		DefaultBranchResolved: defaultBranchResolved,
		DefaultBranchCurrent:  defaultBranchCurrent,
		ReviewBaseResolved:    reviewBaseResolved,
		ReviewBaseHint:        reviewBaseHint,
		SourceBase:            resolvedSource.Base,
		SourceTip:             resolvedSource.Tip,
		ExplicitSourceRange:   source != nil,
		ChangedPaths:          changedPaths,
		FenceEntries:          fenceEntries,
		DeclaredRowIDs:        ids,
		TicketTokens:          tokens,
		SpecTag:               specTag(ids),
		TicketsDirExists:      ticketsDirExists,
	}, nil
}

// AuthorizeReviewedSource returns the one shared range fact after checking its
// committed paths against the staged spec's existing ownership-fence owner.
// Landing consumes this narrower final authorization rather than re-parsing fences.
func AuthorizeReviewedSource(root, slug, base string) (diff.SourceRange, error) {
	facts, failure := Gather(root, "review", slug, base)
	if failure != nil {
		return diff.SourceRange{}, fmt.Errorf("%s: %s", failure.Kind, failure.Hint)
	}
	check := pathsAuthorizedCheck(facts)
	if check.Verdict == verdictRed {
		return diff.SourceRange{}, errors.New(check.Detail)
	}
	if facts.SourceBase == "" || facts.SourceTip == "" {
		return diff.SourceRange{}, errors.New("reviewed source range is unresolved")
	}
	return diff.SourceRange{Base: facts.SourceBase, Tip: facts.SourceTip, CommittedPaths: facts.ChangedPaths}, nil
}

// specStatus resolves the typed Status: value for slug via the spec package's Facts —
// the one source of typed spec status — matched by slug rather than re-parsing the
// content this package already holds.
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
	// resolved may point at a path Facts' specs/*/spec.md glob does not cover (e.g. a
	// non-standard argument); either way, no typed status is available to trust.
	return "", &BootstrapFailure{"spec status not readable", resolved + " did not resolve through the folder-spec glob"}
}

// specTag is the alphabetic prefix shared by a spec's declared row IDs, e.g. "PF" for
// PF1..n. A spec with no declared IDs (already refused upstream by the coverage
// bootstrap check) has no tag; callers past bootstrap never observe that case.
func specTag(ids []string) string {
	if len(ids) == 0 {
		return ""
	}
	return tagOf(ids[0])
}

// fenceTokens extracts every backticked token in the `## Ownership fences` section
// that is not inside parentheses — parenthetical prose is annotation, never
// authorization. Paren depth and backtick state carry across line boundaries: a
// parenthetical that opens on one line and closes on a later one still shields every
// token inside it, and depth returns to zero once it closes so a later real entry
// authorizes normally.
func fenceTokens(content []byte) []string {
	var tokens []string
	inSection := false
	depth := 0
	inTick := false
	depthAtOpen := 0
	var cur strings.Builder
	for _, raw := range strings.Split(string(content), "\n") {
		line := strings.TrimSuffix(raw, "\r")
		if strings.TrimSpace(line) == "## Ownership fences" {
			inSection = true
			continue
		}
		if inSection && fencesEndRe.MatchString(line) {
			inSection = false
		}
		if inSection {
			fenceTokensInLine(line, &depth, &inTick, &depthAtOpen, &cur, &tokens)
		}
	}
	return tokens
}

// fenceTokensInLine is one line's pass through the fence-section state machine:
// paren depth, backtick state, and the token under construction are threaded in by
// pointer so the caller can carry them across every line of the section. A
// backtick-quoted token is captured into tokens only when the depth at the moment its
// opening backtick appeared was zero — inside an open paren, whether opened on this
// line or an earlier one, never authorizes.
func fenceTokensInLine(line string, depth *int, inTick *bool, depthAtOpen *int, cur *strings.Builder, tokens *[]string) {
	for _, r := range line {
		switch r {
		case '(':
			*depth++
		case ')':
			if *depth > 0 {
				*depth--
			}
		case '`':
			if *inTick {
				if *depthAtOpen == 0 {
					*tokens = append(*tokens, cur.String())
				}
				cur.Reset()
				*inTick = false
			} else {
				*inTick = true
				*depthAtOpen = *depth
			}
		default:
			if *inTick {
				cur.WriteRune(r)
			}
		}
	}
}

// gatherTicketTokens enumerates specs/<slug>/tickets/, recursing into subdirectories
// so a token cited only under tickets/sub/ is found the same as one at the top level.
// Every entry — file or subdirectory, at every depth — is lstat-classified before it
// is opened or descended into, so a FIFO or other special file is refused rather than
// blocking no matter how deep it sits. Review mode requires the top-level directory to
// exist; build mode instead reports whether it exists at all (the second return value)
// so the verdict core can tell an absent directory (row checks not-applicable) from a
// present-but-empty one (row checks run for real and read as unowned rows).
func gatherTicketTokens(dir, mode string) (tokens []string, exists bool, err *BootstrapFailure) {
	d := bounds.ClassifyDir(dir)
	switch d.State {
	case bounds.StateAbsent:
		if mode == "review" {
			return nil, false, &BootstrapFailure{"tickets directory absent", dir + " does not exist"}
		}
		return nil, false, nil
	case bounds.StateEmpty:
		return nil, true, nil
	case bounds.StateParsed:
		// fall through to enumeration below
	default:
		return nil, false, &BootstrapFailure{"tickets directory not readable", dir + " is " + string(d.State) + ": " + d.Reason}
	}

	tokens, ticketErr := scanTicketEntries(dir, d.Entries)
	if ticketErr != nil {
		return nil, false, ticketErr
	}
	return tokens, true, nil
}

// scanTicketEntries walks one already-classified directory listing, scanning files
// for tokens and recursing into subdirectories with the same lstat-first
// classification gatherTicketTokens applies at the top level — so the special-file
// refusal holds at every depth, not only the first.
func scanTicketEntries(dir string, entries []fs.DirEntry) ([]string, *BootstrapFailure) {
	var tokens []string
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			sub := bounds.ClassifyDir(path)
			switch sub.State {
			case bounds.StateEmpty:
				// nothing to scan
			case bounds.StateParsed:
				subTokens, subErr := scanTicketEntries(path, sub.Entries)
				if subErr != nil {
					return nil, subErr
				}
				tokens = append(tokens, subTokens...)
			default:
				return nil, &BootstrapFailure{"tickets directory not readable", path + " is " + string(sub.State) + ": " + sub.Reason}
			}
			continue
		}
		c := bounds.Classify(path, bounds.ControlRecordLimit)
		switch c.State {
		case bounds.StateParsed:
			tokens = append(tokens, tokenRe.FindAllString(string(c.Data), -1)...)
		case bounds.StateEmpty:
			// nothing to scan
		default:
			return nil, &BootstrapFailure{"ticket file not readable", path + " is " + string(c.State) + ": " + c.Reason}
		}
	}
	return tokens, nil
}

// baseCurrentFacts backs the base-current check: it resolves the default branch and
// reports whether its tip is an ancestor of HEAD — merge-base(default, HEAD) equal to
// rev-parse(default). An unresolved default branch answers (false, false): the check
// itself renders that as red without a separate bootstrap failure, since map #7 names
// this a per-check red rather than a bootstrap precondition.
func baseCurrentFacts(root string) (resolved, current bool) {
	def, ok := git.ResolvedDefault(root)
	if !ok {
		return false, false
	}
	mergeBase, err1 := git.Output("merge-base", def, "HEAD")
	tip, err2 := git.Output("rev-parse", def)
	if err1 != nil || err2 != nil {
		return true, false
	}
	return true, mergeBase == tip
}

// reviewBaseFacts wraps the exported diff-base resolution — the single source `bench
// diff` itself consumes — so preflight and `bench diff` can never disagree about the
// base.
func reviewBaseFacts(root string) (base string, resolved bool, hint string) {
	base, _, errKind, errHint := diff.ResolveReviewBase(root)
	if errKind != "" {
		return "", false, errKind + ": " + errHint
	}
	return base, true, ""
}
