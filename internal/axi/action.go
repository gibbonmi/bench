// Package axi owns typed, executable follow-up actions for AXI query responses.
package axi

import (
	"errors"
	"regexp"
	"strings"

	"github.com/gibbonmi/bench/internal/toon"
)

type actionKind uint8

const (
	actionInspectFull actionKind = iota + 1
	actionRetryDiff
)

// Action is one typed executable follow-up. Its fields stay private so a caller
// cannot replace a known argument with prose or a placeholder.
type Action struct {
	kind       actionKind
	commit     string
	invocation []string
}

// InspectFull returns the bounded action for a live diff when commit is empty, or
// for the named resolved commit otherwise.
func InspectFull(commit string) Action {
	return Action{kind: actionInspectFull, commit: commit}
}

// RetryDiff returns the exact diff invocation that a drift refusal could not satisfy.
func RetryDiff(invocation []string) Action {
	return Action{kind: actionRetryDiff, invocation: append([]string(nil), invocation...)}
}

// RenderHelp renders the terminal help block, including the definitive empty state.
func RenderHelp(actions []Action) (string, error) {
	rows := make([][]string, 0, len(actions))
	for _, action := range actions {
		args, why, err := action.render()
		if err != nil {
			return "", err
		}
		rows = append(rows, []string{"bench " + strings.Join(args, " "), why})
	}
	return toon.Table("help", []string{"cmd", "why"}, rows)
}

var fullSHA = regexp.MustCompile(`^[0-9a-f]{40}$`)

func (action Action) render() ([]string, string, error) {
	switch action.kind {
	case actionInspectFull:
		if len(action.invocation) != 0 {
			return nil, "", errors.New("full diff action arguments are derived by the owner")
		}
		args := []string{"diff", "--full"}
		if action.commit != "" {
			if !fullSHA.MatchString(action.commit) {
				return nil, "", errors.New("full diff action requires a resolved commit sha")
			}
			args = append(args, "--commit", action.commit)
		}
		return args, "inspect the complete patch", nil
	case actionRetryDiff:
		if action.commit != "" {
			return nil, "", errors.New("retry action cannot guess a commit")
		}
		if !validDiffInvocation(action.invocation) {
			return nil, "", errors.New("retry action requires an exact executable diff invocation")
		}
		return action.invocation, "retry after the repository stopped moving", nil
	default:
		return nil, "", errors.New("unknown action kind")
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
	if len(rest) == 3 {
		return (rest[0] == "--full" && rest[1] == "--commit" && validValue(rest[2])) ||
			(rest[0] == "--commit" && validValue(rest[1]) && rest[2] == "--full")
	}
	return false
}

func validValue(value string) bool {
	return value != "" && !strings.ContainsAny(value, "<>\t\r\n ")
}
