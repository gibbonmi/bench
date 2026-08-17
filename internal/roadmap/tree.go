package roadmap

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/spec"
)

// RoadmapDir is RoadmapFile's sibling: the directory holding one detail owner per
// index row, `roadmap/<ID>.md`, whose first line repeats that row's index line.
const RoadmapDir = "roadmap"

// RowFile is one classified entry of the RoadmapDir listing — the basename as the
// directory reported it, the state the classifier graded it, and its bytes when the
// read completed.
type RowFile struct {
	Name, Reason string
	State        bounds.FileState
	Data         []byte
}

// Tree is the classified split board: RoadmapFile's read and the RoadmapDir listing.
// It is the one shape the parse consumes, so a test drives the parser with in-memory
// bytes and every command reaches the same parse through LoadTree.
type Tree struct {
	Index     bounds.Classified
	DirState  bounds.FileState
	DirReason string
	DirBytes  int
	Files     []RowFile
}

// LoadTree is the one filesystem read of the split board. Every roadmap surface —
// the board command, the context snapshot, the owner check, status, the dashboard,
// and the conformance check — reaches the parse through it, so no caller re-derives
// where a row's detail lives.
func LoadTree(root string) Tree {
	tree := Tree{Index: bounds.Classify(filepath.Join(root, RoadmapFile), bounds.ControlRecordLimit)}
	dir := bounds.ClassifyDir(filepath.Join(root, RoadmapDir))
	tree.DirState, tree.DirReason, tree.DirBytes = dir.State, dir.Reason, dirBytes(dir.Entries)
	for _, entry := range dir.Entries {
		// The producer form, not Classify: a row file is authoritative input the board
		// is graded from, so a link here would grade bytes nobody put under roadmap/.
		// Both link kinds land in the wrong-type state the listing pass already reports.
		c := bounds.ClassifyNoFollow(filepath.Join(root, RoadmapDir, entry.Name()))
		tree.Files = append(tree.Files, RowFile{Name: entry.Name(), Reason: c.Reason, State: c.State, Data: c.Data})
	}
	return tree
}

// rowFilePath is the repo-relative detail owner of one row ID. Diagnostics, the
// migration, and every reader that names a row's file go through it.
func rowFilePath(id string) string { return RoadmapDir + "/" + id + ".md" }

// Diagnostic is one integrity fault ParseDocument found in the split board, carried
// as its own path and reason rather than a formatted string a caller would have to
// re-parse. A legal basename may itself contain ": " (the string String returns), so
// a reader that wants the path or the reason separately takes the field, never a cut
// of the rendered text.
type Diagnostic struct {
	Path, Reason string
}

// String renders a Diagnostic in the one format every roadmap surface has always
// shown a reader: the path, then the reason. ValidateRoadmapTree's conformance
// binding and every diagnostic-string test read this text.
func (d Diagnostic) String() string { return d.Path + ": " + d.Reason }

