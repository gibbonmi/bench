package worktree

import "testing"

// The race phase probes for this exact declaration before it materializes, and its
// -run filter exits 0 on a package that never declares it, so the fixture has to carry
// the real name. It fails outright rather than racing: the phase reds either way, and a
// deliberate data race would report nondeterministically.
func TestConcurrentCleanupRecordsOneTransaction(t *testing.T) {
	t.Fatal("canary: intentional race-phase failure")
}
