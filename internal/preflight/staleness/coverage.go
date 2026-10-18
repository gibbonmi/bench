package staleness

import (
	"path"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/tickets"
)

// ticketsDir is the tickets folder of the spec at specPath.
func ticketsDir(specPath string) string {
	return path.Join(path.Dir(specPath), "tickets")
}

// coverageStates is the coverage of the pair: one state for each coverage row that either
// spec maps, with its predicate bytes, and one state for each ticket file in paths, with its
// Covers tokens. Another file in the tickets folder is an asset that the ticket grammar
// ignores. The coverage parser and the ticket grammar own both readings, so the audit splits
// no table and no field of its own. A ticket object that does not read refuses the
// collection.
func coverageStates(root, specPath string, baseline, current []byte, paths []PathState) ([]CoverageState, error) {
	baselineRows, currentRows := coverage.MapRows(baseline), coverage.MapRows(current)
	before, after := predicates(baselineRows), predicates(currentRows)
	var ids []string
	for _, row := range append(currentRows, baselineRows...) {
		if !slices.Contains(ids, row.ID) {
			ids = append(ids, row.ID)
		}
	}
	var states []CoverageState
	for _, id := range ids {
		states = append(states, CoverageState{Path: specPath, Kind: "coverage " + id + " predicate", Before: side(before, id), After: side(after, id)})
	}
	for _, state := range paths {
		if !tickets.Covers(ticketsDir(specPath), state.Path) || path.Ext(state.Path) != tickets.Ext {
			continue
		}
		covers := CoverageState{Path: state.Path, Kind: "ticket " + path.Base(state.Path) + " covers"}
		var err error
		if covers.Before, err = coversAt(root, state.Path, state.Before); err != nil {
			return nil, err
		}
		if covers.After, err = coversAt(root, state.Path, state.After); err != nil {
			return nil, err
		}
		states = append(states, covers)
	}
	return states, nil
}

// predicates is the behavior cell of each coverage row, by row ID.
func predicates(rows []coverage.MapRow) map[string]string {
	cells := map[string]string{}
	for _, row := range rows {
		cells[row.ID] = row.Behavior
	}
	return cells
}

// side is the value that values holds for id, or Absent.
func side(values map[string]string, id string) string {
	if value, ok := values[id]; ok {
		return value
	}
	return Absent
}

// coversAt is the Covers tokens of the ticket object at one identity, in field order, or
// Absent for the zero object.
func coversAt(root, ticketPath string, object Object) (string, error) {
	if object.ID == "" {
		return Absent, nil
	}
	content, err := git.ReadControlBlob(root, object.ID)
	if err != nil {
		return "", &Error{Stage: StageTicket, Subject: ticketPath, Cause: err}
	}
	parsed, _ := tickets.ParseTicket(path.Base(ticketPath), content, nil, "")
	return strings.Join(parsed.Covers, " "), nil
}

// coverageEntry is the drift entry of one changed coverage subject: its path on each side that
// holds the subject, and the subject's value at each identity as its object.
func coverageEntry(state CoverageState) DriftEntry {
	e := DriftEntry{Before: Absent, After: Absent, Kind: state.Kind, BeforeObject: state.Before, AfterObject: state.After}
	if state.Before != Absent {
		e.Before = state.Path
	}
	if state.After != Absent {
		e.After = state.Path
	}
	return e
}
