//go:build linux

package runbinary

import "syscall"

// builderProcAttr adds a parent-death signal to the group placement. An owner killed
// mid-build never reaches its drain, so the kernel is the only remaining path that
// removes the builder child.
func builderProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setpgid: true, Pdeathsig: syscall.SIGKILL}
}
