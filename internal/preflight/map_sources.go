package preflight

import (
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/maps"
)

// mapSourceDiagnostics grades the Sources records of every candidate decision map at tip
// against the published tree, in candidate path order. It lists the tip tree under the
// active and compiled map roots, keeps each path that the shared candidate rule admits, and
// reads each candidate through the bounded blob read. The published-tree resolver also reds
// each source that the predicted delivery closure of specPath deletes. A prediction error
// is one diagnostic, and the maps are still graded. Only build mode fills the list.
func mapSourceDiagnostics(root, mode, tip, specPath string) []string {
	if mode != modeBuild {
		return nil
	}
	entries, err := git.ListTree(root, "-r", tip, "--", maps.DecisionsDir, "specs")
	if err != nil {
		return []string{"decision maps not readable at the source tip: " + err.Error()}
	}
	var diagnostics []string
	deleted, err := (commitrepo.Store{Root: root}).ClosedRowPaths(tip, specPath)
	if err != nil {
		diagnostics = append(diagnostics, "delivery closure of "+specPath+" not predictable: "+err.Error())
	}
	resolve := maps.PublishedTreeResolver(root, tip, specPath, deleted)
	for _, entry := range entries {
		if _, admitted := maps.DecisionMapCandidatePath(entry.Path); !admitted {
			continue
		}
		// A link entry reads as its target string, which holds no Sources section.
		if !entry.IsRegularFile() {
			diagnostics = append(diagnostics, entry.Path+": candidate map has tree mode "+entry.Mode+", not a regular file")
			continue
		}
		data, err := git.ReadControlBlob(root, entry.Object)
		if err != nil {
			diagnostics = append(diagnostics, entry.Path+": not readable at the source tip: "+err.Error())
			continue
		}
		for _, diagnostic := range maps.SourceDiagnostics(entry.Path, data, resolve) {
			diagnostics = append(diagnostics, diagnostic.Message)
		}
	}
	return diagnostics
}
