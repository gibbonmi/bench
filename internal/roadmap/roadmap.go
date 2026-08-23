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
	return "usage: bench roadmap | bench roadmap --context [--full] | bench roadmap --context --row <ID,...> | bench roadmap --flow\n"
}

// ideaGrammar declares the argument shape usage.Parse enforces for this subcommand.
// Arity, flag recognition, `--`, and help all come from this grammar, not a local switch.
// The text is variadic, so MaxArgs is unbounded. `--` lets idea text begin with a dash.
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
// dispatches to the --context form, which declares its own grammar. This grammar
// therefore takes no arguments.
var roadmapGrammar = usage.Grammar{
	Cmd:  "bench roadmap",
	Help: strings.TrimSuffix(roadmapUsage(), "\n"),
}

// IdeasFile and RoadmapFile are the two repo-relative control records this package owns.
// Both are exported because the status board names them in the rows it prints when their
// reads fail. A literal repeated there would be a second derivation of a name this
// package decides.
const (
	IdeasFile   = "capture/IDEAS.md"
	RoadmapFile = "ROADMAP.md"
)

// IdeaCommand implements `bench idea <text...>` and appends a dated line to
// capture/IDEAS.md. It joins the args with single spaces. An empty or all-whitespace text
// yields the usage string on exit 2 and does not touch the file. Otherwise it resolves
// the repo root and normalizes a missing trailing newline, so a hand-edited last line
// does not absorb the new entry. It then appends `- <ISO date>  <text>`, with two spaces
// between date and text, and creates the file if absent.
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
	// The inbox lives in a directory a fresh repo may not have yet. A parked idea is often
	// the first thing to touch this directory.
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
		// This code performs one write, not two. An interrupt between two separate writes would
		// leave a bare blank line with no entry behind it.
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

// needsNewline reports whether the file is non-empty and its last byte is not a newline.
// In that case, an appended line would merge onto a hand-edited last line.
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
		// The classifier carries no diagnostic for a clean read of nothing. The error line
		// supplies the one fact that separates this state from absence.
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

// Drain is the typed capture-source snapshot that `bench roadmap` and `bench status`
// project. It carries each pending count with its source state. Absence and empty stay
// the ordinary quiet posture, while a failed read renders fail-closed unknown evidence.
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
// working document. The dashboard renders this text. A false present flag drives the
// dashboard's definitive empty state. Every state but a clean non-empty read lands there,
// so the page degrades instead of failing, unlike the `bench roadmap` command.
func RoadmapText(root string) (text string, present bool) {
	tree := LoadTree(root)
	if tree.Index.State != bounds.StateParsed {
		return "", false
	}
	doc, _, _ := ParseDocument(tree, nil, true)
	return doc.Text, true
}

// ParkedIdeas returns the parked idea lines from capture/IDEAS.md, every line beginning
// `- `. DrainCounts tallies these same lines, since both functions read through
// ideaLines, one source. A file that did not read as a usable document yields nil, which
// the dashboard renders as its empty state.
func ParkedIdeas(root string) []string {
	parked, _ := ideaLines(filepath.Join(root, IdeasFile))
	return parked
}

// RecommendedSequence extracts the `## Recommended sequence` section, from its heading,
// with trailing whitespace tolerated, to the next `## ` heading or EOF. This function
// tracks fence state throughout. A heading inside a fenced code block neither starts the
// section nor ends it. Both `bench roadmap`'s next-action callout and the dashboard's
// sequence block read this section, so the two share one parser.
func RecommendedSequence(roadmap string) string {
	_, text, _ := parseSequence(strings.Split(roadmap, "\n"))
	return text
}

// learningCount classifies capture/learnings.md and counts its parser-approved open rows.
// Absent is the quiet-journal posture. Every present non-document state is retained, so
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

// ideaLines is the one reader of capture/IDEAS.md-style parked lines, every line
// beginning `- `. It returns the lines with the file's own readability state. DrainCounts
// tallies these lines, and ParkedIdeas returns them for the dashboard, so the count and
// the rendered list never disagree. Absent or empty yields no lines at
// bounds.StateParsed, the ordinary quiet-inbox posture. This is the same
// absent-and-empty-is-quiet, failed-read-is-reported contract as learningCount.
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

// rowNextTokens is the ordered set of values a row's `Next:` line may carry, one token
// per phase a row can reach. It is the one source for this grammar. The split-board
// parser validates against it, and the drain's token table is graded against it. The
// enforcement and its documentation therefore cannot drift apart.
var rowNextTokens = []string{"shape", "spec", "ticket", "decide", "kit-edit"}

// RowNextTokens returns the ordered row-token set. It hands back a fresh slice, so a
// caller that sorts or trims its copy cannot rewrite the grammar for everyone else.
func RowNextTokens() []string { return append([]string(nil), rowNextTokens...) }
