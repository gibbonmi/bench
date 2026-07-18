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

func readOfflineArchive(data []byte, archiveName, target, version string) (map[string]tarFile, error) {
	if int64(len(data)) > maxArchiveCompressedSize {
		return nil, errors.New("offline archive compressed size exceeds inspection limit")
	}
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	gz.Multistream(false)
	defer gz.Close()
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
	if err := validateOfflineArchiveDirs(dirs, target, version); err != nil {
		return nil, err
	}
	if err := validateOfflineArchiveFiles(files, target, version); err != nil {
		return nil, err
	}
	return files, nil
}

func validateOfflineArchiveDirs(dirs map[string]bool, target, version string) error {
	files := map[string]bool{}
	wrapper := "redbench-" + version + ".tgz"
	platform := "redbench-" + target + "-" + version + ".tgz"
	for _, name := range []string{
		"bin/bench",
		"packages/" + wrapper,
		"packages/" + platform,
		"OFFLINE.md",
		"evidence/components/wrapper-component-manifest.json",
		"evidence/components/platform-component-manifest.json",
	} {
		files[name] = true
	}
	for _, evidence := range PackageEvidenceRegistry() {
		files["evidence/"+evidence.Path] = true
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

func validateOfflineArchiveFiles(files map[string]tarFile, target, version string) error {
	wrapper := "redbench-" + version + ".tgz"
	platform := "redbench-" + target + "-" + version + ".tgz"
	want := map[string]int64{
		"bin/bench":            0o755,
		"packages/" + wrapper:  0o644,
		"packages/" + platform: 0o644,
		"OFFLINE.md":           0o644,
		"evidence/components/wrapper-component-manifest.json":  0o644,
		"evidence/components/platform-component-manifest.json": 0o644,
	}
	for _, evidence := range PackageEvidenceRegistry() {
		want["evidence/"+evidence.Path] = 0o644
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
	if !bytes.Contains(files["OFFLINE.md"].data, []byte("npm --offline")) || !bytes.Contains(files["OFFLINE.md"].data, []byte("platform tarball first")) {
		return errors.New("offline archive instructions are incomplete")
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
