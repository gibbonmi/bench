package releaseevidence

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type releasePlan struct {
	SchemaVersion     int                `json:"schema_version"`
	TargetCardinality int                `json:"target_cardinality"`
	Targets           []targetEvidence   `json:"targets"`
	ArchiveEntries    []archivePlanEntry `json:"archive_entries"`
}

type archivePlanEntry struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Kind string `json:"kind"`
}

type releaseArtifact struct {
	Name   string `json:"name"`
	Target string `json:"target"`
	Kind   string `json:"kind"`
}

func readReleasePlan(root string) (releasePlan, error) {
	data, err := releasePlanOutput(root, "normalized-json")
	if err != nil {
		return releasePlan{}, fmt.Errorf("release plan is unreadable: %w", err)
	}
	var plan releasePlan
	if err := decodeStrict(data, &plan); err != nil {
		return releasePlan{}, fmt.Errorf("release plan is malformed: %w", err)
	}
	return plan, nil
}

func releasePlanOutput(root string, arguments ...string) ([]byte, error) {
	args := append([]string{filepath.Join(root, "scripts", "release-plan.mjs"), root}, arguments...)
	return exec.Command("node", args...).Output()
}

func readReleaseArtifacts(root, version string) ([]releaseArtifact, error) {
	data, err := releasePlanOutput(root, "artifact-records", version)
	if err != nil {
		return nil, fmt.Errorf("release artifact inventory is unavailable: %w", err)
	}
	var artifacts []releaseArtifact
	if err := decodeStrict(data, &artifacts); err != nil {
		return nil, fmt.Errorf("release artifact inventory is malformed: %w", err)
	}
	return artifacts, nil
}

func archiveInventory(plan releasePlan, target, version string) (map[string]int64, error) {
	files := map[string]int64{}
	for _, entry := range plan.ArchiveEntries {
		paths := []string{entry.Path}
		if entry.Kind == "package_evidence" {
			paths = paths[:0]
			for _, evidence := range PackageEvidenceRegistry() {
				paths = append(paths, strings.Replace(entry.Path, "{package_evidence}", evidence.Path, 1))
			}
		}
		for _, name := range paths {
			name = strings.ReplaceAll(name, "{version}", version)
			name = strings.ReplaceAll(name, "{target}", target)
			if _, exists := files[name]; exists {
				return nil, fmt.Errorf("release plan archive inventory duplicates %s", name)
			}
			if entry.Mode == "0755" {
				files[name] = 0o755
			} else {
				files[name] = 0o644
			}
		}
	}
	return files, nil
}
