package ownership

import (
	"slices"

	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/tickets"
)

// UncitedIDs is every row ID of ids that no parsed ticket's `Covers:` line cites, in ids
// order. A row ID in ticket prose cites nothing. The rows-owned check and the breakdown both
// read this one rule.
func UncitedIDs(ids []string, parsed []tickets.Ticket) []string {
	var open []string
	for _, id := range ids {
		if !slices.ContainsFunc(parsed, func(ticket tickets.Ticket) bool { return slices.Contains(ticket.Covers, id) }) {
			open = append(open, id)
		}
	}
	return open
}

// uncited is every coverage parser row whose row ID UncitedIDs names, in map order.
func uncited(rows []coverage.MapRow, parsed []tickets.Ticket) []coverage.MapRow {
	ids := make([]string, len(rows))
	for i, row := range rows {
		ids[i] = row.ID
	}
	open := UncitedIDs(ids, parsed)
	var result []coverage.MapRow
	for _, row := range rows {
		if slices.Contains(open, row.ID) {
			result = append(result, row)
		}
	}
	return result
}
