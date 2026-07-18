package releaseevidence

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type tarFile struct {
	mode int64
	data []byte
}

const (
	maxArchiveCompressedSize int64 = 128 << 20
	maxArchiveMemberCount          = 10_000
	maxArchiveExpandedSize   int64 = 64 << 20
)

var maxArchiveMemberSize int64 = 256 << 20

func SetArchiveMemberLimitForTesting(limit int64) func() {
	previous := maxArchiveMemberSize
	maxArchiveMemberSize = limit
	return func() { maxArchiveMemberSize = previous }
}

func ValidateTarballForTesting(data []byte) error {
	_, _, err := readTarball(data)
	return err
}

type componentManifest struct {
	SchemaVersion int
	Component     componentIdentity
	Files         []manifestFile
}

type componentIdentity struct {
	Name    string
	Version string
	Target  targetName
}

type targetName struct {
	OS   string
	Arch string
}

type manifestFile struct {
	Path   string
	Mode   string
	Size   int64
	SHA256 string
}

type packageIdentity struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func mustCanonicalPayload(raw json.RawMessage) []byte {
	var value any
	if json.Unmarshal(raw, &value) != nil {
		return nil
	}
	data, _ := json.Marshal(value)
	return data
}

func digest(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func ReadRegular(path string) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%s is not a regular file", filepath.Base(path))
	}
	return os.ReadFile(path)
}

func ReadPackageVersion(root string) (string, error) {
	data, err := ReadRegular(filepath.Join(root, "package.json"))
	if err != nil {
		return "", err
	}
	var pkg struct {
		Version string `json:"version"`
	}
	if json.Unmarshal(data, &pkg) != nil || pkg.Version == "" {
		return "", errors.New("invalid package version")
	}
	return pkg.Version, nil
}

type platformDefinition struct {
	OS     string `json:"os"`
	Arch   string `json:"arch"`
	GOOS   string `json:"goos"`
	GOArch string `json:"goarch"`
	Runner string `json:"runner"`
}

