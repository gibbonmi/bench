package spec

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/tickets"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// widenTicketPath is the repository-relative ticket that the widen fixtures edit.
const widenTicketPath = "specs/s/tickets/a.md"

// wrongType is the refusal kind of a ticket that is not a regular file.
const wrongType = "is wrong-type"

// widenTicket is a complete ticket whose Writes: value is writes.
func widenTicket(writes string) string {
	return "# A\n\nBlocked by: none\nWrites: " + writes + "\nCovers: CR1\n\n## What to build\n\nx\n\n## Acceptance\n\n- [ ] x\n"
}

// widenRepo commits a.md with the signed-off writes and one file at each tree path on
// main, and returns a linked worktree of that repository.
func widenRepo(t *testing.T, writes string, tree ...string) string {
	t.Helper()
	return retireWorktree(t, widenPrimary(t, map[string]string{widenTicketPath: widenTicket(writes)}, tree...))
}

// widenPrimary commits files and one file at each tree path on main, and returns the
// primary checkout.
func widenPrimary(t *testing.T, files map[string]string, tree ...string) string {
	t.Helper()
	for _, path := range tree {
		files[path] = "x\n"
	}
	return retirePrimary(t, "s", "Status: staged\n", files)
}

// writeTicket replaces the worktree bytes of one repository-relative ticket.
func writeTicket(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readTicket(t *testing.T, root string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(widenTicketPath)))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// widenTable is the expected widen table, rendered through the producer's encoder.
func widenTable(t *testing.T, rows ...[]string) string {
	t.Helper()
	out, err := toon.Table("widen", []string{"ticket", "entry", "state"}, rows)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// runWiden widens a.md of slug s from dir.
func runWiden(t *testing.T, dir string, args ...string) (string, int) {
	t.Helper()
	return runSpec(t, dir, append([]string{"widen", "s", "--ticket", "a.md"}, args...)...)
}

// requireRefusal asserts exit 1, unchanged ticket bytes, and a refusal that holds each
// of names, so a refusal from a later guard cannot stand in for the expected one.
func requireRefusal(t *testing.T, root, before, out string, code int, names ...string) {
	t.Helper()
	if code != 1 {
		t.Fatalf("code = %d, want 1; out = %q", code, out)
	}
	if after := readTicket(t, root); after != before {
		t.Fatalf("ticket changed on refusal: %q -> %q", before, after)
	}
	for _, name := range names {
		if !strings.Contains(out, name) {
			t.Errorf("out = %q, want it to name %q", out, name)
		}
	}
}

// requireWidened asserts exit 0, the table want unless it is empty, and the ticket bytes after.
func requireWidened(t *testing.T, root, out string, code int, want, after string) {
	t.Helper()
	if code != 0 || (want != "" && out != want) {
		t.Fatalf("widen = (%d, %q), want exit 0 and %q", code, out, want)
	}
	if got := readTicket(t, root); got != after {
		t.Errorf("ticket = %q, want %q", got, after)
	}
}

// linkInPlace replaces link with a symbolic link to target.
func linkInPlace(t *testing.T, target, link string) {
	t.Helper()
	if err := os.RemoveAll(link); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, link); err != nil {
		capability.Unavailable(t, capability.Symlink, err.Error())
	}
}

func TestSpecWidenAddsAnInScopeEntry(t *testing.T) {
	root := widenRepo(t, "internal/x/y.go", "internal/x/y.go")
	head := gittest.Output(t, root, "rev-parse", "HEAD")
	out, code := runWiden(t, root, "internal/x/y_test.go")
	requireWidened(t, root, out, code, widenTable(t, []string{widenTicketPath, "internal/x/y_test.go", "added"}), widenTicket("internal/x/y.go, internal/x/y_test.go"))
	if after := gittest.Output(t, root, "rev-parse", "HEAD"); after != head {
		t.Errorf("HEAD moved from %s to %s", head, after)
	}
}

