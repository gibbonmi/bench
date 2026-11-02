// Package rgguard decides whether one Bash command holds an rg call that the guard
// refuses. A call with a pattern and no path reads its stdin, so it can wait on an open
// pipe. A call with a live command substitution runs a command that the reader does not
// see. The package starts no process and reads only the command string.
package rgguard

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/shellcommand"
)

// Verdict is the guard's decision for one command. The zero Verdict permits the command.
// Segment holds the projected words of the refused rg call, and reason names its cause.
type Verdict struct {
	Blocked bool
	Segment []string
	reason  reason
}

// reason names the cause of a refusal. Each cause has its own repair.
type reason uint8

const (
	pathless reason = iota
	substitution
)

// sentences holds the refusal sentence of each cause: the cause, then its repair.
var sentences = [...]string{
	pathless:     "This rg call has a pattern and no path, so it reads stdin and can wait on an open pipe. Name the path to search, such as `.` for the current directory.",
	substitution: "This rg call holds a command substitution, which the shell runs before rg starts. Put the pattern in single quotes, or run the other command as its own step.",
}

// Message returns the complete refusal text for a blocked verdict: the refusal sentence
// of its cause, then the words of the refused call.
func (v Verdict) Message() string {
	return "BLOCKED: " + sentences[v.reason] + " segment=" + strings.Join(v.Segment, " ")
}

// CommandFromEnvelope returns the command of a PreToolUse Bash envelope, or the error of
// an envelope that does not decode.
func CommandFromEnvelope(data []byte) (string, error) {
	envelope, err := shellcommand.DecodeEnvelope(data, shellcommand.CommandStruct)
	if err != nil {
		return "", err
	}
	return envelope.Command(), nil
}

// Classify returns the verdict for one Bash command. It reads each simple command, the -c
// children of a shell one level down, and the child of a Bench worktree exec call. It
// returns the refusal of the first rg call that it refuses.
func Classify(command string) Verdict {
	return judge(shellcommand.Parse(command), false, true)
}

// judge reads the spans of one stream in order. A stream that the lexer refused permits
// each call, because the fallback split cannot show the operands. The fed flag marks a
// child whose outer call has a stdin supplier, and the first rg call of the child reads
// it. The guard reads children of the top-level stream only, which the outer flag marks.
func judge(stream shellcommand.Stream, fed, outer bool) Verdict {
	if stream.Unlexed {
		return Verdict{}
	}
	for _, span := range stream.Commands {
		tokens := stream.Tokens[span.Start:span.End]
		words := shellcommand.ProjectCommandWords(tokens)
		prefix := shellcommand.ResolveRoutinePrefix(words)
		if !prefix.Executes || prefix.Index >= len(words) {
			continue
		}
		at := site{marks: tokens, end: endOperator(stream, span), fed: fed || supplied(stream, span)}
		switch {
		case isRg(words[prefix.Index]):
			fed = false
			if verdict := at.judge(words, prefix); verdict.Blocked {
				return verdict
			}
		case !outer:
		case shellcommand.IsBenchName(words[prefix.Index]):
			child := shellcommand.WorktreeExec(words, prefix.Index).Child
			if inner := shellcommand.ResolveRoutinePrefix(child); inner.Executes && inner.Index < len(child) && isRg(child[inner.Index]) {
				at.marks = afterSeparator(tokens, len(words)-len(child)-1)
				if verdict := at.judge(child, inner); verdict.Blocked {
					return verdict
				}
			}
		default:
			for _, wrapper := range shellcommand.ShellChildren(words, prefix.Index) {
				if verdict := judge(wrapper.Stream, at.fed, false); verdict.Blocked {
					return verdict
				}
			}
		}
	}
	return Verdict{}
}

// site holds the stream facts around one rg call: the tokens whose marks reach the call,
// the operator at the end of its span, and whether a pipe or an input redirection
// supplies its stdin.
type site struct {
	marks []shellcommand.Token
	end   shellcommand.Token
	fed   bool
}

// judge decides one rg call. words holds the call with its prefix, and prefix locates rg
// in it. A live substitution refuses first. An uncertain call, a stdin supplier, and a
// mode that needs no path each permit next. A pattern with no path refuses last.
func (s site) judge(words []string, prefix shellcommand.RoutinePrefix) Verdict {
	if s.end.Substitution || slices.ContainsFunc(s.marks, func(token shellcommand.Token) bool { return token.Substitution }) {
		return Verdict{Blocked: true, Segment: words, reason: substitution}
	}
	call, known := readOperands(words[prefix.Index+1:])
	if !known || prefix.ViaXargs || processSubstitution(s.end) || slices.ContainsFunc(words[:prefix.Index], configAssignment) {
		return Verdict{}
	}
	if s.fed || call.stdin || call.mode || !call.pathless() {
		return Verdict{}
	}
	return Verdict{Blocked: true, Segment: words, reason: pathless}
}

