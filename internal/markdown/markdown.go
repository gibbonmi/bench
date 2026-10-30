// Package markdown is the one block reader for authored Markdown. It classifies each line of a
// document by four block rules: frontmatter, HTML comments, fenced blocks, and H2 headings. A
// grammar module reads its lines from here and keeps only its own grammar, so the block rules
// have one source.
//
// The reader never refuses. It reports at most one fault, and each caller decides what a fault
// means for its own grammar.
package markdown

import (
	"sort"
	"strings"
)

// FenceMarker is the shortest backtick fence marker, and H2Prefix opens an H2 heading. A writer
// renders a fence or a heading through them, so the writer keeps no block literal of its own.
const (
	FenceMarker = "```"
	H2Prefix    = "## "
)

// frontmatterMarker opens and closes the frontmatter block.
const frontmatterMarker = "---"

// Class names the block that holds one line.
type Class uint8

const (
	// Body is a line outside every block.
	Body Class = iota
	// Frontmatter is a line of the leading frontmatter block, its two marker lines included.
	Frontmatter
	// FenceOpener opens a fenced block.
	FenceOpener
	// Fenced is a line inside a fenced block.
	Fenced
	// FenceCloser closes a fenced block.
	FenceCloser
)

// FaultKind names the block that a fault leaves open.
type FaultKind uint8

const (
	FrontmatterFault FaultKind = iota
	CommentFault
	FenceFault
)

// Fault is one unterminated block. Line is the 1-based line of its opener.
type Fault struct {
	Kind FaultKind
	Line int
}

// Span is a range of source bytes, from Start up to but not including End.
type Span struct {
	Start, End int
}

// Line is one physical line of a document and its classification.
type Line struct {
	// Number is the 1-based physical line number, and Offset is the byte offset of the line
	// start in the source.
	Number, Offset int
	// Raw is the line without its newline and without one trailing carriage return.
	Raw string
	// Text is Raw with each removed comment byte cut out.
	Text string
	// Info is the trimmed text after the marker run of a fence opener or a fence closer. On an
	// opener, it is the info string.
	Info string
	// Section is the title of the open H2 section, and InSection reports that a section is
	// open. A heading line opens its own section, so its Section is its title.
	Section string
	Class   Class
	// Heading reports that the line opens an H2 section.
	Heading   bool
	InSection bool
}

// Document is one reading: the classified lines, the removed comment spans in source order, and
// at most one fault. The span of a comment that a removal rejoins holds the removed comment too.
type Document struct {
	Lines    []Line
	Comments []Span
	Fault    *Fault
}

// Read reads a whole document. Content with no newline is one line, and empty content is one
// empty body line.
func Read(content string) Document {
	return read(content, true)
}

// ReadFragment reads a section body or a State body. Frontmatter is a rule of document position,
// so a fragment gets no frontmatter rule, and a leading marker line stays body.
func ReadFragment(content string) Document {
	return read(content, false)
}

// read applies the rules in a fixed order: frontmatter, comments, fences, and H2 headings. The
// comment rule reads the text outside the frontmatter. A comment hides a fence marker, so the
// fence rule reads the comment-free text, and the heading rule reads only body lines.
//
// The fault order is frontmatter, comment, then fence. An unterminated comment removes the rest
// of the document, and with it the closer of any fence that is open, so its fault comes first.
func read(content string, frontmatter bool) Document {
	doc := Document{Lines: splitLines(content)}
	lines := doc.Lines
	held := 0
	if frontmatter {
		held, doc.Fault = frontmatterEnd(lines)
	}
	from := 0
	if held > 0 {
		closer := lines[held-1]
		from = closer.Offset + len(closer.Raw)
	}
	cut, open := cutComments(content, from)
	doc.Comments = cut
	var block fence
	section, inSection := "", false
	for i := range lines {
		line := &lines[i]
		for len(cut) > 0 && cut[0].End <= line.Offset {
			cut = cut[1:]
		}
		line.Text = keptText(content, line.Offset, line.Offset+len(line.Raw), cut)
		if i < held {
			line.Class = Frontmatter
			continue
		}
		block.classify(line)
		if title, ok := strings.CutPrefix(line.Text, H2Prefix); ok && line.Class == Body {
			line.Heading = true
			section, inSection = strings.TrimSpace(title), true
		}
		line.Section, line.InSection = section, inSection
	}
	switch {
	case doc.Fault != nil:
	case open >= 0:
		doc.Fault = &Fault{Kind: CommentFault, Line: lineAt(lines, open)}
	case block.line > 0:
		doc.Fault = &Fault{Kind: FenceFault, Line: block.line}
	}
	return doc
}

// splitLines splits content at each newline, so n newlines give n+1 lines.
func splitLines(content string) []Line {
	lines := make([]Line, 0, strings.Count(content, "\n")+1)
	for start := 0; ; {
		end := strings.IndexByte(content[start:], '\n')
		if end < 0 {
			end = len(content) - start
		}
		raw := strings.TrimSuffix(content[start:start+end], "\r")
		lines = append(lines, Line{Number: len(lines) + 1, Offset: start, Raw: raw})
		start += end + 1
		if start > len(content) {
			return lines
		}
	}
}

// frontmatterEnd returns the number of leading lines that the frontmatter block holds. The block
// opens on a first line that is exactly the marker and closes on the next such line. A block
// with no closer holds no line and gives a fault on line 1.
func frontmatterEnd(lines []Line) (int, *Fault) {
	if lines[0].Raw != frontmatterMarker {
		return 0, nil
	}
	for i := 1; i < len(lines); i++ {
		if lines[i].Raw == frontmatterMarker {
			return i + 1, nil
		}
	}
	return 0, &Fault{Kind: FrontmatterFault, Line: 1}
}

// fence is the open fenced block: the number of its opener line, or 0, and the character and
// the length of its marker run.
type fence struct {
	line int
	char byte
	run  int
}

// classify sets the fence class of line from its comment-free text. A closer uses the opener's
// character and a run at least as long, so a shorter run or the other character stays inside
// the block. Text after the closer's run does not stop the close.
func (f *fence) classify(line *Line) {
	char, run, info := markerRun(line.Text)
	switch {
	case f.line == 0 && run == 0:
	case f.line == 0:
		line.Class, line.Info = FenceOpener, info
		*f = fence{line: line.Number, char: char, run: run}
	case char == f.char && run >= f.run:
		line.Class, line.Info = FenceCloser, info
		*f = fence{}
	default:
		line.Class = Fenced
	}
}

// markerRun returns the fence marker run that starts text after an indent of spaces and tabs,
// and the trimmed text after the run. A run of fewer than three backticks or tildes is no
// marker, and it gives a zero run.
func markerRun(text string) (char byte, run int, info string) {
	rest := strings.TrimLeft(text, " \t")
	if rest == "" || rest[0] != '`' && rest[0] != '~' {
		return 0, 0, ""
	}
	for run < len(rest) && rest[run] == rest[0] {
		run++
	}
	if run < len(FenceMarker) {
		return 0, 0, ""
	}
	return rest[0], run, strings.TrimSpace(rest[run:])
}

// lineAt returns the number of the line that holds the byte at offset.
func lineAt(lines []Line, offset int) int {
	return sort.Search(len(lines), func(i int) bool { return lines[i].Offset > offset })
}
