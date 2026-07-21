package releaseevidence

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"
)

func readTarball(data []byte) (map[string]tarFile, componentManifest, error) {
	files, _, err := readSecureArchive(data, func(header *tar.Header) (string, error) {
		return archiveRelativePath(header.Name)
	})
	if err != nil {
		return nil, componentManifest{}, err
	}
	manifestFile, ok := files[requirements.ComponentManifest.Path]
	if !ok {
		return nil, componentManifest{}, errors.New("component manifest is missing")
	}
	manifest, err := decodeComponentManifest(manifestFile.data)
	if err != nil {
		return nil, componentManifest{}, fmt.Errorf("component manifest is malformed: %w", err)
	}
	return files, manifest, nil
}

func readSecureArchive(data []byte, nameFor func(*tar.Header) (string, error)) (map[string]tarFile, map[string]bool, error) {
	if int64(len(data)) > maxArchiveCompressedSize {
		return nil, nil, errors.New("archive compressed size exceeds inspection limit")
	}
	source := bytes.NewReader(data)
	gz, err := gzip.NewReader(source)
	if err != nil {
		return nil, nil, err
	}
	gz.Multistream(false)
	reader := tar.NewReader(gz)
	files, dirs := map[string]tarFile{}, map[string]bool{}
	var expanded int64
	for members := 1; ; members++ {
		header, err := reader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		if members > maxArchiveMemberCount {
			return nil, nil, errors.New("archive member count exceeds inspection limit")
		}
		if header.Size < 0 || header.Size > maxArchiveExpandedSize-expanded {
			return nil, nil, errors.New("archive expanded size exceeds inspection limit")
		}
		expanded += header.Size
		name, err := nameFor(header)
		if err != nil {
			return nil, nil, err
		}
		mode := header.Mode & 0o777
		if header.Typeflag == tar.TypeDir {
			if mode != 0o755 || header.Mode&0o7000 != 0 || dirs[name] {
				return nil, nil, fmt.Errorf("archive contains duplicate or unsafe directory %s", name)
			}
			dirs[name] = true
			continue
		}
		if header.Typeflag != tar.TypeReg && header.Typeflag != tar.TypeRegA {
			return nil, nil, fmt.Errorf("archive contains special file %s", name)
		}
		if header.Mode&0o7000 != 0 || mode != 0o644 && mode != 0o755 {
			return nil, nil, fmt.Errorf("archive contains unsafe mode %o for %s", mode, name)
		}
		if _, exists := files[name]; exists {
			return nil, nil, fmt.Errorf("archive contains duplicate path %s", name)
		}
		body, err := io.ReadAll(io.LimitReader(reader, maxArchiveMemberSize+1))
		if err != nil {
			return nil, nil, err
		}
		if int64(len(body)) > maxArchiveMemberSize {
			return nil, nil, fmt.Errorf("archive member %s exceeds inspection limit", name)
		}
		if len(body) == 0 {
			return nil, nil, fmt.Errorf("archive contains empty member %s", name)
		}
		files[name] = tarFile{mode: mode, data: body}
	}
	if err := rejectGzipSuffix(gz, source); err != nil {
		return nil, nil, err
	}
	return files, dirs, nil
}

func rejectGzipSuffix(gz *gzip.Reader, source *bytes.Reader) error {
	if _, err := io.Copy(io.Discard, gz); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}
	if source.Len() == 0 {
		return nil
	}
	if err := gz.Reset(source); err == nil {
		return errors.New("archive has concatenated gzip members")
	} else if errors.Is(err, io.EOF) {
		return nil
	} else {
		return fmt.Errorf("archive has invalid trailing gzip data: %w", err)
	}
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
	for _, evidence := range PackageEvidenceRegistry() {
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
	if _, ok := manifestPaths[requirements.ComponentManifest.Path]; ok {
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
