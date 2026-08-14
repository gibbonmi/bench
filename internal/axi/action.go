// Package axi owns typed, executable follow-up actions for AXI query responses.
package axi

import (
	"errors"
	"regexp"
	"strings"
	"unicode"

	"github.com/gibbonmi/bench/internal/toon"
)

type actionKind uint8

const (
	actionInspectFull actionKind = iota + 1
	actionRetryDiff
	actionInvocation
	actionHarnessPhase
)

// Action is one typed follow-up. Its fields stay private so a caller
// cannot replace a known argument with prose or a placeholder.
type Action struct {
	kind       actionKind
	commit     string
	base       string
	invocation []string
	arguments  []InvocationArgument
	phase      string
	why        string
}

// InvocationArgument is one declared part of an executable invocation.
// Known arguments carry the exact value already available to the query; future
// inputs render as an explicit slot instead of a guessed value.
type InvocationArgument struct {
	known         string
	knownDeclared bool
	future        string
}

// KnownArgument returns an argument whose value is already known.
func KnownArgument(value string) InvocationArgument {
	return InvocationArgument{known: value, knownDeclared: true}
}

// FutureInput returns one explicitly unknown argument slot.
func FutureInput(name string) InvocationArgument {
	return InvocationArgument{future: name}
}

// ExecutableInvocation returns an exact executable follow-up action.
func ExecutableInvocation(why string, arguments ...InvocationArgument) Action {
	declared := append([]InvocationArgument(nil), arguments...)
	invocation := make([]string, 0, len(declared))
	for _, argument := range declared {
		if argument.knownDeclared {
			if rendered, err := renderKnownArgument(argument.known); err == nil {
				invocation = append(invocation, rendered)
			} else {
				invocation = append(invocation, argument.known)
			}
		} else {
			invocation = append(invocation, "<"+argument.future+">")
		}
	}
	return Action{kind: actionInvocation, arguments: declared, invocation: invocation, why: why}
}

// HarnessPhase returns the canonical follow-up phase rather than a shell command.
func HarnessPhase(phase, why string) Action {
	return Action{kind: actionHarnessPhase, phase: phase, why: why}
}

// InspectFull returns the bounded action for a live diff when commit is empty, or
// for the named resolved commit otherwise.
func InspectFull(commit string) Action {
	return Action{kind: actionInspectFull, commit: commit}
}

// InspectFullBase retains an explicit source base in a bounded-diff follow-up.
func InspectFullBase(base string) Action {
	return Action{kind: actionInspectFull, base: base}
}

// RetryAfterMovement is the one why-line every drift refusal advertises. It is
// owned here so a command's own refusal cannot drift from the action that renders it.
const RetryAfterMovement = "retry after the repository stopped moving"

// RetryDiff returns the exact diff invocation that a drift refusal could not satisfy.
func RetryDiff(invocation []string) Action {
	return Action{kind: actionRetryDiff, invocation: append([]string(nil), invocation...)}
}

// RetryInvocation is the same drift retry for any command whose invocation is not a
// diff: the caller declares its arguments and this owner supplies the why-line.
func RetryInvocation(arguments ...InvocationArgument) Action {
	return ExecutableInvocation(RetryAfterMovement, arguments...)
}

// RenderHelp renders the terminal help block, including the definitive empty state.
func RenderHelp(actions []Action) (string, error) {
	rows := make([][]string, 0, len(actions))
	seen := make(map[string]struct{}, len(actions))
	for _, action := range actions {
		if action.hasUnsupportedDisclosureValue() {
			return toon.Table("help", []string{"cmd", "why"}, nil)
		}
		command, why, err := action.render()
		if err != nil {
			return "", err
		}
		key := command + "\x00" + why
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		rows = append(rows, []string{command, why})
	}
	return toon.Table("help", []string{"cmd", "why"}, rows)
}

func (action Action) hasUnsupportedDisclosureValue() bool {
	if hasUnsupportedControl(action.why) {
		return true
	}
	for _, argument := range action.arguments {
		if argument.knownDeclared && !validKnownArgument(argument.known) {
			return true
		}
	}
	return false
}

func hasUnsupportedControl(value string) bool {
	for _, r := range value {
		if unicode.IsControl(r) && r != '\t' && r != '\n' && r != '\r' {
			return true
		}
	}
	return false
}

var fullSHA = regexp.MustCompile(`^[0-9a-f]{40}$`)