func inspectArtifacts(root string) ([]artifactEvidence, []targetEvidence, string, error) {
	matrixData, err := ReadRegular(filepath.Join(root, "scripts", "platforms.json"))
	if err != nil {
		return nil, nil, "", fmt.Errorf("platform matrix is unreadable: %w", err)
	}
	var matrix []platformDefinition
	if err := decodeStrict(matrixData, &matrix); err != nil || len(matrix) != 4 {
		return nil, nil, "", errors.New("platform matrix must contain exactly four supported targets")
	}
	seenTargets := map[string]bool{}
	version, err := ReadPackageVersion(root)
	if err != nil {
		return nil, nil, "", fmt.Errorf("package identity is unreadable: %w", err)
	}
	artifactDir := filepath.Join(root, "dist", "artifacts")
	entries, err := os.ReadDir(artifactDir)
	if err != nil {
		return nil, nil, "", fmt.Errorf("artifact directory is unreadable: %w", err)
	}
	want := map[string]string{"redbench-" + version + ".tgz": "wrapper"}
	targets := make([]targetEvidence, 0, len(matrix))
	for _, item := range matrix {
		key := item.OS + "-" + item.Arch
		if (item.OS != "darwin" && item.OS != "linux") || (item.Arch != "arm64" && item.Arch != "x64") || item.Runner == "" || seenTargets[key] {
			return nil, nil, "", fmt.Errorf("platform matrix contains an invalid or duplicate target: %s", key)
		}
		wantGOArch := map[string]string{"arm64": "arm64", "x64": "amd64"}[item.Arch]
		if item.GOOS != item.OS || item.GOArch != wantGOArch {
			return nil, nil, "", fmt.Errorf("platform matrix target %s has inconsistent Go target", key)
		}
		seenTargets[key] = true
		name := fmt.Sprintf("redbench-%s-%s-%s.tgz", item.OS, item.Arch, version)
		want[name] = item.OS + "-" + item.Arch
		want[fmt.Sprintf("redbench-%s-%s-%s.tar.gz", version, item.OS, item.Arch)] = item.OS + "-" + item.Arch
		targets = append(targets, targetEvidence{OS: item.OS, Arch: item.Arch, GOOS: item.GOOS, GOArch: item.GOArch, Runner: item.Runner})
	}
	for _, key := range []string{"darwin-arm64", "darwin-x64", "linux-arm64", "linux-x64"} {
		if !seenTargets[key] {
			return nil, nil, "", fmt.Errorf("platform matrix omits %s", key)
		}
	}
	if len(entries) != len(want) {
		return nil, nil, "", fmt.Errorf("artifact set has %d entries, want %d", len(entries), len(want))
	}
	artifacts := make([]artifactEvidence, 0, len(want))
	packageFiles := map[string]map[string]tarFile{}
	artifactBytes := map[string][]byte{}
	setHash := sha256.New()
	for _, entry := range entries {
		if !entry.Type().IsRegular() || entry.Type()&os.ModeSymlink != 0 {
			return nil, nil, "", fmt.Errorf("artifact input is not a regular file: %s", entry.Name())
		}
		target, ok := want[entry.Name()]
		if !ok {
			return nil, nil, "", fmt.Errorf("unknown artifact: %s", entry.Name())
		}
		artifactPath := filepath.Join(artifactDir, entry.Name())
		info, err := os.Lstat(artifactPath)
		if err != nil {
			return nil, nil, "", fmt.Errorf("artifact is unreadable: %s", entry.Name())
		}
		if info.Size() > maxArchiveCompressedSize {
			return nil, nil, "", fmt.Errorf("artifact %s compressed size exceeds inspection limit", entry.Name())
		}
		data, err := os.ReadFile(artifactPath)
		if err != nil || len(data) == 0 {
			return nil, nil, "", fmt.Errorf("artifact is unreadable or empty: %s", entry.Name())
		}
		fmt.Fprintf(setHash, "%s:%d:", entry.Name(), len(data))
		_, _ = setHash.Write(data)
		_, _ = setHash.Write([]byte{0})
		artifactBytes[entry.Name()] = data
		if strings.HasSuffix(entry.Name(), ".tar.gz") {
			files, err := readOfflineArchive(data, entry.Name(), target, version)
			if err != nil {
				return nil, nil, "", fmt.Errorf("offline archive %s is invalid: %w", entry.Name(), err)
			}
			inventory, err := canonicalArchiveInventory(files)
			if err != nil {
				return nil, nil, "", err
			}
			artifacts = append(artifacts, artifactEvidence{Name: entry.Name(), Target: target, Size: int64(len(data)), SHA256: digest(data), ComponentDigest: digest(files["evidence/components/platform-component-manifest.json"].data), SBOMDigest: digest(files["evidence/governance/sbom.spdx.json"].data), InventoryDigest: digest(inventory)})
			continue
		}
		files, manifest, err := readTarball(data)
		if err != nil {
			return nil, nil, "", fmt.Errorf("artifact %s is invalid: %w", entry.Name(), err)
		}
		if err := validatePackageEvidence(files, manifest, target, version); err != nil {
			return nil, nil, "", fmt.Errorf("artifact %s evidence is invalid: %w", entry.Name(), err)
		}
		inventory, err := canonicalManifestFiles(manifest.Files)
		if err != nil {
			return nil, nil, "", err
		}
		packageFiles[entry.Name()] = files
		artifacts = append(artifacts, artifactEvidence{Name: entry.Name(), Target: target, Size: int64(len(data)), SHA256: digest(data), ComponentDigest: digest(files[requirements.ComponentManifest.Path].data), SBOMDigest: digest(files["governance/sbom.spdx.json"].data), InventoryDigest: digest(inventory)})
	}
	for _, item := range matrix {
		target := item.OS + "-" + item.Arch
		platformName := fmt.Sprintf("redbench-%s-%s-%s.tgz", item.OS, item.Arch, version)
		archiveName := fmt.Sprintf("redbench-%s-%s-%s.tar.gz", version, item.OS, item.Arch)
		archiveFiles, err := readOfflineArchive(artifactBytes[archiveName], archiveName, target, version)
		if err != nil {
			return nil, nil, "", err
		}
		if !bytes.Equal(archiveFiles["packages/redbench-"+version+".tgz"].data, artifactBytes["redbench-"+version+".tgz"]) || !bytes.Equal(archiveFiles["packages/"+platformName].data, artifactBytes[platformName]) {
			return nil, nil, "", fmt.Errorf("offline archive %s does not carry the approved npm tarball bytes", archiveName)
		}
		if !bytes.Equal(archiveFiles["bin/bench"].data, packageFiles[platformName]["bin/bench"].data) {
			return nil, nil, "", fmt.Errorf("offline archive %s binary differs from platform package", archiveName)
		}
		if !bytes.Equal(archiveFiles["evidence/components/wrapper-component-manifest.json"].data, packageFiles["redbench-"+version+".tgz"][requirements.ComponentManifest.Path].data) || !bytes.Equal(archiveFiles["evidence/components/platform-component-manifest.json"].data, packageFiles[platformName][requirements.ComponentManifest.Path].data) {
			return nil, nil, "", fmt.Errorf("offline archive %s component evidence differs from package evidence", archiveName)
		}
	}
	return artifacts, targets, hex.EncodeToString(setHash.Sum(nil)), nil
}

func fingerprintArtifactSet(root string) (string, error) {
	entries, err := os.ReadDir(filepath.Join(root, "dist", "artifacts"))
	if err != nil {
		return "", err
	}
	h := sha256.New()
	for _, entry := range entries {
		if !entry.Type().IsRegular() || entry.Type()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("artifact input is not a regular file: %s", entry.Name())
		}
		data, err := os.ReadFile(filepath.Join(root, "dist", "artifacts", entry.Name()))
		if err != nil {
			return "", err
		}
		fmt.Fprintf(h, "%s:%d:", entry.Name(), len(data))
		_, _ = h.Write(data)
		_, _ = h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
