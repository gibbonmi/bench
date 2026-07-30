package canary

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/gibbonmi/bench/internal/subprocess"
)

const canaryAbortHelperEnv = "BENCH_CANARY_ABORT_HELPER"

func TestSweepAttributesPanickingTest(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "runtime", "root-panic")

	err := Sweep(root, sweepAbortRunner(authenticGoAbortOutput(t, "root")))
	if err == nil {
		t.Fatal("Sweep err = nil, want the panicking test reported")
	}
	got := err.Error()
	if !strings.Contains(got, "inner test abort") || !strings.Contains(got, "TestSubjectHash") || strings.Contains(got, "did not bite") {
		t.Fatalf("Sweep err = %q, want an inner test abort naming TestSubjectHash", got)
	}
}

func TestSweepAttributesPanickingSubtest(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "runtime", "subtest-panic")

	err := Sweep(root, sweepAbortRunner(authenticGoAbortOutput(t, "subtest")))
	if err == nil {
		t.Fatal("Sweep err = nil, want the panicking subtest reported")
	}
	got := err.Error()
	if !strings.Contains(got, "inner test abort") || !strings.Contains(got, "TestSubjectHash/empty_subject/hash") || strings.Contains(got, "did not bite") {
		t.Fatalf("Sweep err = %q, want an inner test abort naming the deepest subtest", got)
	}
}

func TestSweepAttributesPanickingRepeatedSubtest(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "runtime", "repeated-subtest-panic")

	err := Sweep(root, sweepAbortRunner(authenticGoAbortOutput(t, "repeated-subtest")))
	if err == nil {
		t.Fatal("Sweep err = nil, want the panicking repeated subtest reported")
	}
	got := err.Error()
	if !strings.Contains(got, "inner test abort") || !strings.Contains(got, "TestSubjectHash/hash#01") || strings.Contains(got, "did not bite") {
		t.Fatalf("Sweep err = %q, want an inner test abort naming the repeated subtest", got)
	}
}

func TestSweepAttributesPanickingPunctuatedSubtest(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "runtime", "punctuated-subtest-panic")

	err := Sweep(root, sweepAbortRunner(authenticGoAbortOutput(t, "punctuated-subtest")))
	if err == nil {
		t.Fatal("Sweep err = nil, want the panicking punctuated subtest reported")
	}
	got := err.Error()
	if !strings.Contains(got, "inner test abort") || !strings.Contains(got, "TestSubjectHash/empty-subject") || strings.Contains(got, "did not bite") {
		t.Fatalf("Sweep err = %q, want an inner test abort naming the punctuated subtest", got)
	}
}

