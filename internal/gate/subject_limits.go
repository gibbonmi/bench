package gate

import (
	"os"
	"strconv"
)

// defaultManifestEntryLimit bounds declared-input fingerprinting.
const defaultManifestEntryLimit = 100000

// manifestEntryLimit accepts only a lower fail-safe test override.
func manifestEntryLimit() int {
	if v := os.Getenv("BENCH_GATE_ENTRY_LIMIT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 && n < defaultManifestEntryLimit {
			return n
		}
	}
	return defaultManifestEntryLimit
}
