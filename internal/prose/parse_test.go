package prose

import (
	"fmt"
	"strings"
	"testing"
)

// words returns n distinct one-word tokens, so a caller builds a sentence of an exact
// length without counting by hand.
func words(n int) string {
	out := make([]string, n)
	for i := range out {
		out[i] = fmt.Sprintf("w%d", i+1)
	}
	return strings.Join(out, " ")
}

// repeatLine returns n copies of one line, each on its own physical line.
func repeatLine(line string, n int) string {
	out := make([]string, n)
	for i := range out {
		out[i] = line
	}
	return strings.Join(out, "\n") + "\n"
}

// TestFindings grades one document per row. Each row is one acceptance row of the
// ticket, so a parser that drops a rule fails the row that names it.
func TestFindings(t *testing.T) {
	for _, tt := range []struct {
		name string
		doc  string
		want []Finding
	}{
		{
			name: "PD9 sentence over the bound",
			doc:  words(26) + ".\n",
			want: []Finding{{Kind: KindSentence, Line: 1, Count: 26}},
		},
		{
			name: "PD10 sentence at the bound",
			doc:  words(25) + ".\n",
		},
		{
			name: "PD11 paragraph over the bound",
			doc:  "One. Two. Three. Four. Five. Six. Seven.\n",
			want: []Finding{{Kind: KindParagraph, Line: 1, Count: 7}},
		},
		{
			name: "PD12 paragraph at the bound",
			doc:  "One. Two. Three. Four. Five. Six.\n",
		},
		{
			name: "PD13 code span is one token",
			doc:  "`" + words(40) + "` is short.\n",
		},
		{
			name: "PD13 fenced block",
			doc:  "```text\n" + words(40) + ".\n```\n",
		},
		{
			name: "PD13 table row",
			doc:  "| " + strings.Join(strings.Fields(words(40)), " | ") + " |\n",
		},
		{
			name: "PD13 heading",
			doc:  "# " + words(30) + "\n",
		},
		{
			name: "PD13 frontmatter value",
			doc:  "---\ndescription: " + words(30) + ".\n---\n\nShort prose.\n",
		},
		{
			name: "PD13 HTML comment",
			doc:  words(20) + ". <!-- " + words(30) + " -->\n",
		},
		{
			name: "PD13 link target",
			doc:  words(22) + " [guide](https://example.com/" + strings.Repeat("a", 40) + " \"one two three four five\") ends here.\n",
		},
		{
			name: "PD14 indented block after a blank line",
			doc:  "Intro prose.\n\n    " + words(30) + ".\n",
		},
		{
			name: "PD14 list continuation after a non-blank line",
			doc:  "- item start\n    " + words(30) + ".\n",
			want: []Finding{{Kind: KindSentence, Line: 1, Count: 32}},
		},
		{
			name: "PD15 list item over the paragraph bound",
			doc:  "- One. Two. Three. Four. Five. Six. Seven.\n",
			want: []Finding{{Kind: KindParagraph, Line: 1, Count: 7}},
		},
		{
			name: "PD15 adjacent list items are separate paragraphs",
			doc:  "- One. Two. Three. Four.\n- Five. Six. Seven. Eight.\n",
		},
		{
			name: "PD15 a bullet marker is not a word",
			doc:  "- " + words(25) + ".\n",
		},
		{
			name: "PD15 an ordered marker is not a sentence",
			doc:  "3. One. Two. Three. Four. Five. Six.\n",
		},
		{
			name: "PD15 ordered marker is not a word",
			doc:  "3. " + words(25) + ".\n",
		},
		{
			name: "PD16 in-token periods and abbreviations do not split",
			doc:  "See spec.md and 1.2 and e.g. and i.e. and etc. and vs. and cf. and one…more " + words(10) + ".\n",
			want: []Finding{{Kind: KindSentence, Line: 1, Count: 26}},
		},
		{
			name: "PD16 a terminator splits",
			doc:  words(26) + ". " + words(26) + ".\n",
			want: []Finding{{Kind: KindSentence, Line: 1, Count: 26}, {Kind: KindSentence, Line: 1, Count: 26}},
		},
		{
			name: "PD17 a no-break space splits a word",
			doc:  words(24) + " last word.\n",
			want: []Finding{{Kind: KindSentence, Line: 1, Count: 26}},
		},
		{
			name: "PD17 a zero-width space does not split a word",
			doc:  words(24) + " last​word.\n",
		},
		{
			name: "PD21 unterminated fence",
			doc:  "Intro prose.\n\n```text\ncode\n",
			want: []Finding{{Kind: KindFence, Line: 3}},
		},
		{
			name: "PD21 unterminated HTML comment",
			doc:  "Intro prose.\n\n<!-- open\nmore\n",
			want: []Finding{{Kind: KindComment, Line: 3}},
		},
		{
			name: "PD21 unterminated frontmatter",
			doc:  "---\ndescription: open\nbody\n",
			want: []Finding{{Kind: KindFrontmatter, Line: 1}},
		},
		{
			name: "PD41 named field lines are not graded",
			doc:  "Writes: internal/prose/ (new)\nBlocked by: none\nNext: run the gate\nFeeds: ticket 02\n",
		},
		{
			name: "PD41 a four-word label line with no terminator is not graded",
			doc:  "One two three four: " + words(30) + "\n",
		},
		{
			name: "PD41 a long prefix before a colon is prose",
			doc:  words(6) + ": " + words(20) + "\n",
			want: []Finding{{Kind: KindSentence, Line: 1, Count: 26}},
		},
		{
			name: "PD41 an Occurrence line with a terminator is graded",
			doc:  "Occurrence: " + words(30) + ".\n",
			want: []Finding{{Kind: KindSentence, Line: 1, Count: 31}},
		},
		{
			name: "PD41 a five-word label line is prose",
			doc:  "One two three four five: " + words(30) + ".\n",
			want: []Finding{{Kind: KindSentence, Line: 1, Count: 35}},
		},
		{
			name: "PD41 twenty label lines are twenty paragraphs",
			doc:  repeatLine("Occurrence: FT100 names this spec.", 20),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := Findings(tt.doc)
			if len(got) != len(tt.want) {
				t.Fatalf("Findings() = %+v, want %+v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("finding %d = %+v, want %+v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
