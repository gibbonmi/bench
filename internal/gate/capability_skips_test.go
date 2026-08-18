package gate

import (
	"bytes"
	"os"
	"path/filepath"
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
// contracts enforceable: a test that found its staging absent has no verdict, and the
// gate is what stages it, so folding that into a green footer count is what let the
// conformance suite go unenforced. The red message names the test and the reason —
// a reader who has to go hunt for which check stopped running got a number, not a
// diagnosis. Making environmentFailure return "" is the mutant this pins against.
func TestEnvironmentSkipIsRedInsideTheOracle(t *testing.T) {
	t.Setenv(requireCapabilitiesEnv, "")
	path := writeSkipLog(t,
		capability.Skip{Kind: capability.KindEnvironment, Name: "TestRootConformance", Reason: "BENCH_CONFORMANCE_ROOT not set"},
		capability.Skip{Kind: capability.KindCapability, Class: capability.Fifo, Name: "TestFifoRefusal", Reason: "requires a host fifo"},
	)
	var stdout, stderr bytes.Buffer
	if !reportCapabilitySkips(path, &stdout, &stderr) {
		t.Fatal("an environment skip inside the oracle did not report red")
	}
	for _, want := range []string{"TestRootConformance", "BENCH_CONFORMANCE_ROOT not set"} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("red message %q does not name %q", stderr.String(), want)
		}
	}
	wantRow := "capability-skips class=environment: 1 (TestRootConformance: BENCH_CONFORMANCE_ROOT not set)"
	if !strings.Contains(stdout.String(), wantRow) {
		t.Fatalf("skip rows =\n%s\nwant a row %q", stdout.String(), wantRow)
	}
}

// TestCapabilitySkipsStayInformational holds the other half of the posture: a host that
// genuinely cannot make a security assertion is not the gate's staging failure, so those
// rows stay green outside strict mode. A blanket red on any skip would turn this red.
func TestCapabilitySkipsStayInformational(t *testing.T) {
	t.Setenv(requireCapabilitiesEnv, "")
	path := writeSkipLog(t,
		capability.Skip{Kind: capability.KindCapability, Class: capability.Fifo, Name: "TestFifoRefusal", Reason: "requires a host fifo"},
	)
	var stdout, stderr bytes.Buffer
	if reportCapabilitySkips(path, &stdout, &stderr) {
		t.Fatalf("capability-only skips reported red: %s", stderr.String())
	}
	if want := "capability-skips: 1 (capability=1 environment=0)"; !strings.Contains(stdout.String(), want) {
		t.Fatalf("skip rows =\n%s\nwant %q", stdout.String(), want)
	}
	if strings.Contains(stdout.String(), "class=environment") {
		t.Fatalf("skip rows =\n%s\nwant no environment row when none were tallied", stdout.String())
	}
}
