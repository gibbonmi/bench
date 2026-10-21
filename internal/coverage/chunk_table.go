package coverage

import (
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/tickets"
)

// chunksHeading opens the spec section whose first pipe table is the chunk table.
const chunksHeading = "## Implementation chunks"

// acceptanceHeader is the lowercased header cell of the chunk table's acceptance column.
const acceptanceHeader = "acceptance rows"

// The two pinned details of a chunk table that does not read.
const (
	noChunkTable       = "no Implementation chunks table"
	noAcceptanceColumn = "chunk table has no acceptance rows column"
)

var (
	// chunksEndRe is the next `##` or `###` heading, which closes the chunks section.
	chunksEndRe = regexp.MustCompile(`^#{2,3} `)
	// rowRangeRe is one range token: two row IDs joined by a hyphen or an en dash. Its
	// submatches are the first ID, its tag, the last ID, and its tag.
	rowRangeRe = regexp.MustCompile(`^(` + tickets.RowIDPattern + `)(?:-|–)(` + tickets.RowIDPattern + `)$`)
)

// ChunkTable is the acceptance column of a spec's Implementation chunks table, read against
// the spec's own coverage map. Missing is the pinned detail of an absent table or an absent
// acceptance column, and it is empty when both read. Listed is every row ID the acceptance
// cells name, in table order, with each range expanded. Faults names each token that lists no
// row.
type ChunkTable struct {
	Missing string
	Listed  []string
	Faults  []string
}

// ReadChunkTable reads the first pipe table under the Implementation chunks heading of content,
// before the next `##` or `###` heading. It finds the acceptance column by its header cell at
// any position. It splits each acceptance cell on commas and reads each token against the
// declared row IDs of the coverage map in content.
func ReadChunkTable(content []byte) ChunkTable {
	declared := rowIDs(parse(content))
	var table ChunkTable
	column, inSection, inTable := -1, false, false
	for _, raw := range strings.Split(string(content), "\n") {
		line := strings.TrimSuffix(raw, "\r")
		if !inSection {
			inSection = line == chunksHeading
			continue
		}
		if chunksEndRe.MatchString(line) {
			break
		}
		cells, ok := pipeCells(line)
		if !ok {
			if inTable {
				break
			}
			continue
		}
		if !inTable {
			inTable = true
			column = slices.IndexFunc(cells, func(c string) bool { return strings.ToLower(c) == acceptanceHeader })
			if column < 0 {
				return ChunkTable{Missing: noAcceptanceColumn}
			}
			continue
		}
		if separatorRow(cells) || column >= len(cells) {
			continue
		}
		table.readCell(cells[column], declared)
	}
	if !inTable {
		return ChunkTable{Missing: noChunkTable}
	}
	return table
}

// readCell reads one acceptance cell. It trims spaces and backticks from each comma-separated
// token and skips an empty token. A row ID names one row, and a range names its endpoints and
// each declared row of their tag between them. Any other token is a fault.
func (t *ChunkTable) readCell(cell string, declared []string) {
	for _, token := range strings.Split(cell, ",") {
		token = strings.Trim(token, " \t`")
		switch {
		case token == "":
		case rowIDRe.MatchString(token):
			t.Listed = append(t.Listed, token)
		case rowRangeRe.MatchString(token):
			rows, fault := expandRange(token, declared)
			t.Listed = append(t.Listed, rows...)
			if fault != "" {
				t.Faults = append(t.Faults, fault)
			}
		default:
			t.Faults = append(t.Faults, "token "+token+" is no row ID or range")
		}
	}
}

// expandRange lists the declared rows that a range token names, or the fault that names the
// token: two tags, a first number above the last, or an undeclared endpoint.
func expandRange(token string, declared []string) ([]string, string) {
	m := rowRangeRe.FindStringSubmatch(token)
	first, last := m[1], m[3]
	if m[2] != m[4] {
		return nil, "range " + token + " spans two tags"
	}
	low, high := rowNumber(first), rowNumber(last)
	if low > high {
		return nil, "range " + token + " runs from a higher number to a lower one"
	}
	for _, endpoint := range []string{first, last} {
		if !slices.Contains(declared, endpoint) {
			return nil, "range " + token + " names the undeclared endpoint " + endpoint
		}
	}
	var rows []string
	for _, id := range declared {
		if n := rowNumber(id); tickets.TagOf(id) == m[2] && n >= low && n <= high {
			rows = append(rows, id)
		}
	}
	return rows, ""
}

// rowNumber is the number half of a row ID that matches the row-ID grammar.
func rowNumber(id string) int {
	n, _ := strconv.Atoi(strings.TrimPrefix(id, tickets.TagOf(id)))
	return n
}
