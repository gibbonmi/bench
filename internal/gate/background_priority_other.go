//go:build !linux

package gate

import "syscall"

// lowerPriority sets this process to backgroundNice. Outside Linux the nice value is an
// attribute of the process, and the platform has no I/O class that a process sets for itself.
func lowerPriority() error {
	return syscall.Setpriority(syscall.PRIO_PROCESS, 0, backgroundNice)
}
