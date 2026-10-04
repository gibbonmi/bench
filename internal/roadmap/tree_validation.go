package roadmap

import (
	"fmt"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"path/filepath"
	"strings"
)

// ValidateRoadmapTree grades the split board at root and returns the loader's ordered
// integrity diagnostics, rendered to strings. It is the conformance check's whole
// implementation. The parse already derives every fault class, so the gate reads the same
// tree the board command, the context snapshot, and the owner check read. This avoids a
// second reading that could disagree with them. A repo with neither ROADMAP.md nor
// roadmap/ yields nothing, so a tree with no board is quiet rather than red. The registry
// binding this feeds wants strings, not the typed Diagnostic every in-package caller
// carries, so the conversion happens once, at this one public boundary.
func ValidateRoadmapTree(root string) []string {
	_, _, diagnostics := ParseDocument(LoadTree(root), nil, false)
	if len(diagnostics) == 0 {
		return nil
	}
	rendered := make([]string, len(diagnostics))
	for i, d := range diagnostics {
		rendered[i] = d.String()
	}
	return rendered
}

// RequirementBytes excludes only a validated occurrence ledger from a detail owner.
func RequirementBytes(name string, data []byte) ([]byte, error) {
	if filepath.ToSlash(filepath.Dir(name)) != RoadmapDir {
		return data, nil
	}
	if _, ok := rowFileID(filepath.Base(name)); !ok {
		return data, nil
	}
	lines := strings.Split(string(data), "\n")
	if _, _, valid := parseOccurrenceLedger(lines); !valid {
		return nil, fmt.Errorf("%s has a malformed occurrence ledger", name)
	}
	kept := make([]string, 0, len(lines))
	for _, line := range lines {
		if !isOccurrenceLedgerLine(line) {
			kept = append(kept, line)
		}
	}
	return []byte(strings.Join(kept, "\n")), nil
}

// RevisionDocument parses an immutable roadmap through the canonical split-board parser.
func RevisionDocument(root, revision string) (Document, error) {
	tree := Tree{Index: bounds.Classified{State: bounds.StateAbsent}, DirState: bounds.StateAbsent}
	index, err := git.Output("-C", root, "ls-tree", "-z", revision, "--", RoadmapFile)
	if err != nil {
		return Document{}, err
	}
	if index != "" {
		data, err := git.ReadTreeFile(root, revision, RoadmapFile)
		if err != nil {
			return Document{}, err
		}
		tree.Index = bounds.Classified{State: bounds.StateParsed, Data: data}
	}
	directory, err := git.Output("-C", root, "ls-tree", "-z", revision, "--", RoadmapDir)
	if err != nil {
		return Document{}, err
	}
	if directory != "" {
		names, err := git.Output("-C", root, "ls-tree", "--name-only", "-z", revision+":"+RoadmapDir)
		if err != nil {
			return Document{}, err
		}
		tree.DirState = bounds.StateParsed
		for _, name := range strings.Split(names, "\x00") {
			if name == "" {
				continue
			}
			data, err := git.ReadTreeFile(root, revision, RoadmapDir+"/"+name)
			if err != nil {
				return Document{}, err
			}
			tree.Files = append(tree.Files, RowFile{Name: name, State: bounds.StateParsed, Data: data})
		}
	}
	document, failures, diagnostics := ParseDocument(tree, nil, true)
	if len(failures) != 0 || len(diagnostics) != 0 {
		return Document{}, fmt.Errorf("roadmap at %s is structurally untrusted", revision)
	}
	return document, nil
}

// RowPath returns the canonical detail owner of an existing row identity.
func RowPath(id string) string { return rowFilePath(id) }
