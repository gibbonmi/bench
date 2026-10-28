//go:build !linux

package gate

import "errors"

// ioClass answers no I/O class, because outside Linux a thread has none.
func ioClass(int) (string, error) {
	return "", errors.New("the I/O class is a Linux priority")
}
