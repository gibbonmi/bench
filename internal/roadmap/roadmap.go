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
	"github.com/gibbonmi/bench/internal/retros"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

func roadmapUsage() string {
	return "usage: bench roadmap | bench roadmap --context [--full] | bench roadmap --context --row <ID,...>\n"
}

// ideaGrammar is the declared argument shape usage.Parse enforces for this subcommand —
// arity, flag recognition, `--`, and help all come from there rather than a local switch.
// The text is variadic, so MaxArgs is unbounded; `--` is what makes idea text that
// begins with a dash expressible.
var ideaGrammar = usage.Grammar{
	Cmd:  "bench idea",
	Help: `usage: bench idea "<text>"`,
	Flags: []usage.Flag{
		{Name: "--owner", HasValue: true, NoEmptyValue: true},
		{Name: "--incident", HasValue: true, NoEmptyValue: true},
	},
	MaxArgs: -1,
}

var occurrenceOwner = regexp.MustCompile(`^FT[1-9][0-9]*$`)

// ValidOccurrenceOwner reports whether owner has the shared occurrence-owner grammar.
func ValidOccurrenceOwner(owner string) bool { return occurrenceOwner.MatchString(owner) }

// roadmapGrammar is the bare `bench roadmap` form. Every argument-bearing invocation
// is dispatched to the --context form, which declares its own grammar, so this one
// takes nothing at all.
var roadmapGrammar = usage.Grammar{
	Cmd:  "bench roadmap",
	Help: strings.TrimSuffix(roadmapUsage(), "\n"),
}

// IdeasFile and RoadmapFile are the two repo-relative control records this package
// owns. Both are exported because the status board names them in the rows it prints
// when their reads fail, and a literal repeated there is a second derivation of a name
// this package decides.
const (
	IdeasFile   = "capture/IDEAS.md"
	RoadmapFile = "ROADMAP.md"
)

// IdeaCommand implements `bench idea <text...>`: it appends a dated line to capture/IDEAS.md.
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
	displayText := text
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	owner, hasOwner := parsed.Flags["--owner"]
	incident, hasIncident := parsed.Flags["--incident"]
	if hasOwner != hasIncident {
		return toon.MissingArg(ideaGrammar.Cmd, "--owner and --incident") + "\n", 2
	}
	if hasOwner {
		if !ValidOccurrenceOwner(owner) {
			return toon.Usage(ideaGrammar.Cmd, "--owner "+owner) + "\n", 2
		}
		if !ValidOccurrenceIncident(incident) {
			return toon.Usage(ideaGrammar.Cmd, "--incident "+incident) + "\n", 2
		}
		if err := validateOccurrenceOwner(root, owner); err != nil {
			return toon.Errorf("cannot validate occurrence owner", err.Error()) + "\n", 1
		}
		text += " [occurrence:" + owner + "/" + incident + "]"
	}
	file := filepath.Join(root, IdeasFile)
	// The inbox lives in a directory a fresh repo may not have yet, and parking an idea
	// is often the first thing that touches it.
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		return cannotWriteIdeas(err), 1
	}

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
	return "parked: " + displayText + "\n", 0
}

func validateOccurrenceOwner(root, owner string) error {
	tree := LoadTree(root)
	if tree.Index.State != bounds.StateParsed {
		return fmt.Errorf("%s is not a trusted current roadmap", RoadmapFile)
	}
	doc, failures, diagnostics := ParseDocument(tree, nil, true)
	if len(failures) != 0 || len(diagnostics) != 0 || len(doc.OccurrenceDiscrepancies) != 0 {
		return fmt.Errorf("%s is structurally untrusted", RoadmapFile)
	}
	for _, row := range doc.Rows {
		if row.ID != owner {
			continue
		}
		return nil
	}
	return fmt.Errorf("owner %s is absent from the current roadmap", owner)
}

