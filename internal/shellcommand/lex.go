package shellcommand

import "strings"

// lex folds quotes and splits s into shell tokens. Beside each token it records the
// Token.Substitution mark. The mark adds information only and never changes the
// folded text.
func lex(s string) (tokens []string, marks []bool, ok bool) {
	runes := []rune(s)
	var current strings.Builder
	active, marked, opener := false, false, false
	flush := func() {
		if active {
			tokens = append(tokens, current.String())
			marks = append(marks, marked)
			current.Reset()
			active, marked = false, false
		}
	}
	for i := 0; i < len(runes); {
		c := runes[i]
		switch {
		case isSpace(c):
			flush()
			i++
		case isPunct(c):
			flush()
			start := i
			for i < len(runes) && isPunct(runes[i]) {
				i++
			}
			tokens = append(tokens, string(runes[start:i]))
			marks = append(marks, opener)
			opener = false
		case c == '\'':
			active = true
			end := closingQuote(runes, i+1, c)
			if end >= len(runes) {
				return nil, nil, false
			}
			current.WriteString(string(runes[i+1 : end]))
			i = end + 1
		case c == '"':
			active = true
			end := closingQuote(runes, i+1, c)
			if end >= len(runes) {
				return nil, nil, false
			}
			for i++; i < end; {
				if doubleQuoteEscape(runes, i) {
					current.WriteRune(runes[i+1])
					i += 2
					continue
				}
				if runes[i] == '\\' && i+1 < len(runes) && (runes[i+1] == '$' || runes[i+1] == '`') {
					current.WriteString(string(runes[i : i+2]))
					i += 2
					continue
				}
				marked = marked || runes[i] == '`' || dollarParen(runes, i)
				current.WriteRune(runes[i])
				i++
			}
			i = end + 1
		case c == '\\':
			if i+1 >= len(runes) {
				return nil, nil, false
			}
			active = true
			current.WriteRune(runes[i+1])
			i += 2
		default:
			active = true
			marked = marked || c == '`'
			// An unquoted `(` is an operator, so the next token carries the `$(` mark.
			opener = dollarParen(runes, i)
			current.WriteRune(c)
			i++
		}
	}
	flush()
	return tokens, marks, true
}

// closingQuote is the one source for where a quote ends. The quoted text starts at
// runes[i], and quote is the opening rune. It returns the index of the closing quote,
// or len(runes) when no quote closes the text. Inside double quotes, a backslash
// escapes a double quote or a backslash.
func closingQuote(runes []rune, i int, quote rune) int {
	for i < len(runes) && runes[i] != quote {
		if quote == '"' && doubleQuoteEscape(runes, i) {
			i += 2
			continue
		}
		i++
	}
	return i
}

// doubleQuoteEscape reports whether runes[i] is a backslash that escapes the next
// rune inside double quotes.
func doubleQuoteEscape(runes []rune, i int) bool {
	return runes[i] == '\\' && i+1 < len(runes) && (runes[i+1] == '"' || runes[i+1] == '\\')
}

func dollarParen(runes []rune, i int) bool {
	return runes[i] == '$' && i+1 < len(runes) && runes[i+1] == '('
}
