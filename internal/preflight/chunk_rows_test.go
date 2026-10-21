package preflight

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// listing is the chunk table whose acceptance cells list exactly ids.
func listing(ids ...string) coverage.ChunkTable { return coverage.ChunkTable{Listed: ids} }

// chunkRow decides baseFacts with the chunk table that the coverage reader reads from the
// shared spec, whose map declares PF1 and PF2, after its chunk table becomes table.
func chunkRow(t *testing.T, table string) CheckResult {
	t.Helper()
	f := baseFacts()
	f.ChunkTable = coverage.ReadChunkTable([]byte(preflighttest.WithChunkTable(preflighttest.SpecBody("example"), table)))
	row, ok := checkRow(Decide(f), chunkRowsName)
	if !ok {
		t.Fatalf("Decide renders no %s row", chunkRowsName)
	}
	return row
}

// requireChunkRed fails unless row is the red chunk-rows row whose detail holds each name.
func requireChunkRed(t *testing.T, row CheckResult, names ...string) {
	t.Helper()
	if row.Verdict != verdictRed {
		t.Fatalf("chunk-rows = %+v, want red", row)
	}
	for _, name := range names {
		if !strings.Contains(row.Detail, name) {
			t.Errorf("chunk-rows detail %q does not name %s", row.Detail, name)
		}
	}
}

// TestChunkRowsRedsAnUnlistedRow covers PC1: a declared row that no acceptance cell lists
// reds the row, and the detail names that row.
func TestChunkRowsRedsAnUnlistedRow(t *testing.T) {
	requireChunkRed(t, chunkRow(t, preflighttest.ChunkTable("PF1")), "PF2")
}

// TestChunkRowsRedsAnUndeclaredID covers PC2: an acceptance-cell ID that the map does not
// declare reds the row, and the detail names that ID.
func TestChunkRowsRedsAnUndeclaredID(t *testing.T) {
	requireChunkRed(t, chunkRow(t, preflighttest.ChunkTable("PF1", "PF2", "PF9")), "PF9")
}

// TestChunkRowsAcceptsARepeatedRow covers PC9: a row that two chunk rows list stays green.
func TestChunkRowsAcceptsARepeatedRow(t *testing.T) {
	table := preflighttest.ChunkTable("PF1", "PF2") + "| c2 / two.md | PF2 |\n"
	if row := chunkRow(t, table); row.Verdict != verdictGreen {
		t.Fatalf("chunk-rows over a repeated row = %+v, want green", row)
	}
}

// TestChunkRowsRedsAnAbsentTable covers PC11 and PC12: an absent table and an absent
// acceptance column each red the row with its own detail.
func TestChunkRowsRedsAnAbsentTable(t *testing.T) {
	for _, test := range []struct{ name, table, detail string }{
		{"PC11 absent table", "", "no Implementation chunks table"},
		{"PC12 absent column", "\n## Implementation chunks\n\n| chunk | rows |\n|---|---|\n| c1 | PF1, PF2 |\n", "chunk table has no acceptance rows column"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if row := chunkRow(t, test.table); row.Verdict != verdictRed || row.Detail != test.detail {
				t.Fatalf("chunk-rows = %+v, want red with the detail %q", row, test.detail)
			}
		})
	}
}

// TestChunkRowsNotApplicableWithoutTickets covers PC13: build mode with no tickets folder
// reads the row not-applicable with no detail, even over an absent table.
func TestChunkRowsNotApplicableWithoutTickets(t *testing.T) {
	f := baseFacts()
	f.Mode, f.TicketsDirExists, f.ChunkTable = modeBuild, false, coverage.ChunkTable{}
	if row, ok := checkRow(Decide(f), chunkRowsName); !ok || row.Verdict != verdictNA || row.Detail != "" {
		t.Fatalf("chunk-rows without tickets = %+v (rendered %t), want not-applicable with no detail", row, ok)
	}
}

// TestChunkRowsRefusesTheBuildCharge covers PC15: a red chunk-rows row refuses the build
// charge, and the refusal names the row.
func TestChunkRowsRefusesTheBuildCharge(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	spec := preflighttest.WithChunkTable(preflighttest.SpecBody(slug), preflighttest.ChunkTable("PF1"))
	args := preflighttest.CommitChargeSource(t, root, slug, map[string]string{"specs/" + slug + "/spec.md": spec})
	out, code := Command(args)
	if want := "error: preflight required: chunk-rows: "; code != 1 || !strings.Contains(out, want) || !strings.Contains(out, "PF2") {
		t.Fatalf("build charge = (%d):\n%s\nwant the refusal %q naming PF2", code, out, want)
	}
}
