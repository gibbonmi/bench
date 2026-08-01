package gate

// The full gate's enforcement of the reduced-run declaration. Every outer run grades
// its excludable phases against a materialized worktree of the current tree with the
// declared allowlist absent, so a phase whose excludability is a lie fails loudly on an
// ordinary gate at the moment the dependency appears — instead of passing in the dark
// and letting a later reduced run inherit evidence nothing produced.
//
// The subject is materialized through git — a throwaway index snapshots the working
// tree, the declared paths are dropped from it, and the resulting commit is checked out
// as a detached worktree — because the excludable phases include contract tests that
// stage their subject through `git ls-files`: a plain directory copy is the obvious
// construction and it fails that staging outright. Registering a real git worktree also
// hands orphan retirement to the existing worktree machinery; the only cleanup here is
// the run's own deferred removal.

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/canary"
	benchgit "github.com/gibbonmi/bench/internal/git"
)

// runOuterPhases is the outer gate's entry: it splits the resolved table into the
// phases that must see the real root and the declared-excludable ones, which run only
// against the stripped subject. Three shapes run unsplit — every phase against the real
// root: a table with nothing excludable, a root that is not a git work tree (the
// declaration's paths are repository-relative, so outside a repository there is no
// membership to enforce), and a root that is not the kit. The declaration this
// construction enforces is the kit's own — ReducedScope() is compiled from the kit's
// source, and only the kit's root may inherit against it (the reduction guard in
// reducedInheritance rules the same way) — so a linked repo pays no stripped worktree
// for an allowlist it never declared, and the identity check is by file, not spelling,
// with any stat failure running unsplit. Those are the deliberately open cases; a
// materialization failure inside the kit's own repository is red, because silently
// skipping the construction would de-enforce the declaration on exactly the run that
// needed it.
func runOuterPhases(ctx context.Context, root, kit string, phases []Phase, stdout, stderr io.Writer) int {
	scope := ReducedScope()
	primary := make([]Phase, 0, len(phases))
	excludable := false
	for _, phase := range phases {
		if scope.Excludable(phase.Name) {
			excludable = true
			continue
		}
		primary = append(primary, phase)
	}
	if !excludable || !insideGitWorkTree(root) || !sameDirectory(root, kit) {
		return runPhases(ctx, kit, phases, outerMode, stdout, stderr)
	}
	subjectRoot, cleanup, err := materializeStrippedSubject(root, scope)
	if err != nil {
		fmt.Fprintf(stderr, "gate: cannot materialize the stripped subject: %v\n", err)
		fmt.Fprintln(stderr, "gate: red")
		return 1
	}
	defer cleanup()
	stripped, err := strippedPhaseSet(subjectRoot, kit, scope)
	if err != nil {
		fmt.Fprintln(stderr, err)
		fmt.Fprintln(stderr, "gate: red")
		return 1
	}
	return runSplitPhases(ctx, kit, primary, stripped, subjectRoot, stdout, stderr)
}

// strippedPhaseSet resolves the phase table a second time against the stripped subject
// and keeps the excludable phases plus the build phase. Re-resolving is what keeps one
// source for every phase definition — the same table constructor, pointed at the other
// root — and the build phase rides along because it produces the dist/ binary the
// excludable phases exec, exactly as it does on the primary root.
func strippedPhaseSet(subjectRoot, kit string, scope Scope) ([]Phase, error) {
	table, err := phaseTable(subjectRoot, kit)
	if err != nil {
		return nil, err
	}
	stripped := make([]Phase, 0, len(table))
	for _, phase := range table {
		if scope.Excludable(phase.Name) || phase.Name == canary.PhaseBuild {
			stripped = append(stripped, phase)
		}
	}
	return stripped, nil
}

