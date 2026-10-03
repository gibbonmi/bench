package worktree

// This file owns the serial ceiling check over the parallel census facts.

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

// serialSet returns one line for every serial test in facts, sorted, with the
// reason the census classified it serial.
func serialSet(facts []testFact) []string {
	var serial []string
	for _, fact := range facts {
		if fact.serialReason == "" {
			continue
		}
		serial = append(serial, fmt.Sprintf("%s:%d: %s (%s)", fact.file, fact.line, fact.name, fact.serialReason))
	}
	sort.Strings(serial)
	return serial
}

// serialCeilingBreach returns the refusal for a serial set above or below ceiling,
// and is empty when the set equals it. The refusal above lists the whole set with
// each reason, so the reader sees which test is new. The refusal below names the
// lower ceiling, so a test that leaves the serial set lowers the pin with it.
func serialCeilingBreach(facts []testFact, ceiling int) string {
	serial := serialSet(facts)
	switch {
	case len(serial) < ceiling:
		return belowCeilingRefusal(len(serial), ceiling)
	case len(serial) > ceiling:
		return fmt.Sprintf("the package holds %d serial tests, above the ceiling of %d:\n%s",
			len(serial), ceiling, strings.Join(serial, "\n"))
	}
	return ""
}

// belowCeilingRefusal renders the refusal for n serial tests below the ceiling c.
func belowCeilingRefusal(n, c int) string {
	return fmt.Sprintf("the package holds %d serial tests, below the ceiling of %d: lower worktreeSerialCeiling to %d in this change", n, c, n)
}

// TestCensusRefusesASerialSetBelowTheCeiling proves the ceiling check is exact: a
// test that leaves the serial set turns the check red until the ceiling drops.
func TestCensusRefusesASerialSetBelowTheCeiling(t *testing.T) {
	t.Parallel()
	dir := plantTestFiles(t, map[string]string{"serial_test.go": `package worktree

import "testing"

func TestBoundOne(t *testing.T) {
	bindEnv(t, "BENCH_HOME", "value")
}
`})
	facts, err := censusFacts(dir)
	if err != nil {
		t.Fatalf("census: %v", err)
	}
	breach, want := serialCeilingBreach(facts, 2), belowCeilingRefusal(1, 2)
	if breach == "" || breach != want {
		t.Fatalf("a set of one under a ceiling of two = %q, want %q", breach, want)
	}
}
