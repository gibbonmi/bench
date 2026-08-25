// Package shellcommand parses the shell subset that Bench guards inspect.
package shellcommand

import (
	"regexp"
	"strings"
)

// TokenKind identifies the shell role of a token.
type TokenKind uint8

const (
	// Word is an argument or command word after quote folding.
	Word TokenKind = iota
	// ControlOperator separates simple commands.
	ControlOperator
	// Redirection changes a simple command's input or output.
	Redirection
)

// Token is one parsed shell token.
type Token struct {
	Kind TokenKind
	Text string
}

// SimpleCommand identifies the half-open token range for one simple command.
type SimpleCommand struct {
	Start int
	End   int
}

// Stream is the shell token stream and its simple-command spans.
type Stream struct {
	Tokens   []Token
	Commands []SimpleCommand
}

const punctChars = "();<>|&\n"
const spaceChars = " \t\r"

var redirectRe = regexp.MustCompile(`^(?:[0-9]+)?(?:>>?|<<?<?)(?:[|&])?$|^&>>?$`)

// Parse folds quotes, removes heredoc bodies, and returns the remaining shell tokens.
// A heredoc's operator stays in the stream because it is an outer redirection.
func Parse(command string) Stream {
	command = stripHeredocBodies(command)
	raw, ok := lex(command)
	if !ok {
		raw = fallbackSplit(command)
	}
	stream := Stream{Tokens: make([]Token, 0, len(raw))}
	start := 0
	for _, rawToken := range raw {
		token := Token{Kind: Word, Text: rawToken}
		if allPunct(rawToken) {
			token.Text = strings.ReplaceAll(rawToken, "\n", "")
			if token.Text == "" {
				token.Text = ";"
			}
			if redirectRe.MatchString(token.Text) {
				token.Kind = Redirection
			} else {
				token.Kind = ControlOperator
			}
		}
		if token.Kind == ControlOperator {
			if start < len(stream.Tokens) {
				stream.Commands = append(stream.Commands, SimpleCommand{Start: start, End: len(stream.Tokens)})
			}
			stream.Tokens = append(stream.Tokens, token)
			start = len(stream.Tokens)
			continue
		}
		stream.Tokens = append(stream.Tokens, token)
	}
	if start < len(stream.Tokens) {
		stream.Commands = append(stream.Commands, SimpleCommand{Start: start, End: len(stream.Tokens)})
	}
	return stream
}

func isSpace(r rune) bool { return strings.ContainsRune(spaceChars, r) }
func isPunct(r rune) bool { return strings.ContainsRune(punctChars, r) }

func allPunct(token string) bool {
	if token == "" {
		return false
	}
	for _, r := range token {
		if !isPunct(r) {
			return false
		}
	}
	return true
}

func lex(s string) (tokens []string, ok bool) {
	runes := []rune(s)
	var current strings.Builder
	active := false
	flush := func() {
		if active {
			tokens = append(tokens, current.String())
			current.Reset()
			active = false
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
		case c == '\'':
			active = true
			i++
			for i < len(runes) && runes[i] != '\'' {
				current.WriteRune(runes[i])
				i++
			}
			if i >= len(runes) {
				return nil, false
			}
			i++
		case c == '"':
			active = true
			i++
			for i < len(runes) && runes[i] != '"' {
				if runes[i] == '\\' && i+1 < len(runes) && (runes[i+1] == '"' || runes[i+1] == '\\') {
					current.WriteRune(runes[i+1])
					i += 2
					continue
				}
				current.WriteRune(runes[i])
				i++
			}
			if i >= len(runes) {
				return nil, false
			}
			i++
		case c == '\\':
			if i+1 >= len(runes) {
				return nil, false
			}
			active = true
			current.WriteRune(runes[i+1])
			i += 2
		default:
			active = true
			current.WriteRune(c)
			i++
		}
	}
	flush()
	return tokens, true
}

func fallbackSplit(s string) []string {
	var tokens []string
	for index, line := range strings.Split(s, "\n") {
		if index > 0 {
			tokens = append(tokens, ";")
		}
		tokens = append(tokens, strings.Fields(line)...)
	}
	return tokens
}

var heredocOpRe = regexp.MustCompile(`<<<|<<(-?)[ \t]*(?:'([^']+)'|"([^"]+)"|\\?([A-Za-z0-9_][A-Za-z0-9_-]*))`)

func stripHeredocBodies(s string) string {
	type opened struct {
		delim string
		dash  bool
	}
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	var pending []opened
	for _, line := range lines {
		if len(pending) > 0 {
			candidate := line
			if pending[0].dash {
				candidate = strings.TrimLeft(candidate, "\t")
			}
			if candidate == pending[0].delim {
				pending = pending[1:]
			}
			continue
		}
		out = append(out, line)
		for _, match := range heredocOpRe.FindAllStringSubmatch(line, -1) {
			if match[0] == "<<<" {
				continue
			}
			pending = append(pending, opened{delim: match[2] + match[3] + match[4], dash: match[1] == "-"})
		}
	}
	return strings.Join(out, "\n")
}
