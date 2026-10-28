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

// lowerPriority sets each thread of this process to backgroundNice, and then each thread to the
// idle I/O class. Each step is best-effort: a refused thread does not stop its pass, and the
// refusals of both passes reach the caller as one error.
func lowerPriority() error {
	return errors.Join(eachThread(niceThread), eachThread(idleThread))
}

// eachThread calls set once for each thread of this process. On Linux a priority is an attribute
// of a thread, and a new thread or a child process takes it from the thread that creates it. So
// the pass repeats until it finds no thread that it has not set. A thread that exits during the
// pass needs no priority, and each other refusal joins the answer.
func eachThread(set func(tid int) error) error {
	done := map[int]bool{}
	var refused []error
	for {
		entries, err := os.ReadDir("/proc/self/task")
		if err != nil {
			return errors.Join(append(refused, err)...)
		}
		found := false
		for _, entry := range entries {
			tid, err := strconv.Atoi(entry.Name())
			if err != nil || done[tid] {
				continue
			}
			done[tid], found = true, true
			if err := set(tid); err != nil && !errors.Is(err, syscall.ESRCH) {
				refused = append(refused, err)
			}
		}
		if !found {
			return errors.Join(refused...)
		}
	}
}

// niceThread sets the thread tid to backgroundNice.
func niceThread(tid int) error {
	return syscall.Setpriority(syscall.PRIO_PROCESS, tid, backgroundNice)
}

// idleThread sets the thread tid to the idle I/O class.
func idleThread(tid int) error {
	if _, _, errno := syscall.Syscall(syscall.SYS_IOPRIO_SET, ioprioWhoProcess, uintptr(tid), ioprioClassIdle<<ioprioClassShift); errno != 0 {
		return errno
	}
	return nil
}
