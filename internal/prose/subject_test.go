package prose

import (
	"os"
	"strings"
	"testing"
)

// TestGradeNamedHonorsExclusions is OG08: a named Markdown path listed in
// .bench/prose-exclusions produces no finding, even though its content is over the
// sentence bound.
func TestGradeNamedHonorsExclusions(t *testing.T) {
	root := t.TempDir()
	write(t, root, "planted/notes.md", longSentence())
	write(t, root, ".bench/prose-exclusions", "planted/notes.md fixture carries planted content\n")

	if got := GradeNamed(root, []string{"planted/notes.md"}); got != nil {
		t.Fatalf("GradeNamed on an excluded named path = %v, want no finding", got)
	}
}

// TestGradeNamedGradesAnUnexcludedSubject shows the exclusion row in the fixture above
// narrows the grade rather than emptying it: an un-excluded subject with the same fault
// still reports.
func TestGradeNamedGradesAnUnexcludedSubject(t *testing.T) {
	root := t.TempDir()
	write(t, root, "planted/other.md", longSentence())
	write(t, root, ".bench/prose-exclusions", "planted/notes.md fixture carries planted content\n")

	if got := GradeNamed(root, []string{"planted/other.md"}); len(got) == 0 {
		t.Fatalf("GradeNamed on an un-excluded named path = %v, want a finding", got)
	}
}

// TestGradeNamedSkipsAbsentAndRefusesSymlinkPaths covers the edge inventory: a named path
// that does not exist is skipped, and a symbolic link is refused rather than followed.
func TestGradeNamedSkipsAbsentAndRefusesSymlinkPaths(t *testing.T) {
	root := t.TempDir()
	write(t, root, ".bench/prose-exclusions", "")
	write(t, root, "real.md", longSentence())
	link := root + "/link.md"
	if err := os.Symlink(root+"/real.md", link); err != nil {
		t.Fatalf("symlink fixture: %v", err)
	}

	if got := GradeNamed(root, []string{"missing.md"}); got != nil {
		t.Fatalf("GradeNamed on an absent path = %v, want no finding", got)
	}
	got := GradeNamed(root, []string{"link.md"})
	if len(got) != 1 || !strings.HasPrefix(got[0], `prose: "link.md": refused subject:`) {
		t.Fatalf("GradeNamed on a symlink = %v, want one refused-subject line", got)
	}
}
