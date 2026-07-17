//go:build darwin

package releaseevidence

import (
	"syscall"
	"unsafe"
)

func atomicExchangeDirs(left, right string) error {
	l, err := syscall.BytePtrFromString(left)
	if err != nil {
		return err
	}
	r, err := syscall.BytePtrFromString(right)
	if err != nil {
		return err
	}
	const (
		sysRenameatxNP = 488
		atFDCWD        = ^uintptr(1)
		renameSwap     = 2
	)
	_, _, errno := syscall.Syscall6(sysRenameatxNP, atFDCWD, uintptr(unsafe.Pointer(l)), atFDCWD, uintptr(unsafe.Pointer(r)), renameSwap, 0)
	if errno != 0 {
		return errno
	}
	return nil
}
