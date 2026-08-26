package gate

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
)

// writeSkipLog renders skips through the writer half of the line shape, so this test
// cannot pass against a reader that drifted from what capability.Render emits.
func writeSkipLog(t *testing.T, skips ...capability.Skip) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "skip.log")
	var buf bytes.Buffer
	for _, skip := range skips {
		line, err := capability.Render(skip)
		if err != nil {
			t.Fatalf("Render(%#v): %v", skip, err)
		}
		buf.WriteString(line)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// TestEnvironmentSkipIsRedInsideTheOracle pins the posture that makes the gate's own
// contracts enforceable. A test that found its staging absent has no verdict, and the
// gate is what stages it. Folding that fact into a green footer count leaves the
// conformance suite unenforced. The red is unconditional here: strict mode is unset, and
// the capability skip beside the environment skip is the control that proves the red
// belongs to the environment skip alone. The diagnosis names the test and the reason,
// because a reader who must hunt for which check stopped running needs a diagnosis, not
// a number. This test pins against a report that answers no diagnosis and against one
// that answers not-red.
func TestEnvironmentSkipIsRedInsideTheOracle(t *testing.T) {
	t.Setenv(requireCapabilitiesEnv, "")
	path := writeSkipLog(t,
		capability.Skip{Kind: capability.KindEnvironment, Name: "TestRootConformance", Reason: "BENCH_CONFORMANCE_ROOT not set"},
		capability.Skip{Kind: capability.KindCapability, Class: capability.Fifo, Name: "TestFifoRefusal", Reason: "requires a host fifo"},
	)
	var stdout bytes.Buffer
	rows, red := reportCapabilitySkips(path, &stdout)
	if !red {
		t.Fatal("an environment skip inside the oracle did not report red")
	}
	want := []string{"TestRootConformance: BENCH_CONFORMANCE_ROOT not set"}
	if !slices.Equal(rows, want) {
		t.Fatalf("diagnosis rows = %q, want %q", rows, want)
	}
	wantRow := "capability-skips class=environment: 1 (TestRootConformance: BENCH_CONFORMANCE_ROOT not set)"
	if !strings.Contains(stdout.String(), wantRow) {
		t.Fatalf("skip rows =\n%s\nwant a row %q", stdout.String(), wantRow)
	}
}

// TestCapabilitySkipsStayInformational holds the other half of the posture. A host that
// genuinely cannot make a security assertion is not the gate's staging failure, so those
// skips stay green outside strict mode and enter no diagnosis row. A blanket red on any
// skip turns this test red, and so does a report that files an informational skip as a
// failure the reader is asked to fix.
func TestCapabilitySkipsStayInformational(t *testing.T) {
	t.Setenv(requireCapabilitiesEnv, "")
	path := writeSkipLog(t,
		capability.Skip{Kind: capability.KindCapability, Class: capability.Fifo, Name: "TestFifoRefusal", Reason: "requires a host fifo"},
	)
	var stdout bytes.Buffer
	rows, red := reportCapabilitySkips(path, &stdout)
	if red {
		t.Fatalf("capability-only skips reported red: %q", rows)
	}
	if len(rows) != 0 {
		t.Fatalf("diagnosis rows = %q, want none", rows)
	}
	if want := "capability-skips: 1 (capability=1 environment=0)"; !strings.Contains(stdout.String(), want) {
		t.Fatalf("skip rows =\n%s\nwant %q", stdout.String(), want)
	}
	if strings.Contains(stdout.String(), "class=environment") {
		t.Fatalf("skip rows =\n%s\nwant no environment row when none were tallied", stdout.String())
	}
}
