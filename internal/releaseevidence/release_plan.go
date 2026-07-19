package releaseevidence

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

type releasePlan struct {
	SchemaVersion  int              `json:"schema_version"`
	Targets        []targetEvidence `json:"targets"`
	ArchiveEntries []map[string]any `json:"archive_entries"`
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

func archiveInventory(root, target, version string) (map[string]int64, error) {
	data, err := releasePlanOutput(root, "archive-inventory", target, version)
	if err != nil {
		return nil, fmt.Errorf("release plan archive inventory is unavailable: %w", err)
	}
	var files map[string]int64
	if err := decodeStrict(data, &files); err != nil {
		return nil, fmt.Errorf("release plan archive inventory is malformed: %w", err)
	}
	return files, nil
}

func archiveEntryPath(root, kind, target, version string) (string, error) {
	return releasePlanPath(root, "archive-entry-path", kind, target, version)
}

func archiveEvidencePath(root, requirementKey, target, version string) (string, error) {
	return releasePlanPath(root, "archive-evidence-path", requirementKey, target, version)
}

func releasePlanPath(root string, arguments ...string) (string, error) {
	data, err := releasePlanOutput(root, arguments...)
	if err != nil {
		return "", fmt.Errorf("release plan archive path is unavailable: %w", err)
	}
	path := string(data)
	if len(path) == 0 || path[len(path)-1] != '\n' {
		return "", fmt.Errorf("release plan archive path is malformed")
	}
	return path[:len(path)-1], nil
}
