package packagesurface

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
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
		words, err := embedPatterns(match[1])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", file, err)
		}
		for _, pattern := range words {
			// The all: prefix only admits hidden files; the pattern after it names the paths.
			pattern = strings.TrimPrefix(pattern, "all:")
			rel, err := filepath.Rel(root, filepath.Join(dir, filepath.FromSlash(pattern)))
			if err != nil {
				return nil, err
			}
			patterns = append(patterns, filepath.ToSlash(rel))
		}
	}
	return patterns, nil
}

// embedPatterns splits the arguments of one embed directive as Go reads them. A pattern ends
// at a space, unless it is a double-quoted or back-quoted Go string, which can hold one.
func embedPatterns(args string) ([]string, error) {
	var patterns []string
	for args = strings.TrimSpace(args); args != ""; args = strings.TrimSpace(args) {
		var pattern string
		switch args[0] {
		case '`':
			quoted, rest, closed := strings.Cut(args[1:], "`")
			if !closed {
				return nil, invalidEmbedQuote(args)
			}
			pattern, args = quoted, rest
		case '"':
			end := 1
			for end < len(args) && args[end] != '"' {
				if args[end] == '\\' {
					end++
				}
				end++
			}
			if end >= len(args) {
				return nil, invalidEmbedQuote(args)
			}
			unquoted, err := strconv.Unquote(args[:end+1])
			if err != nil {
				return nil, invalidEmbedQuote(args[:end+1])
			}
			pattern, args = unquoted, args[end+1:]
		default:
			end := strings.IndexFunc(args, unicode.IsSpace)
			if end < 0 {
				end = len(args)
			}
			pattern, args = args[:end], args[end:]
		}
		if args != "" && strings.TrimLeftFunc(args, unicode.IsSpace) == args {
			return nil, invalidEmbedQuote(args)
		}
		patterns = append(patterns, pattern)
	}
	return patterns, nil
}

// invalidEmbedQuote is the refusal of directive text that Go does not read as a pattern.
func invalidEmbedQuote(text string) error {
	return fmt.Errorf("invalid quoted string in //go:embed: %s", text)
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
