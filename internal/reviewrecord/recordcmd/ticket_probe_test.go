package recordcmd_test

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// BC29: a ticket that owes no requirement writes the chunk step alone. The call exits 0, freezes
// the chunk at the tip, and prints an empty verifications cell.
func TestTicketFreezesWithoutRequirements(t *testing.T) {
	f, _ := pooled(t, 1, recordtest.Owing(nil))
	base, tip := advance(f, "chunk 1")
	out := succeed(t, f, ticketArgs(t, "1", "1.md", "p", "--base", base))
	if chunk := chunkEntry(t, f, "1"); chunk.Base != base || chunk.Tip != tip || len(chunk.Verification) != 0 {
		t.Errorf("chunk 1 = %s..%s with %d results, want %s..%s with none", chunk.Base, chunk.Tip, len(chunk.Verification), base, tip)
	}
	if want := `  "1",1.md,none,created,` + base + "," + tip + `,"",none` + "\n"; !strings.HasSuffix(out, want) {
		t.Errorf("ticket output = %q, want the row %q", out, want)
	}
}
