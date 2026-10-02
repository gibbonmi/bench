package recordcmd_test

import (
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordcmd"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/toon"
)

// excerpted is a recorded fixture and the edit that names path as the excerpt.
func excerpted(t *testing.T, path string) (*recordtest.Fixture, map[string]string) {
	t.Helper()
	return recorded(t, 1), map[string]string{"--excerpt": path}
}

func TestRecordVerificationEmbedsTheExcerptBytes(t *testing.T) {
	const text = "line\tone\rtwo\x1bthree\n"
	f, edit := excerpted(t, excerptFile(t, "excerpt.txt", text))
	verify(t, f, edit)
	if got := onlyResult(t, f).NativeRef.Excerpt; got != text {
		t.Fatalf("excerpt = %q, want %q", got, text)
	}
}

func TestRecordVerificationDigestsTheExcerpt(t *testing.T) {
	const text = "  padded result\n\n"
	f, edit := excerpted(t, excerptFile(t, "excerpt.txt", text))
	verify(t, f, edit)
	if got := onlyResult(t, f).NativeRef.Digest; got != rr.Digest([]byte(text)) {
		t.Fatalf("excerpt digest = %s, want %s", got, rr.Digest([]byte(text)))
	}
}

func TestRecordRefusesAnAbsentExcerpt(t *testing.T) {
	path := filepath.Join(t.TempDir(), "absent.txt")
	f, edit := excerpted(t, path)
	if out := refuseVerify(t, f, "is absent", edit); out != toon.RecordError(path, bounds.StateAbsent, "")+"\n" {
		t.Fatalf("output = %q, want the toon.RecordError line for %s", out, path)
	}
}

func TestRecordRefusesALinkedExcerpt(t *testing.T) {
	path := filepath.Join(t.TempDir(), "link.txt")
	symlink(t, excerptFile(t, "target.txt", "linked bytes\n"), path)
	f, edit := excerpted(t, path)
	refuseVerify(t, f, "is wrong-type", edit)
}

func TestRecordRefusesASpecialExcerpt(t *testing.T) {
	path := filepath.Join(t.TempDir(), "fifo")
	if err := syscall.Mkfifo(path, 0o644); err != nil {
		capability.Capability(t, capability.Fifo, "FIFOs unavailable on this filesystem: "+err.Error())
	}
	f, edit := excerpted(t, path)
	type result struct {
		out  string
		code int
	}
	args, done := verifyArgs(t, edit), make(chan result, 1)
	go func() {
		out, code := recordcmd.Command(f.Root, args)
		done <- result{out, code}
	}()
	select {
	case got := <-done:
		if got.code != 1 || !strings.Contains(got.out, "is wrong-type") {
			t.Fatalf("FIFO excerpt = exit %d, output %q; want exit 1 naming is wrong-type", got.code, got.out)
		}
	case <-time.After(bounds.TestDeadline(0)):
		t.Fatal("verification form blocked on the FIFO excerpt")
	}
}

func TestRecordRefusesAnEmptyExcerpt(t *testing.T) {
	f, edit := excerpted(t, excerptFile(t, "empty.txt", ""))
	refuseVerify(t, f, "is empty", edit)
}

func TestRecordRefusesAMalformedExcerpt(t *testing.T) {
	f, edit := excerpted(t, excerptFile(t, "malformed.txt", "bad \xff byte\n"))
	refuseVerify(t, f, "is malformed", edit)
}

func TestRecordRefusesAnOversizedExcerpt(t *testing.T) {
	f, edit := excerpted(t, excerptFile(t, "oversized.txt", strings.Repeat("o", int(bounds.ControlRecordLimit)+1)))
	refuseVerify(t, f, "is unreadable", edit)
}

func TestRecordReadsAnExcerptPathWithGlobCharacters(t *testing.T) {
	const text = "glob path result\n"
	f, edit := excerpted(t, excerptFile(t, "my [glob]*.txt", text))
	verify(t, f, edit)
	if got := onlyResult(t, f).NativeRef.Excerpt; got != text {
		t.Fatalf("excerpt = %q, want %q", got, text)
	}
}
