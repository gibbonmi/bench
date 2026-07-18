package releaseevidence

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
)

func readOfflineArchive(rootPath string, data []byte, archiveName, target, version string) (map[string]tarFile, error) {
	if int64(len(data)) > maxArchiveCompressedSize {
		return nil, errors.New("offline archive compressed size exceeds inspection limit")
	}
	source := bytes.NewReader(data)
	gz, err := gzip.NewReader(source)
	if err != nil {
		return nil, err
	}
	gz.Multistream(false)
	tr := tar.NewReader(gz)
	root := strings.TrimSuffix(archiveName, ".tar.gz") + "/"
	files := map[string]tarFile{}
	dirs := map[string]bool{"": true}
	rootSeen := false
	var expandedSize int64
	memberCount := 0
	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		memberCount++
		if memberCount > maxArchiveMemberCount {
			return nil, errors.New("offline archive has too many members")
		}
		if header.Size < 0 || header.Size > maxArchiveExpandedSize-expandedSize {
			return nil, errors.New("offline archive expanded size exceeds inspection limit")
		}
		if header.Name == root {
			if rootSeen || header.Typeflag != tar.TypeDir || header.Mode&0o777 != 0o755 {
				return nil, errors.New("offline archive root is duplicated or has an invalid mode")
			}
			rootSeen = true
			continue
		}
		if !strings.HasPrefix(header.Name, root) || hasControlBytes(header.Name) || strings.Contains(header.Name, "\\") {
			return nil, fmt.Errorf("unsafe offline archive path %q", header.Name)
		}
		rel := strings.TrimPrefix(header.Name, root)
		if rel == "" {
			return nil, errors.New("offline archive has an invalid root")
		}
		cleanRel := strings.TrimSuffix(rel, "/")
		clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(cleanRel)))
		if clean != cleanRel || clean == "." || strings.HasPrefix(clean, "../") || filepath.IsAbs(clean) {
			return nil, fmt.Errorf("unsafe offline archive path %q", header.Name)
		}
		if header.Typeflag == tar.TypeDir {
			if header.Mode&0o777 != 0o755 || !strings.HasSuffix(header.Name, "/") {
				return nil, fmt.Errorf("offline archive directory mode is invalid: %s", header.Name)
			}
			if dirs[cleanRel] {
				return nil, fmt.Errorf("offline archive contains duplicate directory %s", cleanRel)
			}
			dirs[cleanRel] = true
			continue
		}
		if header.Typeflag != tar.TypeReg && header.Typeflag != tar.TypeRegA {
			return nil, fmt.Errorf("offline archive contains special file %s", clean)
		}
		mode := header.Mode & 0o777
		if mode != 0o644 && mode != 0o755 || header.Mode&0o7000 != 0 {
			return nil, fmt.Errorf("offline archive contains unsafe mode %o for %s", mode, clean)
		}
		if _, exists := files[clean]; exists {
			return nil, fmt.Errorf("offline archive contains duplicate path %s", clean)
		}
		body, err := io.ReadAll(io.LimitReader(tr, maxArchiveMemberSize+1))
		if err != nil {
			return nil, err
		}
		if int64(len(body)) > maxArchiveMemberSize {
			return nil, fmt.Errorf("offline archive member %s exceeds inspection limit", clean)
		}
		expandedSize += int64(len(body))
		if expandedSize > maxArchiveExpandedSize {
			return nil, errors.New("offline archive expanded size exceeds inspection limit")
		}
		files[clean] = tarFile{mode: mode, data: body}
	}
	if !rootSeen {
		return nil, fmt.Errorf("offline archive root is missing: %s", root)
	}
	if err := rejectGzipSuffix(gz, source); err != nil {
		return nil, err
	}
	plan, err := readReleasePlan(rootPath)
	if err != nil {
		return nil, err
	}
	if err := validateOfflineArchiveDirs(dirs, plan, target, version); err != nil {
		return nil, err
	}
	if err := validateOfflineArchiveFiles(files, plan, target, version); err != nil {
		return nil, err
	}
	return files, nil
}

func validateOfflineArchiveDirs(dirs map[string]bool, plan releasePlan, target, version string) error {
	files, err := archiveInventory(plan, target, version)
	if err != nil {
		return err
	}
	want := map[string]bool{"": true}
	for name := range files {
		for prefix := name; ; {
			if slash := strings.LastIndexByte(prefix, '/'); slash >= 0 {
				prefix = prefix[:slash]
			} else {
				break
			}
			want[prefix] = true
		}
	}
	if len(dirs) != len(want) {
		return fmt.Errorf("offline archive directory inventory has %d entries, want %d", len(dirs), len(want))
	}
	for name := range dirs {
		if !want[name] {
			return fmt.Errorf("offline archive contains unexpected directory %s", name)
		}
	}
	return nil
}

func validateOfflineArchiveFiles(files map[string]tarFile, plan releasePlan, target, version string) error {
	want, err := archiveInventory(plan, target, version)
	if err != nil {
		return err
	}
	if len(files) != len(want) {
		return fmt.Errorf("offline archive inventory has %d files, want %d", len(files), len(want))
	}
	for name, mode := range want {
		file, ok := files[name]
		if !ok {
			return fmt.Errorf("offline archive is missing %s", name)
		}
		if file.mode != mode || len(file.data) == 0 {
			return fmt.Errorf("offline archive file %s has invalid mode or empty bytes", name)
		}
		if !bytes.HasSuffix(file.data, []byte("\n")) && (name == "OFFLINE.md" || strings.HasPrefix(name, "evidence/")) {
			return fmt.Errorf("offline archive file %s is missing a final newline", name)
		}
	}
	if !bytes.Contains(files["OFFLINE.md"].data, []byte("--offline")) || !bytes.Contains(files["OFFLINE.md"].data, []byte("sha256sum -c SHA256SUMS")) || !bytes.Contains(files["OFFLINE.md"].data, []byte("npm publish ./packages/")) {
		return errors.New("offline archive instructions are incomplete")
	}
	return validateArchiveManifest(files, target, version)
}

func validateArchiveManifest(files map[string]tarFile, target, version string) error {
	manifestFile := files["evidence/component-manifest.json"]
	manifest, err := decodeComponentManifest(manifestFile.data)
	if err != nil {
		return fmt.Errorf("offline archive component manifest is malformed: %w", err)
	}
	parts := strings.Split(target, "-")
	if len(parts) != 2 || manifest.SchemaVersion != 1 || manifest.Component.Name != "redbench-offline-"+target || manifest.Component.Version != version || manifest.Component.Target.OS != parts[0] || manifest.Component.Target.Arch != parts[1] {
		return errors.New("offline archive component manifest identity is invalid")
	}
	if len(manifest.Files) != len(files)-1 {
		return errors.New("offline archive component manifest does not enumerate every internal file")
	}
	last := ""
	for _, item := range manifest.Files {
		file, ok := files[item.Path]
		if !ok || item.Path <= last || item.Path == "evidence/component-manifest.json" || item.Mode != fmt.Sprintf("%o", file.mode) || item.Size != int64(len(file.data)) || item.SHA256 != digest(file.data) {
			return fmt.Errorf("offline archive component manifest disagrees with %s", item.Path)
		}
		last = item.Path
	}
	return nil
}

func canonicalArchiveInventory(files map[string]tarFile) ([]byte, error) {
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	items := make([]map[string]any, 0, len(names))
	for _, name := range names {
		file := files[name]
		items = append(items, map[string]any{"path": name, "mode": fmt.Sprintf("%o", file.mode), "size": len(file.data), "sha256": digest(file.data)})
	}
	return canonicalJSON(items)
}
