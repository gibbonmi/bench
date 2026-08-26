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

// Classify reports whether command invokes Bench with an outer shell operator or redirection.
func Classify(command string, resolver Resolver) bool {
	return scan(shellcommand.Parse(command), resolver, true)
}

func scan(stream shellcommand.Stream, resolver Resolver, wrapper bool) bool {
	for _, span := range stream.Commands {
		words := shellcommand.ProjectCommandWords(stream.Tokens[span.Start:span.End])
		if len(words) == 0 {
			continue
		}
		prefix := shellcommand.ResolveRoutinePrefix(words)
		if prefix.Executes && prefix.Index < len(words) && isBench(words[prefix.Index], resolver) {
			return hasOuterSyntax(stream)
		}
		if wrapper && prefix.Executes && prefix.Index < len(words) && isWrapper(words[prefix.Index]) {
			for i := prefix.Index + 1; i+1 < len(words); i++ {
				child := shellcommand.Parse(words[i+1])
				if wrapperFlag.MatchString(words[i]) && invokes(child, resolver, false) && (hasOuterSyntax(stream) || scan(child, resolver, false)) {
					return true
				}
			}
		}
	}
	return false
}

// InvokesBench reports whether command invokes Bench, in a simple command or in a
// string that a shell wrapper runs. The scan goes one wrapper level deep, which is
// the depth Classify reads, so both tests name the same set of Bench calls.
func InvokesBench(command string, r Resolver) bool {
	return invokes(shellcommand.Parse(command), r, true)
}

func invokes(stream shellcommand.Stream, resolver Resolver, wrapper bool) bool {
	for _, span := range stream.Commands {
		words := shellcommand.ProjectCommandWords(stream.Tokens[span.Start:span.End])
		prefix := shellcommand.ResolveRoutinePrefix(words)
		if !prefix.Executes || prefix.Index >= len(words) {
			continue
		}
		if isBench(words[prefix.Index], resolver) {
			return true
		}
		if !wrapper || !isWrapper(words[prefix.Index]) {
			continue
		}
		for i := prefix.Index + 1; i+1 < len(words); i++ {
			if wrapperFlag.MatchString(words[i]) && invokes(shellcommand.Parse(words[i+1]), resolver, false) {
				return true
			}
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
func BlockMessage() string {
	return "BLOCKED: Bench response is bounded, complete, and self-contained. Run the Bench command without a shell follow-on."
}
