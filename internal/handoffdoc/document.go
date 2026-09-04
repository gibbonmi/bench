// Package handoffdoc owns the grammar of the session handoff document: a header
// block, one "## main" section, one "## request <digest>" section per live
// assignment, and a trailing "## Shape". It parses, renders, rewrites one
// section, removes a section by key, and ensures "main". It also owns the locked
// read-modify-write that lets two live phases write the file at once.
//
// The package is a leaf and imports the standard library only. Three packages
// read this grammar — the handoff command, the worktree retirement path, and the
// status row — and internal/worktree cannot import internal/handoff. A grammar
// that lived in either consumer would be re-derived in the other.
package handoffdoc

import (
	"fmt"
	"strings"
)

// MainKey is the key of the section every document carries. It is the section a
// primary-checkout close owns, and the one the file falls back to after the last
// assignment section is removed.
const MainKey = "main"

// The section headings the renderer emits and the parser reads back.
// RequestHeadingPrefix carries the trailing space, so the digest begins at the
// byte after it. The three exported spellings are what the handoff's Shape text
// describes, so the single-source check grades that text against these rather
// than against a second copy of them.
const (
	MainHeading          = "## " + MainKey
	RequestHeadingPrefix = "## request "
	StateHeading         = "### State"
	shapeHeading         = "## Shape"
)

// The label lines a section carries. The six pins are exported because the
// command that fills a section names them, and a second spelling there would
// render a field this parser reads back under a different label.
const (
	LabelRequestToken = "Request token"
	LabelLabel        = "Label"
	LabelWorktreeTip  = "Worktree tip"
	LabelRecordedBase = "Recorded base"
	LabelSpec         = "Spec"
	LabelSpecStatus   = "Spec status"
	LabelNextCommand  = "Next command"
)

// Field is one label line. The value is rendered verbatim after "Label: ", and a
// line carries no sentence terminator of its own. That is what makes the prose
// lane skip it: an unterminated label-shaped line is a template field, not prose.
type Field struct {
	Label string
	Value string
}

// Section is one assignment's block. Fields holds the pins in render order, and
// Spec with Spec status may repeat as pairs when a worktree holds two live specs.
// Next is the "Next command" value, kept apart because the command preserves a
// non-empty one byte for byte. State is the reviewer-owned body below "### State",
// which nothing in this package interprets.
type Section struct {
	Key    string
	Fields []Field
	Next   string
	State  string
}

// Document is the whole file. Header is everything above the first section, and
// Shape is the body of the trailing "## Shape". Sections render in slice order,
// with main first.
type Document struct {
	Header   string
	Sections []Section
	Shape    string
}

// ParseError names the file and the line a refusal points at. Ambiguity is
// refused rather than resolved: the alternative is a rewrite that writes over a
// section this parser guessed wrong about.
type ParseError struct {
	Path   string
	Line   int
	Reason string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%s:%d: %s", e.Path, e.Line, e.Reason)
}

// New returns a fresh document that holds main alone. It is what a first run in a
// repo starts from, so an absent file and an empty one take the same path.
func New() *Document {
	return &Document{Sections: []Section{{Key: MainKey}}}
}

// Section returns the section under key.
func (d *Document) Section(key string) (Section, bool) {
	for _, section := range d.Sections {
		if section.Key == key {
			return section, true
		}
	}
	return Section{}, false
}

// Put replaces the section under section.Key, or appends it. An appended request
// section lands after every existing one, so main keeps its place at the top.
func (d *Document) Put(section Section) {
	for i := range d.Sections {
		if d.Sections[i].Key == section.Key {
			d.Sections[i] = section
			return
		}
	}
	d.Sections = append(d.Sections, section)
}

// Remove drops the section under key and reports whether one was there.
func (d *Document) Remove(key string) bool {
	for i := range d.Sections {
		if d.Sections[i].Key != key {
			continue
		}
		d.Sections = append(d.Sections[:i], d.Sections[i+1:]...)
		return true
	}
	return false
}

// EnsureMain puts an empty main section at the top when the document holds none.
// The retirement path calls it after a removal, so the file still carries a State
// once the last assignment is gone.
func (d *Document) EnsureMain() {
	if _, found := d.Section(MainKey); found {
		return
	}
	d.Sections = append([]Section{{Key: MainKey}}, d.Sections...)
}

