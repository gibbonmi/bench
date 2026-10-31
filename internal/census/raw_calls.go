package census

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/benchguard"
	"github.com/gibbonmi/bench/internal/gitguard"
	"github.com/gibbonmi/bench/internal/poolkey"
	"github.com/gibbonmi/bench/internal/shellcommand"
)

// pathBreak holds the characters that end a path in the raw command text.
const pathBreak = "/ \t\r\n'\";&|<>()`,"

// Record appends one line for a raw call: a command text that names a path under the
// repository's worktree pool, where no simple command invokes Bench. A text that names
// no such path records nothing. The caller ignores the error, so a failed write never
// changes a verdict.
func Record(command, root, home string, now time.Time) error {
	return recordWith(command, root, home, now, benchguard.DefaultResolver())
}

// recordWith is Record with an injected executable resolver, which lets a test name a
// Bench wrapper without the machine's own PATH.
func recordWith(command, root, home string, now time.Time, resolver benchguard.Resolver) error {
	return recordReadWith(command, root, home, "", now, resolver)
}

// poolAttribution answers the assignment id and the verb head of the record that a raw call
// naming a pool path makes, and whether it makes one.
func poolAttribution(command, root, home string, resolver benchguard.Resolver) (string, string, bool) {
	if benchguard.InvokesBench(command, resolver) {
		return "", "", false
	}
	prefix := poolkey.Pool(home, root) + string(filepath.Separator)
	id, ok := assignment(command, prefix)
	if !ok {
		return "", "", false
	}
	head := verbHead(command, prefix)
	return id, head, head != ""
}

// assignment returns the assignment id of the first pool path in the raw command text.
// The scan reads the raw text, so a path inside a quoted word or a heredoc body counts.
// The prefix ends with the separator, so the sibling directory `<pool>x/` never matches.
func assignment(command, prefix string) (string, bool) {
	for rest := command; ; {
		start := strings.Index(rest, prefix)
		if start < 0 {
			return "", false
		}
		rest = rest[start+len(prefix):]
		if id, ok := poolkey.SplitAssignmentSegment(firstSegment(rest)); ok {
			return id, true
		}
	}
}

// firstSegment returns the text up to the first character that ends a path.
func firstSegment(text string) string {
	if end := strings.IndexAny(text, pathBreak); end >= 0 {
		return text[:end]
	}
	return text
}

// verbHead returns the resolved head of the first simple command whose words name a
// pool path. A text that names the pool only where the parser cannot see it, such as a
// heredoc body, falls back to the first command's resolved head.
func verbHead(command, prefix string) string {
	stream := shellcommand.Parse(command)
	commands := make([][]string, 0, len(stream.Commands))
	for _, span := range stream.Commands {
		words := shellcommand.ProjectCommandWords(stream.Tokens[span.Start:span.End])
		if len(words) == 0 {
			continue
		}
		commands = append(commands, words)
	}
	first := ""
	for index, words := range commands {
		head := headAt(commands, index, prefix)
		if first == "" {
			first = head
		}
		for _, word := range words {
			if strings.Contains(word, prefix) {
				return head
			}
		}
	}
	return first
}

// headAt names the executable the simple command at index runs. A command of only
// assignments runs none of its own, so the head comes from the next command in the same
// text that has a command word; the assignment text itself is never a head, because it
// names a variable and its operand rather than a verb. A text whose commands are all
// assignments has no executable at all, and the head degenerates to a key.
func headAt(commands [][]string, index int, prefix string) string {
	for i := index; i < len(commands); i++ {
		if !onlyAssignments(commands[i]) {
			return resolvedHead(commands[i])
		}
	}
	return degenerateHead(commands[index], prefix)
}

// onlyAssignments reports whether a simple command's words set variables and run
// nothing. ResolveRoutinePrefix cannot answer this alone, because it also runs off the
// end of a command whose last word is a routine prefix such as a bare `env`.
func onlyAssignments(words []string) bool {
	for _, word := range words {
		if !shellcommand.IsAssignment(word) {
			return false
		}
	}
	return len(words) > 0
}

// degenerateHead names an assignment-only text by the key it assigns, with the `=` kept
// and the value dropped, such as `W=`. The key of the assignment whose value carries the
// pool path is the one recorded, because an unrelated co-assignment says nothing about
// the call. The value stays out, because the record holds a head and never an operand.
func degenerateHead(words []string, prefix string) string {
	chosen := ""
	for _, word := range words {
		if !shellcommand.IsAssignment(word) {
			continue
		}
		if chosen == "" {
			chosen = word
		}
		if strings.Contains(word, prefix) {
			chosen = word
			break
		}
	}
	if chosen == "" {
		return ""
	}
	key, _, _ := strings.Cut(chosen, "=")
	return key + "="
}

// resolvedHead names the executable a simple command's words actually run, stepping
// over the routine prefixes (an assignment, `env`, `timeout`, `xargs`) that would
// otherwise be recorded in the verb's place. A `git` head also carries the first
// subcommand word, because a bare `git` head hides which verb the drain should
// propose. The caller filters an assignment-only command out first, so the fallback to
// the first word here answers a prefix that runs off the end, such as a bare `env`.
func resolvedHead(words []string) string {
	prefix := shellcommand.ResolveRoutinePrefix(words)
	index := prefix.Index
	if index >= len(words) {
		index = 0
	}
	head := words[index]
	if filepath.Base(head) == "git" {
		if sub, _, ok := gitguard.FindSubcommand(words, index+1, len(words)); ok {
			return head + " " + sub
		}
	}
	return head
}
