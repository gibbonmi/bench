package ownership

import (
	"slices"

	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/tickets"
)

// UncitedIDs is every row ID of ids that no parsed ticket's `Covers:` line cites, in ids
// order. A row ID in ticket prose cites nothing. The rows-owned check and the breakdown both
// read this one rule over the one declared row-ID set.
func UncitedIDs(ids []string, parsed []tickets.Ticket) []string {
	var open []string
	for _, id := range ids {
		if !slices.ContainsFunc(parsed, func(ticket tickets.Ticket) bool { return slices.Contains(ticket.Covers, id) }) {
			open = append(open, id)
		}
	}
	return open
}

// uncited is the coverage parser row of each declared row ID that UncitedIDs names, in ids
// order. rows supply only each row's cells, so an ID with no parser row keeps its ID alone.
func uncited(ids []string, rows []coverage.MapRow, parsed []tickets.Ticket) []coverage.MapRow {
	var result []coverage.MapRow
	for _, id := range UncitedIDs(ids, parsed) {
		row := coverage.MapRow{ID: id}
		if i := slices.IndexFunc(rows, func(row coverage.MapRow) bool { return row.ID == id }); i >= 0 {
			row = rows[i]
		}
		result = append(result, row)
	}
	return result
}
