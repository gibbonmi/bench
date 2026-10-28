//go:build linux

package gate

import "syscall"

// ioClassNames names each I/O class as ionice(1) prints it.
var ioClassNames = map[int]string{0: "none", 1: "realtime", ioprioClassBestEffort: "best-effort", ioprioClassIdle: "idle"}

// ioClass answers the name of the I/O class of the thread tid.
func ioClass(tid int) (string, error) {
	value, _, errno := syscall.Syscall(syscall.SYS_IOPRIO_GET, ioprioWhoProcess, uintptr(tid), 0)
	if errno != 0 {
		return "", errno
	}
	return ioClassNames[int(value)>>ioprioClassShift], nil
}
