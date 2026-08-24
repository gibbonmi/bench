// Movement types and review-range resolution for package diff.
package diff

import (
	"sort"

	"github.com/gibbonmi/bench/internal/git"
)

// diffRange is the resolved review range every rendering section shares: the base
// and method preamble lines, and the argument shapes changedFiles/commitLog/diffBody
// need for either the branch-relative (resolved base through worktree) range or the
// commit-relative (exact first-parent) range.
type diffRange struct {
	base      string
	head      string
	method    string
	filesArgs []string
	logRange  string
	bodyArgs  []string
}

// SourceRange is the frozen, inclusive-ancestor source identity shared by explicit
// diff and preflight. Tip is captured with Base so callers never pair a path set
// from one source revision with the identity of another.
type SourceRange struct {
	Base, Tip      string
	CommittedPaths []string
}

// MovementSnapshot is the typed checkout state one movement-checked read exposes to
// its caller. The root and facts stay coupled so derived source ranges and paths
// cannot accidentally inspect the process working directory.
type MovementSnapshot struct {
	root  string
	Facts git.DiffFacts
}

// MovementResult is the single outcome of a root-bound attempt. Drift remains
// distinct from a read failure so callers retry only repository movement.
type MovementResult struct {
	DriftKind, DriftHint string
	Kind, Hint           string
}

// MovementChecked captures one root-bound checkout identity, runs read against its
// typed facts, then verifies the identity again. A non-empty drift asks the caller to
// retry; kind and hint preserve an ordinary read failure without reclassifying it.
func MovementChecked(root string, read func(MovementSnapshot) (kind, hint string)) MovementResult {
	facts, err := git.AllFilesFacts(root)
	if err != nil {
		return MovementResult{Kind: "checkout facts failed", Hint: err.Error()}
	}
	before, err := capturedIdentity(root, facts)
	if err != nil {
		return MovementResult{Kind: "snapshot identity failed", Hint: err.Error()}
	}
	if kind, hint := read(MovementSnapshot{root: root, Facts: facts}); kind != "" {
		return MovementResult{Kind: kind, Hint: hint}
	}
	snapshotAfterRead()
	after, err := recapturedIdentity(root, facts)
	if err != nil {
		return MovementResult{Kind: "snapshot identity failed", Hint: err.Error()}
	}
	if drift := before.drifted(after); drift != "" {
		return MovementResult{DriftKind: drift, DriftHint: "the " + drift + " changed while reading; retry the exact invocation"}
	}
	return MovementResult{}
}

// MovementCheckedRetry is the movement-retry policy every drift-sensitive read shares:
// one retry when the only thing that failed was repository movement, then the second
// attempt's answer stands, including its terminal drift hint, which no caller spells
// for itself. A read failure is never retried, so a caller cannot turn a broken read
// into a drift refusal.
func MovementCheckedRetry(root string, read func(MovementSnapshot) (kind, hint string)) MovementResult {
	var result MovementResult
	for attempt := 0; attempt < 2; attempt++ {
		result = MovementChecked(root, read)
		if result.Kind != "" || result.DriftKind == "" {
			return result
		}
	}
	return result
}

// ResolveSourceRange resolves base against the snapshot's captured source tip.
func (snapshot MovementSnapshot) ResolveSourceRange(base string) (SourceRange, string, string) {
	return ResolveSourceRange(snapshot.root, base, snapshot.Facts.Head)
}

// SourceSnapshotPaths returns source paths derived from the same captured checkout.
func (snapshot MovementSnapshot) SourceSnapshotPaths(source SourceRange) ([]string, error) {
	return sourceSnapshotPaths(snapshot.root, source, snapshot.Facts)
}

