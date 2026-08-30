// Package benchguard classifies shell follow-ons attached to Bench calls.
package benchguard

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gibbonmi/bench/internal/shellcommand"
)

var wrapperFlag = regexp.MustCompile(`^-[A-Za-z]*c[A-Za-z]*$`)

type Resolver struct {
	Getwd        func() (string, error)
	LookPath     func(string) (string, error)
	EvalSymlinks func(string) (string, error)
}

func DefaultResolver() Resolver {
	return Resolver{Getwd: os.Getwd, LookPath: exec.LookPath, EvalSymlinks: filepath.EvalSymlinks}
}

func CommandFromEnvelope(data []byte) (string, error) {
	var envelope struct {
		ToolInput map[string]json.RawMessage `json:"tool_input"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return "", err
	}
	raw, ok := envelope.ToolInput["command"]
	if !ok {
		return "", errors.New("command field missing")
	}
	var command string
	if err := json.Unmarshal(raw, &command); err != nil {
		return "", err
	}
	if command == "" {
		return "", errors.New("command field empty")
	}
	if strings.IndexByte(command, 0) >= 0 {
		return "", errors.New("command field has control byte")
	}
	return command, nil
}

// Side names where the refused operator sits relative to the Bench segment. The
// reader repairs a different thing on each side, so the refusal sentence follows it.
type Side uint8

const (
	// SideAfter is a control operator after the Bench segment, the default side.
	SideAfter Side = iota
	// SideBefore is a control operator before the Bench segment.
	SideBefore
	// SideRedirection is a redirection inside the Bench segment.
	SideRedirection
)

// Verdict is the guard's decision for one command. Segment holds the Bench
// segment's projected words, Operator names the token that caused a refusal, and
// Side names which side of the segment that token sits on.
type Verdict struct {
	Blocked  bool
	Segment  []string
	Operator string
	Side     Side
}

// Message returns the refusal line: the fixed prefix, the repair sentence for the
// operator's side, then the Bench segment and the operator that caused the refusal.
func (v Verdict) Message() string {
	return blockedPrefix + " " + repairSentence(v.Side) + " segment=" + strings.Join(v.Segment, " ") + " operator=" + v.Operator
}

// repairSentence names the one repair that removes the refused token.
func repairSentence(side Side) string {
	switch side {
	case SideBefore:
		return "Run the Bench command from the current directory; it resolves the worktree itself."
	case SideRedirection:
		return "Run the Bench command without a redirection."
	default:
		return followOnSentence
	}
}

// Classify returns the verdict, the Bench segment's projected words, and the adjacent
// operator token. The rule is span-scoped: it reads the first Bench-headed simple
// command. A `bench worktree exec` head allows a heredoc redirection inside its
// span, and a `;` or `&&` after the span when every later simple command is
// non-Bench. Any other Bench head refuses on an operator or a redirection anywhere
// in the stream. The named token follows one precedence: a redirection inside the
// span, else the control operator before the span, else the one after it.
func Classify(command string, resolver Resolver) Verdict {
	return judge(shellcommand.Parse(command), resolver, true)
}

// judge walks the spans in order. The outer flag marks the top-level stream: a
// wrapper string reads one level deep from there, and the exec exception covers a
// direct head only, so it never applies inside that string.
func judge(stream shellcommand.Stream, resolver Resolver, outer bool) Verdict {
	for index, span := range stream.Commands {
		words := shellcommand.ProjectCommandWords(stream.Tokens[span.Start:span.End])
		prefix := shellcommand.ResolveRoutinePrefix(words)
		if !prefix.Executes || prefix.Index >= len(words) {
			continue
		}
		if isBench(words[prefix.Index], resolver) {
			if outer && isExecHead(words, prefix.Index) {
				return judgeExec(stream, index, words, resolver)
			}
			blocked := refusal(stream, index, words, false)
			blocked.Blocked = hasOuterSyntax(stream)
			return blocked
		}
		if !outer || !isWrapper(words[prefix.Index]) {
			continue
		}
		for i := prefix.Index + 1; i+1 < len(words); i++ {
			child := shellcommand.Parse(words[i+1])
			if !wrapperFlag.MatchString(words[i]) || !invokes(child, resolver, false) {
				continue
			}
			if hasOuterSyntax(stream) {
				return refusal(stream, index, words, false)
			}
			if inner := judge(child, resolver, false); inner.Blocked {
				return inner
			}
		}
	}
	return Verdict{}
}

// judgeExec reads the exec segment's own span. A segment before it shapes the Bench
// call, a non-heredoc redirection inside it changes the response, a control operator
// other than `;` or `&&` after it consumes the response, and a later Bench segment
// makes two responses share one call.
func judgeExec(stream shellcommand.Stream, index int, words []string, resolver Resolver) Verdict {
	refused := refusal(stream, index, words, true)
	if index > 0 {
		return refused
	}
	span := stream.Commands[index]
	for i := span.Start; i < span.End; i++ {
		if stream.Tokens[i].Kind == shellcommand.Redirection && !shellcommand.IsHeredoc(stream.Tokens[i]) {
			return refused
		}
	}
	if span.End < len(stream.Tokens) {
		switch stream.Tokens[span.End].Text {
		case ";", "&&":
		default:
			return refused
		}
	}
	for _, later := range stream.Commands[index+1:] {
		if spanInvokesBench(shellcommand.ProjectCommandWords(stream.Tokens[later.Start:later.End]), resolver, true) {
			return refused
		}
	}
	return Verdict{}
}

// isExecHead reports whether the segment's head is a direct `bench worktree exec`
// call, which is the only head the exception covers.
func isExecHead(words []string, index int) bool {
	return index+2 < len(words) && words[index+1] == "worktree" && words[index+2] == "exec"
}

// operatorFor names the token adjacent to one segment, and the side it sits on, by the
// refusal's precedence. The heredocAllowed flag marks a segment the exec exception
// covers: a heredoc there is legal, so it never names the cause, and the next
// redirection does.
func operatorFor(stream shellcommand.Stream, index int, heredocAllowed bool) (string, Side) {
	span := stream.Commands[index]
	tokens := stream.Tokens[span.Start:span.End]
	for i := range tokens {
		if tokens[i].Kind != shellcommand.Redirection {
			continue
		}
		if heredocAllowed && shellcommand.IsHeredoc(tokens[i]) {
			continue
		}
		return shellcommand.RedirectionText(tokens, i), SideRedirection
	}
	if span.Start > 0 && stream.Tokens[span.Start-1].Kind == shellcommand.ControlOperator {
		return stream.Tokens[span.Start-1].Text, SideBefore
	}
	if span.End < len(stream.Tokens) && stream.Tokens[span.End].Kind == shellcommand.ControlOperator {
		return stream.Tokens[span.End].Text, SideAfter
	}
	return "", SideAfter
}

// refusal builds the verdict for one refused segment from the adjacent token.
func refusal(stream shellcommand.Stream, index int, words []string, heredocAllowed bool) Verdict {
	operator, side := operatorFor(stream, index, heredocAllowed)
	return Verdict{Blocked: true, Segment: words, Operator: operator, Side: side}
}

// PoolCd names the pool path that a `cd` in command targets, or the empty string when
// no simple command changes directory into the pool. The scan reads the command text
// only: it never resolves a path, so an unexpanded variable and a relative target stay
// allowed. A wrapper string is one word here, so a relative `cd` inside an exec child
// stays allowed too.
func PoolCd(command, pools string) string {
	if pools == "" {
		return ""
	}
	prefix := pools + string(filepath.Separator)
	stream := shellcommand.Parse(command)
	for _, span := range stream.Commands {
		words := shellcommand.ProjectCommandWords(stream.Tokens[span.Start:span.End])
		routine := shellcommand.ResolveRoutinePrefix(words)
		if !routine.Executes || routine.Index >= len(words) || words[routine.Index] != "cd" {
			continue
		}
		for _, argument := range words[routine.Index+1:] {
			if strings.HasPrefix(argument, prefix) {
				return argument
			}
		}
	}
	return ""
}

// PoolCdMessage returns the refusal line for a `cd` into the pool path. It names the one
// command form a Bench worktree takes, and the target it read.
func PoolCdMessage(target string) string {
	return `BLOCKED: a Bench worktree runs through bench worktree exec "<label>" -- <command>; never cd into the pool path. target=` + target
}

// InvokesBench reports whether command invokes Bench, in a simple command or in a
// string that a shell wrapper runs. The scan goes one wrapper level deep, which is
// the depth Classify reads, so both tests name the same set of Bench calls.
func InvokesBench(command string, r Resolver) bool {
	return invokes(shellcommand.Parse(command), r, true)
}

func invokes(stream shellcommand.Stream, resolver Resolver, wrapper bool) bool {
	for _, span := range stream.Commands {
		if spanInvokesBench(shellcommand.ProjectCommandWords(stream.Tokens[span.Start:span.End]), resolver, wrapper) {
			return true
		}
	}
	return false
}

// spanInvokesBench reports whether one simple command runs Bench, either as its own
// head or through a string that a shell wrapper runs.
func spanInvokesBench(words []string, resolver Resolver, wrapper bool) bool {
	prefix := shellcommand.ResolveRoutinePrefix(words)
	if !prefix.Executes || prefix.Index >= len(words) {
		return false
	}
	if isBench(words[prefix.Index], resolver) {
		return true
	}
	if !wrapper || !isWrapper(words[prefix.Index]) {
		return false
	}
	for i := prefix.Index + 1; i+1 < len(words); i++ {
		if wrapperFlag.MatchString(words[i]) && invokes(shellcommand.Parse(words[i+1]), resolver, false) {
			return true
		}
	}
	return false
}

func isWrapper(word string) bool {
	switch filepath.Base(word) {
	case "sh", "bash", "zsh":
		return true
	}
	return false
}
func isBench(word string, r Resolver) bool {
	if word == "bench" {
		return true
	}
	if base := filepath.Base(word); base == "bench" || base == "bench.sh" {
		return true
	}
	if r.Getwd == nil || r.LookPath == nil || r.EvalSymlinks == nil {
		return false
	}
	target, err := r.resolve(word)
	if err != nil {
		return false
	}
	base := filepath.Base(target)
	return base == "bench" || base == "bench.sh"
}

func (r Resolver) resolve(word string) (string, error) {
	var candidate string
	var err error
	if filepath.Dir(word) == "." {
		candidate, err = r.LookPath(word)
	} else {
		cwd, err := r.Getwd()
		if err != nil {
			return "", err
		}
		candidate = filepath.Join(cwd, word)
	}
	if err != nil {
		return "", err
	}
	return r.EvalSymlinks(candidate)
}
func hasOuterSyntax(stream shellcommand.Stream) bool {
	for _, token := range stream.Tokens {
		if token.Kind == shellcommand.Redirection || token.Kind == shellcommand.ControlOperator {
			return true
		}
	}
	return false
}

// blockedPrefix opens every refusal, and followOnSentence is the repair for the
// default side. BlockMessage joins them, so a caller that reads the exported text and
// a caller that reads a verdict read one source.
const (
	blockedPrefix    = "BLOCKED: Bench response is bounded, complete, and self-contained."
	followOnSentence = "Run the Bench command without a shell follow-on."
)

func BlockMessage() string {
	return blockedPrefix + " " + followOnSentence
}
