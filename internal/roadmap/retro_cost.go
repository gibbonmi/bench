package roadmap

import (
	"path"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/otelrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/tickets"
)

// costBudget is one per-ticket budget of the cost section. value selects the build total,
// in the unit of the line, that the line divides by the ticket count.
type costBudget struct {
	label, unit string
	limit       float64
	value       func(costFacts) float64
}

// costBudgets is the one budget table. No guidance file states a budget number, so a
// changed budget changes here and nowhere else.
var costBudgets = []costBudget{
	{label: "wall per ticket", unit: " min", limit: 60, value: func(f costFacts) float64 { return float64(f.wall) / 60 }},
	{label: "bookkeeping per ticket", limit: 7.1, value: func(f costFacts) float64 { return float64(f.bookkeeping) }},
}

// overBudget ends a per-ticket line whose value is above its budget. The drain finds each
// breach by this one marker.
const overBudget = ", over budget"

// costFacts holds the derived facts of one landing. ranged marks the facts of the build
// range, graphed the facts of the plan graph, and dispatched the version 2 dispatch counts.
// wall is in seconds, so a per-ticket value keeps the part of a minute.
type costFacts struct {
	landing                             string
	ranged                              bool
	wall, total, bookkeeping, change    int
	graphed                             bool
	tickets, chain, chunks, checkpoints int
	dispatched                          bool
	authors, repairs                    int
	assignments                         map[string]bool
}

// traceSeams are the seams that the trace line sums, in line order, each with its label.
var traceSeams = []struct{ label, seam string }{{"commit", otelrecord.SeamCommit}, {"lane", otelrecord.SeamLane}, {"landing", otelrecord.SeamLanding}}

// costSection prints the build cost of slug from its landing commit L. It reads git
// history and the tree of L, so a retired spec folder loses no fact. Each fact that the
// tree does not hold reads unknown, so a missing measurement never blocks the draft.
func costSection(home, root, slug string) string {
	f := deriveCost(root, slug)
	rangeText := unknownFact
	if f.landing != "" {
		rangeText = buildRange(f.landing)
	}
	lines := []string{
		"- landing: " + orUnknown(f.landing) + ", range " + rangeText,
		"- wall: " + costCount(f.wall/60, f.ranged) + " min",
		"- commits: " + costCount(f.total, f.ranged) + " total, " + costCount(f.bookkeeping, f.ranged) + " bookkeeping, " + costCount(f.change, f.ranged) + " change",
		"- graph: " + costCount(f.tickets, f.graphed) + " tickets, chain " + costCount(f.chain, f.graphed) + ", " + costCount(f.chunks, f.graphed) + " chunks, " + costCount(f.checkpoints, f.graphed) + " serial checkpoints",
		"- dispatch: " + costCount(f.authors, f.dispatched) + " authors, " + costCount(f.repairs, f.dispatched) + " repairs",
		"- trace: " + traceTimes(home, root, f.assignments),
	}
	return strings.Join(append(lines, budgetLines(f)...), "\n")
}

// traceTimes states the count and the summed milliseconds of each seam in traceSeams. It
// keeps each trace that holds a span of an assignment in ids. An absent record and a record
// without a span of those assignments state unknown, because a zero would read as measured.
func traceTimes(home, root string, ids map[string]bool) string {
	spans, err := otelrecord.ReadSpans(home, root)
	if err != nil {
		return unknownFact
	}
	traces := map[string]bool{}
	for _, span := range spans {
		if ids[span.Attributes[otelrecord.AttrAssignmentID]] {
			traces[span.TraceID] = true
		}
	}
	if len(traces) == 0 {
		return unknownFact
	}
	parts := make([]string, 0, len(traceSeams))
	for _, seam := range traceSeams {
		count, total := 0, int64(0)
		for _, span := range spans {
			if traces[span.TraceID] && span.Seam == seam.seam {
				count, total = count+1, total+span.Elapsed().Milliseconds()
			}
		}
		parts = append(parts, seam.label+" "+strconv.Itoa(count)+" spans in "+strconv.FormatInt(total, 10)+" ms")
	}
	return strings.Join(parts, ", ")
}

// budgetLines compares each per-ticket value with its budget in costBudgets. The compare
// reads the printed value, so a line never marks a value that prints equal to its budget.
func budgetLines(f costFacts) []string {
	lines := make([]string, 0, len(costBudgets))
	for _, budget := range costBudgets {
		value, over := unknownFact, false
		if f.ranged && f.graphed && f.tickets > 0 {
			value = strconv.FormatFloat(budget.value(f)/float64(f.tickets), 'f', 1, 64)
			printed, err := strconv.ParseFloat(value, 64)
			over = err == nil && printed > budget.limit
		}
		line := "- " + budget.label + ": " + value + budget.unit + " against " + strconv.FormatFloat(budget.limit, 'f', -1, 64) + budget.unit
		if over {
			line += overBudget
		}
		lines = append(lines, line)
	}
	return lines
}

func costCount(value int, known bool) string {
	if !known {
		return unknownFact
	}
	return strconv.Itoa(value)
}