func TestSpecWidenRefusesAnOutOfScopeEntry(t *testing.T) {
	root := widenRepo(t, "internal/x/y.go", "internal/x/y.go")
	before := readTicket(t, root)
	out, code := runWiden(t, root, "internal/other/z.go")
	requireRefusal(t, root, before, out, code, "internal/other/z.go", "--decision <reference>")
	// The rerun names the whole request, with each argument as one shell token.
	out, code = runWiden(t, root, "internal/x/w.go", "internal/other/z.go (new)", "internal/other/a b.go")
	requireRefusal(t, root, before, out, code, "rerun bench spec widen s --ticket a.md --decision <reference> internal/x/w.go 'internal/other/z.go (new)' 'internal/other/a b.go'\n")
	// A slug and a ticket basename with a space each stay one shell token.
	spaced := retireWorktree(t, retirePrimary(t, "s t", "Status: staged\n", map[string]string{"specs/s t/tickets/a b.md": widenTicket("internal/x/y.go")}))
	out, code = runSpec(t, spaced, "widen", "s t", "--ticket", "a b.md", "internal/other/z.go")
	if want := "rerun bench spec widen 's t' --ticket 'a b.md' --decision <reference> internal/other/z.go\n"; code != 1 || !strings.Contains(out, want) {
		t.Errorf("spaced widen = (%d, %q), want exit 1 and %q", code, out, want)
	}
}

func TestSpecWidenAdmitsADecidedEntry(t *testing.T) {
	root := widenRepo(t, "internal/x/y.go", "internal/x/y.go")
	out, code := runWiden(t, root, "--decision", "r1", "internal/x/w.go", "internal/other/z.go")
	want := widenTable(t, []string{widenTicketPath, "internal/x/w.go", "added"}, []string{widenTicketPath, "internal/other/z.go", "decided"})
	requireWidened(t, root, out, code, want, widenTicket("internal/x/y.go, internal/x/w.go, internal/other/z.go"))
}

func TestSpecWidenRefusesANewEntryThatGitTracks(t *testing.T) {
	for _, tc := range []struct {
		name, entry, path string
	}{
		{"a tracked file", "internal/x/z.go (new)", "internal/x/z.go"},
		{"a tracked directory", "internal/x/ (new)", "internal/x"},
		{"a file that only the worktree branch tracks", "internal/x/b.go (new)", "internal/x/b.go"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := widenRepo(t, "internal/x/y.go", "internal/x/y.go", "internal/x/z.go")
			writeTicket(t, root, "internal/x/b.go", "x\n")
			gittest.Output(t, root, "add", "-A")
			gittest.Output(t, root, "commit", "-qm", "branch file")
			before := readTicket(t, root)
			out, code := runWiden(t, root, "internal/x/w.go", tc.entry)
			requireRefusal(t, root, before, out, code, trackedNew+": "+tc.path+" —", "drop the "+tickets.NewMarker+" marker")
		})
	}
	// An untracked working file is no tracked path, so its (new) entry widens.
	root := widenRepo(t, "internal/x/y.go", "internal/x/y.go")
	writeTicket(t, root, "internal/x/w.go", "x\n")
	out, code := runWiden(t, root, "internal/x/w.go (new)")
	requireWidened(t, root, out, code, widenTable(t, []string{widenTicketPath, "internal/x/w.go (new)", "added"}), widenTicket("internal/x/y.go, internal/x/w.go (new)"))
}

func TestSpecWidenUsage(t *testing.T) {
	root := widenRepo(t, "internal/x/y.go", "internal/x/y.go")
	before := readTicket(t, root)
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"empty decision", []string{"widen", "s", "--ticket", "a.md", "--decision", "", "internal/other/z.go"}},
		{"empty ticket", []string{"widen", "s", "--ticket", "", "internal/x/z.go"}},
		{"missing ticket", []string{"widen", "s", "internal/x/z.go"}},
		{"missing entry", []string{"widen", "s", "--ticket", "a.md"}},
		{"unknown flag", []string{"widen", "s", "--ticket", "a.md", "--base", "HEAD", "internal/x/z.go"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, code := runSpec(t, root, tc.args...)
			if code != 2 || !strings.HasPrefix(out, "usage: bench spec widen") {
				t.Fatalf("widen = (%d, %q), want exit 2 and the widen usage", code, out)
			}
			if after := readTicket(t, root); after != before {
				t.Fatalf("ticket changed on a usage error: %q", after)
			}
		})
	}
	out, code := runSpec(t, root, "widen", "--help")
	if code != 0 || !strings.HasPrefix(out, "usage: bench spec widen <slug> --ticket <basename>") {
		t.Errorf("widen --help = (%d, %q), want the declared usage", code, out)
	}
}

