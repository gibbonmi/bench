//go:build linux

package preflight

import (
	"errors"
	"syscall"
	"testing"
)

func TestChargeDoesNotOpenLinkedInputTarget(t *testing.T) {
	for _, input := range []string{"spec", "ticket", "tickets directory"} {
		t.Run(input, func(t *testing.T) {
			root, slug, target := seedLinkedChargeInput(t, input)
			opened := watchOpen(t, target)
			_, _ = Command(chargeArgs(t, root, slug, false))
			if opened() {
				t.Fatalf("charge opened the target of linked %s input before refusing it", input)
			}
		})
	}
}

func watchOpen(t *testing.T, path string) func() bool {
	t.Helper()
	fd, err := syscall.InotifyInit1(syscall.IN_CLOEXEC | syscall.IN_NONBLOCK)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = syscall.Close(fd) })
	if _, err := syscall.InotifyAddWatch(fd, path, syscall.IN_OPEN); err != nil {
		t.Fatal(err)
	}
	return func() bool {
		var events [syscall.SizeofInotifyEvent * 4]byte
		n, err := syscall.Read(fd, events[:])
		if errors.Is(err, syscall.EAGAIN) {
			return false
		}
		if err != nil {
			t.Fatal(err)
		}
		return n > 0
	}
}
