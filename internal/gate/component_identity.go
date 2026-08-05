package gate

// The content address a scoped component's ancestor evidence is stored under. A component
// may inherit that evidence only while its identity is unmoved, so an identity blind to
// something the component reads credits work nobody graded. Everything a component reads
// therefore reaches the hash: its declared input contents, selected positively out of the
// same `git add -A` snapshot the stripped identity reads, and the execution closure that
// decides what running it would actually do.
//
// Nothing here has a fallback. An identity that cannot be computed is returned as an error,
// and the decision that consumes it runs the component instead — an identity assembled from
// the inputs that happened to resolve is exactly the shape that buys a wrong skip.

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"sort"
)

// componentPolicyVersion is the per-component identity's hash domain. It extends the
// whole-tree policy rather than standing beside it, so a policy change moves every
// component's identity with it.
const componentPolicyVersion = policyVersion + "/component-v1"

// componentPolicyDomain is the domain one component hashes under. The component's name is
// part of the domain, so two components whose declarations cover identical files still
// compute different addresses and neither one's slot can answer for the other.
func componentPolicyDomain(component string) string {
	return componentPolicyVersion + "/" + component
}

func resolveComponentIdentities(root string, generation *treeGeneration) (map[string]string, error) {
	sets, err := ResolveComponentInputs(root)
	if err != nil {
		return nil, err
	}
	snapshot := generation.snapshot
	phases := map[string]Phase{}
	for _, phase := range BenchkitPhases(root, kitRoot(root)) {
		phases[phase.Name] = phase
	}
	identities := make(map[string]string, len(sets))
	for name, inputs := range sets {
		phase, materialized := phases[name]
		if !materialized {
			continue
		}
		identity, err := componentIdentity(root, inputs, phase, snapshot)
		if err != nil {
			return nil, fmt.Errorf("identify %s: %w", name, err)
		}
		identities[name] = identity
	}
	return identities, nil
}

// componentIdentity is the address inputs and phase resolve to against snapshot. Every
// snapshot entry a declaration covers contributes its metadata — the mode, type, and object
// id git recorded — so the address answers for the declared content and not for the names
// it was declared under.
//
// A declared file the snapshot has no entry for is a refusal rather than a silently
// unhashed input: the declaration names that file exactly, so its absence is the
// declaration and the tree disagreeing at a named point.
//
// A declared directory covering nothing contributes nothing, and that is not the same
// judgment. Git tracks no empty directory, so a directory that is empty and one whose
// surfaces have not landed yet are indistinguishable in the snapshot, and refusing would
// leave a component permanently unable to compute an identity. Nothing is bought by the
// silence either: the first file to land under it joins the hash and moves the identity,
// so an empty directory can never carry an ungraded change.
func componentIdentity(root string, inputs ComponentInputs, phase Phase, snapshot treeSnapshot) (string, error) {
	h := sha256.New()
	frame(h, componentPolicyDomain(inputs.Component()))
	for _, declared := range inputs.Paths() {
		if err := confinedRepoPath(root, declared); err != nil {
			return "", fmt.Errorf("declared input %q: %w", declared, err)
		}
		covered := snapshotEntriesUnder(snapshot, declared)
		if len(covered) == 0 && !declaresDirectory(declared) {
			return "", fmt.Errorf("declared input %q is absent from the snapshot", declared)
		}
		for _, entry := range covered {
			frameEach(h, "path", entry.Path)
			frameEach(h, "content", entry.Metadata)
		}
	}
	// The execution closure: what running the component would do, for a tree its inputs
	// call unmoved. An argv or an environment contract that moved runs a different check,
	// and a seal digest that moved execs a different binary.
	frameEach(h, "argv", phase.Argv...)
	frameEach(h, "env", phase.Env...)
	frameEach(h, "dir", phase.Dir)
	frameEach(h, "seal", inputs.Digests()...)
	return hex.EncodeToString(h.Sum(nil)), nil
}

// snapshotEntriesUnder returns the snapshot entries a declared entry covers, sorted by path.
// The order is the hash's, so it is fixed here rather than inherited from the listing — a
// component whose files were framed in a different order each call would address a different
// slot each call.
func snapshotEntriesUnder(snapshot treeSnapshot, declared string) []treeEntry {
	// A file entry is resolved by lookup, not by scanning: what declaredEntryCovers states
	// for it is byte-exact equality, and an index lookup is that comparison. Scanning the
	// whole snapshot for each of the hundreds of files a closure declares is the same answer
	// paid for once per entry in the tree.
	if !declaresDirectory(declared) {
		entry, tracked := snapshot.entry(declared)
		if !tracked {
			return nil
		}
		return []treeEntry{entry}
	}
	var covered []treeEntry
	for _, entry := range snapshot.entries {
		if declaredEntryCovers(declared, entry.Path) {
			covered = append(covered, entry)
		}
	}
	sort.Slice(covered, func(i, j int) bool { return covered[i].Path < covered[j].Path })
	return covered
}

// confinedRepoPath reports whether rel names a file inside root. A declaration is
// repository-relative and slash-separated by contract; anything absolute or reaching above
// root names content no changeset under root can be compared against.
func confinedRepoPath(root, rel string) error {
	native := filepath.FromSlash(rel)
	if rel == "" || filepath.IsAbs(native) || !withinRoot(root, filepath.Join(root, native)) {
		return errors.New("outside the repository")
	}
	return nil
}

// frameEach writes each value under its own tag, so a value can never be mistaken for the
// tag that opens the next group.
func frameEach(w io.Writer, tag string, values ...string) {
	for _, value := range values {
		frame(w, tag)
		frame(w, value)
	}
}
