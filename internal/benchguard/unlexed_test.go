package benchguard

import "testing"

// unlexedSentenceWant is the repair for a command the lexer cannot parse, written here
// independently of the production constant.
const unlexedSentenceWant = "Close every quote and escape in the command, then run it again."

// TestClassifyRefusesAnUnlexableBenchCommand proves the guard fails closed. The fallback
// split cannot show which operators follow a Bench call, so a command the lexer refuses
// is denied when it holds a Bench word, and allowed when it holds none. The refusal
// names the segment that holds the word.
func TestClassifyRefusesAnUnlexableBenchCommand(t *testing.T) {
	resolver := notBenchResolver()
	for _, tc := range []struct{ name, command, want string }{
		{"follow-on after a separator, quote on the next line", "git status;bench gate | tail -20\necho 'x", blockedPrefixWant + " " + unlexedSentenceWant + " segment=bench gate operator="},
		{"bench with an open quote", "bench status 'x", blockedPrefixWant + " " + unlexedSentenceWant + " segment=bench status 'x operator="},
	} {
		verdict := Classify(tc.command, resolver)
		if !verdict.Blocked {
			t.Errorf("%s: Classify(%q) allowed the call", tc.name, tc.command)
			continue
		}
		if got := verdict.Message(); got != tc.want {
			t.Errorf("%s: message = %q, want %q", tc.name, got, tc.want)
		}
	}
	for _, command := range []string{"echo 'x\nls -la", "cat benchmark.txt 'x"} {
		if Classify(command, resolver).Blocked {
			t.Errorf("Classify(%q) refused a command with no Bench word", command)
		}
	}
}
