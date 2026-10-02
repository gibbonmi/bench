package recordcmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordcmd"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// refuse runs the chunk form for chunk 1 of a fresh commit and requires exit 1.
func refuse(t *testing.T, f *recordtest.Fixture) string {
	t.Helper()
	base, tip := advance(f, "chunk 1")
	out, code := recordcmd.Command(f.Root, chunkArgs(t, "1", base, tip))
	if code != 1 {
		t.Fatalf("chunk form = exit %d, output %q; want exit 1", code, out)
	}
	return out
}

// placed prepares the reviews directory of f and returns the record path in it.
func placed(t *testing.T, f *recordtest.Fixture) string {
	t.Helper()
	path := recordFile(t, f)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func symlink(t *testing.T, target, path string) {
	t.Helper()
	if err := os.Symlink(target, path); err != nil {
		capability.Capability(t, capability.Symlink, "symlinks unavailable: "+err.Error())
	}
}

func TestRecordRefusesADanglingRecordLink(t *testing.T) {
	f := linked(t, 1)
	target := filepath.Join(t.TempDir(), "target.md")
	symlink(t, target, placed(t, f))
	refuse(t, f)
	if _, err := os.Lstat(target); !os.IsNotExist(err) {
		t.Fatalf("link target after refusal: %v, want absent", err)
	}
}

func TestRecordRefusesALiveRecordLink(t *testing.T) {
	f := linked(t, 1)
	target := filepath.Join(t.TempDir(), "target.md")
	if err := os.WriteFile(target, []byte("outside the tree\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	symlink(t, target, placed(t, f))
	refuse(t, f)
	if got := string(bytesOf(t, target)); got != "outside the tree\n" {
		t.Fatalf("link target = %q, want its bytes unchanged", got)
	}
}

func TestRecordRefusesASpecialRecordFile(t *testing.T) {
	f := linked(t, 1)
	if err := syscall.Mkfifo(placed(t, f), 0o644); err != nil {
		capability.Capability(t, capability.Fifo, "FIFOs unavailable on this filesystem: "+err.Error())
	}
	base, tip := advance(f, "chunk 1")
	args, exit := chunkArgs(t, "1", base, tip), make(chan int, 1)
	go func() {
		_, code := recordcmd.Command(f.Root, args)
		exit <- code
	}()
	select {
	case code := <-exit:
		if code != 1 {
			t.Fatalf("chunk form over a FIFO = exit %d, want 1", code)
		}
	case <-time.After(bounds.TestDeadline(0)):
		t.Fatal("chunk form blocked on the FIFO")
	}
}

func TestRecordRefusesAnEmptyRecordFile(t *testing.T) {
	f := linked(t, 1)
	path := placed(t, f)
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	refuse(t, f)
	if data := bytesOf(t, path); len(data) != 0 {
		t.Fatalf("record = %q, want it empty", data)
	}
}

// refuseDocument writes a valid record for chunk 1, rewrites it through edit, and
// requires that the next chunk form refuses and leaves the edited bytes unchanged.
func refuseDocument(t *testing.T, edit func(string) string) {
	t.Helper()
	f := linked(t, 1)
	base, tip := advance(f, "chunk 1")
	record(t, f, "1", base, tip)
	path := recordFile(t, f)
	edited := edit(string(bytesOf(t, path)))
	if err := os.WriteFile(path, []byte(edited), 0o644); err != nil {
		t.Fatal(err)
	}
	refuse(t, f)
	if got := string(bytesOf(t, path)); got != edited {
		t.Fatalf("record = %q, want its bytes unchanged %q", got, edited)
	}
}

func TestRecordRefusesAnUnterminatedFence(t *testing.T) {
	refuseDocument(t, func(document string) string { return strings.TrimSuffix(document, "```\n") })
}

func TestRecordRefusesAnInvalidRecord(t *testing.T) {
	refuseDocument(t, func(document string) string {
		return strings.Replace(document, "\"version\": 2,", "\"version\": 2,\n  \"version\": 2,", 1)
	})
}

func TestRecordRefusesARenderedRecordOverTheBound(t *testing.T) {
	f := linked(t, 2)
	base, tip := advance(f, "chunk 1")
	record(t, f, "1", base, tip)
	path := recordFile(t, f)
	document := string(bytesOf(t, path)) + "\n"
	padded := document + strings.Repeat("p", int(bounds.ControlRecordLimit)-len(document)-1) + "\n"
	if err := os.WriteFile(path, []byte(padded), 0o644); err != nil {
		t.Fatal(err)
	}
	read(t, f)
	base, tip = advance(f, "chunk 2")
	out, code := recordcmd.Command(f.Root, chunkArgs(t, "2", base, tip))
	if code != 1 {
		t.Fatalf("chunk form over the bound = exit %d, output %q; want exit 1", code, out)
	}
	if got := string(bytesOf(t, path)); got != padded {
		t.Fatalf("refusal changed the record")
	}
}

func TestRecordChunkRefusesAControlCharacter(t *testing.T) {
	f := linked(t, 1)
	out, code := recordcmd.Command(f.Root, chunkArgs(t, "1\x1b", f.Tip(), "no-such-rev"))
	if code != 1 || !strings.Contains(out, "--chunk") || strings.Contains(out, "no-such-rev") || strings.Contains(out, "\x1b") {
		t.Fatalf("control character in --chunk = exit %d, output %q; want exit 1 naming --chunk and no revision refusal", code, out)
	}
}

func TestRecordRefusesThePrimaryCheckout(t *testing.T) {
	f, primary := recordtest.NewLinked(t, 1, recordtest.Delegate)
	out, code := recordcmd.Command(primary, chunkArgs(t, "1", f.Tip(), f.Tip()))
	if want := usage.PrimaryCheckoutRefusal() + "\n"; code != 1 || out != want {
		t.Fatalf("primary checkout = exit %d, output %q; want exit 1 and %q", code, out, want)
	}
	if _, err := os.Lstat(filepath.Join(primary, "reviews")); !os.IsNotExist(err) {
		t.Fatalf("primary reviews directory: %v, want absent", err)
	}
}

func TestRecordPrimaryRefusalComesFirst(t *testing.T) {
	f, primary := recordtest.NewLinked(t, 1, recordtest.Delegate)
	out, code := recordcmd.Command(primary, chunkArgs(t, "1", f.Tip(), "no-such-rev"))
	if code != 1 || !strings.Contains(out, usage.PrimaryCheckoutRefusal()) || strings.Contains(out, "no-such-rev") {
		t.Fatalf("primary checkout with an unknown revision = exit %d, output %q; want only the primary-checkout refusal", code, out)
	}
}

func TestRecordFailedTemporaryWriteChangesNothing(t *testing.T) {
	f := linked(t, 1)
	base, tip := advance(f, "chunk 1")
	record(t, f, "1", base, tip)
	path := recordFile(t, f)
	dir := filepath.Dir(path)
	before := bytesOf(t, path)
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })
	probe := filepath.Join(dir, "probe")
	if err := os.WriteFile(probe, nil, 0o644); err == nil {
		_ = os.Remove(probe)
		capability.Capability(t, capability.Privilege, "the test process writes into a mode 0555 directory, so a failed temporary write is unobservable")
	}
	listing := func() []string {
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatal(err)
		}
		var names []string
		for _, entry := range entries {
			names = append(names, entry.Name())
		}
		return names
	}
	names := listing()
	refuse(t, f)
	if after := bytesOf(t, path); !bytes.Equal(before, after) {
		t.Fatalf("refusal changed the record")
	}
	if got := listing(); !reflect.DeepEqual(got, names) {
		t.Fatalf("reviews listing = %q, want %q", got, names)
	}
}

func TestRecordRefusesOutsideARepository(t *testing.T) {
	out, code := recordcmd.Command("", chunkArgs(t, "1", "HEAD", "HEAD"))
	if want := toon.NotInRepo() + "\n"; code != 1 || out != want {
		t.Fatalf("empty root = exit %d, output %q; want exit 1 and %q", code, out, want)
	}
}
