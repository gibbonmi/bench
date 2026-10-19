package staleness

import (
	"maps"
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

// CommittedTicket is one ticket of a spec's tickets folder at one commit: its repository
// path and its parse.
type CommittedTicket struct {
	Path   string
	Parsed tickets.Ticket
}

// TicketEntriesAt is the ticket files of the spec's tickets folder at commit, read from
// committed objects through the one tracked-ticket enumeration, with its duplicate-basename
// diagnostics. A nested ticket counts, and a later copy of a basename does not, as in the
// checkout enumeration. Each reader of the ticket set of one commit reads this one answer. A
// tickets-folder member that is no regular file, or a ticket object that does not read,
// refuses with an *Error.
func TicketEntriesAt(root, commit, specPath string) ([]tickets.Entry, []string, error) {
	dir := ticketsDir(specPath)
	found, err := members(root, commit, []string{dir})
	if err != nil {
		return nil, nil, err
	}
	return tickets.EnumerateTracked(dir, slices.Collect(maps.Keys(found)), func(file string) ([]byte, error) {
		content, err := git.ReadControlBlob(root, found[file].ID)
		if err != nil {
			return nil, &Error{Stage: StageTicket, Subject: file, Cause: err}
		}
		return content, nil
	})
}

// TicketsAt is TicketEntriesAt parsed, with each ticket at its repository path.
func TicketsAt(root, commit, specPath string) ([]CommittedTicket, error) {
	entries, _, err := TicketEntriesAt(root, commit, specPath)
	if err != nil {
		return nil, err
	}
	dir := ticketsDir(specPath)
	set := make([]CommittedTicket, 0, len(entries))
	for _, entry := range entries {
		parsed, _ := tickets.ParseTicket(entry.Name, entry.Data, nil, "")
		set = append(set, CommittedTicket{Path: path.Join(dir, entry.Rel), Parsed: parsed})
	}
	return set, nil
}

// coverageStates is the coverage of the pair: one state for each coverage row that either
// spec maps, with its predicate bytes, and one state for each ticket file in paths that the
// base or the tip ticket set holds, with its Covers tokens. The coverage parser and the ticket
// grammar own both readings, so the audit splits no table and no field of its own.
func coverageStates(specPath string, baseline, current []byte, paths []PathState, base, tip []CommittedTicket) []CoverageState {
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
		beforeCovers, inBase := coversOf(base, state.Path)
		afterCovers, inTip := coversOf(tip, state.Path)
		if inBase || inTip {
			states = append(states, CoverageState{Path: state.Path, Kind: "ticket " + path.Base(state.Path) + " covers", Before: beforeCovers, After: afterCovers})
		}
	}
	return states
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

// coversOf is the Covers tokens of the ticket at ticketPath in set, in field order, and
// whether set holds that ticket. A set without the ticket answers Absent.
func coversOf(set []CommittedTicket, ticketPath string) (string, bool) {
	at := slices.IndexFunc(set, func(ticket CommittedTicket) bool { return ticket.Path == ticketPath })
	if at < 0 {
		return Absent, false
	}
	return strings.Join(set[at].Parsed.Covers, " "), true
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
