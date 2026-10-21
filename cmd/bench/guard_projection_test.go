package main

import (
	"bytes"
	"strings"
	"testing"
)

// TestGuardGitEnvelopePolicy grades guard-git exits for envelope family E0 and the Git decoder rows.
// Each git reset case blocks only when the standard struct decode keeps the command.
func TestGuardGitEnvelopePolicy(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want int
	}{
		{"empty stdin", ``, 0},
		{"invalid JSON", `not json`, 0},
		{"absent tool_input", `{}`, 0},
		{"null tool_input", `{"tool_input":null}`, 0},
		{"absent command", `{"tool_input":{}}`, 0},
		{"null command", `{"tool_input":{"command":null}}`, 0},
		{"numeric command", `{"tool_input":{"command":7}}`, 0},
		{"array command", `{"tool_input":{"command":["git reset --hard"]}}`, 0},
		{"empty command", `{"tool_input":{"command":""}}`, 0},
		{"NUL command", `{"tool_input":{"command":"\u0000"}}`, 0},
		{"nested key casing", `{"tool_input":{"COMMAND":"git reset --hard"}}`, 2},
		{"duplicate command null", `{"tool_input":{"command":"git reset --hard","command":null}}`, 2},
		{"duplicate input null", `{"tool_input":{"command":"git reset --hard","file_path":"README.md"},"tool_input":null,"cwd":"/repo"}`, 2},
		{"earlier command type error", `{"tool_input":{"command":7,"command":"git reset --hard"}}`, 0},
		{"root key casing", `{"TOOL_INPUT":{"command":"git reset --hard"}}`, 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := guardGit(nil, strings.NewReader(tc.in), &stdout, &stderr)
			if code != tc.want {
				t.Errorf("guardGit(%q) exit = %d, want %d", tc.in, code, tc.want)
			}
			if stdout.Len() != 0 {
				t.Errorf("guardGit(%q) stdout = %q, want empty", tc.in, stdout.String())
			}
			if blocked := strings.HasPrefix(stderr.String(), "BLOCKED:"); blocked != (tc.want == 2) || (tc.want == 0 && stderr.Len() != 0) {
				t.Errorf("guardGit(%q) stderr = %q, want BLOCKED only on exit 2", tc.in, stderr.String())
			}
		})
	}
}