// ParseDocument projects a classified tree into the Document every roadmap surface
// renders, the parse failures the snapshot reports, and the ordered integrity
// diagnostics the conformance check returns. A degraded RoadmapDir comes first because
// it is the precondition every row-level finding rests on; the rest run in index order
// and then in directory order, and each begins with the repo-relative path at fault.
//
// Row disposition is fixed per fault class rather than left to each caller, so
// rows_total, the row selector, and the status board agree on what a faulted board
// holds: a missing detail owner, a heading mismatch, an inline body, and every row of a
// degraded directory keep the index row, while a wrapped heading, the second position of
// a duplicated ID, an orphan, and an unrecognized file yield no row at all.
func ParseDocument(tree Tree, statuses map[string]string, full bool) (Document, []ParseFailure, []Diagnostic) {
	content := tree.Index.Data
	lines := strings.Split(string(content), "\n")
	doc := Document{Text: string(content)}
	var failures []ParseFailure
	var diagnostics []Diagnostic
	owners, unread := rowFileOwners(tree.Files)
	// rowed is the ID a row was built for, so the next index line carrying it is a
	// duplicate; claimed is every ID the index names at all, faulted lines included, so
	// the directory pass never calls a claimed row's file an orphan.
	rowed, claimed := map[string]bool{}, map[string]bool{}
	// A directory the classifier could not read yields no listing, so every row would
	// otherwise be told its detail owner is missing. Name the directory once instead:
	// nobody looked, so no row's owner is known to be absent.
	dirDegraded := degradedState(tree.DirState)
	if dirDegraded {
		diagnostics = append(diagnostics, Diagnostic{RoadmapDir + "/", fmt.Sprintf("%s detail directory: %s", tree.DirState, tree.DirReason)})
	}

	for i := 0; i < len(lines); {
		line := lines[i]
		m := roadmapStartRe.FindStringSubmatch(line)
		if m == nil {
			if strings.HasPrefix(line, "**") {
				raw, n := projectBody(line, full)
				failures = append(failures, ParseFailure{RoadmapFile, "malformed roadmap row", raw, n})
			}
			i++
			continue
		}
		id, at := m[1], i+1
		claimed[id] = true
		closeAt := strings.Index(m[2], "**")
		i++
		// Whatever sits under the heading is consumed either way: the split shape has
		// no body in the index, so these lines are only ever evidence of a fault.
		under := i
		for i < len(lines) && !strings.HasPrefix(lines[i], "**") && !strings.HasPrefix(lines[i], "## ") {
			i++
		}
		if closeAt < 0 {
			diagnostics = append(diagnostics, Diagnostic{RoadmapFile, fmt.Sprintf("wrapped heading at line %d; a row heading is one physical line", at)})
			continue
		}
		if rowed[id] {
			diagnostics = append(diagnostics, Diagnostic{RoadmapFile, fmt.Sprintf("duplicate row %s at line %d", id, at)})
			continue
		}
		rowed[id] = true

		row := RoadmapRow{ID: id, Title: strings.Trim(m[2][:closeAt], " —:-\t")}
		rowText := line
		switch file, owned := owners[id]; {
		case owned:
			// The row file is the whole row as a reader sees it: its first line
			// repeats the index line, so the ledger, the spec path, and the
			// trigger words are all read from this text.
			rowText = string(file.Data)
			first, body, _ := strings.Cut(rowText, "\n")
			if first != line {
				diagnostics = append(diagnostics, Diagnostic{rowFilePath(id), fmt.Sprintf("heading does not match %s row %s", RoadmapFile, id)})
			}
			row.Body, row.BodyBytes = projectBody(strings.TrimSpace(body), full)
		case unread[id], dirDegraded:
			// The listing pass has already named the file — or the directory — it could
			// not read; the row keeps its place on the board with the body nobody was
			// able to load.
		default:
			diagnostics = append(diagnostics, Diagnostic{rowFilePath(id), fmt.Sprintf("missing detail owner for %s row %s", RoadmapFile, id)})
		}
		if strings.TrimSpace(m[2][closeAt+2:]) != "" || anyNonBlank(lines[under:i]) {
			diagnostics = append(diagnostics, Diagnostic{RoadmapFile, fmt.Sprintf("row %s carries an inline body; move it to %s", id, rowFilePath(id))})
		}

		if slugs := spec.LiveSpecSlugs([]byte(rowText)); len(slugs) > 0 {
			row.Spec = slugs[0]
			row.SpecStatus = statuses[slugs[0]]
		}
		lower := strings.ToLower(rowText)
		row.ExternalTrigger = strings.Contains(lower, "pending ") || strings.Contains(lower, "graduate on") || strings.Contains(lower, "scheduled")
		if keys, count, valid := parseOccurrenceLedger(strings.Split(rowText, "\n")); valid {
			row.OccurrenceKeys, row.OccurrenceCount = keys, count
		} else {
			doc.OccurrenceDiscrepancies = append(doc.OccurrenceDiscrepancies, OccurrenceDiscrepancy{Source: rowFilePath(id), CaptureUnit: id, Kind: "malformed-ledger", Owner: id, Structural: true})
			failures = append(failures, ParseFailure{rowFilePath(id), "malformed-ledger", "", 0})
		}
		doc.Rows = append(doc.Rows, row)
	}

	diagnostics = append(diagnostics, listingDiagnostics(tree.Files, claimed)...)
	sequence, sequenceText, hasSection := parseSequence(lines)
	doc.Sequence, doc.SequenceText = sequence, sequenceText
	if len(content) > 0 && len(doc.Rows) == 0 && len(failures) == 0 && len(diagnostics) == 0 && !hasSection {
		raw, n := projectBody(string(content), full)
		failures = append(failures, ParseFailure{RoadmapFile, noRoadmapRowsReason, raw, n})
	}
	return doc, failures, diagnostics
}

