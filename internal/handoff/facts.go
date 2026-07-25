package handoff

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/status"
)

// facts is the pin block's field set, already reduced to the strings the renderer places.
// Each field is either a fact or an explicit statement that the fact is unavailable —
// never an empty value, which a cold reader cannot tell apart from a deliberate blank.
type facts struct {
	Repo     string // the repository's directory name
	Origin   string // the origin remote's URL, or "" when the repo has none
	Path     string // the git root as a reader on another machine should see it
	Branch   string // the branch, or a phrase naming why there is none
	Head     string // the short HEAD sha, or a phrase naming why there is none
	Dirty    string // the working-tree clause
	Unpushed string // the unpushed-commit clause
	Specs    []string
	Gate     string
	Action   string // the next command, derived or overridden
	Signal   string // the board signal Action came from; "" when Action was overridden

	// NoInvocable marks the board that had signals but no command among them. It is a
	// third state, distinct from an empty Action on a clean board: collapsing the two
	// would report a clean board to a session that has work waiting on it.
	NoInvocable bool
}

// unknown phrases name an absent fact in the reader's terms. They are values, not error
// paths: a degenerate repository is a state to report, not a reason to refuse.
const (
	unknownBranch   = "detached HEAD (no branch)"
	unknownHead     = "HEAD unknown (no commits yet)"
	unknownDirty    = "dirty state unknown"
	unknownUnpushed = "unpushed count unknown"
)

// collect gathers every pin fact under root. It never fails: each query that cannot answer
// degrades to its explicit unknown, because a session resuming into a half-initialized repo
// still needs the fields that do resolve.
func collect(root string) facts {
	f := facts{
		Repo:   filepath.Base(root),
		Path:   renderPath(root, os.Getenv("HOME")),
		Branch: branch(root),
		Head:   head(root),
		Specs:  liveSpecs(root),
		Gate:   gateField(status.GateVerdict(root)),
	}
	if origin, err := git.Output("-C", root, "remote", "get-url", "origin"); err == nil {
		f.Origin = origin
	}
	f.Dirty, f.Unpushed = landedState(root)
	return f
}

// branch names the checked-out branch. `rev-parse --abbrev-ref` answers "HEAD" when
// detached and fails outright on an unborn branch, so the symbolic ref settles the second
// case: a repo with no commits still has a named branch, and reporting it as detached would
// state something false.
func branch(root string) string {
	if name, err := git.Output("-C", root, "rev-parse", "--abbrev-ref", "HEAD"); err == nil && name != "" && name != "HEAD" {
		return "`" + name + "`"
	}
	if name, err := git.Output("-C", root, "symbolic-ref", "--quiet", "--short", "HEAD"); err == nil && name != "" {
		return "`" + name + "`"
	}
	return unknownBranch
}

func head(root string) string {
	sha, err := git.Output("-C", root, "rev-parse", "HEAD")
	if err != nil || sha == "" {
		return unknownHead
	}
	return "HEAD `" + status.Short(sha) + "`"
}

// landedState renders the dirty and unpushed clauses. git.LandedState answers both from one
// repository-wide pass and fails as a unit — a repo whose default branch does not resolve
// has no unpushed count to report — so both clauses degrade together.
//
// The count excludes the handoff itself. The command rewrites that file on every run, so a
// count including it is not a fact about the tree the reader inherited: on a tracked clean
// tree the first run prints "clean tree" and its own write makes the second run print
// "1 dirty path", which is the confident-wrong-fact this block exists to remove and breaks
// the byte-identical guarantee repeated application rests on. Counting the pending write
// instead of excluding it does not close the circle — the count is part of the content
// whose presence makes the file dirty, so that reading can never report a clean tree even
// when the tree genuinely is one — and the reader loses nothing, because "the handoff
// changed" is precisely what this invocation just did.
func landedState(root string) (dirty, unpushed string) {
	fact, err := git.LandedState(root)
	if err != nil {
		return unknownDirty, unknownUnpushed
	}
	paths := fact.DirtyPaths
	if handoffIsDirty(root) && paths > 0 {
		paths--
	}
	dirty = "clean tree"
	if paths > 0 {
		dirty = status.Plural(paths, "dirty path", "dirty paths")
	}
	return dirty, status.Plural(fact.UnpushedCommits, "unpushed commit", "unpushed commits")
}

// handoffIsDirty reports whether the handoff is among the dirty paths the repository-wide
// pass counted. A failed query answers false, leaving the count as git reported it: an
// unanswerable query is not evidence about any path, and overstating the tree's cleanliness
// is the worse of the two errors.
func handoffIsDirty(root string) bool {
	out, err := git.Output("-C", root, "status", "--porcelain", "--", status.HandoffFile)
	return err == nil && out != ""
}

