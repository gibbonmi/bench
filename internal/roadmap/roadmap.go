// Package roadmap owns idea capture, roadmap display, and parked-line counting.
package roadmap

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

const ideasFile = "IDEAS.md"

// IdeaCommand implements `bench idea <text...>`: it appends a dated line to IDEAS.md.
// The args are joined with single spaces; an empty or all-whitespace text yields the
// usage string on exit 2 without touching the file. Otherwise it resolves the repo
// root, normalizes a missing trailing newline (so a hand-edited last line without one
// does not swallow the new entry onto its physical line), then appends
// `- <ISO date>  <text>` — two spaces between date and text — creating the file if absent.
func IdeaCommand(args []string) (string, int) {
	text := strings.Join(args, " ")
	if strings.TrimSpace(text) == "" {
		return "usage: bench idea \"<text>\"\n", 2
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	file := filepath.Join(root, ideasFile)

	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return cannotWriteIdeas(err), 1
	}
	defer f.Close()

	if needsNewline(file) {
		if _, err := f.WriteString("\n"); err != nil {
			return cannotWriteIdeas(err), 1
		}
	}
	line := "- " + time.Now().Format("2006-01-02") + "  " + text + "\n"
	if _, err := f.WriteString(line); err != nil {
		return cannotWriteIdeas(err), 1
	}
	return "parked: " + text + "\n", 0
}

func cannotWriteIdeas(err error) string {
	return toon.Errorf("cannot write "+ideasFile, err.Error()) + "\n"
}

// needsNewline reports whether the file is non-empty and its last byte is not a
// newline — the case where an appended line would merge onto a hand-edited last line.
func needsNewline(file string) bool {
	info, err := os.Stat(file)
	if err != nil || info.Size() == 0 {
		return false
	}
	data, err := os.ReadFile(file)
	if err != nil || len(data) == 0 {
		return false
	}
	return data[len(data)-1] != '\n'
}

// RoadmapCommand implements `bench roadmap`: it prints ROADMAP.md verbatim, or the
// literal `roadmap empty` when the file is absent or zero bytes.
func RoadmapCommand(args []string) (string, int) {
	root, err := git.Root()
	if err != nil {
		return "roadmap empty\n", 0
	}
	data, err := os.ReadFile(filepath.Join(root, "ROADMAP.md"))
	if err != nil || len(data) == 0 {
		return "roadmap empty\n", 0
	}
	return string(data), 0
}

// ParkedCount returns the number of `^- ` lines (a hyphen then a space at line start)
// in <root>/ROADMAP.md — the parked-idea figure the `bench status` footer shows. A
// missing or unreadable file counts as zero.
func ParkedCount(root string) int {
	data, err := os.ReadFile(filepath.Join(root, "ROADMAP.md"))
	if err != nil {
		return 0
	}
	count := 0
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "- ") {
			count++
		}
	}
	return count
}
