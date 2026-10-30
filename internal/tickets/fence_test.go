// Tests for the fence fault that ParseTicket reads from the block reader.
package tickets

import (
	"strings"
	"testing"
)

// TestParseTicketReadsTheBlockReaderFault grades each fence shape by the block reader's
// fault. A tilde block refuses like a backtick block, a shorter run stays inside a longer
// block, and an unclosed frontmatter marker on the first line changes no diagnostic.
func TestParseTicketReadsTheBlockReaderFault(t *testing.T) {
	const body = "One parser owns the schema."
	fourRun := "````\nquoted\n```\nCovers: TG9\n````"
	cases := []struct {
		name, document string
		want, covers   []string
	}{
		{
			name:     "MB35 an unterminated tilde block",
			document: wellFormed() + "~~~\nquoted\n",
			want:     []string{"unterminated fence"},
			covers:   []string{"TG1", "TG2"},
		},
		{
			name:     "MB36 a four-backtick block that holds a three-backtick line",
			document: strings.Replace(wellFormed(), body, strings.Replace(fourRun, "Covers: TG9\n", "", 1), 1),
			covers:   []string{"TG1", "TG2"},
		},
		{
			name:     "MB37 a Covers line inside the four-backtick block",
			document: strings.Replace(strings.Replace(wellFormed(), "Covers: TG1, TG2\n", "", 1), body, fourRun, 1),
			want:     []string{"missing Covers"},
		},
		{
			name:     "MB95 an unclosed frontmatter marker on the first line",
			document: "---\n" + wellFormed(),
			covers:   []string{"TG1", "TG2"},
		},
		{
			name:     "MB95 an unclosed frontmatter marker and an unterminated block",
			document: "---\n" + wellFormed() + "```\nquoted\n",
			want:     []string{"unterminated fence"},
			covers:   []string{"TG1", "TG2"},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			ticket, got := ParseTicket("parser.md", []byte(testCase.document), []string{"parser.md"}, "TG")
			assertDiagnostics(t, got, testCase.want)
			if strings.Join(ticket.Covers, ",") != strings.Join(testCase.covers, ",") {
				t.Fatalf("Covers = %#v, want %#v", ticket.Covers, testCase.covers)
			}
		})
	}
}
