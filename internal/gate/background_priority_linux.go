//go:build linux

package gate

import (
	"errors"
	"os"
	"strconv"
	"syscall"
)

// The I/O priority value of the idle class, as ioprio_set(2) encodes it for one thread.
const (
	ioprioWhoProcess = 1
	ioprioClassIdle  = 3
	ioprioClassShift = 13
)

// lowerPriority sets each thread of this process to backgroundNice and to the idle I/O class.
// On Linux both are attributes of a thread, and a new thread or a child process takes them from
// the thread that creates it. So the pass repeats until it finds no thread that it has not set.
// A thread that exits during the pass needs no priority.
func lowerPriority() error {
	set := map[int]bool{}
	for {
		entries, err := os.ReadDir("/proc/self/task")
		if err != nil {
			return err
		}
		found := false
		for _, entry := range entries {
			tid, err := strconv.Atoi(entry.Name())
			if err != nil || set[tid] {
				continue
			}
			set[tid], found = true, true
			if err := lowerThread(tid); err != nil && !errors.Is(err, syscall.ESRCH) {
				return err
			}
		}
		if !found {
			return nil
		}
	}
}

// lowerThread sets the thread tid to backgroundNice and to the idle I/O class.
func lowerThread(tid int) error {
	if err := syscall.Setpriority(syscall.PRIO_PROCESS, tid, backgroundNice); err != nil {
		return err
	}
	if _, _, errno := syscall.Syscall(syscall.SYS_IOPRIO_SET, ioprioWhoProcess, uintptr(tid), ioprioClassIdle<<ioprioClassShift); errno != 0 {
		return errno
	}
	return nil
}
