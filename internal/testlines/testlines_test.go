package testlines

import (
	"strings"
	"testing"
)

func TestRunnerLine(t *testing.T) {
	cases := []struct {
		line string
		want bool
	}{
		{"=== RUN   TestX", true},
		{"=== PAUSE TestX", true},
		{"=== CONT  TestX", true},
		{"=== NAME  TestX", true},
		{"--- FAIL: TestX (0.00s)", true},
		{"--- PASS: TestX (0.00s)", true},
		{"--- SKIP: TestX (0.00s)", true},
		{"# github.com/x/y", true},
		{"bench-skip reason", true},
		{"FAIL", true},
		{"FAIL\tgithub.com/x/y 0.1s", true},
		{"FAIL github.com/x/y", true},
		{"ok\tgithub.com/x/y", true},
		{"ok github.com/x/y", true},
		{"?\tgithub.com/x/y\t[no test files]", true},
		{"? github.com/x/y", true},
		{"    x.go:12: want 1, got 2", false},
		{"panic: boom", false},
		{"WARNING: DATA RACE", false},
	}
	for _, c := range cases {
		if got := RunnerLine(c.line); got != c.want {
			t.Errorf("RunnerLine(%q) = %v, want %v", c.line, got, c.want)
		}
	}
}

func TestFailureRows(t *testing.T) {
	cases := []struct {
		name  string
		lines []string
		want  []string
	}{
		{
			name:  "two failure names yield two rows",
			lines: []string{"=== RUN   TestA", "--- FAIL: TestA (0.00s)", "=== RUN   TestB", "--- FAIL: TestB (0.00s)", "FAIL"},
			want:  []string{"--- FAIL: TestA (0.00s)", "--- FAIL: TestB (0.00s)"},
		},
		{
			name:  "diagnostics travel with the failure name row",
			lines: []string{"--- FAIL: TestA (0.00s)", "    a_test.go:12: want 1, got 2", "    a_test.go:13: second detail"},
			want:  []string{"--- FAIL: TestA (0.00s)", "a_test.go:12: want 1, got 2", "a_test.go:13: second detail"},
		},
		{
			name:  "build-failed terminal and panic are rows",
			lines: []string{"FAIL\tgithub.com/x/y [build failed]", "panic: boom"},
			want:  []string{"FAIL\tgithub.com/x/y [build failed]", "panic: boom"},
		},
		{
			name:  "build-error block carries its compiler line",
			lines: []string{"# github.com/x/y", "./x.go:12:3: undefined: y"},
			want:  []string{"# github.com/x/y", "./x.go:12:3: undefined: y"},
		},
		{
			name:  "an indented failure name opens a block",
			lines: []string{"    --- FAIL: TestA/sub (0.00s)", "        a_test.go:12: detail"},
			want:  []string{"--- FAIL: TestA/sub (0.00s)", "a_test.go:12: detail"},
		},
		{
			name:  "a pass line ends the block and adds no row",
			lines: []string{"--- FAIL: TestA (0.00s)", "    a_test.go:12: detail", "--- PASS: TestB (0.00s)", "    stray line"},
			want:  []string{"--- FAIL: TestA (0.00s)", "a_test.go:12: detail"},
		},
		{
			name:  "a run line ends the block and adds no row",
			lines: []string{"--- FAIL: TestA (0.00s)", "    a_test.go:12: detail", "=== RUN   TestB", "    stray line"},
			want:  []string{"--- FAIL: TestA (0.00s)", "a_test.go:12: detail"},
		},
		{
			name:  "ok and no-test lines are never rows",
			lines: []string{"ok\tgithub.com/x/y\t0.10s", "?\tgithub.com/x/z\t[no test files]"},
			want:  []string{},
		},
		{
			name:  "an unmatched stream yields an empty slice",
			lines: []string{"WARNING: DATA RACE", "WARNING: DATA RACE"},
			want:  []string{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := FailureRows(c.lines)
			if got == nil {
				t.Fatal("FailureRows returned nil, want a slice")
			}
			if len(got) != len(c.want) {
				t.Fatalf("FailureRows = %q, want %q", got, c.want)
			}
			for i := range got {
				if got[i] != c.want[i] {
					t.Errorf("row %d = %q, want %q", i, got[i], c.want[i])
				}
			}
		})
	}
}

func TestFailureRowsReadsSplitStream(t *testing.T) {
	stream := "=== RUN   TestA\n--- FAIL: TestA (0.00s)\n    a_test.go:9: boom\nFAIL\nFAIL\tgithub.com/x/y\t0.10s\n"
	got := FailureRows(strings.Split(stream, "\n"))
	want := []string{"--- FAIL: TestA (0.00s)", "a_test.go:9: boom", "FAIL\tgithub.com/x/y\t0.10s"}
	if len(got) != len(want) {
		t.Fatalf("FailureRows = %q, want %q", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("row %d = %q, want %q", i, got[i], want[i])
		}
	}
}
