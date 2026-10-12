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
	stream := shellcommand.Parse(command)
	return slices.ContainsFunc(stream.Commands, func(span shellcommand.SimpleCommand) bool {
		return runsReviewerVerb(shellcommand.ProjectCommandWords(stream.Tokens[span.Start:span.End]), verb, resolver)
	})
}

// runsReviewerVerb reports whether one simple command's words run verb.
func runsReviewerVerb(words []string, verb [2]string, resolver Resolver) bool {
	prefix := shellcommand.ResolveRoutinePrefix(words)
	if !prefix.Executes || prefix.Index >= len(words) {
		return false
	}
	head, args := words[prefix.Index], words[prefix.Index+1:]
	switch {
	case isBench(head, resolver) && isExecHead(words, prefix.Index):
		if child := slices.Index(args, "--"); child >= 0 {
			return runsReviewerVerb(args[child+1:], verb, resolver)
		}
	case isBench(head, resolver):
		return namesVerb(args, verb)
	case isWrapper(head):
		for i := 0; i+1 < len(args); i++ {
			if wrapperFlag.MatchString(args[i]) && ReviewerVerb(args[i+1], verb, resolver) {
				return true
			}
		}
	}
	return false
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
