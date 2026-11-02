package rgguard

import (
	"reflect"
	"strings"
	"testing"
)

// The verdicts that the coverage rows name. A verdict row compares the verdict without
// its segment, and the message rows grade the segment.
var (
	permits             = Verdict{}
	refusesPathless     = Verdict{Blocked: true, reason: pathless}
	refusesSubstitution = Verdict{Blocked: true, reason: substitution}
)

// TestClassifyRgCalls drives each Classify row of the coverage map with exact command
// strings. The expectations come from ripgrep 15.1.0 help and probes and from POSIX
// quoting, not from the option table.
func TestClassifyRgCalls(t *testing.T) {
	for _, tc := range []struct {
		row, command string
		want         Verdict
	}{
		{"R08", "rg foo", refusesPathless},
		{"R09", "rg foo .", permits},
		{"R09", "rg foo a b", permits},
		{"R10", "cat f | rg foo", permits},
		{"R10", "cat f |& rg foo", permits},
		{"R11", "rg foo | wc -l", refusesPathless},
		{"R12", "rg foo < f", permits},
		{"R12", "rg foo <<< s", permits},
		{"R12", "rg foo <<EOF\nbody\nEOF", permits},
		{"R13", "rg -g '*.go' foo", refusesPathless},
		{"R13", "rg --glob=*.go foo", refusesPathless},
		{"R13", "rg -A 3 foo", refusesPathless},
		{"R14", "rg -nA3 foo", refusesPathless},
		{"R14", "rg -nA 3 foo .", permits},
		{"R15", "rg -e foo src", permits},
		{"R15", "rg -e foo", refusesPathless},
		{"R16", "rg --files src", permits},
		{"R16", "rg --version foo", permits},
		{"R16", "rg -nh foo", permits},
		{"R17", "rg -- -foo .", permits},
		{"R17", "rg -- -foo", refusesPathless},
		{"R18", "rg --zzz foo", permits},
		{"R18", "rg -Zq foo", permits},
		{"R19", "rg foo 'x", permits},
		{"R20", "ls | xargs rg foo", permits},
		{"R20", "xargs rg foo", permits},
		{"R21", "rg foo <(cmd)", permits},
		{"R22", "RIPGREP_CONFIG_PATH=c rg foo", permits},
		{"R23", "rg", permits},
		{"R23", "rg -i", permits},
		{"R24", `rg "$(x)" .`, refusesSubstitution},
		{"R24", "rg foo $(git ls-files)", refusesSubstitution},
		{"R26", "rg 'a$(b)' .", permits},
		{"R26", "rg \\`x\\` .", permits},
		{"R27", "rg $(x)", refusesSubstitution},
		{"R29", "env X=1 rg foo", refusesPathless},
		{"R29", "timeout 5 rg foo", refusesPathless},
		{"R30", "bash -c 'rg foo'", refusesPathless},
		{"R31", "bench worktree exec X -- rg foo", refusesPathless},
		{"R32", "cat f | bench worktree exec X -- rg foo", permits},
		{"R32", "cat f | bash -c 'rg x .; rg foo'", refusesPathless},
		{"R33", "cd x && rg foo", refusesPathless},
		{"R34", "command -v rg", permits},
		{"R34", "command -v rg foo", permits},
		{"R35", "rg foo -", permits},
		{"R35", "rg -f -", permits},
		{"R36", `bash -c "bash -c 'rg foo'"`, permits},
		{"R36", "sudo rg foo", permits},
		{"R58", `bench worktree exec X -- rg "$(x)" .`, refusesSubstitution},
		{"R58", "bench worktree exec X -- rg `x` .", refusesSubstitution},
		{"R59", "cat f | (rg foo)", permits},
		{"R59", "cat f | (cd x; rg foo)", refusesPathless},
	} {
		got := Classify(tc.command)
		got.Segment = nil
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s: Classify(%q) = %+v, want %+v", tc.row, tc.command, got, tc.want)
		}
	}
}

// TestClassifyRgMessages grades the two refusal messages. Each message names its cause
// and its repair, and ends with the projected words of the refused call.
func TestClassifyRgMessages(t *testing.T) {
	for _, tc := range []struct {
		row, command, suffix string
		names                []string
	}{
		{"R25", `rg "$(x)" .`, "segment=rg $(x) .", []string{"command substitution", "single quotes"}},
		{"R28", "rg foo", "segment=rg foo", []string{"no path", "`.`"}},
	} {
		message := Classify(tc.command).Message()
		if !strings.HasSuffix(message, " "+tc.suffix) {
			t.Errorf("%s: message %q does not end with %q", tc.row, message, tc.suffix)
		}
		for _, name := range tc.names {
			if !strings.Contains(message, name) {
				t.Errorf("%s: message %q does not contain %q", tc.row, message, name)
			}
		}
	}
}
