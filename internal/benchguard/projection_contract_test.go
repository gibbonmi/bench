package benchguard

import "testing"

// benchEnvelopeCase is one ordered raw envelope and the CommandFromEnvelope result it keeps.
// The error texts are written here independently, so a respelled production error turns the row red.
type benchEnvelopeCase struct {
	name, in, command, err string
}

func checkBenchEnvelope(t *testing.T, cases []benchEnvelopeCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			command, err := CommandFromEnvelope([]byte(tc.in))
			got := ""
			if err != nil {
				got = err.Error()
			}
			if command != tc.command || got != tc.err {
				t.Errorf("CommandFromEnvelope(%s) = (%q, %q), want (%q, %q)", tc.in, command, got, tc.command, tc.err)
			}
		})
	}
}

// TestBenchEnvelopeRequiredField grades the missing-command error, including exact map keys and map-null clearing.
func TestBenchEnvelopeRequiredField(t *testing.T) {
	checkBenchEnvelope(t, []benchEnvelopeCase{
		{name: "absent command", in: `{"tool_input":{}}`, err: "command field missing"},
		{name: "absent tool_input", in: `{}`, err: "command field missing"},
		{name: "nested key casing", in: `{"tool_input":{"COMMAND":"git reset --hard"}}`, err: "command field missing"},
		{name: "duplicate input null", in: `{"tool_input":{"command":"git reset --hard","file_path":"README.md"},"tool_input":null,"cwd":"/repo"}`, err: "command field missing"},
	})
}

// TestBenchEnvelopeType grades the non-string error and the final raw value under folded root keys.
func TestBenchEnvelopeType(t *testing.T) {
	checkBenchEnvelope(t, []benchEnvelopeCase{
		{name: "numeric command", in: `{"tool_input":{"command":7}}`, err: "json: cannot unmarshal number into Go value of type string"},
		{name: "array command", in: `{"tool_input":{"command":["git status"]}}`, err: "json: cannot unmarshal array into Go value of type string"},
		{name: "earlier command type error", in: `{"tool_input":{"command":7,"command":"git reset --hard"}}`, command: "git reset --hard"},
		{name: "root key casing", in: `{"TOOL_INPUT":{"command":"git reset --hard"}}`, command: "git reset --hard"},
		{name: "unrelated cwd type", in: `{"cwd":7,"tool_input":{"command":"git reset --hard"}}`, command: "git reset --hard"},
	})
}

// TestBenchEnvelopeValue grades the empty and control-byte errors, including a map entry replaced by null.
func TestBenchEnvelopeValue(t *testing.T) {
	checkBenchEnvelope(t, []benchEnvelopeCase{
		{name: "empty command", in: `{"tool_input":{"command":""}}`, err: "command field empty"},
		{name: "null command", in: `{"tool_input":{"command":null}}`, err: "command field empty"},
		{name: "duplicate command null", in: `{"tool_input":{"command":"git reset --hard","command":null}}`, err: "command field empty"},
		{name: "NUL command", in: `{"tool_input":{"command":"git\u0000status"}}`, err: "command field has control byte"},
	})
}

// TestGuardWrapperSelection is GP04. Classify searches every matching `-c` child of one
// wrapper, so a follow-on in the second child refuses, and a nested wrapper stays unread.
// The nested rows read one level only, in the outer judge and in the exec later-segment
// scan, and the exec exception covers a direct head only, so an exec inside a wrapper refuses.
func TestGuardWrapperSelection(t *testing.T) {
	for _, tc := range []struct {
		command, operator string
		blocked           bool
	}{
		{command: "bash -c 'echo ok' -c 'bench status && echo hi'", operator: "&&", blocked: true},
		{command: "bash -c 'bench status | cat' -c 'echo ok'", operator: "|", blocked: true},
		{command: `bash -c "sh -c 'bench status && echo hi'"`},
		{command: `bash -c "sh -c 'bench status'" && echo hi`},
		{command: `bench worktree exec L -- true; bash -c "sh -c 'bench status'"`},
		{command: "bash -c 'bench worktree exec L -- cp a b; cp b a'", operator: ";", blocked: true},
	} {
		got := Classify(tc.command, notBenchResolver())
		if got.Blocked != tc.blocked || got.Operator != tc.operator {
			t.Errorf("Classify(%q) = (%t, %q), want (%t, %q)", tc.command, got.Blocked, got.Operator, tc.blocked, tc.operator)
		}
	}
}

// TestReviewerVerbKeepsRecursiveReach is GP62, GP63, and GP73 with edge family R. The reviewer
// walk recurses through every wrapper and exec child, and it reads an exec child only under a Bench head.
func TestReviewerVerbKeepsRecursiveReach(t *testing.T) {
	for _, tc := range []struct {
		command string
		denied  bool
	}{
		{`bash -c "sh -c 'bench commitment override --assignment x --reason y'"`, true},
		{"bench worktree exec t -- bench worktree exec u -- bench commitment override", true},
		{`bench worktree exec t -- bash -c "bench worktree exec u -- sh -c 'bench commitment override'"`, true},
		{"bench worktree exec t bench commitment override", false},
		{"x worktree exec t -- bench commitment override", false},
	} {
		if got := ReviewerVerb(tc.command, [2]string{"commitment", "override"}, notBenchResolver()); got != tc.denied {
			t.Errorf("ReviewerVerb(%q) = %t, want %t", tc.command, got, tc.denied)
		}
	}
}