// liveSpecs names the staged spec with its Status line, in path order. An implemented spec
// is finished work and tells a resuming session nothing about what to build, and a spec
// carrying no Status line is malformed rather than staged — naming it with an invented
// status would state something the file does not say.
func liveSpecs(root string) []string {
	all, err := spec.Facts(root)
	if err != nil {
		return nil
	}
	var live []string
	for _, s := range all {
		if s.Status == "" || s.Status == spec.StatusImplemented {
			continue
		}
		live = append(live, "`specs/"+s.Slug+".md` (Status: "+s.Status+")")
	}
	return live
}

// gateField renders the verdict together with the tree it was computed on and whether that
// tree is still the work tree. The three parts travel as one: a bare verdict invites a cold
// session to read a cached green as a statement about the tree it actually inherited.
//
// Staleness is the inspection's own, never a comparison made here — internal/status owns
// that rule, and a second derivation of it is how the two surfaces come to disagree. A
// verdict in any state other than ready is a statement about the gate run rather than about
// a tree, so it carries its cached tree without a staleness clause.
//
// An invalid or unavailable inspection classified the cache without ever parsing a record,
// so it has no cached tree to name. That case says so and falls back to the work tree the
// reader is actually on, which the inspection resolves before it reaches the cache.
func gateField(gv status.GateInfo) string {
	if !gv.Present {
		return "no gate has run."
	}
	verdict := gv.State
	if gv.State == string(gate.Ready) {
		verdict = gv.Status
	}
	if gv.CachedTree == "" {
		return verdict + " — no cached tree survives, work tree " + treeRef(gv.WorkTree)
	}
	field := verdict + " at " + treeRef(gv.CachedTree)
	if gv.State != string(gate.Ready) {
		return field
	}
	if gv.Stale {
		return field + " — stale, work tree " + treeRef(gv.WorkTree)
	}
	return field + " — current"
}

// treeRef renders a tree hash as an inline reference, or names its absence. Every tree the
// gate field prints goes through here, so no path can emit the empty inline-code span a
// reader would take for a deliberate blank rather than for a fact that does not exist.
func treeRef(tree string) string {
	if tree == "" {
		return "unknown"
	}
	return "`" + status.Short(tree) + "`"
}

// benchCommandPrefix opens a `bench` subcommand action, the second of the two shapes a
// board action takes when it is something a session can type.
const benchCommandPrefix = "bench "

// IsInvocable reports whether a board action is something a session can type. It is
// exported as the one source of that rule: the runtime contract's expectation reads it
// rather than restating the prefixes, so a change to what qualifies cannot leave a second
// copy asserting the superseded rule.
//
// Qualifying is syntactic, against the canonical form the board writes: an action is a
// command when it opens a phase invocation or a `bench` subcommand. The opening is
// necessary but not sufficient — a row may join several steps with status.StepSeparator
// ("/bench-final-check / push" once the tree is clean), and that string opens a phase
// invocation while being two commands, which is not something a reader can run. Splitting
// it and taking an arm would be this package deciding what the board meant by a sequence,
// so a compound action does not qualify.
func IsInvocable(action string) bool {
	if strings.Contains(action, status.StepSeparator) {
		return false
	}
	return strings.HasPrefix(action, harnessPrefix[harnessClaude]) || strings.HasPrefix(action, benchCommandPrefix)
}

// nextAction selects the board's next command under root, and reports whether a board
// carrying signals offered none. The board is the one source of what to do next, so the
// handoff and `bench status` cannot disagree about it.
func nextAction(root string) (action, signal string, noneInvocable bool) {
	signals := status.Signals(root)
	if action, signal, ok := firstInvocable(signals); ok {
		return action, signal, false
	}
	return "", "", len(signals) > 0
}

// firstInvocable takes the leading signal whose action is a command rather than a hint,
// walking the signals in the order they arrive. That order is the board's severity ladder
// and internal/status owns it, so this re-ranks nothing: the choice is only whether a row
// qualifies — which IsInvocable decides — never whether it outranks another. Most board
// actions are prose describing a situation ("fix before commit", "split (craft-seams)"),
// and a field promising an invocation cannot render one of those.
func firstInvocable(signals []status.Signal) (action, name string, ok bool) {
	for _, s := range signals {
		if IsInvocable(s.Action) {
			return s.Action, s.Name, true
		}
	}
	return "", "", false
}
