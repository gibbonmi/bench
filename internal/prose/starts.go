package prose

import (
	"strconv"
	"strings"
)

// maxStartWords is how many words a sentence start carries. Three words identify the
// sentence inside its paragraph without quoting the document.
const maxStartWords = 3

// Paragraphs returns the sentences of every paragraph of doc, in document order. Each
// sentence is its words joined by one space, and every inline code span reads as the author
// wrote it, so a wrap inside a sentence and a run of spaces both arrive collapsed. A caller
// that pins a document against another source reads this projection, and the pin then shares
// the grade's paragraph rule and sentence rule. A document whose frontmatter block, HTML
// comment, or fenced block never closes returns nothing, because past that delimiter the
// parser cannot tell prose from code.
func Paragraphs(doc string) [][]string {
	lines, fault := prepare(doc)
	if fault != nil {
		return nil
	}
	var out [][]string
	walkParagraphs(lines, func(_ int, toks []token) {
		var sentences []string
		for _, span := range sentenceSpans(toks) {
			sentences = append(sentences, sentenceText(span))
		}
		out = append(out, sentences)
	})
	return out
}

// gradeBlocks grades every paragraph of the remaining lines.
func gradeBlocks(lines []string) []Finding {
	var out []Finding
	walkParagraphs(lines, func(start int, toks []token) {
		out = append(out, gradeParagraph(start, toks)...)
	})
	return out
}

// sentenceSpans splits one paragraph's tokens into its sentences, in document order. A
// sentence closes at a boundary token, and the paragraph's last token closes the sentence it
// is in. A run of tokens with no word closes nothing, so bare punctuation is never a
// sentence. The grade and Paragraphs both split here, so a sentence means one thing.
func sentenceSpans(toks []token) [][]token {
	var out [][]token
	words, from := 0, 0
	for i, t := range toks {
		if isWord(t.text) {
			words++
		}
		if !isBoundaryToken(t.text) && i != len(toks)-1 {
			continue
		}
		if words > 0 {
			out = append(out, toks[from:i+1])
			words, from = 0, i+1
		}
	}
	return out
}

// gradeParagraph grades both bounds over one paragraph and records the start of every
// sentence. A paragraph fault carries the starts, so a reader finds the sentences of a deep
// paragraph without a second pass over the file.
func gradeParagraph(start int, toks []token) []Finding {
	var out []Finding
	var starts []SentenceStart
	spans := sentenceSpans(toks)
	for _, span := range spans {
		words, line := 0, span[0].line
		for _, t := range span {
			if !isWord(t.text) {
				continue
			}
			if words == 0 {
				line = t.line
			}
			words++
		}
		if words > MaxSentenceWords {
			out = append(out, Finding{Kind: KindSentence, Line: line, Count: words})
		}
		starts = append(starts, SentenceStart{Line: line, Text: sentenceStart(span)})
	}
	if len(spans) > MaxParagraphSentences {
		out = append(out, Finding{Kind: KindParagraph, Line: start, Count: len(spans), Starts: starts})
	}
	return out
}

// sentenceText joins one sentence's tokens with one space. Every whitespace run of the
// document therefore reads as one space, which is the collapse a pinning needle carries.
func sentenceText(toks []token) string {
	out := make([]string, 0, len(toks))
	for _, t := range toks {
		out = append(out, t.word())
	}
	return strings.Join(out, " ")
}

// sentenceStart joins the first words of one sentence, as the author wrote them. It counts
// the same words the sentence bound counts, so bare punctuation never fills a start, and a
// sentence of fewer words gives the words it has.
func sentenceStart(toks []token) string {
	out := make([]string, 0, maxStartWords)
	for _, t := range toks {
		if !isWord(t.text) {
			continue
		}
		out = append(out, t.word())
		if len(out) == maxStartWords {
			break
		}
	}
	return strings.Join(out, " ")
}

// word answers one token as the author wrote it, with every code span it folded restored.
func (t token) word() string {
	if t.raw != "" {
		return t.raw
	}
	return t.text
}

// unfoldCodeSpans restores the code spans one token folded, so the token reads as the
// author wrote it. It returns an empty string for a token that folded no span, and the
// caller then keeps the token text itself.
func unfoldCodeSpans(field string, spans []string) string {
	if !strings.Contains(field, codeSpanToken) {
		return ""
	}
	var out strings.Builder
	rest := field
	for {
		at := strings.Index(rest, codeSpanToken)
		if at < 0 {
			out.WriteString(rest)
			return out.String()
		}
		out.WriteString(rest[:at])
		rest = rest[at+len(codeSpanToken):]
		end := strings.IndexByte(rest, 0)
		if end < 0 {
			out.WriteString(rest)
			return out.String()
		}
		index, err := strconv.Atoi(rest[:end])
		if err == nil && index < len(spans) {
			out.WriteString(spans[index])
		}
		rest = rest[end+1:]
	}
}
