package canary

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
)

// TestReadMarkerRefusesSpecialFile plants a FIFO where a marker file belongs — the
// hostile input a discovered path can carry — and grades the reader on refusing it by
// name. The read runs off the test goroutine behind a deadline because an unguarded open
// of a FIFO never returns: without the deadline the failure is a wedged suite that only
// the sweep's own timeout would ever report.
func TestReadMarkerRefusesSpecialFile(t *testing.T) {
	root := t.TempDir()
	fx := canaryFixture(root, mappedFamily(t), "fifo")
	mkdir(t, fx)
	path := filepath.Join(fx, checkFileName)
	if err := exec.Command("mkfifo", path).Run(); err != nil {
		t.Fatal(err)
	}

	type answer struct {
		present bool
		err     error
	}
	done := make(chan answer, 1)
	go func() {
		_, present, err := readMarker(fx, checkFileName)
		done <- answer{present: present, err: err}
	}()

	select {
	case got := <-done:
		if got.err == nil {
			t.Fatalf("readMarker accepted a FIFO (present=%v), want a refusal", got.present)
		}
		if !strings.Contains(got.err.Error(), path) {
			t.Errorf("refusal %q names no offending path, want %q", got.err, path)
		}
	case <-time.After(bounds.TestDeadline(bounds.TestDeadlineFloor)):
		t.Fatal("readMarker blocked on the FIFO: the path is opened before it is typed")
	}
}

// TestReadMarkerSeparatesAbsentFromEmpty pins the two states a marker's consumers have to
// tell apart — no file at all against a file holding nothing — since one is a default and
// the other is a defect.
func TestReadMarkerSeparatesAbsentFromEmpty(t *testing.T) {
	root := t.TempDir()
	absent := canaryFixture(root, mappedFamily(t), "absent")
	mkdir(t, absent)
	empty := canaryFixture(root, mappedFamily(t), "empty")
	write(t, filepath.Join(empty, checkFileName), "")

	if name, present, err := readMarker(absent, checkFileName); err != nil || present || name != "" {
		t.Errorf("absent marker read as (%q, %v, %v), want (\"\", false, nil)", name, present, err)
	}
	if name, present, err := readMarker(empty, checkFileName); err != nil || !present || name != "" {
		t.Errorf("empty marker read as (%q, %v, %v), want (\"\", true, nil)", name, present, err)
	}
}

// TestReadMarkerTrimsSurroundingSpace grades the trim every marker consumer depends on: a
// file written with a trailing newline and one written without name the same thing.
func TestReadMarkerTrimsSurroundingSpace(t *testing.T) {
	root := t.TempDir()
	for _, content := range []string{"TestOwner", "TestOwner\n", "  TestOwner \t\n"} {
		fx := canaryFixture(root, mappedFamily(t), "trim")
		write(t, filepath.Join(fx, checkFileName), content)
		name, present, err := readMarker(fx, checkFileName)
		if err != nil || !present || name != "TestOwner" {
			t.Errorf("marker %q read as (%q, %v, %v), want (\"TestOwner\", true, nil)", content, name, present, err)
		}
	}
}
