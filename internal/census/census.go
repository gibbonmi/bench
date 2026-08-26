// Package census owns the record of raw shell calls that touch a Bench worktree.
package census

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/benchguard"
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
	return write(Dir(home, root), id, now.UTC().Format(time.RFC3339)+"\t"+sanitize.Controls(head)+"\n")
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

// verbHead returns the head of the first simple command whose words name a pool path.
// A text that names the pool only where the parser cannot see it, such as a heredoc
// body, falls back to the first command's head.
func verbHead(command, prefix string) string {
	stream := shellcommand.Parse(command)
	first := ""
	for _, span := range stream.Commands {
		words := shellcommand.ProjectCommandWords(stream.Tokens[span.Start:span.End])
		if len(words) == 0 {
			continue
		}
		if first == "" {
			first = words[0]
		}
		for _, word := range words {
			if strings.Contains(word, prefix) {
				return words[0]
			}
		}
	}
	return first
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
