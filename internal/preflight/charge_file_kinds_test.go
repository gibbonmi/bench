package preflight

import (
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

func TestChargeRefusesLinkedSpecAndTicketInputs(t *testing.T) {
	for _, input := range []string{"spec", "ticket", "tickets directory"} {
		t.Run(input, func(t *testing.T) {
			root, slug, _ := seedLinkedChargeInput(t, input)
			out, code := Command(chargeArgs(t, root, slug, false))
			if code != 1 || !strings.Contains(out, "source required") ||
				!strings.Contains(out, "wrong-type") ||
				!strings.Contains(out, "restore the named canonical source") ||
				strings.Contains(out, "complete,next}") {
				t.Fatalf("linked %s = (%d):\n%s", input, code, out)
			}
		})
	}
}

func seedLinkedChargeInput(t *testing.T, input string) (root, slug, target string) {
	t.Helper()
	root, slug = seedConformant(t)
	spec := specBody(slug, "- `specs/"+slug+"/` (input fixtures)")
	mustWriteFile(t, "specs/"+slug+"/spec.md", spec)
	outside := t.TempDir()
	switch input {
	case "spec":
		target = filepath.Join(outside, "spec-target.txt")
		if err := os.WriteFile(target, []byte("invalid spec bytes\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove("specs/" + slug + "/spec.md"); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, "specs/"+slug+"/spec.md"); err != nil {
			t.Fatal(err)
		}
	case "ticket":
		target = filepath.Join(outside, "ticket-target.txt")
		if err := os.WriteFile(target, []byte("invalid ticket bytes\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove("specs/" + slug + "/tickets/one.md"); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, "specs/"+slug+"/tickets/one.md"); err != nil {
			t.Fatal(err)
		}
	case "tickets directory":
		target = filepath.Join(outside, "tickets-target")
		if err := os.Rename("specs/"+slug+"/tickets", target); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(target, "one.md"), []byte("invalid ticket bytes\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, "specs/"+slug+"/tickets"); err != nil {
			t.Fatal(err)
		}
	default:
		t.Fatalf("unknown linked input %q", input)
	}
	runGit(t, "add", "-A")
	runGit(t, "commit", "-q", "-m", "linked "+input+" input")
	return root, slug, target
}

func TestChargeRefusesSpecialTicketBeforeRead(t *testing.T) {
	for _, kind := range []string{"dangling link", "FIFO", "socket"} {
		t.Run(kind, func(t *testing.T) {
			root, slug := seedConformant(t)
			path := "specs/" + slug + "/tickets/one.md"
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			switch kind {
			case "dangling link":
				if err := os.Symlink("missing", path); err != nil {
					t.Fatal(err)
				}
			case "FIFO":
				if err := syscall.Mkfifo(path, 0o600); err != nil {
					t.Fatal(err)
				}
			case "socket":
				listener, err := net.Listen("unix", path)
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() { _ = listener.Close() })
			}
			out, code := Command(chargeArgs(t, root, slug, false))
			if code != 1 || !strings.Contains(out, "ticket file not readable") ||
				!strings.Contains(out, "one.md") ||
				strings.Contains(out, "tickets-parse") || strings.Contains(out, "complete,next}") {
				t.Fatalf("%s ticket = (%d):\n%s", kind, code, out)
			}
		})
	}
}
