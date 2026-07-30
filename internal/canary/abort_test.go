package canary

import (
	"os"
	"os/exec"
	"strings"
	"testing"
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
	output := strings.Join([]string{
		"--- FAIL: TestConcurrentWrites (0.00s)",
		"fatal error: concurrent map writes",
	}, "\n")

	err := Sweep(root, sweepAbortResultRunner(RunResult{ExitCode: 2, Output: output}))
	if err == nil {
		t.Fatal("Sweep err = nil, want the runtime fatal reported")
	}
	got := err.Error()
	if !strings.Contains(got, "inner test abort") || !strings.Contains(got, "TestConcurrentWrites") || strings.Contains(got, "did not bite") {
		t.Fatalf("Sweep err = %q, want the runtime fatal attributed without a panic marker", got)
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
	}
}

func authenticGoAbortOutput(t *testing.T, mode string) string {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run", "^TestSubjectHash$")
	cmd.Env = append(os.Environ(), canaryAbortHelperEnv+"="+mode)
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
	return func(call RunCall) RunResult {
		if toolResult, done := stubToolchain(call); done {
			return toolResult
		}
		if isBaseline(call) {
			return RunResult{ExitCode: 1, Output: "baseline noise\n"}
		}
		return result
	}
}
