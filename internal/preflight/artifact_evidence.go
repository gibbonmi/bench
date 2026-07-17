package preflight

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

var exchangeEvidenceDirs = atomicExchangeDirs

type tarFile struct {
	mode int64
	data []byte
}

type componentManifest struct {
	SchemaVersion int               `json:"schema_version"`
	Component     componentIdentity `json:"component"`
	Files         []manifestFile    `json:"files"`
}

type componentIdentity struct {
	Name    string     `json:"name"`
	Version string     `json:"version"`
	Target  targetName `json:"target"`
}

type targetName struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

type manifestFile struct {
	Path   string `json:"path"`
	Mode   string `json:"mode"`
	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`
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

type platformDefinition struct {
	OS     string `json:"os"`
	Arch   string `json:"arch"`
	GOOS   string `json:"goos"`
	GOArch string `json:"goarch"`
	Runner string `json:"runner"`
}

func inspectArtifacts(root string) ([]artifactEvidence, []targetEvidence, error) {
	matrixData, err := readRegular(filepath.Join(root, "scripts", "platforms.json"))
	if err != nil {
		return nil, nil, fmt.Errorf("platform matrix is unreadable: %w", err)
	}
	var matrix []platformDefinition
	if err := decodeStrict(matrixData, &matrix); err != nil || len(matrix) != 4 {
		return nil, nil, errors.New("platform matrix must contain exactly four targets")
	}
	version, err := readPackageVersion(root)
	if err != nil {
		return nil, nil, fmt.Errorf("package identity is unreadable: %w", err)
	}
	artifactDir := filepath.Join(root, "dist", "artifacts")
	entries, err := os.ReadDir(artifactDir)
	if err != nil {
		return nil, nil, fmt.Errorf("artifact directory is unreadable: %w", err)
	}
	want := map[string]string{"redbench-" + version + ".tgz": "wrapper"}
	targets := make([]targetEvidence, 0, len(matrix))
	for _, item := range matrix {
		name := fmt.Sprintf("redbench-%s-%s-%s.tgz", item.OS, item.Arch, version)
		want[name] = item.OS + "-" + item.Arch
		targets = append(targets, targetEvidence{OS: item.OS, Arch: item.Arch, GOOS: item.GOOS, GOArch: item.GOArch, Runner: item.Runner})
	}
	if len(entries) != len(want) {
		return nil, nil, fmt.Errorf("artifact set has %d entries, want %d", len(entries), len(want))
	}
	artifacts := make([]artifactEvidence, 0, len(want))
	for _, entry := range entries {
		if !entry.Type().IsRegular() || entry.Type()&os.ModeSymlink != 0 {
			return nil, nil, fmt.Errorf("artifact input is not a regular file: %s", entry.Name())
		}
		target, ok := want[entry.Name()]
		if !ok {
			return nil, nil, fmt.Errorf("unknown artifact: %s", entry.Name())
		}
		data, err := os.ReadFile(filepath.Join(artifactDir, entry.Name()))
		if err != nil || len(data) == 0 {
			return nil, nil, fmt.Errorf("artifact is unreadable or empty: %s", entry.Name())
		}
		files, manifest, err := readTarball(data)
		if err != nil {
			return nil, nil, fmt.Errorf("artifact %s is invalid: %w", entry.Name(), err)
		}
		if err := validatePackageEvidence(files, manifest, target, version); err != nil {
			return nil, nil, fmt.Errorf("artifact %s evidence is invalid: %w", entry.Name(), err)
		}
		inventory, err := canonicalJSON(manifest.Files)
		if err != nil {
			return nil, nil, err
		}
		artifacts = append(artifacts, artifactEvidence{Name: entry.Name(), Target: target, Size: int64(len(data)), SHA256: digest(data), ComponentDigest: digest(files["component-manifest.json"].data), SBOMDigest: digest(files["governance/sbom.spdx.json"].data), InventoryDigest: digest(inventory)})
	}
	return artifacts, targets, nil
}

func readTarball(data []byte) (map[string]tarFile, componentManifest, error) {
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, componentManifest{}, err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	files := map[string]tarFile{}
	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, componentManifest{}, err
		}
		name, err := archiveRelativePath(header.Name)
		if err != nil {
			return nil, componentManifest{}, err
		}
		if header.Typeflag == tar.TypeDir {
			continue
		}
		if header.Typeflag != tar.TypeReg && header.Typeflag != tar.TypeRegA {
			return nil, componentManifest{}, fmt.Errorf("archive contains special file %s", name)
		}
		mode := header.Mode & 0o777
		if header.Mode&0o7000 != 0 || mode != 0o644 && mode != 0o755 {
			return nil, componentManifest{}, fmt.Errorf("archive contains unsafe mode %o for %s", mode, name)
		}
		if _, exists := files[name]; exists {
			return nil, componentManifest{}, fmt.Errorf("archive contains duplicate path %s", name)
		}
		body, err := io.ReadAll(io.LimitReader(tr, 256<<20))
		if err != nil {
			return nil, componentManifest{}, err
		}
		files[name] = tarFile{mode: header.Mode & 0o777, data: body}
	}
	manifestFile, ok := files["component-manifest.json"]
	if !ok {
		return nil, componentManifest{}, errors.New("component manifest is missing")
	}
	var manifest componentManifest
	if err := decodeStrict(manifestFile.data, &manifest); err != nil {
		return nil, componentManifest{}, fmt.Errorf("component manifest is malformed: %w", err)
	}
	return files, manifest, nil
}

func validatePackageEvidence(files map[string]tarFile, manifest componentManifest, target, version string) error {
	if manifest.SchemaVersion != 1 || manifest.Component.Version != version || manifest.Component.Name == "" {
		return errors.New("component identity or schema is invalid")
	}
	if target == "wrapper" {
		if manifest.Component.Name != "redbench" || manifest.Component.Target.OS != "all" || manifest.Component.Target.Arch != "all" {
			return errors.New("wrapper component identity is invalid")
		}
	} else {
		if manifest.Component.Name != "@redbench/"+target {
			return errors.New("platform component identity is invalid")
		}
		parts := strings.Split(target, "-")
		if len(parts) != 2 || manifest.Component.Target.OS != parts[0] || manifest.Component.Target.Arch != parts[1] {
			return errors.New("platform component target is invalid")
		}
	}
	for _, evidence := range packageEvidenceRegistry() {
		file, ok := files[evidence.Path]
		if !ok {
			return fmt.Errorf("required package evidence is missing: %s", evidence.Path)
		}
		mode, err := strconv.ParseInt(evidence.Mode, 8, 32)
		if err != nil || file.mode != mode {
			return fmt.Errorf("package evidence mode is invalid for %s", evidence.Path)
		}
		if len(file.data) == 0 || !bytes.HasSuffix(file.data, []byte("\n")) {
			return fmt.Errorf("package evidence is empty or missing a final newline: %s", evidence.Path)
		}
		switch evidence.Schema {
		case "license/v1":
		case "notices/v1", "spdx-json/2.3", "governance-policy/v1":
			record, found := requirementForPath(evidence.Path)
			if !found {
				return fmt.Errorf("package evidence has no requirement record: %s", evidence.Path)
			}
			if evidence.Schema == "spdx-json/2.3" {
				if err := validateSPDXDocument(file.data, manifest.Component.Name, version); err != nil {
					return err
				}
			} else if err := validateRequirementBytes(record, file.data, Identity{}); err != nil {
				return err
			}
		default:
			return fmt.Errorf("package evidence has unsupported schema: %s", evidence.Schema)
		}
	}
	pkg, ok := files["package.json"]
	if !ok {
		return errors.New("package identity is missing")
	}
	var identity packageIdentity
	if err := rejectDuplicateJSONKeys(pkg.data); err != nil || json.Unmarshal(pkg.data, &identity) != nil || identity.Version != version {
		return errors.New("package identity is malformed")
	}
	if target == "wrapper" && identity.Name != "redbench" || target != "wrapper" && identity.Name != "@redbench/"+target {
		return errors.New("package name does not match component")
	}
	manifestPaths := map[string]manifestFile{}
	last := ""
	for _, item := range manifest.Files {
		if item.Path == "" || item.Path <= last || strings.Contains(item.Path, "\\") || strings.HasPrefix(item.Path, "../") || filepath.IsAbs(item.Path) {
			return errors.New("component inventory is not sorted or contains an unsafe path")
		}
		last = item.Path
		if _, exists := manifestPaths[item.Path]; exists {
			return errors.New("component inventory contains a duplicate path")
		}
		manifestPaths[item.Path] = item
		actual, exists := files[item.Path]
		if !exists || item.Size != int64(len(actual.data)) || item.Mode != fmt.Sprintf("%o", actual.mode) || item.SHA256 != digest(actual.data) {
			return fmt.Errorf("component inventory disagrees with observed bytes: %s", item.Path)
		}
	}
	if len(manifestPaths) != len(files)-1 {
		return errors.New("component inventory does not enumerate every package file")
	}
	if _, ok := manifestPaths["component-manifest.json"]; ok {
		return errors.New("component inventory self-references its manifest")
	}
	return nil
}

func requirementForPath(path string) (Requirement, bool) {
	for _, record := range requirements.Records {
		if record.Path == path {
			return record, true
		}
	}
	return Requirement{}, false
}

func archiveRelativePath(name string) (string, error) {
	if !strings.HasPrefix(name, "package/") || strings.Contains(name, "\\") || hasControlBytes(name) {
		return "", fmt.Errorf("unsafe archive path %q", name)
	}
	rel := strings.TrimPrefix(name, "package/")
	clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(rel)))
	if rel == "" || clean != rel || clean == "." || strings.HasPrefix(clean, "../") || filepath.IsAbs(clean) {
		return "", fmt.Errorf("unsafe archive path %q", name)
	}
	return clean, nil
}

func PromoteEvidence(root string, mode Mode, results []Result, manifest Manifest) error {
	return PromoteEvidenceFiles(root, mode, results, manifest, nil)
}

func PromoteEvidenceFiles(root string, mode Mode, results []Result, manifest Manifest, files map[string][]byte) error {
	dist := filepath.Join(root, "dist")
	if info, err := os.Lstat(dist); err == nil && (!info.IsDir() || info.Mode()&os.ModeSymlink != 0) {
		return fmt.Errorf("dist output target is not a real directory")
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(dist, 0o755); err != nil {
		return err
	}
	stage, err := os.MkdirTemp(dist, ".preflight-stage-")
	if err != nil {
		return err
	}
	keep := false
	defer func() {
		if !keep {
			_ = os.RemoveAll(stage)
		}
	}()
	for _, result := range results {
		record := Record{SchemaVersion: 1, Phase: result.Name, Mode: mode, Status: result.Status, ExitCode: result.ExitCode, Error: result.Failure}
		if err := writeJSONSync(filepath.Join(stage, result.Name+".json"), record); err != nil {
			return err
		}
	}
	if err := writeJSONSync(filepath.Join(stage, "manifest.json"), manifest); err != nil {
		return err
	}
	for name, data := range files {
		if filepath.Base(name) != name || name == "" {
			return fmt.Errorf("invalid promoted evidence file name: %s", name)
		}
		if err := writeBytesSync(filepath.Join(stage, name), data); err != nil {
			return err
		}
	}
	if err := syncDir(stage); err != nil {
		return err
	}
	target := filepath.Join(dist, "preflight")
	if info, err := os.Lstat(target); err == nil && (!info.IsDir() || info.Mode()&os.ModeSymlink != 0) {
		return fmt.Errorf("preflight output target is not a real directory")
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	if _, err := os.Lstat(target); os.IsNotExist(err) {
		if err := os.Rename(stage, target); err != nil {
			return err
		}
		keep = true
	} else if err != nil {
		return err
	} else {
		if err := exchangeEvidenceDirs(stage, target); err != nil {
			return err
		}
		keep = true
	}
	if err := syncDir(dist); err != nil {
		return err
	}
	if err := os.RemoveAll(stage); err != nil {
		return err
	}
	return nil
}

func atomicExchangeDirs(left, right string) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("atomic evidence replacement is unsupported on %s", runtime.GOOS)
	}
	var trap uintptr
	switch runtime.GOARCH {
	case "amd64":
		trap = 316
	case "arm64":
		trap = 276
	default:
		return fmt.Errorf("atomic evidence replacement is unsupported on linux/%s", runtime.GOARCH)
	}
	l, err := syscall.BytePtrFromString(left)
	if err != nil {
		return err
	}
	r, err := syscall.BytePtrFromString(right)
	if err != nil {
		return err
	}
	atFDCWD := ^uintptr(99)
	_, _, errno := syscall.Syscall6(trap, atFDCWD, uintptr(unsafe.Pointer(l)), atFDCWD, uintptr(unsafe.Pointer(r)), 2, 0)
	if errno != 0 {
		return errno
	}
	return nil
}

func writeJSONSync(path string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return writeBytesSync(path, append(data, '\n'))
}

func writeBytesSync(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	if _, err = f.Write(data); err == nil {
		err = f.Sync()
	}
	closeErr := f.Close()
	if err != nil {
		return err
	}
	return closeErr
}

func syncDir(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := f.Sync(); err != nil {
		return fmt.Errorf("sync %s: %w", path, err)
	}
	return nil
}