// runSplitPhases schedules both sets concurrently and merges one verdict. The stripped
// set writes to its own capability-skip log, kept apart from the primary log so the two
// postures cannot blur: the primary run keeps the dev tier's informational reading.
func runSplitPhases(ctx context.Context, root string, primary, stripped []Phase, subjectRoot string, stdout, stderr io.Writer) int {
	primaryLog, primaryCleanup, err := newSkipLog()
	if err != nil {
		fmt.Fprintf(stderr, "gate: cannot open the capability skip log: %v\n", err)
		fmt.Fprintln(stderr, "gate: red")
		return 1
	}
	defer primaryCleanup()
	strippedLog, strippedCleanup, err := newSkipLog()
	if err != nil {
		fmt.Fprintf(stderr, "gate: cannot open the stripped-subject skip log: %v\n", err)
		fmt.Fprintln(stderr, "gate: red")
		return 1
	}
	defer strippedCleanup()

	names := make([]string, 0, len(stripped))
	for _, phase := range stripped {
		names = append(names, phase.Name)
	}
	// A skipped grading surface that says nothing reads as a gate that never ran, so
	// the construction announces itself and what it grades.
	fmt.Fprintf(stdout, "gate: stripped subject %s grades: %s\n", subjectRoot, strings.Join(names, ", "))

	open := prefixedPhaseWriters(stdout, stderr)
	type outcome struct {
		results   []phaseResult
		cancelled bool
	}
	primaryDone := make(chan outcome, 1)
	go func() {
		results, cancelled := schedule(ctx, root, withSkipLog(primary, primaryLog), false, open)
		primaryDone <- outcome{results, cancelled}
	}()
	strippedResults, strippedCancelled := schedule(ctx, root, withSkipLog(stripped, strippedLog), false, open)
	primaryOutcome := <-primaryDone

	return aggregateAndReport(append(primaryOutcome.results, strippedResults...),
		primaryOutcome.cancelled || strippedCancelled, stdout, stderr,
		func() bool { return reportCapabilitySkips(primaryLog, stdout, stderr) },
		func() bool { return reportStrippedSkips(strippedLog, stderr) })
}

// reportStrippedSkips is the posture half of the enforcement. It reds on one signature:
// a kind=environment skip whose reason names a path this construction stripped. Three
// doors have to stay shut here, and this predicate has already been wrong through two
// of them.
//
// Too narrow — it cannot be "simplified" into reusing strict mode. Strict mode counts
// only kind=capability, and the kit's idiom for a missing subject file,
// skipIfSubjectFileMissing, degrades absence into a structured skip of kind=environment
// and exits green. A strict-mode tally here would let a phase whose subject file this
// construction removed go permanently, silently green: the failure that killed this
// design's first draft.
//
// Too broad by kind — it cannot be widened to red on any kind. kind=capability skips
// are host limitations, a filesystem without FIFOs or a run that cannot drop privilege,
// which stripping cannot induce and which this repository's own suite emits on ordinary
// developer hosts. They stay informational here exactly as the dev tier leaves them.
//
// Too broad by reason — the environment kind alone is not the signature either. The
// excludable phases emit environment skips for an absent non-declared subject file, an
// unmaterialized fixture, a release plan with no target for this host, an unset
// conformance root; a real run of this repository produced 146 of them, none related to
// stripping. Counting the kind reds every full gate. Only a reason naming a declared
// path can have been induced here, because the declaration is the only thing this
// construction removed.
//
// The membership half is structural rather than textual: the reason's whitespace-
// separated fields go through Scope.Member, the single source of what "declared" means,
// so this predicate never encodes the wording of any skip site and follows the
// declaration when it changes. The textual half is only the tokenizing, which the wire
// format forces — reason is free text. That also makes the predicate independent of
// skipIfSubjectFileMissing's phrasing: any future degradation that names the path it
// could not find is caught by the same rule.
//
// An unreadable log is red for the same reason it is under strict mode — a tally that
// proves nothing reads exactly like a clean run.
func reportStrippedSkips(path string, stderr io.Writer) bool {
	tally, err := readSkipTally(path)
	if err != nil {
		fmt.Fprintf(stderr, "gate: the stripped-subject skip log %s is unreadable, so the excludable phases prove nothing: %v\n", path, err)
		return true
	}
	induced := strippingInducedSkips(tally, ReducedScope())
	if len(induced) == 0 {
		return false
	}
	fmt.Fprintf(stderr, "gate: %d environment skip(s) are fatal against the stripped subject: an excludable phase degraded on a declared path this run stripped instead of running:\n  %s\n",
		len(induced), strings.Join(induced, "\n  "))
	return true
}