// rowFileOwners splits the listing into the row files a row may be read from and the
// row IDs whose file was there but could not be read. An empty file is an owner, not an
// absence: its heading is missing rather than its detail, which is the heading mismatch
// the row already reports.
func rowFileOwners(files []RowFile) (owners map[string]RowFile, unread map[string]bool) {
	owners, unread = map[string]RowFile{}, map[string]bool{}
	for _, file := range files {
		id, ok := rowFileID(file.Name)
		if !ok {
			continue
		}
		if file.State == bounds.StateParsed || file.State == bounds.StateEmpty {
			owners[id] = file
			continue
		}
		unread[id] = true
	}
	return owners, unread
}

// listingDiagnostics reports what the roadmap/ listing holds that no index row claimed,
// in directory order: a basename outside the row-ID grammar, a row file the classifier
// could not read, and a detail file whose row is not on the board. claimed carries every
// ID the index named, not only the ones that became rows: a file whose index line is
// itself faulted has an owner on the board and is not an orphan.
func listingDiagnostics(files []RowFile, claimed map[string]bool) []Diagnostic {
	var diagnostics []Diagnostic
	for _, file := range files {
		id, ok := rowFileID(file.Name)
		switch {
		case !ok:
			diagnostics = append(diagnostics, Diagnostic{RoadmapDir + "/" + file.Name, fmt.Sprintf("unrecognized file under %s/; expected <row ID>.md", RoadmapDir)})
		case file.State != bounds.StateParsed && file.State != bounds.StateEmpty:
			diagnostics = append(diagnostics, Diagnostic{rowFilePath(id), fmt.Sprintf("%s detail file: %s", file.State, file.Reason)})
		case !claimed[id]:
			diagnostics = append(diagnostics, Diagnostic{rowFilePath(id), fmt.Sprintf("orphan detail file with no %s row %s", RoadmapFile, id)})
		}
	}
	return diagnostics
}

// rowFileID reports the row a detail file owns, or false when its basename is not
// `<row ID>.md` under the row-ID grammar `--row` and the index line share.
func rowFileID(name string) (string, bool) {
	id := strings.TrimSuffix(name, ".md")
	if id == name || !ValidRowID(id) {
		return "", false
	}
	return id, true
}

func anyNonBlank(lines []string) bool {
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			return true
		}
	}
	return false
}

// parseSequence extracts the `## Recommended sequence` section, from its heading
// (trailing whitespace tolerated) to the next `## ` heading or EOF, and reports
// whether the index carries any unfenced `## ` heading at all. Fence state is tracked
// throughout: a heading inside a fenced code block neither starts the section nor
// terminates it.
func parseSequence(lines []string) (rows []SequenceRow, text string, hasSection bool) {
	inSequence, inFence, sequenceStart, sequenceEnd := false, false, -1, len(lines)
	for idx, line := range lines {
		trimmed := strings.TrimRight(line, " \t\r")
		if strings.HasPrefix(trimmed, "```") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		// Any unfenced `## ` heading — not only Recommended sequence — is roadmap
		// structure: a document with sections is a working roadmap the reader
		// recognizes, possibly one with nothing due yet, never "not the document
		// you think".
		if strings.HasPrefix(trimmed, "## ") {
			hasSection = true
		}
		if !inSequence && trimmed == "## Recommended sequence" {
			inSequence = true
			sequenceStart = idx
			continue
		}
		if inSequence && strings.HasPrefix(trimmed, "## ") {
			sequenceEnd = idx
			break
		}
		if !inSequence {
			continue
		}
		parts := strings.SplitN(line, ". ", 2)
		if len(parts) != 2 {
			continue
		}
		rank, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		cmd := ""
		if m := commandRe.FindString(parts[1]); m != "" {
			cmd = m
		}
		rows = append(rows, SequenceRow{Rank: rank, Text: parts[1], Command: cmd})
	}
	if sequenceStart >= 0 {
		text = strings.Join(lines[sequenceStart:sequenceEnd], "\n")
		if !strings.HasSuffix(text, "\n") {
			text += "\n"
		}
	}
	return rows, text, hasSection
}
