package roadmap

// ValidateRoadmapTree grades the split board at root and returns the loader's ordered
// integrity diagnostics. It is the conformance check's whole implementation: the parse
// already derives every fault class, so the gate reads the same tree the board command,
// the context snapshot, and the owner check read rather than a second reading that could
// disagree with them. A repo with neither ROADMAP.md nor roadmap/ yields nothing, so a
// tree with no board is quiet rather than red.
func ValidateRoadmapTree(root string) []string {
	_, _, diagnostics := ParseDocument(LoadTree(root), nil, false)
	return diagnostics
}