func (action Action) render() (string, string, error) {
	switch action.kind {
	case actionInspectFull:
		if len(action.invocation) != 0 {
			return "", "", errors.New("full diff action arguments are derived by the owner")
		}
		args := []string{"diff", "--full"}
		if action.commit != "" && action.base != "" {
			return "", "", errors.New("full diff action cannot combine commit and base")
		}
		if action.commit != "" {
			if !fullSHA.MatchString(action.commit) {
				return "", "", errors.New("full diff action requires a resolved commit sha")
			}
			args = append(args, "--commit", action.commit)
		}
		if action.base != "" {
			if !fullSHA.MatchString(action.base) {
				return "", "", errors.New("full diff action requires a resolved base sha")
			}
			args = append(args, "--base", action.base)
		}
		return "bench " + strings.Join(args, " "), "inspect the complete patch", nil
	case actionRetryDiff:
		if action.commit != "" || action.base != "" {
			return "", "", errors.New("retry action cannot guess a commit")
		}
		if !validDiffInvocation(action.invocation) {
			return "", "", errors.New("retry action requires an exact executable diff invocation")
		}
		return "bench " + strings.Join(action.invocation, " "), RetryAfterMovement, nil
	case actionInvocation:
		if action.phase != "" || action.commit != "" || action.base != "" {
			return "", "", errors.New("executable action has undeclared fields")
		}
		if action.why == "" {
			return "", "", errors.New("executable action requires a reason")
		}
		args := make([]string, 0, len(action.arguments))
		for i, argument := range action.arguments {
			if argument.knownDeclared == (argument.future != "") {
				return "", "", errors.New("executable action requires known arguments or declared future inputs")
			}
			if argument.knownDeclared {
				rendered, err := renderKnownArgument(argument.known)
				if err != nil || (i == 0 && !shellSafeToken(argument.known)) {
					return "", "", errors.New("executable action has an invalid known argument")
				}
				args = append(args, rendered)
				continue
			}
			if i == 0 || !validFutureInput(argument.future) {
				return "", "", errors.New("executable action has an undeclared future input")
			}
			args = append(args, "<"+argument.future+">")
		}
		if len(args) == 0 || args[0] == "bench" {
			return "", "", errors.New("executable action requires a command without prose")
		}
		if !sameStrings(action.invocation, args) {
			return "", "", errors.New("executable action dropped or changed a declared argument")
		}
		return "bench " + strings.Join(action.invocation, " "), action.why, nil
	case actionHarnessPhase:
		if action.commit != "" || len(action.invocation) != 0 || len(action.arguments) != 0 {
			return "", "", errors.New("harness phase action has undeclared fields")
		}
		if !validHarnessPhase(action.phase) || action.why == "" {
			return "", "", errors.New("harness phase action requires a canonical phase and reason")
		}
		return action.phase, action.why, nil
	default:
		return "", "", errors.New("unknown action kind")
	}
}

func validDiffInvocation(args []string) bool {
	if len(args) == 0 || args[0] != "diff" {
		return false
	}
	rest := args[1:]
	if len(rest) == 0 || (len(rest) == 1 && rest[0] == "--full") {
		return true
	}
	if len(rest) == 2 && rest[0] == "--commit" {
		return validValue(rest[1])
	}
	if len(rest) == 2 && rest[0] == "--base" {
		return validValue(rest[1])
	}
	if len(rest) == 3 {
		return (rest[0] == "--full" && rest[1] == "--commit" && validValue(rest[2])) ||
			(rest[0] == "--commit" && validValue(rest[1]) && rest[2] == "--full") ||
			(rest[0] == "--full" && rest[1] == "--base" && validValue(rest[2])) ||
			(rest[0] == "--base" && validValue(rest[1]) && rest[2] == "--full")
	}
	return false
}

func validValue(value string) bool {
	return value != "" && !strings.ContainsAny(value, "<>\t\r\n ")
}

var futureInput = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

func validFutureInput(name string) bool {
	return futureInput.MatchString(name)
}

func validHarnessPhase(phase string) bool {
	switch phase {
	case "/bench-shape-idea", "/bench-what-next":
		return true
	default:
		return false
	}
}

func renderKnownArgument(value string) (string, error) {
	if !validKnownArgument(value) {
		return "", errors.New("known argument is invalid")
	}
	if shellSafeToken(value) {
		return value, nil
	}
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'", nil
}

func validKnownArgument(value string) bool {
	if value == "" || strings.ContainsAny(value, "<>") {
		return false
	}
	return !hasUnsupportedControl(value)
}

func shellSafeToken(value string) bool {
	for _, r := range value {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && !strings.ContainsRune("_@%+=:,./-", r) {
			return false
		}
	}
	return true
}

func sameStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}
