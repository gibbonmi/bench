package conformance

import "testing"

// Deliberately not named TestRootConformance: the filtered suite skips that one by
// contract (it is the entry point the phase must not re-enter), so a fixture whose only
// test carried that name would go green and prove nothing.
func TestCanaryConformanceSuiteFails(t *testing.T) {
	t.Fatal("canary: intentional conformance-suite failure")
}
