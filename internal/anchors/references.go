package anchors

import (
	"fmt"
	"go/scanner"
	"go/token"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
)

// RegistryDir is the repository-relative anchor registry directory.
const RegistryDir = "internal/anchors"

// ReferencingFiles maps each Go string literal in root's anchor registry directory to
// the sorted repository-relative files that hold it. The scan reads every top-level
// `.go` file, test files included, so an independent expectation joins its registry
// row. It does not recurse. A file names a path by a key, a constant, or a positional
// value alike, so the scan reads every interpreted and raw literal rather than one
// field.
//
// The scan classifies each file without a link followed and before any open. A link,
// a special file, an unreadable file, or source that does not tokenize refuses the
// whole scan and names the path, so a partial answer never reads as a complete one.
// An absent directory is an empty answer, because a linked repository keeps no anchor
// registry.
func ReferencingFiles(root string) (map[string][]string, error) {
	refs := map[string][]string{}
	dir := bounds.ClassifyDirNoFollow(filepath.Join(root, filepath.FromSlash(RegistryDir)))
	switch dir.State {
	case bounds.StateAbsent, bounds.StateEmpty:
		return refs, nil
	case bounds.StateParsed:
	default:
		return nil, fmt.Errorf("anchor registry %s refused: %s", RegistryDir, dir.Reason)
	}
	for _, entry := range dir.Entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}
		rel := path.Join(RegistryDir, entry.Name())
		literals, err := fileLiterals(filepath.Join(root, filepath.FromSlash(rel)), rel)
		if err != nil {
			return nil, err
		}
		for _, literal := range literals {
			refs[literal] = append(refs[literal], rel)
		}
	}
	for literal, files := range refs {
		sort.Strings(files)
		refs[literal] = files
	}
	return refs, nil
}

// fileLiterals returns the distinct unquoted string literals of one registry file.
func fileLiterals(abs, rel string) ([]string, error) {
	classified := bounds.ClassifyNoFollow(abs)
	switch classified.State {
	case bounds.StateParsed:
	case bounds.StateEmpty:
		return nil, nil
	default:
		return nil, fmt.Errorf("anchor registry file %s refused: %s %s", rel, classified.State, classified.Reason)
	}
	var fault error
	var s scanner.Scanner
	fset := token.NewFileSet()
	s.Init(fset.AddFile(rel, -1, len(classified.Data)), classified.Data, func(pos token.Position, msg string) {
		if fault == nil {
			fault = fmt.Errorf("anchor registry file %s refused: %s: %s", rel, pos, msg)
		}
	}, 0)
	seen := map[string]bool{}
	var literals []string
	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		if tok != token.STRING {
			continue
		}
		value, err := strconv.Unquote(lit)
		if err != nil || seen[value] {
			continue
		}
		seen[value] = true
		literals = append(literals, value)
	}
	if fault != nil {
		return nil, fault
	}
	return literals, nil
}
