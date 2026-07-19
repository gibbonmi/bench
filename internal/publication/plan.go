// Package publication owns governed npm publication: a registry port with two
// adapters (a hermetic fixture adapter the gate exercises, and a public-npm
// adapter that shells the npm CLI for the runbook), a resumable first-publication
// state machine, and the durable publication-record.json that makes a retry
// idempotent. The package never reads a credential into evidence or the record.
package publication

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
)

// ArtifactRecord mirrors one row of `release-plan.mjs artifact-records`: the
// exact file name, its target (a platform pair or "wrapper"), and its kind. The
// naming and target-matrix logic has exactly one source, scripts/release-plan.mjs;
// this package only ever calls out to it, never reimplements it.
type ArtifactRecord struct {
	Name   string `json:"name"`
	Target string `json:"target"`
	Kind   string `json:"kind"`
}

// PublicationSet returns the exact four platform packages plus one wrapper —
// the build set a first publication must publish — filtered from the full
// release-plan artifact-records (which also lists offline .tar.gz archives).
// Order is whatever release-plan.mjs emits (sorted by file name); callers that
// need the platform-then-wrapper publication order impose it themselves so the
// ordering policy lives in the state machine, not here.
func PublicationSet(root, version string) ([]ArtifactRecord, error) {
	data, err := exec.Command("node", filepath.Join(root, "scripts", "release-plan.mjs"), root, "artifact-records", version).Output()
	if err != nil {
		return nil, fmt.Errorf("release artifact inventory is unavailable: %w", err)
	}
	var all []ArtifactRecord
	if err := json.Unmarshal(data, &all); err != nil {
		return nil, fmt.Errorf("release artifact inventory is malformed: %w", err)
	}
	var packages []ArtifactRecord
	for _, record := range all {
		if record.Kind == "wrapper" || record.Kind == "platform" {
			packages = append(packages, record)
		}
	}
	return packages, nil
}
