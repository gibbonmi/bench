package spec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/tickets"
	"github.com/gibbonmi/bench/internal/toon"
)

func TestSpecWidenRefusesWithoutSignedOffTickets(t *testing.T) {
	for _, tc := range []struct {
		name  string
		files map[string]string
		link  string
	}{
		{"no tickets folder of the slug", nil, ""},
		{"only a file that is not a ticket", map[string]string{"specs/s/tickets/notes.txt": widenTicket("internal/x/y.go")}, ""},
		{"only a ticket that is a symbolic link", nil, "specs/s/tickets/b.md"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			files := map[string]string{"specs/other/tickets/a.md": widenTicket("internal/x/y.go")}
			for path, content := range tc.files {
				files[path] = content
			}
			primary := widenPrimary(t, files, "internal/x/y.go")
			if tc.link != "" {
				link := filepath.Join(primary, filepath.FromSlash(tc.link))
				if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
					t.Fatal(err)
				}
				linkInPlace(t, "../../other/tickets/a.md", link)
				gittest.Output(t, primary, "add", "-A")
				gittest.Output(t, primary, "commit", "-qm", "link ticket")
			}
			root := retireWorktree(t, primary)
			writeTicket(t, root, widenTicketPath, widenTicket("internal/x/y.go"))
			before := readTicket(t, root)
			out, code := runWiden(t, root, "internal/x/z.go")
			requireRefusal(t, root, before, out, code, "holds no ticket of specs/s")
		})
	}
}

// An unresolved default branch refuses with its own message, so an absent tip never reads
// as an unreadable one. Two local branches, neither named main, and no origin leave the
// default branch unresolved.
func TestSpecWidenRefusesAnUnresolvedDefaultBranch(t *testing.T) {
	primary := widenPrimary(t, map[string]string{widenTicketPath: widenTicket("internal/x/y.go")}, "internal/x/y.go")
	root := retireWorktree(t, primary)
	gittest.MustRun(t, primary, "branch", "-m", "main", "trunk")
	before := readTicket(t, root)
	out, code := runWiden(t, root, "internal/x/z.go")
	want := toon.Errorf("default branch is unresolved", "widen reads the signed-off scope from the default-branch tip") + "\n"
	if out != want {
		t.Fatalf("out = %q, want %q", out, want)
	}
	requireRefusal(t, root, before, out, code)
}

// The signed-off ticket set is the tracked-ticket enumeration of the tip, so a later copy
// of a basename grants no scope and a path that spec-TOON cannot represent refuses.
func TestSpecWidenReadsSignedOffTicketsThroughTheOwner(t *testing.T) {
	for _, tc := range []struct {
		name, path, want string
	}{
		{"a duplicate basename", "specs/s/tickets/sub/a.md", outsideScope},
		{"an unrepresentable path", "specs/s/tickets/b\x1b.md", tickets.Refusal{Path: "b\x1b.md"}.Message("")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			files := map[string]string{widenTicketPath: widenTicket("internal/x/y.go"), tc.path: widenTicket("internal/other/o.go")}
			root := retireWorktree(t, widenPrimary(t, files, "internal/x/y.go", "internal/other/o.go"))
			before := readTicket(t, root)
			out, code := runWiden(t, root, "internal/other/z.go")
			requireRefusal(t, root, before, out, code, tc.want)
		})
	}
}

// A signed-off entry is a literal path, so pathspec magic in one entry cannot fail the
// scope read for the slug.
func TestSpecWidenReadsSignedOffEntriesLiterally(t *testing.T) {
	root := widenRepo(t, "internal/x/y.go, :(glob)internal/q", "internal/x/y.go")
	out, code := runWiden(t, root, "internal/x/z.go")
	requireWidened(t, root, out, code, widenTable(t, []string{widenTicketPath, "internal/x/z.go", "added"}), widenTicket("internal/x/y.go, :(glob)internal/q, internal/x/z.go"))
}
