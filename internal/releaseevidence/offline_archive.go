package releaseevidence

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

func readOfflineArchive(rootPath string, data []byte, archiveName, target, version string) (map[string]tarFile, error) {
	root := strings.TrimSuffix(archiveName, ".tar.gz") + "/"
	files, dirs, err := readSecureArchive(data, func(header *tar.Header) (string, error) {
		if header.Name == root {
			return "", nil
		}
		if !strings.HasPrefix(header.Name, root) || hasControlBytes(header.Name) || strings.Contains(header.Name, "\\") {
			return "", fmt.Errorf("unsafe offline archive path %q", header.Name)
		}
		rel := strings.TrimSuffix(strings.TrimPrefix(header.Name, root), "/")
		clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(rel)))
		if rel == "" || clean != rel || clean == "." || strings.HasPrefix(clean, "../") || filepath.IsAbs(clean) || header.Typeflag == tar.TypeDir && !strings.HasSuffix(header.Name, "/") {
			return "", fmt.Errorf("unsafe offline archive path %q", header.Name)
		}
		return clean, nil
	})
	if err != nil {
		return nil, err
	}
	if !dirs[""] {
		return nil, fmt.Errorf("offline archive root is missing: %s", root)
	}
	if err := validateOfflineArchiveDirs(rootPath, dirs, target, version); err != nil {
		return nil, err
	}
	if err := validateOfflineArchiveFiles(rootPath, files, target, version); err != nil {
		return nil, err
	}
	return files, nil
}

func validateOfflineArchiveDirs(root string, dirs map[string]bool, target, version string) error {
	files, err := archiveInventory(root, target, version)
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

func validateOfflineArchiveFiles(root string, files map[string]tarFile, target, version string) error {
	want, err := archiveInventory(root, target, version)
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
	if !bytes.Contains(files["OFFLINE.md"].data, []byte("--offline")) || !bytes.Contains(files["OFFLINE.md"].data, []byte("SHA256SUMS | sha256sum -c -")) || bytes.Contains(files["OFFLINE.md"].data, []byte("sha256sum -c SHA256SUMS")) || !bytes.Contains(files["OFFLINE.md"].data, []byte("npm publish ./packages/")) {
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
