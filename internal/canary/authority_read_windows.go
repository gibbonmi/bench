package canary

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var peekNamedPipe = syscall.NewLazyDLL("kernel32.dll").NewProc("PeekNamedPipe")

func readCanaryAuthority(file *os.File, limit int) ([]byte, error) {
	handle := syscall.Handle(file.Fd())
	data := make([]byte, 0, limit)
	for {
		var available uint32
		ok, _, callErr := peekNamedPipe.Call(
			uintptr(handle),
			0,
			0,
			0,
			uintptr(unsafe.Pointer(&available)),
			0,
		)
		if ok == 0 {
			if errors.Is(callErr, syscall.ERROR_BROKEN_PIPE) {
				return data, nil
			}
			return nil, fmt.Errorf("inspect canary authority: %w", callErr)
		}
		if available == 0 {
			return nil, errors.New("canary authority is not complete")
		}
		remaining := limit + 1 - len(data)
		if remaining <= 0 {
			return nil, errors.New("canary authority exceeds read limit")
		}
		if uint32(remaining) > available {
			remaining = int(available)
		}
		buffer := make([]byte, remaining)
		n, err := syscall.Read(handle, buffer)
		if err != nil {
			return nil, fmt.Errorf("read canary authority: %w", err)
		}
		data = append(data, buffer[:n]...)
		if len(data) > limit {
			return nil, errors.New("canary authority exceeds read limit")
		}
	}
}
