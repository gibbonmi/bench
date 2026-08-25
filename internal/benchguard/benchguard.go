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

var assignment = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*=`)
var fileDescriptor = regexp.MustCompile(`^\d+$`)
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

// Classify reports whether command invokes Bench with an outer shell operator or redirection.
func Classify(command string, resolver Resolver) bool {
	return scan(shellcommand.Parse(command), resolver, true)
}

func scan(stream shellcommand.Stream, resolver Resolver, wrapper bool) bool {
	for _, span := range stream.Commands {
		words := commandWords(stream.Tokens[span.Start:span.End])
		if len(words) == 0 {
			continue
		}
		index := resolvePrefix(words)
		if index < len(words) && isBench(words[index], resolver) {
			return hasOuterSyntax(stream)
		}
		if wrapper && index < len(words) && isWrapper(words[index]) {
			for i := index + 1; i+1 < len(words); i++ {
				child := shellcommand.Parse(words[i+1])
				if wrapperFlag.MatchString(words[i]) && containsBench(child, resolver) && (hasOuterSyntax(stream) || scan(child, resolver, false)) {
					return true
				}
			}
		}
	}
	return false
}

func containsBench(stream shellcommand.Stream, resolver Resolver) bool {
	for _, span := range stream.Commands {
		words := commandWords(stream.Tokens[span.Start:span.End])
		index := resolvePrefix(words)
		if index < len(words) && isBench(words[index], resolver) {
			return true
		}
	}
	return false
}

func commandWords(tokens []shellcommand.Token) []string {
	words := make([]string, 0, len(tokens))
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Kind == shellcommand.Redirection {
			if len(words) > 0 && fileDescriptor.MatchString(words[len(words)-1]) {
				words = words[:len(words)-1]
			}
			i++
			continue
		}
		if tokens[i].Kind == shellcommand.Word {
			words = append(words, tokens[i].Text)
		}
	}
	return words
}
func resolvePrefix(words []string) int {
	i := 0
	for i < len(words) && assignment.MatchString(words[i]) {
		i++
	}
	for i < len(words) {
		switch filepath.Base(words[i]) {
		case "env":
			i++
			for i < len(words) && (strings.HasPrefix(words[i], "-") || assignment.MatchString(words[i])) {
				i++
			}
		case "command", "nohup":
			i++
			for i < len(words) && strings.HasPrefix(words[i], "-") {
				i++
			}
		case "timeout":
			i++
			for i < len(words) && strings.HasPrefix(words[i], "-") {
				i++
			}
			if i < len(words) {
				i++
			}
		case "xargs":
			i++
			for i < len(words) && strings.HasPrefix(words[i], "-") {
				i++
			}
		default:
			return i
		}
	}
	return i
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
func BlockMessage() string {
	return "BLOCKED: Bench response is bounded, complete, and self-contained. Run the Bench command without a shell follow-on."
}
