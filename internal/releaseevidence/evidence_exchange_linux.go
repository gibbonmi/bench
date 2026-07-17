//go:build linux

package releaseevidence

import (
	"fmt"
	"runtime"
	"syscall"
	"unsafe"
)

func atomicExchangeDirs(left, right string) error {
	var trap uintptr
	switch runtime.GOARCH {
	case "amd64":
		trap = 316
	case "arm64":
		trap = 276
	default:
		return fmt.Errorf("atomic evidence replacement is unsupported on linux/%s", runtime.GOARCH)
	}
	return renameExchange(trap, ^uintptr(99), left, right)
}

func renameExchange(trap, atFDCWD uintptr, left, right string) error {
	l, err := syscall.BytePtrFromString(left)
	if err != nil {
		return err
	}
	r, err := syscall.BytePtrFromString(right)
	if err != nil {
		return err
	}
	_, _, errno := syscall.Syscall6(trap, atFDCWD, uintptr(unsafe.Pointer(l)), atFDCWD, uintptr(unsafe.Pointer(r)), 2, 0)
	if errno != 0 {
		return errno
	}
	return nil
}