// isRg reports whether a word names the ripgrep executable, bare or by path.
func isRg(word string) bool { return filepath.Base(word) == "rg" }

// endOperator returns the control operator that ends one span, or the zero Token when the
// span ends the stream.
func endOperator(stream shellcommand.Stream, span shellcommand.SimpleCommand) shellcommand.Token {
	if span.End < len(stream.Tokens) {
		return stream.Tokens[span.End]
	}
	return shellcommand.Token{}
}

// supplied reports whether a pipe into one span or an input redirection inside it
// supplies the stdin of the span. A pipe after the span feeds the next command instead.
// A pipe into a subshell group feeds the first span of the group. The lexer gives `| (`
// as two operators and `|(` as one, so the read joins a lone `(` to the token before it
// and then drops one trailing `(` from the lead.
func supplied(stream shellcommand.Stream, span shellcommand.SimpleCommand) bool {
	var lead string
	if span.Start > 0 {
		lead = stream.Tokens[span.Start-1].Text
	}
	if lead == "(" && span.Start > 1 {
		lead = stream.Tokens[span.Start-2].Text + lead
	}
	switch strings.TrimSuffix(lead, "(") {
	case "|", "|&":
		return true
	}
	return slices.ContainsFunc(stream.Tokens[span.Start:span.End], func(token shellcommand.Token) bool {
		return token.Kind == shellcommand.Redirection && strings.HasPrefix(strings.TrimLeft(token.Text, "0123456789"), "<")
	})
}

// processSubstitution reports whether an operator opens a process substitution, whose
// path the command text does not name.
func processSubstitution(operator shellcommand.Token) bool {
	return operator.Text == "<(" || operator.Text == ">("
}

// configAssignment reports whether a word names a ripgrep config file, which can add a
// pattern or a path that the command text does not show.
func configAssignment(word string) bool { return strings.HasPrefix(word, "RIPGREP_CONFIG_PATH=") }

// afterSeparator returns the span tokens after the separator of a WorktreeExec child,
// which is the projected word at index separator. WorktreeExec returns plain argv, so the
// marks of its child come from these outer tokens. A redirection can sit between words,
// so the first token prefix that projects past the separator index ends at the separator.
func afterSeparator(tokens []shellcommand.Token, separator int) []shellcommand.Token {
	for k := range tokens {
		if len(shellcommand.ProjectCommandWords(tokens[:k+1])) > separator {
			return tokens[k+1:]
		}
	}
	return nil
}

// operands is the reading of the words after rg.
type operands struct {
	positionals int
	flagPattern bool // -e or -f supplies the pattern
	stdin       bool // -f - reads the patterns from stdin
	mode        bool // a mode needs no pattern and no path
}

// pathless reports whether the call has a pattern and no path. After -e or -f every
// positional is a path; otherwise the first positional is the pattern. A `-` path also
// asks for stdin, and as a path it already permits the call.
func (o operands) pathless() bool {
	if o.flagPattern {
		return o.positionals == 0
	}
	return o.positionals == 1
}

// apply records what one option word and its value mean for the call.
func (o *operands) apply(r role, value string) {
	switch r {
	case pattern:
		o.flagPattern = true
	case patternFile:
		o.flagPattern = true
		o.stdin = o.stdin || value == "-"
	case mode, valuedMode:
		o.mode = true
	}
}

// readOperands reads the words after rg against the option table. A long value option
// takes its value after `=` or from the next word, and ripgrep has no prefix
// abbreviation. In a short bundle, the first value letter takes the rest of the bundle
// without a leading `=`, or the next word when the bundle ends. After `--`, every word is
// a positional. It reports false for an option word that the table does not list.
func readOperands(args []string) (operands, bool) {
	var call operands
	for i := 0; i < len(args); i++ {
		word := args[i]
		switch {
		case word == "--":
			call.positionals += len(args) - i - 1
			return call, true
		case word == "-" || !strings.HasPrefix(word, "-"):
			call.positionals++
		case strings.HasPrefix(word, "--"):
			name, value, attached := strings.Cut(word, "=")
			r, known := options[name]
			if !known || (attached && !r.takesValue()) {
				return operands{}, false
			}
			if r.takesValue() && !attached && i+1 < len(args) {
				i++
				value = args[i]
			}
			call.apply(r, value)
		default:
			for j := 1; j < len(word); j++ {
				r, known := options["-"+word[j:j+1]]
				if !known {
					return operands{}, false
				}
				if !r.takesValue() {
					call.apply(r, "")
					continue
				}
				value := strings.TrimPrefix(word[j+1:], "=")
				if j+1 == len(word) && i+1 < len(args) {
					i++
					value = args[i]
				}
				call.apply(r, value)
				break
			}
		}
	}
	return call, true
}
