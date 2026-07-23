package gate

// The capability-skip collector: aggregation from a scripted log, the definitive
// zero state, and the cross-process concurrency the real phases produce.

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
)

func writeSkipLog(t *testing.T, skips ...capability.Skip) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "skips.log")
	var buf bytes.Buffer
	for _, skip := range skips {
		line, err := capability.Render(skip)
		if err != nil {
			t.Fatalf("render %#v: %v", skip, err)
		}
		buf.WriteString(line)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatalf("write skip log: %v", err)
	}
	return path
}

// TestSkipRowsReportEveryClass drives the collector directly against a scripted log,
// which pins the aggregation itself rather than the far weaker claim that some rows
// appeared. Classes report in the package's declared order and each carries its own
// count, so a collector that folded every class into one total, or dropped the class
// token, fails here.
func TestSkipRowsReportEveryClass(t *testing.T) {
	path := writeSkipLog(t,
		capability.Skip{Kind: capability.KindCapability, Class: capability.Symlink, Reason: "no unprivileged symlinks"},
		capability.Skip{Kind: capability.KindEnvironment, Reason: "subject root has no bin/bench.sh"},
		capability.Skip{Kind: capability.KindCapability, Class: capability.Privilege, Reason: "cannot drop privilege"},
		capability.Skip{Kind: capability.KindCapability, Class: capability.Symlink, Reason: "no unprivileged symlinks"},
	)

	want := []string{
		"capability-skips: 4 (capability=3 environment=1)",
		"capability-skips class=symlink: 2",
		"capability-skips class=privilege: 1",
	}
	if got := skipRows(readSkipTally(path)); !reflect.DeepEqual(got, want) {
		t.Fatalf("skipRows = %#v, want %#v", got, want)
	}
}

// TestSkipRowsStateZeroExplicitly proves silence and zero stay distinguishable: a
// fully capable host and a collector that silently broke would otherwise read alike.
// Both an untouched log and a log no phase ever created report the same explicit zero.
func TestSkipRowsStateZeroExplicitly(t *testing.T) {
	want := []string{"capability-skips: 0 (capability=0 environment=0)"}
	for name, path := range map[string]string{
		"empty log":  writeSkipLog(t),
		"absent log": filepath.Join(t.TempDir(), "never-created.log"),
		"noise only": writeNoiseLog(t),
	} {
		if got := skipRows(readSkipTally(path)); !reflect.DeepEqual(got, want) {
			t.Fatalf("%s: skipRows = %#v, want %#v", name, got, want)
		}
	}
}

func writeNoiseLog(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "noise.log")
	if err := os.WriteFile(path, []byte("ok\tinternal/gate\t0.4s\nPASS\n"), 0o644); err != nil {
		t.Fatalf("write noise log: %v", err)
	}
	return path
}

// TestCapabilitySkipsCountEveryConcurrentPhase runs the real phase runner with phases
// that all append to the one log the gate hands them. The concurrency this exercises
// is cross-process at write time — separate phase processes appending to a shared
// path, which is what the real `go test` phases do — while the read happens once,
// single-threaded, after the phases join. A collector that gave each phase its own
// destination, or read back only one of them, loses counts on exactly the multi-phase
// runs that matter and fails here.
func TestCapabilitySkipsCountEveryConcurrentPhase(t *testing.T) {
	const phases, linesPerPhase = 6, 40

	symlink, err := capability.Render(capability.Skip{Kind: capability.KindCapability, Class: capability.Symlink, Reason: "no unprivileged symlinks"})
	if err != nil {
		t.Fatalf("render symlink skip: %v", err)
	}
	environment, err := capability.Render(capability.Skip{Kind: capability.KindEnvironment, Reason: "fixture not materialized"})
	if err != nil {
		t.Fatalf("render environment skip: %v", err)
	}

	table := make([]Phase, 0, phases)
	for i := 0; i < phases; i++ {
		table = append(table, Phase{
			Name: "appender" + strconv.Itoa(i),
			Argv: []string{"bash", "-c", `for _ in $(seq "$2"); do printf '%s' "$1" >> "$BENCH_SKIP_LOG"; done`, "bash", symlink + environment, strconv.Itoa(linesPerPhase)},
		})
	}

	var stdout, stderr bytes.Buffer
	if code := runPhases(context.Background(), t.TempDir(), table, outerMode, &stdout, &stderr); code != 0 {
		t.Fatalf("runPhases = %d, want 0; stdout:\n%s\nstderr:\n%s", code, stdout.String(), stderr.String())
	}

	total := phases * linesPerPhase
	for _, want := range []string{
		fmt.Sprintf("capability-skips: %d (capability=%d environment=%d)", 2*total, total, total),
		fmt.Sprintf("capability-skips class=symlink: %d", total),
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout missing %q:\n%s", want, stdout.String())
		}
	}
}

// TestStrictCapabilityFailure pins the strict verdict to the capability population
// only, and to the exact flag value. An environment skip is a staging fact, not a
// security class, so it never contributes; any flag value other than "1" leaves the
// rows informational.
func TestStrictCapabilityFailure(t *testing.T) {
	mixed := readSkipTally(writeSkipLog(t,
		capability.Skip{Kind: capability.KindCapability, Class: capability.Fifo, Reason: "no fifo"},
		capability.Skip{Kind: capability.KindEnvironment, Reason: "fixture not materialized"},
	))
	environmentOnly := readSkipTally(writeSkipLog(t,
		capability.Skip{Kind: capability.KindEnvironment, Reason: "fixture not materialized"},
	))

	t.Setenv(requireCapabilitiesEnv, "1")
	failure := strictFailure(mixed)
	if !strings.Contains(failure, "fifo") {
		t.Fatalf("strictFailure = %q, want a message naming the fifo class", failure)
	}
	if got := strictFailure(environmentOnly); got != "" {
		t.Fatalf("strictFailure on environment skips = %q, want no failure", got)
	}

	for _, value := range []string{"", "0", "true", "yes"} {
		t.Setenv(requireCapabilitiesEnv, value)
		if got := strictFailure(mixed); got != "" {
			t.Fatalf("strictFailure with %s=%q = %q, want no failure", requireCapabilitiesEnv, value, got)
		}
	}
}
