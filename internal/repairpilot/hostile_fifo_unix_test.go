//go:build linux || darwin

package repairpilot

import (
	"syscall"
	"testing"
)

func platformStoredHostileFixtures() []storedHostileFixture {
	return []storedHostileFixture{{
		name: "fifo",
		make: func(t *testing.T, path string, _ Options) {
			if err := syscall.Mkfifo(path, 0o600); err != nil {
				t.Fatal(err)
			}
		},
	}}
}