func TestSpecWidenReadsTheSignedOffScope(t *testing.T) {
	root := widenRepo(t, "internal/x/y.go", "internal/x/y.go")
	// The build widens a.md and adds b.md on its own branch; neither reaches the default tip.
	writeTicket(t, root, widenTicketPath, widenTicket("internal/x/y.go, internal/new/n.go"))
	writeTicket(t, root, "specs/s/tickets/b.md", widenTicket("internal/x/y.go"))
	gittest.Output(t, root, "add", "-A")
	gittest.Output(t, root, "commit", "-qm", "self-widened plan")
	before := readTicket(t, root)
	out, code := runWiden(t, root, "internal/new/other.go")
	requireRefusal(t, root, before, out, code, outsideScope)

	out, code = runSpec(t, root, "widen", "s", "--ticket", "b.md", "internal/x/z.go")
	if code != 0 {
		t.Fatalf("worktree-only ticket widen = (%d, %q), want exit 0", code, out)
	}
	if want := widenTable(t, []string{"specs/s/tickets/b.md", "internal/x/z.go", "added"}); out != want {
		t.Errorf("out = %q, want %q", out, want)
	}
}

func TestSpecWidenDirectoryPrefix(t *testing.T) {
	const signedOff = "internal/x/y.go, internal/status, CHANGELOG.md, internal/n/new.go (new), internal/git"
	tree := []string{"internal/x/y.go", "internal/status/s.go", "CHANGELOG.md", "internal/git/g.go", "README.md"}
	for _, tc := range []struct {
		name, entry string
		code        int
	}{
		{"a file entry admits its sibling subtree", "internal/x/sub/z.go", 0},
		{"a directory entry refuses a sibling directory", "internal/other/x.go", 1},
		{"an absent new entry admits its parent directory", "internal/n/other.go", 0},
		{"a directory entry refuses a sibling with its prefix spelling", "internal/git2/x.go", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := widenRepo(t, signedOff, tree...)
			before := readTicket(t, root)
			out, code := runWiden(t, root, tc.entry)
			if tc.code == 1 {
				requireRefusal(t, root, before, out, code, outsideScope)
				return
			}
			requireWidened(t, root, out, code, "", widenTicket(signedOff+", "+tc.entry))
		})
	}
	// A root-level file is its own prefix: another ticket can widen into that file, and
	// no ticket can widen into another root file.
	t.Run("a root file entry is its own prefix", func(t *testing.T) {
		root := widenRepo(t, signedOff, tree...)
		before := readTicket(t, root)
		out, code := runWiden(t, root, "README.md")
		requireRefusal(t, root, before, out, code, outsideScope)
		writeTicket(t, root, "specs/s/tickets/b.md", widenTicket("internal/x/y.go"))
		out, code = runSpec(t, root, "widen", "s", "--ticket", "b.md", "CHANGELOG.md")
		if want := widenTable(t, []string{"specs/s/tickets/b.md", "CHANGELOG.md", "added"}); code != 0 || out != want {
			t.Errorf("root file widen = (%d, %q), want exit 0 and %q", code, out, want)
		}
	})
}

func TestSpecWidenWritesAllOrNothing(t *testing.T) {
	root := widenRepo(t, "internal/x/y.go", "internal/x/y.go")
	before := readTicket(t, root)
	out, code := runWiden(t, root, "internal/x/z.go", "internal/other/z.go")
	requireRefusal(t, root, before, out, code, outsideScope)
}

func TestSpecWidenRefusesThePrimaryCheckout(t *testing.T) {
	root := widenPrimary(t, map[string]string{widenTicketPath: widenTicket("internal/x/y.go")}, "internal/x/y.go")
	before := readTicket(t, root)
	out, code := runWiden(t, root, "internal/x/z.go")
	requireRefusal(t, root, before, out, code)
	if shared := usage.PrimaryCheckoutRefusal(); !strings.HasPrefix(out, shared) {
		t.Errorf("out = %q, want the shared refusal as its prefix (%q)", out, shared)
	}
}

func TestSpecWidenRefusesHostileEntries(t *testing.T) {
	for _, tc := range []struct {
		name, entry string
	}{
		{"a dot-dot segment", "internal/x/../../cmd/y.go"},
		{"an absolute path", "/internal/x/z.go"},
		{"the entry separator", "internal/x/a.go, internal/x/b.go"},
		{"an ESC byte", "internal/x/a\x1b.go"},
		{"a newline that TOON permits", "internal/x/a\nWrites: internal/other/b.go"},
		{"a carriage return that TOON permits", "internal/x/a\r.go"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := widenRepo(t, "internal/x/y.go", "internal/x/y.go")
			before := readTicket(t, root)
			out, code := runWiden(t, root, "--decision", "r1", tc.entry)
			requireRefusal(t, root, before, out, code, "error: entry ")
		})
	}
	root := widenRepo(t, "internal/x/y.go", "internal/x/y.go")
	const legal = "internal/x/a b[1]\t.go"
	out, code := runWiden(t, root, legal, " internal/x/padded.go ")
	requireWidened(t, root, out, code, "", widenTicket("internal/x/y.go, "+legal+", internal/x/padded.go"))
}