// Render writes the document back. The output is what Parse accepts, and a second
// render of a parsed document is byte-identical to the first: every blank line is
// re-derived rather than carried, so spacing cannot accumulate across runs.
func (d *Document) Render() string {
	var out strings.Builder
	if header := strings.Trim(d.Header, "\n"); header != "" {
		out.WriteString(header)
		out.WriteString("\n\n")
	}
	for _, section := range d.Sections {
		out.WriteString(heading(section.Key))
		out.WriteString("\n\n")
		for _, field := range section.Fields {
			out.WriteString(labelLine(field.Label, field.Value))
			out.WriteString("\n")
		}
		out.WriteString(labelLine(LabelNextCommand, section.Next))
		out.WriteString("\n\n")
		out.WriteString(StateHeading)
		out.WriteString("\n\n")
		if state := strings.Trim(section.State, "\n"); state != "" {
			out.WriteString(state)
			out.WriteString("\n\n")
		}
	}
	if shape := strings.Trim(d.Shape, "\n"); shape != "" {
		out.WriteString(shapeHeading)
		out.WriteString("\n\n")
		out.WriteString(shape)
		out.WriteString("\n")
	}
	return out.String()
}

// heading renders a section's own heading. main is spelled plainly; every other
// key is a request digest.
func heading(key string) string {
	if key == MainKey {
		return MainHeading
	}
	return RequestHeadingPrefix + key
}

// labelLine renders one field. An empty value drops the trailing space, because a
// line that ends in whitespace is one an editor silently rewrites, and the rewrite
// would make an unchanged run non-idempotent.
func labelLine(label, value string) string {
	if value == "" {
		return label + ":"
	}
	return label + ": " + value
}

// Parse reads a rendered document. It refuses an unknown or repeated section
// heading, a body line that is neither a label line nor part of State, and a
// document with no main section. Each refusal names the file and the line.
//
// Fenced blocks are skipped when locating a heading. A heading inside a fence is
// prose about the document, not a section of it.
//
// A document in the pre-section shape reads as main alone. See convertLegacy.
func Parse(path string, content []byte) (*Document, error) {
	lines := strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")
	doc := &Document{}
	blocks, header := splitSections(lines)
	doc.Header = strings.Trim(strings.Join(header, "\n"), "\n")
	blocks, err := convertLegacy(path, blocks)
	if err != nil {
		return nil, err
	}
	seen := map[string]int{}
	for _, block := range blocks {
		key, err := sectionKey(path, block.line, block.heading)
		if err != nil {
			return nil, err
		}
		if key == shapeHeading {
			doc.Shape = strings.Trim(strings.Join(block.body, "\n"), "\n")
			continue
		}
		if first, repeated := seen[key]; repeated {
			return nil, &ParseError{path, block.line, fmt.Sprintf("section %q is declared twice, first at line %d; leave one section per key", key, first)}
		}
		seen[key] = block.line
		section, err := parseSection(path, key, block)
		if err != nil {
			return nil, err
		}
		doc.Sections = append(doc.Sections, section)
	}
	if _, found := doc.Section(MainKey); !found {
		return nil, &ParseError{path, len(lines), fmt.Sprintf("document carries no %q section; every handoff holds one", MainHeading)}
	}
	return doc, nil
}

// block is one level-two heading, as the file spells it, and the raw lines below
// it. The heading stays unresolved here so that the legacy conversion can read a
// spelling this grammar has no key for.
type block struct {
	heading string
	line    int
	body    []string
}

// splitSections cuts the lines into the header and one block per level-two
// heading. It returns the blocks in file order, so main keeps whatever place the
// file gave it and Shape stays last.
func splitSections(lines []string) ([]block, []string) {
	var blocks []block
	var header []string
	fenced := false
	for i, raw := range lines {
		trimmed := strings.TrimRight(raw, " \t")
		if isFence(trimmed) {
			fenced = !fenced
		}
		if fenced || !strings.HasPrefix(trimmed, "## ") {
			if len(blocks) == 0 {
				header = append(header, raw)
			} else {
				blocks[len(blocks)-1].body = append(blocks[len(blocks)-1].body, raw)
			}
			continue
		}
		blocks = append(blocks, block{heading: trimmed, line: i + 1})
	}
	return blocks, header
}

// The level-two headings the handoff carried before it held one section per
// assignment. Both are derived from the spellings this package already owns, so a
// rename of either carries here rather than leaving a stale legacy spelling.
var (
	legacyStateHeading = strings.TrimPrefix(StateHeading, "#")
	legacyNextHeading  = "## " + LabelNextCommand
)

