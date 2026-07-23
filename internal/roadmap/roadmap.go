// Package roadmap owns idea capture, roadmap display, and the drain counts the
// status board and roadmap command share.
package roadmap

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/learnings"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// grammar is the declared argument shape usage.Parse enforces for this subcommand —
// arity, flag recognition, `--`, and help all come from there rather than a local switch.
// The text is variadic, so MaxArgs is unbounded; `--` is what makes idea text that
// begins with a dash expressible.
var ideaGrammar = usage.Grammar{
	Cmd:     "bench idea",
	Help:    `usage: bench idea "<text>"`,
	MaxArgs: -1,
}

// roadmapGrammar is the bare `bench roadmap` form. Every argument-bearing invocation
// is dispatched to the --context form, which declares its own grammar, so this one
// takes nothing at all.
var roadmapGrammar = usage.Grammar{
	Cmd:  "bench roadmap",
	Help: "usage: bench roadmap",
}

const ideasFile = "IDEAS.md"

// IdeaCommand implements `bench idea <text...>`: it appends a dated line to IDEAS.md.
// The args are joined with single spaces; an empty or all-whitespace text yields the
// usage string on exit 2 without touching the file. Otherwise it resolves the repo
// root, normalizes a missing trailing newline (so a hand-edited last line without one
// does not swallow the new entry onto its physical line), then appends
// `- <ISO date>  <text>` — two spaces between date and text — creating the file if absent.
func IdeaCommand(args []string) (string, int) {
	parsed, line, code := usage.Parse(ideaGrammar, args)
	if line != "" {
		return line + "\n", code
	}
	text := strings.Join(parsed.Positionals, " ")
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

	entry := "- " + time.Now().Format("2006-01-02") + "  " + text + "\n"
	if needsNewline(file) {
		// One write, not two: an interrupt between a separate newline write and the
		// entry write would leave a bare blank line with no entry behind it.
		entry = "\n" + entry
	}
	if _, err := f.WriteString(entry); err != nil {
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

// missingRoadmap is the absent-or-zero-byte posture: no working document is a
// maintenance prompt, never a crash or silent empty output.
const missingRoadmap = "no ROADMAP.md — run /bench-what-next to create the working roadmap\n"

// RoadmapCommand implements `bench roadmap`: it prints ROADMAP.md verbatim followed
// by the drain-status block when capture sources need draining, or by the
// `## Recommended sequence` callout when nothing does. An absent or zero-byte
// roadmap, a missing sequence section, or a section without two-or-three numbered
// items each get an explicit message pointing at /bench-what-next, exit 0.
func RoadmapCommand(args []string) (string, int) {
	if _, line, code := usage.Parse(roadmapGrammar, args); line != "" {
		return line + "\n", code
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	data, err := os.ReadFile(filepath.Join(root, "ROADMAP.md"))
	if err != nil || len(data) == 0 {
		return missingRoadmap, 0
	}
	doc, _ := ParseDocument(data, nil, true)
	text := doc.Text
	if status := drainStatus(root); status != "" {
		return text + status, 0
	}
	return text + nextAction(text), 0
}

// DrainCounts returns the maintenance inbox counts `bench roadmap` reports before a
// reviewer trusts the roadmap sequence. Missing or unreadable files count as zero.
func DrainCounts(root string) (ideas, openLearnings int) {
	return lineCount(filepath.Join(root, ideasFile)), learningCount(root)
}

// RoadmapText returns ROADMAP.md's raw contents and whether it is present and non-empty —
// the same absent-or-zero-byte boundary RoadmapCommand keys its empty-state posture on.
// The dashboard renders this text; a false present flag drives its definitive empty state.
func RoadmapText(root string) (text string, present bool) {
	data, err := os.ReadFile(filepath.Join(root, "ROADMAP.md"))
	if err != nil || len(data) == 0 {
		return "", false
	}
	doc, _ := ParseDocument(data, nil, true)
	return doc.Text, true
}

// ParkedIdeas returns the parked idea lines from IDEAS.md — every line beginning `- `, the
// same lines DrainCounts tallies (both go through ideaLines, one source). A missing or
// unreadable file yields nil, which the dashboard renders as its empty state.
func ParkedIdeas(root string) []string {
	return ideaLines(filepath.Join(root, ideasFile))
}

func drainStatus(root string) string {
	ideas, openLearnings := DrainCounts(root)
	if ideas == 0 && openLearnings == 0 {
		return ""
	}
	return fmt.Sprintf("\n\n## Drain status\n\n- ideas: %d parked in %s\n- learnings: %d open in .bench/learnings.md\n\nRun /bench-what-next before trusting the sequence.\n", ideas, ideasFile, openLearnings)
}

// numberedItem matches one sequence entry; the format contract wants two or three.
var numberedItem = regexp.MustCompile(`(?m)^[0-9]+\. `)

// nextAction renders the no-drain call to action: the sequence section verbatim, or
// an explicit gap message so a broken format contract never yields silent output.
func nextAction(roadmap string) string {
	section := RecommendedSequence(roadmap)
	if section == "" {
		return "\n\nROADMAP.md has no ## Recommended sequence section — run /bench-what-next to restore it.\n"
	}
	if n := len(numberedItem.FindAllString(section, -1)); n < 2 || n > 3 {
		return fmt.Sprintf("\n\nmalformed ## Recommended sequence: %d numbered item(s), expected two or three — run /bench-what-next to repair it.\n", n)
	}
	return "\n\n## Next action\n\n" + section
}

// RecommendedSequence extracts the `## Recommended sequence` section, from its
// heading (trailing whitespace tolerated) to the next `## ` heading or EOF. Fence
// state is tracked throughout: a heading inside a fenced code block neither starts
// the section nor terminates it. Both `bench roadmap`'s next-action callout and the
// dashboard's sequence block read it, so the two share one parser.
func RecommendedSequence(roadmap string) string {
	doc, _ := ParseDocument([]byte(roadmap), nil, true)
	return doc.SequenceText
}

func learningCount(root string) int {
	data, err := os.ReadFile(filepath.Join(root, ".bench", "learnings.md"))
	if err != nil {
		return 0
	}
	return len(learnings.Rows(data))
}

func lineCount(file string) int {
	return len(ideaLines(file))
}

// ideaLines is the one reader of IDEAS.md-style parked lines: every line beginning `- `.
// lineCount tallies them for the drain counts and ParkedIdeas returns them for the
// dashboard, so the count and the rendered list can never disagree. Missing or unreadable
// file → nil.
func ideaLines(file string) []string {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil
	}
	_, _, out := parseIdeas(data, true)
	return out
}
