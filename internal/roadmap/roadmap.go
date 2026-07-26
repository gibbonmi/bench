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

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/learnings"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// ideaGrammar is the declared argument shape usage.Parse enforces for this subcommand —
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
		return ideaGrammar.Help + "\n", 2
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

// missingRoadmap is the absent-file posture: no working document at all is a maintenance
// prompt, never a crash or silent empty output. A ROADMAP.md that exists but holds no
// bytes is not this case — someone created it and left it unwritten, which the empty
// state reports instead.
const missingRoadmap = "no ROADMAP.md — run /bench-what-next to create the working roadmap\n"

// RoadmapCommand implements `bench roadmap`: it prints ROADMAP.md verbatim followed
// by the drain-status block when capture sources need draining, or by the
// `## Recommended sequence` callout when nothing does. Absence renders the
// maintenance prompt on exit 0 — the only authoritative empty state — while every
// other state exits 1 with a structured error naming it: a present-but-empty file,
// any classifier failure, or the parser's own unsupported-schema verdict. A read that
// did not produce a working document can therefore never print as one.
func RoadmapCommand(args []string) (string, int) {
	if _, line, code := usage.Parse(roadmapGrammar, args); line != "" {
		return line + "\n", code
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	c := bounds.Classify(filepath.Join(root, "ROADMAP.md"), controlRecordLimit)
	switch {
	case c.State == bounds.StateAbsent:
		return missingRoadmap, 0
	case c.State == bounds.StateEmpty:
		// The classifier carries no diagnostic for a clean read of nothing, so the
		// error line supplies the one fact that separates this from absence.
		return roadmapReadError(c.State, "the file exists but holds no bytes"), 1
	case c.State.Failed():
		return roadmapReadError(c.State, c.Reason), 1
	}
	doc, failures := ParseDocument(c.Data, nil, true)
	for _, f := range failures {
		if f.Reason == noRoadmapRowsReason {
			return roadmapReadError(bounds.StateUnsupportedSchema, f.Reason), 1
		}
	}
	text := doc.Text
	if status := drainStatus(root); status != "" {
		return text + status, 0
	}
	return text + nextAction(text), 0
}

// roadmapReadError renders the AXI error line every fail-closed classifier state
// produces for `bench roadmap`: `error: <path> is <state> — <reason>`, one grammar
// with every other migrated surface's error line.
func roadmapReadError(state bounds.FileState, reason string) string {
	return toon.Errorf("ROADMAP.md is "+string(state), reason) + "\n"
}

// DrainCounts returns the maintenance inbox counts `bench roadmap` reports before a
// reviewer trusts the roadmap sequence, plus each source's own readability state so a
// caller can distinguish a genuinely empty inbox from a failed read. Absent or empty is
// the ordinary quiet-inbox posture (count 0, no fail-closed state); only an unreadable
// or wrong-type source marks a failed read, matching maps.UnresolvedCount's contract.
func DrainCounts(root string) (ideas int, ideasState bounds.FileState, openLearnings int, learningsState bounds.FileState) {
	parked, ideasState := ideaLines(filepath.Join(root, ideasFile))
	openLearnings, learningsState = learningCount(root)
	return len(parked), ideasState, openLearnings, learningsState
}

// RoadmapText returns ROADMAP.md's rendered contents and whether the file yielded a
// working document at all. The dashboard renders this text; a false present flag drives
// its definitive empty state, which is where every state but a clean non-empty read
// lands — the page degrades rather than failing, unlike the `bench roadmap` command.
func RoadmapText(root string) (text string, present bool) {
	c := bounds.Classify(filepath.Join(root, "ROADMAP.md"), controlRecordLimit)
	if c.State != bounds.StateParsed {
		return "", false
	}
	doc, _ := ParseDocument(c.Data, nil, true)
	return doc.Text, true
}

// ParkedIdeas returns the parked idea lines from IDEAS.md — every line beginning `- `, the
// same lines DrainCounts tallies (both go through ideaLines, one source). A file that did
// not read as a usable document yields nil, which the dashboard renders as its empty state.
func ParkedIdeas(root string) []string {
	parked, _ := ideaLines(filepath.Join(root, ideasFile))
	return parked
}

func drainStatus(root string) string {
	// bench roadmap's status callout only needs the counts: a source that failed to
	// read renders as 0 here, and RoadmapCommand's own classified read above already
	// fails closed on ROADMAP.md itself, so an unreadable IDEAS.md or learnings.md
	// degrades this callout to quiet rather than the command's own read failing.
	ideas, _, openLearnings, _ := DrainCounts(root)
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

// learningCount classifies .bench/learnings.md and counts its open rows. Absent or
// empty is 0 rows at bounds.StateParsed (the ordinary quiet-journal posture); a journal
// whose read failed reports 0 with that failed state, so the caller renders unknown
// instead of a fabricated clean journal.
func learningCount(root string) (int, bounds.FileState) {
	c := bounds.Classify(filepath.Join(root, ".bench", "learnings.md"), controlRecordLimit)
	switch {
	case c.State.Failed():
		return 0, c.State
	case c.State == bounds.StateAbsent || c.State == bounds.StateEmpty:
		return 0, bounds.StateParsed
	default:
		return len(learnings.Rows(c.Data)), bounds.StateParsed
	}
}

// ideaLines is the one reader of IDEAS.md-style parked lines: every line beginning `- `,
// returned with the file's own readability state. DrainCounts tallies them and
// ParkedIdeas returns them for the dashboard, so the count and the rendered list can
// never disagree about either the lines or what was readable. Absent or empty is no
// lines at bounds.StateParsed (the ordinary quiet-inbox posture), the same
// absent/empty-is-quiet, failed-read-is-reported contract as learningCount.
func ideaLines(file string) ([]string, bounds.FileState) {
	c := bounds.Classify(file, controlRecordLimit)
	switch {
	case c.State.Failed():
		return nil, c.State
	case c.State == bounds.StateAbsent || c.State == bounds.StateEmpty:
		return nil, bounds.StateParsed
	}
	_, _, out := parseIdeas(c.Data, true)
	return out, bounds.StateParsed
}
