package conformance

import "testing"

// The conformance-suite phase materializes only for a root that declares the entry
// point, and the filtered run then skips that one by contract, so this fixture carries
// the name and never executes it. Without the declaration the phase does not exist and
// the fixture grades nothing.
func TestRootConformance(t *testing.T) {}

// The failure the fixture exists for, deliberately under some other name: the entry
// point is skipped, so a fixture whose only failing test carried that name would go
// green and prove nothing.
func TestCanaryConformanceSuiteFails(t *testing.T) {
	t.Fatal("canary: intentional conformance-suite failure")
}
