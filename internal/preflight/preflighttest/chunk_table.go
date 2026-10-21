package preflighttest

import "strings"

// specRows are the coverage-map rows that specBody declares, in map order: the row ID, the
// behavior cell, and the why cell.
var specRows = [][3]string{{"PF1", "does x", "catches z"}, {"PF2", "does y", "catches w"}}

// specRowIDs are the row IDs of specRows, in map order.
func specRowIDs() []string {
	ids := make([]string, len(specRows))
	for i, row := range specRows {
		ids[i] = row[0]
	}
	return ids
}

// ChunkTable renders an Implementation chunks section whose one chunk lists rows in its
// acceptance cell. Every spec fixture that reaches a graded verdict renders its table here.
func ChunkTable(rows ...string) string {
	return "\n## Implementation chunks\n\n| stable chunk ID / tickets | acceptance rows |\n|---|---|\n| c1 / one.md | " +
		strings.Join(rows, ", ") + " |\n"
}

// WithChunkTable replaces the chunk table that SpecBody renders in spec with table. An empty
// table removes the section.
func WithChunkTable(spec, table string) string {
	return strings.Replace(spec, ChunkTable(specRowIDs()...), table, 1)
}
