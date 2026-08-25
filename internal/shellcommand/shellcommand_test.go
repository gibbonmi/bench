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
