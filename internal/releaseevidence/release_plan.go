package releaseevidence

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
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

func readReleasePlan(root string) (releasePlan, error) {
	data, err := ReadRegular(filepath.Join(root, "scripts", "release-plan.json"))
	if err != nil {
		return releasePlan{}, fmt.Errorf("release plan is unreadable: %w", err)
	}
	var plan releasePlan
	if err := decodeStrict(data, &plan); err != nil {
		return releasePlan{}, fmt.Errorf("release plan is malformed: %w", err)
	}
	if plan.SchemaVersion != 1 || plan.TargetCardinality < 1 || len(plan.Targets) != plan.TargetCardinality || len(plan.ArchiveEntries) == 0 {
		return releasePlan{}, errors.New("release plan cardinality or inventory is invalid")
	}
	seenTargets := map[string]bool{}
	for _, target := range plan.Targets {
		key := target.OS + "-" + target.Arch
		if (target.OS != "darwin" && target.OS != "linux") || (target.Arch != "arm64" && target.Arch != "x64") || target.GOOS != target.OS || target.GOArch != map[string]string{"arm64": "arm64", "x64": "amd64"}[target.Arch] || target.Runner == "" || seenTargets[key] {
			return releasePlan{}, fmt.Errorf("release plan target is invalid or duplicate: %s", key)
		}
		seenTargets[key] = true
	}
	seenEntries := map[string]bool{}
	for _, entry := range plan.ArchiveEntries {
		clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(entry.Path)))
		if entry.Path == "" || clean != entry.Path || strings.Contains(entry.Path, "\\") || strings.HasPrefix(clean, "../") || filepath.IsAbs(clean) || (entry.Mode != "0644" && entry.Mode != "0755") || entry.Kind == "" || seenEntries[entry.Path] {
			return releasePlan{}, fmt.Errorf("release plan archive entry is invalid: %s", entry.Path)
		}
		seenEntries[entry.Path] = true
	}
	return plan, nil
}

func artifactNames(plan releasePlan, version string) []string {
	names := []string{"redbench-" + version + ".tgz"}
	for _, target := range plan.Targets {
		names = append(names, "redbench-"+target.OS+"-"+target.Arch+"-"+version+".tgz", "redbench-"+version+"-"+target.OS+"-"+target.Arch+".tar.gz")
	}
	sort.Strings(names)
	return names
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
