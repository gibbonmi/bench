package benchguard

import (
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/shellcommand"
)

// PoolReference names the pool path that command reaches by a route other than
// `bench worktree exec`, or the empty string when no simple command reaches it. Three
// shapes reach it: a `cd` target, an assignment value in any command position, and a
// git repository or worktree option. The scan reads the command text only: it never
// resolves a path, so an unexpanded variable and a relative target stay allowed. A
// wrapper string is one word here, so a relative `cd` inside an exec child stays
// allowed too. A pool path as a plain argument stays allowed, because a read of a file
// under the pool is not a route into the worktree.
func PoolReference(command, pools string) string {
	if pools == "" {
		return ""
	}
	prefix := pools + string(filepath.Separator)
	stream := shellcommand.Parse(command)
	for _, span := range stream.Commands {
		words := shellcommand.ProjectCommandWords(stream.Tokens[span.Start:span.End])
		if target := assignedPool(words, prefix); target != "" {
			return target
		}
		routine := shellcommand.ResolveRoutinePrefix(words)
		if !routine.Executes || routine.Index >= len(words) {
			continue
		}
		arguments := words[routine.Index+1:]
		switch words[routine.Index] {
		case "cd":
			for _, argument := range arguments {
				if strings.HasPrefix(argument, prefix) {
					return argument
				}
			}
		case "git":
			if target := gitPool(arguments, prefix); target != "" {
				return target
			}
		}
	}
	return ""
}

// assignedPool names the pool path an assignment word holds. The scan reads every word,
// because an assignment before `export` or `env` sits outside the routine prefix, and a
// later command expands the variable the guard never sees.
func assignedPool(words []string, prefix string) string {
	for _, word := range words {
		if !shellcommand.IsAssignment(word) {
			continue
		}
		if value := word[strings.IndexByte(word, '=')+1:]; strings.HasPrefix(value, prefix) {
			return value
		}
	}
	return ""
}

// gitPathOptions are the git options that name a repository or a worktree directory.
// Each value follows its option word, or attaches to it after the option's joiner.
var gitPathOptions = []struct{ option, joiner string }{{"-C", ""}, {"--git-dir", "="}, {"--work-tree", "="}}

// gitPool names the pool path a git option directs the command at.
func gitPool(arguments []string, prefix string) string {
	for index, argument := range arguments {
		for _, flag := range gitPathOptions {
			var value string
			switch {
			case argument == flag.option:
				if index+1 < len(arguments) {
					value = arguments[index+1]
				}
			case strings.HasPrefix(argument, flag.option+flag.joiner):
				value = argument[len(flag.option)+len(flag.joiner):]
			}
			if strings.HasPrefix(value, prefix) {
				return value
			}
		}
	}
	return ""
}

// PoolReferenceMessage returns the refusal line for a command that reaches the pool path
// outside the exec verb. It names the one command form a Bench worktree takes, and the
// target it read.
func PoolReferenceMessage(target string) string {
	return `BLOCKED: a Bench worktree runs through bench worktree exec "<label>" -- <command>; never cd, assign, or git -C into the pool path. target=` + target
}
