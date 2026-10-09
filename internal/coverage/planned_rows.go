package coverage

import (
	"fmt"
	"strings"
)

// citationFormat is the seam-cell form that cites one test: a backticked test file path
// and, in parentheses after it, the backticked test name that each caller supplies.
const citationFormat = "`<path>_test.go` (`%s`)"

// CitationForm is the citation form with a placeholder for the test name.
var CitationForm = fmt.Sprintf(citationFormat, "<Name>")

// PlannedRows returns, in map order, the row ID of each opted-in map row whose seam cell
// holds the planned marker. A map without row IDs gives no row ID, because a positional
// row number cannot join a ticket's `Covers:` token. A review-owned cell is never a
// planned row: its marker is a prefix, and the planned marker can match a later word.
func PlannedRows(spec []byte) []string {
	p := parse(spec)
	var planned []string
	for i, id := range rowIDs(p) {
		cell := strings.TrimSpace(p.sch.cell(p.dataRows[i], fieldSeam))
		if !strings.HasPrefix(cell, reviewOwnedPrefix) && plannedRe.MatchString(cell) {
			planned = append(planned, id)
		}
	}
	return planned
}
