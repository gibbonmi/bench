package gitguard

import (
	"regexp"
	"strings"
)

// The tokenizer reproduces the shell lexing of
// shlex(posix=True, punctuation_chars="();<>|&\n", whitespace_split=True). Its job is
// not a general shell parser: it splits a command line into the word and control-op
// tokens scan() consumes. It folds quotes and escapes into words, so a wrapper's inner
// `-c` string survives as one token. And it lexes a bare newline as an operator, so a
// multi-line block scans as separate commands (the common way an agent batches git).

// punctChars are the shell operator characters. A run of them forms one operator token,
// separately from word characters — matching shlex's punctuation_chars grouping. The
// newline is here (not in whitespace) so a bare newline emits as its own operator.
//
// Deliberate divergence: `#` is NOT a comment char here (shlex's commenters would drop
// `#`-to-end-of-line). Honoring it is unnecessary for an honest-mistake guard, and it
// only fails safe when skipped — a `#`-commented destructive verb over-blocks rather
// than slipping through. shlex's own comment handling has the opposite hole: it eats
// the newline separator, so `git status # x`⏎`git push` sneaks a push past the guard.
const punctChars = "();<>|&\n"

// spaceChars separate tokens without emitting one. shlex's default whitespace is
// " \t\r\n" with the newline removed, so CR still splits but a newline does not.
const spaceChars = " \t\r"

// redirectRe matches a redirection operator token (optionally fd-prefixed / fd-dup
// suffixed). The token and its following target are dropped by stripRedirections, and a
// bare leading fd digit is popped.
var redirectRe = regexp.MustCompile(`^(?:[0-9]+)?(?:>>?|<<?<?)(?:[|&])?$|^&>>?$`)

func isSpace(r rune) bool { return strings.ContainsRune(spaceChars, r) }
func isPunct(r rune) bool { return strings.ContainsRune(punctChars, r) }

// allPunct reports whether every rune of tok is an operator character (a non-empty
// operator-only token) — the guard for collapsing a newline-bearing operator token.
func allPunct(tok string) bool {
	if tok == "" {
		return false
	}
	for _, r := range tok {
		if !isPunct(r) {
			return false
		}
	}
	return true
}

// tokenize splits s into scan()'s tokens. It lexes with quote/escape handling. On
// unbalanced quoting (or a trailing escape) it falls back to a plain split. That split
// still honors newlines as boundaries, so a multi-line block can't slip through unclassified.
// It then collapses any operator-only token carrying a newline to a plain control op
// (`\n`→`;`, `;\n`→`;`, `&&\n`→`&&`) and strips redirections.
func tokenize(s string) []string {
	raw, ok := lex(s)
	if !ok {
		raw = fallbackSplit(s)
	}
	out := make([]string, 0, len(raw))
	for _, tok := range raw {
		if strings.Contains(tok, "\n") && allPunct(tok) {
			tok = strings.ReplaceAll(tok, "\n", "")
			if tok == "" {
				tok = ";"
			}
		}
		out = append(out, tok)
	}
	return stripRedirections(out)
}

// lex is the shlex-equivalent scanner. It returns the raw token list, or ok=false when
// a quote is left open or a trailing backslash has nothing to escape (shlex's
// ValueError). That case routes the caller to fallbackSplit. Words accumulate
// non-space, non-operator runs with single/double quotes and backslash escapes folded
// in; a run of operator characters emits as one operator token.
func lex(s string) (tokens []string, ok bool) {
	runes := []rune(s)
	var cur strings.Builder
	active := false // a word token is in progress (may be the empty string, e.g. '')
	flush := func() {
		if active {
			tokens = append(tokens, cur.String())
			cur.Reset()
			active = false
		}
	}
	i, n := 0, len(runes)
	for i < n {
		c := runes[i]
		switch {
		case isSpace(c):
			flush()
			i++
		case isPunct(c):
			flush()
			start := i
			for i < n && isPunct(runes[i]) {
				i++
			}
			tokens = append(tokens, string(runes[start:i]))
		case c == '\'':
			active = true
			i++
			for i < n && runes[i] != '\'' {
				cur.WriteRune(runes[i])
				i++
			}
			if i >= n {
				return nil, false // unterminated single quote
			}
			i++ // closing quote
		case c == '"':
			active = true
			i++
			for i < n && runes[i] != '"' {
				// Inside double quotes shlex's escape char escapes only the quote and
				// the backslash itself; before any other char the backslash is literal.
				if runes[i] == '\\' && i+1 < n && (runes[i+1] == '"' || runes[i+1] == '\\') {
					cur.WriteRune(runes[i+1])
					i += 2
					continue
				}
				cur.WriteRune(runes[i])
				i++
			}
			if i >= n {
				return nil, false // unterminated double quote
			}
			i++ // closing quote
		case c == '\\':
			if i+1 >= n {
				return nil, false // trailing escape with nothing to escape
			}
			active = true
			cur.WriteRune(runes[i+1])
			i += 2
		default:
			active = true
			cur.WriteRune(c)
			i++
		}
	}
	flush()
	return tokens, true
}

// fallbackSplit is the malformed-quoting path. It splits on newlines (each a `;`
// boundary, so a multi-line block still can't hide a destructive verb), then
// whitespace-splits each line. This is the shlex-ValueError fallback: it drops quote
// handling but never a newline boundary.
func fallbackSplit(s string) []string {
	var raw []string
	for idx, line := range strings.Split(s, "\n") {
		if idx > 0 {
			raw = append(raw, ";")
		}
		raw = append(raw, strings.Fields(line)...)
	}
	return raw
}

// stripRedirections drops redirection operator tokens and their targets, popping a bare
// leading fd digit (`2` in `2> file`).
func stripRedirections(tokens []string) []string {
	out := []string{}
	i := 0
	for i < len(tokens) {
		if redirectRe.MatchString(tokens[i]) {
			if len(out) > 0 && isDigits(out[len(out)-1]) {
				out = out[:len(out)-1]
			}
			i += 2
			continue
		}
		out = append(out, tokens[i])
		i++
	}
	return out
}

// isDigits reports whether s is a non-empty run of ASCII digits (the fd-prefix pop; the
// empty string is not digits).
func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
