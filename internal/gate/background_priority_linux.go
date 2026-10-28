//go:build linux

package gate

import (
	"errors"
	"os"
	"strconv"
	"syscall"
)

// The I/O priority values of ioprio_set(2): the targets of one thread and of one process group,
// the idle class, and the best-effort class at the level that the nice value backgroundNice gives
// a thread with no class of its own.
const (
	ioprioWhoProcess       = 1
	ioprioWhoGroup         = 2
	ioprioClassBestEffort  = 2
	ioprioClassIdle        = 3
	ioprioClassShift       = 13
	ioprioBestEffortLevel  = (backgroundNice + 20) / 5
	ioprioPromotedPriority = ioprioClassBestEffort<<ioprioClassShift | ioprioBestEffortLevel
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

// promoteIO returns each process group of the session that leader leads to the best-effort I/O
// class. A background grade leads its own session, and its oracle runs in a process group of that
// session, so the pass reaches the grade and each process that it started. The nice value stays,
// because Linux refuses an unprivileged lower nice value. A group that is gone needs no class, and
// each other refusal joins the answer.
func promoteIO(leader int) error {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return err
	}
	groups := map[int]bool{}
	for _, entry := range entries {
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}
		session, _, errno := syscall.RawSyscall(syscall.SYS_GETSID, uintptr(pid), 0, 0)
		if errno != 0 || int(session) != leader {
			continue
		}
		if group, err := syscall.Getpgid(pid); err == nil {
			groups[group] = true
		}
	}
	var refused []error
	for group := range groups {
		if _, _, errno := syscall.Syscall(syscall.SYS_IOPRIO_SET, ioprioWhoGroup, uintptr(group), ioprioPromotedPriority); errno != 0 && errno != syscall.ESRCH {
			refused = append(refused, errno)
		}
	}
	return errors.Join(refused...)
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
