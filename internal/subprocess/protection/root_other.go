//go:build !linux && !darwin

package protection

import (
	"errors"
	"syscall"
)

// systemBase is the system temporary directory.
func systemBase() string { return "/tmp" }

// BootID refuses, because this system names no boot identifier that Bench reads.
func BootID() (string, error) {
	return "", errors.New("protection: this system names no boot identifier")
}

// systemProbe sends signal to target. This system has no zombie rule, so an unreaped zombie
// answers signal zero as present.
func systemProbe(target int, signal syscall.Signal) error { return syscall.Kill(target, signal) }
