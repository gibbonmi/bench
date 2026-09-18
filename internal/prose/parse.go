package prose

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// codeSpanToken opens the placeholder that replaces every inline code span before the
// label test and before the split into sentences, so a colon or a period inside a span
// never ends a field line or a sentence, and a long span counts as one word. The whole
// placeholder is the prefix, the span's index in the fold, and a closing NUL, so a
// diagnostic restores the span as the author wrote it. The NUL bytes keep the placeholder
// outside any authored token, and the letters keep it a word.
const codeSpanToken = "\x00code"

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

// templateFields are the label prefixes of the ticket header, the roadmap ledger, and the
// decision map's evidence ledger, in lowercase and with one space between words. Such a
// field keeps its own paragraph even when it ends with a period, because the template
// repeats the field and the repetition is not one deep paragraph. The list is closed: the
// rule file names these fields, and a new template field needs a rule edit and an edit
// here. TestTemplateFieldNamesMatchTheRuleFile holds the two halves together.
var templateFields = map[string]bool{
	"blocked by":  true,
	"covers":      true,
	"drift":       true,
	"occurrence":  true,
	"occurrences": true,
	"source":      true,
	"sources":     true,
	"supports":    true,
	"writes":      true,
}

// foldCodeSpans replaces every inline code span with one token and returns the spans it
// folded, in order. A span opens at a run of backticks and closes at the next run of the
// same length, so a shorter run inside the span stays part of the span. An unclosed run is
// literal text.
func foldCodeSpans(content string) (string, []string) {
	var out strings.Builder
	var spans []string
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
		spans = append(spans, content[i:end+open])
		out.WriteString(codeSpanToken + strconv.Itoa(len(spans)-1) + "\x00")
		i = end + open
	}
	return out.String(), spans
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
// the token, because a sentence reports the line of its first token. raw is the token as
// the author wrote it, and it is empty for a token that folded no code span.
type token struct {
	text string
	line int
	raw  string
}

// Findings grades one document and returns every fault in document order. An
// unterminated frontmatter block, HTML comment, or fenced code block ends the grade at
// once and returns one finding: past that delimiter the parser cannot tell prose from
// code, and a truncated grade reports a clean file.
func Findings(doc string) []Finding {
	lines, fault := prepare(doc)
	if fault != nil {
		return []Finding{*fault}
	}
	return gradeBlocks(lines)
}

// walkParagraphs splits the remaining lines into paragraphs and calls visit with each
// paragraph's first physical line and its tokens. A blank line, a skipped line, a list
// marker, and a field line each start a new paragraph. The grade and the exported
// projection both walk here, so one paragraph rule serves both.
func walkParagraphs(lines []string, visit func(start int, toks []token)) {
	var current []token
	start := 0

	flush := func() {
		if len(current) > 0 {
			visit(start, current)
		}
		current, start = nil, 0
	}
	add := func(content string, line int, spans []string) {
		toks := tokenize(content, line, spans)
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
		content, spans := foldCodeSpans(trimmed)
		if m := listMarkerPattern.FindString(content); m != "" {
			flush()
			content = content[len(m):]
		}
		terminated := hasTerminator(content, number)
		if isFieldLine(content, terminated) {
			flush()
			// A field line is its own paragraph, so a run of such lines never forms one deep
			// paragraph. A field line with no terminator holds no sentence to grade.
			if terminated {
				add(content, number, spans)
				flush()
			}
			continue
		}
		add(content, number, spans)
	}
	flush()
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

// isFieldLine reports whether the line is a template field rather than prose. The line
// must first be label-shaped: a prefix of one to four words that ends at the first
// colon, which is why a URL scheme is not a label. Such a line is a field when its
// prefix is a template field name, or when the whole line carries no sentence
// terminator. A label-shaped prose line that carries a terminator is therefore prose,
// and it stays in its paragraph for the paragraph bound to count.
func isFieldLine(content string, terminated bool) bool {
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
	if !terminated {
		return true
	}
	return templateFields[strings.ToLower(strings.Join(words, " "))]
}

// hasTerminator reports whether the line carries a sentence terminator.
func hasTerminator(content string, line int) bool {
	for _, t := range tokenize(content, line, nil) {
		if isBoundaryToken(t.text) {
			return true
		}
	}
	return false
}

// tokenize turns one folded line of prose into tokens. It keeps link text, drops link
// targets, and removes the emphasis marks around a word. spans are the code spans the fold
// of that line returned, so each token can restore the span its placeholder stands for.
func tokenize(content string, line int, spans []string) []token {
	s := inlineLinkPattern.ReplaceAllString(content, "$1")
	s = refLinkPattern.ReplaceAllString(s, "$1")
	var out []token
	for _, field := range strings.FieldsFunc(s, unicode.IsSpace) {
		field = strings.Trim(field, emphasisMarks)
		if field == "" {
			continue
		}
		out = append(out, token{text: field, line: line, raw: unfoldCodeSpans(field, spans)})
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
