// Package roadmapflow derives the board's flow from the detail files each commit touched.
// The flow is the rows opened, fed, and retired over a window of recent history. This
// package lives beside package roadmap, not inside it. Package roadmap sits at 11 source
// files against a reviewer-owned directory budget of 12. The open mass this package
// reports is read back through roadmap.LoadTree, so the row count keeps one source.
package roadmapflow

import (
	"regexp"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/roadmap"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// windowDrains is how many drain commits the window spans. A drain lands about weekly,
// so three of them frame roughly the last fortnight of the board's movement.
const windowDrains = 3

// unknownMass is the open-mass cell for a board that could not be read. A degraded board
// must never render as an empty one, so the count degrades to a word rather than to zero.
const unknownMass = "unknown"

// isDetailPath reports whether a `git log --name-status` path names one row's detail
// owner. A detail owner is a `.md` file directly under the roadmap directory whose
// basename is a row ID. The row-ID grammar stays with package roadmap. A file whose
// basename falls outside it contributes to no count, by the board parser's own rule
// rather than a second one.
func isDetailPath(path string) bool {
	name, found := strings.CutPrefix(path, roadmap.RoadmapDir+"/")
	if !found || strings.Contains(name, "/") {
		return false
	}
	id, found := strings.CutSuffix(name, ".md")
	return found && roadmap.ValidOccurrenceOwner(id)
}

// commitID matches the `%H` line that opens each record of the history query. Every
// other line of that output is a status letter and a path, so a hexadecimal identity is
// the one unambiguous record boundary.
var commitID = regexp.MustCompile(`^[0-9a-f]{40}$`)

var flowGrammar = usage.Grammar{
	Cmd:   "bench roadmap --flow",
	Help:  "usage: bench roadmap --flow",
	Flags: []usage.Flag{{Name: "--flow", Required: true}},
}

// event is one commit that touched at least one detail file. It carries its identity and
// the per-status counts of the detail files it touched. A commit that only modifies
// detail files is an event too, since it feeds rows even though it opens none.
type event struct {
	id                   string
	opened, fed, retired int
}

// span is the window's summed evidence. found is false when the history holds no event at
// all. That is the definitive empty flow, not a window of zeroes.
type span struct {
	opened, fed, retired int
	drains               int
	from, to             string
	found                bool
}

// Command implements `bench roadmap --flow`: one flow block over the window that spans
// the last three drain commits, then the AXI help envelope.
func Command(args []string) (string, int) {
	if _, line, code := usage.Parse(flowGrammar, args); line != "" {
		return line + "\n", code
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	events, err := history(root)
	if err != nil {
		return toon.Errorf("roadmap flow unavailable", err.Error()) + "\n", 1
	}
	return render(selectWindow(events), openMass(root))
}

// history returns every commit that touched a detail file, newest first. The query reads
// file status rather than commit subjects, so no subject grammar can inflate or hide a
// count. --no-renames keeps a move inside the directory as one add plus one delete,
// rather than a rename git would report as neither.
func history(root string) ([]event, error) {
	// An unborn HEAD is a repository with no history, not a failed query. `git log` exits
	// nonzero on it, so this code answers it before the query runs.
	if !git.OK("-C", root, "rev-parse", "--verify", "--quiet", "HEAD") {
		return nil, nil
	}
	out, err := git.Output("-C", root, "log", "--name-status", "--no-renames", "--format=%H", "--", roadmap.RoadmapDir)
	if err != nil {
		return nil, err
	}
	var events []event
	var current *event
	for _, line := range strings.Split(out, "\n") {
		if commitID.MatchString(line) {
			events = append(events, event{id: line})
			current = &events[len(events)-1]
			continue
		}
		status, path, found := strings.Cut(line, "\t")
		if current == nil || !found || !isDetailPath(path) {
			continue
		}
		switch status {
		case "A":
			current.opened++
		case "M":
			current.fed++
		case "D":
			current.retired++
		}
	}
	kept := events[:0]
	for _, e := range events {
		if e.opened+e.fed+e.retired > 0 {
			kept = append(kept, e)
		}
	}
	return kept, nil
}

// selectWindow walks the events newest first and sums them through the third drain commit
// inclusive. A drain commit is one that adds a detail file, so the boundary derives from
// the same file evidence as the counts. A history holding fewer than three drains yields
// the whole history and reports the drain count it found. A young board still has a flow
// to report.
func selectWindow(events []event) span {
	var window span
	for _, e := range events {
		window.found = true
		window.opened += e.opened
		window.fed += e.fed
		window.retired += e.retired
		window.from = e.id
		if window.to == "" {
			window.to = e.id
		}
		if e.opened > 0 {
			window.drains++
			if window.drains == windowDrains {
				break
			}
		}
	}
	return window
}

// openMass returns the number of index rows the current board holds, or the unknown cell
// when the board did not read. An unreadable index and a detail directory that could not
// be listed both leave the count underived, not zero.
func openMass(root string) any {
	tree := roadmap.LoadTree(root)
	if tree.Index.State != bounds.StateParsed || tree.DirState.Failed() {
		return unknownMass
	}
	doc, _, _ := roadmap.ParseDocument(tree, nil, false)
	return len(doc.Rows)
}

// render emits the flow block and the help envelope. The block's cells are counts, a
// boolean, and two hexadecimal identities. No git-sourced text, such as a branch name, a
// path, or a subject, reaches a field.
func render(window span, mass any) (string, int) {
	fields := []string{"opened", "fed", "retired", "net", "open_mass", "target_met", "drains", "window_from", "window_to"}
	var rows [][]any
	var actions []axi.Action
	if window.found {
		net := window.opened - window.retired
		rows = [][]any{{window.opened, window.fed, window.retired, net, mass, net <= 0, window.drains, window.from, window.to}}
		if net > 0 {
			actions = append(actions, axi.HarnessPhase("/bench-drain", "the flow target does not hold: retire rows in the next drain"))
		}
	}
	block, err := toon.TableTyped("flow", fields, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	help, err := axi.RenderHelp(actions)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return block + help, 0
}
