package handoff

import (
	"os"
	"path/filepath"
	"strconv"
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
	return "HEAD `" + short(sha) + "`"
}

// landedState renders the dirty and unpushed clauses. git.LandedState answers both from one
// repository-wide pass and fails as a unit — a repo whose default branch does not resolve
// has no unpushed count to report — so both clauses degrade together.
func landedState(root string) (dirty, unpushed string) {
	fact, err := git.LandedState(root)
	if err != nil {
		return unknownDirty, unknownUnpushed
	}
	dirty = "clean tree"
	if fact.DirtyPaths > 0 {
		dirty = plural(fact.DirtyPaths, "dirty path", "dirty paths")
	}
	return dirty, plural(fact.UnpushedCommits, "unpushed commit", "unpushed commits")
}

// liveSpecs names every spec that is not yet implemented, in path order. An implemented
// spec is finished work and tells a resuming session nothing about what to build.
func liveSpecs(root string) []string {
	all, err := spec.Facts(root)
	if err != nil {
		return nil
	}
	var live []string
	for _, s := range all {
		if s.Status == "implemented" {
			continue
		}
		state := s.Status
		if state == "" {
			state = "unknown"
		}
		live = append(live, "`specs/"+s.Slug+".md` (Status: "+state+")")
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
	return "`" + short(tree) + "`"
}

// benchCommandPrefix opens a `bench` subcommand action, the second of the two shapes a
// board action takes when it is something a session can type.
const benchCommandPrefix = "bench "

// boardStepSeparator is how a board row joins the steps of a sequence into one action.
// internal/status writes it; this package only recognizes it, to tell an action naming one
// command apart from an action naming several.
const boardStepSeparator = " / "

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
// qualifies, never whether it outranks another.
//
// Most board actions are prose describing a situation — "fix before commit",
// "split (craft-seams)" — and a field promising an invocation cannot render one of those.
// So qualifying is syntactic, against the canonical form the board writes: an action is a
// command when it opens a phase invocation or a `bench` subcommand and never otherwise.
//
// The opening is necessary but not sufficient. A board row may join several steps with the
// separator below — the git row reads "/bench-final-check / push" once the tree is clean —
// and that string opens a phase invocation while being two commands, which is not something
// a reader can run. Splitting it and taking an arm would be this package deciding what the
// board meant by a sequence, so a compound action simply does not qualify and the walk
// continues.
func firstInvocable(signals []status.Signal) (action, name string, ok bool) {
	for _, s := range signals {
		if strings.Contains(s.Action, boardStepSeparator) {
			continue
		}
		if strings.HasPrefix(s.Action, harnessPrefix[harnessClaude]) || strings.HasPrefix(s.Action, benchCommandPrefix) {
			return s.Action, s.Name, true
		}
	}
	return "", "", false
}

// short is the seven-byte tree/commit prefix the board renders, guarding a value shorter
// than the slice.
func short(s string) string { return s[:min(7, len(s))] }

func plural(n int, one, many string) string {
	unit := many
	if n == 1 {
		unit = one
	}
	return strconv.Itoa(n) + " " + unit
}
