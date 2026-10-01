package shellcommand

import "strings"

// UnlexedAdvice is the one repair sentence for a command the lexer refused. Each guard
// that refuses such a command appends it.
const UnlexedAdvice = "Close every quote and escape in the command, then run it again."

// fallbackSplit tokenizes a command the lexer refused. It cannot know where a quote
// ends, so it folds no quote and splits on whitespace. A run of operator characters,
// a newline included, stays its own token, so a separator never joins two commands
// into one word.
func fallbackSplit(s string) []string {
	var tokens []string
	runes := []rune(s)
	start := -1
	flush := func(end int) {
		if start >= 0 {
			tokens = append(tokens, string(runes[start:end]))
			start = -1
		}
	}
	for i := 0; i < len(runes); {
		switch c := runes[i]; {
		case isSpace(c):
			flush(i)
			i++
		case isPunct(c):
			flush(i)
			run := i
			for i < len(runes) && isPunct(runes[i]) {
				i++
			}
			tokens = append(tokens, string(runes[run:i]))
		default:
			if start < 0 {
				start = i
			}
			i++
		}
	}
	flush(len(runes))
	return tokens
}

// unquote removes the quote and escape characters the fallback split keeps in a word.
var unquote = strings.NewReplacer(`'`, "", `"`, "", `\`, "")

// UnlexedCommand returns the first simple command that holds a word match accepts, when
// the lexer refused the command. A guard refuses that command, because the fallback
// split cannot show where the shell runs the word. The query reads each word without
// its quote and escape characters, so a quoted command name still matches. A command
// the lexer parsed has no such span.
func (s Stream) UnlexedCommand(match func(string) bool) (SimpleCommand, bool) {
	if !s.Unlexed {
		return SimpleCommand{}, false
	}
	for _, span := range s.Commands {
		for _, token := range s.Tokens[span.Start:span.End] {
			if token.Kind == Word && match(unquote.Replace(token.Text)) {
				return span, true
			}
		}
	}
	return SimpleCommand{}, false
}
