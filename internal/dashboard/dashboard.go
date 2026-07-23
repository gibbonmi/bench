// Package dashboard ports `bench dashboard`: it renders one self-contained static HTML
// snapshot of the project board — the gate verdict, the ambient signals, the roadmap and
// its recommended sequence, the parked ideas, the open-learnings count, and the worktree
// pool — for a human to open in a browser. It is a consumer of the existing readers, never
// a new source: the Snapshot is composed from status.Signals / status.GateVerdict, the
// roadmap readers, and the worktree classifier, so the page can never diverge from what
// the CLI surfaces already report.
//
// Render is a pure function of its Snapshot — no IO, no clock, no git — so every case,
// including hostile input, is testable without a browser or a repo. Command gathers the
// Snapshot (the one place IO and the wall clock live), renders it, and either writes the
// file atomically under the git dir or emits it on stdout for `--stdout`.
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

// grammar is the declared argument shape usage.Parse enforces for this subcommand —
// arity, flag recognition, `--`, and help all come from there rather than a local switch.
var grammar = usage.Grammar{
	Cmd:   "bench dashboard",
	Help:  "usage: bench dashboard [--stdout]",
	Flags: []usage.Flag{{Name: "--stdout"}},
}

// Snapshot is the complete data the page renders — everything time- or environment-
// dependent gathered up front so Render stays pure. It is assembled by gather from the
// existing readers and is hand-constructible in tests.
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
	// WorktreesErr carries the git worktree-list query's error when the classify
	// itself failed, so Render can show the failure rather than an empty pool pane
	// that would otherwise read as "no worktrees" — the false-empty class FT29 swept.
	WorktreesErr string
}

// Command implements `bench dashboard [--stdout]`. With no argument it writes the rendered
// page atomically to `<git-dir>/bench-dashboard.html` and prints that path (exit 0);
// `--stdout` emits the HTML on stdout and writes nothing. `--stdout` is the only accepted
// flag — any other argument is a usage error (exit 2); outside a git repository is the
// structured not-in-repo error (exit 1), so a bad invocation never writes a stray file.
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
	gitDir, err := git.Output("-C", root, "rev-parse", "--absolute-git-dir")
	if err != nil {
		return toon.Errorf("cannot resolve git dir", err.Error()) + "\n", 1
	}
	target := filepath.Join(gitDir, "bench-dashboard.html")
	if err := atomicWrite(target, page); err != nil {
		return toon.Errorf("cannot write bench-dashboard.html", err.Error()) + "\n", 1
	}
	return target + "\n", 0
}

// gather composes the Snapshot from the existing readers — the one place IO and the wall
// clock enter. It never re-parses a source: the signals ladder and the gate verdict come
// from internal/status, the roadmap text/sequence, ideas, and learnings count from
// internal/roadmap, and the worktree pool from the shared classifier.
func gather(root string) Snapshot {
	text, present := roadmap.RoadmapText(root)
	_, openLearnings := roadmap.DrainCounts(root)
	registered, werr := worktree.ClassifyRegisteredWorktrees(root)
	snap := Snapshot{
		GeneratedAt:    time.Now(),
		Gate:           status.GateVerdict(root),
		Signals:        status.Signals(root),
		RoadmapText:    text,
		RoadmapPresent: present,
		Sequence:       roadmap.RecommendedSequence(text),
		Ideas:          roadmap.ParkedIdeas(root),
		OpenLearnings:  openLearnings,
		Worktrees:      poolWorktrees(registered),
	}
	if werr != nil {
		snap.WorktreesErr = werr.Error()
	}
	return snap
}

// poolWorktrees drops the repo root — the expected primary checkout — and keeps the pool
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
// over the target — so an interrupt mid-write never leaves a half-written page at the
// published path and a reader never sees a partial file. On any error the temp is removed,
// leaving no leftover scratch.
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