// strippingInducedSkips returns the environment skip reasons that name a declared path.
// Fields are trimmed of the punctuation a reason wraps a path in — quotes and a trailing
// clause separator — before membership is asked, so a path quoted or followed by a colon
// still reads as the path it is.
func strippingInducedSkips(tally skipTally, scope Scope) []string {
	var induced []string
	for _, reason := range tally.environmentReasons {
		for _, field := range strings.Fields(reason) {
			if scope.Member(strings.Trim(field, `"'`+"`"+`,;:`)) {
				induced = append(induced, reason)
				break
			}
		}
	}
	return induced
}

// materializeStrippedSubject builds a detached git worktree holding the current working
// tree with the declared paths absent, and returns its path with a cleanup that retires
// it. Concurrent runs get distinct temporary roots. The commit's ident and dates are
// fixed so an unchanged stripped tree maps to one object across runs rather than
// accreting a dangling commit per gate.
func materializeStrippedSubject(root string, scope Scope) (string, func(), error) {
	tree := benchgit.TreeHash(root)
	if tree == "" || tree == "none" {
		return "", nil, fmt.Errorf("cannot hash the working tree of %s", root)
	}
	indexDir, err := os.MkdirTemp("", "bench-stripped-index-")
	if err != nil {
		return "", nil, err
	}
	defer os.RemoveAll(indexDir)
	index := []string{"GIT_INDEX_FILE=" + filepath.Join(indexDir, "index")}
	if _, err := gitAt(root, index, "read-tree", tree); err != nil {
		return "", nil, err
	}
	// --ignore-unmatch: a declared path absent from the tree is simply not there to
	// strip, which is the fixture and fresh-clone case, not a defect.
	strip := append([]string{"rm", "-r", "-q", "--cached", "--ignore-unmatch", "--"}, scopePathspecs(scope)...)
	if _, err := gitAt(root, index, strip...); err != nil {
		return "", nil, err
	}
	strippedTree, err := gitAt(root, index, "write-tree")
	if err != nil {
		return "", nil, err
	}
	ident := []string{
		"GIT_AUTHOR_NAME=bench-gate", "GIT_AUTHOR_EMAIL=gate@bench.invalid",
		"GIT_COMMITTER_NAME=bench-gate", "GIT_COMMITTER_EMAIL=gate@bench.invalid",
		"GIT_AUTHOR_DATE=1970-01-01T00:00:00Z", "GIT_COMMITTER_DATE=1970-01-01T00:00:00Z",
	}
	commit, err := gitAt(root, ident, "commit-tree", strippedTree, "-m", "bench stripped gate subject")
	if err != nil {
		return "", nil, err
	}
	parent, err := os.MkdirTemp("", "bench-stripped-subject-")
	if err != nil {
		return "", nil, err
	}
	subjectRoot := filepath.Join(parent, "subject")
	if _, err := gitAt(root, nil, "worktree", "add", "-q", "--detach", subjectRoot, commit); err != nil {
		_ = os.RemoveAll(parent)
		return "", nil, err
	}
	cleanup := func() {
		_, _ = gitAt(root, nil, "worktree", "remove", "--force", subjectRoot)
		_ = os.RemoveAll(parent)
	}
	return subjectRoot, cleanup, nil
}

// scopePathspecs renders the declaration as git pathspecs: files verbatim, directories
// without the membership-marking trailing slash.
func scopePathspecs(scope Scope) []string {
	specs := scope.Files()
	for _, dir := range scope.Directories() {
		specs = append(specs, strings.TrimSuffix(dir, "/"))
	}
	return specs
}

func gitAt(root string, env []string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	cmd.Env = append(os.Environ(), env...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %v: %s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}

func insideGitWorkTree(root string) bool {
	out, err := benchgit.Output("-C", root, "rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}