func TestSweepAbortGrammar(t *testing.T) {
	cases := []struct {
		name      string
		result    RunResult
		wantAbort bool
		want      string
	}{
		{
			name: "outer phase prefix",
			result: RunResult{ExitCode: 2, Output: strings.Join([]string{
				"[contract] --- FAIL: TestPrefixedPanic (0.00s)",
				"[contract] panic: prefixed panic",
			}, "\n")},
			wantAbort: true,
		},
		{
			name:      "incidental diagnostic",
			result:    RunResult{ExitCode: 1, Output: "diagnostic mentions panic: but the test completed\n"},
			wantAbort: false,
			want:      `canary 'incidental diagnostic' did not bite (want red + "target-incidental diagnostic"; got exit 1)`,
		},
		{
			name: "green panic-shaped text",
			result: RunResult{ExitCode: 0, Output: strings.Join([]string{
				"--- FAIL: TestRecoveredPanic (0.00s)",
				"panic: text from a completed run",
			}, "\n")},
			wantAbort: false,
			want:      `canary 'green panic-shaped text' did not bite (want red + "target-green panic-shaped text"; got exit 0)`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			contractFixture(t, root, "runtime", tc.name)

			err := Sweep(root, sweepAbortResultRunner(tc.result))
			if err == nil {
				t.Fatal("Sweep err = nil, want the fixture reported")
			}
			got := err.Error()
			if tc.wantAbort {
				if !strings.Contains(got, "inner test abort") || !strings.Contains(got, "TestPrefixedPanic") || strings.Contains(got, "did not bite") {
					t.Fatalf("Sweep err = %q, want the prefixed Go panic attributed", got)
				}
				return
			}
			if got != tc.want {
				t.Fatalf("Sweep err = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestSweepAttributesTruncatedPanic(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "runtime", "truncated-panic")
	output := "panic: " + strings.Repeat("x", 240) + "\x1b[31m"

	err := Sweep(root, sweepAbortResultRunner(RunResult{ExitCode: 2, Output: output}))
	if err == nil {
		t.Fatal("Sweep err = nil, want the truncated panic reported")
	}
	got := err.Error()
	if !strings.Contains(got, "inner test abort") || !strings.Contains(got, "unknown test") || !strings.Contains(got, `contract package "runtime"`) || strings.Contains(got, "did not bite") {
		t.Fatalf("Sweep err = %q, want an attributed unknown-test abort", got)
	}
	if strings.ContainsRune(got, '\x1b') || !strings.Contains(got, "… (") || len(got) > 280 {
		t.Fatalf("Sweep err = %q, want a bounded control-safe diagnostic", got)
	}
}

func TestSweepAttributesRuntimeFatal(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "runtime", "runtime-fatal")
	output := authenticGoAbortOutput(t, "runtime-fatal")

	err := Sweep(root, sweepAbortResultRunner(RunResult{ExitCode: 2, Output: output}))
	if err == nil {
		t.Fatal("Sweep err = nil, want the runtime fatal reported")
	}
	want := `canary 'runtime-fatal' inner test abort in unknown test in contract package "runtime": fatal error: sync: unlock of unlocked mutex`
	if got := err.Error(); got != want {
		t.Fatalf("Sweep err = %q, want the bounded authentic runtime-fatal attribution %q", got, want)
	}
}

func TestSweepAttributesAbortAssociatedHeader(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "runtime", "competing-headers")
	output := strings.Join([]string{
		"[other] --- FAIL: TestWrongFirst (0.00s)",
		"[contract] --- FAIL: TestParent (0.00s)",
		"[contract]     --- FAIL: TestParent/aborting (0.00s)",
		"[other] --- FAIL: TestWrongNearest (0.00s)",
		"[contract] panic: nested panic",
		"[other] --- FAIL: TestWrongLast (0.00s)",
	}, "\n")

	err := Sweep(root, sweepAbortResultRunner(RunResult{ExitCode: 2, Output: output}))
	if err == nil {
		t.Fatal("Sweep err = nil, want the panic reported")
	}
	got := err.Error()
	if !strings.Contains(got, "TestParent/aborting") || strings.Contains(got, "TestWrongFirst") || strings.Contains(got, "TestWrongNearest") || strings.Contains(got, "TestWrongLast") {
		t.Fatalf("Sweep err = %q, want the deepest header associated with the abort", got)
	}
}

func TestSweepAbortPrecedesBite(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "runtime", "partial-expect-panic")
	output := strings.Join([]string{
		"target-partial-expect-panic",
		"--- FAIL: TestPartialOutput (0.00s)",
		"panic: abort after EXPECT",
	}, "\n")

	err := Sweep(root, sweepAbortResultRunner(RunResult{ExitCode: 2, Output: output}))
	if err == nil {
		t.Fatal("Sweep err = nil, want the panic to win over partial EXPECT output")
	}
	got := err.Error()
	if !strings.Contains(got, "inner test abort") || !strings.Contains(got, "TestPartialOutput") || strings.Contains(got, "did not bite") {
		t.Fatalf("Sweep err = %q, want the panic to win over partial EXPECT output", got)
	}
}

func TestSweepAttributesProcessAbort(t *testing.T) {
	scopeFamily := mappedFamily(t)
	scope := boundCheck(t, scopeFamily)
	cases := []struct {
		name    string
		result  RunResult
		setup   func(*testing.T, string)
		wantErr string
	}{
		{
			name: "contract package spawn failure",
			result: RunResult{
				ExitCode:    1,
				Termination: subprocess.TerminationSpawnFailed,
			},
			setup: func(t *testing.T, root string) {
				contractFixture(t, root, "runtime", "contract package spawn failure")
			},
			wantErr: `canary 'contract package spawn failure' process abort in contract package "runtime"`,
		},
		{
			name: "contract package signal termination",
			result: RunResult{
				ExitCode:    -1,
				Termination: subprocess.TerminationSignaled,
			},
			setup: func(t *testing.T, root string) {
				contractFixture(t, root, "runtime", "contract package signal termination")
			},
			wantErr: `canary 'contract package signal termination' process abort in contract package "runtime"`,
		},
		{
			name: "conformance scope",
			result: RunResult{
				ExitCode:    1,
				Termination: subprocess.TerminationSpawnFailed,
			},
			setup: func(t *testing.T, root string) {
				fixture(t, canaryFixture(root, scopeFamily, "conformance scope"), "")
			},
			wantErr: `canary 'conformance scope' process abort in conformance check "` + scope + `"`,
		},
		{
			name: "unscoped inner gate",
			result: RunResult{
				ExitCode:    1,
				Termination: subprocess.TerminationSpawnFailed,
			},
			setup: func(t *testing.T, root string) {
				fixture(t, filepath.Join(root, "tests", "canary", "unscoped inner gate"), "")
			},
			wantErr: `canary 'unscoped inner gate' process abort in the inner gate`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			tc.setup(t, root)

			err := Sweep(root, sweepAbortResultRunner(tc.result))
			if err == nil {
				t.Fatal("Sweep err = nil, want the process abort reported")
			}
			if got := err.Error(); got != tc.wantErr {
				t.Fatalf("Sweep err = %q, want %q", got, tc.wantErr)
			}
		})
	}
}

