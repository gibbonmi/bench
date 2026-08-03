package gate

// Shared verdict-record fixtures for the loader and reuse tests.

import (
	"strings"
	"time"
)

const verdictTestTree = "0123456789abcdef0123456789abcdef01234567"

var verdictTestOracle = strings.Repeat("a", 64)

var (
	partialTestIdentity = strings.Repeat("b", 64)
	partialTestSeal     = strings.Repeat("c", 64)
)

// partialTestRecord carries both evidence forms at once — an ancestor slot for a scoped
// component and a reused seal for build — so a case that drops one form is refused by the
// same record every other case starts from.
func partialTestRecord(now time.Time) verdictRecord {
	return verdictRecord{
		Schema:     verdictSchema,
		State:      Ready,
		Status:     "green",
		Tree:       verdictTestTree,
		Oracle:     verdictTestOracle,
		RecordedAt: now.Format(time.RFC3339),
		Executed:   []string{"conformance", "conformance-suite"},
		Skipped:    []string{"build", "vet"},
		SkipEvidence: map[string]skipEvidence{
			"build": {Seal: partialTestSeal},
			"vet":   {Identity: partialTestIdentity, AuthoredAt: now.Add(-90 * time.Minute).Format(time.RFC3339)},
		},
	}
}

func fullTestRecord(now time.Time) verdictRecord {
	return verdictRecord{
		Schema:     verdictSchema,
		State:      Ready,
		Status:     "green",
		Tree:       verdictTestTree,
		Oracle:     verdictTestOracle,
		RecordedAt: now.Format(time.RFC3339),
	}
}
