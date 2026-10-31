package maps

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
)

// DecisionMapCandidate identifies a directly discovered map and its ownership.
type DecisionMapCandidate struct {
	Path     string
	Compiled bool
}

func isDirectoryDoc(name string) bool {
	return strings.EqualFold(strings.TrimSuffix(name, filepath.Ext(name)), "README")
}

// DecisionMapCandidatePath is the one candidate rule. It admits a repository-relative slash
// path as an active map, decisions/<name>.md, or as a compiled map,
// specs/<folder>/decisions/<name>.md, and reports which. Neither the folder nor the name
// starts with a dot, and the name is not a README.
func DecisionMapCandidatePath(path string) (compiled, admitted bool) {
	parts := strings.Split(path, "/")
	switch {
	case len(parts) == 2 && parts[0] == DecisionsDir:
	case len(parts) == 4 && parts[0] == "specs" && !strings.HasPrefix(parts[1], ".") && parts[2] == DecisionsDir:
		compiled = true
	default:
		return false, false
	}
	name := parts[len(parts)-1]
	return compiled, !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".md") && !isDirectoryDoc(name)
}

// discoverDirectoryCandidates lists the direct-child candidates of one active or compiled
// decision-map directory. DecisionMapCandidatePath decides each file entry, so callers
// choose only which directories to compose.
func discoverDirectoryCandidates(root, dir string) ([]DecisionMapCandidate, bounds.FileState, string) {
	classified := bounds.ClassifyDir(dir)
	if classified.State != bounds.StateParsed {
		return nil, classified.State, classified.Reason
	}
	var candidates []DecisionMapCandidate
	for _, entry := range classified.Entries {
		if entry.IsDir() {
			continue
		}
		path, err := filepath.Rel(root, filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, bounds.StateUnreadable, err.Error()
		}
		path = filepath.ToSlash(path)
		if compiled, admitted := DecisionMapCandidatePath(path); admitted {
			candidates = append(candidates, DecisionMapCandidate{Path: path, Compiled: compiled})
		}
	}
	return candidates, bounds.StateParsed, ""
}

// DiscoverDecisionMapCandidates finds active and compiled direct Markdown candidates.
func DiscoverDecisionMapCandidates(root string) ([]DecisionMapCandidate, error) {
	candidates, diagnostics := discoverDecisionMapCandidates(root)
	if len(diagnostics) > 0 {
		return nil, fmt.Errorf("%s", diagnostics[0])
	}
	return candidates, nil
}

// discoverDecisionMapCandidates is the sole active-and-compiled candidate traversal.
// Callers receive every readable candidate plus any independent directory diagnostics.
func discoverDecisionMapCandidates(root string) ([]DecisionMapCandidate, []string) {
	var candidates []DecisionMapCandidate
	var diagnostics []string
	appendDirectory := func(dir string) {
		discovered, state, reason := discoverDirectoryCandidates(root, dir)
		if state == bounds.StateAbsent {
			return
		}
		if state != bounds.StateParsed && state != bounds.StateEmpty {
			rel, err := filepath.Rel(root, dir)
			if err != nil {
				rel = dir
			}
			diagnostics = append(diagnostics, fmt.Sprintf("%s: %s: %s", filepath.ToSlash(rel), state, reason))
			return
		}
		candidates = append(candidates, discovered...)
		diagnostics = append(diagnostics, orphanTicketFolders(root, dir)...)
	}
	appendDirectory(filepath.Join(root, DecisionsDir))
	specs := filepath.Join(root, "specs")
	classified := bounds.ClassifyDir(specs)
	if classified.State == bounds.StateAbsent {
		sort.Slice(candidates, func(i, j int) bool { return candidates[i].Path < candidates[j].Path })
		return candidates, diagnostics
	}
	if classified.State != bounds.StateParsed && classified.State != bounds.StateEmpty {
		diagnostics = append(diagnostics, fmt.Sprintf("specs: %s: %s", classified.State, classified.Reason))
		return candidates, diagnostics
	}
	for _, spec := range classified.Entries {
		if !spec.IsDir() || strings.HasPrefix(spec.Name(), ".") {
			continue
		}
		appendDirectory(filepath.Join(specs, spec.Name(), DecisionsDir))
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].Path < candidates[j].Path })
	return candidates, diagnostics
}

// ValidateDecisionMapTree validates every discovered active and compiled decision map.
func ValidateDecisionMapTree(root string) []string {
	candidates, diagnostics := discoverDecisionMapCandidates(root)
	active := make(map[string]string)
	for _, candidate := range candidates {
		if !candidate.Compiled {
			active[filepath.Base(candidate.Path)] = candidate.Path
		}
	}
	for _, candidate := range candidates {
		if path, exists := active[filepath.Base(candidate.Path)]; candidate.Compiled && exists {
			diagnostics = append(diagnostics, fmt.Sprintf("%s: active map also compiled at %s; move the map instead of copying it", path, candidate.Path))
		}

		file := bounds.Classify(filepath.Join(root, filepath.FromSlash(candidate.Path)), bounds.ControlRecordLimit)
		if file.State != bounds.StateParsed && file.State != bounds.StateEmpty {
			diagnostics = append(diagnostics, fmt.Sprintf("%s: %s: %s", candidate.Path, file.State, file.Reason))
			continue
		}
		_, mapDiagnostics := ValidateDecisionMap(root, candidate.Path, candidate.Compiled, file.Data)
		for _, diagnostic := range mapDiagnostics {
			diagnostics = append(diagnostics, diagnostic.Message)
		}
	}
	return diagnostics
}
