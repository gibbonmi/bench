package benchguard

import (
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/shellcommand"
)

// ReviewerVerb reports whether command runs verb, a Bench family and form that only the
// reviewer runs. The verb counts as the head of a simple command, as the child of `bench
// worktree exec`, and inside a shell wrapper string. A tree target from a worktree passes
// the verb's primary-checkout test, so this rule is the only layer that stops that form.
// The scan reads text only, so a fabricated pseudo-terminal stays outside it.
func ReviewerVerb(command string, verb [2]string, resolver Resolver) bool {
	return streamRunsVerb(shellcommand.Parse(command), verb, resolver)
}

// streamRunsVerb reports whether any simple command of stream runs verb.
func streamRunsVerb(stream shellcommand.Stream, verb [2]string, resolver Resolver) bool {
	return slices.ContainsFunc(stream.Commands, func(span shellcommand.SimpleCommand) bool {
		return runsReviewerVerb(shellcommand.ProjectCommandWords(stream.Tokens[span.Start:span.End]), verb, resolver)
	})
}

// runsReviewerVerb reports whether one simple command's words run verb. It recurses into
// each exec child and each shell child, so the walk has no depth limit.
func runsReviewerVerb(words []string, verb [2]string, resolver Resolver) bool {
	prefix := shellcommand.ResolveRoutinePrefix(words)
	if !prefix.Executes || prefix.Index >= len(words) {
		return false
	}
	if isBench(words[prefix.Index], resolver) {
		if exec := shellcommand.WorktreeExec(words, prefix.Index); exec.Head {
			return runsReviewerVerb(exec.Child, verb, resolver)
		}
		return namesVerb(words[prefix.Index+1:], verb)
	}
	return slices.ContainsFunc(shellcommand.ShellChildren(words, prefix.Index), func(child shellcommand.WrapperChild) bool {
		return streamRunsVerb(child.Stream, verb, resolver)
	})
}

// namesVerb reports whether the arguments of a Bench call name verb. The dispatcher reads a
// leading flag and its value, the tree target, between the family and the form.
func namesVerb(args []string, verb [2]string) bool {
	if len(args) == 0 || args[0] != verb[0] {
		return false
	}
	form := args[1:]
	for len(form) >= 2 && strings.HasPrefix(form[0], "-") {
		form = form[2:]
	}
	return len(form) > 0 && form[0] == verb[1]
}

// ReviewerVerbMessage is the refusal line for a tool call that runs verb.
func ReviewerVerbMessage(verb [2]string) string {
	return "BLOCKED: bench " + verb[0] + " " + verb[1] + " is the reviewer's verb; the reviewer runs it from the primary checkout in their own terminal."
}
