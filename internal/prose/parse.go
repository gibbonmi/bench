package prose

import (
	"regexp"
	"strings"
	"unicode"
)

// codeSpanToken replaces every inline code span before the label test and before the
// split into sentences, so a colon or a period inside a span never ends a field line or
// a sentence, and a long span counts as one word. The NUL bytes keep the placeholder
// outside any authored token.
const codeSpanToken = "\x00code\x00"

var (
	inlineLinkPattern  = regexp.MustCompile(`!?\[([^\]]*)\]\([^)]*\)`)
	refLinkPattern     = regexp.MustCompile(`!?\[([^\]]*)\]\[[^\]]*\]`)
	linkDefPattern     = regexp.MustCompile(`^\[[^\]]+\]:`)
	listMarkerPattern  = regexp.MustCompile(`^([-*+]|[0-9]+[.)])\s+`)
	emphasisMarks      = "*_~"
	sentenceCloseMarks = "\"'’”)]}»›"
)

// abbreviations are the tokens that end with a period and never end a sentence. The
// list is closed: the rule file names these five, and a sixth needs a rule edit.
var abbreviations = map[string]bool{
	"e.g.": true,
	"i.e.": true,
	"etc.": true,
	"vs.":  true,
	"cf.":  true,
}

// foldCodeSpans replaces every inline code span with one token. A span opens at a run
// of backticks and closes at the next run of the same length, so a shorter run inside
// the span stays part of the span. An unclosed run is literal text.
func foldCodeSpans(content string) string {
	var out strings.Builder
	for i := 0; i < len(content); {
		if content[i] != '`' {
			out.WriteByte(content[i])
			i++
			continue
		}
		open := backtickRun(content, i)
		end := -1
		for j := i + open; j < len(content); {
			if content[j] != '`' {
				j++
				continue
			}
			run := backtickRun(content, j)
			if run == open {
				end = j
				break
			}
			j += run
		}
		if end < 0 {
			out.WriteString(content[i : i+open])
			i += open
			continue
		}
		out.WriteString(codeSpanToken)
		i = end + open
	}
	return out.String()
}

// backtickRun returns the length of the run of backticks that starts at index i.
func backtickRun(content string, i int) int {
	n := 0
	for i+n < len(content) && content[i+n] == '`' {
		n++
	}
	return n
}

// token is one word candidate and the physical line it came from. The line travels with
// the token, because a sentence reports the line of its first token.
type token struct {
	text string
	line int
}

// Findings grades one document and returns every fault in document order. An
// unterminated frontmatter block, HTML comment, or fenced code block ends the grade at
// once and returns one finding: past that delimiter the parser cannot tell prose from
// code, and a truncated grade reports a clean file.
func Findings(doc string) []Finding {
	lines := strings.Split(doc, "\n")
	if f := stripFrontmatter(lines); f != nil {
		return []Finding{*f}
	}
	if f := stripComments(lines); f != nil {
		return []Finding{*f}
	}
	if f := stripFences(lines); f != nil {
		return []Finding{*f}
	}
	return gradeBlocks(lines)
}

// stripFrontmatter blanks a leading YAML block. Frontmatter values are trigger text for
// a loader, not prose for a reader.
func stripFrontmatter(lines []string) *Finding {
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return nil
	}
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			for j := 0; j <= i; j++ {
				lines[j] = ""
			}
			return nil
		}
	}
	return &Finding{Kind: KindFrontmatter, Line: 1}
}

// stripComments removes every HTML comment span in place and keeps the line count, so
// each later finding still names the physical line of the document.
func stripComments(lines []string) *Finding {
	open, openLine := false, 0
	for i, line := range lines {
		var kept strings.Builder
		rest := line
		for {
			if open {
				end := strings.Index(rest, "-->")
				if end < 0 {
					rest = ""
					break
				}
				open, rest = false, rest[end+3:]
				continue
			}
			start := strings.Index(rest, "<!--")
			if start < 0 {
				kept.WriteString(rest)
				break
			}
			kept.WriteString(rest[:start])
			open, openLine, rest = true, i+1, rest[start+4:]
		}
		lines[i] = kept.String()
	}
	if open {
		return &Finding{Kind: KindComment, Line: openLine}
	}
	return nil
}

// stripFences blanks every fenced code block and both of its markers. A closing fence
// uses the same marker character and is at least as long as the opening one.
func stripFences(lines []string) *Finding {
	openLine, marker := 0, ""
	for i, line := range lines {
		found := fenceMarker(strings.TrimLeft(line, " \t"))
		if openLine == 0 {
			if found != "" {
				openLine, marker, lines[i] = i+1, found, ""
			}
			continue
		}
		lines[i] = ""
		if found != "" && found[0] == marker[0] && len(found) >= len(marker) {
			openLine, marker = 0, ""
		}
	}
	if openLine != 0 {
		return &Finding{Kind: KindFence, Line: openLine}
	}
	return nil
}

// fenceMarker returns the leading run of backticks or tildes when the run is a fence
// marker, and an empty string otherwise.
func fenceMarker(trimmed string) string {
	if len(trimmed) < 3 {
		return ""
	}
	c := trimmed[0]
	if c != '`' && c != '~' {
		return ""
	}
	n := 0
	for n < len(trimmed) && trimmed[n] == c {
		n++
	}
	if n < 3 {
		return ""
	}
	return trimmed[:n]
}

