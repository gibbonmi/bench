// Package shellcommand parses the shell subset that Bench guards inspect.
package shellcommand

import (
	"path/filepath"
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

// RoutinePrefix identifies the executable behind shell assignments and routine prefixes.
type RoutinePrefix struct {
	Index    int
	ViaXargs bool
	Executes bool
}

const punctChars = "();<>|&\n"
const spaceChars = " \t\r"

var redirectRe = regexp.MustCompile(`^(?:[0-9]+)?(?:>>?|<<?<?)(?:[|&])?$|^&>>?$`)
var assignmentRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*=`)

// IsAssignment reports whether word is a shell assignment: a portable KEY, then "=",
// then any VALUE including the empty one. It is the one source for that shape, so a
// prefix scan and an --env check agree on what a caller may write.
func IsAssignment(word string) bool { return assignmentRe.MatchString(word) }

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

// ProjectCommandWords removes redirections and their operands from a simple command.
func ProjectCommandWords(tokens []Token) []string {
	words := make([]string, 0, len(tokens))
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Kind == Redirection {
			if len(words) > 0 && isDigits(words[len(words)-1]) {
				words = words[:len(words)-1]
			}
			i++
			continue
		}
		if tokens[i].Kind == Word {
			words = append(words, tokens[i].Text)
		}
	}
	return words
}

// RedirectionText rebuilds a redirection's source text from one simple command's
// tokens. The fd digits, the operator, and the operand join without spaces, so the
// three tokens of `2>&1` print as `2>&1`.
func RedirectionText(tokens []Token, index int) string {
	text := tokens[index].Text
	if index > 0 && tokens[index-1].Kind == Word && isDigits(tokens[index-1].Text) {
		text = tokens[index-1].Text + text
	}
	if index+1 < len(tokens) && tokens[index+1].Kind == Word {
		text += tokens[index+1].Text
	}
	return text
}

// IsHeredoc reports whether a token opens a heredoc. Parse removes the body and
// keeps this operator, and a here-string `<<<` stays a plain redirection.
func IsHeredoc(token Token) bool { return token.Kind == Redirection && token.Text == "<<" }

// ResolveRoutinePrefix finds the command word after shell assignments and routine prefixes.
func ResolveRoutinePrefix(words []string) RoutinePrefix {
	prefix := RoutinePrefix{Executes: true}
	for prefix.Index < len(words) && IsAssignment(words[prefix.Index]) {
		prefix.Index++
	}
	for prefix.Index < len(words) {
		if IsAssignment(words[prefix.Index]) {
			prefix.Index++
			continue
		}
		switch filepath.Base(words[prefix.Index]) {
		case "env":
			prefix.Index = skipEnv(words, prefix.Index+1)
		case "command":
			var query bool
			prefix.Index, query = skipCommand(words, prefix.Index+1)
			if query {
				prefix.Executes = false
				return prefix
			}
		case "nohup":
			prefix.Index = skipFlagOnly(words, prefix.Index+1)
		case "timeout":
			prefix.Index = skipTimeout(words, prefix.Index+1)
		case "xargs":
			prefix.ViaXargs = true
			prefix.Index = skipXargs(words, prefix.Index+1)
		default:
			return prefix
		}
	}
	return prefix
}

func skipEnv(words []string, i int) int {
	for i < len(words) {
		word := words[i]
		if word == "--" {
			return i + 1
		}
		if IsAssignment(word) || !strings.HasPrefix(word, "-") {
			if IsAssignment(word) {
				i++
				continue
			}
			return i
		}
		i++
		if word == "-u" || word == "--unset" || word == "-C" || word == "--chdir" {
			if i < len(words) {
				i++
			}
		}
	}
	return i
}

func skipCommand(words []string, i int) (int, bool) {
	query := false
	for i < len(words) {
		word := words[i]
		if word == "--" {
			return i + 1, query
		}
		if !strings.HasPrefix(word, "-") || word == "-" {
			return i, query
		}
		if strings.ContainsAny(strings.TrimLeft(word, "-"), "vV") {
			query = true
		}
		i++
	}
	return i, query
}

func skipFlagOnly(words []string, i int) int {
	for i < len(words) {
		if words[i] == "--" {
			return i + 1
		}
		if !strings.HasPrefix(words[i], "-") || words[i] == "-" {
			return i
		}
		i++
	}
	return i
}

func skipTimeout(words []string, i int) int {
	for i < len(words) {
		word := words[i]
		if word == "--" {
			i++
			break
		}
		if !strings.HasPrefix(word, "-") || word == "-" {
			break
		}
		i++
		if word == "-s" || word == "--signal" || word == "-k" || word == "--kill-after" {
			if i < len(words) {
				i++
			}
		}
	}
	if i < len(words) {
		i++
	}
	return i
}

func skipXargs(words []string, i int) int {
	for i < len(words) {
		word := words[i]
		if word == "--" {
			return i + 1
		}
		if !strings.HasPrefix(word, "-") || word == "-" {
			return i
		}
		i++
		if xargsValueOption(word) && i < len(words) {
			i++
		}
	}
	return i
}

func xargsValueOption(word string) bool {
	switch word {
	case "-E", "-I", "-L", "-P", "-d", "-n", "-s", "--eof", "--replace", "--max-lines", "--max-procs", "--delimiter", "--max-args", "--max-chars":
		return true
	}
	return false
}

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
