package evidencecmd_test

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// The per-path membership rule in this file is stated independently of the review
// collector, so a changed descriptor or a joined diff source turns a case red.

// fileMember is one file patch source of a published manifest: its pack-local identifier,
// its source digest, and its ordered page digests.
type fileMember struct {
	id, sha string
	pages   []string
}

// fileMembers maps each declared file path to its ordered patch sources, read from the
// published manifest through the strict reader.
func fileMembers(t *testing.T, root, identity string) map[string][]fileMember {
	t.Helper()
	manifest := preflighttest.IdentifiedPack(t, root, identity).Manifest()
	members := map[string][]fileMember{}
	for _, source := range manifest.Sources {
		if source.Role != preflighttest.DiffFileRole {
			continue
		}
		member := fileMember{id: source.ID, sha: source.SHA256}
		for _, page := range manifest.Pages {
			if page.Source == source.ID {
				member.pages = append(member.pages, page.SHA256)
			}
		}
		members[source.Path] = append(members[source.Path], member)
	}
	if len(members) == 0 {
		t.Fatalf("manifest declares no %s source", preflighttest.DiffFileRole)
	}
	return members
}

// sameDigests reports whether two members hold equal source and page digests.
func sameDigests(left, right []fileMember) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i].sha != right[i].sha || !slices.Equal(left[i].pages, right[i].pages) {
			return false
		}
	}
	return true
}

// commitReview writes files, commits every change, and repins the review tip to the commit.
func commitReview(t *testing.T, args []string, files map[string]string, message string) {
	t.Helper()
	for path, body := range files {
		preflighttest.MustWriteFile(t, path, body)
	}
	preflighttest.RunGit(t, "add", "-A")
	preflighttest.RunGit(t, "commit", "-q", "-m", message)
	args[6] = preflighttest.RunGit(t, "rev-parse", "HEAD")
}

// goSource is a valid Go file in package notes with lines comment lines of text.
func goSource(text string, lines int) string {
	return "package notes\n\n" + strings.Repeat("// "+text+"\n", lines)
}

// preparedMembers prepares the review over args and returns its file membership.
func preparedMembers(t *testing.T, root string, args []string) (string, map[string][]fileMember) {
	t.Helper()
	identity, _, _ := prepareEvidence(t, args)
	return identity, fileMembers(t, root, identity)
}

// selectedPatch reads every source of one declared path through its own source stream, in
// manifest order, and returns the joined bytes.
func selectedPatch(t *testing.T, identity string, members []fileMember) string {
	t.Helper()
	var body strings.Builder
	for _, member := range members {
		for _, page := range traverseSource(t, identity, member.id) {
			if page.row["source"] != member.id {
				t.Fatalf("source %s stream returned a page of %v", member.id, page.row["source"])
			}
			body.WriteString(page.row["content"].(string))
		}
	}
	return body.String()
}

// TestReviewFilePageStability is RE1 and RE2. Growing an earlier patch across a page boundary
// leaves every page digest of a later patch unchanged. A separate spec, ticket, or review
// record edit leaves every untouched code patch unchanged, and it changes its own patch.
func TestReviewFilePageStability(t *testing.T) {
	codePaths := []string{"notes/a.go", "notes/z.go", "target/target.go", "edited/user.go"}
	seed := func(t *testing.T) (string, []string) {
		root, _, args := preflighttest.SeedReviewEvidence(t, false)
		commitReview(t, args, map[string]string{
			"notes/a.go": goSource("first patch", 10),
			"notes/z.go": goSource("later patch that spans pages", 600),
		}, "code patches")
		return root, args
	}
	t.Run("earlier patch crosses a page boundary", func(t *testing.T) {
		root, args := seed(t)
		_, before := preparedMembers(t, root, args)
		commitReview(t, args, map[string]string{"notes/a.go": goSource("first patch", 900)}, "grow a.go")
		_, after := preparedMembers(t, root, args)
		if len(before["notes/a.go"]) != 1 || len(before["notes/a.go"][0].pages) != 1 || len(after["notes/a.go"][0].pages) < 2 {
			t.Fatalf("a.go pages before=%v after=%v, want one page growing past a page boundary", before["notes/a.go"], after["notes/a.go"])
		}
		if len(after["notes/z.go"]) != 1 || len(after["notes/z.go"][0].pages) < 2 {
			t.Fatalf("z.go members = %v, want one multi-page patch", after["notes/z.go"])
		}
		if !sameDigests(before["notes/z.go"], after["notes/z.go"]) {
			t.Fatalf("z.go digests moved after a.go grew:\nbefore %v\nafter  %v", before["notes/z.go"], after["notes/z.go"])
		}
	})
	for _, document := range []struct{ name, path, suffix string }{
		{"spec", "specs/example/spec.md", "\n## Further notes\n\nEdited.\n"},
		{"ticket", "specs/example/tickets/one.md", "- [ ] It stays built.\n"},
		{"review record", "reviews/example.md", "# Review record\n"},
	} {
		t.Run(document.name, func(t *testing.T) {
			root, args := seed(t)
			_, before := preparedMembers(t, root, args)
			// The review record starts absent, so its edit is the file's first content.
			current, _ := os.ReadFile(document.path)
			commitReview(t, args, map[string]string{document.path: string(current) + document.suffix}, "edit "+document.name)
			_, after := preparedMembers(t, root, args)
			for _, path := range codePaths {
				if len(before[path]) == 0 || !sameDigests(before[path], after[path]) {
					t.Errorf("%s edit moved %s digests:\nbefore %v\nafter  %v", document.name, path, before[path], after[path])
				}
			}
			if len(after[document.path]) != 1 || sameDigests(before[document.path], after[document.path]) {
				t.Errorf("%s patch = %v before and %v after, want one changed patch", document.path, before[document.path], after[document.path])
			}
		})
	}
}