// convertLegacy rewrites a pre-section document into the blocks this grammar
// reads, and passes any other document through untouched. The legacy file held
// one State, one Next command, and further reviewer blocks above Shape; all of it
// belongs to main, because that shape predates the per-assignment section.
//
// The conversion emits blocks rather than a Section, so the rendered legacy body
// goes back through parseSection and no second reader of a section body exists.
func convertLegacy(path string, blocks []block) ([]block, error) {
	legacy, section := 0, ""
	for _, b := range blocks {
		if b.heading == legacyStateHeading && legacy == 0 {
			legacy = b.line
		}
		if (b.heading == MainHeading || strings.HasPrefix(b.heading, RequestHeadingPrefix)) && section == "" {
			section = b.heading
		}
	}
	if legacy == 0 {
		return blocks, nil
	}
	if section != "" {
		return nil, &ParseError{path, legacy, fmt.Sprintf("heading %q sits in a document that also carries %q; the legacy shape and the section shape do not mix", legacyStateHeading, section)}
	}

	var state, extras, next string
	var shape *block
	for _, b := range blocks {
		body := strings.Trim(strings.Join(b.body, "\n"), "\n")
		switch {
		case b.heading == shapeHeading:
			copied := b
			shape = &copied
		case b.heading == legacyStateHeading:
			state = body
		case b.heading == legacyNextHeading:
			if strings.Contains(body, "\n") {
				return nil, &ParseError{path, b.line, fmt.Sprintf("heading %q carries more than one line; a next command is one line", b.heading)}
			}
			next = body
		default:
			extras += "\n\n### " + strings.TrimPrefix(b.heading, "## ") + "\n\n" + body
		}
	}

	main := block{heading: MainHeading, line: legacy}
	main.body = append(main.body, labelLine(LabelNextCommand, next), "", StateHeading, "")
	main.body = append(main.body, strings.Split(state+extras, "\n")...)
	converted := []block{main}
	if shape != nil {
		converted = append(converted, *shape)
	}
	return converted, nil
}

// sectionKey resolves one heading to its key. Shape keeps its heading as the key,
// because it is the one block that holds no section.
func sectionKey(path string, line int, trimmed string) (string, error) {
	switch {
	case trimmed == MainHeading:
		return MainKey, nil
	case trimmed == shapeHeading:
		return shapeHeading, nil
	case strings.HasPrefix(trimmed, RequestHeadingPrefix):
		digest := strings.TrimSpace(strings.TrimPrefix(trimmed, RequestHeadingPrefix))
		if digest == "" || strings.ContainsAny(digest, " \t") {
			return "", &ParseError{path, line, fmt.Sprintf("heading %q names no single request digest", trimmed)}
		}
		return digest, nil
	}
	return "", &ParseError{path, line, fmt.Sprintf("heading %q is not %q, %q, or %q", trimmed, MainHeading, RequestHeadingPrefix+"<digest>", shapeHeading)}
}

// parseSection reads one block's label lines and its State body. Everything below
// the State heading is the reviewer's, and it passes through untouched.
func parseSection(path, key string, b block) (Section, error) {
	section := Section{Key: key}
	fenced := false
	for i, raw := range b.body {
		trimmed := strings.TrimRight(raw, " \t")
		if isFence(trimmed) {
			fenced = !fenced
		}
		if !fenced && trimmed == StateHeading {
			section.State = strings.Trim(strings.Join(b.body[i+1:], "\n"), "\n")
			return section, nil
		}
		if strings.TrimSpace(trimmed) == "" {
			continue
		}
		label, value, ok := splitLabel(trimmed)
		if !ok {
			return Section{}, &ParseError{path, b.line + i + 1, fmt.Sprintf("line %q is neither a label line nor part of %q", trimmed, StateHeading)}
		}
		if label == LabelNextCommand {
			section.Next = value
			continue
		}
		section.Fields = append(section.Fields, Field{Label: label, Value: value})
	}
	return Section{}, &ParseError{path, b.line, fmt.Sprintf("section %q carries no %q heading", key, StateHeading)}
}

// splitLabel cuts a label line at its first colon. The label is one to four words
// so that a sentence with a colon in it never reads as a field, which is the same
// bound the prose lane applies.
func splitLabel(line string) (string, string, bool) {
	colon := strings.IndexByte(line, ':')
	if colon <= 0 {
		return "", "", false
	}
	label := line[:colon]
	if n := len(strings.Fields(label)); n == 0 || n > 4 || strings.TrimSpace(label) != label {
		return "", "", false
	}
	value := strings.TrimSpace(line[colon+1:])
	return label, value, true
}

// isFence reports whether a line opens or closes a fenced block. Both markdown
// fence characters count, and leading indentation does not disqualify one.
func isFence(line string) bool {
	trimmed := strings.TrimLeft(line, " \t")
	return strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~")
}
