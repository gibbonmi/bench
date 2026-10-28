//go:build !linux

package gate

import (
	"errors"
	"syscall"
)

// lowerPriority sets this process to backgroundNice. Outside Linux the nice value is an
// attribute of the process, and the platform has no I/O class that a process sets for itself.
func lowerPriority() error {
	return syscall.Setpriority(syscall.PRIO_PROCESS, 0, backgroundNice)
}

// promoteIO has no I/O class to return outside Linux, and the nice value stays.
func promoteIO(int) error { return nil }

// ioClass answers no I/O class, because outside Linux a thread has none.
func ioClass(int) (string, error) {
	return "", errors.New("the I/O class is a Linux priority")
}