func cannotWriteIdeas(err error) string {
	return toon.Errorf("cannot write "+IdeasFile, err.Error()) + "\n"
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

// RoadmapCommand implements the bounded top-of-board projection for `bench roadmap`.
func RoadmapCommand(args []string) (string, int) {
	if _, line, code := usage.Parse(roadmapGrammar, args); line != "" {
		return line + "\n", code
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	tree := LoadTree(root)
	switch {
	case tree.Index.State == bounds.StateAbsent:
		return renderRoadmapBoard(Document{}, nil, DrainCounts(root), true, tree.DirState)
	case tree.Index.State == bounds.StateEmpty:
		// The classifier carries no diagnostic for a clean read of nothing, so the
		// error line supplies the one fact that separates this from absence.
		return toon.RecordError(RoadmapFile, tree.Index.State, "the file exists but holds no bytes") + "\n", 1
	case tree.Index.State.Failed():
		return toon.RecordError(RoadmapFile, tree.Index.State, tree.Index.Reason) + "\n", 1
	}
	doc, failures, diagnostics := ParseDocument(tree, nil, true)
	for _, f := range failures {
		if f.Reason == noRoadmapRowsReason {
			return toon.RecordError(RoadmapFile, bounds.StateUnsupportedSchema, f.Reason) + "\n", 1
		}
	}
	return renderRoadmapBoard(doc, diagnostics, DrainCounts(root), false, tree.DirState)
}

func renderRoadmapBoard(doc Document, diagnostics []Diagnostic, drain Drain, absent bool, dirState bounds.FileState) (string, int) {
	rowsShown := len(doc.Rows)
	if rowsShown > 10 {
		rowsShown = 10
	}
	roadmapRows := make([][]any, rowsShown)
	for i := range roadmapRows {
		r := doc.Rows[i]
		roadmapRows[i] = []any{r.ID, r.Title, r.Spec, r.SpecStatus, r.ExternalTrigger, r.OccurrenceCount, r.OccurrenceKeys}
	}
	sequenceRows := make([][]any, len(doc.Sequence))
	for i, r := range doc.Sequence {
		sequenceRows[i] = []any{r.Rank, r.Text, r.Command}
	}
	sources := []SourceFact{
		{Source: RoadmapFile, State: string(bounds.StateParsed)},
		{Source: RoadmapDir + "/", State: string(dirState)},
		{Source: IdeasFile, State: string(drain.IdeasState)},
		{Source: learnings.JournalPath, State: string(drain.LearningsState)},
		{Source: retros.Directory + "/", State: string(drain.RetrosState)},
	}
	trusted := occurrenceSequenceTrusted(doc.OccurrenceDiscrepancies, diagnostics, sources)
	blocks := []struct {
		name   string
		fields []string
		rows   [][]any
	}{
		{"roadmap", []string{"id", "title", "spec", "spec_status", "external_trigger", "occurrence_count", "occurrence_keys"}, roadmapRows},
		{"board", []string{"rows_shown", "rows_total", "sequence_trusted"}, [][]any{{rowsShown, len(doc.Rows), trusted}}},
		{"sequence", []string{"rank", "text", "command"}, sequenceRows},
		{"drain", []string{"ideas", "ideas_state", "learnings", "learnings_state", "retros", "retros_state"}, [][]any{{drain.Ideas, string(drain.IdeasState), drain.OpenLearnings, string(drain.LearningsState), drain.Retros, string(drain.RetrosState)}}},
	}
	if absent || drain.Ideas != 0 || drain.OpenLearnings != 0 || drain.Retros != 0 {
		blocks = append(blocks, struct {
			name   string
			fields []string
			rows   [][]any
		}{"help", []string{"cmd", "why"}, [][]any{{"/bench-drain", "create or drain the working roadmap"}}})
	} else {
		blocks = append(blocks, struct {
			name   string
			fields []string
			rows   [][]any
		}{"help", []string{"cmd", "why"}, nil})
	}
	var out strings.Builder
	for _, block := range blocks {
		text, err := toon.TableTyped(block.name, block.fields, block.rows)
		if err != nil {
			return toon.RenderError(err) + "\n", 1
		}
		out.WriteString(text)
	}
	return out.String(), 0
}

// Drain is the typed capture-source snapshot `bench roadmap` and `bench status` project.
// It carries each pending count and its source state so absence and empty remain the
// ordinary quiet posture while failed reads render fail-closed unknown evidence.
type Drain struct {
	Ideas, OpenLearnings, Retros            int
	IdeasState, LearningsState, RetrosState bounds.FileState
}

// DrainCounts gathers every capture source used by the roadmap and status surfaces.
func DrainCounts(root string) Drain {
	parked, ideasState := ideaLines(filepath.Join(root, IdeasFile))
	openLearnings, learningsState := learningCount(root)
	retroFacts := retros.Facts(root)
	return Drain{
		Ideas: len(parked), OpenLearnings: openLearnings, Retros: len(retroFacts.Entries),
		IdeasState: ideasState, LearningsState: learningsState, RetrosState: retroFacts.State,
	}
}

// RoadmapText returns ROADMAP.md's rendered contents and whether the file yielded a
// working document at all. The dashboard renders this text; a false present flag drives
// its definitive empty state, which is where every state but a clean non-empty read
// lands — the page degrades rather than failing, unlike the `bench roadmap` command.
func RoadmapText(root string) (text string, present bool) {
	tree := LoadTree(root)
	if tree.Index.State != bounds.StateParsed {
		return "", false
	}
	doc, _, _ := ParseDocument(tree, nil, true)
	return doc.Text, true
}

// ParkedIdeas returns the parked idea lines from capture/IDEAS.md — every line beginning `- `, the
// same lines DrainCounts tallies (both go through ideaLines, one source). A file that did
// not read as a usable document yields nil, which the dashboard renders as its empty state.
func ParkedIdeas(root string) []string {
	parked, _ := ideaLines(filepath.Join(root, IdeasFile))
	return parked
}

// RecommendedSequence extracts the `## Recommended sequence` section, from its
// heading (trailing whitespace tolerated) to the next `## ` heading or EOF. Fence
// state is tracked throughout: a heading inside a fenced code block neither starts
// the section nor terminates it. Both `bench roadmap`'s next-action callout and the
// dashboard's sequence block read it, so the two share one parser.
func RecommendedSequence(roadmap string) string {
	_, text, _ := parseSequence(strings.Split(roadmap, "\n"))
	return text
}

// learningCount classifies capture/learnings.md and counts its parser-approved open rows.
// Absent is the quiet-journal posture; every present non-document state is retained so
// the drain surfaces unknown evidence instead of a fabricated clean journal.
func learningCount(root string) (int, bounds.FileState) {
	c := bounds.Classify(filepath.Join(root, filepath.FromSlash(learnings.JournalPath)), bounds.ControlRecordLimit)
	switch {
	case c.State == bounds.StateAbsent:
		return 0, bounds.StateParsed
	case c.State != bounds.StateParsed:
		return 0, c.State
	}
	_, malformed := learnings.Parse(c.Data)
	if len(malformed) > 0 {
		return 0, bounds.StateMalformed
	}
	return len(learnings.Rows(c.Data)), bounds.StateParsed
}

// ideaLines is the one reader of capture/IDEAS.md-style parked lines: every line beginning `- `,
// returned with the file's own readability state. DrainCounts tallies them and
// ParkedIdeas returns them for the dashboard, so the count and the rendered list can
// never disagree about either the lines or what was readable. Absent or empty is no
// lines at bounds.StateParsed (the ordinary quiet-inbox posture), the same
// absent/empty-is-quiet, failed-read-is-reported contract as learningCount.
func ideaLines(file string) ([]string, bounds.FileState) {
	c := bounds.Classify(file, bounds.ControlRecordLimit)
	switch {
	case c.State.Failed():
		return nil, c.State
	case c.State == bounds.StateAbsent || c.State == bounds.StateEmpty:
		return nil, bounds.StateParsed
	}
	_, _, out := parseIdeas(c.Data, true)
	return out, bounds.StateParsed
}
