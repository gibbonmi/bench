package canary

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type fixtureMutation struct {
	Path string `json:"path"`
	Old  string `json:"old"`
	New  string `json:"new"`
}

func regularFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

func materializeMutationFixture(root, fixture, dst string) error {
	basePath := filepath.Join(fixture, "BASE")
	if regularFile(basePath) {
		rels, err := basePaths(root, basePath)
		if err != nil {
			return err
		}
		for _, rel := range rels {
			if err := copyBaseFile(filepath.Join(root, rel), filepath.Join(dst, rel)); err != nil {
				return err
			}
		}
	}
	filesDir := filepath.Join(fixture, "files")
	if info, err := os.Stat(filesDir); err == nil && info.IsDir() {
		if err := materialize(filesDir, dst); err != nil {
			return err
		}
	}
	mutatePath := filepath.Join(fixture, "MUTATE.json")
	if !regularFile(mutatePath) {
		return nil
	}
	data, err := os.ReadFile(mutatePath)
	if err != nil {
		return err
	}
	var mutations []fixtureMutation
	if err := json.Unmarshal(data, &mutations); err != nil {
		return err
	}
	for _, mutation := range mutations {
		path := filepath.Join(dst, filepath.Clean(mutation.Path))
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.Count(string(body), mutation.Old) != 1 {
			return fmt.Errorf("mutation anchor in %s did not occur exactly once", mutation.Path)
		}
		body = []byte(strings.Replace(string(body), mutation.Old, mutation.New, 1))
		if err := os.WriteFile(path, body, 0); err != nil {
			return err
		}
	}
	return nil
}

func basePaths(root, basePath string) ([]string, error) {
	seenIncludes := map[string]bool{}
	paths := []string{}
	var walk func([]byte) error
	walk = func(data []byte) error {
		for _, rel := range strings.Fields(string(data)) {
			if strings.HasPrefix(rel, "@") {
				include := strings.TrimPrefix(rel, "@")
				if !safeBasePath(include) {
					return fmt.Errorf("invalid BASE include %q", include)
				}
				if seenIncludes[include] {
					return fmt.Errorf("recursive BASE include %q", include)
				}
				seenIncludes[include] = true
				included, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(include)))
				if err != nil {
					return err
				}
				if err := walk(included); err != nil {
					return err
				}
				continue
			}
			if !safeBasePath(rel) {
				return fmt.Errorf("invalid BASE path %q", rel)
			}
			paths = append(paths, rel)
		}
		return nil
	}
	data, err := os.ReadFile(basePath)
	if err != nil {
		return nil, err
	}
	if err := walk(data); err != nil {
		return nil, err
	}
	return paths, nil
}

func safeBasePath(value string) bool {
	if value == "" || strings.Contains(value, "\\") || strings.Contains(value, "\x00") {
		return false
	}
	clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(value)))
	return clean == value && clean != "." && clean != ".." && !strings.HasPrefix(clean, "../") && !filepath.IsAbs(filepath.FromSlash(value))
}

func copyBaseFile(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("BASE input is not a regular file: %s", src)
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dst, data, info.Mode().Perm())
}
