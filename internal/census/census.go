// Package census owns the record of raw shell calls that touch a Bench worktree.
package census

import (
	"cmp"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/poolkey"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// Dir returns the directory that holds one repository's census records. The
// directory is a sibling of the worktree pool and never a child of it, because
// `bench worktree reclaim` enumerates `<home>/worktrees` and a foreign entry
// there changes its plan.
func Dir(home, root string) string {
	return filepath.Join(home, "census", poolkey.Key(root))
}

// recordSeparator divides the fields of a record line.
const recordSeparator = "\t"

// composeRecord renders one record line: the time, the separator, the head, each later
// field after its own separator, and the newline that closes the record. The head passes
// through sanitize.Controls here, because a raw tab or newline in it would forge a second
// field or a second record. This function and parseRecord are the one source of the
// line's layout, so no other call site states where the head sits.
func composeRecord(timestamp, head string, later ...string) string {
	return strings.Join(append([]string{timestamp, sanitize.Controls(head)}, later...), recordSeparator) + "\n"
}

// parseRecord returns the head one record line carries. The head is the second field,
// and a later field is not part of it. A line with no separator, and a line whose head
// field is empty, carry no head at all; only a foreign writer makes either.
func parseRecord(line string) (string, bool) {
	_, head, ok := recordFields(line)
	return head, ok
}
func recordFields(line string) (string, string, bool) {
	timestamp, head, _, ok := splitRecord(line)
	return timestamp, head, ok
}

// splitRecord answers the time, the head, and the later fields of one record line. Every
// census reader parses here, so no reader restates where the head sits or when it is absent.
func splitRecord(line string) (string, string, []string, bool) {
	fields := strings.Split(line, recordSeparator)
	if len(fields) < 2 || fields[1] == "" {
		return "", "", nil, false
	}
	return fields[0], fields[1], fields[2:], true
}

// write appends one line to the named record file of an assignment. The directory is
// never followed through a symlink, because a redirected census writes outside the Bench
// home. The single append write keeps concurrent writers' lines intact.
func write(dir, name, line string) error {
	if info, err := os.Lstat(dir); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("census directory %s is a symlink", dir)
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create census directory: %w", err)
	}
	file, err := os.OpenFile(filepath.Join(dir, name), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
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
		if !poolkey.IsAssignmentID(entry.Name()) {
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
	counts := heads(home, root, assignment)
	return renderBreakdown(counts, func(head string) string { return strconv.Itoa(counts[head]) })
}

// heads returns the number of records each verb head holds in one assignment's file.
// The head comes from parseRecord, the same codec the writer composes with, so the two
// ends never disagree on where the head sits. A line with no head, which only a foreign
// writer makes, counts under no head.
func heads(home, root, assignment string) map[string]int {
	counts, _ := tally(Dir(home, root), assignment, func(line string) (string, int64, bool) {
		head, ok := parseRecord(line)
		return head, 0, ok
	})
	return counts
}

// tally reads one record file and answers, for each verb head, the number of lines that
// parse reads and the sum of the weight parse reads from them. A line that parse refuses
// counts under no head.
func tally(dir, name string, parse func(string) (string, int64, bool)) (map[string]int, map[string]int64) {
	calls, weights := map[string]int{}, map[string]int64{}
	text, ok := readRecords(dir, name)
	if !ok {
		return calls, weights
	}
	for _, line := range recordLines(text) {
		head, weight, ok := parse(line)
		if !ok {
			continue
		}
		calls[head]++
		weights[head] += weight
	}
	return calls, weights
}

// renderBreakdown joins the heads into one breakdown text, `<head>=<value>,...`. The heads
// sort by rank, largest first, and a tie sorts by the head name. The head text passes
// through sanitize.Controls here rather than at the writer alone, because a record file
// that a foreign writer changed must still print safely.
func renderBreakdown[Rank int | int64](ranks map[string]Rank, value func(string) string) string {
	names := slices.Collect(maps.Keys(ranks))
	slices.SortFunc(names, func(a, b string) int {
		if ranks[a] != ranks[b] {
			return cmp.Compare(ranks[b], ranks[a])
		}
		return strings.Compare(a, b)
	})
	parts := make([]string, 0, len(names))
	for _, name := range names {
		parts = append(parts, escapeDelimiters(sanitize.Controls(name))+"="+value(name))
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

// isAssignmentID is the poolkey predicate under the name the event reader calls.
var isAssignmentID = poolkey.IsAssignmentID

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

// Drop removes one assignment's raw-call file and output file, which the assignment's
// retirement calls. An absent file is not an error, so an assignment that made no call
// and a retirement that runs twice both complete. An identifier that is not an
// assignment id is refused rather than composed into a path, because the removal is
// unrecoverable.
func Drop(home, root, assignment string) error {
	if err := checkAssignmentID(assignment); err != nil {
		return err
	}
	for _, name := range []string{assignment, assignment + outputSuffix} {
		if err := os.Remove(filepath.Join(Dir(home, root), name)); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove census record: %w", err)
		}
	}
	return nil
}

// checkAssignmentID refuses an identifier that is not an assignment id. The output writer
// and the drop call it before they compose a file name, so no operand leaves the directory.
func checkAssignmentID(assignment string) error {
	if !poolkey.IsAssignmentID(assignment) {
		return fmt.Errorf("census assignment id is malformed: %s", sanitize.Controls(assignment))
	}
	return nil
}
