package releaseevidence

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func readReproducibility(root string, artifacts []artifactEvidence) (reproducibilityEvidence, error) {
	data, err := ReadRegular(filepath.Join(root, "dist", "reproducibility.json"))
	if err != nil {
		return reproducibilityEvidence{}, errors.New("reproducibility comparison record is missing")
	}
	var record reproducibilityEvidence
	if err := decodeStrict(data, &record); err != nil {
		return reproducibilityEvidence{}, fmt.Errorf("reproducibility comparison record is malformed: %w", err)
	}
	if record.SchemaVersion != 1 || record.Status != "green" || record.Builds != 2 || len(record.Artifacts) != len(artifacts) {
		return reproducibilityEvidence{}, errors.New("reproducibility comparison is incomplete")
	}
	want := make(map[string]artifactEvidence, len(artifacts))
	for _, artifact := range artifacts {
		want[artifact.Name] = artifact
	}
	seen := map[string]bool{}
	for _, item := range record.Artifacts {
		artifact, ok := want[item.Name]
		if !ok || seen[item.Name] || !item.Match || item.Size != artifact.Size || item.SHA256 != artifact.SHA256 {
			return reproducibilityEvidence{}, fmt.Errorf("reproducibility comparison does not match inspected artifact: %s", item.Name)
		}
		seen[item.Name] = true
	}
	if len(seen) != len(want) {
		return reproducibilityEvidence{}, errors.New("reproducibility comparison omits an inspected artifact")
	}
	return record, nil
}

func inspectNativeProofs(root string, targets []targetEvidence, artifacts []artifactEvidence) ([]nativeProofEvidence, error) {
	dir := filepath.Join(root, "dist", "native-proofs")
	info, err := os.Lstat(dir)
	if os.IsNotExist(err) {
		return []nativeProofEvidence{}, nil
	}
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return nil, errors.New("native proof directory is unsafe")
	}
	expected := make(map[string]bool, len(targets))
	for _, target := range targets {
		expected[target.OS+"-"+target.Arch+".json"] = true
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.New("native proof directory is unreadable")
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil || entry.Type()&os.ModeSymlink != 0 || !info.Mode().IsRegular() || !expected[entry.Name()] {
			return nil, fmt.Errorf("native proof directory contains an unexpected entry: %s", entry.Name())
		}
	}
	digestByName := make(map[string]string, len(artifacts))
	for _, artifact := range artifacts {
		digestByName[artifact.Name] = artifact.SHA256
	}
	proofs := make([]nativeProofEvidence, 0, len(targets))
	for _, target := range targets {
		name := target.OS + "-" + target.Arch
		data, err := ReadRegular(filepath.Join(dir, name+".json"))
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("native proof %s is unreadable", name)
		}
		var proof nativeProofEvidence
		if err := decodeStrict(data, &proof); err != nil {
			return nil, fmt.Errorf("native proof %s is malformed: %w", name, err)
		}
		platform := fmt.Sprintf("redbench-%s-%s-%s.tgz", target.OS, target.Arch, mustPackageVersion(root))
		archive := fmt.Sprintf("redbench-%s-%s-%s.tar.gz", mustPackageVersion(root), target.OS, target.Arch)
		packageFiles, _, err := readTarball(mustArtifactBytes(root, platform))
		if err != nil {
			return nil, fmt.Errorf("native proof %s cannot inspect platform binary: %w", name, err)
		}
		binaryDigest := digest(packageFiles["bin/bench"].data)
		if proof.SchemaVersion != 1 || proof.Target != name || proof.Runner != target.Runner || proof.Status != "green" || proof.RebuiltSHA256 == "" || proof.RebuiltSHA256 != binaryDigest || proof.BinarySHA256 != binaryDigest || proof.PackageSHA256 != digestByName[platform] || proof.ArchiveSHA256 != digestByName[archive] || (target.OS == "linux" && proof.MuslStatus != "green") || (target.OS == "darwin" && proof.MuslStatus != "not_applicable") {
			return nil, fmt.Errorf("native proof %s does not match inspected artifacts", name)
		}
		proofs = append(proofs, proof)
	}
	sort.Slice(proofs, func(i, j int) bool { return proofs[i].Target < proofs[j].Target })
	return proofs, nil
}

func mustArtifactBytes(root, name string) []byte {
	data, err := ReadRegular(filepath.Join(root, "dist", "artifacts", name))
	if err != nil {
		return nil
	}
	return data
}

func mustPackageVersion(root string) string {
	version, err := ReadPackageVersion(root)
	if err != nil {
		return ""
	}
	return version
}
