// Package census owns the record of raw shell calls that touch a Bench worktree.
package census

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/benchguard"
	"github.com/gibbonmi/bench/internal/gitguard"
	"github.com/gibbonmi/bench/internal/poolkey"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/shellcommand"
)

// Dir returns the directory that holds one repository's census records. The
// directory is a sibling of the worktree pool and never a child of it, because
// `bench worktree reclaim` enumerates `<home>/worktrees` and a foreign entry
// there changes its plan.
func Dir(home, root string) string {
	return filepath.Join(home, "census", poolkey.Key(root))
}

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
	if benchguard.InvokesBench(command, resolver) {
		return nil
	}
	prefix := poolkey.Pool(home, root) + string(filepath.Separator)
	id, ok := assignment(command, prefix)
	if !ok {
		return nil
	}
	head := verbHead(command, prefix)
	if head == "" {
		return nil
	}
	return write(Dir(home, root), id, composeRecord(now.UTC().Format(time.RFC3339), head))
}

// recordSeparator divides a record line's two fields.
const recordSeparator = "\t"

// composeRecord renders one record line: the time, the separator, the head, and the
// newline that closes the record. The head passes through sanitize.Controls here,
// because a raw tab or newline in it would forge a second field or a second record.
// This function and parseRecord are the one source of the line's layout, so no other
// call site states where the head sits.
func composeRecord(timestamp, head string) string {
	return timestamp + recordSeparator + sanitize.Controls(head) + "\n"
}

// parseRecord returns the head one record line carries. The head is the second field,
// and a later field is not part of it. A line with no separator, and a line whose head
// field is empty, carry no head at all; only a foreign writer makes either.
func parseRecord(line string) (string, bool) {
	_, rest, ok := strings.Cut(line, recordSeparator)
	if !ok {
		return "", false
	}
	head, _, _ := strings.Cut(rest, recordSeparator)
	if head == "" {
		return "", false
	}
	return head, true
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

// write appends one line to the assignment's record file. The directory is never
// followed through a symlink, because a redirected census writes outside the Bench
// home. The single append write keeps concurrent writers' lines intact.
func write(dir, id, line string) error {
	if info, err := os.Lstat(dir); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("census directory %s is a symlink", dir)
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create census directory: %w", err)
	}
	file, err := os.OpenFile(filepath.Join(dir, id), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("open census record: %w", err)
	}
	defer file.Close()
	if _, err := file.WriteString(line); err != nil {
		return fmt.Errorf("append census record: %w", err)
	}
	return nil
}

// Counts returns the number of records each assignment holds under root's census
// directory. The ambient board is the caller, so no condition on the disk becomes a
// board failure: an absent directory, an absent or empty file, and a file type
// readRecords refuses each count zero. A last line with no newline still counts as one
// record, because a concurrent writer can be between its two writes.
func Counts(home, root string) (map[string]int, error) {
	counts := map[string]int{}
	entries, err := os.ReadDir(Dir(home, root))
	if err != nil {
		if os.IsNotExist(err) {
			return counts, nil
		}
		return counts, fmt.Errorf("read census directory: %w", err)
	}
	for _, entry := range entries {
		if !isAssignmentID(entry.Name()) {
			continue
		}
		text, ok := readRecords(Dir(home, root), entry.Name())
		if !ok {
			continue
		}
		counts[entry.Name()] = lineCount(text)
	}
	return counts, nil
}

// HeadBreakdown renders one assignment's raw calls per verb head, in the shape
// `<head>=<count>,...`. The heads sort by count, largest first, and a tie sorts by
// the head name, so the reader sees the heaviest verb first. An assignment with no
// records renders no text, because the landing prints the line only where there is
// evidence to print. The reader has the same posture as Counts, because both read
// through readRecords: an absent directory, an absent file, and a file type the
// reader refuses each render no text, because the census is evidence beside the
// landing and never a condition on it.
func HeadBreakdown(home, root, assignment string) string {
	return renderHeads(heads(home, root, assignment))
}

// heads returns the number of records each verb head holds in one assignment's file.
// The head comes from parseRecord, the same codec the writer composes with, so the two
// ends never disagree on where the head sits. A line with no head, which only a foreign
// writer makes, counts under no head.
func heads(home, root, assignment string) map[string]int {
	counts := map[string]int{}
	text, ok := readRecords(Dir(home, root), assignment)
	if !ok {
		return counts
	}
	for _, line := range recordLines(text) {
		head, ok := parseRecord(line)
		if !ok {
			continue
		}
		counts[head]++
	}
	return counts
}

// renderHeads joins the counted heads into the breakdown text. The head text passes
// through sanitize.Controls here rather than at the writer alone, because a record
// file that a foreign writer changed must still print safely.
func renderHeads(counts map[string]int) string {
	names := slices.Collect(maps.Keys(counts))
	slices.SortFunc(names, func(a, b string) int {
		if counts[a] != counts[b] {
			return counts[b] - counts[a]
		}
		return strings.Compare(a, b)
	})
	parts := make([]string, 0, len(names))
	for _, name := range names {
		parts = append(parts, fmt.Sprintf("%s=%d", escapeDelimiters(sanitize.Controls(name)), counts[name]))
	}
	return strings.Join(parts, ",")
}

// escapeDelimiters backslash-escapes the two characters the breakdown line uses as its
// own grammar. sanitize.Controls runs first and this escape second, never the other way
// around: Controls has already doubled every literal backslash by then, so each
// backslash this step adds is an escape introducer and a left-to-right reader recovers
// the head exactly. The reverse order doubles the added backslash, and the head forges a
// delimiter again.
var delimiterEscaper = strings.NewReplacer(",", `\,`, "=", `\=`)

func escapeDelimiters(head string) string { return delimiterEscaper.Replace(head) }

// readRecords returns one assignment's record text under dir. Only a regular file is
// read: a directory, a symlink, a FIFO, or a device at the record's name is foreign,
// and an open of a FIFO blocks the caller forever. An absent file and a refused file
// type each read as no text at all, which every census reader states as no records.
func readRecords(dir, name string) (string, bool) {
	path := filepath.Join(dir, name)
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() {
		return "", false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	return string(data), true
}

// isAssignmentID reports whether name is one 32-hex assignment id. The test composes
// the segment pair rather than restating the hexadecimal shape, so poolkey stays the
// one source of what an identifier looks like.
func isAssignmentID(name string) bool {
	_, ok := poolkey.SplitAssignmentSegment(poolkey.AssignmentSegment(name, name))
	return ok
}

// lineCount returns the number of records in one file's text.
func lineCount(text string) int {
	return len(recordLines(text))
}

// recordLines splits one file's text into its record lines. A trailing newline closes
// the last record and opens none of its own. A last line with no newline is still one
// record, because a concurrent writer can be between its two writes.
func recordLines(text string) []string {
	if text == "" {
		return nil
	}
	return strings.Split(strings.TrimSuffix(text, "\n"), "\n")
}

// Drop removes one assignment's record file, which the assignment's retirement calls.
// An absent file is not an error, so an assignment that made no raw call and a
// retirement that runs twice both complete. An identifier that is not an assignment
// id is refused rather than composed into a path, because the removal is
// unrecoverable.
func Drop(home, root, assignment string) error {
	if !isAssignmentID(assignment) {
		return fmt.Errorf("census assignment id is malformed: %s", sanitize.Controls(assignment))
	}
	if err := os.Remove(filepath.Join(Dir(home, root), assignment)); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove census record: %w", err)
	}
	return nil
}