// ResolveSourceRange resolves an explicit source base and its exact source tip.
func ResolveSourceRange(root, base, tip string) (SourceRange, string, string) {
	b, err := git.Output("-C", root, "rev-parse", "--verify", base+"^{commit}")
	if err != nil {
		return SourceRange{}, "cannot resolve --base", "'" + base + "' does not name a commit reachable in this repository"
	}
	if !git.OK("-C", root, "merge-base", "--is-ancestor", b, tip) {
		return SourceRange{}, "--base is not an ancestor", "'" + b + "' is not an ancestor of '" + tip + "'"
	}
	paths, err := changedFilesAt(root, b, tip)
	if err != nil {
		return SourceRange{}, "committed source paths not readable", err.Error()
	}
	committed := make([]string, 0, len(paths))
	for _, path := range paths {
		committed = append(committed, path[1])
	}
	return SourceRange{Base: b, Tip: tip, CommittedPaths: committed}, "", ""
}

// resolveCommitRange builds the diffRange for `--commit <sha>`: base is <sha>'s
// resolved first parent. The sha is verified before anything renders. An unresolvable
// sha and a root commit's missing parent are each their own structured error (kind,
// hint), never a leaked git failure.
func resolveCommitRange(root, commitArg string) (dr diffRange, errKind, errHint string) {
	headSha, err := git.Output("-C", root, "rev-parse", "--verify", commitArg+"^{commit}")
	if err != nil {
		return diffRange{}, "cannot resolve --commit",
			"'" + commitArg + "' does not name a commit reachable in this repository"
	}
	baseSha, err := git.Output("-C", root, "rev-parse", "--verify", commitArg+"^")
	if err != nil {
		return diffRange{}, "--commit has no parent",
			"'" + commitArg + "' is a root commit — there is no first parent to diff against"
	}
	return diffRange{
		base:      baseSha,
		head:      headSha,
		method:    "commit " + headSha,
		filesArgs: []string{baseSha, headSha},
		logRange:  baseSha + ".." + headSha,
		bodyArgs:  []string{baseSha, headSha},
	}, "", ""
}

// resolveBranchRange builds the diffRange for bare `bench diff`/`--full`: the
// recorded-key base when it names a reachable ancestor, else merge-base with the
// default branch, byte-identical to the pre-`--commit` behavior.
func resolveBranchRange(root string) (dr diffRange, errKind, errHint string) {
	base, method, errKind, errHint := ResolveReviewBase(root)
	if errKind != "" {
		return diffRange{}, errKind, errHint
	}
	return diffRange{
		base:      base,
		method:    method,
		filesArgs: []string{base},
		logRange:  base + "..HEAD",
		bodyArgs:  []string{base},
	}, "", ""
}

func resolveExplicitRange(root, base, head string) (diffRange, string, string) {
	source, kind, hint := ResolveSourceRange(root, base, head)
	if kind != "" {
		return diffRange{}, kind, hint
	}
	return diffRange{base: source.Base, head: source.Tip, method: "explicit", filesArgs: []string{source.Base}, logRange: source.Base + ".." + source.Tip, bodyArgs: []string{source.Base}}, "", ""
}

// resolvePairRange builds the immutable diffRange for the explicit --base/--source-tip
// pair — the same pair `bench preflight review` takes. The tip is resolved before the
// base so its refusal names the flag the caller got wrong.
func resolvePairRange(root, base, tip string) (diffRange, string, string) {
	resolvedTip, err := git.Output("-C", root, "rev-parse", "--verify", tip+"^{commit}")
	if err != nil {
		return diffRange{}, "cannot resolve --source-tip", "'" + tip + "' does not name a commit reachable in this repository"
	}
	source, kind, hint := ResolveSourceRange(root, base, resolvedTip)
	if kind != "" {
		return diffRange{}, kind, hint
	}
	return diffRange{base: source.Base, head: source.Tip, method: "explicit pair", filesArgs: []string{source.Base, source.Tip}, logRange: source.Base + ".." + source.Tip, bodyArgs: []string{source.Base, source.Tip}}, "", ""
}

// SourceSnapshotPaths returns the complete explicit-source inventory: committed,
// index, tracked-worktree, and untracked paths. The committed half belongs to SourceRange.
func SourceSnapshotPaths(root string, source SourceRange) ([]string, error) {
	facts, err := git.AllFilesFacts(root)
	if err != nil {
		return nil, err
	}
	return sourceSnapshotPaths(root, source, facts)
}

