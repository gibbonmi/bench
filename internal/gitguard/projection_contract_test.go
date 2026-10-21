package gitguard

import "testing"

// TestGitProjectionAuthority grades the Git policy around the shared exec projection.
// A child in another checkout and a redirected repository leave push unresolved.
// A non-Bench head projects no child.
func TestGitProjectionAuthority(t *testing.T) {
	cases := []struct {
		name string
		cmd  string
		chk  Checker
		want string
	}{
		{"child push", "bench worktree exec X -- git push origin topic", pushMain, "git push with an unresolved destination"},
		{"redirected push", "git -C /other push origin topic", pushMain, "git push with an unresolved destination"},
		{"non-Bench exec head", "mybench worktree exec X -- git reset --hard", refYes, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Classify(tc.cmd, tc.chk); got != tc.want {
				t.Errorf("Classify(%q) = %q, want %q", tc.cmd, got, tc.want)
			}
		})
	}
}

// TestGitEnvelopeContract grades that an unreadable, missing, null, or empty command reads as "".
func TestGitEnvelopeContract(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"readable command", `{"tool_input":{"command":"git status"}}`, "git status"},
		{"empty input", ``, ""},
		{"invalid JSON", `not json`, ""},
		{"array root", `[]`, ""},
		{"array tool_input", `{"tool_input":[]}`, ""},
		{"numeric command", `{"tool_input":{"command":7}}`, ""},
		{"array command", `{"tool_input":{"command":["git","reset"]}}`, ""},
		{"missing tool_input", `{}`, ""},
		{"missing command", `{"tool_input":{}}`, ""},
		{"null tool_input", `{"tool_input":null}`, ""},
		{"null command", `{"tool_input":{"command":null}}`, ""},
		{"empty command", `{"tool_input":{"command":""}}`, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := CommandFromEnvelope([]byte(tc.in)); got != tc.want {
				t.Errorf("CommandFromEnvelope(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

// TestGuardWrapperSelection grades that Git reads only the first matching shell child.
func TestGuardWrapperSelection(t *testing.T) {
	cases := []struct {
		name string
		cmd  string
		want string
	}{
		{"later reset child", "bash -c 'echo ok' -c 'git reset --hard'", ""},
		{"first reset child", "bash -c 'git reset --hard' -c 'echo ok'", "git reset"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Classify(tc.cmd, refYes); got != tc.want {
				t.Errorf("Classify(%q) = %q, want %q", tc.cmd, got, tc.want)
			}
		})
	}
}
