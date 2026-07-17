package preflight

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func inputFingerprint(root string, run RunEvidence) (string, error) {
	h := sha256.New()
	paths := append([]string{"internal/preflight/requirements.json"}, releaseInputPaths...)
	for _, record := range requirements.Records {
		paths = append(paths, record.Path)
	}
	paths = append(paths, "dist/artifacts")
	sort.Strings(paths)
	seen := map[string]bool{}
	for _, rel := range paths {
		if seen[rel] {
			continue
		}
		seen[rel] = true
		if err := fingerprintPath(h, root, rel); err != nil {
			return "", err
		}
	}
	for _, result := range run.Phases {
		_, _ = h.Write([]byte(result.Name))
		_, _ = h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func fingerprintPath(h io.Writer, root, rel string) error {
	path := filepath.Join(root, filepath.FromSlash(rel))
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		_, _ = io.WriteString(h, "absent:\x00"+rel+"\n")
		return nil
	}
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 || (!info.Mode().IsRegular() && !info.IsDir()) {
		return fmt.Errorf("unsafe release evidence input: %s", rel)
	}
	if info.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		names := make([]string, 0, len(entries))
		for _, entry := range entries {
			names = append(names, entry.Name())
		}
		sort.Strings(names)
		for _, name := range names {
			if err := fingerprintPath(h, root, filepath.ToSlash(filepath.Join(rel, name))); err != nil {
				return err
			}
		}
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, _ = io.WriteString(h, fmt.Sprintf("file:%s:%o:%d:", rel, info.Mode().Perm(), len(data)))
	_, _ = h.Write(data)
	_, _ = h.Write([]byte{0})
	return nil
}
