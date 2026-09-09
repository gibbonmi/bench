package prose

import (
	"strconv"
	"strings"
)

// maxStartWords is how many words a sentence start carries. Three words identify the
// sentence inside its paragraph without quoting the document.
const maxStartWords = 3

// gradeParagraph splits one paragraph into sentences, grades both bounds, and records the
// start of every sentence. A paragraph fault carries the starts, so a reader finds the
// sentences of a deep paragraph without a second pass over the file.
func gradeParagraph(start int, toks []token) []Finding {
	var out []Finding
	var starts []SentenceStart
	sentences := 0
	words, first, from := 0, start, 0
	closeSentence := func(end int) {
		if words > MaxSentenceWords {
			out = append(out, Finding{Kind: KindSentence, Line: first, Count: words})
		}
		starts = append(starts, SentenceStart{Line: first, Text: sentenceStart(toks[from:end])})
		sentences++
		words, from = 0, end
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
				closeSentence(i + 1)
			}
		}
	}
	if sentences > MaxParagraphSentences {
		out = append(out, Finding{Kind: KindParagraph, Line: start, Count: sentences, Starts: starts})
	}
	return out
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
		word := t.raw
		if word == "" {
			word = t.text
		}
		out = append(out, word)
		if len(out) == maxStartWords {
			break
		}
	}
	return strings.Join(out, " ")
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
