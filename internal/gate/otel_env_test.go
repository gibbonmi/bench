package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/otelrecord"
)

// TestGateEnvStripsTheRecordVariables covers the record's env pair at the parent side.
// The two variables address the gate run that started this process. A phase child that
// inherited them from the operator's shell would attach its lines to a run it is not
// part of, so the composer sets them on the one child that runs the phase table and the
// stripper removes every inherited value first.
func TestGateEnvStripsTheRecordVariables(t *testing.T) {
	names := otelrecord.HandoffVariables()
	t.Setenv(names[0], "/some/other/repository")
	t.Setenv(names[1], "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01")

	composed, err := gateEnv()
	if err != nil {
		t.Fatalf("gate environment: %v", err)
	}
	for _, name := range names {
		for _, entry := range composed {
			if strings.HasPrefix(entry, name+"=") {
				t.Errorf("the phases environment kept %q", entry)
			}
		}
	}
}

// TestTheGateSpellsNoHandoffVariable holds row LE12: the record package owns the handoff
// names, so a surviving copy in a non-test gate file reds the literal search. The literal
// is built from two parts, so this file does not match itself.
func TestTheGateSpellsNoHandoffVariable(t *testing.T) {
	literal := "BENCH_" + "OTEL_"
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("list the gate package: %v", err)
	}
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		raw, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}
		if strings.Contains(string(raw), literal) {
			t.Errorf("%s spells a handoff variable; call the record package's handoff instead", file)
		}
	}
}
