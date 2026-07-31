// Package maps owns the active decision-map query model.
package maps

import (
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// DecisionsDir is the active decision-map directory, relative to the repository root.
const DecisionsDir = "decisions"

var grammar = usage.Grammar{
	Cmd:  "bench maps",
	Help: "usage: bench maps [--count|--template]",
	Flags: []usage.Flag{
		{Name: "--count"},
		{Name: "--template"},
	},
}

type activeScan struct {
	rows   [][]any
	count  int
	state  bounds.FileState
	reason string
}

// ActiveRows projects one active-directory scan into query rows and its distinct map count.
func ActiveRows(root string) ([][]any, int, bounds.FileState) {
	s := scanActive(root)
	return s.rows, s.count, s.state
}

func scanActive(root string) activeScan {
	dir := filepath.Join(root, DecisionsDir)
	candidates, state, reason := discoverDirectoryCandidates(root, dir, false)
	if state != bounds.StateParsed {
		if state == bounds.StateAbsent {
			return activeScan{state: bounds.StateParsed}
		}
		return activeScan{state: state, reason: reason}
	}

	s := activeScan{state: bounds.StateParsed}
	for _, candidate := range candidates {
		mapName := strings.TrimSuffix(filepath.Base(candidate.Path), filepath.Ext(candidate.Path))
		file := bounds.Classify(filepath.Join(root, filepath.FromSlash(candidate.Path)), bounds.ControlRecordLimit)
		if file.State != bounds.StateParsed {
			s.invalid(mapName, string(file.State)+": "+file.Reason)
			continue
		}
		m, diagnostics := ValidateDecisionMap(root, candidate.Path, false, file.Data)
		if len(diagnostics) > 0 {
			s.invalid(mapName, diagnostics[0].Message)
			continue
		}
		rows := projectedRows(mapName, m)
		if m.Status == "shaping" {
			s.count++
		}
		if len(rows) > 0 {
			s.rows = append(s.rows, rows...)
		}
	}
	sort.SliceStable(s.rows, func(i, j int) bool { return s.rows[i][0].(string) < s.rows[j][0].(string) })
	return s
}

func (s *activeScan) invalid(name, reason string) {
	s.rows = append(s.rows, []any{name, "invalid", "map", "invalid", reason})
	s.count++
}

func projectedRows(name string, m DecisionMap) [][]any {
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
		rows = append(rows, []any{name, ticket.Title, ticket.Type, state, blockers})
	}
	if len(rows) == 0 && m.Status == "shaping" {
		rows = append(rows, []any{name, "Not yet specified", "fog", "shaping", ""})
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
	if _, count := parsed.Flags["--count"]; count {
		if _, template := parsed.Flags["--template"]; template {
			return grammar.Help + " (--count and --template are mutually exclusive)\n", 2
		}
	}
	if _, template := parsed.Flags["--template"]; template {
		return DecisionMapTemplate(), 0
	}
	root, err := git.Root()
	if err != nil {
		if _, count := parsed.Flags["--count"]; count {
			return "0\n", 0
		}
		return toon.NotInRepo() + "\n", 1
	}
	s := scanActive(root)
	if s.state.Failed() {
		return toon.RecordError(DecisionsDir, s.state, s.reason) + "\n", 1
	}
	if _, count := parsed.Flags["--count"]; count {
		return strconv.Itoa(s.count) + "\n", 0
	}
	out, err := toon.TableTyped("maps", []string{"map", "title", "type", "state", "blockers"}, s.rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	for _, row := range s.rows {
		if row[3] == "invalid" {
			return out, 1
		}
	}
	return out, 0
}
