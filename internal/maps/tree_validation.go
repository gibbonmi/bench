package maps

import (
	"fmt"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/bounds"
)

// ValidateDecisionMapTree validates every discovered active and compiled decision map.
func ValidateDecisionMapTree(root string) []string {
	candidates, diagnostics := discoverDecisionMapCandidates(root)
	for _, candidate := range candidates {
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
