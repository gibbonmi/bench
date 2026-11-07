package shellcommand

import (
	"reflect"
	"testing"
)

func TestParsePublishesTokenClassesAndSimpleCommands(t *testing.T) {
	got := Parse("env X=1 bench help > log && cat <<'EOF'\nbench gate | tail -20\nEOF\ngit push")
	want := Stream{
		Tokens: []Token{
			{Kind: Word, Text: "env"}, {Kind: Word, Text: "X=1"}, {Kind: Word, Text: "bench"}, {Kind: Word, Text: "help"},
			{Kind: Redirection, Text: ">"}, {Kind: Word, Text: "log"}, {Kind: ControlOperator, Text: "&&"},
			{Kind: Word, Text: "cat"}, {Kind: Redirection, Text: "<<"}, {Kind: Word, Text: "EOF"}, {Kind: ControlOperator, Text: ";"},
			{Kind: Word, Text: "git"}, {Kind: Word, Text: "push"},
		},
		Commands: []SimpleCommand{{Start: 0, End: 6}, {Start: 7, End: 10}, {Start: 11, End: 13}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse() = %#v, want %#v", got, want)
	}
}

func TestParseFoldsQuotedOperators(t *testing.T) {
	got := Parse("bash -c 'git push && git reset --hard'")
	want := Stream{
		Tokens:   []Token{{Kind: Word, Text: "bash"}, {Kind: Word, Text: "-c"}, {Kind: Word, Text: "git push && git reset --hard"}},
		Commands: []SimpleCommand{{Start: 0, End: 3}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse() = %#v, want %#v", got, want)
	}
}

// An unbalanced quote or a trailing backslash stops the lexer. Parse then splits on
// whitespace and operator runs, newlines included, and marks the stream Unlexed. A
// guard still sees each later command as its own simple command and does not read it
// as quoted data.
func TestParseKeepsLaterCommandsVisibleAfterAnUnbalancedInput(t *testing.T) {
	for _, tc := range []struct {
		name    string
		command string
		want    Stream
	}{
		{
			name:    "unterminated single quote",
			command: "echo 'oops && git push --force",
			want: Stream{
				Tokens: []Token{
					{Kind: Word, Text: "echo"}, {Kind: Word, Text: "'oops"}, {Kind: ControlOperator, Text: "&&"},
					{Kind: Word, Text: "git"}, {Kind: Word, Text: "push"}, {Kind: Word, Text: "--force"},
				},
				Commands: []SimpleCommand{{Start: 0, End: 2}, {Start: 3, End: 6}},
				Unlexed:  true,
			},
		},
		{
			name:    "unterminated double quote across a newline",
			command: "echo \"oops\ngit push",
			want: Stream{
				Tokens:   []Token{{Kind: Word, Text: "echo"}, {Kind: Word, Text: "\"oops"}, {Kind: ControlOperator, Text: ";"}, {Kind: Word, Text: "git"}, {Kind: Word, Text: "push"}},
				Commands: []SimpleCommand{{Start: 0, End: 2}, {Start: 3, End: 5}},
				Unlexed:  true,
			},
		},
		{
			name:    "trailing backslash",
			command: "git status | git push origin\\",
			want: Stream{
				Tokens: []Token{
					{Kind: Word, Text: "git"}, {Kind: Word, Text: "status"}, {Kind: ControlOperator, Text: "|"},
					{Kind: Word, Text: "git"}, {Kind: Word, Text: "push"}, {Kind: Word, Text: "origin\\"},
				},
				Commands: []SimpleCommand{{Start: 0, End: 2}, {Start: 3, End: 6}},
				Unlexed:  true,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := Parse(tc.command); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Parse(%q) = %#v, want %#v", tc.command, got, tc.want)
			}
		})
	}
}

// A `<<` inside quotes is data, not a heredoc operator, so Parse removes no line after
// it. A heredoc body starts only after the line on which every quote closes. A guard
// still sees the later command.
func TestParseKeepsLinesAfterAQuotedHeredocOperator(t *testing.T) {
	for _, tc := range []struct {
		name    string
		command string
	}{
		{"double-quoted operator", "echo \"<<X\"\ngit reset --hard\nX"},
		{"single-quoted operator", "echo '<<X'\ngit reset --hard\nX"},
		{"operator in a quote that spans lines", "echo \"a\n<<X\"\ngit reset --hard\nX"},
		{"quote open after a real operator", "cat <<X \"\n\" ; git reset --hard\nX"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stream := Parse(tc.command)
			for _, command := range stream.Commands {
				words := ProjectCommandWords(stream.Tokens[command.Start:command.End])
				if reflect.DeepEqual(words, []string{"git", "reset", "--hard"}) {
					return
				}
			}
			t.Errorf("Parse(%q) = %#v, want a simple command git reset --hard", tc.command, stream)
		})
	}
}

func TestProjectCommandWordsDropsRedirectionsAndDescriptors(t *testing.T) {
	stream := Parse("2>/dev/null env -u X bench help")
	if got, want := ProjectCommandWords(stream.Tokens), []string{"env", "-u", "X", "bench", "help"}; !reflect.DeepEqual(got, want) {
		t.Errorf("ProjectCommandWords() = %#v, want %#v", got, want)
	}
}

func TestResolveRoutinePrefix(t *testing.T) {
	for _, tc := range []struct {
		name     string
		words    []string
		index    int
		viaXargs bool
		executes bool
	}{
		{"env unset", []string{"env", "-u", "X", "bench"}, 3, false, true},
		{"timeout signal and kill after", []string{"timeout", "-s", "KILL", "-k", "1", "5", "bench"}, 6, false, true},
		{"xargs max args", []string{"xargs", "-n", "1", "bench"}, 3, true, true},
		{"command query", []string{"command", "-V", "bench"}, 2, false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := ResolveRoutinePrefix(tc.words)
			if got.Index != tc.index || got.ViaXargs != tc.viaXargs || got.Executes != tc.executes {
				t.Errorf("ResolveRoutinePrefix(%q) = %#v, want index %d, viaXargs %t, executes %t", tc.words, got, tc.index, tc.viaXargs, tc.executes)
			}
		})
	}
}