// gradeBlocks splits the remaining lines into paragraphs and grades each one. A blank
// line, a skipped line, a list marker, and a label line each start a new paragraph.
func gradeBlocks(lines []string) []Finding {
	var out []Finding
	var current []token
	start := 0

	flush := func() {
		if len(current) > 0 {
			out = append(out, gradeParagraph(start, current)...)
		}
		current, start = nil, 0
	}
	add := func(content string, line int) {
		toks := tokenize(content, line)
		if len(toks) == 0 {
			return
		}
		if len(current) == 0 {
			start = line
		}
		current = append(current, toks...)
	}

	afterBlank, inSkipBlock := true, false
	for i, line := range lines {
		number := i + 1
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			flush()
			afterBlank, inSkipBlock = true, false
			continue
		}
		if inSkipBlock {
			continue
		}
		// An indented block that opens after a blank line is code. The same indent after a
		// non-blank line continues a list item, and that text is prose.
		if afterBlank && indentWidth(line) >= 4 {
			flush()
			inSkipBlock = true
			continue
		}
		if afterBlank && strings.HasPrefix(trimmed, "<") {
			flush()
			inSkipBlock = true
			continue
		}
		afterBlank = false
		if isSkippedLine(trimmed) {
			flush()
			continue
		}
		// The fold comes before the label test, so a colon inside a code span never makes a
		// field line. Only a field line is skipped.
		content := foldCodeSpans(trimmed)
		if m := listMarkerPattern.FindString(content); m != "" {
			flush()
			content = content[len(m):]
		}
		if isLabelLine(content) {
			flush()
			// A label line is its own paragraph, so a run of such lines never forms one deep
			// paragraph. A label line with no sentence terminator is a template field, not prose.
			if hasTerminator(content, number) {
				add(content, number)
				flush()
			}
			continue
		}
		add(content, number)
	}
	flush()
	return out
}

// isSkippedLine reports the line classes that carry no prose: a heading, a table row, a
// thematic break, and a link reference definition.
func isSkippedLine(trimmed string) bool {
	switch {
	case strings.HasPrefix(trimmed, "#"):
		return true
	case strings.HasPrefix(trimmed, "|"):
		return true
	case isThematicBreak(trimmed):
		return true
	case linkDefPattern.MatchString(trimmed):
		return true
	}
	return false
}

// isThematicBreak reports whether a line is a horizontal rule: three or more of one
// mark, with spaces allowed between them and nothing else on the line.
func isThematicBreak(trimmed string) bool {
	if trimmed == "" {
		return false
	}
	mark := rune(trimmed[0])
	if mark != '-' && mark != '_' && mark != '*' {
		return false
	}
	count := 0
	for _, r := range trimmed {
		switch {
		case r == mark:
			count++
		case r == ' ' || r == '\t':
		default:
			return false
		}
	}
	return count >= 3
}

// indentWidth counts the leading whitespace of a line, with one tab as four columns.
func indentWidth(line string) int {
	width := 0
	for _, r := range line {
		switch r {
		case ' ':
			width++
		case '\t':
			width += 4
		default:
			return width
		}
	}
	return width
}

// isLabelLine reports whether the line opens with a label of at most four words and a
// colon. The colon must end the label, so a URL scheme is not a label.
func isLabelLine(content string) bool {
	colon := strings.IndexByte(content, ':')
	if colon <= 0 {
		return false
	}
	if colon+1 < len(content) && content[colon+1] != ' ' && content[colon+1] != '\t' {
		return false
	}
	words := strings.Fields(content[:colon])
	if len(words) == 0 || len(words) > 4 {
		return false
	}
	for _, w := range words {
		if isBoundaryToken(w) {
			return false
		}
	}
	return true
}

// hasTerminator reports whether the line carries a sentence terminator.
func hasTerminator(content string, line int) bool {
	for _, t := range tokenize(content, line) {
		if isBoundaryToken(t.text) {
			return true
		}
	}
	return false
}

// tokenize turns one line of prose into tokens. It folds each code span into one token,
// keeps link text, drops link targets, and removes the emphasis marks around a word.
func tokenize(content string, line int) []token {
	s := foldCodeSpans(content)
	s = inlineLinkPattern.ReplaceAllString(s, "$1")
	s = refLinkPattern.ReplaceAllString(s, "$1")
	var out []token
	for _, field := range strings.FieldsFunc(s, unicode.IsSpace) {
		field = strings.Trim(field, emphasisMarks)
		if field == "" {
			continue
		}
		out = append(out, token{text: field, line: line})
	}
	return out
}

// gradeParagraph splits one paragraph into sentences and grades both bounds.
func gradeParagraph(start int, toks []token) []Finding {
	var out []Finding
	sentences := 0
	words, first := 0, start
	closeSentence := func() {
		if words > MaxSentenceWords {
			out = append(out, Finding{Kind: KindSentence, Line: first, Count: words})
		}
		sentences++
		words = 0
	}
	for i, t := range toks {
		if words == 0 {
			first = t.line
		}
		if isWord(t.text) {
			words++
		}
		if isBoundaryToken(t.text) || i == len(toks)-1 {
			if words > 0 {
				closeSentence()
			}
		}
	}
	if sentences > MaxParagraphSentences {
		out = append(out, Finding{Kind: KindParagraph, Line: start, Count: sentences})
	}
	return out
}

// isWord reports whether a token counts toward the sentence bound. A token counts when
// it holds a letter or a digit, so bare punctuation is not a word.
func isWord(tok string) bool {
	for _, r := range tok {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

// isBoundaryToken reports whether a token ends a sentence. The test looks at the end of
// the token, so a period inside a token never splits. A closing quote or bracket after
// the terminator is ignored, and the five listed abbreviations never end a sentence.
func isBoundaryToken(tok string) bool {
	t := strings.TrimRight(tok, sentenceCloseMarks)
	if t == "" {
		return false
	}
	if abbreviations[strings.ToLower(t)] {
		return false
	}
	return strings.HasSuffix(t, ".") || strings.HasSuffix(t, "!") || strings.HasSuffix(t, "?") || strings.HasSuffix(t, "…")
}