func sourceSnapshotPaths(root string, source SourceRange, facts git.DiffFacts) ([]string, error) {
	paths := append([]string(nil), source.CommittedPaths...)
	tracked, err := ChangedFilePathsAt(root, source.Base)
	if err != nil {
		return nil, err
	}
	paths = append(paths, tracked...)
	for _, entry := range facts.Changes {
		if entry.Status == "??" {
			paths = append(paths, entry.Path)
		}
	}
	sort.Strings(paths)
	unique := paths[:0]
	for _, path := range paths {
		if len(unique) == 0 || unique[len(unique)-1] != path {
			unique = append(unique, path)
		}
	}
	return unique, nil
}

func resolveBranchRangeFromFacts(root string, facts git.DiffFacts) (dr diffRange, errKind, errHint string) {
	base, method := "", ""
	if facts.RecordedBase != "" {
		switch {
		case !git.OK("-C", root, "cat-file", "-e", facts.RecordedBase+"^{commit}"):
			method = "merge-base (recorded sha unreachable)"
		case !git.OK("-C", root, "merge-base", "--is-ancestor", facts.RecordedBase, facts.Head):
			method = "merge-base (recorded sha not an ancestor)"
		default:
			base, method = facts.RecordedBase, "recorded"
		}
	}
	if base == "" {
		if !facts.DefaultResolved {
			return diffRange{}, "cannot resolve a review base", "this repository has no resolvable default branch; record one with: git config branch.<name>.benchBase <sha>"
		}
		var err error
		base, err = git.Output("-C", root, "merge-base", facts.DefaultTip, facts.Head)
		if err != nil {
			return diffRange{}, "cannot resolve a review base", "no merge-base with '" + facts.DefaultBranch + "'; record one with: git config branch.<name>.benchBase <sha>"
		}
		if method == "" {
			method = "merge-base"
		}
	}
	return diffRange{
		base:      base,
		head:      facts.Head,
		method:    method,
		filesArgs: []string{base},
		logRange:  base + ".." + facts.Head,
		bodyArgs:  []string{base},
	}, "", ""
}

// ResolveReviewBase is the single source of the resolved review base for
// bench diff and its consumers: the recorded `branch.<name>.benchBase` key when
// it names a reachable ancestor of HEAD, else merge-base with the resolved
// default branch. method carries which path answered: `recorded`, `merge-base`, or one
// of the loud fallback labels when a recorded key is present but unusable. A non-empty
// errKind/errHint with an empty base is the only absence shape; base is never empty on
// a nil error.
func ResolveReviewBase(root string) (base, method, errKind, errHint string) {
	base, method = resolveBase()
	if base != "" {
		return base, method, "", ""
	}
	def, ok := git.ResolvedDefault(root)
	if !ok {
		return "", "", "cannot resolve a review base",
			"this repository has no resolvable default branch; record one with: git config branch.<name>.benchBase <sha>"
	}
	mb, err := git.Output("merge-base", def, "HEAD")
	if err != nil {
		return "", "", "cannot resolve a review base",
			"no merge-base with '" + def + "'; record one with: git config branch.<name>.benchBase <sha>"
	}
	if method == "" {
		method = "merge-base"
	}
	return mb, method, "", ""
}

// resolveBase returns the recorded-key base and `recorded` when the key names a
// reachable ancestor, or ("", <loud fallback method>) when the key is present but
// unreachable/divergent, or ("","") when there is no usable recorded key.
func resolveBase() (base, method string) {
	branch, _ := git.Output("symbolic-ref", "--quiet", "--short", "HEAD")
	if branch == "" {
		return "", ""
	}
	key, _ := git.Output("config", "branch."+branch+".benchBase")
	if key == "" {
		return "", ""
	}
	switch {
	case !git.OK("cat-file", "-e", key+"^{commit}"):
		return "", "merge-base (recorded sha unreachable)"
	case !git.OK("merge-base", "--is-ancestor", key, "HEAD"):
		return "", "merge-base (recorded sha not an ancestor)"
	default:
		return key, "recorded"
	}
}
