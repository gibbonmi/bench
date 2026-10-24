package gate

import (
	"testing"

	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// TestChunkCheckpointRefusesADirtyCheckout covers RC44: the chunk checkpoint refuses a working
// tree that differs from HEAD before the evidence pass and before the oracle. The review
// record has no exemption, so a valid record change alone also refuses.
func TestChunkCheckpointRefusesADirtyCheckout(t *testing.T) {
	for _, tc := range []struct {
		name  string
		dirty func(*testing.T, *recordtest.Fixture)
	}{
		{"tracked edit", func(_ *testing.T, f *recordtest.Fixture) { f.Write("tracked.txt", "edited\n") }},
		{"uncommitted record", annotateRecord},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := checkpointFixture(t)
			tc.dirty(t, f)
			code, out := runCheckpoint(t, f)
			requireCleanCheckoutRefusal(t, f, code, out, "stale reviewed source")
		})
	}
}
