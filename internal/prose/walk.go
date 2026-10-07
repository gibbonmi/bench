package prose

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
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
func Grade(root string) []string { return GradeTree(root).Findings }

// TreeGrade is the answer of one whole-tree grade. Subjects holds each collected subject
// that the exclusion set does not exclude, in walk order. A walk or exclusion refusal
// grades no subject, so its diagnostics are the Findings of an empty Subjects list.
type TreeGrade struct {
	Subjects []string
	Findings []string
}

// GradeTree walks root as Grade does and also answers the graded subjects.
func GradeTree(root string) TreeGrade {
	collected, diags := collect(root)
	if len(diags) > 0 {
		return TreeGrade{Findings: diags}
	}
	if len(collected) == 0 {
		return TreeGrade{}
	}
	g, exDiags := NewGrader(root)
	if len(exDiags) > 0 {
		return TreeGrade{Findings: exDiags}
	}
	var out TreeGrade
	for _, rel := range collected {
		if g.ex.excluded(rel) {
			continue
		}
		out.Subjects = append(out.Subjects, rel)
		out.Findings = append(out.Findings, g.GradeSubject(rel)...)
	}
	return out
}

// collect returns every graded subject under root as a repository-relative slash path.
// A link to a directory is not descended and not reported, because a linked tree is
// graded where it lives.
//
// When root is the top of its own git work tree, the walk grades only a subject that the
// git index lists, so an ignored or untracked file is outside the grade. The walk grades
// any other root whole: a root that git cannot list, or a root below the top of an outer
// work tree, has no tracked set of its own.
func collect(root string) ([]string, []string) {
	var tracked map[string]bool
	if files, err := outline.TrackedFiles(root); err == nil && workTreeTop(root) {
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

// workTreeTop reports whether root is the top of its own git work tree. Git prints an
// empty prefix only at the top, so the check needs no path comparison.
func workTreeTop(root string) bool {
	prefix, err := git.Output("-C", root, "rev-parse", "--show-prefix")
	return err == nil && prefix == ""
}
