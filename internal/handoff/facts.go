package handoff

import (
	"os"
	"path/filepath"

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

// branch renders the checked-out branch for the pin block. git.CheckedOutBranch owns the
// probe; this adds the two things the pin block wants and the probe does not answer — the
// literal "HEAD" is detachment rather than a branch name, and a real name is backticked
// for the markdown the block is written in.
func branch(root string) string {
	if name, err := git.CheckedOutBranch(root); err == nil && name != "" && name != "HEAD" {
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

// landedState renders the dirty and unpushed clauses while excluding this command's capture
// file. Rewriting the handoff is not inherited work for the next session to commit.
func landedState(root string) (dirty, unpushed string) {
	fact, err := git.LandedState(root, status.HandoffFile)
	if err != nil {
		return unknownDirty, unknownUnpushed
	}
	dirty = "clean tree"
	if fact.DirtyPaths > 0 {
		dirty = status.Plural(fact.DirtyPaths, "dirty path", "dirty paths")
	}
	return dirty, status.Plural(fact.UnpushedCommits, "unpushed commit", "unpushed commits")
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
		live = append(live, "`"+s.Path+"` (Status: "+s.Status+")")
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

func applyRoute(f *facts, route status.RouteResult) {
	f.Action = route.Lead.Action
	f.Signal = route.Lead.Name
	f.NoInvocable = route.NoCommand
}
