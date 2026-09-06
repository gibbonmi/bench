// Package maps owns the active decision-map query model.
package maps

import (
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// DecisionsDir is the active decision-map directory, relative to the repository root.
const DecisionsDir = "decisions"

var grammar = usage.Grammar{
	Cmd:  "bench maps",
	Help: "usage: bench maps [<map>] [--count|--template|--ticket-template]",
	Flags: []usage.Flag{
		{Name: "--count"},
		{Name: "--template"},
		{Name: "--ticket-template"},
	},
	MaxArgs: 1,
}

// exclusiveMapFlags is the mutually exclusive flag set, in the order the help line
// spells it. The first conflicting pair names the refusal.
var exclusiveMapFlags = []string{"--count", "--template", "--ticket-template"}

// mapOperandValid reports whether the operand is a bare map name. A slash, a
// ".md" suffix, or a byte below 0x20 names a path or carries a terminal escape,
// so the grammar refuses it here instead of resolving it against the tree.
func mapOperandValid(name string) bool {
	if strings.Contains(name, "/") || strings.HasSuffix(name, ".md") {
		return false
	}
	for i := 0; i < len(name); i++ {
		if name[i] < 0x20 {
			return false
		}
	}
	return true
}

type activeScan struct {
	rows         [][]any
	invalidPaths map[string]string
	// names holds every discovered active map, filtered or not, so the operand
	// refusal separates "no such map" from "this map projects no row".
	names      map[string]bool
	count      int
	readyCount int
	state      bounds.FileState
	reason     string
}

// ActiveRows projects one active-directory scan into query rows and its distinct map count.
func ActiveRows(root string) ([][]any, int, bounds.FileState) {
	s := scanActive(root)
	return s.rows, s.count, s.state
}

// ActiveCounts returns the unresolved-or-invalid and ready active-map counts.
func ActiveCounts(root string) (unresolved, ready int, state bounds.FileState) {
	s := scanActive(root)
	return s.count, s.readyCount, s.state
}

func scanActive(root string) activeScan { return scanActiveMap(root, "") }

// scanActiveMap scans the active decision maps. A non-empty only names the one map
// the operand selects: every other map contributes no row and no count, so the
// filter reaches `--count` as well as the table.
func scanActiveMap(root, only string) activeScan {
	dir := filepath.Join(root, DecisionsDir)
	candidates, state, reason := discoverDirectoryCandidates(root, dir, false)
	if state != bounds.StateParsed {
		if state == bounds.StateAbsent {
			return activeScan{state: bounds.StateParsed, names: map[string]bool{}}
		}
		return activeScan{state: state, reason: reason, names: map[string]bool{}}
	}

	s := activeScan{state: bounds.StateParsed, invalidPaths: make(map[string]string), names: map[string]bool{}}
	for _, candidate := range candidates {
		mapName := strings.TrimSuffix(filepath.Base(candidate.Path), filepath.Ext(candidate.Path))
		s.names[mapName] = true
		if only != "" && mapName != only {
			continue
		}
		file := bounds.Classify(filepath.Join(root, filepath.FromSlash(candidate.Path)), bounds.ControlRecordLimit)
		if file.State != bounds.StateParsed {
			s.invalid(mapName, candidate.Path, string(file.State)+": "+file.Reason)
			continue
		}
		m, diagnostics := ValidateDecisionMap(root, candidate.Path, false, file.Data)
		if len(diagnostics) > 0 {
			s.invalid(mapName, candidate.Path, diagnostics[0].Message)
			continue
		}
		rows := projectedRows(mapName, candidate.Path, m)
		if m.Status == "shaping" {
			s.count++
		} else if m.Status == "ready" {
			s.readyCount++
		}
		if len(rows) > 0 {
			s.rows = append(s.rows, rows...)
		}
		s.rows = append(s.rows, staleRows(root, mapName, candidate.Path, m.Sources)...)
	}
	sort.SliceStable(s.rows, func(i, j int) bool { return s.rows[i][0].(string) < s.rows[j][0].(string) })
	return s
}

func (s *activeScan) invalid(name, path, reason string) {
	row := []any{name, "invalid", "map", "invalid", reason, path}
	s.rows = append(s.rows, row)
	s.invalidPaths[invalidRowKey(row)] = path
	s.count++
}

// ticketPath names the file one projected ticket row points the reader at.
func ticketPath(indexPath, id string) string {
	return strings.TrimSuffix(indexPath, ".md") + "/" + ticketsDirName + "/" + id + ".md"
}

func projectedRows(name, path string, m DecisionMap) [][]any {
	byID := make(map[string]DecisionTicket, len(m.Tickets))
	for _, ticket := range m.Tickets {
		byID[ticket.ID] = ticket
	}
	var rows [][]any
	for _, ticket := range m.Tickets {
		if resolved(ticket) {
			continue
		}
		state := unresolvedState(ticket)
		blockers := unresolvedBlockerTitles(ticket, byID)
		if state == "frontier" && blockers != "" {
			state = "blocked"
		}
		rows = append(rows, []any{name, ticket.Title, ticket.Type, state, blockers, ticketPath(path, ticket.ID)})
	}
	if len(rows) == 0 && m.Status == "shaping" {
		rows = append(rows, []any{name, "Not yet specified", "fog", "shaping", "", path})
	}
	// A ready map has no unresolved ticket left to project, so without this row the
	// map that is ready to be specified is the one the default view omits.
	if m.Status == "ready" {
		rows = append(rows, []any{name, m.Title, "map", "ready", "", path})
	}
	return rows
}

func unresolvedState(ticket DecisionTicket) string {
	return ticketAnswerState(ticket)
}

func unresolvedBlockerTitles(ticket DecisionTicket, byID map[string]DecisionTicket) string {
	var titles []string
	for _, id := range blockers(ticket.BlockedBy) {
		if blocker, ok := byID[id]; ok && !resolved(blocker) {
			titles = append(titles, blocker.Title)
		}
	}
	return strings.Join(titles, ", ")
}

// Rows returns the active-map query rows.
func Rows(root string) [][]any {
	rows, _, _ := ActiveRows(root)
	return rows
}

// UnresolvedCount returns the distinct active-map count and active scan state.
func UnresolvedCount(root string) (int, bounds.FileState) {
	_, count, state := ActiveRows(root)
	return count, state
}

// Command implements `bench maps`.
func Command(args []string) (string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	var chosen []string
	for _, flag := range exclusiveMapFlags {
		if _, set := parsed.Flags[flag]; set {
			chosen = append(chosen, flag)
		}
	}
	if len(chosen) > 1 {
		return grammar.Help + " (" + chosen[0] + " and " + chosen[1] + " are mutually exclusive)\n", 2
	}
	if _, template := parsed.Flags["--template"]; template {
		return DecisionMapTemplate(), 0
	}
	if _, template := parsed.Flags["--ticket-template"]; template {
		return DecisionTicketTemplate(), 0
	}
	only := ""
	if len(parsed.Positionals) == 1 {
		only = parsed.Positionals[0]
		if !mapOperandValid(only) {
			// The refusal prints the grammar rather than echoing the operand, because a
			// path-shaped or control-carrying operand is exactly what must not be replayed.
			return grammar.Help + "\n", 2
		}
	}
	root, err := git.Root()
	if err != nil {
		if _, count := parsed.Flags["--count"]; count {
			return "0\n", 0
		}
		return toon.NotInRepo() + "\n", 1
	}
	s := scanActiveMap(root, only)
	if s.state.Failed() {
		return toon.RecordError(DecisionsDir, s.state, s.reason) + "\n", 1
	}
	if only != "" && !s.names[only] {
		return "maps: no active map named " + strconv.Quote(only) + "\n", 1
	}
	if _, count := parsed.Flags["--count"]; count {
		return strconv.Itoa(s.count) + "\n", 0
	}
	out, err := toon.TableTyped("maps", []string{"map", "title", "type", "state", "blockers", "path"}, s.rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	help, err := axi.RenderHelp(actionsForRows(s.rows, s.invalidPaths))
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	out += help
	for _, row := range s.rows {
		if row[3] == "invalid" {
			return out, 1
		}
	}
	return out, 0
}

func actionsForRows(rows [][]any, invalidPaths map[string]string) []axi.Action {
	actions := make([]axi.Action, 0, len(rows))
	for _, row := range rows {
		mapName, title, state := row[0].(string), row[1].(string), row[3].(string)
		switch state {
		case "frontier":
			actions = append(actions, axi.HarnessPhase("/bench-shape-idea", "shape "+mapName+": "+title))
		case "ready":
			actions = append(actions, axi.HarnessPhaseOn("/bench-write-spec", row[5].(string), "spec "+mapName+": "+title))
		case "invalid":
			path := invalidPaths[invalidRowKey(row)]
			actions = append(actions, axi.ExecutableInvocation("repair "+path, axi.KnownArgument("maps"), axi.KnownArgument("--template")))
		}
	}
	return actions
}

func invalidRowKey(row []any) string {
	return row[0].(string) + "\x00" + row[4].(string)
}
