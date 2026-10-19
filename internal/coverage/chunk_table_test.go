package coverage

import (
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// chunkSpec renders a spec whose coverage map declares ids, followed by chunks and a later
// section that closes the chunks section.
func chunkSpec(ids []string, chunks string) []byte {
	body := "# c\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReducedID
	for _, id := range ids {
		body += "| " + id + " | 1 | b | s | w |\n"
	}
	return []byte(body + chunks + "\n## Testing decisions\n")
}

// acceptanceSpec is chunkSpec with the shared fixture chunk table whose acceptance cell is cell.
func acceptanceSpec(ids []string, cell string) []byte {
	return chunkSpec(ids, preflighttest.ChunkTable(cell))
}

// requireListed fails unless table lists exactly want with no fault and no missing detail.
func requireListed(t *testing.T, table ChunkTable, want ...string) {
	t.Helper()
	if table.Missing != "" || len(table.Faults) != 0 || !slices.Equal(table.Listed, want) {
		t.Fatalf("chunk table = %+v, want exactly %v with no fault", table, want)
	}
}

// TestChunkTableExpandsRanges covers PC3, PC4, and PC37. A hyphen or an en dash range
// names its endpoints and each declared row of its tag between them, and only those.
func TestChunkTableExpandsRanges(t *testing.T) {
	for _, test := range []struct {
		name     string
		declared []string
		cell     string
		want     []string
	}{
		{"PC3 hyphen range", []string{"PF1", "PF2", "PF3", "PF4"}, "PF1, PF2-PF4", []string{"PF1", "PF2", "PF3", "PF4"}},
		{"PC4 en dash range", []string{"PF1", "PF2", "PF3", "PF4"}, "PF1, PF2–PF4", []string{"PF1", "PF2", "PF3", "PF4"}},
		{"PC37 range over an omitted row", []string{"PF1", "PF2", "PF4"}, "PF1, PF2–PF4", []string{"PF1", "PF2", "PF4"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			requireListed(t, ReadChunkTable(acceptanceSpec(test.declared, test.cell)), test.want...)
		})
	}
}

// requireFault fails unless table lists no row and holds one fault that names token.
func requireFault(t *testing.T, table ChunkTable, token string) {
	t.Helper()
	if table.Missing != "" || len(table.Listed) != 0 || len(table.Faults) != 1 || !strings.Contains(table.Faults[0], token) {
		t.Fatalf("chunk table = %+v, want one fault naming %q and no listed row", table, token)
	}
}

// TestChunkTableRangeFaults covers PC5, PC6, and PC7. Each bad range is a fault that names
// its token, and it lists no row. The map declares both endpoints of the two-tag and the
// reversed range, so only their own rule can fault them.
func TestChunkTableRangeFaults(t *testing.T) {
	declared := []string{"PF1", "PF2", "PF3", "QF3"}
	for _, test := range []struct{ name, token string }{
		{"PC5 undeclared endpoint", "PF2-PF9"},
		{"PC6 two tags", "PF1-QF3"},
		{"PC7 reversed bounds", "PF3–PF1"},
	} {
		t.Run(test.name, func(t *testing.T) {
			requireFault(t, ReadChunkTable(acceptanceSpec(declared, test.token)), test.token)
		})
	}
}

// TestChunkTableTokenFaults covers PC8 and PC38. A token that is neither a row ID nor a range
// is a fault that names it, and the token rule trims backticks and skips an empty token.
func TestChunkTableTokenFaults(t *testing.T) {
	declared := []string{"PF1", "PF2"}
	t.Run("PC8 free text", func(t *testing.T) {
		requireFault(t, ReadChunkTable(acceptanceSpec(declared, "all rows")), "all rows")
	})
	t.Run("PC38 backticks and a trailing comma", func(t *testing.T) {
		requireListed(t, ReadChunkTable(acceptanceSpec(declared, "`PF1`, PF2,")), "PF1", "PF2")
	})
}

// TestChunkTableFindsTheAcceptanceColumn covers PC10: the reader finds the acceptance column
// by its header cell in any case and at any position.
func TestChunkTableFindsTheAcceptanceColumn(t *testing.T) {
	spec := chunkSpec([]string{"PF1", "PF2"}, "\n## Implementation chunks\n\n| chunk | tickets | outcome | blocked by | Acceptance rows | tests |\n"+
		"|---|---|---|---|---|---|\n| C1 | one.md | done | none | PF1, PF2 | unit |\n")
	requireListed(t, ReadChunkTable(spec), "PF1", "PF2")
}
