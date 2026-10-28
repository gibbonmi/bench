package gate

import (
	"io"
	"path/filepath"
)

// noRunBinary is the grade identity member of a run that selects no run binary. It is one
// stable value, so the build inputs of a source that the run never builds do not move its key.
const noRunBinary = "no run binary"

// runBinarySelection answers the source root of the run binary that a run selects, and false
// when the run selects none. A run selects one only when it has an owner and its resolved gate
// routes through the phase table. The source is the graded tree when a prospective gate grades
// a tree that builds its own run binary, and the kit otherwise. The run's selection and the
// grade identity both read it.
func runBinarySelection(owned bool, res Resolution, runtimeRoot, storageRoot string) (string, bool) {
	if !owned || !phaseTableGate(runtimeRoot, res) {
		return "", false
	}
	if res.Kind == ProspectiveGateSh &&
		isRegularFile(filepath.Join(runtimeRoot, "scripts", "go-build.sh")) &&
		GoModuleRoot(runtimeRoot) {
		return runtimeRoot, true
	}
	return kitRoot(storageRoot), true
}

// runBinaryMember answers the grade identity member of a run: the build inputs of the source
// root that the run selects, or noRunBinary when it selects none.
func runBinaryMember(owned bool, res Resolution, runtimeRoot, storageRoot string) string {
	source, selected := runBinarySelection(owned, res, runtimeRoot, storageRoot)
	if !selected {
		return noRunBinary
	}
	return buildInputsIdentity(source)
}

// frameRunBinary sets the resolution that the run of s executes, and frames the run-binary
// member of that run on id. An ordinary run executes its resolved gate under the production
// owner. A prospective run always has an owner, and it executes the prospective script when
// root holds one.
func (s *subject) frameRunBinary(id io.Writer, root, identityRoot string, prospective bool) error {
	owned := true
	if !prospective {
		owned = productionRunBinaryOwner() != nil
	} else if resolved, err := prospectiveResolution(root, s.Resolution); err != nil {
		return err
	} else {
		s.Resolution = resolved
	}
	frame(id, "run binary")
	frame(id, runBinaryMember(owned, s.Resolution, root, identityRoot))
	return nil
}
