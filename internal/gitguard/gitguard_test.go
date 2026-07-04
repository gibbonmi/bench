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

// TestDescribeClassesEqualsDenyTable is the one-source check: the advertised deny
// surface is exactly the table classify denies from, in order.
func TestDescribeClassesEqualsDenyTable(t *testing.T) {
	want := "git push, git reset --hard, git clean -f, git branch -f, git branch -D, " +
		"git checkout path, git switch --force, git restore path, history rewrite, " +
		"git stash drop, git commit --amend, git update-ref -d, git tag -d, " +
		"git reflog expire, git worktree remove --force"
	if got := DescribeClasses(); got != want {
		t.Errorf("DescribeClasses() =\n  %q\nwant\n  %q", got, want)
	}
}

func TestBlockMessageNamesLabel(t *testing.T) {
	msg := BlockMessage("git push")
	if !strings.HasPrefix(msg, "BLOCKED: `git push`") {
		t.Errorf("BlockMessage did not lead with BLOCKED + label: %q", msg)
	}
	if !strings.Contains(msg, "hand back") {
		t.Errorf("BlockMessage lost the hand-back instruction: %q", msg)
	}
}

func TestClassifyEmptyAllows(t *testing.T) {
	if got := Classify("", refYes); got != "" {
		t.Errorf("Classify(\"\") = %q, want \"\"", got)
	}
}
