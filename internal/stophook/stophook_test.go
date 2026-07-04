package stophook

import (
	"fmt"
	"strings"
	"testing"
)

func TestActive(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"json true is active", `{"stop_hook_active":true}`, true},
		{"json false is not active", `{"stop_hook_active":false}`, false},
		{"top-level false is not active", `false`, false},
		{"absent field is not active", `{}`, false},
		{"string True is active", `{"stop_hook_active":"True"}`, true},
		{"number is not active", `{"stop_hook_active":1}`, false},
		{"other string is not active", `{"stop_hook_active":"yes"}`, false},
		{"invalid json is not active", `not json`, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Active([]byte(c.in)); got != c.want {
				t.Errorf("Active(%q) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

// linesOf builds an n-line string joined by "\n" (no trailing newline), each line
// numbered so the boundary of a tail is checkable.
func linesOf(n int) string {
	xs := make([]string, n)
	for i := range xs {
		xs[i] = fmt.Sprintf("line %d", i+1)
	}
	return strings.Join(xs, "\n")
}

func TestTail(t *testing.T) {
	cases := []struct {
		name string
		in   string
		n    int
		want string
	}{
		{"31 lines yields last 30", linesOf(31), 30, linesOf(31)[len("line 1\n"):]},
		{"exactly 30 lines is unchanged", linesOf(30), 30, linesOf(30)},
		{"5 lines is unchanged", linesOf(5), 30, linesOf(5)},
		{"empty is empty", "", 30, ""},
		{"trailing newline is not a blank line", "a\nb\n", 30, "a\nb"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Tail(c.in, c.n); got != c.want {
				t.Errorf("Tail(%q, %d) = %q, want %q", c.in, c.n, got, c.want)
			}
		})
	}
}

// TestTailBoundaryCount is the sharp boundary check: 31 lines under a 30-line tail
// drops exactly the first line and keeps 30.
func TestTailBoundaryCount(t *testing.T) {
	got := Tail(linesOf(31), 30)
	lines := strings.Split(got, "\n")
	if len(lines) != 30 {
		t.Fatalf("Tail(31 lines, 30) produced %d lines, want 30", len(lines))
	}
	if lines[0] != "line 2" {
		t.Errorf("first kept line = %q, want %q", lines[0], "line 2")
	}
	if lines[29] != "line 31" {
		t.Errorf("last kept line = %q, want %q", lines[29], "line 31")
	}
}

func TestBlockMessage(t *testing.T) {
	msg := BlockMessage(linesOf(40))

	if !strings.HasPrefix(msg, "BLOCKED: the gate is red, so this shift is not done.\n") {
		t.Errorf("BlockMessage missing header, got:\n%s", msg)
	}
	if !strings.Contains(msg, "do not weaken or skip a check") {
		t.Errorf("BlockMessage missing middle header line, got:\n%s", msg)
	}
	if !strings.HasSuffix(msg, "Gate output:\n"+Tail(linesOf(40), 30)) {
		t.Errorf("BlockMessage tail mismatch, got:\n%s", msg)
	}

	// Only the last 30 of the 40 gate lines survive into the message.
	if strings.Contains(msg, "line 10\n") {
		t.Errorf("BlockMessage kept a line beyond the 30-line tail (line 10), got:\n%s", msg)
	}
	if !strings.Contains(msg, "line 11") {
		t.Errorf("BlockMessage dropped the first kept line (line 11), got:\n%s", msg)
	}
	if !strings.Contains(msg, "line 40") {
		t.Errorf("BlockMessage dropped the last line (line 40), got:\n%s", msg)
	}
}
