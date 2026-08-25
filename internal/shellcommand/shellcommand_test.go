package shellcommand

import (
	"reflect"
	"testing"
)

func TestParsePublishesTokenClassesAndSimpleCommands(t *testing.T) {
	got := Parse("env X=1 bench help > log && cat <<'EOF'\nbench gate | tail -20\nEOF\ngit push")
	want := Stream{
		Tokens: []Token{
			{Word, "env"}, {Word, "X=1"}, {Word, "bench"}, {Word, "help"},
			{Redirection, ">"}, {Word, "log"}, {ControlOperator, "&&"},
			{Word, "cat"}, {Redirection, "<<"}, {Word, "EOF"}, {ControlOperator, ";"},
			{Word, "git"}, {Word, "push"},
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
		Tokens:   []Token{{Word, "bash"}, {Word, "-c"}, {Word, "git push && git reset --hard"}},
		Commands: []SimpleCommand{{Start: 0, End: 3}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse() = %#v, want %#v", got, want)
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
