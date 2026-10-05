package roadmap

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/roadmap/roadmaptest"
)

// indexTree is the index-only tree a test drives when the row files are not the subject.
// It holds parsed index bytes and no roadmap/ directory at all. Absence is the state a
// missing directory classifies to, and the state the parse must read here. The zero value
// is no state at all and would grade as a degraded directory.
func indexTree(index string) Tree {
	state := bounds.StateParsed
	if index == "" {
		state = bounds.StateEmpty
	}
	return Tree{Index: bounds.Classified{State: state, Data: []byte(index)}, DirState: bounds.StateAbsent}
}

// splitTree is the split board a test drives. It holds parsed index bytes plus one parsed
// row file per named basename, in the directory order os.ReadDir would report.
func splitTree(index string, files map[string]string) Tree {
	tree := indexTree(index)
	if files == nil {
		return tree
	}
	tree.DirState = bounds.StateParsed
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		state := bounds.StateParsed
		if files[name] == "" {
			state = bounds.StateEmpty
		}
		tree.Files = append(tree.Files, RowFile{Name: name, State: state, Data: []byte(files[name])})
		tree.DirBytes += len(files[name])
	}
	return tree
}

// Row is one split-board row, a heading line and its body. board and writeBoard each take
// one Row per row.
type Row struct {
	Heading string
	Body    string
}

// writeBoard writes the on-disk split board built from one heading-and-body Row per row.
// It writes the index text as ROADMAP.md and one roadmap/<name> row file per row. It is
// the on-disk twin of splitTree, for the tests that drive a command rather than the
// parse.
func writeBoard(t *testing.T, root string, rows ...Row) {
	t.Helper()
	index, files := board(rows...)
	roadmaptest.WriteSplitBoard(t, root, index, files)
}

// board renders an index and its row files from rows, the shape most fixtures need. That
// shape is `**FT1 — one.**` in ROADMAP.md and a heading plus body in roadmap/FT1.md.
func board(rows ...Row) (string, map[string]string) {
	var index strings.Builder
	files := map[string]string{}
	for _, row := range rows {
		index.WriteString(row.Heading + "\n\n")
		space := strings.Index(row.Heading, " ")
		if space < 0 {
			panic(fmt.Sprintf("board: heading %q has no space after the row ID", row.Heading))
		}
		id := row.Heading[2:space]
		content := row.Heading + "\n"
		if row.Body != "" {
			content += row.Body
		}
		files[id+".md"] = content
	}
	return index.String(), files
}

// rowNextTree is the one-row split board the marker tests drive. It holds FT1's index
// line under section, and a detail file whose body is the given text.
func rowNextTree(section, body string) Tree {
	const heading = "**FT1 — one.**"
	index := heading + "\n"
	if section != "" {
		index = section + "\n\n" + index
	}
	return splitTree(index, map[string]string{"FT1.md": heading + "\n" + body})
}

// Closing A's row removes its heading and every sequence entry for A, and keeps B's row
// and B's sequence entry byte for byte.
func TestCommitmentSequenceClosure(t *testing.T) {
	const index = "# Roadmap\n\n## Parked\n\n**FT1 — A**\n\n**FT2 — B**\n\n## Recommended sequence\n\n1. A\n2. B\n"
	got, err := Close([]byte(index), []string{"FT1"}, []string{"B"})
	const want = "# Roadmap\n\n## Parked\n\n**FT2 — B**\n\n## Recommended sequence\n\n1. B\n"
	if err != nil || string(got) != want {
		t.Fatalf("closure = %q, %v; want %q", got, err, want)
	}
	document, _, diagnostics := ParseDocument(splitTree(string(got), map[string]string{"FT2.md": "**FT2 — B**\n"}), nil, true)
	if len(diagnostics) != 0 || len(document.Rows) != 1 || document.Rows[0].ID != "FT2" || len(document.Sequence) != 1 || document.Sequence[0].Text != "B" {
		t.Fatalf("closed board = %+v, %v", document, diagnostics)
	}
	for _, rows := range [][]string{{"FT3"}, {"FT1", "FT1"}} {
		if _, err := Close([]byte(index+"**FT1 — A again**\n"), rows, nil); err == nil {
			t.Fatalf("closure of %v accepted an absent or repeated row", rows)
		}
	}
}

// Closing FT1 unblocks its dependents in both dependency tables. A row that names FT1 as
// its dependent goes, FT1 leaves each dependency list, and a row whose list empties goes.
// FT10 shares FT1's prefix and stays, and the table headings stay byte for byte.
func TestCommitmentDependencyClosure(t *testing.T) {
	const tables = "## Dependencies\n\nThe dependent FT is named first.\n\n### Literal\n\n| FT | Depends on | Why |\n|---|---|---|\n"
	const index = "# Roadmap\n\n**FT1 — A**\n\n**FT2 — B**\n\n" + tables +
		"| FT2 | FT1 | B needs A. |\n| FT3 | FT10 | C needs J. |\n\n### Recommended\n\n| FT | Better specified after | Why |\n|---|---|---|\n" +
		"| FT4 | FT5, FT1, FT6 | D follows. |\n| FT1 | FT7 | A follows G. |\n\n## Recommended sequence\n\n1. A\n2. B\n"
	got, err := Close([]byte(index), []string{"FT1"}, []string{"B"})
	const want = "# Roadmap\n\n**FT2 — B**\n\n" + tables +
		"| FT3 | FT10 | C needs J. |\n\n### Recommended\n\n| FT | Better specified after | Why |\n|---|---|---|\n" +
		"| FT4 | FT5, FT6 | D follows. |\n\n## Recommended sequence\n\n1. B\n"
	if err != nil || string(got) != want {
		t.Fatalf("dependency closure = %q, %v; want %q", got, err, want)
	}
}

func TestCommitmentOccurrenceUpdate(t *testing.T) {
	const original = "**FT1 — Preserve the obligation**\n\nThe requirement stays unchanged.\n"
	got, err := RequirementBytes("roadmap/FT1.md", []byte(original+"Occurrences: first, second\n"))
	if err != nil || string(got) != original {
		t.Fatalf("occurrence requirement = %q, %v", got, err)
	}
	for _, value := range []string{"Occurrences:", "Occurrences: second, first", "Occurrences: changed requirement prose"} {
		if _, err := RequirementBytes("roadmap/FT1.md", []byte(original+value+"\n")); err == nil {
			t.Fatalf("accepted malformed event %q", value)
		}
	}
	got, err = RequirementBytes("roadmap/FT1.md", []byte(original+"Event: a replacement requirement\n"))
	if err != nil || string(got) == original {
		t.Fatalf("unknown label erased requirement: %q %v", got, err)
	}
}