func TestSpecWidenRefusesHostileTickets(t *testing.T) {
	t.Run("a FIFO", func(t *testing.T) {
		root := widenRepo(t, "internal/x/y.go", "internal/x/y.go")
		ticket := filepath.Join(root, filepath.FromSlash(widenTicketPath))
		if err := os.Remove(ticket); err != nil {
			t.Fatal(err)
		}
		if err := syscall.Mkfifo(ticket, 0o644); err != nil {
			capability.Unavailable(t, capability.Fifo, err.Error())
		}
		t.Chdir(root)
		type result struct {
			out  string
			code int
		}
		done := make(chan result, 1)
		go func() {
			out, code, _ := Command([]string{"widen", "s", "--ticket", "a.md", "internal/x/z.go"})
			done <- result{out, code}
		}()
		select {
		case got := <-done:
			if got.code != 1 || !strings.Contains(got.out, wrongType) {
				t.Fatalf("widen = (%d, %q), want exit 1 and %q", got.code, got.out, wrongType)
			}
		case <-time.After(bounds.TestDeadline(0)):
			t.Fatal("widen blocked on a FIFO ticket, so it read before classifying the path")
		}
	})
	// Each refusal reads a.md through readTicket, which follows a link, so an unchanged
	// read also proves that the link target kept its bytes.
	for _, tc := range []struct {
		name, want string
		args       []string
		setup      func(t *testing.T, root string)
	}{
		{"a symbolic link", wrongType, nil, func(t *testing.T, root string) {
			writeTicket(t, root, "outside.md", widenTicket("internal/x/y.go"))
			linkInPlace(t, filepath.Join(root, "outside.md"), filepath.Join(root, filepath.FromSlash(widenTicketPath)))
		}},
		{"a linked tickets folder", wrongType, nil, func(t *testing.T, root string) {
			outside := filepath.Join(t.TempDir(), "tickets")
			writeTicket(t, outside, "a.md", widenTicket("internal/x/y.go"))
			linkInPlace(t, outside, filepath.Join(root, "specs", "s", "tickets"))
		}},
		{"no Writes: line", "has no Writes: line", nil, func(t *testing.T, root string) {
			writeTicket(t, root, widenTicketPath, "# A\n\nBlocked by: none\nCovers: CR1\n")
		}},
		{"a ticket basename with a path", "ticket address", []string{"widen", "s", "--ticket", "../tickets/a.md", "internal/x/z.go"}, nil},
		{"a slug with a path", "ticket address", []string{"widen", "s/../s", "--ticket", "a.md", "internal/x/z.go"}, nil},
		{"a clean ticket basename below a folder", "ticket address", []string{"widen", "s", "--ticket", "sub/a.md", "internal/x/z.go"}, nil},
		{"a clean slug below a folder", "ticket address", []string{"widen", "a/b", "--ticket", "a.md", "internal/x/z.go"}, nil},
		{"a ticket that is not Markdown", "ticket address", []string{"widen", "s", "--ticket", "a.txt", "internal/x/z.go"}, func(t *testing.T, root string) {
			writeTicket(t, root, "specs/s/tickets/a.txt", widenTicket("internal/x/y.go"))
		}},
		{"an absent ticket", "is absent", []string{"widen", "s", "--ticket", "c.md", "internal/x/z.go"}, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := widenRepo(t, "internal/x/y.go", "internal/x/y.go")
			if tc.setup != nil {
				tc.setup(t, root)
			}
			args := tc.args
			if args == nil {
				args = []string{"widen", "s", "--ticket", "a.md", "internal/x/z.go"}
			}
			before := readTicket(t, root)
			out, code := runSpec(t, root, args...)
			requireRefusal(t, root, before, out, code, tc.want)
		})
	}
	out, code := runWiden(t, t.TempDir(), "internal/x/z.go")
	if code != 1 || out != toon.NotInRepo()+"\n" {
		t.Errorf("outside a repository = (%d, %q), want exit 1 and the not-in-repo error", code, out)
	}
}

func TestSpecWidenFromADeeperDirectory(t *testing.T) {
	root := widenRepo(t, "internal/x/y.go", "internal/x/y.go")
	out, code := runWiden(t, filepath.Join(root, "internal", "x"), "internal/x/z.go")
	requireWidened(t, root, out, code, "", widenTicket("internal/x/y.go, internal/x/z.go"))
}
