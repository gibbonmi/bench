package markdown

import (
	"fmt"
	"runtime"
	"slices"
	"strings"
	"testing"
)

// classNames and faultNames spell a class and a fault kind in an expected layout.
var (
	classNames = map[Class]string{Body: "body", Frontmatter: "frontmatter", FenceOpener: "opener", Fenced: "fenced", FenceCloser: "closer"}
	faultNames = map[FaultKind]string{FrontmatterFault: "frontmatter", CommentFault: "comment", FenceFault: "fence"}
)

// layout renders a reading as one row for each line, then one row for each comment span and one
// for the fault. A line row is its number, its byte offset, its class, and its comment-free text
// between bars. Its raw text follows when the two differ, then its info string, an H2 mark, and
// its open section.
func layout(doc Document) []string {
	var rows []string
	for _, line := range doc.Lines {
		row := fmt.Sprintf("%d@%d %s |%s|", line.Number, line.Offset, classNames[line.Class], line.Text)
		if line.Raw != line.Text {
			row += " raw=|" + line.Raw + "|"
		}
		if line.Info != "" {
			row += " info=|" + line.Info + "|"
		}
		if line.Heading {
			row += " h2"
		}
		if line.InSection {
			row += " section=|" + line.Section + "|"
		}
		rows = append(rows, row)
	}
	for _, span := range doc.Comments {
		rows = append(rows, fmt.Sprintf("comment %d-%d", span.Start, span.End))
	}
	if doc.Fault != nil {
		rows = append(rows, fmt.Sprintf("fault %s line %d", faultNames[doc.Fault.Kind], doc.Fault.Line))
	}
	return rows
}

// readingRow is one document and the layout of its reading. A fragment row reads through
// ReadFragment, and every other row reads through Read.
type readingRow struct {
	name     string
	content  string
	fragment bool
	want     []string
}

// checkReadings reads each row's content and compares the whole layout, so a row also fails
// when its rule moves a line, a span, or a fault that the row does not name.
func checkReadings(t *testing.T, rows []readingRow) {
	t.Helper()
	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			read := Read
			if row.fragment {
				read = ReadFragment
			}
			if got := layout(read(row.content)); !slices.Equal(got, row.want) {
				t.Errorf("reading of %q:\n%s\nwant:\n%s", row.content, strings.Join(got, "\n"), strings.Join(row.want, "\n"))
			}
		})
	}
}

// TestReadFenceEdges holds every fence edge of the block reader in one table: the marker
// characters, the run lengths, the indent, the info string, and the unterminated block.
func TestReadFenceEdges(t *testing.T) {
	checkReadings(t, []readingRow{
		{
			name:    "MB1 a tilde block is fenced",
			content: "~~~\nx\n~~~",
			want:    []string{"1@0 opener |~~~|", "2@4 fenced |x|", "3@6 closer |~~~|"},
		},
		{
			name:    "MB2 a shorter run stays inside a longer block",
			content: "````\n```\nx\n````",
			want:    []string{"1@0 opener |````|", "2@5 fenced |```|", "3@9 fenced |x|", "4@11 closer |````|"},
		},
		{
			name:    "MB3 a backtick run stays inside a tilde block",
			content: "~~~\n```\n~~~",
			want:    []string{"1@0 opener |~~~|", "2@4 fenced |```|", "3@8 closer |~~~|"},
		},
		{
			name:    "MB4 an opener after spaces or a tab opens a block",
			content: "  ```\nx\n```\n\t```\ny\n```",
			want:    []string{"1@0 opener |  ```|", "2@6 fenced |x|", "3@8 closer |```|", "4@12 opener |\t```|", "5@17 fenced |y|", "6@19 closer |```|"},
		},
		{
			name:    "MB5 an opener carries its trimmed info string",
			content: "```go  \nx\n```",
			want:    []string{"1@0 opener |```go  | info=|go|", "2@8 fenced |x|", "3@10 closer |```|"},
		},
		{
			name:    "MB6 two backticks open no block",
			content: "``code``\ntext",
			want:    []string{"1@0 body |``code``|", "2@9 body |text|"},
		},
		{
			name:    "MB7 an unterminated block fences each later line",
			content: "```\nx\ny",
			want:    []string{"1@0 opener |```|", "2@4 fenced |x|", "3@6 fenced |y|", "fault fence line 1"},
		},
		{
			name:    "MB73 text after the closer run does not stop the close",
			content: "```\nx\n```go\ny",
			want:    []string{"1@0 opener |```|", "2@4 fenced |x|", "3@6 closer |```go| info=|go|", "4@12 body |y|"},
		},
		{
			name:    "MB78 a no-break space indent opens no block",
			content: "\u00a0```\nx",
			want:    []string{"1@0 body |\u00a0```|", "2@6 body |x|"},
		},
	})
}

