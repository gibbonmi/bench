package benchguard

import (
	"strings"
	"testing"
)

// substitutionSentenceWant is the repair for a Bench call beside a live command
// substitution, written here independently of the production constant.
const substitutionSentenceWant = "Run the Bench command with no unquoted command substitution in it or around it."

// TestClassifyNamesCommandSubstitution proves the follow-on refusal names command
// substitution at an unquoted `$(` and refuses nothing new. A quoted substitution in an
// argument stays allowed, and a bare subshell keeps the plain operator.
func TestClassifyNamesCommandSubstitution(t *testing.T) {
	resolver := notBenchResolver()
	for _, tc := range []struct{ row, command, suffix string }{
		{"R51", "bench gate $(x)", " segment=bench gate $ operator=$("},
		{"R52", "echo $(bench status)", " segment=bench status operator=$("},
	} {
		verdict := Classify(tc.command, resolver)
		if !verdict.Blocked {
			t.Errorf("%s: Classify(%q) allowed the call", tc.row, tc.command)
			continue
		}
		message := verdict.Message()
		if !strings.Contains(message, substitutionSentenceWant) || !strings.HasSuffix(message, tc.suffix) {
			t.Errorf("%s: message = %q, want %q and the suffix %q", tc.row, message, substitutionSentenceWant, tc.suffix)
		}
	}
	if command := "bench learning t --what \"$(cat f)\""; Classify(command, resolver).Blocked {
		t.Errorf("R53: Classify(%q) refused a quoted substitution argument", command)
	}
	const subshell = "true; (bench status)"
	verdict := Classify(subshell, resolver)
	want := blockedPrefixWant + " " + isolationSentenceWant + " segment=bench status operator=("
	if got := verdict.Message(); !verdict.Blocked || got != want {
		t.Errorf("R54: Classify(%q) = blocked %v, message %q, want %q", subshell, verdict.Blocked, got, want)
	}
}
