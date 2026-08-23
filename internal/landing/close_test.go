package landing

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestTicketsOnlyFolder pins the predicate both the close step and `bench status`
// read. The shapes that matter are the two the close step branches on: tickets with
// no spec.md, against a folder carrying one. It also pins the names that must never
// resolve to a direct child of specs/.
func TestTicketsOnlyFolder(t *testing.T) {
	root := t.TempDir()
	write := func(rel, content string) {
		t.Helper()
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("specs/tickets-only/tickets/one.md", "# One\n")
	write("specs/spec-backed/tickets/one.md", "# One\n")
	write("specs/spec-backed/spec.md", "Status: staged\n")
	write("specs/odd name [x]*/tickets/one.md", "# One\n")
	write("specs/loose.md", "# Loose\n")
	if err := os.MkdirAll(filepath.Join(root, "specs", "empty"), 0o755); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name string
		want bool
		why  string
	}{
		{name: "tickets-only", want: true, why: "a direct child of specs/ holding no spec.md"},
		{name: "odd name [x]*", want: true, why: "a name a shell would expand still names its own folder"},
		{name: "empty", want: true, why: "no spec.md is no spec.md"},
		{name: "spec-backed", want: false, why: "a folder carrying spec.md takes the status flip"},
		{name: "absent", want: false, why: "a slug naming no folder is not tickets-only"},
		{name: "loose.md", want: false, why: "a file is not a folder"},
		{name: "spec-backed/tickets", want: false, why: "not a direct child of specs/"},
		{name: "..", want: false, why: "a traversal names no direct child"},
		{name: ".", want: false, why: "specs/ itself is not a spec folder"},
		{name: "", want: false, why: "an empty slug names nothing"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := TicketsOnlyFolder(root, tc.name); got != tc.want {
				t.Fatalf("TicketsOnlyFolder(%q) = %v, want %v: %s", tc.name, got, tc.want, tc.why)
			}
		})
	}

	got, err := TicketsOnlyFolders(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"empty", "odd name [x]*", "tickets-only"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("TicketsOnlyFolders = %v, want %v", got, want)
	}
}

// A repository with no specs/ directory counts zero rather than failing, so the status
// row it feeds is a signal and not an error path.
func TestTicketsOnlyFoldersWithoutSpecsDirectory(t *testing.T) {
	got, err := TicketsOnlyFolders(t.TempDir())
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if len(got) != 0 {
		t.Fatalf("got %v, want none", got)
	}
}
