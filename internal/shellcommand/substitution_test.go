package shellcommand

import (
	"reflect"
	"testing"
)

// TestParseMarksLiveSubstitution proves Parse marks each token that opens a live command
// substitution and leaves each literal unmarked. The expectations come from POSIX shell
// quoting: single quotes and a backslash make `$(` and a backtick literal, double quotes
// keep them live, and a `(` operator opens a substitution only directly after a bare `$`.
func TestParseMarksLiveSubstitution(t *testing.T) {
	type mark struct {
		Index int
		Token Token
	}
	for _, tc := range []struct {
		name    string
		command string
		want    []mark
	}{
		{"R01 double-quoted dollar paren", `rg "$(x)" .`, []mark{{1, Token{Kind: Word, Text: "$(x)", Substitution: true}}}},
		{"R02 double-quoted backtick", "rg \"`x`\" .", []mark{{1, Token{Kind: Word, Text: "`x`", Substitution: true}}}},
		{"R02 bare backtick", "rg `x` .", []mark{{1, Token{Kind: Word, Text: "`x`", Substitution: true}}}},
		{"R03 bare dollar paren", "rg $(x) .", []mark{{2, Token{Kind: ControlOperator, Text: "(", Substitution: true}}}},
		{"R04 single-quoted dollar paren", "rg 'lit$(x)' .", nil},
		{"R04 single-quoted backtick", "rg 'a`b' .", nil},
		{"R05 escaped dollar paren", `rg \$\(x\) .`, nil},
		{"R05 escaped backticks", "rg \\`x\\` .", nil},
		{"R05 escaped dollar in double quotes", `rg "\$(x)" .`, nil},
		{"R06 subshell", "(rg x .)", nil},
		{"R06 spaced dollar", "a $ (b)", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got []mark
			for i, token := range Parse(tc.command).Tokens {
				if token.Substitution {
					got = append(got, mark{i, token})
				}
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("marked tokens of Parse(%q) = %+v, want %+v", tc.command, got, tc.want)
			}
		})
	}
}
