package preflight

import (
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
func Gather(root, mode, slug string) (Facts, *BootstrapFailure) {
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

	defaultBranchResolved, defaultBranchCurrent := baseCurrentFacts(root)
	reviewBase, reviewBaseResolved, reviewBaseHint := reviewBaseFacts(root)
	var changedPaths []string
	if reviewBaseResolved {
		changedPaths, err = changedFilePaths(reviewBase)
		if err != nil {
			return Facts{}, &BootstrapFailure{"changed files not readable", err.Error()}
		}
	}

	return Facts{
		Mode:                  mode,
		SpecPath:              filepath.ToSlash(specref.RelTo(root, resolved)),
		DefaultBranchResolved: defaultBranchResolved,
		DefaultBranchCurrent:  defaultBranchCurrent,
		ReviewBaseResolved:    reviewBaseResolved,
		ReviewBaseHint:        reviewBaseHint,
		ChangedPaths:          changedPaths,
		FenceEntries:          fenceEntries,
		DeclaredRowIDs:        ids,
		TicketTokens:          tokens,
		SpecTag:               specTag(ids),
		TicketsDirExists:      ticketsDirExists,
	}, nil
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
// authorization.
func fenceTokens(content []byte) []string {
	var tokens []string
	inSection := false
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
			tokens = append(tokens, fenceTokensInLine(line)...)
		}
	}
	return tokens
}

// fenceTokensInLine is a small state machine over one line: it tracks parenthesis
// depth and backtick state, capturing a backtick-quoted token only when the depth at
// the moment its opening backtick appeared was zero.
func fenceTokensInLine(line string) []string {
	var tokens []string
	depth := 0
	inTick := false
	depthAtOpen := 0
	var cur strings.Builder
	for _, r := range line {
		switch r {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case '`':
			if inTick {
				if depthAtOpen == 0 {
					tokens = append(tokens, cur.String())
				}
				cur.Reset()
				inTick = false
			} else {
				inTick = true
				depthAtOpen = depth
			}
		default:
			if inTick {
				cur.WriteRune(r)
			}
		}
	}
	return tokens
}

// gatherTicketTokens enumerates specs/<slug>/tickets/, lstat-classifying every entry
// before it is opened so a FIFO or other special file is refused rather than blocking.
// Review mode requires the directory to exist; build mode instead reports whether it
// exists at all (the second return value) so the verdict core can tell an absent
// directory (row checks not-applicable) from a present-but-empty one (row checks run
// for real and read as unowned rows).
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

	for _, entry := range d.Entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		c := bounds.Classify(path, bounds.ControlRecordLimit)
		switch c.State {
		case bounds.StateParsed:
			tokens = append(tokens, tokenRe.FindAllString(string(c.Data), -1)...)
		case bounds.StateEmpty:
			// nothing to scan
		default:
			return nil, false, &BootstrapFailure{"ticket file not readable", path + " is " + string(c.State) + ": " + c.Reason}
		}
	}
	return tokens, true, nil
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

// changedFilePaths mirrors bench diff's exact changed-file semantics — `git diff
// --name-status --no-renames -z <base>`, which folds in committed, index, and tracked
// worktree changes since base — without importing the diff package's unexported
// parser. Only the path half of each row is needed here.
func changedFilePaths(base string) ([]string, error) {
	raw, err := git.Raw("diff", "--name-status", "--no-renames", "-z", base)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(string(raw), "\x00")
	var paths []string
	for i := 0; i+1 < len(parts); i += 2 {
		status, path := parts[i], parts[i+1]
		if status == "" && path == "" {
			continue
		}
		paths = append(paths, path)
	}
	return paths, nil
}
