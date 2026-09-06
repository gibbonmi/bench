// Package dashboard implements `bench dashboard`. It renders one self-contained static
// HTML snapshot of the project board. The page shows the gate verdict, the ambient
// signals, the roadmap and its sequence, the parked ideas, the open-learnings count, and
// the worktree pool. A human opens this page in a browser.
//
// The package consumes the existing readers. It never adds a new source. The Snapshot is
// composed from status.Signals, status.GateVerdict, the roadmap readers, and the worktree
// classifier. So the page can never diverge from what the CLI surfaces already report.
//
// Render is a pure function of its Snapshot: no IO, no clock, no git. So every case,
// including hostile input, is testable without a browser or a repo. Command gathers the
// Snapshot, the one place IO and the wall clock live. Command then renders it. It either
// writes the file atomically under the git dir, or emits it on stdout for `--stdout`.
package dashboard

import (
	"os"
	"path/filepath"
	"time"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/roadmap"
	"github.com/gibbonmi/bench/internal/status"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
	"github.com/gibbonmi/bench/internal/worktree"
)

// grammar is the declared argument shape usage.Parse enforces for this subcommand.
// Arity, flag recognition, `--`, and help all come from there, instead of a local
// switch.
var grammar = usage.Grammar{
	Cmd:   "bench dashboard",
	Help:  "usage: bench dashboard [--stdout]",
	Flags: []usage.Flag{{Name: "--stdout"}},
}

// Snapshot is the complete data the page renders. It gathers everything time- or
// environment-dependent up front, so Render stays pure. gather assembles it from the
// existing readers. A test can construct it by hand.
type Snapshot struct {
	GeneratedAt    time.Time
	Gate           status.GateInfo
	Signals        []status.Signal
	RoadmapText    string
	RoadmapPresent bool
	Sequence       string
	Ideas          []string
	OpenLearnings  int
	Worktrees      []worktree.Registered
	// WorktreesErr carries the git worktree-list query's error, for when the classify
	// step itself fails. Render then shows the failure, instead of an empty pool pane
	// that would otherwise read as "no worktrees" (the false-empty class from FT29).
	WorktreesErr string
}

// Command implements `bench dashboard [--stdout]`. With no argument, Command writes the
// rendered page atomically to `<git-dir>/bench-dashboard.html`. Command then prints that
// path, exit 0. `--stdout` emits the HTML on stdout instead. It writes nothing else.
//
// `--stdout` is the only accepted flag. Any other argument gives a usage error, exit 2.
// Outside a git repository, Command gives the structured not-in-repo error, exit 1. So a
// bad invocation never writes a stray file.
func Command(args []string) (string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	_, toStdout := parsed.Flags["--stdout"]
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	page := Render(gather(root))
	if toStdout {
		return page, 0
	}
	gitDir, err := git.AdminDir(root)
	if err != nil {
		return toon.Errorf("cannot resolve git dir", err.Error()) + "\n", 1
	}
	target := filepath.Join(gitDir, "bench-dashboard.html")
	if err := atomicWrite(target, page); err != nil {
		return toon.Errorf("cannot write bench-dashboard.html", err.Error()) + "\n", 1
	}
	return target + "\n", 0
}

// gather composes the Snapshot from the existing readers. This is the one place IO and
// the wall clock enter. gather never re-parses a source. The signals ladder and the gate
// verdict come from internal/status. The roadmap text, sequence, ideas, and learnings
// count come from internal/roadmap. The worktree pool comes from the shared classifier.
func gather(root string) Snapshot {
	text, present := roadmap.RoadmapText(root)
	// The dashboard renders the count only. A failed learnings read degrades to 0 here,
	// matching drainStatus's posture. The status board owns the fail-closed unknown row.
	drain := roadmap.DrainCounts(root)
	registered, werr := worktree.ClassifyRegisteredWorktrees(root)
	snap := Snapshot{
		GeneratedAt:    time.Now(),
		Gate:           status.GateVerdict(root),
		Signals:        status.Signals(root),
		RoadmapText:    text,
		RoadmapPresent: present,
		Sequence:       roadmap.RecommendedSequence(text),
		Ideas:          roadmap.ParkedIdeas(root),
		OpenLearnings:  drain.OpenLearnings,
		Worktrees:      poolWorktrees(registered),
	}
	if werr != nil {
		snap.WorktreesErr = werr.Error()
	}
	return snap
}

// poolWorktrees drops the repo root, the expected primary checkout. It keeps the pool
// entries a reviewer wants to see: out-of-pool, leased, and warm worktrees.
func poolWorktrees(regs []worktree.Registered) []worktree.Registered {
	var out []worktree.Registered
	for _, r := range regs {
		if r.Class == worktree.ClassRoot {
			continue
		}
		out = append(out, r)
	}
	return out
}

// atomicWrite renders to a sibling temp file in the target's directory, then renames it
// over the target. So an interrupt mid-write never leaves a half-written page at the
// published path. A reader never sees a partial file. On any error, atomicWrite removes
// the temp file, leaving no leftover scratch.
func atomicWrite(target, content string) error {
	dir := filepath.Dir(target)
	tmp, err := os.CreateTemp(dir, "bench-dashboard-*.html.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	if _, err := tmp.WriteString(content); err != nil {
		tmp.Close()
		os.Remove(name)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(name)
		return err
	}
	if err := os.Rename(name, target); err != nil {
		os.Remove(name)
		return err
	}
	return nil
}
