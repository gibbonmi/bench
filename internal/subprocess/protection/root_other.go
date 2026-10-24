//go:build !linux && !darwin

package protection

import "errors"

// systemBase is the system temporary directory.
func systemBase() string { return "/tmp" }

// bootID refuses, because this system names no boot identifier that Bench reads.
func bootID() (string, error) {
	return "", errors.New("protection: this system names no boot identifier")
}
