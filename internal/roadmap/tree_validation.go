package roadmap

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