// deriveCost reads the facts of the landing that spec.Landing names. A slug without a
// landing answers no fact.
func deriveCost(root, slug string) costFacts {
	var f costFacts
	landing, found, err := spec.Landing(root, slug)
	if err != nil || !found {
		return f
	}
	f.landing = landing
	rangeCost(root, slug, &f)
	planCost(root, slug, &f)
	return f
}

// rangeCost counts the non-merge commits of L^1..L^2 and splits them by their changed
// paths. The wall time runs from the committer time of the oldest of them to the committer
// time of L, so a folded default-branch commit and the first parent of L stay outside it.
func rangeCost(root, slug string, f *costFacts) {
	end, err := committerTime(root, f.landing)
	if err != nil {
		return
	}
	raw, err := git.Output("-C", root, "log", "--no-merges", "--format=%H %ct", buildRange(f.landing))
	if err != nil {
		return
	}
	oldest := end
	for _, line := range strings.Split(raw, "\n") {
		hash, stamp, ok := strings.Cut(line, " ")
		if !ok {
			continue
		}
		at, err := strconv.ParseInt(stamp, 10, 64)
		if err != nil {
			return
		}
		paths, err := git.Raw("-C", root, "diff-tree", "-z", "-r", "--root", "--no-commit-id", "--name-only", "--no-renames", hash)
		if err != nil {
			return
		}
		oldest = min(oldest, at)
		f.total++
		if bookkeeping(slug, strings.Split(strings.TrimRight(string(paths), "\x00"), "\x00")) {
			f.bookkeeping++
		} else {
			f.change++
		}
	}
	f.wall, f.ranged = int(end-oldest), true
}

// buildRange is the revision range of the build that landing L merged: the commits of its
// second parent that its first parent does not hold.
func buildRange(landing string) string { return landing + "^1.." + landing + "^2" }

func committerTime(root, commit string) (int64, error) {
	out, err := git.Output("-C", root, "show", "-s", "--format=%ct", commit)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(out, 10, 64)
}

// bookkeeping reports whether every changed path of one commit is below the spec folder of
// slug or is its review record.
func bookkeeping(slug string, paths []string) bool {
	specPath := spec.LiveSpecPath(slug)
	record, err := reviewrecord.RecordPath(specPath)
	if err != nil {
		return false
	}
	folder := path.Dir(specPath) + "/"
	for _, changed := range paths {
		if changed != "" && changed != record && !strings.HasPrefix(changed, folder) {
			return false
		}
	}
	return true
}

// planCost reads the plan and its tickets from the tree of L through the decode-only
// reader, so a plan that the current grade refuses still gives its facts.
func planCost(root, slug string, f *costFacts) {
	data, err := git.ReadTreeFile(root, f.landing, spec.LiveSpecPath(slug))
	if err != nil {
		return
	}
	plan, err := reviewrecord.DecodePlan(data)
	if err != nil {
		return
	}
	if plan.Version == 2 && plan.Execution != nil {
		f.assignments = map[string]bool{}
		for _, history := range plan.Execution.Assignments {
			for _, assignment := range history {
				if assignment.Assignment != "" {
					f.assignments[assignment.Assignment] = true
				}
				if assignment.Predecessor != "" {
					f.repairs++
				} else {
					f.authors++
				}
			}
		}
		f.dispatched = true
	}
	owner := map[string]int{}
	var names []string
	for i, chunk := range plan.Chunks {
		for _, name := range chunk.Tickets {
			owner[name] = i
			names = append(names, name)
		}
	}
	blockers := map[string][]string{}
	for _, name := range names {
		data, err := git.ReadTreeFile(root, f.landing, path.Join(ticketsDir(slug), name))
		if err != nil {
			return
		}
		ticket, _ := tickets.ParseTicket(name, data, names, "")
		for _, blocker := range ticket.Blockers {
			if _, planned := owner[blocker]; planned {
				blockers[name] = append(blockers[name], blocker)
			}
		}
	}
	chunkBlockers := map[int][]int{}
	chunkIDs := make([]int, len(plan.Chunks))
	for i := range chunkIDs {
		chunkIDs[i] = i
	}
	for name, list := range blockers {
		for _, blocker := range list {
			if owner[blocker] != owner[name] {
				chunkBlockers[owner[name]] = append(chunkBlockers[owner[name]], owner[blocker])
			}
		}
	}
	f.tickets, f.chunks = len(names), len(plan.Chunks)
	f.chain = longestChain(names, func(name string) []string { return blockers[name] })
	f.checkpoints = longestChain(chunkIDs, func(chunk int) []int { return chunkBlockers[chunk] })
	f.graphed = true
}

// longestChain counts the nodes of the longest blocker chain. An ungraded plan can hold a
// cycle, so a node on the open walk reads as depth 0 and the walk ends.
func longestChain[K comparable](nodes []K, next func(K) []K) int {
	depth := map[K]int{}
	var walk func(K) int
	walk = func(node K) int {
		if known, seen := depth[node]; seen {
			return known
		}
		depth[node] = 0
		deepest := 0
		for _, blocker := range next(node) {
			deepest = max(deepest, walk(blocker))
		}
		depth[node] = deepest + 1
		return deepest + 1
	}
	longest := 0
	for _, node := range nodes {
		longest = max(longest, walk(node))
	}
	return longest
}
