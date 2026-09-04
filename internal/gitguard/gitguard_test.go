package gitguard

import (
	"strings"
	"testing"
)

func TestCommandFromEnvelope(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"well-formed", `{"tool_input":{"command":"git push"}}`, "git push"},
		{"non-JSON allows (empty)", `not json at all`, ""},
		{"missing tool_input", `{"other":1}`, ""},
		{"missing command", `{"tool_input":{"foo":"bar"}}`, ""},
		{"empty command", `{"tool_input":{"command":""}}`, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := CommandFromEnvelope([]byte(c.in)); got != c.want {
				t.Errorf("CommandFromEnvelope(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

// TestBlockMessageNamesLabel pins the refusal shape for every label, including each of
// the seven push classes, so a shared or dropped label cannot hide which rule fired.
func TestBlockMessageNamesLabel(t *testing.T) {
	labels := []string{
		"git reset --hard",
		"git push to the default branch",
		"git push --force",
		"git push --delete",
		"git push --all",
		"git push --mirror",
		"git push --tags",
		"git push with an unresolved destination",
	}
	for _, label := range labels {
		msg := BlockMessage(label)
		if !strings.HasPrefix(msg, "BLOCKED: `"+label+"`") {
			t.Errorf("BlockMessage(%q) did not lead with BLOCKED + label: %q", label, msg)
		}
		if !strings.Contains(msg, "hand back") {
			t.Errorf("BlockMessage(%q) lost the hand-back instruction: %q", label, msg)
		}
	}
}

// TestBlockMessageCarriesUnresolvedAdvice pins the whole advice: the unresolved push is
// the refusal an agent can act on, and it has two fixes. The redirected push moves back
// into the repository, and the destination-less push names the remote and the branch.
func TestBlockMessageCarriesUnresolvedAdvice(t *testing.T) {
	const advice = "Run the push from inside the repository. Name the remote and the branch: git push <remote> <branch>."
	msg := BlockMessage("git push with an unresolved destination")
	if !strings.HasSuffix(msg, advice) {
		t.Errorf("BlockMessage for the unresolved label did not end with the advice: %q", msg)
	}
	if other := BlockMessage("git push --force"); strings.Contains(other, advice) {
		t.Errorf("BlockMessage for the force label carried the unresolved advice: %q", other)
	}
}

func TestClassifyEmptyAllows(t *testing.T) {
	if got := Classify("", refYes); got != "" {
		t.Errorf("Classify(\"\") = %q, want \"\"", got)
	}
}
