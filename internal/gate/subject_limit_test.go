package gate

import (
	"os"
	"testing"
)

// TestManifestEntryLimitConstant pins the production entry-limit value. The runtime
// boundary proof runs at a lowered limit (BENCH_GATE_ENTRY_LIMIT) for speed, so this is
// the independent expectation that a change to the real ceiling turns red — without it,
// the cheap proof could pass while the shipped limit silently drifted.
func TestManifestEntryLimitConstant(t *testing.T) {
	if defaultManifestEntryLimit != 100000 {
		t.Fatalf("defaultManifestEntryLimit = %d, want 100000 (the shipped gate identity ceiling)", defaultManifestEntryLimit)
	}
}

// TestManifestEntryLimitOverrideIsTightenOnly pins the fail-safe: BENCH_GATE_ENTRY_LIMIT
// may only lower the ceiling. A value at or above the default, or a malformed one, is
// ignored — so the override can never raise the limit and can never enable a false green.
func TestManifestEntryLimitOverrideIsTightenOnly(t *testing.T) {
	cases := []struct {
		value string
		want  int
	}{
		{"", defaultManifestEntryLimit},
		{"10", 10},
		{"0", 0},
		{"99999", 99999},
		{"100000", defaultManifestEntryLimit},   // equal — ignored, cannot raise
		{"100001", defaultManifestEntryLimit},   // above — ignored
		{"-1", defaultManifestEntryLimit},       // negative — ignored
		{"nonsense", defaultManifestEntryLimit}, // malformed — ignored
	}
	for _, tc := range cases {
		t.Setenv("BENCH_GATE_ENTRY_LIMIT", tc.value)
		if tc.value == "" {
			os.Unsetenv("BENCH_GATE_ENTRY_LIMIT")
		}
		if got := manifestEntryLimit(); got != tc.want {
			t.Errorf("manifestEntryLimit() with %q = %d, want %d", tc.value, got, tc.want)
		}
	}
}
