//go:build !linux && !darwin

package runbinary

import "syscall"

// builderProcAttr places the builder child in its own process group. A platform outside
// the release matrix gets the group placement alone.
func builderProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setpgid: true}
}
