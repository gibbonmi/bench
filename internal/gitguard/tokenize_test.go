package gitguard

import (
	"reflect"
	"testing"
)

// TestTokenizeEdges pins the subtle lexing the 79-case matrix samples but cannot
// exhaust. It covers newline-as-separator, operator-token collapsing, malformed-quote
// fallback, redirection stripping, and a quoted newline surviving into a wrapper word.
func TestTokenizeEdges(t *testing.T) {
	nl := "\n"
	cases := []struct {
		name string
		in   string
		want []string
	}{
		{"plain split", "git push", []string{"git", "push"}},
		{"single-quote is one word", "bash -c 'git push'", []string{"bash", "-c", "git push"}},
		{"operators as separators", "git status && git push", []string{"git", "status", "&&", "git", "push"}},
		{"semicolon separator", "git add -A; git push", []string{"git", "add", "-A", ";", "git", "push"}},
		{"newline collapses to control op", "git add -A" + nl + "git push", []string{"git", "add", "-A", ";", "git", "push"}},
		{"trailing operator run with newline", "git status &&" + nl + "git push", []string{"git", "status", "&&", "git", "push"}},
		{"redirection with target stripped", "git checkout README.md > log", []string{"git", "checkout", "README.md"}},
		{"fd-prefixed redirection stripped", "git checkout main 2>/dev/null", []string{"git", "checkout", "main"}},
		{"dup redirection stripped", "git checkout main 2>&1", []string{"git", "checkout", "main"}},
		{"quoted newline survives in word", "bash -c 'cd r" + nl + "git push'", []string{"bash", "-c", "cd r" + nl + "git push"}},
		{"malformed quote falls back, newline still a boundary", "git 'unbalanced" + nl + "git push", []string{"git", "'unbalanced", ";", "git", "push"}},
		{"env assignment stays one token", "GIT_TRACE=1 git push", []string{"GIT_TRACE=1", "git", "push"}},
		// `#` is an ordinary word char here, NOT a shlex comment — a deliberate,
		// fail-safe divergence from the Python original (see tokenize.go). The tokens
		// after `#` survive, so a `#`-commented destructive verb over-blocks rather than
		// slipping through. This pins that the divergence stays intentional.
		{"hash is a word, not a comment", "echo hi # x", []string{"echo", "hi", "#", "x"}},
		{"tilde and slashes stay in word", "git checkout HEAD~1 /tmp/x", []string{"git", "checkout", "HEAD~1", "/tmp/x"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := tokenize(c.in)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("tokenize(%q) = %#v, want %#v", c.in, got, c.want)
			}
		})
	}
}

func TestStripRedirectionsPopsFdDigit(t *testing.T) {
	// `2` `>` `file` → the fd digit is popped and the redirect+target dropped.
	got := stripRedirections([]string{"echo", "hi", "2", ">", "file"})
	want := []string{"echo", "hi"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("stripRedirections = %#v, want %#v", got, want)
	}
}