func TestDefaultRunnerPropagatesSpawnFailure(t *testing.T) {
	result := defaultRunner(RunCall{Kind: RunBite, Binary: "bench-no-such-binary"})
	if result.ExitCode != 1 || result.Termination != subprocess.TerminationSpawnFailed {
		t.Fatalf("default runner result = exit %d, termination %v, want 1/spawn failure", result.ExitCode, result.Termination)
	}
	if result.Output != "" {
		t.Fatalf("default runner output = %q, want no raw spawn error", result.Output)
	}
}

func TestSweepKeepsNumericFailuresCompleted(t *testing.T) {
	cases := []struct {
		name string
		exit int
		want string
	}{
		{
			name: "numeric exit 1",
			exit: 1,
			want: `canary 'numeric exit 1' did not bite (want red + "target-numeric exit 1"; got exit 1)`,
		},
		{
			name: "numeric exit 2",
			exit: 2,
			want: `canary 'numeric exit 2' did not bite (want red + "target-numeric exit 2"; got exit 2)`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			contractFixture(t, root, "runtime", tc.name)

			err := Sweep(root, sweepAbortResultRunner(RunResult{ExitCode: tc.exit}))
			if err == nil {
				t.Fatal("Sweep err = nil, want the completed failure reported")
			}
			if got := err.Error(); got != tc.want {
				t.Fatalf("Sweep err = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestSweepProcessAbortPrecedesBite(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "runtime", "partial-expect-process")

	err := Sweep(root, sweepAbortResultRunner(RunResult{
		ExitCode:    1,
		Termination: subprocess.TerminationSpawnFailed,
		Output:      "target-partial-expect-process\nchild failure after EXPECT\n",
	}))
	if err == nil {
		t.Fatal("Sweep err = nil, want the process abort to win over partial EXPECT output")
	}
	got := err.Error()
	if !strings.Contains(got, "process abort") || strings.Contains(got, "did not bite") {
		t.Fatalf("Sweep err = %q, want the process abort to win over partial EXPECT output", got)
	}
}

func TestSweepOrdersMixedAbortFailures(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"a-panic", "b-process", "c-did-not-bite"} {
		contractFixture(t, root, "runtime", name)
	}

	err := Sweep(root, sweepAbortResultsRunner(map[string]RunResult{
		"a-panic": {
			ExitCode: 2,
			Output: strings.Join([]string{
				"--- FAIL: TestFirstAbort (0.00s)",
				"panic: first failure",
			}, "\n"),
		},
		"b-process": {
			ExitCode:    1,
			Termination: subprocess.TerminationSpawnFailed,
		},
		"c-did-not-bite": {ExitCode: 1, Output: "ordinary failure\n"},
	}))
	if err == nil {
		t.Fatal("Sweep err = nil, want all fixture failures reported")
	}
	want := strings.Join([]string{
		"canary 'a-panic' inner test abort in TestFirstAbort: panic: first failure",
		`canary 'b-process' process abort in contract package "runtime"`,
		`canary 'c-did-not-bite' did not bite (want red + "target-c-did-not-bite"; got exit 1)`,
	}, "\n")
	if got := err.Error(); got != want {
		t.Fatalf("Sweep err = %q, want fixture order %q", got, want)
	}
}

func TestSweepBoundsProcessAbortDiagnostic(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "runtime", "bounded-process")
	result := RunResult{
		ExitCode:    1,
		Termination: subprocess.TerminationSpawnFailed,
		Output: "goroutine 123 [running]:\n" +
			strings.Repeat("unbounded child error ", 32) + "\x1b[31m",
	}
	want := `canary 'bounded-process' process abort in contract package "runtime"`

	for range 3 {
		err := Sweep(root, sweepAbortResultRunner(result))
		if err == nil {
			t.Fatal("Sweep err = nil, want the process abort reported")
		}
		got := err.Error()
		if got != want {
			t.Fatalf("Sweep err = %q, want byte-stable %q", got, want)
		}
		if strings.Contains(got, "goroutine 123") || strings.Contains(got, "unbounded child error") || strings.ContainsRune(got, '\x1b') {
			t.Fatalf("Sweep err = %q, want no raw child output", got)
		}
	}
}

func TestSubjectHash(t *testing.T) {
	switch os.Getenv(canaryAbortHelperEnv) {
	case "root":
		panic("subject hash panic")
	case "subtest":
		t.Run("empty_subject", func(t *testing.T) {
			t.Run("hash", func(t *testing.T) {
				panic("subtest hash panic")
			})
		})
	case "repeated-subtest":
		t.Run("hash", func(t *testing.T) {})
		t.Run("hash", func(t *testing.T) {
			panic("repeated subtest panic")
		})
	case "punctuated-subtest":
		t.Run("empty-subject", func(t *testing.T) {
			panic("punctuated subtest panic")
		})
	case "runtime-fatal":
		var mutex sync.Mutex
		mutex.Unlock()
	}
}

func authenticGoAbortOutput(t *testing.T, mode string) string {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run", "^TestSubjectHash$")
	env := os.Environ()
	if mode == "runtime-fatal" {
		// The fatal probe needs only the runtime marker; suppressing stack printing keeps
		// the captured compatibility sample bounded and independent of paths and goroutine IDs.
		bounded := make([]string, 0, len(env)+1)
		for _, value := range env {
			if strings.HasPrefix(value, "GOTRACEBACK=") {
				continue
			}
			bounded = append(bounded, value)
		}
		env = append(bounded, "GOTRACEBACK=none")
	}
	cmd.Env = append(env, canaryAbortHelperEnv+"="+mode)
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("panic helper exited 0; output:\n%s", output)
	}
	return string(output)
}

func sweepAbortRunner(output string) Runner {
	return sweepAbortResultRunner(RunResult{ExitCode: 2, Output: output})
}

func sweepAbortResultRunner(result RunResult) Runner {
	return sweepAbortRunnerWith(func(RunCall) RunResult { return result })
}

func sweepAbortResultsRunner(results map[string]RunResult) Runner {
	return sweepAbortRunnerWith(func(call RunCall) RunResult {
		return results[filepath.Base(call.FixtureDir)]
	})
}

func sweepAbortRunnerWith(run func(RunCall) RunResult) Runner {
	return func(call RunCall) RunResult {
		if toolResult, done := stubToolchain(call); done {
			return toolResult
		}
		if isBaseline(call) {
			return RunResult{ExitCode: 1, Output: "baseline noise\n"}
		}
		return run(call)
	}
}
