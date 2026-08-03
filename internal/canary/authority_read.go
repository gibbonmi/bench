//go:build !windows

package canary

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

func readCanaryAuthority(file *os.File, limit int) ([]byte, error) {
	fd := int(file.Fd())
	if err := syscall.SetNonblock(fd, true); err != nil {
		return nil, err
	}
	data := make([]byte, 0, limit)
	buffer := make([]byte, limit+1)
	for {
		n, err := syscall.Read(fd, buffer)
		if n > 0 {
			data = append(data, buffer[:n]...)
			if len(data) > limit {
				return nil, errors.New("canary authority exceeds read limit")
			}
		}
		if err == nil {
			if n == 0 {
				return data, nil
			}
			continue
		}
		if errors.Is(err, syscall.EINTR) {
			continue
		}
		if errors.Is(err, syscall.EAGAIN) || errors.Is(err, syscall.EWOULDBLOCK) {
			return nil, errors.New("canary authority is not complete")
		}
		return nil, fmt.Errorf("read canary authority: %w", err)
	}
}
