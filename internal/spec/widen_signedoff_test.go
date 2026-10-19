package spec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
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

// A signed-off entry is a literal path, so pathspec magic in one entry cannot fail the
// scope read for the slug.
func TestSpecWidenReadsSignedOffEntriesLiterally(t *testing.T) {
	root := widenRepo(t, "internal/x/y.go, :(glob)internal/q", "internal/x/y.go")
	out, code := runWiden(t, root, "internal/x/z.go")
	requireWidened(t, root, out, code, widenTable(t, []string{widenTicketPath, "internal/x/z.go", "added"}), widenTicket("internal/x/y.go, :(glob)internal/q, internal/x/z.go"))
}
