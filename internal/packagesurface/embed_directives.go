package packagesurface

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// EmbedDirectives are the patterns of the //go:embed directives of a tree, each resolved
// against the directory of its file to a repo-relative slash path.
type EmbedDirectives []string

// ScanEmbedDirectives reads every Go file below root, test files included, and answers the
// patterns of their embed directives. The scan skips the Git entry of root. EmbedTargets
// reads only the build sources, so it is no substitute for a caller that must see each
// directive of the tree.
func ScanEmbedDirectives(root string) (EmbedDirectives, error) {
	var directives EmbedDirectives
	err := filepath.WalkDir(root, func(file string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Name() == ".git" && filepath.Dir(file) == filepath.Clean(root) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !entry.Type().IsRegular() || !strings.HasSuffix(entry.Name(), ".go") {
			return nil
		}
		patterns, err := directivePatterns(root, file)
		directives = append(directives, patterns...)
		return err
	})
	return directives, err
}

// directivePatterns answers the pattern of each embed directive of the Go file at file,
// resolved against the file's directory to a slash path relative to root.
func directivePatterns(root, file string) ([]string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var patterns []string
	dir := filepath.Dir(file)
	for _, line := range strings.Split(string(data), "\n") {
		match := goEmbedDirective.FindStringSubmatch(strings.TrimSpace(line))
		if match == nil {
			continue
		}
		for _, pattern := range strings.Fields(match[1]) {
			// The all: prefix only admits hidden files; the pattern after it names the paths.
			pattern = strings.TrimPrefix(strings.Trim(pattern, "\"`"), "all:")
			rel, err := filepath.Rel(root, filepath.Join(dir, filepath.FromSlash(pattern)))
			if err != nil {
				return nil, err
			}
			patterns = append(patterns, filepath.ToSlash(rel))
		}
	}
	return patterns, nil
}

// Matches reports whether a directive embeds the repo-relative slash path rel: a directive
// equals rel or one of its parent directories, or matches one of them as a glob.
func (d EmbedDirectives) Matches(rel string) bool {
	for p := rel; p != "." && p != "/"; p = path.Dir(p) {
		for _, pattern := range d {
			if matched, _ := path.Match(pattern, p); matched || pattern == p {
				return true
			}
		}
	}
	return false
}
