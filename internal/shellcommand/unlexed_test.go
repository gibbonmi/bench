package shellcommand

import (
	"reflect"
	"testing"
)

// TestParseMarksALexerFailure proves Parse tells its caller when the tokens come from the
// fallback split, so a guard can refuse rather than trust a guess.
func TestParseMarksALexerFailure(t *testing.T) {
	for _, tc := range []struct {
		command string
		want    bool
	}{
		{"echo 'x", true},
		{"git status\necho \"x", true},
		{"git push origin\\", true},
		{"git status", false},
		{"bash -c 'git push && git reset --hard'", false},
	} {
		if got := Parse(tc.command).Unlexed; got != tc.want {
			t.Errorf("Parse(%q).Unlexed = %t, want %t", tc.command, got, tc.want)
		}
	}
}

// TestFallbackSplitsShellOperators proves the fallback keeps each operator run as its
// own token, so `status;git` never welds two commands into one word.
func TestFallbackSplitsShellOperators(t *testing.T) {
	got := Parse("git status;git push --force\necho 'x")
	want := Stream{
		Tokens: []Token{
			{Word, "git"}, {Word, "status"}, {ControlOperator, ";"},
			{Word, "git"}, {Word, "push"}, {Word, "--force"}, {ControlOperator, ";"},
			{Word, "echo"}, {Word, "'x"},
		},
		Commands: []SimpleCommand{{Start: 0, End: 2}, {Start: 3, End: 6}, {Start: 7, End: 9}},
		Unlexed:  true,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse() = %#v, want %#v", got, want)
	}
}

// TestUnlexedCommandNamesTheHolder proves the query names the first simple command that
// holds a matching word, reads a word without its quote and escape characters, and
// answers nothing for a command the lexer parsed.
func TestUnlexedCommandNamesTheHolder(t *testing.T) {
	isGit := func(word string) bool { return word == "git" }
	for _, tc := range []struct {
		name    string
		command string
		want    SimpleCommand
		ok      bool
	}{
		{"holder after a separator", "ls;git status\necho 'x", SimpleCommand{Start: 2, End: 4}, true},
		{"quoted word", "'git' push\necho 'x", SimpleCommand{Start: 0, End: 2}, true},
		{"escaped word", "\\git push 'x", SimpleCommand{Start: 0, End: 3}, true},
		{"no holder", "ls -la\necho 'x", SimpleCommand{}, false},
		{"word that only contains the name", "cat .gitignore 'x", SimpleCommand{}, false},
		{"parsed command", "git push --force", SimpleCommand{}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := Parse(tc.command).UnlexedCommand(isGit)
			if got != tc.want || ok != tc.ok {
				t.Errorf("UnlexedCommand(%q) = %#v, %t, want %#v, %t", tc.command, got, ok, tc.want, tc.ok)
			}
		})
	}
}
