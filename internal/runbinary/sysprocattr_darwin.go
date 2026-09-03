//go:build darwin

package runbinary

import "syscall"

// builderProcAttr places the builder child in its own process group. Darwin carries no
// parent-death signal, so the drain is the whole removal path here.
func builderProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setpgid: true}
}
