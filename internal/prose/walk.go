package prose

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// skippedDirs are the directory names the walk never enters. A fixture tree, a
// dependency tree, and a build output hold planted or generated text that no author owns.
var skippedDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
	"dist":         true,
	"testdata":     true,
}

// Grade walks root and returns one diagnostic for each fault it finds. It returns no
// diagnostic for a clean root and none at all for a root that holds no `*.md` subject:
// a tree with nothing to grade cannot fail an exclusion rule it has no use for.
//
// The walk keys on the `*.md` name, so a comment inside a `.go` or shell file is outside
// the grade. Grade composes the same per-subject Grader that GradeNamed uses, so the
// exclusion list, the byte classifier, and the finding render have one source.
func Grade(root string) []string {
	subjects, diags := collect(root)
	if len(diags) > 0 {
		return diags
	}
	if len(subjects) == 0 {
		return nil
	}
	g, exDiags := NewGrader(root)
	if len(exDiags) > 0 {
		return exDiags
	}
	var out []string
	for _, rel := range subjects {
		out = append(out, g.GradeSubject(rel)...)
	}
	return out
}

// collect returns every graded subject under root as a repository-relative slash path.
// A link to a directory is not descended and not reported, because a linked tree is
// graded where it lives.
func collect(root string) ([]string, []string) {
	var subjects []string
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if p != root && skippedDirs[d.Name()] {
				return fs.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		if d.Type()&fs.ModeSymlink != 0 {
			if info, statErr := os.Stat(p); statErr == nil && info.IsDir() {
				return nil
			}
		}
		rel, relErr := filepath.Rel(root, p)
		if relErr != nil {
			return relErr
		}
		subjects = append(subjects, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		return nil, []string{fmt.Sprintf("prose: %q: the walk of the graded root failed: %s", root, err)}
	}
	return subjects, nil
}
