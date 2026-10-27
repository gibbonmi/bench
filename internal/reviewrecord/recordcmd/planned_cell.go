package recordcmd

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/reviewrecord"
)

// refusePlanned reads the spec and its plan at commit, a full commit ID, and refuses when a
// row that rows selects from that plan still holds a planned seam cell. It reads the
// commit, not the working copy, because the record names the commit. An empty output lets
// the write go on. The refusal comes before the write, so it writes nothing.
func refusePlanned(f form, root, spec, commit string, rows func(reviewrecord.Plan) []string) (string, int) {
	data, plan, err := reviewrecord.ReadSpecAtCommit(root, spec, commit)
	if err != nil {
		return f.refuse(err.Error())
	}
	selected := rows(plan)
	var held []string
	for _, id := range coverage.PlannedRows(data) {
		if slices.Contains(selected, id) {
			held = append(held, id)
		}
	}
	if len(held) == 0 {
		return "", 0
	}
	return f.refuse(fmt.Sprintf("rows %s hold a planned seam cell at %s; cite each test as %s", strings.Join(held, ", "), commit, coverage.CitationForm))
}

// chunkRows selects the planned rows of chunk id.
func chunkRows(id string) func(reviewrecord.Plan) []string {
	return inChunk(id, func(chunk *reviewrecord.PlannedChunk) []string { return chunk.Rows })
}

// ticketRows selects the rows that the Covers: line of ticket names in chunk id. A ticket that
// the chunk does not list selects no row.
func ticketRows(id, ticket string) func(reviewrecord.Plan) []string {
	return inChunk(id, func(chunk *reviewrecord.PlannedChunk) []string { return chunk.Covers[ticket] })
}

// inChunk selects the rows that pick answers for chunk id. A chunk that the plan does not name
// selects no row, and the writer then refuses it with its own cause.
func inChunk(id string, pick func(*reviewrecord.PlannedChunk) []string) func(reviewrecord.Plan) []string {
	return func(plan reviewrecord.Plan) []string {
		if chunk := plan.Chunk(id); chunk != nil {
			return pick(chunk)
		}
		return nil
	}
}