// TestReviewFileIdentity is RE3, RE4, and RE6. File identity is the declared path, never the
// source ordinal: an inserted earlier file keeps each retained path's digests under both
// rename settings, a changed patch changes its digest, and equal content keeps one
// membership per path.
func TestReviewFileIdentity(t *testing.T) {
	for _, renames := range []bool{true, false} {
		t.Run(fmt.Sprintf("inserted earlier file renames=%v", renames), func(t *testing.T) {
			root, _, args := preflighttest.SeedReviewEvidence(t, false)
			preflighttest.RunGit(t, "config", "diff.renames", fmt.Sprint(renames))
			preflighttest.RunGit(t, "mv", "outside/user.go", "outside/moved.go")
			commitReview(t, args, map[string]string{"notes/z.txt": "retained\n"}, "rename and add")
			_, before := preparedMembers(t, root, args)
			commitReview(t, args, map[string]string{"notes/0-inserted.txt": "inserted\n"}, "insert")
			_, after := preparedMembers(t, root, args)
			_, basePath := before["outside/user.go"]
			if _, tipPath := before["outside/moved.go"]; !tipPath || basePath == renames {
				t.Fatalf("renames=%v rename membership = %v", renames, before)
			}
			moved := 0
			for path, members := range before {
				if !sameDigests(members, after[path]) {
					t.Errorf("retained %s digests changed: %v -> %v", path, members, after[path])
				}
				if len(after[path]) > 0 && after[path][0].id != members[0].id {
					moved++
				}
			}
			if len(after["notes/0-inserted.txt"]) != 1 || moved == 0 {
				t.Fatalf("inserted file members = %v and %d moved ordinals, want an inserted source that moves ordinals", after["notes/0-inserted.txt"], moved)
			}
		})
	}
	t.Run("changed patch", func(t *testing.T) {
		root, _, args := preflighttest.SeedReviewEvidence(t, false)
		commitReview(t, args, map[string]string{"notes/y.txt": "neighbor\n", "notes/z.txt": "first\n"}, "patches")
		_, before := preparedMembers(t, root, args)
		commitReview(t, args, map[string]string{"notes/z.txt": "second\n"}, "edit z")
		identity, after := preparedMembers(t, root, args)
		if len(after["notes/z.txt"]) != 1 || sameDigests(before["notes/z.txt"], after["notes/z.txt"]) {
			t.Fatalf("changed z.txt members = %v -> %v, want a changed digest", before["notes/z.txt"], after["notes/z.txt"])
		}
		if got, want := selectedPatch(t, identity, after["notes/z.txt"]), preflighttest.RawPatch(t, args, "notes/z.txt"); got != want {
			t.Fatalf("changed z.txt patch =\n%s\nwant\n%s", got, want)
		}
	})
	t.Run("duplicate content", func(t *testing.T) {
		root, _, args := preflighttest.SeedReviewEvidence(t, false)
		content := strings.Repeat("equal duplicated content line\n", 900)
		commitReview(t, args, map[string]string{"notes/dup-1.txt": content, "notes/dup-2.txt": content}, "duplicates")
		identity, members := preparedMembers(t, root, args)
		first, second := members["notes/dup-1.txt"], members["notes/dup-2.txt"]
		if len(first) != 1 || len(second) != 1 || first[0].id == second[0].id || first[0].sha == second[0].sha {
			t.Fatalf("duplicate members = %v and %v, want one distinct source per path", first, second)
		}
		if len(first[0].pages) < 3 || len(first[0].pages) != len(second[0].pages) || first[0].pages[0] == second[0].pages[0] {
			t.Fatalf("duplicate pages = %v and %v, want equal-length multi-page patches with distinct headers", first[0].pages, second[0].pages)
		}
		if !slices.Equal(first[0].pages[1:], second[0].pages[1:]) {
			t.Fatalf("duplicate later pages differ: %v and %v", first[0].pages[1:], second[0].pages[1:])
		}
		for _, path := range []string{"notes/dup-1.txt", "notes/dup-2.txt"} {
			if got, want := selectedPatch(t, identity, members[path]), preflighttest.RawPatch(t, args, path); got != want {
				t.Errorf("%s patch =\n%.200s\nwant\n%.200s", path, got, want)
			}
		}
	})
}
