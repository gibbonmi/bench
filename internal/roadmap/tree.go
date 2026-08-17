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
		c := bounds.Classify(filepath.Join(root, RoadmapDir, entry.Name()), bounds.ControlRecordLimit)
		tree.Files = append(tree.Files, RowFile{Name: entry.Name(), Reason: c.Reason, State: c.State, Data: c.Data})
	}
	return tree
}

// rowFilePath is the repo-relative detail owner of one row ID. Diagnostics, the
// migration, and every reader that names a row's file go through it.
func rowFilePath(id string) string { return RoadmapDir + "/" + id + ".md" }

// ParseDocument projects a classified tree into the Document every roadmap surface
// renders, the parse failures the snapshot reports, and the ordered integrity
// diagnostics the conformance check returns. Diagnostics run in index order and then in
// directory order, and each begins with the repo-relative path of the file at fault.
//
// Row disposition is fixed per fault class rather than left to each caller, so
// rows_total, the row selector, and the status board agree on what a faulted board
// holds: a missing detail owner, a heading mismatch, and an inline body keep the index
// row, while a wrapped heading, the second position of a duplicated ID, an orphan, and
// an unrecognized file yield no row at all.
func ParseDocument(tree Tree, statuses map[string]string, full bool) (Document, []ParseFailure, []string) {
	content := tree.Index.Data
	lines := strings.Split(string(content), "\n")
	doc := Document{Text: string(content)}
	var failures []ParseFailure
	var diagnostics []string
	owners, unread := rowFileOwners(tree.Files)
	indexed := map[string]bool{}

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
		closeAt := strings.Index(m[2], "**")
		i++
		// Whatever sits under the heading is consumed either way: the split shape has
		// no body in the index, so these lines are only ever evidence of a fault.
		under := i
		for i < len(lines) && !strings.HasPrefix(lines[i], "**") && !strings.HasPrefix(lines[i], "## ") {
			i++
		}
		if closeAt < 0 {
			diagnostics = append(diagnostics, fmt.Sprintf("%s: wrapped heading at line %d; a row heading is one physical line", RoadmapFile, at))
			continue
		}
		if indexed[id] {
			diagnostics = append(diagnostics, fmt.Sprintf("%s: duplicate row %s at line %d", RoadmapFile, id, at))
			continue
		}
		indexed[id] = true

		row := RoadmapRow{ID: id, Title: strings.Trim(m[2][:closeAt], " —:-\t")}
		rowText := line
		switch file, owned := owners[id]; {
		case owned:
			// The row file is the whole row as a reader sees it: its first line
			// repeats the index line, so the ledger, the spec path, and the
			// trigger words read exactly the text they read before the split.
			rowText = string(file.Data)
			first, body, _ := strings.Cut(rowText, "\n")
			if first != line {
				diagnostics = append(diagnostics, fmt.Sprintf("%s: heading does not match %s row %s", rowFilePath(id), RoadmapFile, id))
			}
			row.Body, row.BodyBytes = projectBody(strings.TrimSpace(body), full)
		case unread[id]:
			// The directory pass names the file it could not read; the row keeps its
			// place on the board with the body nobody was able to load.
		default:
			diagnostics = append(diagnostics, fmt.Sprintf("%s: missing detail owner for %s row %s", rowFilePath(id), RoadmapFile, id))
		}
		if strings.TrimSpace(m[2][closeAt+2:]) != "" || anyNonBlank(lines[under:i]) {
			diagnostics = append(diagnostics, fmt.Sprintf("%s: row %s carries an inline body; move it to %s", RoadmapFile, id, rowFilePath(id)))
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

	diagnostics = append(diagnostics, listingDiagnostics(tree.Files, indexed)...)
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
// could not read, and a detail file whose row is not on the board.
func listingDiagnostics(files []RowFile, indexed map[string]bool) []string {
	var diagnostics []string
	for _, file := range files {
		id, ok := rowFileID(file.Name)
		switch {
		case !ok:
			diagnostics = append(diagnostics, fmt.Sprintf("%s/%s: unrecognized file under %s/; expected <row ID>.md", RoadmapDir, file.Name, RoadmapDir))
		case file.State != bounds.StateParsed && file.State != bounds.StateEmpty:
			diagnostics = append(diagnostics, fmt.Sprintf("%s: %s detail file: %s", rowFilePath(id), file.State, file.Reason))
		case !indexed[id]:
			diagnostics = append(diagnostics, fmt.Sprintf("%s: orphan detail file with no %s row %s", rowFilePath(id), RoadmapFile, id))
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