// TestReadBlockRules holds the frontmatter, comment, heading, and line rules of the block
// reader, their fault order, and the bound on the comment pass.
func TestReadBlockRules(t *testing.T) {
	checkReadings(t, []readingRow{
		{
			name:    "MB8 frontmatter opens only on line 1",
			content: "---\nindex: <!-- x\n---\n\n---",
			want:    []string{"1@0 frontmatter |---|", "2@4 frontmatter |index: <!-- x|", "3@18 frontmatter |---|", "4@22 body ||", "5@23 body |---|"},
		},
		{
			name:    "MB74 a leading space opens no frontmatter",
			content: " ---\nx\n---",
			want:    []string{"1@0 body | ---|", "2@5 body |x|", "3@7 body |---|"},
		},
		{
			name:    "MB75 a trailing space opens no frontmatter",
			content: "--- \nx\n---",
			want:    []string{"1@0 body |--- |", "2@5 body |x|", "3@7 body |---|"},
		},
		{
			name:    "MB9 unclosed frontmatter faults on line 1 and reads as body",
			content: "---\na\n## B",
			want:    []string{"1@0 body |---|", "2@4 body |a|", "3@6 body |## B| h2 section=|B|", "fault frontmatter line 1"},
		},
		{
			name:    "MB90 a trailing space closes no frontmatter",
			content: "---\na: b\n--- \n<!--",
			want:    []string{"1@0 body |---|", "2@4 body |a: b|", "3@9 body |--- |", "4@14 body || raw=|<!--|", "comment 14-18", "fault frontmatter line 1"},
		},
		{
			name:    "MB10 a comment is cut from its lines and the line numbers stay",
			content: "a <!-- x --> b\nc <!-- d\ne\nf --> g\nh",
			want: []string{
				"1@0 body |a  b| raw=|a <!-- x --> b|",
				"2@15 body |c | raw=|c <!-- d|",
				"3@24 body || raw=|e|",
				"4@26 body | g| raw=|f --> g|",
				"5@34 body |h|",
				"comment 2-12",
				"comment 17-31",
			},
		},
		{
			name:    "MB11 a fence marker inside a comment opens no block",
			content: "<!--\n```\n-->\ntext",
			want:    []string{"1@0 body || raw=|<!--|", "2@5 body || raw=|```|", "3@9 body || raw=|-->|", "4@13 body |text|", "comment 0-12"},
		},
		{
			name:    "MB12 a rejoined comment is removed too",
			content: "<!<!-- z -->-- needle -->",
			want:    []string{"1@0 body || raw=|<!<!-- z -->-- needle -->|", "comment 0-25"},
		},
		{
			name:    "MB71 an open comment inside an open fence reports the comment",
			content: "```\n<!--",
			want:    []string{"1@0 opener |```|", "2@4 fenced || raw=|<!--|", "comment 4-8", "fault comment line 2"},
		},
		{
			name:    "MB88 the first opener runs through the first closer",
			content: "<!-- a <!-- b --> c",
			want:    []string{"1@0 body | c| raw=|<!-- a <!-- b --> c|", "comment 0-17"},
		},
		{
			name:    "MB89 the closer search starts after the opener",
			content: "<!-->",
			want:    []string{"1@0 body || raw=|<!-->|", "comment 0-5", "fault comment line 1"},
		},
		{
			name:     "MB92 a fragment has no frontmatter",
			content:  "---\n```\nx\n---",
			fragment: true,
			want:     []string{"1@0 body |---|", "2@4 opener |```|", "3@8 fenced |x|", "4@10 fenced |---|", "fault fence line 2"},
		},
		{
			name:    "MB13 an unterminated comment removes the rest",
			content: "a\n<!-- b\nc\nd",
			want:    []string{"1@0 body |a|", "2@2 body || raw=|<!-- b|", "3@9 body || raw=|c|", "4@11 body || raw=|d|", "comment 2-12", "fault comment line 2"},
		},
		{
			name:    "MB14 an H2 title is trimmed",
			content: "## Title  ",
			want:    []string{"1@0 body |## Title  | h2 section=|Title|"},
		},
		{
			name:    "MB15 an indented line and other levels are no H2",
			content: "  ## Title\n##Title\n### Title",
			want:    []string{"1@0 body |  ## Title|", "2@11 body |##Title|", "3@19 body |### Title|"},
		},
		{
			name:    "MB16 a fenced heading is no H2",
			content: "```\n## Title\n```",
			want:    []string{"1@0 opener |```|", "2@4 fenced |## Title|", "3@13 closer |```|"},
		},
		{
			name:    "MB17 each line carries its open section",
			content: "x\n## A\ny\n## B\nz",
			want:    []string{"1@0 body |x|", "2@2 body |## A| h2 section=|A|", "3@7 body |y| section=|A|", "4@9 body |## B| h2 section=|B|", "5@14 body |z| section=|B|"},
		},
		{
			name:    "MB18 a carriage return is stripped and offsets stay",
			content: "a\r\n```\r\n",
			want:    []string{"1@0 body |a|", "2@3 opener |```|", "3@8 fenced ||", "fault fence line 2"},
		},
		{
			name:    "MB19 empty content is one empty body line",
			content: "",
			want:    []string{"1@0 body ||"},
		},
		{
			name:    "MB72 a last line without a newline stays",
			content: "a",
			want:    []string{"1@0 body |a|"},
		},
		{
			name:    "MB76 a no-break space trims from a title",
			content: "## Title\u00a0\n## \u00a0Title",
			want:    []string{"1@0 body |## Title\u00a0| h2 section=|Title|", "2@11 body |## \u00a0Title| h2 section=|Title|"},
		},
		{
			name:    "MB77 a zero-width space stays in a title",
			content: "## Title\u200b",
			want:    []string{"1@0 body |## Title\u200b| h2 section=|Title\u200b|"},
		},
	})

	t.Run("MB79 a comment-dense input allocates linearly", func(t *testing.T) {
		const size, comments = 1 << 20, 10000
		comment := "<!-- c -->"
		line := comment + strings.Repeat("x", size/comments-len(comment)-1) + "\n"
		content := strings.Repeat(line, comments)
		content += strings.Repeat("x", size-len(content))
		var before, after runtime.MemStats
		runtime.ReadMemStats(&before)
		doc := Read(content)
		runtime.ReadMemStats(&after)
		if len(doc.Comments) != comments {
			t.Fatalf("Read cut %d comments, want %d", len(doc.Comments), comments)
		}
		if allocated := after.TotalAlloc - before.TotalAlloc; allocated > 4*size {
			t.Errorf("Read of %d bytes allocated %d bytes, want at most %d", size, allocated, 4*size)
		}
	})
}
