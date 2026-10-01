package gitguard

import (
	"strings"
	"testing"
)

// unlexedLabelWant and unlexedAdviceWant are the refusal for a command the lexer cannot
// parse, written here independently of the production table.
const (
	unlexedLabelWant  = "git in a command the guard cannot parse"
	unlexedAdviceWant = "Close every quote and escape in the command, then run it again."
)

// TestClassifyRefusesAnUnlexableGitCommand proves the guard fails closed. The shell runs
// each line before the one that opens a quote, and the fallback split cannot show where
// a later line ends. So a command the lexer refuses is denied when it holds a git word,
// and allowed when it holds none.
func TestClassifyRefusesAnUnlexableGitCommand(t *testing.T) {
	for _, tc := range []struct{ name, command, want string }{
		{"push after a separator, quote on the next line", "git status;git push --force\necho 'x", unlexedLabelWant},
		{"read-only git with an open quote", "git status 'x", unlexedLabelWant},
		{"quoted git word", "'git' push --force\necho 'x", unlexedLabelWant},
		{"git by path after a trailing backslash", "/usr/bin/git log\\", unlexedLabelWant},
		{"no git word", "echo 'x\nls -la", ""},
		{"a word that only contains git", "cat .gitignore 'x", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := Classify(tc.command, refYes); got != tc.want {
				t.Errorf("Classify(%q) = %q, want %q", tc.command, got, tc.want)
			}
		})
	}
}

// TestBlockMessageCarriesUnlexedAdvice proves the refusal names the repair.
func TestBlockMessageCarriesUnlexedAdvice(t *testing.T) {
	if msg := BlockMessage(unlexedLabelWant); !strings.HasSuffix(msg, " "+unlexedAdviceWant) {
		t.Errorf("BlockMessage(%q) did not end with the advice: %q", unlexedLabelWant, msg)
	}
}
