package worktree

import (
	"errors"
	"fmt"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	refreshop "github.com/gibbonmi/bench/internal/refresh"
	"github.com/gibbonmi/bench/internal/toon"
	"io"
	"time"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/usage"
)

// ErrTreeTargetPath is the outcome of a tree-target value that names no label and that the
// worktree target grammar reads as a path. A tree target is never a path, so the caller
// answers this outcome as a grammar refusal.
var ErrTreeTargetPath = errors.New("tree target is a path")

// TreeTarget answers the worktree of the one assignment whose label is exactly label. The
// lookup takes no id, no prefix, and no path, so a tree target is always a name the
// operator chose.
//
// One active match wins over any inactive match of the same label. With no active match,
// one inactive match gives its state refusal. Two or more matches of one class give the
// ambiguity refusal. The label lookup runs before the path-shape test, so a label with a
// separator still names its worktree. The selected assignment then passes the state,
// missing-tree, and creation-bundle checks of every other target-taking verb.
func TreeTarget(root, label string) (string, error) {
	if !lineSafe(label) {
		return "", errTargetControls
	}
	assignments, err := intent.Assignments(root)
	if err != nil {
		return "", err
	}
	labeled := matchingAssignments(assignments, func(a intent.Assignment) bool { return a.Label == label })
	if active := matchingAssignments(labeled, func(a intent.Assignment) bool { return landingActiveState(a.State) }); len(active) > 0 {
		labeled = active
	}
	switch {
	case len(labeled) == 0 && pathShaped(label):
		return "", ErrTreeTargetPath
	case len(labeled) == 0:
		return "", errTargetUnassigned
	case len(labeled) > 1:
		return "", ambiguousAssignments(labeled)
	}
	// The id is the one address that names exactly this record, so the shared resolver
	// applies its checks to the record that the label selected.
	selected, err := resolveAssignmentIn(root, labeled[0].ID, landingActiveState)
	if err != nil {
		return "", err
	}
	return selected.Worktree, nil
}

// pathShaped reports whether targetPath, the owner of the worktree target grammar, reads
// value as a path. A shape that the grammar refuses, such as `~user` or a relative path,
// is a path shape too.
func pathShaped(value string) bool {
	_, isPath, err := targetPath(value)
	return isPath || err != nil
}

// PrintTreeTargetRefusal prints err as the worktree target refusal of verb and answers
// exit 1.
func PrintTreeTargetRefusal(stderr io.Writer, verb string, err error) int {
	return printTargetRefusal(stderr, verb, err)
}

// PrintTreeBuildRefusal prints detail as the refusal of verb for the kit worktree target
// label, whose own build cannot start, and answers exit 1. The repair is the build verb on
// label, with the label as one shell word. Neither line names the executable path.
func PrintTreeBuildRefusal(stderr io.Writer, verb, label, detail string) int {
	next := usage.WorktreeBuildFor(axi.ShellQuote(label))
	return printTargetRefusal(stderr, verb, refusalError{refusal{detail: detail, next: next}})
}

var createGrammar = usage.Grammar{
	Cmd:  "bench worktree create",
	Help: "usage: " + usage.WorktreeCreate,
	Flags: []usage.Flag{
		{Name: "--request", HasValue: true, NoEmptyValue: true, Required: true},
		{Name: "--label", HasValue: true, NoEmptyValue: true, Required: true},
		{Name: "--refresh", HasValue: false},
		{Name: "--from", HasValue: true, NoEmptyValue: true},
	},
}

// createSiblingStart resolves `--from` through the one sibling lookup the merge verb also
// composes. The flag reaches no commit lookup, so a spelling that names no active
// assignment is a refusal rather than a fallthrough to the default tip. The caller runs
// fromRepresentable before the creation starts, so this lookup reads a value a line can
// carry.
func createSiblingStart(root, from string) (creationStart, error) {
	assignments, err := intent.Assignments(root)
	if err != nil {
		return creationStart{}, err
	}
	tip, id, ok, err := siblingTip(root, assignments, "", from)
	if err != nil {
		return creationStart{}, err
	}
	if !ok {
		return creationStart{}, refusalError{refusal{detail: "--from names no active assignment", observed: from}}
	}
	binding, err := (commitrepo.Store{Root: root}).Inheritance(id)
	return creationStart{ref: tip, binding: binding}, err
}

// CreateCommand owns the worktree create grammar and creates or replays one owned
// assignment after parsing succeeds. Grammar answers perform no creation or tracing;
// --from selects an active sibling tip and cannot be combined with --refresh.
func CreateCommand(root, home string, args []string, stdout, stderr io.Writer) int {
	parsed, line, code := usage.Parse(createGrammar, args)
	if line != "" {
		if code == 0 {
			fmt.Fprintln(stdout, line)
			return 0
		}
		fmt.Fprintln(stderr, line)
		return code
	}
	// The seam record opens once the grammar has answered, the way the commit boundary
	// opens once the repository is known: a grammar answer creates nothing to record.
	var assignment string
	finishSpan := beginVerbSpan(home, root, otelCreateSeam)
	exit := createAttributed(&assignment, parsed, root, home, currentTime(), args, stdout, stderr)
	finishSpan(exit, assignment)
	return exit
}

// createAttributed is the create verb's own work at the entry's instant, with the
// assignment the record names written to assignment once the creation resolves it.
func createAttributed(assignment *string, parsed usage.Result, root, home string, now time.Time, args []string, stdout, stderr io.Writer) int {
	from := parsed.Flags["--from"]
	// The two flags name two starts, so the pair refuses before the refresh runs: a fetch
	// that moved the default branch would already have taken effect by the refusal.
	if _, refresh := parsed.Flags["--refresh"]; refresh && from != "" {
		fmt.Fprintln(stderr, toon.Usage(createGrammar.Cmd, "--from with --refresh"))
		return 2
	}
	// The value's representability is a grammar fact, so it refuses before the creation
	// reads anything at all.
	if from != "" {
		if err := fromRepresentable(from); err != nil {
			return printTargetRefusal(stderr, createGrammar.Cmd, err)
		}
	}
	_, startRef := refreshop.Consume(root, args, stdout)
	// The sibling lookup is deferred, because the creation resolves the request replay
	// first. A replay returns its existing record, so the sibling's state gates nothing a
	// no-op run would act on.
	var fromErr error
	resolveStart := func() (creationStart, error) { return creationStart{ref: startRef}, nil }
	if from != "" {
		resolveStart = func() (creationStart, error) {
			tip, err := createSiblingStart(root, from)
			fromErr = err
			return tip, err
		}
	}
	request, label := parsed.Flags["--request"], parsed.Flags["--label"]
	creation, err := createAt(defaultJoins(), root, home, request, label, nil, now, resolveStart)
	// A failed creation returns the zero Creation, so its empty ID leaves the record unnamed.
	*assignment = creation.Assignment.ID
	if err != nil {
		if fromErr != nil {
			return printTargetRefusal(stderr, createGrammar.Cmd, err)
		}
		fmt.Fprintf(stderr, "bench worktree create: %v\n", err)
		return 1
	}
	out, err := toon.Table(createTable, []string{"path", "assignment", "state"}, [][]string{{creation.Path, creation.Assignment.ID, string(creation.Assignment.State)}})
	if err != nil {
		fmt.Fprintf(stderr, "bench worktree create: %v\n", err)
		return 1
	}
	fmt.Fprint(stdout, out)
	fmt.Fprintf(stdout, "next[2]:\n  bench worktree exec \"%s\" -- <command>\n  bench worktree path \"%s\"\n", label, label)
	return 0
}
