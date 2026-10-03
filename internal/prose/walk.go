package prose

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/outline"
)

// skippedDirs are the directory names the walk never enters. A fixture tree, a
// dependency tree, and a build output hold planted or generated text that no author owns.
var skippedDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
	"dist":         true,
	"testdata":     true,
}

// IsSubjectName reports whether a file name is a prose subject. The rule is the `.md`
// suffix, so a comment inside a `.go` or shell file is outside the grade. The walk keys on
// it for a name it finds on disk, and a staged caller keys on it for an index entry path,
// so the two selections cannot disagree about what prose grades.
func IsSubjectName(name string) bool { return strings.HasSuffix(name, ".md") }

// Grade walks root and returns one diagnostic for each fault it finds. It returns no
// diagnostic for a clean root and none at all for a root that holds no `*.md` subject:
// a tree with nothing to grade cannot fail an exclusion rule it has no use for.
//
// Grade composes the same per-subject Grader that GradeNamed uses, so the exclusion list,
// the byte classifier, and the finding render have one source.
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
//
// In a git work tree a subject is graded only when git tracks it, so an ignored or
// untracked file that no commit carries is outside the grade. A root that git cannot list
// is graded whole: it has no tracked set to select by.
func collect(root string) ([]string, []string) {
	var tracked map[string]bool
	if files, err := outline.TrackedFiles(root); err == nil {
		tracked = make(map[string]bool, len(files))
		for _, rel := range files {
			tracked[rel] = true
		}
	}
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
		if !IsSubjectName(d.Name()) {
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
		rel = filepath.ToSlash(rel)
		if tracked != nil && !tracked[rel] {
			return nil
		}
		subjects = append(subjects, rel)
		return nil
	})
	if err != nil {
		return nil, []string{fmt.Sprintf("prose: %q: the walk of the graded root failed: %s", root, err)}
	}
	return subjects, nil
}
