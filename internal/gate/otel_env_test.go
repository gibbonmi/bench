package gate

import (
	"strings"
	"testing"
)

// TestGateEnvStripsTheRecordVariables covers the record's env pair at the parent side.
// The two variables address the gate run that started this process. A phase child that
// inherited them from the operator's shell would attach its lines to a run it is not
// part of, so the composer sets them on the one child that runs the phase table and the
// stripper removes every inherited value first.
func TestGateEnvStripsTheRecordVariables(t *testing.T) {
	t.Setenv(otelRootEnv, "/some/other/repository")
	t.Setenv(otelTraceparentEnv, "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01")

	composed, err := gateEnv()
	if err != nil {
		t.Fatalf("gate environment: %v", err)
	}
	for _, name := range otelGateEnv {
		for _, entry := range composed {
			if strings.HasPrefix(entry, name+"=") {
				t.Errorf("the phases environment kept %q", entry)
			}
		}
	}
}
